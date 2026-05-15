package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type OrderReadyEvent struct {
	OrderID    string    `json:"order_id"`
	ClientID   string    `json:"client_id"`
	SellerID   string    `json:"seller_id"`
	Status     string    `json:"status"`
	OccurredAt time.Time `json:"occurred_at"`
}

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 1,
			MaxBytes: 10 << 20,
		}),
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	log.Printf("notification-service: subscribed to topic=%s group=%s", c.reader.Config().Topic, c.reader.Config().GroupID)
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			return err
		}
		var evt OrderReadyEvent
		if err := json.Unmarshal(msg.Value, &evt); err != nil {
			log.Printf("malformed event at offset=%d partition=%d: %v", msg.Offset, msg.Partition, err)
			continue
		}
		log.Printf("ORDER READY: order_id=%s client_id=%s seller_id=%s occurred_at=%s",
			evt.OrderID, evt.ClientID, evt.SellerID, evt.OccurredAt.Format(time.RFC3339))
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
