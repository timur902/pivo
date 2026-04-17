package seller

import (
	"beer/internal/model"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
	"time"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetSellers(ctx context.Context) ([]model.Seller, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id,name,login,password_hash,created_at,updated_at
		FROM sellers
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	sellers := make([]model.Seller, 0)
	for rows.Next() {
		var s model.Seller
		if err := rows.Scan(&s.ID, &s.Name, &s.Login, &s.PasswordHash, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sellers = append(sellers, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sellers, nil
}

func (r *Repository) GetSellerByID(ctx context.Context, id uuid.UUID) (*model.Seller, error) {
	var s model.Seller
	err := r.pool.QueryRow(ctx, `
		SELECT id,name,login,password_hash,created_at,updated_at
		FROM sellers
		WHERE id = $1
	`, id).Scan(&s.ID, &s.Name, &s.Login, &s.PasswordHash, &s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrSellerNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) AddSeller(ctx context.Context, s model.Seller) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO sellers (id,name,login,password_hash,created_at,updated_at)
		VALUES ($1,$2,$3,$4,$5,$6)
	`, s.ID, s.Name, s.Login, s.PasswordHash, s.CreatedAt, s.UpdatedAt)
	if err != nil {
		return mapSellerWriteError(err)
	}
	return nil
}

func (r *Repository) DeleteSellerByID(ctx context.Context, id uuid.UUID) (bool, error) {
	tag, err := r.pool.Exec(ctx, `DELETE FROM sellers WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func (r *Repository) PatchSellerByID(ctx context.Context, id uuid.UUID, patch model.SellerPatch) (*model.Seller, error) {
	current, err := r.GetSellerByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if patch.Name != nil {
		current.Name = *patch.Name
	}
	if patch.Login != nil {
		current.Login = *patch.Login
	}
	if patch.PasswordHash != nil {
		current.PasswordHash = *patch.PasswordHash
	}
	current.UpdatedAt = time.Now()
	_, err = r.pool.Exec(ctx, `
		UPDATE sellers
		SET name=$2,login=$3,password_hash=$4,updated_at=$5
		WHERE id=$1
	`, id, current.Name, current.Login, current.PasswordHash, current.UpdatedAt)
	if err != nil {
		return nil, mapSellerWriteError(err)
	}
	return current, nil
}

func mapSellerWriteError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		if pgErr.ConstraintName == "sellers_login_key" || (pgErr.TableName == "sellers" && strings.Contains(pgErr.Detail, "(login)")) {
			return ErrLoginAlreadyExists
		}
	}
	return err
}
