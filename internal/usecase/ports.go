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

type Input struct {
	Title       string
	Description string
	Done        bool
	UserID      int
}

type Service struct {
	repo Repository
}