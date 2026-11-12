-include .env

.PHONY: up up-d down logs restart ps health-postgres

up:
	docker compose up --build

up-d:
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f movie-service

restart:
	docker compose restart movie-service

ps:
	docker compose ps

health-postgres:
	docker compose exec postgres pg_isready -U postgres -d movie_service