#!/bin/bash

# Nabd Build Script
# This script builds the complete Nabd application

set -e

echo "Building Nabdﮩ٨ـﮩﮩ٨ـ♡ﮩ٨ـﮩﮩ٨ـ..."

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
echo ""
echo "To run with Docker Compose:"
echo "  docker-compose up -d"