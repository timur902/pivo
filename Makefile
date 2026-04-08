GOOSE_VERSION ?= v3.26.0
GOOSE ?= go run github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)
GOOSE_DRIVER ?= postgres
GOOSE_DBSTRING ?= postgres://beer_user:beer_password@localhost:5432/beer?sslmode=disable
GOOSE_DIR ?= ./migrations

.PHONY: migrate-status migrate-up migrate-down migrate-create

migrate-status:
	$(GOOSE) -dir $(GOOSE_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" status

migrate-up:
	$(GOOSE) -dir $(GOOSE_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" up

migrate-down:
	$(GOOSE) -dir $(GOOSE_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" down

migrate-create:
	$(GOOSE) -dir $(GOOSE_DIR) create $(name) sql
