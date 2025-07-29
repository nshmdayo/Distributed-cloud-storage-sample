.PHONY: help setup build run-node run-api test clean docker-build docker-run

# Variables
BINARY_NAME=dcs
BUILD_DIR=bin
GO_FILES=$(shell find . -name "*.go" -type f)

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## Setup development environment
	@echo "Setting up development environment..."
	go mod tidy
	go mod download
	mkdir -p $(BUILD_DIR)
	mkdir -p data/storage
	mkdir -p logs

build: ## Build all binaries
	@echo "Building binaries..."
	go build -o $(BUILD_DIR)/node ./cmd/node
	go build -o $(BUILD_DIR)/client ./cmd/client
	go build -o $(BUILD_DIR)/api ./cmd/api

run-node: build ## Run storage node
	@echo "Starting storage node..."
	./$(BUILD_DIR)/node

run-api: build ## Run API server
	@echo "Starting API server..."
	./$(BUILD_DIR)/api

run-client: build ## Run client CLI
	@echo "Running client..."
	./$(BUILD_DIR)/client

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	gofmt -s -w .
	go mod tidy

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -rf data/storage/*
	rm -f coverage.out coverage.html
	rm -f *.log

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t distributed-cloud-storage .

docker-run: docker-build ## Run in Docker
	@echo "Running in Docker..."
	docker-compose up

install-deps: ## Install development dependencies
	@echo "Installing development dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
