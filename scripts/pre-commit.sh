#!/bin/bash

# Pre-commit hook for Go project
# This script runs before each commit to ensure code quality

set -e

echo "🔍 Running pre-commit checks..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null; then
    print_error "golangci-lint is not installed!"
    echo "Please install it with:"
    echo "curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$(go env GOPATH)/bin v1.54.2"
    exit 1
fi

# 1. Format check (gofmt)
print_status "Running gofmt check..."
if ! gofmt -l . | grep -q .; then
    print_error "Code is not properly formatted!"
    echo "Files that need formatting:"
    gofmt -l .
    echo ""
    echo "Run 'gofmt -w .' to fix formatting issues"
    exit 1
fi
print_status "gofmt check passed"

# 2. Lint check (golangci-lint)
print_status "Running golangci-lint..."
if ! golangci-lint run; then
    print_error "Linting failed!"
    echo "Please fix the linting issues before committing"
    exit 1
fi
print_status "golangci-lint check passed"

# 3. Test check
print_status "Running tests..."
if ! go test ./... -v; then
    print_error "Tests failed!"
    echo "Please fix the failing tests before committing"
    exit 1
fi
print_status "All tests passed"

# 4. Build check
print_status "Running build check..."
if ! go build ./cmd/main.go; then
    print_error "Build failed!"
    echo "Please fix the build issues before committing"
    exit 1
fi
print_status "Build check passed"

# Clean up build artifact
rm -f main

print_status "All pre-commit checks passed! 🎉"
echo "You can now commit your changes."
