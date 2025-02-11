package controllers

import (
	"chat/models"
	"chat/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketController struct {
	Service  services.WebSocketService
	Upgrader websocket.Upgrader
}

func NewWebSocketController(service services.WebSocketService) *WebSocketController {
	return &WebSocketController{
		Service: service,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// **WebSocket接続を処理**
func (c *WebSocketController) HandleConnections(ctx *gin.Context) {
	ws, err := c.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("WebSocket接続エラー:", err)
		return
	}
	defer ws.Close()

	c.Service.AddClient(ws)

	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			c.Service.RemoveClient(ws)
			break
		}

		if msg.Username == "" {
			msg.Username = "匿名ユーザー"
		}

		err = c.Service.SaveMessage(msg)
		if err != nil {
			log.Println("メッセージ保存エラー:", err)
		}

		c.Service.BroadcastMessage(msg)
	}
}
