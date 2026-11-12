-include .env

.PHONY: run build up down logs lint

run:
	go run main.go

up:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f