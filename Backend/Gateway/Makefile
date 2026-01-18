.PHONY: help build run test test_run clean docker-build docker-up docker-down docker-logs dev

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the Go application (cmd/gateway -> main)
	go build -o main ./cmd/gateway

run: ## Run the application locally
	go run ./cmd/gateway/main.go

test: ## Run tests
	go test -v ./...

test_run: ## Run unit tests with coverage
	go test -v -cover ./internal/...

clean: ## Remove build artifacts
	rm -f main
	go clean

docker-build: ## Build Docker image
	docker compose build

docker-up: ## Start Docker containers
	docker compose up -d

docker-down: ## Stop Docker containers
	docker compose down

docker-logs: ## View Docker logs
	docker compose logs -f

docker-restart: docker-down docker-up ## Restart Docker containers

dev: ## Run in development mode with hot reload (requires air)
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air is not installed. Install it with: go install github.com/cosmtrek/air@latest"; \
		echo "Or run: make run"; \
	fi

init: ## Initialize gateway env (ensure configs/gateway.env exists)
	@if [ ! -f configs/gateway.env ]; then \
		echo "Create and configure configs/gateway.env"; \
		printf "SERVER_PORT=8080\n" > configs/gateway.env; \
		echo "configs/gateway.env initialized with defaults"; \
	else \
		echo "configs/gateway.env already exists"; \
	fi

deps: ## Download Go dependencies
	go mod download
	go mod tidy

fmt: ## Format Go code
	go fmt ./...

lint: ## Run Go linter
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint is not installed. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi