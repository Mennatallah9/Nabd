#!/bin/bash

# Nabd Build Script
# This script builds the complete Nabd application

set -e

echo "Building Nabd..."

# Build frontend
echo "Building frontend..."
cd frontend
npm install
npm run build
cd ..

# Build backend
echo "Building backend..."
cd backend
go mod download
CGO_ENABLED=1 go build -o nabd main.go
cd ..

echo "✅ Build completed successfully!"
echo ""
echo "To run locally:"
echo "  cd backend && ./nabd"