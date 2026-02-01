package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mikail-tommard/task-flow/internal/domain"
)

type RepoUser struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *RepoUser {
	return &RepoUser{
		db: db,
	}
}

func (r *RepoUser) Create(ctx context.Context, user *domain.User) (int, error) {
	const q = `
		INSERT INTO users(email, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int
	err := r.db.QueryRowContext(ctx, q, user.Email(), user.PasswordHash()).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *RepoUser) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash
		FROM users
		WHERE email = $1
	`

	var (
		id int
		emailUser string
		hash string
	)

	err := r.db.QueryRowContext(ctx, q, email).Scan(&id, &emailUser, &hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return domain.FromStorageUser(id, emailUser, hash), nil
}