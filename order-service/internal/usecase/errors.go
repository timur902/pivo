package usecase

import "errors"

var (
	ErrNoItems           = errors.New("at least one item is required")
	ErrInvalidQuantity   = errors.New("quantity must be greater than zero")
	ErrInvalidPrice      = errors.New("price must be greater than zero")
	ErrInvalidTargetStat = errors.New("target status must be ready or cancelled")
)
