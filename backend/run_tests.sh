#!/bin/bash

# Test runner script for Nabd backend

echo "Running Nabd Backend Tests..."
echo "=============================="

cd "$(dirname "$0")"

# Install dependencies
echo "Installing test dependencies..."
go mod tidy

# Run all tests with coverage
echo "Running tests with coverage..."
go test -v -race -coverprofile=coverage.out ./tests/...

# Generate coverage report
echo "Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html

# Display coverage summary
echo "Coverage summary:"
go tool cover -func=coverage.out

echo ""
echo "Test run completed!"
echo "Coverage report generated: coverage.html"