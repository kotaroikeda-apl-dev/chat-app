package repositories

import "chat/models"

type MessageRepository interface {
	CreateMessage(msg models.Message) (int, error)
	GetMessages(spaceId int) ([]models.Message, error)
	DeleteMessage(messageID, spaceID int) error
}
