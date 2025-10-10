#!/bin/bash

# Nabd Build Script
# This script builds the complete Nabd application

set -e

echo "🔹 Building Nabd v0.1.0..."

# Build frontend
echo "📦 Building React frontend..."
cd frontend
npm install
npm run build
cd ..

# Build backend
echo "🚀 Building Go backend..."
cd backend
go mod download
CGO_ENABLED=1 go build -o nabd main.go
cd ..

echo "✅ Build completed successfully!"
echo ""
echo "To run locally:"
echo "  cd backend && ./nabd"
echo ""
echo "To build Docker image:"
echo "  docker build -t nabd:v0.1.0 ."
echo ""
echo "To run with Docker Compose:"
echo "  docker-compose up -d"