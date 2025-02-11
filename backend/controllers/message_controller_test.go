package controllers_test

import (
	"bytes"
	"chat/controllers"
	"chat/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMessageService は MessageService をモックする
type MockMessageService struct {
	mock.Mock
}

func (m *MockMessageService) GetMessages(spaceID int) ([]models.Message, error) {
	args := m.Called(spaceID)
	return args.Get(0).([]models.Message), args.Error(1)
}

func (m *MockMessageService) CreateMessage(msg models.Message) (int, error) {
	args := m.Called(msg)
	return args.Int(0), args.Error(1)
}

func (m *MockMessageService) DeleteMessage(messageID, spaceID int) error {
	args := m.Called(messageID, spaceID)
	return args.Error(0)
}

func setupRouterMessage() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}

func TestMessageController_GetMessages(t *testing.T) {
	mockService := new(MockMessageService)
	controller := controllers.NewMessageController(mockService)
	router := setupRouterMessage()
	router.GET("/messages", controller.GetMessages)

	// 正常系: メッセージ取得成功
	mockMessages := []models.Message{
		{ID: 1, SpaceID: 1, Username: "user1", Text: "Hello"},
	}
	mockService.On("GetMessages", 1).Return(mockMessages, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/messages?spaceId=1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)

	// 異常系: spaceId が数値でない
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/messages?spaceId=abc", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 異常系: Service からエラーが返る
	mockService.On("GetMessages", 2).Return([]models.Message{}, errors.New("DBエラー"))

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/messages?spaceId=2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMessageController_CreateMessage(t *testing.T) {
	mockService := new(MockMessageService)
	controller := controllers.NewMessageController(mockService)
	router := setupRouterMessage()
	router.POST("/messages", controller.CreateMessage)

	newMessage := models.Message{SpaceID: 1, Username: "user1", Text: "Hello"}
	mockService.On("CreateMessage", mock.AnythingOfType("models.Message")).Return(1, nil)

	jsonData, _ := json.Marshal(newMessage)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/messages", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/messages", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMessageController_DeleteMessage(t *testing.T) {
	mockService := new(MockMessageService)
	controller := controllers.NewMessageController(mockService)
	router := setupRouterMessage()
	router.DELETE("/messages", controller.DeleteMessage)

	// 正常系: メッセージ削除成功
	mockService.On("DeleteMessage", 1, 1).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/messages?id=1&spaceId=1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)

	// 異常系: ID が数値でない
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/messages?id=abc&spaceId=1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 異常系: Service からエラーが返る
	mockService.On("DeleteMessage", 2, 2).Return(errors.New("DBエラー"))

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/messages?id=2&spaceId=2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
