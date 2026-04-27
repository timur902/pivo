package client

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

func (r *Repository) GetClients(ctx context.Context) ([]model.Client, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id,name,phone,email,login,password_hash,created_at,updated_at
		FROM clients
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	clients := make([]model.Client, 0)
	for rows.Next() {
		var c model.Client
		if err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Email, &c.Login, &c.PasswordHash, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return clients, nil
}

func (r *Repository) GetClientByID(ctx context.Context, id uuid.UUID) (*model.Client, error) {
	var c model.Client
	err := r.pool.QueryRow(ctx, `
		SELECT id,name,phone,email,login,password_hash,created_at,updated_at
		FROM clients
		WHERE id = $1
	`, id).Scan(&c.ID, &c.Name, &c.Phone, &c.Email, &c.Login, &c.PasswordHash, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrClientNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) AddClient(ctx context.Context, c model.Client) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO clients (id,name,phone,email,login,password_hash,created_at,updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, c.ID, c.Name, c.Phone, c.Email, c.Login, c.PasswordHash, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return mapClientWriteError(err)
	}
	return nil
}

func (r *Repository) DeleteClientByID(ctx context.Context, id uuid.UUID) (bool, error) {
	tag, err := r.pool.Exec(ctx, `DELETE FROM clients WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func (r *Repository) PatchClientByID(ctx context.Context, id uuid.UUID, patch model.ClientPatch) (*model.Client, error) {
	current, err := r.GetClientByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if patch.Name != nil {
		current.Name = *patch.Name
	}
	if patch.Phone != nil {
		current.Phone = *patch.Phone
	}
	if patch.Email != nil {
		current.Email = *patch.Email
	}
	if patch.Login != nil {
		current.Login = *patch.Login
	}
	if patch.PasswordHash != nil {
		current.PasswordHash = *patch.PasswordHash
	}
	current.UpdatedAt = time.Now()
	_, err = r.pool.Exec(ctx, `
		UPDATE clients
		SET name=$2,phone=$3,email=$4,login=$5,password_hash=$6,updated_at=$7
		WHERE id=$1
	`, id, current.Name, current.Phone, current.Email, current.Login, current.PasswordHash, current.UpdatedAt)
	if err != nil {
		return nil, mapClientWriteError(err)
	}
	return current, nil
}

func mapClientWriteError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		if pgErr.ConstraintName == "clients_login_key" || (pgErr.TableName == "clients" && strings.Contains(pgErr.Detail, "(login)")) {
			return ErrLoginAlreadyExists
		}
	}
	return err
}
