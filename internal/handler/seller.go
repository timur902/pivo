package handler

import (
	"beer/internal/model"
	sellerrepo "beer/internal/repository/seller"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func GetSellers(c *gin.Context) {
	sellers, err := sellerrepo.GetSellers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, sellers)
}

func GetSellerByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	seller, ok, err := sellerrepo.GetSellerByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "seller not found"})
		return
	}

	c.JSON(http.StatusOK, seller)
}

func CreateSeller(c *gin.Context) {
	var req CreateSellerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validateCreateSellerRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	seller := model.Seller{
		ID:           uuid.New(),
		Name:         req.Name,
		Login:        req.Login,
		PasswordHash: req.PasswordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := sellerrepo.AddSeller(c.Request.Context(), seller); err != nil {
		if errors.Is(err, sellerrepo.ErrLoginAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "login already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, seller)
}

func PatchSellerByID(c *gin.Context) {
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

	seller, ok, err := sellerrepo.PatchSellerByID(
		c.Request.Context(),
		id,
		req.Name,
		req.Login,
		req.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, sellerrepo.ErrLoginAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "login already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "seller not found"})
		return
	}

	c.JSON(http.StatusOK, seller)
}

func DeleteSellerByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	deleted, err := sellerrepo.DeleteSellerByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "seller not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
