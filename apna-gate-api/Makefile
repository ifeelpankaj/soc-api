# Binary output path
BIN_DIR := bin
SERVER_BIN := $(BIN_DIR)/server
MIGRATE_ENV ?= development
RUN_ENV_FILE_ARG :=
COMMIT_MSG_FILE_ARG :=

ifneq ($(strip $(ENV_FILE)),)
RUN_ENV_FILE_ARG := -env-file "$(ENV_FILE)"
endif

ifneq ($(strip $(MSG_FILE)),)
COMMIT_MSG_FILE_ARG := -msg-file "$(MSG_FILE)"
endif

.PHONY: dev prod build build-prod install test test-coverage test-integration test-all test-race clean watch fmt lint quality hooks-install commit-msg-check help run serve sqlc sqlc-vet check kill-port

help:
	@echo "Available targets:"
	@echo "  dev              - Select env file, kill its port, then start development server"
	@echo "  prod             - Select env file, kill its port, then start production server"
	@echo "  run              - Alias for dev"
	@echo "  serve            - Alias for run"
	@echo "  build            - Build application"
	@echo "  build-prod       - Build optimized production binary"
	@echo "  install          - Install dependencies"
	@echo "  test             - Run fast unit and handler tests"
	@echo "  test-integration - Run DB-backed integration tests (requires Docker)"
	@echo "  test-all         - Run fast, integration, race, and coverage checks"
	@echo "  test-race        - Run fast tests with the Go race detector"
	@echo "  test-coverage    - Run all tests and write coverage output to test_tmp"
	@echo "  clean            - Clean build artifacts"
	@echo "  watch            - Run with live reload (requires air)"
	@echo "  fmt              - Format code"
	@echo "  lint             - Lint code (requires golangci-lint)"
	@echo "  quality          - Format, lint, and test API code"
	@echo "  hooks-install    - Install Lefthook Git hooks from the repo root"
	@echo "  commit-msg-check - Check a commit message file with API rules"
	@echo "  sqlc             - Generate sqlc database code"
	@echo "  sqlc-vet         - Vet sqlc queries"
	@echo "  check            - Generate sqlc, format and test"
	@echo "  migrate-up       - Run migrations up (MIGRATE_ENV=development by default)"
	@echo "  migrate-down     - Rollback last migration"
	@echo "  migrate-create   - Create new migration file"
	@echo "  migrate-status   - Show migration status"

# Single cross-platform port killer (Go, works on Windows/Linux/Mac, cmd/PowerShell/Git Bash)
# Usage: make kill-port ENV_FILE=.env.development
kill-port:
	@go run scripts/killport.go "$(ENV_FILE)"

dev:
	@go run scripts/runserver.go -mode development $(RUN_ENV_FILE_ARG)

prod:
	@go run scripts/runserver.go -mode production $(RUN_ENV_FILE_ARG)

run: dev

serve: run

build:
	@go run scripts/build.go $(SERVER_BIN)

build-prod:
	@go run scripts/build.go -prod $(SERVER_BIN)

install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

test:
	@echo "Running fast tests..."
	go test -v ./...

test-coverage:
	@echo "Running test coverage..."
	go run ./scripts/coverage

test-integration:
	@echo "Running integration tests..."
	go test -v -tags=integration ./...

test-race:
	@echo "Running race tests..."
	go test -race ./...

test-all: test test-integration test-race test-coverage

clean:
	@echo "Cleaning build artifacts..."
ifeq ($(OS),Windows_NT)
	@if exist $(BIN_DIR) rmdir /S /Q $(BIN_DIR) 2>nul || echo.
	@if exist test_tmp rmdir /S /Q test_tmp 2>nul || echo.
else
	@rm -rf $(BIN_DIR)
	@rm -rf test_tmp
endif
	go clean
	@echo "Clean complete"

watch:
	@go run scripts/watch.go $(RUN_ENV_FILE_ARG)

fmt:
	@echo "Formatting code..."
	go fmt ./...

lint:
	@echo "Linting code..."
	@echo "Make sure 'golangci-lint' is installed"
	@echo "To install run go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"
	golangci-lint run ./...

quality: fmt lint test

hooks-install:
	@echo "Installing Lefthook hooks..."
	lefthook install

commit-msg-check:
	@go run scripts/commitmsgcheck.go $(COMMIT_MSG_FILE_ARG)

swagger:
	@echo "Generating Swagger documentation..."
	@echo "Make sure 'swag' is installed: go install github.com/swaggo/swag/cmd/swag@latest"
	swag init -g main.go -d cmd/server,internal/handlers/v1,internal/models,internal/server --parseInternal

sqlc:
	@echo "Generating sqlc code..."
	sqlc generate

sqlc-vet:
	@echo "Checking sqlc queries..."
	sqlc vet

check: sqlc fmt test

migrate-up:
	@echo "Running migrations (env=$(MIGRATE_ENV))..."
	sql-migrate up -env="$(MIGRATE_ENV)"

migrate-down:
	@echo "Rolling back last migration (env=$(MIGRATE_ENV))..."
	sql-migrate down -env="$(MIGRATE_ENV)" -limit=1

migrate-create:
	@go run scripts/createmigration.go

migrate-status:
	@echo "Checking migration status (env=$(MIGRATE_ENV))..."
	sql-migrate status -env="$(MIGRATE_ENV)"
