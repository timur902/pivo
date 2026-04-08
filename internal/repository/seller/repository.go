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

var pool *pgxpool.Pool
var ErrLoginAlreadyExists = errors.New("login already exists")

func SetPool(p *pgxpool.Pool) {
	pool = p
}

func checkPool() error {
	if pool == nil {
		return errors.New("seller repository is not initialized")
	}
	return nil
}

func GetSellers(ctx context.Context) ([]model.Seller, error) {
	if err := checkPool(); err != nil {
		return nil, err
	}
	rows, err := pool.Query(ctx, `
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

func GetSellerByID(ctx context.Context, id uuid.UUID) (model.Seller, bool, error) {
	if err := checkPool(); err != nil {
		return model.Seller{}, false, err
	}
	var s model.Seller
	err := pool.QueryRow(ctx, `
		SELECT id,name,login,password_hash,created_at,updated_at
		FROM sellers
		WHERE id = $1
	`, id).Scan(&s.ID, &s.Name, &s.Login, &s.PasswordHash, &s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Seller{}, false, nil
	}
	if err != nil {
		return model.Seller{}, false, err
	}
	return s, true, nil
}

func AddSeller(ctx context.Context, s model.Seller) error {
	if err := checkPool(); err != nil {
		return err
	}
	_, err := pool.Exec(ctx, `
		INSERT INTO sellers (id,name,login,password_hash,created_at,updated_at)
		VALUES ($1,$2,$3,$4,$5,$6)
	`, s.ID, s.Name, s.Login, s.PasswordHash, s.CreatedAt, s.UpdatedAt)
	if err != nil {
		return mapSellerWriteError(err)
	}
	return nil
}

func DeleteSellerByID(ctx context.Context, id uuid.UUID) (bool, error) {
	if err := checkPool(); err != nil {
		return false, err
	}
	tag, err := pool.Exec(ctx, `DELETE FROM sellers WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func PatchSellerByID(ctx context.Context, id uuid.UUID, name *string, login *string, passwordHash *string) (model.Seller, bool, error) {
	if err := checkPool(); err != nil {
		return model.Seller{}, false, err
	}
	current, ok, err := GetSellerByID(ctx, id)
	if err != nil {
		return model.Seller{}, false, err
	}
	if !ok {
		return model.Seller{}, false, nil
	}
	if name != nil {
		current.Name = *name
	}
	if login != nil {
		current.Login = *login
	}
	if passwordHash != nil {
		current.PasswordHash = *passwordHash
	}
	current.UpdatedAt = time.Now()
	_, err = pool.Exec(ctx, `
		UPDATE sellers
		SET name=$2,login=$3,password_hash=$4,updated_at=$5
		WHERE id=$1
	`, id, current.Name, current.Login, current.PasswordHash, current.UpdatedAt)
	if err != nil {
		return model.Seller{}, false, mapSellerWriteError(err)
	}
	return current, true, nil
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
