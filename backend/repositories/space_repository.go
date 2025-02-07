package repositories

import (
	"chat/models"

	"gorm.io/gorm"
)

type spaceRepository struct {
	DB *gorm.DB
}

func NewSpaceRepository(db *gorm.DB) SpaceRepository {
	return &spaceRepository{DB: db}
}

// スペースを作成
func (repo *spaceRepository) CreateSpace(name string) error {
	space := models.Space{Name: name}
	return repo.DB.Create(&space).Error
}

// スペース一覧を取得
func (repo *spaceRepository) GetSpaces() ([]models.Space, error) {
	var spaces []models.Space
	err := repo.DB.Order("created_at ASC").Find(&spaces).Error
	return spaces, err
}
