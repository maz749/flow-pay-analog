.PHONY: help build run test clean docker-build docker-up docker-down migrate

help: ## Показать эту справку
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Собрать приложение
	go build -o bin/flowpay ./cmd/server

run: ## Запустить приложение локально
	go run ./cmd/server/main.go

test: ## Запустить тесты
	go test -v ./...

clean: ## Очистить собранные файлы
	rm -rf bin/

deps: ## Установить зависимости
	go mod download
	go mod tidy

docker-build: ## Собрать Docker образ
	docker-compose build

docker-up: ## Запустить приложение в Docker
	docker-compose up -d

docker-down: ## Остановить Docker контейнеры
	docker-compose down

docker-logs: ## Показать логи Docker контейнеров
	docker-compose logs -f

migrate: ## Применить миграции базы данных (вручную)
	psql -h localhost -U flowpay -d flowpay_db -f migrations/001_init_schema.sql
	psql -h localhost -U flowpay -d flowpay_db -f migrations/002_add_cancellation_instructions.sql
	psql -h localhost -U flowpay -d flowpay_db -f migrations/003_add_teams_and_roles.sql

dev: ## Запустить в режиме разработки с hot reload
	air

.DEFAULT_GOAL := help
