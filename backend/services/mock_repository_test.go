package services_test

import (
	"chat/models"

	"github.com/stretchr/testify/mock"
)

// **MockMessageRepository（共通）**
type MockMessageRepository struct {
	mock.Mock
}

// `CreateMessage` の戻り値を `(int, error)` に統一
func (m *MockMessageRepository) CreateMessage(msg models.Message) (int, error) {
	args := m.Called(msg)
	return args.Int(0), args.Error(1)
}

// `DeleteMessage` を追加（必要な場合）
func (m *MockMessageRepository) DeleteMessage(messageID int, spaceID int) error {
	args := m.Called(messageID, spaceID)
	return args.Error(0)
}
