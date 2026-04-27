package handler

import (
	"beer/internal/repository/client"
	"beer/internal/repository/position"
	"beer/internal/usecase/seller"
)

type Handler struct {
	clientRepo   *client.Repository
	positionRepo *position.Repository
	sellerUC     *sellerusecase.Usecase
}

func NewHandler(clientRepo *client.Repository, positionRepo *position.Repository, sellerUC *sellerusecase.Usecase) *Handler {
	return &Handler{
		clientRepo:   clientRepo,
		positionRepo: positionRepo,
		sellerUC:     sellerUC,
	}
}
