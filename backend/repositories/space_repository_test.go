package repositories_test

import (
	"chat/repositories"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// モックDBセットアップ関数
func setupMockSpaceDB(t *testing.T) (repositories.SpaceRepository, sqlmock.Sqlmock) {
	t.Helper()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm DB with sqlmock: %v", err)
	}

	repo := repositories.NewSpaceRepository(gormDB)
	return repo, mock
}

// ---------------
// CreateSpace のテスト
// ---------------
func TestCreateSpace(t *testing.T) {
	repo, mock := setupMockSpaceDB(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "spaces" \("name"\) VALUES \(\$1\) RETURNING "created_at","id"`).
		WithArgs("Test Space").
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "id"}).
			AddRow(time.Now(), 1))
	mock.ExpectCommit()

	err := repo.CreateSpace("Test Space")
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// 失敗ケース（DBエラー）
func TestCreateSpace_DBError(t *testing.T) {
	repo, mock := setupMockSpaceDB(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "spaces" \("name"\) VALUES \(\$1\) RETURNING "created_at","id"`).
		WithArgs("Test Space").
		WillReturnError(errors.New("mock db error"))
	mock.ExpectRollback()

	err := repo.CreateSpace("Test Space")
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// ---------------
// GetSpaces のテスト
// ---------------
func TestGetSpaces(t *testing.T) {
	repo, mock := setupMockSpaceDB(t)

	now := time.Now().UTC().Truncate(time.Microsecond)

	mock.ExpectQuery(`SELECT \* FROM "spaces" ORDER BY created_at ASC`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at"}).
			AddRow(1, "Space A", now).
			AddRow(2, "Space B", now))

	spaces, err := repo.GetSpaces()
	assert.NoError(t, err)
	assert.Len(t, spaces, 2)

	assert.Equal(t, 1, spaces[0].ID)
	assert.Equal(t, "Space A", spaces[0].Name)
	assert.Equal(t, 2, spaces[1].ID)
	assert.Equal(t, "Space B", spaces[1].Name)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// GetSpaces の失敗ケース
func TestGetSpaces_DBError(t *testing.T) {
	repo, mock := setupMockSpaceDB(t)

	mock.ExpectQuery(`SELECT \* FROM "spaces" ORDER BY created_at ASC`).
		WillReturnError(errors.New("mock db error"))

	spaces, err := repo.GetSpaces()
	assert.Error(t, err)
	assert.Nil(t, spaces)
	assert.Contains(t, err.Error(), "mock db error")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
