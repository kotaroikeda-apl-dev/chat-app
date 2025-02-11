package services

import (
	"chat/models"
	"chat/repositories"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type userService struct {
	Repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{Repo: repo}
}

// ユーザー登録
func (s *userService) RegisterUser(user models.User) error {
	// すでにユーザーが存在するか確認
	existingUser, err := s.Repo.GetUserByUsername(user.Username)
	if err == nil && existingUser.Username != "" {
		return errors.New("ユーザー名が既に使用されています")
	}

	// 新規ユーザーを登録
	return s.Repo.CreateUser(user)
}

func (s *userService) AuthenticateUser(user models.User) (string, error) {
	storedPassword, err := s.Repo.GetPasswordByUsername(user.Username)
	if err != nil {
		return "", errors.New("認証失敗")
	}

	if storedPassword != user.Password {
		return "", errors.New("認証失敗")
	}

	// トークンの作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", errors.New("トークン生成エラー")
	}

	return tokenString, nil
}
