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

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateTask(ctx context.Context, input Input) (*domain.Task, error) {
	t, err := domain.New(input.Title, input.Description, input.UserID)
	if err != nil {
		return nil, err
	}

	id, err := s.repo.Create(ctx, t)
	if err != nil {
		return nil, err
	}

	return domain.FromStorage(id, input.Title, input.Done, input.Description, input.UserID), nil
}

func (s *Service) GetTask(ctx context.Context, id int) (*domain.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListTasks(ctx context.Context, userID int) ([]*domain.Task, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *Service) CompleteTask(ctx context.Context, id int) error {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	t.Complete()

	return nil
}
