package usecase

import (
	"context"

	"github.com/mikail-tommard/task-flow/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, task *domain.Task) (int, error)
	GetByID(ctx context.Context, id int) (*domain.Task, error)
	ListByUser(ctx context.Context, userID int) ([]*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
}

type RepositoryUser interface {
	Create(ctx context.Context, user *domain.User) (int, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

type Service struct {
	repo Repository
}

type AuthService struct {
	repo   RepositoryUser
	hasher PasswordHasher
}

type InputUser struct {
	Email    string
	Password string
}

type Input struct {
	Title       string
	Description string
	Done        bool
	UserID      int
}

type UpdateTaskInput struct {
	ID          int
	Title       *string
	Description *string
	Done        *bool
}
