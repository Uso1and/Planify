package repo

import (
	"context"
	"database/sql"

	"planify/internal/domain/models"
)

type NoteRepoInterface interface {
	CreateNewNote(c context.Context, note *models.Note) error
}

type NoteRepo struct {
	db *sql.DB
}

func NewNoteRepo(db *sql.DB) *NoteRepo {
	return &NoteRepo{db: db}
}

func (r *NoteRepo) CreateNewNote(c context.Context, note *models.Note) error {
	query := `INSERT INTO notes (user_id, category, title, content, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(c, query, note.UserID, note.Category, note.Title, note.Content).Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt)

}
