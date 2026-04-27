package handler

import (
	"beer/internal/repository/client"
	"beer/internal/repository/position"
	"beer/internal/usecase/seller"
	"beer/proto/order"
)

type Handler struct {
	clientRepo   *client.Repository
	positionRepo *position.Repository
	sellerUC     *sellerusecase.Usecase
	orderClient  orderpb.OrderServiceClient
}

func NewHandler(clientRepo *client.Repository, positionRepo *position.Repository, sellerUC *sellerusecase.Usecase, orderClient orderpb.OrderServiceClient) *Handler {
	return &Handler{
		clientRepo:   clientRepo,
		positionRepo: positionRepo,
		sellerUC:     sellerUC,
		orderClient:  orderClient,
	}
}
