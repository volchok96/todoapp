APP_ENV ?= local

run:
	@echo "Running app with APP_ENV=$(APP_ENV)"
	APP_ENV=$(APP_ENV) go run ./cmd/server/main.go

build:
	go build -o todoapp ./cmd/server/main.go

docker-up:
	@echo "Starting Docker containers..."
	sudo docker-compose up -d --build

docker-down:
	@echo "Stopping and removing Docker containers and volumes..."
	sudo docker-compose down -v

docker-psql:
	sudo docker-compose exec db psql -U postgres -d todoapp

docker-logs:
	sudo docker-compose logs -f app

docker-restart:
	make docker-down
	make docker-up
