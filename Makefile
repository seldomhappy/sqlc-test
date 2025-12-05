.PHONY: help build run test cover fmt lint sqlc db-up db-down db-logs clean vet

# Display help information
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Build and run the application"
	@echo "  test          - Run tests"
	@echo "  cover         - Run tests with coverage report"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run go vet"
	@echo "  lint          - Run linting via Docker"
	@echo "  sqlc          - Generate SQL code via Docker"
	@echo "  db-up         - Start PostgreSQL"
	@echo "  db-down       - Stop and remove PostgreSQL"
	@echo "  db-logs       - Show PostgreSQL logs"
	@echo "  clean         - Clean build artifacts"

build: ## Build the application
	@echo "Building application..."
	go build -o ./bin/app ./cmd/app

run: build ## Build and run the application
	@echo "Running application..."
	./bin/app

test: ## Run tests
	@echo "Running tests..."
	go test ./...

cover: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -cover ./...

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

lint: ## Run linting via Docker
	@echo "Running golangci-lint in Docker..."
	docker run --rm -v "$(CURDIR)":/app -w /app golangci/golangci-lint:v2.7.1 golangci-lint run --no-config --timeout 5m

sqlc: ## Generate SQL code via Docker
	@echo "Running sqlc in Docker..."
	docker run --rm -v "$(CURDIR)":/app -w /app sqlc/sqlc generate

db-up: ## Start PostgreSQL
	@echo "Starting Postgres via docker-compose..."
	docker-compose -f docker-compose.yml up -d db

db-down: ## Stop and remove PostgreSQL
	@echo "Stopping Postgres and removing containers..."
	docker-compose -f docker-compose.yml down -v

db-logs: ## Show PostgreSQL logs
	@echo "Tailing Postgres logs (press Ctrl+C to stop)..."
	docker-compose -f docker-compose.yml logs -f db

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf ./bin/
	go clean

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...
