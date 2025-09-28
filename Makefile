# Makefile for URL Shortener Project

.PHONY: help build test lint format clean setup-hooks run docker-build docker-run

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development commands
build: ## Build the application
	@echo "🔨 Building application..."
	go build -o main ./cmd/main.go
	@echo "✅ Build completed!"

test: ## Run all tests
	@echo "🧪 Running tests..."
	go test ./... -v
	@echo "✅ Tests completed!"

lint: ## Run golangci-lint
	@echo "🔍 Running linter..."
	golangci-lint run
	@echo "✅ Linting completed!"

format: ## Format code with gofmt
	@echo "🎨 Formatting code..."
	gofmt -w .
	@echo "✅ Code formatted!"

clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	rm -f main
	go clean -cache
	@echo "✅ Clean completed!"

# Setup commands
setup-hooks: ## Setup pre-commit hooks
	@echo "🔧 Setting up pre-commit hooks..."
	./scripts/setup-hooks.sh
	@echo "✅ Pre-commit hooks setup completed!"

# Pre-commit checks
pre-commit: ## Run pre-commit checks manually
	@echo "🔍 Running pre-commit checks..."
	./scripts/pre-commit.sh
	@echo "✅ Pre-commit checks completed!"

# Run application
run: ## Run the application
	@echo "🚀 Starting application..."
	go run ./cmd/main.go

# Docker commands
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -f Dockerfile.golang -t url-shortener-base:latest .
	docker build -f Dockerfile -t url-shortener:latest .
	@echo "✅ Docker build completed!"

docker-run: ## Run with Docker Compose
	@echo "🐳 Starting with Docker Compose..."
	docker-compose up --build
	@echo "✅ Docker Compose started!"

docker-stop: ## Stop Docker Compose
	@echo "🛑 Stopping Docker Compose..."
	docker-compose down
	@echo "✅ Docker Compose stopped!"

# Development workflow
dev-setup: setup-hooks ## Complete development setup
	@echo "🎯 Development setup completed!"
	@echo "You can now start developing with:"
	@echo "  make run          # Run the application"
	@echo "  make test         # Run tests"
	@echo "  make lint         # Run linter"
	@echo "  make format       # Format code"

# CI/CD simulation
ci: format lint test build ## Run CI pipeline locally
	@echo "✅ CI pipeline completed successfully!"

# Quick development cycle
dev: format lint test ## Quick development cycle (format + lint + test)
	@echo "✅ Development cycle completed!"

# Install dependencies
deps: ## Install Go dependencies
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies installed!"

# Update dependencies
update-deps: ## Update Go dependencies
	@echo "🔄 Updating dependencies..."
	go get -u ./...
	go mod tidy
	@echo "✅ Dependencies updated!"

# Security check
security: ## Run security checks
	@echo "🔒 Running security checks..."
	golangci-lint run --enable gosec
	@echo "✅ Security checks completed!"

# Performance check
perf: ## Run performance checks
	@echo "⚡ Running performance checks..."
	go test -bench=. ./...
	@echo "✅ Performance checks completed!"

# Coverage
coverage: ## Run tests with coverage
	@echo "📊 Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# All checks
check-all: format lint test build security ## Run all checks
	@echo "✅ All checks completed successfully!"
