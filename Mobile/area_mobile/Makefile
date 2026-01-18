.PHONY: help build run test clean dev init deps fmt lint docker-build docker-up docker-down docker-logs docker-restart up docker-run docker-clean docker-shell

SHELL := /bin/bash
MAKEFLAGS += --no-print-directory

IMAGE_NAME ?= area-mobile
BASE_URL ?= http://gateway:8080
APK_OUT ?= $(CURDIR)/dist

help: ## Show this help message
	@printf 'Usage: make [target]\n\n'
	@printf 'Available targets:\n'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-18s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Fetch Flutter dependencies locally
	flutter pub get

init: deps ## Initialize local environment

build: docker-run ## Build APK via Docker

run: build ## Alias for build

dev: ## Dev run (manual)
	@echo "No dev target configured; run 'flutter run' manually if needed."

test: ## Run tests (not configured)
	@echo "No test target configured for mobile."

fmt: ## Format code (not configured)
	@echo "No formatter configured for mobile."

lint: ## Run lints (not configured)
	@echo "No linter configured for mobile."

clean: ## Remove build artifacts and local APK output
	rm -rf build $(APK_OUT)

docker-build: ## Build the mobile Docker image
	docker build -t $(IMAGE_NAME) --build-arg BASE_URL=$(BASE_URL) .

docker-run: ## Build & run container to produce APK in dist/client.apk then exit
	mkdir -p $(APK_OUT)
	$(MAKE) docker-build
	docker run --rm \
		-e BASE_URL=$(BASE_URL) \
		-e APK_OUTPUT_PATH=/apk/client.apk \
		-e EXIT_AFTER_BUILD=true \
		-v $(APK_OUT):/apk \
		$(IMAGE_NAME)

docker-up: docker-run ## Alias for docker-run

docker-down: ## No long-running containers to stop
	@echo "No containers to stop (docker-run is one-shot)."

docker-logs: ## No persistent logs to tail
	@echo "No logs to tail (docker-run is one-shot)."

docker-restart: docker-run ## Alias for docker-run

up: docker-up ## Alias for docker-up

docker-shell: ## Drop into a shell inside a temp container (debug)
	docker run --rm -it \
		-e BASE_URL=$(BASE_URL) \
		-v $(APK_OUT):/apk \
		$(IMAGE_NAME) /bin/bash

docker-clean: ## Remove the built Docker image
	docker rmi $(IMAGE_NAME) || true
