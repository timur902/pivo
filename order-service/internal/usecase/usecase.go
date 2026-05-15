package usecase

import (
	"context"
	"log"
	"time"

	"beer/order-service/internal/events"
	"beer/order-service/internal/repository"

	"github.com/google/uuid"
)

type EventPublisher interface {
	PublishOrderReady(ctx context.Context, evt events.OrderReadyEvent) error
}

type Usecase struct {
	repo      *repository.Repository
	publisher EventPublisher
}

func NewUsecase(repo *repository.Repository, publisher EventPublisher) *Usecase {
	return &Usecase{repo: repo, publisher: publisher}
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
	return u.repo.UpdateStatus(ctx, orderID, repository.StatusNew, repository.StatusCancelled, repository.OwnerFieldClient, clientID)
}

func (u *Usecase) UpdateStatusBySeller(ctx context.Context, orderID, sellerID uuid.UUID, targetStatus string) (*repository.Order, error) {
	if targetStatus != repository.StatusReady && targetStatus != repository.StatusCancelled {
		return nil, ErrInvalidTargetStat
	}
	order, err := u.repo.UpdateStatus(ctx, orderID, repository.StatusNew, targetStatus, repository.OwnerFieldSeller, sellerID)
	if err != nil {
		return nil, err
	}
	if order.Status == repository.StatusReady && u.publisher != nil {
		evt := events.OrderReadyEvent{
			OrderID:    order.ID.String(),
			ClientID:   order.ClientID.String(),
			SellerID:   order.SellerID.String(),
			Status:     order.Status,
			OccurredAt: time.Now().UTC(),
		}
		if err := u.publisher.PublishOrderReady(ctx, evt); err != nil {
			log.Printf("publish order.ready failed: order_id=%s err=%v", order.ID, err)
		}
	}
	return order, nil
}
