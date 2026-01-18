.PHONY: help build run test test_run clean docker-build docker-up docker-down docker-logs dev

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the Go application
	cd app && go build -o main ./cmd

run: ## Run the application locally
	cd app && go run ./cmd/main.go

test: ## Run tests
	cd app && go test -v ./...

test_run: ## Run unit tests with coverage
	cd app && go test -v -cover ./internal/...

clean: ## Remove build artifacts
	rm -f app/main
	cd app && go clean

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
		cd app && air; \
	else \
		echo "Air is not installed. Install it with: go install github.com/cosmtrek/air@latest"; \
		echo "Or run: make run"; \
	fi

init: ## Initialize the project (copy .env.example to .env)
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo ".env file created from .env.example"; \
	else \
		echo ".env file already exists"; \
	fi

deps: ## Download Go dependencies
	cd app && go mod download
	cd app && go mod tidy

fmt: ## Format Go code
	cd app && go fmt ./...

lint: ## Run Go linter
	@if command -v golangci-lint > /dev/null; then \
		cd app && golangci-lint run; \
	else \
		echo "golangci-lint is not installed. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi