package repositories

import (
	"chat/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// ユーザー作成
func (repo *UserRepository) CreateUser(user models.User) error {
	return repo.DB.Create(&user).Error
}

// パスワード取得
func (repo *UserRepository) GetPassword(username string) (string, error) {
	var user models.User
	err := repo.DB.Where("username = ?", username).First(&user).Error
	return user.Password, err
}
