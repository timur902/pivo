package handler

import (
	"errors"
	"strings"
)

type CreateClientRequest struct {
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Login        string `json:"login"`
	PasswordHash string `json:"password_hash"`
}

type UpdateClientRequest struct {
	Name         *string `json:"name"`
	Phone        *string `json:"phone"`
	Email        *string `json:"email"`
	Login        *string `json:"login"`
	PasswordHash *string `json:"password_hash"`
}

type CreateSellerRequest struct {
	Name         string `json:"name"`
	Login        string `json:"login"`
	PasswordHash string `json:"password_hash"`
}

type UpdateSellerRequest struct {
	Name         *string `json:"name"`
	Login        *string `json:"login"`
	PasswordHash *string `json:"password_hash"`
}

type CreatePositionRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	SizeLiters  float32 `json:"size_liters"`
	Quantity    int     `json:"quantity"`
	Price       int64   `json:"price"`
}

type UpdatePositionRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	ImageURL    *string  `json:"image_url"`
	SizeLiters  *float32 `json:"size_liters"`
	Quantity    *int     `json:"quantity"`
	Price       *int64   `json:"price"`
}

func validateCreateClientRequest(req CreateClientRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(req.Phone) == "" {
		return errors.New("phone is required")
	}
	if strings.TrimSpace(req.Email) == "" {
		return errors.New("email is required")
	}
	if strings.TrimSpace(req.Login) == "" {
		return errors.New("login is required")
	}
	if strings.TrimSpace(req.PasswordHash) == "" {
		return errors.New("password_hash is required")
	}
	if len(req.Phone) > 20 {
		return errors.New("phone must be at most 20 characters")
	}
	return nil
}

func validateCreateSellerRequest(req CreateSellerRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(req.Login) == "" {
		return errors.New("login is required")
	}
	if strings.TrimSpace(req.PasswordHash) == "" {
		return errors.New("password_hash is required")
	}
	return nil
}

func validateCreatePositionRequest(req CreatePositionRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name is required")
	}
	if req.Price <= 0 {
		return errors.New("price must be greater than zero")
	}
	if req.SizeLiters <= 0 {
		return errors.New("size_liters must be greater than zero")
	}
	if req.Quantity < 0 {
		return errors.New("quantity must be greater than or equal to zero")
	}
	return nil
}
