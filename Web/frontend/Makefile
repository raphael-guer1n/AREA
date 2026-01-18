.PHONY: help install dev build start run test init deps fmt lint clean docker-build docker-up docker-down docker-logs docker-restart up

SHELL := /bin/bash
MAKEFLAGS += --no-print-directory

COMPOSE_FILE := docker-compose.yml
COMPOSE_ENV := .env
COMPOSE_ARGS := --env-file $(COMPOSE_ENV)

help: ## Show this help message
	@printf 'Usage: make [target]\n'
	@echo ''
	@printf 'Available targets:\n'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install npm dependencies
	npm install

dev: ## Run Next.js dev server
	npm run dev

build: ## Build for production
	npm run build

start: ## Start production server (requires build)
	npm run start

run: start ## Alias for start

test: ## Run tests (not configured)
	@echo "No test script configured in package.json."

init: install ## Initialize local environment

deps: install ## Install dependencies

fmt: ## Format code (not configured)
	@echo "No formatter script configured in package.json."

lint: ## Run eslint
	npm run lint

clean: ## Remove build artifacts
	rm -rf .next out

docker-build: ## Build Docker image via compose
	docker compose $(COMPOSE_ARGS) -f $(COMPOSE_FILE) build

docker-up: ## Start Docker containers
	docker compose $(COMPOSE_ARGS) -f $(COMPOSE_FILE) up -d --build

docker-down: ## Stop Docker containers
	docker compose $(COMPOSE_ARGS) -f $(COMPOSE_FILE) down

docker-logs: ## Tail Docker logs
	docker compose $(COMPOSE_ARGS) -f $(COMPOSE_FILE) logs -f

docker-restart: docker-down docker-up ## Restart Docker containers

up: docker-up ## Alias for docker-up
