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

func TestCreateNewNote(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewNoteRepo(db)

	now := time.Now()
	note := &models.Note{
		UserID:   1,
		Category: "work",
		Title:    "test title",
		Content:  "test content",
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO notes`).
			WithArgs(note.UserID, note.Category, note.Title, note.Content).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(1, now, now))

		err := repo.CreateNewNote(context.Background(), note)

		assert.NoError(t, err)
		assert.Equal(t, 1, note.ID)
		assert.Equal(t, now, note.CreatedAt)
		assert.Equal(t, now, note.UpdatedAt)
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO notes`).
			WithArgs(note.UserID, note.Category, note.Title, note.Content).
			WillReturnError(errors.New("some error"))

		err := repo.CreateNewNote(context.Background(), note)

		assert.Error(t, err)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNotesByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewNoteRepo(db)

	now := time.Now()
	userID := 1

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "category", "title", "content", "created_at", "updated_at"}).
			AddRow(1, "work", "title 1", "content 1", now, now).
			AddRow(2, "personal", "title 2", "content 2", now, now)

		mock.ExpectQuery(`SELECT id, category, title, content, created_at, updated_at FROM notes`).
			WithArgs(userID).
			WillReturnRows(rows)

		notes, err := repo.GetNotesByUserID(context.Background(), userID)

		assert.NoError(t, err)
		assert.Len(t, notes, 2)
		assert.Equal(t, 1, notes[0].ID)
		assert.Equal(t, "work", notes[0].Category)
		assert.Equal(t, 2, notes[1].ID)
		assert.Equal(t, "personal", notes[1].Category)
	})

	t.Run("empty", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "category", "title", "content", "created_at", "updated_at"})

		mock.ExpectQuery(`SELECT id, category, title, content, created_at, updated_at FROM notes`).
			WithArgs(userID).
			WillReturnRows(rows)

		notes, err := repo.GetNotesByUserID(context.Background(), userID)

		assert.NoError(t, err)
		assert.Empty(t, notes)
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, category, title, content, created_at, updated_at FROM notes`).
			WithArgs(userID).
			WillReturnError(errors.New("some error"))

		notes, err := repo.GetNotesByUserID(context.Background(), userID)

		assert.Error(t, err)
		assert.Nil(t, notes)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNoteByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewNoteRepo(db)

	now := time.Now()
	noteID := 1
	userID := 1

	t.Run("success", func(t *testing.T) {
		row := sqlmock.NewRows([]string{"id", "user_id", "category", "title", "content", "created_at", "updated_at"}).
			AddRow(noteID, userID, "work", "title", "content", now, now)

		mock.ExpectQuery(`SELECT id, user_id, category, title, content, created_at, updated_at FROM notes`).
			WithArgs(noteID, userID).
			WillReturnRows(row)

		note, err := repo.GetNoteByID(context.Background(), noteID, userID)

		assert.NoError(t, err)
		assert.Equal(t, noteID, note.ID)
		assert.Equal(t, userID, note.UserID)
		assert.Equal(t, "work", note.Category)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, user_id, category, title, content, created_at, updated_at FROM notes`).
			WithArgs(noteID, userID).
			WillReturnError(sql.ErrNoRows)

		note, err := repo.GetNoteByID(context.Background(), noteID, userID)

		assert.Error(t, err)
		assert.Nil(t, note)
		assert.True(t, errors.Is(err, sql.ErrNoRows))
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, user_id, category, title, content, created_at, updated_at FROM notes`).
			WithArgs(noteID, userID).
			WillReturnError(errors.New("some error"))

		note, err := repo.GetNoteByID(context.Background(), noteID, userID)

		assert.Error(t, err)
		assert.Nil(t, note)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestUpdateNote(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewNoteRepo(db)

	note := &models.Note{
		ID:       1,
		UserID:   1,
		Category: "work",
		Title:    "updated title",
		Content:  "updated content",
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(`UPDATE notes`).
			WithArgs(note.Category, note.Title, note.Content, note.ID, note.UserID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateNote(context.Background(), note)

		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec(`UPDATE notes`).
			WithArgs(note.Category, note.Title, note.Content, note.ID, note.UserID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateNote(context.Background(), note)

		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectExec(`UPDATE notes`).
			WithArgs(note.Category, note.Title, note.Content, note.ID, note.UserID).
			WillReturnError(errors.New("some error"))

		err := repo.UpdateNote(context.Background(), note)

		assert.Error(t, err)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteNote(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewNoteRepo(db)

	noteID := 1
	userID := 1

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM notes`).
			WithArgs(noteID, userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteNote(context.Background(), noteID, userID)

		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM notes`).
			WithArgs(noteID, userID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.DeleteNote(context.Background(), noteID, userID)

		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM notes`).
			WithArgs(noteID, userID).
			WillReturnError(errors.New("some error"))

		err := repo.DeleteNote(context.Background(), noteID, userID)

		assert.Error(t, err)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}
