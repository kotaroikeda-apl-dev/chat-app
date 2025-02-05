package services

import (
	"chat/models"
	"chat/repositories"
	"errors"
)

type MessageService struct {
	Repo *repositories.MessageRepository
}

func NewMessageService(repo *repositories.MessageRepository) *MessageService {
	return &MessageService{Repo: repo}
}

func (s *MessageService) GetMessages(spaceId int) ([]models.Message, error) {
	return s.Repo.GetMessages(spaceId)
}

// メッセージ登録
func (s *MessageService) CreateMessage(msg models.Message) (int, error) {
	// 入力値のバリデーション
	if msg.Text == "" || msg.Username == "" || msg.SpaceID == "0" {
		return 0, errors.New("メッセージまたはユーザー名が空です")
	}

	// メッセージ保存
	return s.Repo.CreateMessage(msg)
}

func (s *MessageService) DeleteMessage(messageID, spaceID int) error {
	return s.Repo.DeleteMessage(messageID, spaceID)
}
