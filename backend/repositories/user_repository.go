package repositories

import (
	"chat/models"

	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

// ユーザー作成
func (repo *userRepository) CreateUser(user models.User) error {
	return repo.DB.Create(&user).Error
}

// ユーザー取得
func (repo *userRepository) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	err := repo.DB.Where("username = ?", username).First(&user).Error
	return user, err
}

// パスワード取得
func (repo *userRepository) GetPasswordByUsername(username string) (string, error) {
	var user models.User
	err := repo.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Password, nil
}
