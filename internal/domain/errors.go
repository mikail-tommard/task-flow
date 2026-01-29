package domain

import "errors"

var (
	ErrInvalidTitle     = errors.New("invalid title")
	ErrAlreadyCompleted = errors.New("task already completed")
	ErrInvalidEmail     = errors.New("invalid email")
)
