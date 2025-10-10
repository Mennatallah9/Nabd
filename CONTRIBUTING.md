# Contributing to Nabd

Thank you for your interest in contributing to Nabd! We welcome contributions from the community and are pleased to have you join us.

## Table of Contents
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Pull Request Process](#pull-request-process)
- [Issue Guidelines](#issue-guidelines)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Set up the development environment
4. Create a new branch for your feature or bugfix
5. Make your changes
6. Test your changes thoroughly
7. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- Docker and Docker Compose
- Git

### Local Development

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Mennatallah9/Nabd.git
   cd nabd
   ```

2. **Build the complete application:**
   
   **For Linux/macOS:**
   ```bash
   ./build.sh
   ```
   
   **For Windows:**
   ```batch
   build.bat
   ```
   
   This will automatically:
   - Install frontend dependencies
   - Build the React frontend
   - Download Go dependencies
   - Build the Go backend

3. **For development with live reload:**
   
   **Backend development:**
   ```bash
   cd backend
   go mod download
   go run main.go
   ```

   **Frontend development (in a new terminal):**
   ```bash
   cd frontend
   npm install
   npm start
   ```

4. **Access the application:**
   - Backend API: http://localhost:8080
   - Frontend: http://localhost:3000 (development) or served by backend (production)
   - Default admin token: `nabd-admin-token`

### Docker Development

```bash
# Build and run with Docker Compose
docker-compose up -d

# Build individual images
docker build -t nabd:dev .
```

## Making Changes

### Branch Naming

Use descriptive branch names:
- `feature/add-kubernetes-support`
- `bugfix/fix-memory-leak`
- `docs/update-api-documentation`
- `refactor/improve-error-handling`

### Commit Messages

Write clear, descriptive commit messages:
- Use present tense ("Add feature" not "Added feature")
- Use imperative mood ("Move cursor to..." not "Moves cursor to...")


## Pull Request Process

1. **Before submitting:**
   - Ensure your branch is up to date with the main branch
   - Run all tests and ensure they pass
   - Update documentation as needed
   - Test your changes in both development and Docker environments

2. **Pull Request Requirements:**
   - Use the pull request template
   - Provide a clear description of the changes
   - Include screenshots for UI changes
   - Reference any related issues
   - Ensure CI checks pass

3. **Review Process:**
   - Maintainers will review your pull request
   - Address any requested changes promptly
   - Be responsive to feedback and questions
   - Once approved, a maintainer will merge your changes

## Issue Guidelines

### Before Creating an Issue

1. Search existing issues to avoid duplicates
2. Check the documentation for answers
3. Try the latest version to see if the issue still exists


## Testing

### Pre-Testing Setup

Before running tests, ensure the application builds correctly:

**Linux/macOS:**
```bash
./build.sh
```

**Windows:**
```batch
build.bat
```

This ensures all dependencies are properly installed and the application compiles successfully.

### Backend Tests

```bash
cd backend
go test ./...
go test -race ./...  # Test for race conditions
go test -cover ./... # Test coverage
```

### Frontend Tests

```bash
cd frontend
npm test
npm run test:coverage
```

### Integration Tests

```bash
# Start test environment
docker-compose -f docker-compose.test.yml up -d

# Run integration tests
go test -tags integration ./tests/...
```

## Recognition

Contributors are recognized in our:
- README.md contributors section
- Release notes
- Project documentation

Thank you for contributing to Nabd! Your efforts help make container monitoring and auto-healing accessible to everyone.