# Makefile for Arya Gomoku Backend

.PHONY: help build run test clean migrate-up migrate-down migrate-status migrate-reset dev

# Default target
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  dev           - Run in development mode with auto-reload"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  migrate-up    - Apply all pending migrations"
	@echo "  migrate-down  - Rollback the last migration"
	@echo "  migrate-status- Show migration status"
	@echo "  migrate-reset - Reset database (WARNING: drops all data)"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/gomoku main.go

# Run the application
run: build
	@echo "Starting application..."
	./bin/gomoku

# Development mode with auto-reload (requires air)
dev:
	@echo "Starting development server..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Falling back to normal run..."; \
		make run; \
	fi

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Database migrations
migrate-up:
	@echo "Applying migrations..."
	go run cmd/migrate/main.go up

migrate-down:
	@echo "Rolling back migration..."
	go run cmd/migrate/main.go down

migrate-status:
	@echo "Checking migration status..."
	go run cmd/migrate/main.go status

migrate-reset:
	@echo "Resetting database..."
	go run cmd/migrate/main.go reset

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Generate documentation
docs:
	@echo "Generating documentation..."
	@if command -v godoc > /dev/null; then \
		echo "Starting documentation server at http://localhost:6060"; \
		godoc -http=:6060; \
	else \
		echo "godoc not found. Install with: go install golang.org/x/tools/cmd/godoc@latest"; \
	fi

# Database setup (create database and run migrations)
db-setup: migrate-up
	@echo "Database setup completed"

# Full setup for new development environment
setup: deps db-setup
	@echo "Development environment setup completed"

# Production build
build-prod:
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/gomoku main.go

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t arya-gomoku-backend .

docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 arya-gomoku-backend

# Benchmark AI
benchmark:
	@echo "Running AI benchmark..."
	go run cmd/ai_demo/main.go