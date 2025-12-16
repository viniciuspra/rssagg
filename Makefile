-include .env
export

.PHONY: sqlc g-up g-down g-status http

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)
MIGRATIONS_DIR=sql/migrations

dev:
	go build && ./rssagg

health:
	http GET :8080/v1/healthz

sqlc:
	docker compose run --rm sqlc generate

p-up:
	docker compose up -d postgres

p-stop:
	docker compose stop postgres

g-up: p-up
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_URL) up

g-down:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_URL) down

g-status:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_URL) status
