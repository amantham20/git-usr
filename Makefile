.PHONY: build test test-unit test-integration clean install help

# Build the binary
build:
	@echo "Building git-usr..."
	go build -o git-usr main.go

# Run all tests
test: test-unit test-integration

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic

# Run integration tests (requires git to be installed)
test-integration:
	@echo "Running integration tests..."
	go test -v -tags=integration

# Run tests with coverage report
coverage: test-unit
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -f git-usr git-usr.exe coverage.out coverage.html

# Install the binary
install: build
	@echo "Installing git-usr..."
	@if [ "$(shell uname)" = "Darwin" ]; then \
		cp git-usr /usr/local/bin/git-usr; \
	else \
		mkdir -p $(HOME)/.local/bin; \
		cp git-usr $(HOME)/.local/bin/git-usr; \
	fi
	@echo "✅ git-usr installed successfully!"

# Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem

# Show help
help:
	@echo "Available targets:"
	@echo "  build              - Build the git-usr binary"
	@echo "  test               - Run all tests (unit + integration)"
	@echo "  test-unit          - Run unit tests only"
	@echo "  test-integration   - Run integration tests (requires git)"
	@echo "  coverage           - Generate test coverage report"
	@echo "  clean              - Remove build artifacts"
	@echo "  install            - Build and install git-usr"
	@echo "  lint               - Run linter (requires golangci-lint)"
	@echo "  fmt                - Format code"
	@echo "  bench              - Run benchmarks"
	@echo "  help               - Show this help message"
