package usecase

import (
	"context"
	"errors"

	"github.com/mikail-tommard/task-flow/internal/domain"
)

func NewAuthService(repo RepositoryUser, hasher PasswordHasher) *AuthService {
	return &AuthService{
		repo:   repo,
		hasher: hasher,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, input InputUser) (*domain.User, error) {
	hash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	u, err := domain.NewUser(input.Email, hash)
	if err != nil {
		return nil, err
	}

	id, err := s.repo.Create(ctx, u)
	if err != nil {
		return nil, err
	}

	return domain.FromStorageUser(id, input.Email, hash), nil
}

func (s *AuthService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := s.hasher.Compare(u.PasswordHash(), password); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := s.tokens.GenerateToken(u.UserID(), u.Email())
	if err != nil {
		return "", err
	}
	return token, nil
}