#!/bin/bash

# Setup script for pre-commit hooks
# Run this script to setup pre-commit hooks for the project

set -e

echo "🔧 Setting up pre-commit hooks for Go project..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed!"
    echo "Please install Go first: https://golang.org/dl/"
    exit 1
fi

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null; then
    print_warning "golangci-lint is not installed. Installing..."
    
    # Detect shell and use appropriate syntax
    if [[ "$SHELL" == *"fish"* ]]; then
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b (go env GOPATH)/bin v1.54.2
    else
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
    fi
    
    print_status "golangci-lint installed successfully"
else
    print_status "golangci-lint is already installed"
fi

# Setup pre-commit hook
if [ -f "scripts/pre-commit.sh" ]; then
    cp scripts/pre-commit.sh .git/hooks/pre-commit
    chmod +x .git/hooks/pre-commit
    print_status "Pre-commit hook installed successfully"
else
    print_error "scripts/pre-commit.sh not found!"
    exit 1
fi

# Test the setup
print_status "Testing pre-commit setup..."
if [ -x ".git/hooks/pre-commit" ]; then
    print_status "Pre-commit hook is executable"
else
    print_error "Pre-commit hook is not executable!"
    exit 1
fi

echo ""
print_status "Setup completed! 🎉"
echo ""
echo "Pre-commit hooks will now run automatically before each commit."
echo "The hooks will check:"
echo "  - Code formatting (gofmt)"
echo "  - Code linting (golangci-lint)"
echo "  - Tests (go test)"
echo "  - Build (go build)"
echo ""
echo "If any check fails, the commit will be blocked."
echo "Fix the issues and try committing again."
