package security

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	Cost int
}

func NewBcryptHasher(cost int) BcryptHasher {
	return BcryptHasher{
		Cost: cost,
	}
}

func (h BcryptHasher) Hash(password string) (string, error) {
	if password == "" {
		return "", errors.New("password empty")
	}

	cost := h.Cost
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}

	b, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (h BcryptHasher) Compare(passwordHash, password string) error {
	if passwordHash == "" || password == "" {
		return errors.New("empty hash or password")
	}
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}