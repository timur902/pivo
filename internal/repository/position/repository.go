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

var pool *pgxpool.Pool

func SetPool(p *pgxpool.Pool) {
	pool = p
}

func checkPool() error {
	if pool == nil {
		return errors.New("position repository is not initialized")
	}
	return nil
}

func GetPositions(ctx context.Context, limit int, offset int) ([]model.Position, error) {
	if err := checkPool(); err != nil {
		return nil, err
	}
	rows, err := pool.Query(ctx, `
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

func GetPositionByID(ctx context.Context, id uuid.UUID) (model.Position, bool, error) {
	if err := checkPool(); err != nil {
		return model.Position{}, false, err
	}
	var p model.Position
	err := pool.QueryRow(ctx, `
		SELECT id,name,description,image_url,size_liters,quantity,price,created_at,updated_at
		FROM positions
		WHERE id = $1
	`, id).Scan(&p.ID, &p.Name, &p.Description, &p.ImageURL, &p.SizeLiters, &p.Quantity, &p.Price, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Position{}, false, nil
	}
	if err != nil {
		return model.Position{}, false, err
	}
	return p, true, nil
}

func AddPosition(ctx context.Context, p model.Position) error {
	if err := checkPool(); err != nil {
		return err
	}
	_, err := pool.Exec(ctx, `
		INSERT INTO positions (id,name,description,image_url,size_liters,quantity,price,created_at,updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, p.ID, p.Name, p.Description, p.ImageURL, p.SizeLiters, p.Quantity, p.Price, p.CreatedAt, p.UpdatedAt)
	return err
}

func DeletePositionByID(ctx context.Context, id uuid.UUID) (bool, error) {
	if err := checkPool(); err != nil {
		return false, err
	}
	tag, err := pool.Exec(ctx, `DELETE FROM positions WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func PatchPositionByID(ctx context.Context, id uuid.UUID, name *string, description *string, imageURL *string, sizeLiters *float32, quantity *int, price *int64) (model.Position, bool, error) {
	if err := checkPool(); err != nil {
		return model.Position{}, false, err
	}
	current, ok, err := GetPositionByID(ctx, id)
	if err != nil {
		return model.Position{}, false, err
	}
	if !ok {
		return model.Position{}, false, nil
	}
	if name != nil {
		current.Name = *name
	}
	if description != nil {
		current.Description = *description
	}
	if imageURL != nil {
		current.ImageURL = *imageURL
	}
	if sizeLiters != nil {
		current.SizeLiters = *sizeLiters
	}
	if quantity != nil {
		current.Quantity = *quantity
	}
	if price != nil {
		current.Price = *price
	}
	current.UpdatedAt = time.Now()
	_, err = pool.Exec(ctx, `
		UPDATE positions
		SET name=$2,description=$3,image_url=$4,size_liters=$5,quantity=$6,price=$7,updated_at=$8
		WHERE id=$1
	`, id, current.Name, current.Description, current.ImageURL, current.SizeLiters, current.Quantity, current.Price, current.UpdatedAt)
	if err != nil {
		return model.Position{}, false, err
	}
	return current, true, nil
}
