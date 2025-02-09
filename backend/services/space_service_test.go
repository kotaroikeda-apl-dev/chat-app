package services_test

import (
	"chat/models"
	"chat/services"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSpaceRepository は SpaceRepository のモック
type MockSpaceRepository struct {
	mock.Mock
}

func (m *MockSpaceRepository) CreateSpace(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockSpaceRepository) GetSpaces() ([]models.Space, error) {
	args := m.Called()
	return args.Get(0).([]models.Space), args.Error(1)
}

func TestCreateSpace_Success(t *testing.T) {
	mockRepo := new(MockSpaceRepository)
	service := services.NewSpaceService(mockRepo)

	mockRepo.On("CreateSpace", "Test Space").Return(nil)

	err := service.CreateSpace("Test Space")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateSpace_Error(t *testing.T) {
	mockRepo := new(MockSpaceRepository)
	service := services.NewSpaceService(mockRepo)

	mockRepo.On("CreateSpace", "Test Space").Return(errors.New("DB error"))

	err := service.CreateSpace("Test Space")

	assert.Error(t, err)
	assert.Equal(t, "DB error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetSpaces_Success(t *testing.T) {
	mockRepo := new(MockSpaceRepository)
	service := services.NewSpaceService(mockRepo)

	spaces := []models.Space{
		{ID: 1, Name: "Space 1"},
		{ID: 2, Name: "Space 2"},
	}

	mockRepo.On("GetSpaces").Return(spaces, nil)

	result, err := service.GetSpaces()

	assert.NoError(t, err)
	assert.Equal(t, spaces, result)
	mockRepo.AssertExpectations(t)
}

func TestGetSpaces_Error(t *testing.T) {
	mockRepo := new(MockSpaceRepository)
	service := services.NewSpaceService(mockRepo)

	mockRepo.On("GetSpaces").Return([]models.Space{}, errors.New("DB error"))

	result, err := service.GetSpaces()

	assert.Error(t, err)
	assert.Equal(t, "DB error", err.Error())
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}
