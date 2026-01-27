package usecase

import (
	"context"

	"github.com/mikail-tommard/task-flow/internal/domain"
)

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

func (s *Service) UpdateTask(ctx context.Context, input UpdateTaskInput) (*domain.Task, error) {
	t, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		if err := t.Rename(*input.Title); err != nil {
			return nil, err
		}
	}
	if input.Description != nil {
		t.ChangeDescription(*input.Description)
	}
	if input.Done != nil && *input.Done {
		t.Complete()
	}

	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}
