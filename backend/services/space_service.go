package services

import (
	"chat/models"
	"chat/repositories"
)

type SpaceService struct {
	Repo *repositories.SpaceRepository
}

func NewSpaceService(repo *repositories.SpaceRepository) *SpaceService {
	return &SpaceService{Repo: repo}
}

// スペースを作成
func (s *SpaceService) CreateSpace(name string) error {
	return s.Repo.CreateSpace(name)
}

// スペース一覧を取得
func (s *SpaceService) GetSpaces() ([]models.Space, error) {
	return s.Repo.GetSpaces()
}
