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

var pool *pgxpool.Pool
var ErrLoginAlreadyExists = errors.New("login already exists")

func SetPool(p *pgxpool.Pool) {
	pool = p
}

func checkPool() error {
	if pool == nil {
		return errors.New("client repository is not initialized")
	}
	return nil
}

func GetClients(ctx context.Context) ([]model.Client, error) {
	if err := checkPool(); err != nil {
		return nil, err
	}
	rows, err := pool.Query(ctx, `
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

func GetClientByID(ctx context.Context, id uuid.UUID) (model.Client, bool, error) {
	if err := checkPool(); err != nil {
		return model.Client{}, false, err
	}
	var c model.Client
	err := pool.QueryRow(ctx, `
		SELECT id,name,phone,email,login,password_hash,created_at,updated_at
		FROM clients
		WHERE id = $1
	`, id).Scan(&c.ID, &c.Name, &c.Phone, &c.Email, &c.Login, &c.PasswordHash, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Client{}, false, nil
	}
	if err != nil {
		return model.Client{}, false, err
	}
	return c, true, nil
}

func AddClient(ctx context.Context, c model.Client) error {
	if err := checkPool(); err != nil {
		return err
	}
	_, err := pool.Exec(ctx, `
		INSERT INTO clients (id,name,phone,email,login,password_hash,created_at,updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, c.ID, c.Name, c.Phone, c.Email, c.Login, c.PasswordHash, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return mapClientWriteError(err)
	}
	return nil
}

func DeleteClientByID(ctx context.Context, id uuid.UUID) (bool, error) {
	if err := checkPool(); err != nil {
		return false, err
	}
	tag, err := pool.Exec(ctx, `DELETE FROM clients WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func PatchClientByID(ctx context.Context, id uuid.UUID, name *string, phone *string, email *string, login *string, passwordHash *string) (model.Client, bool, error) {
	if err := checkPool(); err != nil {
		return model.Client{}, false, err
	}
	current, ok, err := GetClientByID(ctx, id)
	if err != nil {
		return model.Client{}, false, err
	}
	if !ok {
		return model.Client{}, false, nil
	}
	if name != nil {
		current.Name = *name
	}
	if phone != nil {
		current.Phone = *phone
	}
	if email != nil {
		current.Email = *email
	}
	if login != nil {
		current.Login = *login
	}
	if passwordHash != nil {
		current.PasswordHash = *passwordHash
	}
	current.UpdatedAt = time.Now()
	_, err = pool.Exec(ctx, `
		UPDATE clients
		SET name=$2,phone=$3,email=$4,login=$5,password_hash=$6,updated_at=$7
		WHERE id=$1
	`, id, current.Name, current.Phone, current.Email, current.Login, current.PasswordHash, current.UpdatedAt)
	if err != nil {
		return model.Client{}, false, mapClientWriteError(err)
	}
	return current, true, nil
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
