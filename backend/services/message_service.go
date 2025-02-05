package services

import (
	"chat/models"
	"chat/repositories"
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

func (s *MessageService) CreateMessage(msg models.Message) (int, error) {
	return s.Repo.CreateMessage(msg)
}

func (s *MessageService) DeleteMessage(messageID, spaceID int) error {
	return s.Repo.DeleteMessage(messageID, spaceID)
}
