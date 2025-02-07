package repositories_test

import (
	"errors"
	"testing"
	"time"

	"chat/models"
	"chat/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// テスト用のモックDBセットアップ関数
func setupMockMessageDB(t *testing.T) (repositories.MessageRepository, sqlmock.Sqlmock) {
	t.Helper()

	// sqlmock を初期化
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}

	// GORM をモックDBで初期化
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm DB with sqlmock: %v", err)
	}

	// テスト対象のリポジトリ生成
	repo := repositories.NewMessageRepository(gormDB)
	return repo, mock
}

// CreateMessage のテスト
func TestCreateMessage(t *testing.T) {
	repo, mock := setupMockMessageDB(t)

	now := time.Now().UTC().Truncate(time.Microsecond)
	msg := models.Message{
		SpaceID:   1,
		Username:  "alice",
		Text:      "Hello, World!",
		CreatedAt: now,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "messages" \("space_id","username","text","created_at"\) VALUES \(\$1,\$2,\$3,\$4\) RETURNING "created_at","id"`).
		WithArgs(msg.SpaceID, msg.Username, msg.Text, msg.CreatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(10, now))
	mock.ExpectCommit()

	id, err := repo.CreateMessage(msg)
	assert.NoError(t, err)
	assert.Equal(t, 10, id, "返却される ID が正しいこと")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCreateMessage_DBError(t *testing.T) {
	repo, mock := setupMockMessageDB(t)

	now := time.Now().UTC().Truncate(time.Microsecond)
	msg := models.Message{
		SpaceID:   1,
		Username:  "alice",
		Text:      "Hello, World!",
		CreatedAt: now,
	}

	// DBがエラーを返すケースをモック
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "messages" \("space_id","username","text","created_at"\) VALUES \(\$1,\$2,\$3,\$4\) RETURNING "created_at","id"`).
		WithArgs(msg.SpaceID, msg.Username, msg.Text, msg.CreatedAt).
		WillReturnError(errors.New("mock db error"))
	mock.ExpectRollback()

	id, err := repo.CreateMessage(msg)
	// ここでエラーが返るはず
	assert.Error(t, err)
	assert.Equal(t, 0, id, "IDは0が返る")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// ---------------
// GetMessages のテスト
// ---------------
func TestGetMessages(t *testing.T) {
	repo, mock := setupMockMessageDB(t)

	spaceID := 1

	mock.ExpectQuery(`SELECT \* FROM "messages" WHERE space_id = \$1 ORDER BY created_at ASC`).
		WithArgs(spaceID).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "space_id", "username", "text", "created_at"},
		).
			AddRow(1, 1, "bob", "First message", time.Now()).
			AddRow(2, 1, "charlie", "Second message", time.Now()),
		)

	messages, err := repo.GetMessages(spaceID)
	assert.NoError(t, err)
	assert.Len(t, messages, 2)

	assert.Equal(t, 1, messages[0].ID)
	assert.Equal(t, "bob", messages[0].Username)
	assert.Equal(t, "First message", messages[0].Text)

	assert.Equal(t, 2, messages[1].ID)
	assert.Equal(t, "charlie", messages[1].Username)
	assert.Equal(t, "Second message", messages[1].Text)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// DeleteMessage のテスト
func TestDeleteMessage(t *testing.T) {
	repo, mock := setupMockMessageDB(t)

	messageID := 999
	spaceID := 1

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "messages" WHERE id = \$1 AND space_id = \$2`).
		WithArgs(messageID, spaceID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.DeleteMessage(messageID, spaceID)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// DeleteMessage: 該当レコードなしパターン
func TestDeleteMessage_NotFound(t *testing.T) {
	repo, mock := setupMockMessageDB(t)

	messageID := 12345
	spaceID := 999

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "messages" WHERE id = \$1 AND space_id = \$2`).
		WithArgs(messageID, spaceID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // RowsAffected=0
	mock.ExpectCommit()

	err := repo.DeleteMessage(messageID, spaceID)
	assert.Error(t, err, "該当メッセージが存在しない場合はエラー")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeleteMessage_DBError(t *testing.T) {
	repo, mock := setupMockMessageDB(t)

	messageID := 999
	spaceID := 1

	mock.ExpectBegin()
	// DB エラーをモック
	mock.ExpectExec(`DELETE FROM "messages" WHERE id = \$1 AND space_id = \$2`).
		WithArgs(messageID, spaceID).
		WillReturnError(errors.New("mock delete error"))
	mock.ExpectRollback()

	err := repo.DeleteMessage(messageID, spaceID)
	// ここでエラーが返ることを期待
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mock delete error")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
