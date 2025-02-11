package services

import (
	"chat/models"
	"chat/repositories"
)

type spaceService struct {
	Repo repositories.SpaceRepository
}

func NewSpaceService(repo repositories.SpaceRepository) SpaceService {
	return &spaceService{Repo: repo}
}

// スペースを作成
func (s *spaceService) CreateSpace(name string) error {
	return s.Repo.CreateSpace(name)
}

// スペース一覧を取得
func (s *spaceService) GetSpaces() ([]models.Space, error) {
	return s.Repo.GetSpaces()
}
