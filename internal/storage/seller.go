package storage

import (
	"beer/internal/model"
	"time"
	"github.com/google/uuid"
)

var sellers = []model.Seller{}

func GetSellers() []model.Seller {
	return sellers
}

func AddSeller(seller model.Seller) {
	sellers = append(sellers, seller)
}

func GetSellerByID(id uuid.UUID) (model.Seller, bool) {
	for _, seller := range sellers {
		if seller.ID == id {
			return seller, true
		}
	}
	return model.Seller{}, false
}

func DeleteSellerByID(id uuid.UUID) bool {
	for i, seller := range sellers {
		if seller.ID == id {
			sellers = append(sellers[:i], sellers[i+1:]...)
			return true
		}
	}
	return false
}

func PatchSellerByID(
	id uuid.UUID,
	name *string,
	login *string,
	passwordHash *string,
) (model.Seller, bool) {
	for i := range sellers {
		if sellers[i].ID == id {
			if name != nil {
				sellers[i].Name = *name
			}
			if login != nil {
				sellers[i].Login = *login
			}
			if passwordHash != nil {
				sellers[i].PasswordHash = *passwordHash
			}
			sellers[i].UpdatedAt = time.Now()
			return sellers[i], true
		}
	}
	return model.Seller{}, false
}
