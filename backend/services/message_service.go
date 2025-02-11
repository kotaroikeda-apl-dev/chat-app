package services

import (
	"chat/models"
	"chat/repositories"
	"errors"
)

type messageService struct {
	repo repositories.MessageRepository
}

func NewMessageService(repo repositories.MessageRepository) MessageService {
	return &messageService{repo: repo}
}

func (s *messageService) GetMessages(spaceId int) ([]models.Message, error) {
	return s.repo.GetMessages(spaceId)
}

// メッセージ登録
func (s *messageService) CreateMessage(msg models.Message) (int, error) {
	// 入力値のバリデーション
	if msg.Text == "" || msg.Username == "" || msg.SpaceID == 0 {
		return 0, errors.New("メッセージまたはユーザー名が空です")
	}

	// メッセージ保存
	return s.repo.CreateMessage(msg)
}

func (s *messageService) DeleteMessage(messageID, spaceID int) error {
	if messageID == 0 || spaceID == 0 {
		return errors.New("メッセージIDまたはスペースIDが無効です")
	}
	return s.repo.DeleteMessage(messageID, spaceID)
}
