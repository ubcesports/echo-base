.PHONY: test test-unit test-integration test-coverage test-db-up test-db-down test-migrate

# Run all tests
test: test-unit test-integration

# Run unit tests
test-unit: 
	go test -v ./internal/...

# Run integration tests
test-integration: 
	go test -v ./cmd/...

# Run tests with coverage
test-coverage: 
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Start test database
test-db-up:
	docker compose -f compose.test.yaml up -d postgres-test
	sleep 5  # Wait for database to be ready
	go tool sql-migrate up -env=test

# Stop test database
test-db-down:
	docker compose -f compose.test.yaml down

# Clean test artifacts
test-clean:
	rm -f coverage.out coverage.html
	docker compose -f compose.test.yaml down -v
