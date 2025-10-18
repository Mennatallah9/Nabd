# Nabd Backend Tests

This directory contains unit tests for the Nabd backend application. The tests are organized to mirror the main application structure and provide comprehensive coverage of the core functionality.

## Test Structure

```
tests/
├── controllers/           # Controller layer tests
│   ├── auth_controller_test.go
│   ├── container_controller_test.go
│   └── autoheal_controller_test.go
├── models/               # Model tests
│   └── models_test.go
├── services/             # Service layer tests 
│   ├── docker_service_test.go
│   └── metrics_service_test.go
└── utils/                # Utility function tests
    ├── auth_test.go
    ├── config_test.go
    └── database_test.go
```

## Dependencies

The tests use the following testing libraries:

- **testify**: Assertions and mocking (`github.com/stretchr/testify`)
  - `assert`: For test assertions
  - `require`: For test requirements
  - `mock`: For creating mocks

## Running Tests

### Quick Start
Run all tests:
```bash
cd backend
go test ./tests/...
```

### Using the Test Scripts

#### Unix/Linux/macOS (with coverage):
```bash
./run_tests.sh
```

#### Windows (with coverage):
```bash
.\run_tests.bat
```

#### or using the `test_runner.go` file (without coverage):
```
go run test_runner.go
```


### Using Make (if available)

View available test targets:
```bash
make -f Makefile.test help
```

Run tests with coverage:
```bash
make -f Makefile.test test-coverage
```

### Manual Commands

Run tests with verbose output:
```bash
go test -v ./tests/...
```

Run tests with coverage:
```bash
go test -v -race -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out -o coverage.html
```

Run specific test:
```bash
go test -v -run TestAuthController_Login ./tests/controllers/
```

Run tests for specific package:
```bash
go test -v ./tests/utils/
```

## Adding New Tests

When adding new tests:

1. Create test files with `_test.go` suffix
2. Place them in the appropriate subdirectory under `tests/`
3. Use the same package structure as the main code
4. Add comprehensive test cases covering:
   - Success scenarios
   - Error conditions
   - Edge cases
   - Input validation
