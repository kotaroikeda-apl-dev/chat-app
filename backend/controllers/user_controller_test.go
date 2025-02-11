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

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterUser(user models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) AuthenticateUser(user models.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func setupRouterUser() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestUserController_RegisterUser(t *testing.T) {
	mockService := new(MockUserService)
	controller := controllers.NewUserController(mockService)
	router := setupRouterUser()
	router.POST("/register", controller.RegisterUser)

	validUser := models.User{Username: "testuser", Password: "securepass"}
	mockService.On("RegisterUser", validUser).Return(nil).Once()

	jsonData, _ := json.Marshal(validUser)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{"message": "ユーザー登録成功"}`, w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/register", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockService.On("RegisterUser", validUser).Return(errors.New("ユーザーは既に存在します")).Once()

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.JSONEq(t, `{"error": "ユーザーは既に存在します"}`, w.Body.String())

	mockService.AssertExpectations(t)
}

func TestUserController_LoginUser(t *testing.T) {
	mockService := new(MockUserService)
	controller := controllers.NewUserController(mockService)
	router := setupRouterUser()
	router.POST("/login", controller.LoginUser)

	validUser := models.User{Username: "testuser", Password: "securepass"}
	mockService.On("AuthenticateUser", validUser).Return("valid-token", nil).Once()

	jsonData, _ := json.Marshal(validUser)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"token": "valid-token"}`, w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockService.On("AuthenticateUser", validUser).Return("", errors.New("認証失敗")).Once()

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.JSONEq(t, `{"error": "認証失敗"}`, w.Body.String())

	mockService.AssertExpectations(t)
}
