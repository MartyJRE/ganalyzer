.PHONY: build test clean run install fmt vet lint coverage help

# Binary name and paths
BINARY_NAME=ganalyzer
BUILD_DIR=build
CMD_PATH=cmd/ganalyzer

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Version information
VERSION ?= dev
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS=-ldflags "-X ganalyzer/internal/version.Version=$(VERSION) -X ganalyzer/internal/version.GitCommit=$(GIT_COMMIT) -X ganalyzer/internal/version.BuildDate=$(BUILD_DATE) -s -w"

# Default target
all: test build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_PATH)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	$(GOTEST) -race ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Run the application
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Install the application to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) ./$(CMD_PATH)

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Run linting (requires golangci-lint)
lint:
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GOCMD) mod tidy

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOCMD) mod download

# Development workflow - format, vet, test, build
dev: fmt vet test build

# CI workflow - all checks
ci: fmt vet lint test-race build

# Help
help:
	@echo "Available targets:"
	@echo "  build      - Build the application"
	@echo "  test       - Run tests"
	@echo "  coverage   - Run tests with coverage report"
	@echo "  test-race  - Run tests with race detection"
	@echo "  run        - Build and run the application"
	@echo "  install    - Install to GOPATH/bin"
	@echo "  clean      - Clean build artifacts"
	@echo "  fmt        - Format code"
	@echo "  vet        - Run go vet"
	@echo "  lint       - Run golangci-lint"
	@echo "  tidy       - Tidy dependencies"
	@echo "  deps       - Download dependencies"
	@echo "  dev        - Development workflow (fmt, vet, test, build)"
	@echo "  ci         - CI workflow (fmt, vet, lint, test-race, build)"
	@echo "  help       - Show this help"