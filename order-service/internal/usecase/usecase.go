package usecase

import (
	"context"
	"beer/order-service/internal/repository"
	"github.com/google/uuid"
)

type Usecase struct {
	repo *repository.Repository
}

func NewUsecase(repo *repository.Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) CreateOrder(ctx context.Context, clientID, sellerID uuid.UUID, items []repository.NewOrderItem) (*repository.Order, error) {
	if len(items) == 0 {
		return nil, ErrNoItems
	}
	for _, it := range items {
		if it.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}
		if it.Price <= 0 {
			return nil, ErrInvalidPrice
		}
	}
	return u.repo.CreateOrder(ctx, clientID, sellerID, items)
}

func (u *Usecase) GetOrderForClient(ctx context.Context, orderID, clientID uuid.UUID) (*repository.Order, error) {
	order, err := u.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order.ClientID != clientID {
		return nil, repository.ErrOrderNotOwnedByActor
	}
	return order, nil
}

func (u *Usecase) GetOrderForSeller(ctx context.Context, orderID, sellerID uuid.UUID) (*repository.Order, error) {
	order, err := u.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order.SellerID != sellerID {
		return nil, repository.ErrOrderNotOwnedByActor
	}
	return order, nil
}

func (u *Usecase) ListOrdersByClient(ctx context.Context, clientID uuid.UUID, limit, offset int) ([]repository.Order, error) {
	return u.repo.ListOrdersByClient(ctx, clientID, limit, offset)
}

func (u *Usecase) ListOrdersBySeller(ctx context.Context, sellerID uuid.UUID, limit, offset int) ([]repository.Order, error) {
	return u.repo.ListOrdersBySeller(ctx, sellerID, limit, offset)
}

func (u *Usecase) CancelOrderByClient(ctx context.Context, orderID, clientID uuid.UUID) (*repository.Order, error) {
	return u.repo.UpdateStatus(ctx, orderID, "new", "cancelled", "client_id", clientID)
}

func (u *Usecase) UpdateStatusBySeller(ctx context.Context, orderID, sellerID uuid.UUID, targetStatus string) (*repository.Order, error) {
	if targetStatus != "ready" && targetStatus != "cancelled" {
		return nil, ErrInvalidTargetStat
	}
	return u.repo.UpdateStatus(ctx, orderID, "new", targetStatus, "seller_id", sellerID)
}
