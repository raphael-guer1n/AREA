.PHONY: help up down restart logs ps clean build config status health gateway-dashboard

# Colors for output
GREEN  := \033[0;32m
YELLOW := \033[0;33m
RED    := \033[0;31m
NC     := \033[0m # No Color

# Load environment variables from .env if it exists
-include .env
export

help: ## Show this help message
	@echo '$(GREEN)AREA Microservices Management$(NC)'
	@echo ''
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

config: ## Display services configuration
	@echo '$(GREEN)=== AREA Services Configuration ===$(NC)'
	@echo ''
	@echo '$(YELLOW)Database Credentials:$(NC)'
	@echo '  User: ${DB_USER:-postgres}'
	@echo '  Password: ${DB_PASSWORD:-postgres}'
	@echo ''
	@echo '$(YELLOW)Planned Services (from services.yaml):$(NC)'
	@cat services.yaml | grep -A 5 "  # [a-zA-Z]*:" | grep -E "name:|port:|path:" | sed 's/^/  /' || echo '  No services configured yet'

build: ## Build all microservices
	@echo '$(GREEN)Building all microservices...$(NC)'
	docker-compose build

init: ## Initialize environment (copy .env.example to .env)
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "$(GREEN)✓ .env file created from .env.example$(NC)"; \
		echo "$(YELLOW)Please review and update .env with your configuration$(NC)"; \
	else \
		echo "$(YELLOW).env file already exists$(NC)"; \
	fi

up: ## Start all services
	@if [ ! -f .env ]; then \
		echo "$(YELLOW)Warning: .env file not found. Run 'make init' first or using default values...$(NC)"; \
	fi
	@echo '$(GREEN)Starting all AREA microservices...$(NC)'
	docker-compose up -d
	@echo ''
	@echo '$(GREEN)✓ Services started!$(NC)'
	@echo ''
	@echo 'Run $(YELLOW)make logs$(NC) to view logs'
	@echo 'Run $(YELLOW)make ps$(NC) to check service status'

down: ## Stop all services
	@echo '$(RED)Stopping all AREA services...$(NC)'
	docker-compose down

restart: down up ## Restart all services

stop: down ## Alias for down

logs: ## View logs from all services
	docker-compose logs -f

logs-%: ## View logs from specific service (e.g., make logs-auth)
	docker-compose logs -f $*-service

ps: ## Show running services status
	@echo '$(GREEN)=== AREA Services Status ===$(NC)'
	@docker-compose ps

status: ps ## Alias for ps

health: ## Check health of all services
	@echo '$(GREEN)=== Health Check ===$(NC)'
	@echo ''
	@echo '$(YELLOW)Microservices:$(NC)'
	@docker-compose ps 2>/dev/null | grep -v "^NAME" | grep -v "^---" | awk '{print "  " $$1 ": " $$NF}' || echo '  No services running yet'
	@echo ''
	@echo '$(YELLOW)Add services to docker-compose.yml to see them here$(NC)'

clean: ## Stop and remove all containers, networks, and volumes
	@echo '$(RED)Removing all AREA containers, networks, and volumes...$(NC)'
	docker-compose down -v
	@echo '$(GREEN)✓ Cleanup complete$(NC)'

clean-full: clean ## Full cleanup including images
	@echo '$(RED)Removing all AREA images...$(NC)'
	docker-compose down -v --rmi all
	@echo '$(GREEN)✓ Full cleanup complete$(NC)'

prune: ## Remove unused Docker resources
	@echo '$(YELLOW)Pruning Docker system...$(NC)'
	docker system prune -f
	@echo '$(GREEN)✓ Prune complete$(NC)'

# Individual service management
start-%: ## Start specific service (e.g., make start-auth)
	@echo '$(GREEN)Starting $* service...$(NC)'
	docker-compose up -d $*-service $*-db

stop-%: ## Stop specific service (e.g., make stop-auth)
	@echo '$(RED)Stopping $* service...$(NC)'
	docker-compose stop $*-service $*-db

restart-%: ## Restart specific service (e.g., make restart-auth)
	@echo '$(YELLOW)Restarting $* service...$(NC)'
	docker-compose restart $*-service

# Development helpers
shell-%: ## Open shell in service container (e.g., make shell-auth)
	docker-compose exec $*-service sh

db-%: ## Connect to service database (e.g., make db-auth)
	docker-compose exec $*-db psql -U postgres -d $*_db

# Service creation from template
new-service: ## Create new microservice from template (Usage: make new-service NAME=myservice PORT=8088)
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)Error: NAME is required$(NC)"; \
		echo "Usage: make new-service NAME=myservice PORT=8088"; \
		exit 1; \
	fi
	@if [ -z "$(PORT)" ]; then \
		echo "$(RED)Error: PORT is required$(NC)"; \
		echo "Usage: make new-service NAME=myservice PORT=8088"; \
		exit 1; \
	fi
	@echo "$(GREEN)Creating new service: $(NAME) on port $(PORT)$(NC)"
	@mkdir -p Services/$(NAME)
	@cp -r Template/Microservice/* Services/$(NAME)/
	@echo "$(YELLOW)Service created at: Services/$(NAME)$(NC)"
	@echo "$(YELLOW)Next steps:$(NC)"
	@echo "  1. Update Services/$(NAME)/.env.example with PORT=$(PORT)"
	@echo "  2. Add the service to docker-compose.yml"
	@echo "  3. Update services.yaml configuration"

# Monitoring and debugging
top: ## Show resource usage of containers
	docker stats $$(docker-compose ps -q)

networks: ## Show Docker networks
	docker network ls | grep area

volumes: ## Show Docker volumes
	docker volume ls | grep area

inspect-%: ## Inspect service container (e.g., make inspect-auth)
	docker-compose exec $*-service env

# Documentation
list-services: ## List all configured services
	@echo '$(GREEN)=== Configured Services ===$(NC)'
	@echo ''
	@cat services.yaml | grep -A 3 "  [a-z]*:" | grep -E "name:|port:" | paste - - | sed 's/^/  /' || echo '  No services configured yet'
