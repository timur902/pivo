package server

import (
	"context"
	"errors"

	"beer/order-service/internal/repository"
	"beer/order-service/internal/usecase"
	"beer/proto/order"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	orderpb.UnimplementedOrderServiceServer
	uc *usecase.Usecase
}

func NewServer(uc *usecase.Usecase) *Server {
	return &Server{uc: uc}
}

func (s *Server) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.Order, error) {
	clientID, err := uuid.Parse(req.GetClientId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid client_id")
	}
	sellerID, err := uuid.Parse(req.GetSellerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid seller_id")
	}
	items := make([]repository.NewOrderItem, 0, len(req.GetItems()))
	for _, it := range req.GetItems() {
		positionID, err := uuid.Parse(it.GetPositionId())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid position_id")
		}
		items = append(items, repository.NewOrderItem{
			PositionID: positionID,
			Quantity:   int(it.GetQuantity()),
			Price:      it.GetPrice(),
		})
	}
	order, err := s.uc.CreateOrder(ctx, clientID, sellerID, items)
	if err != nil {
		return nil, mapError(err)
	}
	return toProto(order), nil
}

func (s *Server) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	orderID, err := uuid.Parse(req.GetOrderId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order_id")
	}
	clientID, err := uuid.Parse(req.GetClientId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid client_id")
	}
	order, err := s.uc.GetOrderForClient(ctx, orderID, clientID)
	if err != nil {
		return nil, mapError(err)
	}
	return toProto(order), nil
}

func (s *Server) ListOrdersByClient(ctx context.Context, req *orderpb.ListOrdersByClientRequest) (*orderpb.ListOrdersResponse, error) {
	clientID, err := uuid.Parse(req.GetClientId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid client_id")
	}
	limit, offset := normalizePaging(req.GetLimit(), req.GetOffset())
	orders, err := s.uc.ListOrdersByClient(ctx, clientID, limit, offset)
	if err != nil {
		return nil, mapError(err)
	}
	return toListProto(orders), nil
}

func (s *Server) ListOrdersBySeller(ctx context.Context, req *orderpb.ListOrdersBySellerRequest) (*orderpb.ListOrdersResponse, error) {
	sellerID, err := uuid.Parse(req.GetSellerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid seller_id")
	}
	limit, offset := normalizePaging(req.GetLimit(), req.GetOffset())
	orders, err := s.uc.ListOrdersBySeller(ctx, sellerID, limit, offset)
	if err != nil {
		return nil, mapError(err)
	}
	return toListProto(orders), nil
}

func (s *Server) CancelOrderByClient(ctx context.Context, req *orderpb.CancelOrderByClientRequest) (*orderpb.Order, error) {
	orderID, err := uuid.Parse(req.GetOrderId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order_id")
	}
	clientID, err := uuid.Parse(req.GetClientId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid client_id")
	}
	order, err := s.uc.CancelOrderByClient(ctx, orderID, clientID)
	if err != nil {
		return nil, mapError(err)
	}
	return toProto(order), nil
}

func (s *Server) UpdateStatusBySeller(ctx context.Context, req *orderpb.UpdateStatusBySellerRequest) (*orderpb.Order, error) {
	orderID, err := uuid.Parse(req.GetOrderId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order_id")
	}
	sellerID, err := uuid.Parse(req.GetSellerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid seller_id")
	}
	order, err := s.uc.UpdateStatusBySeller(ctx, orderID, sellerID, req.GetStatus())
	if err != nil {
		return nil, mapError(err)
	}
	return toProto(order), nil
}

func normalizePaging(limit, offset int32) (int, int) {
	l := int(limit)
	o := int(offset)
	if l <= 0 || l > 100 {
		l = 20
	}
	if o < 0 {
		o = 0
	}
	return l, o
}

func toProto(o *repository.Order) *orderpb.Order {
	items := make([]*orderpb.OrderItem, 0, len(o.Items))
	for _, it := range o.Items {
		items = append(items, &orderpb.OrderItem{
			Id:         it.ID.String(),
			OrderId:    it.OrderID.String(),
			PositionId: it.PositionID.String(),
			Quantity:   int32(it.Quantity),
			Price:      it.Price,
		})
	}
	return &orderpb.Order{
		Id:        o.ID.String(),
		ClientId:  o.ClientID.String(),
		SellerId:  o.SellerID.String(),
		Status:    o.Status,
		CreatedAt: o.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: o.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Items:     items,
	}
}

func toListProto(orders []repository.Order) *orderpb.ListOrdersResponse {
	out := make([]*orderpb.Order, 0, len(orders))
	for i := range orders {
		out = append(out, toProto(&orders[i]))
	}
	return &orderpb.ListOrdersResponse{Orders: out}
}

func mapError(err error) error {
	switch {
	case errors.Is(err, repository.ErrOrderNotFound):
		return status.Error(codes.NotFound, "order not found")
	case errors.Is(err, repository.ErrClientNotFound):
		return status.Error(codes.NotFound, "client not found")
	case errors.Is(err, repository.ErrSellerNotFound):
		return status.Error(codes.NotFound, "seller not found")
	case errors.Is(err, repository.ErrPositionNotFound):
		return status.Error(codes.NotFound, "position not found")
	case errors.Is(err, repository.ErrOrderNotOwnedByActor):
		return status.Error(codes.PermissionDenied, "order does not belong to actor")
	case errors.Is(err, repository.ErrInvalidStatusChange):
		return status.Error(codes.FailedPrecondition, "invalid status transition")
	case errors.Is(err, usecase.ErrNoItems),
		errors.Is(err, usecase.ErrInvalidQuantity),
		errors.Is(err, usecase.ErrInvalidPrice),
		errors.Is(err, usecase.ErrInvalidTargetStat):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
