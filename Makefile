DB_URL := postgres://app:app@localhost:5432/shop?sslmode=disable

up:
	docker compose -f deploy/compose/compose.yaml up -d

down:
	docker compose -f deploy/compose/compose.yaml down

run:
	DATABASE_URL=$(DB_URL) go run ./cmd/shop
