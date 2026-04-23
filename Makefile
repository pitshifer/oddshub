.PHONY: build run stop clean help dev logs migrate-up migrate-down

# Помощь
help:
	@echo "Available commands:"
	@echo "  make build       - Compile the application"
	@echo "  make run         - Build and start Docker containers"
	@echo "  make stop        - Stop Docker containers"
	@echo "  make clean       - Clean build artifacts and stop containers"
	@echo "  make dev         - Watch for changes and auto-rebuild"
	@echo "  make logs        - View container logs"
	@echo "  make migrate-up  - Run database migrations"
	@echo "  make migrate-down - Rollback migrations"

# Компиляция приложения
build:
	@echo "Building oddshub..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o oddshub ./cmd/api
	@echo "Build complete!"

# Запуск контейнеров
run: build
	@echo "Starting Docker containers..."
	docker compose up -d
	@echo "oddshub is running!"

# Остановка контейнеров
stop:
	@echo "Stopping Docker containers..."
	docker compose down

# Очистка проекта
clean: stop
	@echo "Cleaning up..."
	rm -f oddshub
	@echo "Clean complete!"

# Logs
logs:
	docker compose logs -f

# Migrations
migrate-up:
	@echo "Running migrations..."
	migrate -path ./migrations -database "postgres://oddshub:secret@localhost:5431/oddshub?sslmode=disable" up
	@echo "Migrations complete!"

migrate-down:
	@echo "Rolling back migrations..."
	migrate -path ./migrations -database "postgres://oddshub:secret@localhost:5431/oddshub?sslmode=disable" down
	@echo "Rollback complete!"


