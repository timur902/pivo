package position

import (
	"beer/internal/model"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetPositions(ctx context.Context, limit int, offset int) ([]model.Position, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id,name,description,image_url,size_liters,quantity,price,created_at,updated_at
		FROM positions
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	positions := make([]model.Position, 0, limit)
	for rows.Next() {
		var p model.Position
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.ImageURL, &p.SizeLiters, &p.Quantity, &p.Price, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return positions, nil
}

func (r *Repository) GetPositionByID(ctx context.Context, id uuid.UUID) (*model.Position, error) {
	var p model.Position
	err := r.pool.QueryRow(ctx, `
		SELECT id,name,description,image_url,size_liters,quantity,price,created_at,updated_at
		FROM positions
		WHERE id = $1
	`, id).Scan(&p.ID, &p.Name, &p.Description, &p.ImageURL, &p.SizeLiters, &p.Quantity, &p.Price, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPositionNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) AddPosition(ctx context.Context, p model.Position) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO positions (id,name,description,image_url,size_liters,quantity,price,created_at,updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, p.ID, p.Name, p.Description, p.ImageURL, p.SizeLiters, p.Quantity, p.Price, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *Repository) DeletePositionByID(ctx context.Context, id uuid.UUID) (bool, error) {
	tag, err := r.pool.Exec(ctx, `DELETE FROM positions WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func (r *Repository) PatchPositionByID(ctx context.Context, id uuid.UUID, patch model.PositionPatch) (*model.Position, error) {
	current, err := r.GetPositionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if patch.Name != nil {
		current.Name = *patch.Name
	}
	if patch.Description != nil {
		current.Description = *patch.Description
	}
	if patch.ImageURL != nil {
		current.ImageURL = *patch.ImageURL
	}
	if patch.SizeLiters != nil {
		current.SizeLiters = *patch.SizeLiters
	}
	if patch.Quantity != nil {
		current.Quantity = *patch.Quantity
	}
	if patch.Price != nil {
		current.Price = *patch.Price
	}
	current.UpdatedAt = time.Now()
	_, err = r.pool.Exec(ctx, `
		UPDATE positions
		SET name=$2,description=$3,image_url=$4,size_liters=$5,quantity=$6,price=$7,updated_at=$8
		WHERE id=$1
	`, id, current.Name, current.Description, current.ImageURL, current.SizeLiters, current.Quantity, current.Price, current.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return current, nil
}
