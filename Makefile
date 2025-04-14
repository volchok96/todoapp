.PHONY: go-generate run build test lint \
	docker-up docker-down docker-restart docker-psql docker-logs-app docker-logs-db \
	db-create db-migrate db-seed db-reset \
	docker-clean docker-clean-all docker local-start \
	help

APP_ENV ?= local
DB_NAME ?= todoapp
DB_USER ?= postgres

## ----------- Local Go Commands -----------

go-generate:
	swag init -g cmd/server/main.go --output docs --dir ./ --parseDependency

run:
	@echo "Running app with APP_ENV=$(APP_ENV)"
	APP_ENV=$(APP_ENV) go run ./cmd/server/main.go

build:
	go build -o todoapp ./cmd/server/main.go

test:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

lint:
	golangci-lint run ./...

## ----------- Docker Commands -----------

docker-up:
	@echo "Starting Docker containers..."
	sudo docker-compose up -d --build

docker-down:
	@echo "Stopping and removing Docker containers and volumes..."
	sudo docker-compose down -v

docker-restart: docker-down docker-up

docker-psql:
	sudo docker-compose exec db psql -U $(DB_USER) -d $(DB_NAME)

docker-logs-app:
	sudo docker-compose logs -f app

docker-logs-db:
	sudo docker-compose logs -f db

docker-clean: docker-down
	sudo docker-compose rm -f

docker-clean-all: docker-clean
	sudo docker-compose down --rmi all --volumes --remove-orphans

docker: build docker-up

## ----------- POSTGRESQL Commands (Local) ----------

db-create:
	psql -U $(DB_USER) -f migrations/01_create_db.sql

db-migrate:
	psql -U $(DB_USER) -d $(DB_NAME) -f migrations/02_up.sql

db-seed:
	psql -U $(DB_USER) -d $(DB_NAME) -f migrations/04_test_data.sql

db-reset:
	psql -U $(DB_USER) -d $(DB_NAME) -f migrations/03_down.sql
	psql -U $(DB_USER) -d $(DB_NAME) -f migrations/02_up.sql
	psql -U $(DB_USER) -d $(DB_NAME) -f migrations/04_test_data.sql

local-start: db-reset run

## ----------- Help -----------

help:
	@echo "Makefile commands:"
	@echo ""
	@echo "  go-generate        Generate swagger"
	@echo "  run                Run the app locally"
	@echo "  build              Build the app binary"
	@echo "  test               Run unit tests with coverage"
	@echo "  lint               Run golangci-lint (must be installed)"
	@echo ""
	@echo "  docker-up          Start docker-compose services"
	@echo "  docker-down        Stop and remove containers and volumes"
	@echo "  docker-restart     Restart docker services"
	@echo "  docker-psql        Connect to DB inside docker"
	@echo "  docker-logs-app    Tail app logs"
	@echo "  docker-logs-db     Tail DB logs"
	@echo "  docker-clean              Remove stopped containers"
	@echo "  docker-clean-all          Remove all containers/images/volumes"
	@echo "  docker             Build and run containers"
	@echo ""
	@echo "  db-create          Create database (local)"
	@echo "  db-migrate         Run migration script"
	@echo "  db-seed            Seed the database with test data"
	@echo "  db-reset           Drop, recreate and seed the database"
	@echo "  local-start        Run migrations, seed DB, and start the app locally"
