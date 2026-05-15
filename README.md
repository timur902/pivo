# Database Migrations

Схема БД теперь управляется только миграциями `goose` из директории `/migrations`.

## Быстрый старт
1. Запустить PostgreSQL:
```bash
docker compose up -d postgres
```
2. Накатить миграции:
```bash
make migrate-up
```
3. Запустить в трёх разных терминалах:
``` bash
make run-beer-api
make run-order-service
make run-notification-service
```

## Полезные команды
```bash
make migrate-status
make migrate-down
make migrate-create name=add_new_table
```

По умолчанию используется DSN:
`postgres://beer_user:beer_password@localhost:5432/beer?sslmode=disable`
