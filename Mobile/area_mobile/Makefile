.PHONY: help deps clean docker-build docker-run docker-clean docker-shell

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

docker-shell: ## Drop into a shell inside a temp container (debug)
	docker run --rm -it \
		-e BASE_URL=$(BASE_URL) \
		-v $(APK_OUT):/apk \
		$(IMAGE_NAME) /bin/bash

docker-clean: ## Remove the built Docker image
	docker rmi $(IMAGE_NAME) || true
