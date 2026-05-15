package repository

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID
	ClientID  uuid.UUID
	SellerID  uuid.UUID
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Items     []OrderItem
}

type OrderItem struct {
	ID         uuid.UUID
	OrderID    uuid.UUID
	PositionID uuid.UUID
	Quantity   int
	Price      int64
}

type NewOrderItem struct {
	PositionID uuid.UUID
	Quantity   int
	Price      int64
}
