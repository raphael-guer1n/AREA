.PHONY: help backend-up backend-down frontend-up frontend-down docker-up docker-down

BACKEND_DIR := Backend
FRONTEND_DIR := Web/frontend

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

backend-up: ## Start backend via Backend/Makefile
	@$(MAKE) -C $(BACKEND_DIR) docker-up

frontend-up: ## Start frontend Docker containers
	@$(MAKE) -C $(FRONTEND_DIR) docker-up

docker-up: backend-up frontend-up ## Start backend then frontend

frontend-down: ## Stop frontend Docker containers
	@$(MAKE) -C $(FRONTEND_DIR) docker-down

backend-down: ## Stop backend via Backend/Makefile
	@$(MAKE) -C $(BACKEND_DIR) docker-down

docker-down: frontend-down backend-down ## Stop frontend then backend
