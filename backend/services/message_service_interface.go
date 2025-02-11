package services

import "chat/models"

type MessageService interface {
	GetMessages(spaceId int) ([]models.Message, error)
	CreateMessage(msg models.Message) (int, error)
	DeleteMessage(messageID, spaceID int) error
}
