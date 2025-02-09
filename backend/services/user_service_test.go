package services_test

import (
	"chat/models"
	"chat/services"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository は UserRepository のモック
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByUsername(username string) (models.User, error) {
	args := m.Called(username)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) GetPasswordByUsername(username string) (string, error) {
	args := m.Called(username)
	return args.String(0), args.Error(1)
}

func TestRegisterUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	newUser := models.User{Username: "testuser", Password: "securepassword"}

	mockRepo.On("GetUserByUsername", "testuser").Return(models.User{}, errors.New("not found"))
	mockRepo.On("CreateUser", newUser).Return(nil)

	err := service.RegisterUser(newUser)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_AlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	existingUser := models.User{Username: "testuser", Password: "oldpassword"}

	mockRepo.On("GetUserByUsername", "testuser").Return(existingUser, nil)

	err := service.RegisterUser(existingUser)

	assert.Error(t, err)
	assert.Equal(t, "ユーザー名が既に使用されています", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthenticateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	user := models.User{Username: "testuser", Password: "securepassword"}

	mockRepo.On("GetPasswordByUsername", "testuser").Return("securepassword", nil)

	token, err := service.AuthenticateUser(user)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// トークンを検証
	parsedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})

	assert.NotNil(t, parsedToken)
	mockRepo.AssertExpectations(t)
}

func TestAuthenticateUser_WrongPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	user := models.User{Username: "testuser", Password: "wrongpassword"}

	mockRepo.On("GetPasswordByUsername", "testuser").Return("securepassword", nil)

	token, err := service.AuthenticateUser(user)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "認証失敗", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthenticateUser_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	mockRepo.On("GetPasswordByUsername", "unknownuser").Return("", errors.New("not found"))

	token, err := service.AuthenticateUser(models.User{Username: "unknownuser", Password: "password"})

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "認証失敗", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthenticateUser_EmptyPasswordInDB(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	user := models.User{Username: "testuser", Password: "securepassword"}

	mockRepo.On("GetPasswordByUsername", "testuser").Return("", nil)

	token, err := service.AuthenticateUser(user)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "認証失敗", err.Error())
	mockRepo.AssertExpectations(t)
}
