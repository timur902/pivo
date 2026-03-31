package storage

import (
	"beer/internal/model"
	"beer/internal/money"
	"time"
	"github.com/google/uuid"
)

var positions = []model.Position{
	{
		ID:          uuid.New(),
		Name:        "Heineken",
		Description: "Светлое пиво",
		ImageURL:    "https://example.com/heineken.jpg",
		SizeLiters:  0.5,
		Quantity:    20,
		Price:       money.New(15000),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		ID:          uuid.New(),
		Name:        "Guinness",
		Description: "Темное пиво",
		ImageURL:    "https://example.com/guinness.jpg",
		SizeLiters:  0.44,
		Quantity:    15,
		Price:       money.New(18000),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
}

func GetPositions() []model.Position {
	return positions
}

func AddPosition(position model.Position) {
	positions = append(positions, position)
}

func DeletePositionByID(id uuid.UUID) bool {
	for i, position := range positions {
		if position.ID == id {
			positions = append(positions[:i], positions[i+1:]...)
			return true
		}
	}
	return false
}

func PatchPositionByID(
	id uuid.UUID,
	name *string,
	description *string,
	imageURL *string,
	sizeLiters *float32,
	quantity *int,
	price *int64,
) (model.Position, bool) {
	for i := range positions {
		if positions[i].ID == id {
			if name != nil {
				positions[i].Name = *name
			}
			if description != nil {
				positions[i].Description = *description
			}
			if imageURL != nil {
				positions[i].ImageURL = *imageURL
			}
			if sizeLiters != nil {
				positions[i].SizeLiters = *sizeLiters
			}
			if quantity != nil {
				positions[i].Quantity = *quantity
			}
			if price != nil {
				positions[i].Price = money.New(*price)
			}
			positions[i].UpdatedAt = time.Now()

			return positions[i], true
		}
	}
	return model.Position{}, false
}