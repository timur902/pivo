package handler

import (
	"beer/internal/model"
	"beer/internal/money"
	"beer/internal/storage"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreatePositionRequest struct {
	Name string `json:"name"`
	Description string `json:"description"`
	ImageURL string `json:"image_url"`
	SizeLiters float32 `json:"size_liters"`
	Quantity int `json:"quantity"`
	Price int64 `json:"price"`
}

type UpdatePositionRequest struct {
	Name *string `json:"name"`
	Description *string `json:"description"`
	ImageURL *string `json:"image_url"`
	SizeLiters *float32 `json:"size_liters"`
	Quantity *int `json:"quantity"`
	Price *int64 `json:"price"`
}

func GetPositions(c *gin.Context) {
	positions := storage.GetPositions()
	c.JSON(http.StatusOK, positions)
}

func CreatePosition(c *gin.Context) {
	var req CreatePositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if req.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price must be greater than zero"})
		return
	}
	now := time.Now()
	position := model.Position{
		ID: uuid.New(),
		Name: req.Name,
		Description: req.Description,
		ImageURL: req.ImageURL,
		SizeLiters: req.SizeLiters,
		Quantity: req.Quantity,
		Price: money.New(req.Price),
		CreatedAt: now,
		UpdatedAt: now,
	}
	storage.AddPosition(position)
	c.JSON(http.StatusCreated, position)
}

func GetPositionByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	positions := storage.GetPositions()
	for _, position := range positions {
		if position.ID == id {
			c.JSON(http.StatusOK, position)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "position not found"})
}

func DeletePositionByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	deleted := storage.DeletePositionByID(id)
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "position not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func PatchPositionByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req UpdatePositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	position, ok := storage.PatchPositionByID(
		id,
		req.Name,
		req.Description,
		req.ImageURL,
		req.SizeLiters,
		req.Quantity,
		req.Price,
	)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "position not found"})
		return
	}
	c.JSON(http.StatusOK, position)
}