package repo

import (
	"context"
	"database/sql"

	"planify/internal/domain/models"
)

type NoteRepoInterface interface {
	CreateNewNote(c context.Context, note *models.Note) error
	GetNotesByUserID(c context.Context, userID int) ([]models.Note, error)
	GetNoteByID(c context.Context, id int, userID int) (*models.Note, error)
	UpdateNote(c context.Context, note *models.Note) error
	DeleteNote(c context.Context, noteID, userID int) error
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

func (r *NoteRepo) GetNotesByUserID(c context.Context, userID int) ([]models.Note, error) {
	query := `SELECT id, category, title, content, created_at, updated_at 
              FROM notes WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(c, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note

	for rows.Next() {
		var note models.Note
		note.UserID = userID

		if err := rows.Scan(&note.ID, &note.Category, &note.Title, &note.Content,
			&note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *NoteRepo) GetNoteByID(c context.Context, id int, userID int) (*models.Note, error) {

	query := `SELECT id, user_id, category, title, content, created_at, updated_at FROM notes WHERE id = $1 AND user_id = $2`

	note := &models.Note{}

	err := r.db.QueryRowContext(c, query, id, userID).Scan(
		&note.ID,
		&note.UserID,
		&note.Category,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
		&note.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return note, nil

}

func (r *NoteRepo) UpdateNote(c context.Context, note *models.Note) error {
	query := `UPDATE notes SET category=$1, title=$2, content=$3, updated_at=NOW() WHERE id=$4 AND user_id=$5`
	_, err := r.db.ExecContext(c, query, note.Category, note.Title, note.Content, note.ID, note.UserID)
	return err
}

func (r *NoteRepo) DeleteNote(c context.Context, noteID, userID int) error {
	query := `DELETE FROM notes WHERE id=$1 AND user_id=$2`

	_, err := r.db.ExecContext(c, query, noteID, userID)
	return err
}
