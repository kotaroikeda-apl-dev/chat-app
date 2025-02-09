package services

import (
	"chat/models"
	"chat/repositories"
	"sync"

	"github.com/gorilla/websocket"
)

type webSocketService struct {
	Repo      repositories.MessageRepository
	Clients   map[*websocket.Conn]bool
	Broadcast chan models.Message
	Mutex     sync.Mutex
	Upgrader  websocket.Upgrader
}

func NewWebSocketService(repo repositories.MessageRepository) WebSocketService {
	return &webSocketService{
		Repo:      repo,
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan models.Message),
	}
}

// クライアントを追加
func (s *webSocketService) AddClient(ws *websocket.Conn) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Clients[ws] = true
}

// クライアントを削除
func (s *webSocketService) RemoveClient(ws *websocket.Conn) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	delete(s.Clients, ws)
}

// **メッセージをDBに保存**
func (s *webSocketService) SaveMessage(msg models.Message) error {
	_, err := s.Repo.CreateMessage(msg)
	return err
}

// **メッセージを全クライアントにブロードキャスト**
func (s *webSocketService) BroadcastMessage(msg models.Message) {
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

func (s *webSocketService) HandleMessages() {
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

func (s *webSocketService) GetClients() map[*websocket.Conn]bool {
	return s.Clients
}
