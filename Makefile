# Загружаем .env (если есть) и пробрасываем все переменные в go run / goose.
# Минус перед include = не падать, если файла нет.
-include .env
export

GOOSE_VERSION ?= v3.26.0
GOOSE ?= go run github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)
GOOSE_DRIVER ?= postgres
GOOSE_DBSTRING ?= $(DATABASE_URL)
GOOSE_DIR ?= ./migrations

.PHONY: migrate-status migrate-up migrate-down migrate-create run-beer-api run-order-service run-notification-service

migrate-status:
	$(GOOSE) -dir $(GOOSE_DIR) status

migrate-up:
	$(GOOSE) -dir $(GOOSE_DIR) up

migrate-down:
	$(GOOSE) -dir $(GOOSE_DIR) down

migrate-create:
	$(GOOSE) -dir $(GOOSE_DIR) create $(name) sql

run-beer-api:
	go run ./beer-api/cmd

run-order-service:
	go run ./order-service/cmd

run-notification-service:
	go run ./notification-service/cmd
