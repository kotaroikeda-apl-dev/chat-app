package services_test

import (
	"chat/models"
	"chat/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (m *MockMessageRepository) GetMessages(spaceId int) ([]models.Message, error) {
	args := m.Called(spaceId)
	return args.Get(0).([]models.Message), args.Error(1)
}

func TestGetMessages(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := services.NewMessageService(mockRepo)

	expectedMessages := []models.Message{
		{ID: 1, SpaceID: 1, Username: "alice", Text: "Hello"},
		{ID: 2, SpaceID: 1, Username: "bob", Text: "Hi"},
	}

	mockRepo.On("GetMessages", 1).Return(expectedMessages, nil)

	messages, err := service.GetMessages(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, messages)

	mockRepo.AssertExpectations(t)
}

func TestCreateMessage_Success(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := services.NewMessageService(mockRepo)

	message := models.Message{SpaceID: 1, Username: "alice", Text: "Hello"}
	mockRepo.On("CreateMessage", message).Return(1, nil)

	id, err := service.CreateMessage(message)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)

	mockRepo.AssertExpectations(t)
}

func TestCreateMessage_ValidationError(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := services.NewMessageService(mockRepo)

	message := models.Message{SpaceID: 1, Username: "", Text: "Hello"}

	id, err := service.CreateMessage(message)
	assert.Error(t, err)
	assert.Equal(t, "メッセージまたはユーザー名が空です", err.Error())
	assert.Equal(t, 0, id)
}

func TestDeleteMessage_Success(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := services.NewMessageService(mockRepo)

	mockRepo.On("DeleteMessage", 1, 1).Return(nil)

	err := service.DeleteMessage(1, 1)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteMessage_ValidationError(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	service := services.NewMessageService(mockRepo)

	err := service.DeleteMessage(0, 1)
	assert.Error(t, err)
	assert.Equal(t, "メッセージIDまたはスペースIDが無効です", err.Error())
}
