package storage

import (
	"beer/internal/model"
	"time"
	"github.com/google/uuid"
)

var clients = []model.Client{}

func GetClients() []model.Client {
	return clients
}

func AddClient(client model.Client) {
	clients = append(clients, client)
}

func GetClientByID(id uuid.UUID) (model.Client, bool) {
	for _, client := range clients {
		if client.ID == id {
			return client, true
		}
	}
	return model.Client{}, false
}

func DeleteClientByID(id uuid.UUID) bool {
	for i, client := range clients {
		if client.ID == id {
			clients = append(clients[:i], clients[i+1:]...)
			return true
		}
	}
	return false
}

func PatchClientByID(
	id uuid.UUID,
	name *string,
	phone *string,
	email *string,
	login *string,
	passwordHash *string,
) (model.Client, bool) {
	for i := range clients {
		if clients[i].ID == id {
			if name != nil {
				clients[i].Name = *name
			}
			if phone != nil {
				clients[i].Phone = *phone
			}
			if email != nil {
				clients[i].Email = *email
			}
			if login != nil {
				clients[i].Login = *login
			}
			if passwordHash != nil {
				clients[i].PasswordHash = *passwordHash
			}
			clients[i].UpdatedAt = time.Now()
			return clients[i], true
		}
	}
	return model.Client{}, false
}
