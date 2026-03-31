package handler

import (
	"beer/internal/model"
	"beer/internal/storage"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func GetClients(c *gin.Context) {
	clients := storage.GetClients()
	c.JSON(http.StatusOK, clients)
}

func GetClientByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	client, ok := storage.GetClientByID(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
		return
	}
	c.JSON(http.StatusOK, client)
}

func CreateClient(c *gin.Context) {
	var req CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if req.Name == "" || req.Phone == "" || req.Email == "" || req.Login == "" || req.PasswordHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, phone, email, login and password_hash are required"})
		return
	}
	now := time.Now()
	client := model.Client{
		ID:           uuid.New(),
		Name:         req.Name,
		Phone:        req.Phone,
		Email:        req.Email,
		Login:        req.Login,
		PasswordHash: req.PasswordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	storage.AddClient(client)
	c.JSON(http.StatusCreated, client)
}

func PatchClientByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	client, ok := storage.PatchClientByID(
		id,
		req.Name,
		req.Phone,
		req.Email,
		req.Login,
		req.PasswordHash,
	)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
		return
	}
	c.JSON(http.StatusOK, client)
}

func DeleteClientByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	deleted := storage.DeleteClientByID(id)
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
