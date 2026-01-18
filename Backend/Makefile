# Global Makefile for Backend
# Runs common targets across subprojects

.PHONY: help build run test test_run clean docker-build docker-up docker-down docker-logs docker-restart dev init deps fmt lint

SHELL := /bin/bash
MAKEFLAGS += --no-print-directory

# Subprojects participating in global commands (relative to Backend/)
PROJECTS := Gateway Services/AuthService Services/ServiceService Services/WebhookService Services/AreaService Services/PollingService Services/MailService Services/CronService

# Colors
RESET  := \033[0m
BOLD   := \033[1m
DIM    := \033[2m
BLUE   := \033[34m
GREEN  := \033[32m
YELLOW := \033[33m
RED    := \033[31m
CYAN   := \033[36m

# Helper to run a target for all subprojects
define run_for_all
	@set -e; \
	for dir in $(PROJECTS); do \
		printf "\n$(BOLD)$(BLUE)==> %s$(RESET) $(DIM)make$(RESET) $(YELLOW)%s$(RESET)\n" "$$dir" "$(1)"; \
		$(MAKE) -C "$$dir" "$(1)"; \
	done
endef

help: ## Show this help message and subproject helps
	@printf '$(BOLD)Usage:$(RESET) make [target]\n'
	@echo ''
	@printf '$(BOLD)Available targets:$(RESET)\n'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ''
	@printf '$(BOLD)Subproject helps:$(RESET)\n'
	$(call run_for_all,help)

build: ## Build all subprojects
	$(call run_for_all,build)

run: ## Run all subprojects locally (may block)
	$(call run_for_all,run)

test: ## Run tests for all subprojects
	$(call run_for_all,test)

test_run: ## Run unit tests for all subprojects
	$(call run_for_all,test_run)

test_summary: ## Run unit tests with clean summary output
	@$(MAKE) test_run 2>&1 | ./scripts/test_summary.sh

clean: ## Clean build artifacts in all subprojects
	$(call run_for_all,clean)

docker-build: ## Build Docker images for all subprojects
	$(call run_for_all,docker-build)

docker-up: ## Start Docker for all subprojects
	$(call run_for_all,docker-up)

docker-down: ## Stop Docker for all subprojects
	$(call run_for_all,docker-down)

docker-logs: ## Tail Docker logs for all subprojects (blocks)
	$(call run_for_all,docker-logs)

docker-restart: ## Restart Docker for all subprojects
	$(call run_for_all,docker-restart)

dev: ## Run dev mode with hot reload for all subprojects (may block)
	$(call run_for_all,dev)

init: ## Initialize environments/configs for all subprojects
	$(call run_for_all,init)

deps: ## Download Go dependencies for all subprojects
	$(call run_for_all,deps)

fmt: ## Format Go code for all subprojects
	$(call run_for_all,fmt)

lint: ## Run linter for all subprojects
	$(call run_for_all,lint)
