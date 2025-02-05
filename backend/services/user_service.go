package services

import (
	"chat/models"
	"chat/repositories"
	"errors"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) CreateUser(user models.User) error {
	return s.Repo.CreateUser(user)
}

func (s *UserService) AuthenticateUser(username, password string) (bool, error) {
	storedPassword, err := s.Repo.GetPassword(username)
	if err != nil {
		return false, err
	}
	if storedPassword != password {
		return false, errors.New("認証失敗")
	}
	return true, nil
}
