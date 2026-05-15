package handler

import (
	"net/http"
	"strconv"

	"beer/proto/order"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateOrderItemRequest struct {
	PositionID string `json:"position_id"`
	Quantity   int32  `json:"quantity"`
	Price      int64  `json:"price"`
}

type CreateOrderRequest struct {
	ClientID string                   `json:"client_id"`
	SellerID string                   `json:"seller_id"`
	Items    []CreateOrderItemRequest `json:"items"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	items := make([]*orderpb.OrderItemInput, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, &orderpb.OrderItemInput{
			PositionId: it.PositionID,
			Quantity:   it.Quantity,
			Price:      it.Price,
		})
	}
	resp, err := h.orderClient.CreateOrder(c.Request.Context(), &orderpb.CreateOrderRequest{
		ClientId: req.ClientID,
		SellerId: req.SellerID,
		Items:    items,
	})
	if err != nil {
		writeGRPCError(c, err)
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) GetOrderByID(c *gin.Context) {
	clientID := c.Query("client_id")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id query parameter is required"})
		return
	}
	resp, err := h.orderClient.GetOrder(c.Request.Context(), &orderpb.GetOrderRequest{
		OrderId:  c.Param("id"),
		ClientId: clientID,
	})
	if err != nil {
		writeGRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ListClientOrders(c *gin.Context) {
	clientID := c.Query("client_id")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id query parameter is required"})
		return
	}
	limit, offset := parseLimitOffset(c)
	resp, err := h.orderClient.ListOrdersByClient(c.Request.Context(), &orderpb.ListOrdersByClientRequest{
		ClientId: clientID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		writeGRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp.GetOrders())
}

func (h *Handler) CancelOrderByClient(c *gin.Context) {
	clientID := c.Query("client_id")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id query parameter is required"})
		return
	}
	resp, err := h.orderClient.CancelOrderByClient(c.Request.Context(), &orderpb.CancelOrderByClientRequest{
		OrderId:  c.Param("id"),
		ClientId: clientID,
	})
	if err != nil {
		writeGRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ListSellerOrders(c *gin.Context) {
	sellerID := c.Query("seller_id")
	if sellerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "seller_id query parameter is required"})
		return
	}
	limit, offset := parseLimitOffset(c)
	resp, err := h.orderClient.ListOrdersBySeller(c.Request.Context(), &orderpb.ListOrdersBySellerRequest{
		SellerId: sellerID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		writeGRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp.GetOrders())
}

func (h *Handler) UpdateOrderStatusBySeller(c *gin.Context) {
	sellerID := c.Query("seller_id")
	if sellerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "seller_id query parameter is required"})
		return
	}
	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	resp, err := h.orderClient.UpdateStatusBySeller(c.Request.Context(), &orderpb.UpdateStatusBySellerRequest{
		OrderId:  c.Param("id"),
		SellerId: sellerID,
		Status:   req.Status,
	})
	if err != nil {
		writeGRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func parseLimitOffset(c *gin.Context) (int, int) {
	limit := 20
	offset := 0
	if raw := c.Query("limit"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}
	if raw := c.Query("offset"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v >= 0 {
			offset = v
		}
	}
	return limit, offset
}

func writeGRPCError(c *gin.Context, err error) {
	st, ok := status.FromError(err)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	httpCode := grpcCodeToHTTP(st.Code())
	c.JSON(httpCode, gin.H{"error": st.Message()})
}

func grpcCodeToHTTP(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.FailedPrecondition:
		return http.StatusConflict
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}
