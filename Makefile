.PHONY: build run migrate lint

# Go parameters
GOBASE := $(shell pwd)
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin
GOFILES := $(wildcard *.go)

# Variables
APP_NAME := micomido
MIGRATOR_NAME := migrator
BINARY_DIR := ./bin

# Default target
all: build

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BINARY_DIR)/$(APP_NAME) main.go
	@go build -o $(BINARY_DIR)/$(MIGRATOR_NAME) cmd/migrator/main.go

# Run the application
run:
	@echo "Starting $(APP_NAME)..."
	@DB_CONNECTION_STRING="user:password@tcp(localhost:3306)/micomido?parseTime=true" \
	S3_BUCKET="micomido" \
	AWS_ACCESS_KEY_ID="minioadmin" \
	AWS_SECRET_ACCESS_KEY="minioadmin" \
	AWS_REGION="us-east-1" \
	AWS_ENDPOINT_URL="http://localhost:9000" \
	./$(BINARY_DIR)/$(APP_NAME)

# Run database migrations
migrate:
	@echo "Running database migrations..."
	@DB_CONNECTION_STRING="user:password@tcp(localhost:3306)/micomido" \
	./$(BINARY_DIR)/$(MIGRATOR_NAME)

# Lint the code
lint:
	@echo "Linting code..."
	@golangci-lint run

# Clean up binaries
clean:
	@echo "Cleaning up..."
	@rm -rf $(BINARY_DIR)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy

# Install golangci-lint
install-lint:
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
