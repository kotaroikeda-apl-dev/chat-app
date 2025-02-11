package services

import "chat/models"

type SpaceService interface {
	CreateSpace(name string) error
	GetSpaces() ([]models.Space, error)
}
