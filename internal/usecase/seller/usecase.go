package sellerusecase

import (
	"beer/internal/model"
	"beer/internal/repository/seller"
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

type Usecase struct {
	sellerRepo *seller.Repository
}

func NewUsecase(sellerRepo *seller.Repository) *Usecase {
	return &Usecase{sellerRepo: sellerRepo}
}

func (u *Usecase) GetSellers(ctx context.Context) ([]model.Seller, error) {
	return u.sellerRepo.GetSellers(ctx)
}

func (u *Usecase) GetSellerByID(ctx context.Context, id uuid.UUID) (*model.Seller, error) {
	sellerEntity, err := u.sellerRepo.GetSellerByID(ctx, id)
	if err != nil {
		if errors.Is(err, seller.ErrSellerNotFound) {
			return nil, ErrSellerNotFound
		}
		return nil, err
	}
	return sellerEntity, nil
}

func (u *Usecase) CreateSeller(ctx context.Context, name string, login string, passwordHash string) (*model.Seller, error) {
	now := time.Now()
	sellerEntity := model.Seller{
		ID:           uuid.New(),
		Name:         name,
		Login:        login,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := u.sellerRepo.AddSeller(ctx, sellerEntity); err != nil {
		if errors.Is(err, seller.ErrLoginAlreadyExists) {
			return nil, ErrLoginAlreadyExists
		}
		return nil, err
	}
	return &sellerEntity, nil
}

func (u *Usecase) PatchSellerByID(ctx context.Context, id uuid.UUID, patch model.SellerPatch) (*model.Seller, error) {
	sellerEntity, err := u.sellerRepo.PatchSellerByID(ctx, id, patch)
	if err != nil {
		if errors.Is(err, seller.ErrSellerNotFound) {
			return nil, ErrSellerNotFound
		}
		if errors.Is(err, seller.ErrLoginAlreadyExists) {
			return nil, ErrLoginAlreadyExists
		}
		return nil, err
	}
	return sellerEntity, nil
}

func (u *Usecase) DeleteSellerByID(ctx context.Context, id uuid.UUID) error {
	deleted, err := u.sellerRepo.DeleteSellerByID(ctx, id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrSellerNotFound
	}
	return nil
}
