package controllers

import (
	"chat/models"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock WebSocket Service
type MockWebSocketService struct {
	mock.Mock
}

func (m *MockWebSocketService) AddClient(conn *websocket.Conn) {
	m.Called(conn)
}

func (m *MockWebSocketService) RemoveClient(conn *websocket.Conn) {
	m.Called(conn)
}

func (m *MockWebSocketService) SaveMessage(msg models.Message) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockWebSocketService) BroadcastMessage(msg models.Message) {
	m.Called(msg)
}

func (m *MockWebSocketService) GetClients() map[*websocket.Conn]bool {
	args := m.Called()
	return args.Get(0).(map[*websocket.Conn]bool)
}

func (m *MockWebSocketService) HandleMessages() {
	m.Called()
}

// NewWebSocketController のユニットテスト
func TestNewWebSocketController(t *testing.T) {
	mockService := new(MockWebSocketService)
	controller := NewWebSocketController(mockService)

	assert.NotNil(t, controller, "WebSocketController の生成に失敗")
	assert.NotNil(t, controller.Service, "Service が nil")
}
