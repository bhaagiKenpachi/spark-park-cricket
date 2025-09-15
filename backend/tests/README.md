# Testing Guide

This document describes how to run tests for the Spark Park Cricket backend.

## Test Structure

The test suite is organized into two main categories:

- **Unit Tests** (`tests/unit/`): Test individual components in isolation using mocks
- **Integration Tests** (`tests/integration/`): Test the complete application flow with real database

## Running Tests

### Unit Tests

Unit tests can be run without any external dependencies:

```bash
# Run all unit tests
go test ./tests/unit/... -v

# Run specific test file
go test ./tests/unit/series_service_unit_test.go -v

# Run with coverage
go test ./tests/unit/... -cover
```

### Integration Tests

Integration tests require a properly configured Supabase database:

1. **Set up environment variables** in `.env` file:
   ```
   SUPABASE_URL=your_supabase_project_url
   SUPABASE_API_KEY=your_supabase_anon_key
   ```

2. **Run database migrations** (if not already done):
   ```bash
   go run cmd/migrate/main.go
   ```

3. **Run integration tests**:
   ```bash
   go test ./tests/integration/... -v
   ```

**Note**: Integration tests will be skipped if Supabase credentials are not configured.

## Known Issues

### UUID Generation Issue

There is currently a known issue with UUID generation when creating new records through the Supabase client. The error appears as:

```
invalid input syntax for type uuid: ""
```

**Workaround**: This issue affects the API endpoints but not the test suite, as unit tests use mocks and integration tests validate the expected behavior.

**Status**: This is likely due to the Supabase Go client version (v0.0.4) and may be resolved by upgrading to a newer version when available.

## Test Examples

### Series Creation Test

The integration test demonstrates:
- Valid series creation request
- Validation error handling for invalid data
- Date validation (end date must be after start date)
- List operations

### Unit Test Coverage

Unit tests cover:
- Service layer business logic
- Input validation
- Error handling
- Repository interaction patterns

## CI/CD Considerations

For GitHub Actions or other CI/CD systems:

1. **Unit tests** can run in any environment
2. **Integration tests** require:
   - Supabase project setup
   - Environment variables configuration
   - Database migration execution

Example GitHub Actions workflow:

```yaml
name: Test
on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.23'
      - run: go test ./tests/unit/... -v

  integration-tests:
    runs-on: ubuntu-latest
    needs: unit-tests
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.23'
      - name: Run migrations
        env:
          SUPABASE_URL: ${{ secrets.SUPABASE_URL }}
          SUPABASE_API_KEY: ${{ secrets.SUPABASE_API_KEY }}
        run: go run cmd/migrate/main.go
      - name: Run integration tests
        env:
          SUPABASE_URL: ${{ secrets.SUPABASE_URL }}
          SUPABASE_API_KEY: ${{ secrets.SUPABASE_API_KEY }}
        run: go test ./tests/integration/... -v
```

## Adding New Tests

### Unit Tests

1. Create test files in `tests/unit/`
2. Use mocks for external dependencies
3. Follow the naming convention: `*_test.go`
4. Use testify for assertions and mocks

### Integration Tests

1. Create test files in `tests/integration/`
2. Use real database connections
3. Clean up test data when possible
4. Skip tests gracefully when dependencies are not available

## Running Specific Test Cases

```bash
# Run specific test function
go test ./tests/unit/ -run TestSeriesService_CreateSeries -v

# Run tests matching pattern
go test ./tests/... -run ".*Series.*" -v

# Run with timeout
go test ./tests/integration/... -timeout 30s -v
```
