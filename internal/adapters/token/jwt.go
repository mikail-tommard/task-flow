package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrSecretRequired = errors.New("secret is required")
)

type Claims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type Config struct {
	Secret    []byte
	AccessTTL time.Duration
	Issuer    string
}

type Service struct {
	secret    []byte
	accessTTL time.Duration
	issuer    string
}

func NewServiceJWT(cfg Config) (*Service, error) {
	if cfg.Secret == nil {
		return nil, ErrSecretRequired
	}

	if cfg.AccessTTL == 0 {
		cfg.AccessTTL = 15 * time.Minute
	}

	if cfg.Issuer == "" {
		return nil, errors.New("issuer is required")
	}

	return &Service{
		secret:    cfg.Secret,
		accessTTL: cfg.AccessTTL,
		issuer:    cfg.Issuer,
	}, nil
}

func (s *Service) GenerateToken(id int, email string) (string, error) {
	if id == 0 {
		return "", errors.New("id is required")
	}

	if email == "" {
		return "", errors.New("email is required")
	}
	now := time.Now()
	claims := Claims{
		ID: id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: s.issuer,
			Subject: email,
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}
