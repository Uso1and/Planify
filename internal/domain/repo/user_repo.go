package repo

import (
	"context"
	"database/sql"

	"planify/internal/domain/models"

)

type UserRepoInterface interface {
	CreateUser(c context.Context, user *models.Users) error
	GetUser(c context.Context, userID int) (*models.Users, error)
	GetUserByUsername(c context.Context, username string) (*models.Users, error)
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(c context.Context, user *models.Users) error {

	query := `INSERT INTO users (username, password, email, created_at) VALUES ($1, $2, $3, $4) RETURNING id`

	return r.db.QueryRowContext(c, query, user.Username, user.Password, user.Email, user.CreatedAt).Scan(&user.ID)

}

func (r *UserRepo) GetUser(c context.Context, userID int) (*models.Users, error) {
	user := &models.Users{ID: userID}

	query := `SELECT username, password, email, created_at FROM users WHERE id = $1`

	err := r.db.QueryRowContext(c, query, userID).Scan(
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return user, err
}

func (r *UserRepo) GetUserByUsername(c context.Context, username string) (*models.Users, error) {

	user := &models.Users{}

	query := `SELECT id, username, password, email, created_at FROM users WHERE username = $1`

	err := r.db.QueryRowContext(c, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil

}
