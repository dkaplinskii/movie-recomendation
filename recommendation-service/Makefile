-include .env

.PHONY: up up-d down logs

up:
	docker compose up --build

up-d:
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f $(APP_NAME)