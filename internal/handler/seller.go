package handler

import (
	"beer/internal/model"
	"beer/internal/usecase/seller"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) GetSellers(c *gin.Context) {
	sellers, err := h.sellerUC.GetSellers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, sellers)
}

func (h *Handler) GetSellerByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	sellerEntity, err := h.sellerUC.GetSellerByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sellerusecase.ErrSellerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "seller not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, sellerEntity)
}

func (h *Handler) CreateSeller(c *gin.Context) {
	var req CreateSellerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validateCreateSellerRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sellerEntity, err := h.sellerUC.CreateSeller(c.Request.Context(), req.Name, req.Login, req.PasswordHash)
	if err != nil {
		if errors.Is(err, sellerusecase.ErrLoginAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "login already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, sellerEntity)
}

func (h *Handler) PatchSellerByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateSellerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	patch := model.SellerPatch{
		Name:         req.Name,
		Login:        req.Login,
		PasswordHash: req.PasswordHash,
	}
	sellerEntity, err := h.sellerUC.PatchSellerByID(c.Request.Context(), id, patch)
	if err != nil {
		if errors.Is(err, sellerusecase.ErrSellerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "seller not found"})
			return
		}
		if errors.Is(err, sellerusecase.ErrLoginAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "login already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, sellerEntity)
}

func (h *Handler) DeleteSellerByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.sellerUC.DeleteSellerByID(c.Request.Context(), id); err != nil {
		if errors.Is(err, sellerusecase.ErrSellerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "seller not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.Status(http.StatusNoContent)
}
