package events

import "time"

const TopicOrdersReady = "orders.ready"

type OrderReadyEvent struct {
	OrderID    string    `json:"order_id"`
	ClientID   string    `json:"client_id"`
	SellerID   string    `json:"seller_id"`
	Status     string    `json:"status"`
	OccurredAt time.Time `json:"occurred_at"`
}
