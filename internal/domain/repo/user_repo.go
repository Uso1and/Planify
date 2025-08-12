package repo

import (
	"context"
	"database/sql"
	"planify/internal/domain/models"
)

type UserRepoInterface interface {
	CreateUser(c context.Context, user *models.Users) error
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(c context.Context, user *models.Users) error {

	query := "INSERT INTO users (username, password, email, created_at) VALUES ($1, $2, $3, $4) RETURNING id"

	return r.db.QueryRowContext(c, query, user.Username, user.Password, user.Email, user.CreatedAt).Scan(&user.ID)

}
