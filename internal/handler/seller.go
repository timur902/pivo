package handler

import (
	"beer/internal/model"
	"beer/internal/storage"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

func GetSellers(c *gin.Context) {
	sellers := storage.GetSellers()
	c.JSON(http.StatusOK, sellers)
}

func GetSellerByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	seller, ok := storage.GetSellerByID(id)
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

	if req.Name == "" || req.Login == "" || req.PasswordHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, login and password_hash are required"})
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

	storage.AddSeller(seller)
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

	seller, ok := storage.PatchSellerByID(
		id,
		req.Name,
		req.Login,
		req.PasswordHash,
	)
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

	deleted := storage.DeleteSellerByID(id)
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "seller not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
