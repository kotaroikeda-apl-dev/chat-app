package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"os"
	"fmt"


	_ "github.com/lib/pq" // PostgreSQLドライバをインポート
	"github.com/gorilla/websocket"
	"github.com/golang-jwt/jwt/v4"
)

var db *sql.DB
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // CORS対応
	},
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Message struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}

var clients = make(map[*websocket.Conn]bool) // 接続中のクライアント
var broadcast = make(chan Message)          // ブロードキャスト用チャネル

func main() {
	// 環境変数を取得
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")

	// 接続文字列を作成
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	// データベースに接続
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("データベース接続エラー: %v", err)
	}
	defer db.Close()

	// 接続確認
	err = db.Ping()
	if err != nil {
		log.Fatalf("データベースに接続できません: %v", err)
	}

	fmt.Println("データベースに正常に接続しました")

	// ルーティング
	mux := http.NewServeMux()
	mux.HandleFunc("/register", registerUser)
	mux.HandleFunc("/login", loginUser)
	mux.HandleFunc("/messages", getMessages)
	mux.HandleFunc("/delete", deleteMessage)
	mux.HandleFunc("/ws", handleConnections)
	mux.HandleFunc("/spaces", createSpace)
	mux.HandleFunc("/spaces/list", getSpaces)

	// メッセージブロードキャスト処理をゴルーチンで開始
	go handleMessages()

	// サーバー起動
	log.Println("サーバーを起動中: http://localhost:8080")
	err = http.ListenAndServe(":8080", enableCORS(mux))
	if err != nil {
		log.Fatal("サーバー起動エラー:", err)
	}
}

// CORS対応
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ユーザー登録
func registerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "無効なリクエストメソッドです", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "リクエストのパースに失敗しました", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	if err != nil {
		http.Error(w, "ユーザー登録に失敗しました", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("ユーザー登録成功"))
}

// ユーザーログイン
func loginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "無効なリクエストメソッドです", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "リクエストのパースに失敗しました", http.StatusBadRequest)
		return
	}

	var storedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE username = $1", user.Username).Scan(&storedPassword)
	if err != nil || storedPassword != user.Password {
		http.Error(w, "認証失敗", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		http.Error(w, "トークン生成エラー", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// メッセージ取得
func getMessages(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, username, text, created_at FROM messages ORDER BY created_at ASC")
	if err != nil {
		http.Error(w, "メッセージの取得に失敗しました", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.ID, &msg.Username, &msg.Text, &msg.CreatedAt)
		if err != nil {
			http.Error(w, "メッセージのパースに失敗しました", http.StatusInternalServerError)
			return
		}
		messages = append(messages, msg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// メッセージ削除
func deleteMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
			http.Error(w, "無効なリクエストメソッドです", http.StatusMethodNotAllowed)
			return
	}

	// クエリパラメータからメッセージIDを取得
	messageID := r.URL.Query().Get("id")
	if messageID == "" {
			http.Error(w, "メッセージIDが指定されていません", http.StatusBadRequest)
			return
	}

	// データベースからメッセージを削除
	_, err := db.Exec("DELETE FROM messages WHERE id = $1", messageID)
	if err != nil {
			log.Println("メッセージ削除エラー:", err)
			http.Error(w, "メッセージの削除に失敗しました", http.StatusInternalServerError)
			return
	}

	log.Println("メッセージ削除成功: ID =", messageID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("メッセージが削除されました"))
}

// WebSocket接続
func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
			log.Println("WebSocket接続エラー:", err)
			return
	}
	defer ws.Close()

	clients[ws] = true
	log.Println("クライアント接続:", ws.RemoteAddr())

	for {
			var msg Message
			err := ws.ReadJSON(&msg)
			if err != nil {
					log.Println("クライアント切断:", err)
					delete(clients, ws)
					break
			}

			// デフォルト値設定
			if msg.Username == "" {
					msg.Username = "匿名ユーザー"
			}

			// データベースに保存
			_, err = db.Exec(
					"INSERT INTO messages (username, text, created_at) VALUES ($1, $2, NOW())",
					msg.Username, msg.Text,
			)
			if err != nil {
					log.Println("メッセージ保存エラー:", err)
			} else {
					log.Println("メッセージ保存成功")
			}

			broadcast <- msg
	}
}

// メッセージブロードキャスト
func handleMessages() {
	for {
		msg := <-broadcast
		log.Println("ブロードキャストするメッセージ:", msg) // ブロードキャストログ

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("クライアントへの送信エラー: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

// スペース作成
func createSpace(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
			http.Error(w, "無効なリクエストメソッドです", http.StatusMethodNotAllowed)
			return
	}

	var space struct {
			Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&space); err != nil {
			http.Error(w, "リクエストのパースに失敗しました", http.StatusBadRequest)
			return
	}

	_, err := db.Exec("INSERT INTO spaces (name) VALUES ($1)", space.Name)
	if err != nil {
			http.Error(w, "スペースの作成に失敗しました", http.StatusInternalServerError)
			return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("スペースが作成されました"))
}

// スペース一覧取得
func getSpaces(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, created_at FROM spaces ORDER BY created_at ASC")
	if err != nil {
			http.Error(w, "スペース一覧の取得に失敗しました", http.StatusInternalServerError)
			return
	}
	defer rows.Close()

	var spaces []struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			CreatedAt string `json:"created_at"`
	}
	for rows.Next() {
			var space struct {
					ID        int    `json:"id"`
					Name      string `json:"name"`
					CreatedAt string `json:"created_at"`
			}
			if err := rows.Scan(&space.ID, &space.Name, &space.CreatedAt); err != nil {
					http.Error(w, "スペースデータのパースに失敗しました", http.StatusInternalServerError)
					return
			}
			spaces = append(spaces, space)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spaces)
}