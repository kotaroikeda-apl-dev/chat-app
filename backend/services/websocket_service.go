package services

import (
	"chat/models"
	"chat/repositories"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketService struct {
	Repo      *repositories.MessageRepository
	Clients   map[*websocket.Conn]bool
	Broadcast chan models.Message
	Mutex     sync.Mutex
	Upgrader  websocket.Upgrader
}

func NewWebSocketService(repo *repositories.MessageRepository) *WebSocketService {
	return &WebSocketService{
		Repo:      repo,
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan models.Message),
	}
}

// クライアントを追加
func (s *WebSocketService) AddClient(ws *websocket.Conn) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Clients[ws] = true
}

// クライアントを削除
func (s *WebSocketService) RemoveClient(ws *websocket.Conn) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	delete(s.Clients, ws)
}

// **メッセージをDBに保存**
func (s *WebSocketService) SaveMessage(msg models.Message) error {
	_, err := s.Repo.CreateMessage(msg)
	return err
}

// **メッセージを全クライアントにブロードキャスト**
func (s *WebSocketService) BroadcastMessage(msg models.Message) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	for client := range s.Clients {
		err := client.WriteJSON(msg)
		if err != nil {
			client.Close()
			delete(s.Clients, client)
		}
	}
}

func (s *WebSocketService) HandleMessages() {
	for {
		// メッセージをチャネルから受け取る
		msg := <-s.Broadcast

		// 全クライアントにメッセージを送信（ブロードキャスト）
		for client := range s.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(s.Clients, client) // 接続が切れたクライアントを削除
			}
		}
	}
}
