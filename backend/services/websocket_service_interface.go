package services

import (
	"chat/models"

	"github.com/gorilla/websocket"
)

type WebSocketService interface {
	AddClient(ws *websocket.Conn)
	RemoveClient(ws *websocket.Conn)
	SaveMessage(msg models.Message) error
	BroadcastMessage(msg models.Message)
	GetClients() map[*websocket.Conn]bool
	HandleMessages()
}
