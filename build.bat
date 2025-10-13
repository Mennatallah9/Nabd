@echo off
REM Nabd Build Script for Windows
REM This script builds the complete Nabd application

echo Building Nabd...

REM Build frontend
echo Building frontend...
cd frontend
call npm install
call npm run build
cd ..

REM Build backend
echo Building backend...
cd backend
go mod download
set CGO_ENABLED=1
go build -o nabd.exe main.go
cd ..

echo ✅ Build completed successfully!
echo.
echo To run locally:
echo   cd backend ^&^& nabd.exe