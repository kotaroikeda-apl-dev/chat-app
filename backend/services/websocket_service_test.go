package services_test

import (
	"chat/models"
	"chat/services"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// **ダミーの WebSocket 接続を作成**
func newMockWebSocketConn(t *testing.T) *websocket.Conn {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("WebSocket upgrade error: %v", err)
		}
		_ = conn
	}))

	// クライアントとして WebSocket に接続
	wsURL := "ws" + server.URL[len("http"):] // http → ws に変換
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Logf("WebSocket connection error: %v", err)
	}

	return conn
}

func TestWebSocketService(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := services.NewWebSocketService(mockRepo)

	t.Run("AddClient and RemoveClient", func(t *testing.T) {
		mockConn := newMockWebSocketConn(t)

		// クライアント追加
		service.AddClient(mockConn)

		// **構造体にキャストして Clients へアクセス**
		wsService := service.(services.WebSocketService)
		assert.True(t, wsService.GetClients()[mockConn])

		// クライアント削除
		service.RemoveClient(mockConn)

		// 削除後の確認
		assert.False(t, wsService.GetClients()[mockConn])
	})

	t.Run("SaveMessage", func(t *testing.T) {
		msg := models.Message{
			SpaceID:   1,
			Username:  "testuser",
			Text:      "Test message",
			CreatedAt: time.Now(),
		}

		mockRepo.On("CreateMessage", mock.AnythingOfType("models.Message")).Return(1, nil)

		err := service.SaveMessage(msg)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("BroadcastMessage", func(t *testing.T) {
		mockConn := newMockWebSocketConn(t) // *websocket.Conn を返す
		service.AddClient(mockConn)

		msg := models.Message{
			SpaceID:   1,
			Username:  "testuser",
			Text:      "Broadcast Test",
			CreatedAt: time.Now(),
		}

		// WebSocket の WriteMessage の呼び出しをモック
		mockConn.WriteMessage(websocket.TextMessage, []byte("Test"))

		service.BroadcastMessage(msg)

		// 少し待機して、ブロードキャストされるのを待つ
		time.Sleep(time.Millisecond * 10)

		mockRepo.AssertExpectations(t)
	})
}
