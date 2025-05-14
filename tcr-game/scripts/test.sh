# scripts/test.sh - Test script
#!/bin/bash

set -e

echo "Running tests..."

# Run unit tests
echo "Running unit tests..."
go test ./internal/... -v

# Run integration tests
echo "Running integration tests..."
go test ./tests/integration/... -v

# Check test coverage
echo "Checking test coverage..."
go test ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

echo "Tests complete! Coverage report in coverage.html"