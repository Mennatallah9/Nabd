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

# Copy frontend build files to backend static directory
echo "Copying frontend files..."
if [ -d "static" ]; then
    rm -rf static
fi
mkdir -p static
cp -r ../frontend/build/static/* static/
cp ../frontend/build/index.html static/
cp ../frontend/build/*.png static/ 2>/dev/null || true
cp ../frontend/build/*.ico static/ 2>/dev/null || true

cd ..

echo "✅ Build completed successfully!"
echo ""
echo "To run locally:"
echo "  cd backend && ./nabd"