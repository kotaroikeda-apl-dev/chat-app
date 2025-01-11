package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

var db *sql.DB
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Message struct {
	ID        int    `json:"id"`
	SpaceID   int    `json:"space_id"`
	Username  string `json:"username"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

func main() {
	initDB()
	defer db.Close()

	setupServer()
}

func initDB() {
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("データベース接続エラー: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("データベースに接続できません: %v", err)
	}

	fmt.Println("データベースに正常に接続しました")
}

func setupServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthCheckHandler)
	mux.HandleFunc("/api/register", registerUser)
	mux.HandleFunc("/api/login", loginUser)
	mux.HandleFunc("/api/messages", getMessages)
	mux.HandleFunc("/api/delete", deleteMessage)
	mux.HandleFunc("/api/ws", handleConnections)
	mux.HandleFunc("/api/spaces", createSpace)
	mux.HandleFunc("/api/spaces/list", getSpaces)
	mux.HandleFunc("/api/messages/create", createMessage)

	go handleMessages()

	log.Println("サーバーを起動中: ポート:8080")
	err := http.ListenAndServe(":8080", enableCORS(mux))
	if err != nil {
		log.Fatal("サーバー起動エラー:", err)
	}
}

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
		log.Printf("ユーザー登録エラー: %v", err)
		http.Error(w, "ユーザー登録に失敗しました", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("ユーザー登録成功"))
}

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

func getMessages(w http.ResponseWriter, r *http.Request) {
	// クエリパラメータから spaceId を取得
	spaceId := r.URL.Query().Get("spaceId")
	if spaceId == "" {
		http.Error(w, "spaceId が指定されていません", http.StatusBadRequest)
		return
	}

	// SQL クエリを修正して WHERE 句を追加
	rows, err := db.Query("SELECT id, space_id, username, text, created_at FROM messages WHERE space_id = $1 ORDER BY created_at ASC", spaceId)
	if err != nil {
		http.Error(w, "メッセージの取得に失敗しました", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.ID, &msg.SpaceID, &msg.Username, &msg.Text, &msg.CreatedAt)
		if err != nil {
			http.Error(w, "メッセージのパースに失敗しました", http.StatusInternalServerError)
			return
		}
		messages = append(messages, msg)
	}

	// メッセージを JSON 形式でレスポンス
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func deleteMessage(w http.ResponseWriter, r *http.Request) {
	messageID := r.URL.Query().Get("id")
	spaceID := r.URL.Query().Get("spaceId")

	log.Printf("受信した削除リクエスト - メッセージID: %s, スペースID: %s", messageID, spaceID)

	if messageID == "" || spaceID == "" {
		http.Error(w, "メッセージIDまたはスペースIDが指定されていません", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM messages WHERE id = $1 AND space_id = $2", messageID, spaceID)
	if err != nil {
		log.Printf("メッセージ削除エラー: %v", err)
		http.Error(w, "メッセージの削除に失敗しました", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("メッセージ削除成功"))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket接続エラー:", err)
		return
	}
	defer ws.Close()

	clients[ws] = true
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
			break
		}

		if msg.Username == "" {
			msg.Username = "匿名ユーザー"
		}

		_, err = db.Exec("INSERT INTO messages (username, text, created_at) VALUES ($1, $2, NOW())", msg.Username, msg.Text)
		if err != nil {
			log.Println("メッセージ保存エラー:", err)
		}

		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

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
			log.Printf("JSON デコードエラー: %v", err)
			http.Error(w, "スペースデータのパースに失敗しました", http.StatusInternalServerError)
			return
		}
		spaces = append(spaces, space)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spaces)
}

func createMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "無効なリクエストメソッドです", http.StatusMethodNotAllowed)
		return
	}

	var msg struct {
		ID       int    `json:"id"`
		SpaceID  int    `json:"space_id"`
		Username string `json:"username"`
		Text     string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Printf("JSONパースエラー: %v", err)
		http.Error(w, "リクエストのパースに失敗しました", http.StatusBadRequest)
		return
	}

	if msg.Text == "" || msg.Username == "" || msg.SpaceID == 0 {
		log.Printf("不正なリクエスト: %+v", msg)
		http.Error(w, "メッセージまたはユーザー名が空です", http.StatusBadRequest)
		return
	}

	var newID int
	err := db.QueryRow(
		"INSERT INTO messages (space_id, username, text, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id",
		msg.SpaceID, msg.Username, msg.Text,
	).Scan(&newID)
	if err != nil {
		log.Printf("メッセージ保存エラー: %v", err)
		http.Error(w, "メッセージの保存に失敗しました", http.StatusInternalServerError)
		return
	}

	msg.ID = newID // メッセージに新しいIDを追加
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// レスポンスとしてステータス200とメッセージを返す
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}
