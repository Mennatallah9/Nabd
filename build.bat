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

REM Copy frontend build files to backend static directory
echo Copying frontend files...
if exist static rmdir /s /q static
mkdir static
xcopy /e /i ..\frontend\build\static\* static\
copy ..\frontend\build\index.html static\
copy ..\frontend\build\*.png static\ 2>nul
copy ..\frontend\build\*.ico static\ 2>nul

cd ..

echo ✅ Build completed successfully!
echo.
echo To run locally:
echo   cd backend ^&^& nabd.exe