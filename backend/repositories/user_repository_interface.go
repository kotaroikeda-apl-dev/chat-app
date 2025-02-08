package repositories

import "chat/models"

type UserRepository interface {
	CreateUser(user models.User) error
	GetUserByUsername(username string) (models.User, error)
	GetPasswordByUsername(username string) (string, error)
}
