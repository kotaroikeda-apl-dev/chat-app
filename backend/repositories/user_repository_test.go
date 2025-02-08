package repositories_test

import (
	"chat/models"
	"chat/repositories"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// モックDBセットアップ関数
func setupMockUserDB(t *testing.T) (repositories.UserRepository, sqlmock.Sqlmock) {
	t.Helper()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // SQL ログを表示
	})
	if err != nil {
		t.Fatalf("Failed to open gorm DB with sqlmock: %v", err)
	}

	repo := repositories.NewUserRepository(gormDB)
	return repo, mock
}

func TestCreateUser(t *testing.T) {
	repo, mock := setupMockUserDB(t)

	user := models.User{
		Username: "testuser",
		Password: "securepassword",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users" ("username","password") VALUES ($1,$2)`)).
		WithArgs(user.Username, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.CreateUser(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByUsername(t *testing.T) {
	repo, mock := setupMockUserDB(t)

	user := models.User{
		Username: "testuser",
		Password: "securepassword",
	}

	rows := sqlmock.NewRows([]string{"username", "password"}).
		AddRow(user.Username, user.Password)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."username" LIMIT $2`)).
		WithArgs(user.Username, 1). // `LIMIT 1` を `WithArgs` に明示
		WillReturnRows(rows)

	result, err := repo.GetUserByUsername("testuser")
	assert.NoError(t, err)
	assert.Equal(t, user.Username, result.Username)
	assert.Equal(t, user.Password, result.Password)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPasswordByUsername_Success(t *testing.T) {
	repo, mock := setupMockUserDB(t)

	user := models.User{
		Username: "testuser",
		Password: "hashedpassword123",
	}

	rows := sqlmock.NewRows([]string{"username", "password"}).
		AddRow(user.Username, user.Password)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."username" LIMIT $2`)).
		WithArgs(user.Username, 1).
		WillReturnRows(rows)

	password, err := repo.GetPasswordByUsername("testuser")
	assert.NoError(t, err)
	assert.Equal(t, user.Password, password)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPasswordByUsername_NotFound(t *testing.T) {
	repo, mock := setupMockUserDB(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."username" LIMIT $2`)).
		WithArgs("unknownuser", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	password, err := repo.GetPasswordByUsername("unknownuser")
	assert.Error(t, err)
	assert.Equal(t, "", password)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
