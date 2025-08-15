package repo

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"planify/internal/domain/models"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepo(db)

	now := time.Now()
	user := &models.Users{
		Username:  "testuser",
		Password:  "hashedpassword",
		Email:     "test@example.com",
		CreatedAt: now,
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO users`).
			WithArgs(user.Username, user.Password, user.Email, user.CreatedAt).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		err := repo.CreateUser(context.Background(), user)

		assert.NoError(t, err)
		assert.Equal(t, 1, user.ID)
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO users`).
			WithArgs(user.Username, user.Password, user.Email, user.CreatedAt).
			WillReturnError(errors.New("database error"))

		err := repo.CreateUser(context.Background(), user)

		assert.Error(t, err)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepo(db)

	userID := 1
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		expectedUser := &models.Users{
			ID:        userID,
			Username:  "testuser",
			Password:  "hashedpassword",
			Email:     "test@example.com",
			CreatedAt: now,
		}

		mock.ExpectQuery(`SELECT username, password, email, created_at FROM users WHERE id = \$1`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"username", "password", "email", "created_at"}).
				AddRow(expectedUser.Username, expectedUser.Password, expectedUser.Email, expectedUser.CreatedAt))

		user, err := repo.GetUser(context.Background(), userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT username, password, email, created_at FROM users WHERE id = \$1`).
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUser(context.Background(), userID)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.True(t, errors.Is(err, sql.ErrNoRows))
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT username, password, email, created_at FROM users WHERE id = \$1`).
			WithArgs(userID).
			WillReturnError(errors.New("database error"))

		user, err := repo.GetUser(context.Background(), userID)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepo(db)

	username := "testuser"
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		expectedUser := &models.Users{
			ID:        1,
			Username:  username,
			Password:  "hashedpassword",
			Email:     "test@example.com",
			CreatedAt: now,
		}

		mock.ExpectQuery(`SELECT id, username, password, email, created_at FROM users WHERE username = \$1`).
			WithArgs(username).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "email", "created_at"}).
				AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Password, expectedUser.Email, expectedUser.CreatedAt))

		user, err := repo.GetUserByUsername(context.Background(), username)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, username, password, email, created_at FROM users WHERE username = \$1`).
			WithArgs(username).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByUsername(context.Background(), username)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.True(t, errors.Is(err, sql.ErrNoRows))
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, username, password, email, created_at FROM users WHERE username = \$1`).
			WithArgs(username).
			WillReturnError(errors.New("database error"))

		user, err := repo.GetUserByUsername(context.Background(), username)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}
