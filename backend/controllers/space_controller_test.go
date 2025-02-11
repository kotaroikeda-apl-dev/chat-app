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

// MockSpaceService は SpaceService のモック
type MockSpaceService struct {
	mock.Mock
}

func (m *MockSpaceService) CreateSpace(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

// 修正: `[]models.Space` を返すように変更
func (m *MockSpaceService) GetSpaces() ([]models.Space, error) {
	args := m.Called()
	return args.Get(0).([]models.Space), args.Error(1)
}

// テスト用のルーターをセットアップ
func setupRouterSpace() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestSpaceController_CreateSpace(t *testing.T) {
	mockService := new(MockSpaceService)
	controller := controllers.NewSpaceController(mockService)
	router := setupRouterSpace()
	router.POST("/spaces", controller.CreateSpace)

	mockService.On("CreateSpace", "NewSpace").Return(nil).Once()

	newSpace := map[string]string{"name": "NewSpace"}
	jsonData, _ := json.Marshal(newSpace)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/spaces", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{"message":"スペースが作成されました"}`, w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/spaces", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"リクエストのパースに失敗しました"}`, w.Body.String())

	mockService.On("CreateSpace", "ErrorSpace").Return(errors.New("DBエラー")).Once()

	errorSpace := map[string]string{"name": "ErrorSpace"}
	jsonData, _ = json.Marshal(errorSpace)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/spaces", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error":"スペースの作成に失敗しました"}`, w.Body.String())

	mockService.AssertExpectations(t)
}

func TestSpaceController_GetSpaces(t *testing.T) {
	mockService := new(MockSpaceService)
	controller := controllers.NewSpaceController(mockService)
	router := setupRouterSpace()
	router.GET("/spaces", controller.GetSpaces)

	mockSpaces := []models.Space{
		{ID: 1, Name: "Space1"},
		{ID: 2, Name: "Space2"},
	}

	mockService.On("GetSpaces").Return(mockSpaces, nil).Once()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/spaces", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedJSON, _ := json.Marshal(mockSpaces)
	assert.JSONEq(t, string(expectedJSON), w.Body.String())

	mockService.On("GetSpaces").Return([]models.Space{}, errors.New("DBエラー")).Once()

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/spaces", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error":"スペース一覧の取得に失敗しました"}`, w.Body.String())

	mockService.AssertExpectations(t)
}
