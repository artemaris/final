.PHONY: build-server build-client build test clean run-server run-client coverage html-coverage

# Build variables
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Build flags
LDFLAGS := -X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE)

build-server:
	@echo "Building server..."
	@go build -ldflags "$(LDFLAGS)" -o bin/gophkeeper-server ./cmd/server

build-client:
	@echo "Building client..."
	@go build -ldflags "$(LDFLAGS)" -o bin/gophkeeper-client ./cmd/client

build: build-server build-client
	@echo "Build complete!"

test:
	@echo "Running tests..."
	@go test -v -cover ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

html-coverage:
	@echo "Generating HTML coverage report..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@echo "Open coverage.html in your browser to view"

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

run-server:
	@echo "Starting server..."
	@go run ./cmd/server

run-client:
	@echo "Running client..."
	@go run ./cmd/client

docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down

docker-logs:
	@docker-compose logs -f

swagger-ui:
	@echo "Starting Swagger UI..."
	@docker run -d -p 8081:8080 -e SWAGGER_JSON=/docs/openapi.yaml -v $(PWD)/docs:/docs --name swagger-editor swaggerapi/swagger-editor
	@echo "Swagger UI available at http://localhost:8081"

swagger-stop:
	@docker stop swagger-editor || true
	@docker rm swagger-editor || true
	@echo "Swagger UI stopped"

coverage-summary:
	@echo "Coverage Summary:"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | grep -E "(total:|internal/)" | tail -10

quick-test:
	@echo "Starting quick API test..."
	@./test_api.sh

dev-setup:
	@echo "Setting up development environment..."
	@docker-compose up -d
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 3
	@echo "Setup complete! Run 'make run-server' to start"

help:
	@echo "GophKeeper - Password Manager"
	@echo ""
	@echo "Available targets:"
	@echo "  dev-setup       Setup Docker and start PostgreSQL"
	@echo "  build-server    Build the server binary"
	@echo "  build-client    Build the client binary"
	@echo "  build           Build both server and client"
	@echo "  test            Run all tests"
	@echo "  test-coverage   Run tests with coverage report"
	@echo "  html-coverage   Generate HTML coverage report"
	@echo "  coverage-summary Show coverage summary"
	@echo "  quick-test      Run API integration tests"
	@echo "  clean           Remove build artifacts"
	@echo "  run-server      Run the server"
	@echo "  run-client      Run the client"
	@echo "  docker-up       Start Docker containers"
	@echo "  docker-down     Stop Docker containers"
	@echo "  swagger-ui      Start Swagger UI for API documentation"
	@echo "  swagger-stop    Stop Swagger UI"
	@echo "  help            Show this help message"
