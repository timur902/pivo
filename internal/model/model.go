package model

import (
	"beer/internal/money"
	"time"
	"github.com/google/uuid"
)

type Position struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	ImageURL    string      `json:"image_url"`
	SizeLiters  float32     `json:"size_liters"`
	Quantity    int         `json:"quantity"`
	Price       money.Money `json:"price"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type Client struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Seller struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Order struct {
	ID        uuid.UUID `json:"id"`
	ClientID  uuid.UUID `json:"client_id"`
	SellerID  uuid.UUID `json:"seller_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderItem struct {
	ID         uuid.UUID   `json:"id"`
	OrderID    uuid.UUID   `json:"order_id"`
	PositionID uuid.UUID   `json:"position_id"`
	Quantity   int         `json:"quantity"`
	Price      money.Money `json:"price"`
}