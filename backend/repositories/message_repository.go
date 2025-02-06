package repositories

import (
	"chat/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type MessageRepository struct {
	DB *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{DB: db}
}

// メッセージを作成し、新しいIDを返す
func (repo *MessageRepository) CreateMessage(msg models.Message) (int, error) {
	result := repo.DB.Create(&msg)
	if result.Error != nil {
		fmt.Println("DBエラー:", result.Error)
		return 0, result.Error
	}
	return msg.ID, nil
}

// 指定された spaceId のメッセージ一覧を取得
func (repo *MessageRepository) GetMessages(spaceId int) ([]models.Message, error) {
	var messages []models.Message
	err := repo.DB.Where("space_id = ?", spaceId).Order("created_at ASC").Find(&messages).Error
	return messages, err
}

func (r *MessageRepository) DeleteMessage(messageID, spaceID int) error {
	result := r.DB.Delete(&models.Message{}, "id = ? AND space_id = ?", messageID, spaceID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("メッセージが見つかりませんでした")
	}

	return nil
}
