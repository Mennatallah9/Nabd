@echo off
REM Nabd Build Script for Windows
REM This script builds the complete Nabd application

echo 🔹 Building Nabd v0.1.0...

REM Build frontend
echo 📦 Building React frontend...
cd frontend
call npm install
call npm run build
cd ..

REM Build backend
echo 🚀 Building Go backend...
cd backend
go mod download
set CGO_ENABLED=1
go build -o nabd.exe main.go
cd ..

echo ✅ Build completed successfully!
echo.
echo To run locally:
echo   cd backend ^&^& nabd.exe
echo.
echo To build Docker image:
echo   docker build -t nabd:v0.1.0 .
echo.
echo To run with Docker Compose:
echo   docker-compose up -d