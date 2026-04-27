package repository

import "errors"

var (
	ErrOrderNotFound        = errors.New("order not found")
	ErrClientNotFound       = errors.New("client not found")
	ErrSellerNotFound       = errors.New("seller not found")
	ErrPositionNotFound     = errors.New("position not found")
	ErrInvalidStatusChange  = errors.New("invalid status transition")
	ErrOrderNotOwnedByActor = errors.New("order does not belong to actor")
)
