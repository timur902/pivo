package model

import (
	"github.com/google/uuid"
	"time"
)

type Position struct {
	ID          uuid.UUID
	Name        string
	Description string
	ImageURL    string
	SizeLiters  float32
	Quantity    int
	Price       int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Client struct {
	ID           uuid.UUID
	Name         string
	Phone        string
	Email        string
	Login        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ClientPatch struct {
	Name         *string
	Phone        *string
	Email        *string
	Login        *string
	PasswordHash *string
}

type Seller struct {
	ID           uuid.UUID
	Name         string
	Login        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type SellerPatch struct {
	Name         *string
	Login        *string
	PasswordHash *string
}

type Order struct {
	ID        uuid.UUID
	ClientID  uuid.UUID
	SellerID  uuid.UUID
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderItem struct {
	ID         uuid.UUID
	OrderID    uuid.UUID
	PositionID uuid.UUID
	Quantity   int
	Price      int64
}

type PositionPatch struct {
	Name        *string
	Description *string
	ImageURL    *string
	SizeLiters  *float32
	Quantity    *int
	Price       *int64
}
