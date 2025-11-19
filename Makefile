# ==============================================================================
# CONFIGURATION & IMPORTS
# ==============================================================================

-include .env
export

APP_NAME = pms-core

MIGRATIONS_DIR = migrations

# ==============================================================================
# DATABASE COMMANDS
# ==============================================================================

## db-create: Create a new empty migration (e.g., make db-create name=add_users)
db-create:
	@echo "Creating migration files"
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

## db-up: Run all pending migrations
db-up:
	@echo "Applying migrations to: $(DATABASE_URL)"
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" -verbose up

## db-down: Undo the last migration applied
db-down:
	@echo "Reversing last migration"
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" -verbose down 1

## db-force: Force a specific version (fix ‘dirty’ status)
db-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $(version)

# ==============================================================================
# APLICATION COMMANDS
# ==============================================================================

## run: Start the server
run:
	@echo "Starting server on port $(PORT)"
	@go run cmd/api/main.go

## build: Compile the binary
build:
	go build -o bin/$(APP_NAME) cmd/api/main.go

## test: Run the project tests
test:
	go test -v ./...

.PHONY: db-create db-up db-down db-force run build test
