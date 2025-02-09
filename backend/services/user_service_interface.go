package services

import "chat/models"

type UserService interface {
	RegisterUser(user models.User) error
	AuthenticateUser(user models.User) (string, error)
}
