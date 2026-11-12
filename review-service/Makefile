-include .env

.PHONY: up up-d down logs restart ps health-postgres proto

up:
	docker compose up --build

up-d:
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f review-service

restart:
	docker compose restart review-service

ps:
	docker compose ps

health-postgres:
	docker compose exec postgres pg_isready -U postgres -d review_service

proto:
	protoc --go_out=. --go-grpc_out=. proto/review.proto
