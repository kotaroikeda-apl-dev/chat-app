package repositories

import "chat/models"

type SpaceRepository interface {
	CreateSpace(name string) error
	GetSpaces() ([]models.Space, error)
}
