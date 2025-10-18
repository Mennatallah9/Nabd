@echo off
REM Test runner script for Nabd backend (Windows)

echo Running Nabd Backend Tests...
echo ==============================

cd /d "%~dp0"

REM Install dependencies
echo Installing test dependencies...
go mod tidy

REM Run all tests with coverage
echo Running tests with coverage...
go test -v -race -coverprofile=coverage.out ./tests/...

REM Generate coverage report
echo Generating coverage report...
go tool cover -html=coverage.out -o coverage.html

REM Display coverage summary
echo Coverage summary:
go tool cover -func=coverage.out

echo.
echo Test run completed!
echo Coverage report generated: coverage.html
pause