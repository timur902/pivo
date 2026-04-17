package handler

import (
	"beer/internal/model"
	"beer/internal/repository/position"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) GetPositions(c *gin.Context) {
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
	positions, err := h.positionRepo.GetPositions(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, positions)
}

func (h *Handler) CreatePosition(c *gin.Context) {
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
	positionEntity := model.Position{
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
	if err := h.positionRepo.AddPosition(c.Request.Context(), positionEntity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, positionEntity)
}

func (h *Handler) GetPositionByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	positionEntity, err := h.positionRepo.GetPositionByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, position.ErrPositionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "position not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, positionEntity)
}

func (h *Handler) DeletePositionByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	deleted, err := h.positionRepo.DeletePositionByID(c.Request.Context(), id)
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

func (h *Handler) PatchPositionByID(c *gin.Context) {
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
	patch := model.PositionPatch{
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		SizeLiters:  req.SizeLiters,
		Quantity:    req.Quantity,
		Price:       req.Price,
	}
	positionEntity, err := h.positionRepo.PatchPositionByID(c.Request.Context(), id, patch)
	if err != nil {
		if errors.Is(err, position.ErrPositionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "position not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, positionEntity)
}
