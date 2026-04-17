package handler

import (
	"beer/internal/model"
	"beer/internal/repository/client"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (h *Handler) GetClients(c *gin.Context) {
	clients, err := h.clientRepo.GetClients(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, clients)
}

func (h *Handler) GetClientByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	clientEntity, err := h.clientRepo.GetClientByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, client.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, clientEntity)
}

func (h *Handler) CreateClient(c *gin.Context) {
	var req CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := validateCreateClientRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()
	clientEntity := model.Client{
		ID:           uuid.New(),
		Name:         req.Name,
		Phone:        req.Phone,
		Email:        req.Email,
		Login:        req.Login,
		PasswordHash: req.PasswordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := h.clientRepo.AddClient(c.Request.Context(), clientEntity); err != nil {
		if errors.Is(err, client.ErrLoginAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "login already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, clientEntity)
}

func (h *Handler) PatchClientByID(c *gin.Context) {
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
	patch := model.ClientPatch{
		Name:         req.Name,
		Phone:        req.Phone,
		Email:        req.Email,
		Login:        req.Login,
		PasswordHash: req.PasswordHash,
	}
	clientEntity, err := h.clientRepo.PatchClientByID(c.Request.Context(), id, patch)
	if err != nil {
		if errors.Is(err, client.ErrClientNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
			return
		}
		if errors.Is(err, client.ErrLoginAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "login already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, clientEntity)
}

func (h *Handler) DeleteClientByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	deleted, err := h.clientRepo.DeleteClientByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
