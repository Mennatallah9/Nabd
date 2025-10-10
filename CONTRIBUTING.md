# Contributing to Nabd

Thank you for your interest in contributing to Nabd! We welcome contributions from the community and are pleased to have you join us.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Pull Request Process](#pull-request-process)
- [Issue Guidelines](#issue-guidelines)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code. Please be respectful and constructive in all interactions.

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
   git clone https://github.com/your-username/nabd.git
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

### Building for Production

The project includes automated build scripts that handle both frontend and backend compilation:

**For Linux/macOS:**
```bash
./build.sh
```

**For Windows:**
```batch
build.bat
```

These scripts will:
1. Install and build the React frontend (`npm install && npm run build`)
2. Download Go dependencies (`go mod download`)
3. Compile the backend binary (`go build -o nabd main.go`)
4. Provide instructions for running or containerizing the application

**Manual Production Build (if needed):**
```bash
# Build frontend
cd frontend
npm install
npm run build
cd ..

# Build backend
cd backend
go mod download
go build -o nabd main.go
cd ..

# The backend will serve the frontend static files from frontend/build/
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
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

Example:
```
Add container restart retry mechanism

- Implement exponential backoff for failed restarts
- Add configurable retry limits
- Update auto-heal service to use new retry logic

Fixes #123
```

### Code Changes

1. **Make small, focused changes:** Each pull request should address a single concern
2. **Write tests:** Include unit tests for new functionality
3. **Update documentation:** Keep README and code comments up to date
4. **Follow coding standards:** Maintain consistency with existing code
5. **Test thoroughly:** Ensure your changes don't break existing functionality

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

### Bug Reports

Include the following information:
- Nabd version
- Operating system and version
- Docker version
- Go version (if building from source)
- Steps to reproduce the issue
- Expected vs actual behavior
- Relevant logs or error messages
- Screenshots if applicable

### Feature Requests

Provide:
- Clear description of the proposed feature
- Use case and motivation
- Proposed implementation approach (if applicable)
- Any alternative solutions considered

### Questions and Support

For questions:
- Check the documentation first
- Search existing issues and discussions
- Use GitHub Discussions for general questions
- Use Issues for specific bugs or feature requests

## Coding Standards

### Go Backend

- Follow Go conventions and idioms
- Use `gofmt` for code formatting
- Use meaningful variable and function names
- Include error handling for all operations
- Write unit tests for business logic
- Document exported functions and types
- Use structured logging

Example:
```go
// GetContainerMetrics retrieves current metrics for all running containers
func (s *MetricsService) GetContainerMetrics() ([]models.ContainerMetric, error) {
    containers, err := s.dockerService.GetContainers()
    if err != nil {
        return nil, fmt.Errorf("failed to get containers: %w", err)
    }
    // ... implementation
}
```

### React Frontend

- Use functional components with hooks
- Follow React best practices
- Use TypeScript for new components (gradual migration)
- Implement proper error boundaries
- Use consistent naming conventions
- Write meaningful component tests

Example:
```javascript
const ContainerCard = ({ container, metric, onRestart }) => {
  const [loading, setLoading] = useState(false);
  
  const handleRestart = async () => {
    setLoading(true);
    try {
      await onRestart(container.name);
    } catch (error) {
      console.error('Failed to restart container:', error);
    } finally {
      setLoading(false);
    }
  };
  
  return (
    // JSX implementation
  );
};
```

### Database

- Use prepared statements for SQL queries
- Include proper indexes for performance
- Write migration scripts for schema changes
- Document database schema changes

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

### Manual Testing

1. Test with real Docker containers
2. Verify auto-healing functionality
3. Test alert thresholds
4. Validate UI responsiveness
5. Test authentication flows

## Documentation

### Code Documentation

- Document all exported functions and types
- Include usage examples for complex functionality
- Keep comments up to date with code changes
- Use clear, concise language

### User Documentation

- Update README.md for user-facing changes
- Include configuration examples
- Document new API endpoints
- Provide troubleshooting guides

### API Documentation

- Document all REST endpoints
- Include request/response examples
- Specify authentication requirements
- Document error codes and messages

## Release Process

### Version Numbering

We use Semantic Versioning (SemVer):
- MAJOR version for incompatible API changes
- MINOR version for backwards-compatible functionality
- PATCH version for backwards-compatible bug fixes

### Release Checklist

1. Update version numbers
2. Update CHANGELOG.md
3. Test release candidate
4. Create release tag
5. Build and publish Docker images
6. Update documentation

## Getting Help

- **Documentation:** Check the README and wiki
- **Issues:** Search existing issues before creating new ones
- **Discussions:** Use GitHub Discussions for general questions
- **Code Review:** Participate in pull request reviews
- **Community:** Join our community discussions

## Recognition

Contributors are recognized in our:
- README.md contributors section
- Release notes
- Project documentation

Thank you for contributing to Nabd! Your efforts help make container monitoring and auto-healing accessible to everyone.