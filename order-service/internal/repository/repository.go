package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateOrder(ctx context.Context, clientID, sellerID uuid.UUID, items []NewOrderItem) (*Order, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var clientExists bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM clients WHERE id = $1)`, clientID).Scan(&clientExists); err != nil {
		return nil, err
	}
	if !clientExists {
		return nil, ErrClientNotFound
	}

	var sellerExists bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM sellers WHERE id = $1)`, sellerID).Scan(&sellerExists); err != nil {
		return nil, err
	}
	if !sellerExists {
		return nil, ErrSellerNotFound
	}

	for _, it := range items {
		var positionExists bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM positions WHERE id = $1)`, it.PositionID).Scan(&positionExists); err != nil {
			return nil, err
		}
		if !positionExists {
			return nil, ErrPositionNotFound
		}
	}

	now := time.Now().UTC()
	order := Order{
		ID:        uuid.New(),
		ClientID:  clientID,
		SellerID:  sellerID,
		Status:    "new",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO orders (id, client_id, seller_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, order.ID, order.ClientID, order.SellerID, order.Status, order.CreatedAt, order.UpdatedAt); err != nil {
		return nil, err
	}

	order.Items = make([]OrderItem, 0, len(items))
	for _, it := range items {
		oi := OrderItem{
			ID:         uuid.New(),
			OrderID:    order.ID,
			PositionID: it.PositionID,
			Quantity:   it.Quantity,
			Price:      it.Price,
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO order_items (id, order_id, position_id, quantity, price)
			VALUES ($1, $2, $3, $4, $5)
		`, oi.ID, oi.OrderID, oi.PositionID, oi.Quantity, oi.Price); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, oi)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *Repository) GetOrderByID(ctx context.Context, id uuid.UUID) (*Order, error) {
	var o Order
	err := r.pool.QueryRow(ctx, `
		SELECT id, client_id, seller_id, status, created_at, updated_at
		FROM orders
		WHERE id = $1
	`, id).Scan(&o.ID, &o.ClientID, &o.SellerID, &o.Status, &o.CreatedAt, &o.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	items, err := r.getItems(ctx, o.ID)
	if err != nil {
		return nil, err
	}
	o.Items = items
	return &o, nil
}

func (r *Repository) getItems(ctx context.Context, orderID uuid.UUID) ([]OrderItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, order_id, position_id, quantity, price
		FROM order_items
		WHERE order_id = $1
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]OrderItem, 0)
	for rows.Next() {
		var it OrderItem
		if err := rows.Scan(&it.ID, &it.OrderID, &it.PositionID, &it.Quantity, &it.Price); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) ListOrdersByClient(ctx context.Context, clientID uuid.UUID, limit, offset int) ([]Order, error) {
	return r.listOrders(ctx, `client_id = $1`, clientID, limit, offset)
}

func (r *Repository) ListOrdersBySeller(ctx context.Context, sellerID uuid.UUID, limit, offset int) ([]Order, error) {
	return r.listOrders(ctx, `seller_id = $1`, sellerID, limit, offset)
}

func (r *Repository) listOrders(ctx context.Context, where string, actorID uuid.UUID, limit, offset int) ([]Order, error) {
	query := `
		SELECT id, client_id, seller_id, status, created_at, updated_at
		FROM orders
		WHERE ` + where + `
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, actorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := make([]Order, 0, limit)
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.ClientID, &o.SellerID, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i := range orders {
		items, err := r.getItems(ctx, orders[i].ID)
		if err != nil {
			return nil, err
		}
		orders[i].Items = items
	}
	return orders, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, orderID uuid.UUID, expectedCurrent, next string, ownerField string, ownerID uuid.UUID) (*Order, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var current Order
	err = tx.QueryRow(ctx, `
		SELECT id, client_id, seller_id, status, created_at, updated_at
		FROM orders
		WHERE id = $1
		FOR UPDATE
	`, orderID).Scan(&current.ID, &current.ClientID, &current.SellerID, &current.Status, &current.CreatedAt, &current.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}

	switch ownerField {
	case "client_id":
		if current.ClientID != ownerID {
			return nil, ErrOrderNotOwnedByActor
		}
	case "seller_id":
		if current.SellerID != ownerID {
			return nil, ErrOrderNotOwnedByActor
		}
	}

	if current.Status != expectedCurrent {
		return nil, ErrInvalidStatusChange
	}

	now := time.Now().UTC()
	if _, err := tx.Exec(ctx, `UPDATE orders SET status = $2, updated_at = $3 WHERE id = $1`, orderID, next, now); err != nil {
		return nil, err
	}
	current.Status = next
	current.UpdatedAt = now

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	items, err := r.getItems(ctx, current.ID)
	if err != nil {
		return nil, err
	}
	current.Items = items
	return &current, nil
}
