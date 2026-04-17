package handler

import (
	"beer/internal/model"
	positionrepo "beer/internal/repository/position"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

func GetPositions(c *gin.Context) {
	limit := 20
	offset := 0
	if rawLimit := c.Query("limit"); rawLimit != "" {
		parsedLimit, err := strconv.Atoi(rawLimit)
		if err != nil || parsedLimit <= 0 || parsedLimit > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be an integer between 1 and 100"})
			return
		}
		limit = parsedLimit
	}
	if rawOffset := c.Query("offset"); rawOffset != "" {
		parsedOffset, err := strconv.Atoi(rawOffset)
		if err != nil || parsedOffset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "offset must be an integer greater than or equal to 0"})
			return
		}
		offset = parsedOffset
	}
	positions, err := positionrepo.GetPositions(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, positions)
}

func CreatePosition(c *gin.Context) {
	var req CreatePositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := validateCreatePositionRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()
	position := model.Position{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		SizeLiters:  req.SizeLiters,
		Quantity:    req.Quantity,
		Price:       req.Price,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := positionrepo.AddPosition(c.Request.Context(), position); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, position)
}

func GetPositionByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	position, ok, err := positionrepo.GetPositionByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if ok {
		c.JSON(http.StatusOK, position)
		return
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
	deleted, err := positionrepo.DeletePositionByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
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
	position, ok, err := positionrepo.PatchPositionByID(
		c.Request.Context(),
		id,
		req.Name,
		req.Description,
		req.ImageURL,
		req.SizeLiters,
		req.Quantity,
		req.Price,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "position not found"})
		return
	}
	c.JSON(http.StatusOK, position)
}
