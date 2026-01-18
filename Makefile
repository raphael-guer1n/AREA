# Global Makefile for Frontend Docker
# Runs Docker targets for the Web frontend and mobile app

.PHONY: help docker-build docker-up docker-down docker-logs docker-restart up docker-mobile-build docker-mobile-run docker-mobile-shell docker-mobile-clean

SHELL := /bin/bash
MAKEFLAGS += --no-print-directory

# Subprojects participating in global commands (relative to repo root)
PROJECTS_FRONTEND := Web/frontend
PROJECTS_MOBILE := Mobile/area_mobile
PROJECTS_ALL := $(PROJECTS_FRONTEND) $(PROJECTS_MOBILE)
PROJECTS := $(PROJECTS_FRONTEND)

# Colors
RESET  := \033[0m
BOLD   := \033[1m
DIM    := \033[2m
BLUE   := \033[34m
YELLOW := \033[33m
CYAN   := \033[36m

# Helper to run a target for all subprojects
define run_for_all
	@set -e; \
	for dir in $(PROJECTS); do \
		printf "\n$(BOLD)$(BLUE)==> %s$(RESET) $(DIM)make$(RESET) $(YELLOW)%s$(RESET)\n" "$$dir" "$(1)"; \
		$(MAKE) -C "$$dir" "$(1)"; \
	done
endef

help: ## Show this help message
	@printf '$(BOLD)Usage:$(RESET) make [target]\n'
	@echo ''
	@printf '$(BOLD)Available targets:$(RESET)\n'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

docker-build: PROJECTS := $(PROJECTS_ALL)
docker-build: ## Build Docker images for the frontend and mobile app
	$(call run_for_all,docker-build)

docker-up: ## Start Docker containers for the frontend
	$(call run_for_all,docker-up)

up: docker-up ## Alias for docker-up

docker-down: ## Stop Docker containers for the frontend
	$(call run_for_all,docker-down)

docker-logs: ## Tail Docker logs for the frontend
	$(call run_for_all,docker-logs)

docker-restart: ## Restart Docker containers for the frontend
	$(call run_for_all,docker-restart)

docker-mobile-build: PROJECTS := $(PROJECTS_MOBILE)
docker-mobile-build: ## Build Docker image for the mobile app
	$(call run_for_all,docker-build)

docker-mobile-run: PROJECTS := $(PROJECTS_MOBILE)
docker-mobile-run: ## Build & run the mobile container to produce an APK then exit
	$(call run_for_all,docker-run)

docker-mobile-shell: PROJECTS := $(PROJECTS_MOBILE)
docker-mobile-shell: ## Drop into a shell in the mobile container
	$(call run_for_all,docker-shell)

docker-mobile-clean: PROJECTS := $(PROJECTS_MOBILE)
docker-mobile-clean: ## Remove the mobile Docker image
	$(call run_for_all,docker-clean)
