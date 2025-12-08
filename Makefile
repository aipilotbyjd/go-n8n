.PHONY: help build run test clean setup migrate docker-up docker-down

# Variables
BINARY_NAME=go-n8n
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

help: ## Show this help message
	@echo '${GREEN}n8n Clone - Makefile Commands${NC}'
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${NC} ${GREEN}<command>${NC}'
	@echo ''
	@echo 'Available commands:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  ${YELLOW}%-20s${NC} %s\n", $$1, $$2}'

setup: ## Initial project setup
	@echo "${GREEN}Setting up development environment...${NC}"
	@go mod download
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/google/wire/cmd/wire@latest
	@cp .env.example .env 2>/dev/null || true
	@echo "${GREEN}✓ Setup complete!${NC}"

deps: ## Download and install dependencies
	@echo "${GREEN}Installing dependencies...${NC}"
	@go mod download
	@go mod tidy
	@echo "${GREEN}✓ Dependencies installed!${NC}"

build: ## Build all binaries
	@echo "${GREEN}Building binaries...${NC}"
	@CGO_ENABLED=0 go build ${LDFLAGS} -o bin/api cmd/api/main.go
	@CGO_ENABLED=0 go build ${LDFLAGS} -o bin/worker cmd/worker/main.go
	@CGO_ENABLED=0 go build ${LDFLAGS} -o bin/scheduler cmd/scheduler/main.go
	@CGO_ENABLED=0 go build ${LDFLAGS} -o bin/websocket cmd/websocket/main.go
	@CGO_ENABLED=0 go build ${LDFLAGS} -o bin/migrate cmd/migrate/main.go
	@echo "${GREEN}✓ Build complete!${NC}"

build-api: ## Build API server
	@CGO_ENABLED=0 go build ${LDFLAGS} -o bin/api cmd/api/main.go

build-worker: ## Build worker
	@CGO_ENABLED=0 go build ${LDFLAGS} -o bin/worker cmd/worker/main.go

run-api: ## Run API server
	@go run cmd/api/main.go

run-worker: ## Run worker
	@go run cmd/worker/main.go

run-scheduler: ## Run scheduler
	@go run cmd/scheduler/main.go

run-websocket: ## Run WebSocket server
	@go run cmd/websocket/main.go

dev: ## Run with hot reload (requires air)
	@~/go/bin/air -c .air.toml

test: ## Run tests
	@echo "${GREEN}Running tests...${NC}"
	@go test -v -race -cover ./...

test-coverage: ## Run tests with coverage
	@echo "${GREEN}Running tests with coverage...${NC}"
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "${GREEN}✓ Coverage report generated: coverage.html${NC}"

test-integration: ## Run integration tests
	@go test -v -tags=integration ./test/integration/...

test-e2e: ## Run end-to-end tests
	@go test -v -tags=e2e ./test/e2e/...

benchmark: ## Run benchmarks
	@go test -bench=. -benchmem ./...

lint: ## Run linter
	@echo "${GREEN}Running linter...${NC}"
	@golangci-lint run --deadline=5m

fmt: ## Format code
	@echo "${GREEN}Formatting code...${NC}"
	@go fmt ./...
	@gofmt -s -w .

vet: ## Run go vet
	@go vet ./...

security: ## Run security scan
	@gosec -fmt sarif -out results.sarif ./...

clean: ## Clean build artifacts
	@echo "${GREEN}Cleaning...${NC}"
	@rm -rf bin/ coverage.* vendor/ *.out
	@echo "${GREEN}✓ Cleaned!${NC}"

# Database commands
migrate-up: ## Run database migrations
	@echo "${GREEN}Running migrations...${NC}"
	@migrate -path internal/infrastructure/persistence/postgres/migrations -database "postgres://n8n_user:n8n_password@localhost:5432/n8n_db?sslmode=disable" up

migrate-down: ## Rollback last migration
	@migrate -path internal/infrastructure/persistence/postgres/migrations -database "postgres://n8n_user:n8n_password@localhost:5432/n8n_db?sslmode=disable" down 1

migrate-create: ## Create new migration (usage: make migrate-create name=create_users_table)
	@migrate create -ext sql -dir internal/infrastructure/persistence/postgres/migrations -seq $(name)

seed: ## Seed database
	@go run cmd/migrate/seed.go

# Docker commands
docker-up: ## Start all services with docker-compose
	@echo "${GREEN}Starting Docker services...${NC}"
	@docker-compose up -d
	@echo "${GREEN}✓ Services started!${NC}"

docker-down: ## Stop all services
	@echo "${GREEN}Stopping Docker services...${NC}"
	@docker-compose down
	@echo "${GREEN}✓ Services stopped!${NC}"

docker-logs: ## Show logs from all services
	@docker-compose logs -f

docker-build: ## Build Docker images
	@docker build -f deployments/docker/Dockerfile.api -t ${BINARY_NAME}-api:${VERSION} .
	@docker build -f deployments/docker/Dockerfile.worker -t ${BINARY_NAME}-worker:${VERSION} .

docker-clean: ## Clean Docker resources
	@docker-compose down -v
	@docker system prune -f

# Development helpers
gen: ## Generate code (wire, mocks, etc.)
	@go generate ./...
	@wire ./...

watch: ## Watch for file changes and rebuild
	@air -c .air.toml

logs-api: ## Tail API logs
	@docker-compose logs -f api

logs-worker: ## Tail worker logs
	@docker-compose logs -f worker

psql: ## Connect to PostgreSQL
	@docker-compose exec postgres psql -U n8n_user -d n8n_db

redis-cli: ## Connect to Redis
	@docker-compose exec redis redis-cli

# Quick start commands
quickstart: docker-up migrate-up seed ## Quick start: start Docker, run migrations, seed data
	@echo "${GREEN}✓ Quick start complete! API running on http://localhost:8080${NC}"

reset: docker-down docker-up migrate-up seed ## Reset everything: restart Docker, fresh migrations, seed
	@echo "${GREEN}✓ Reset complete!${NC}"

# Production commands
prod-build: ## Build for production
	@echo "${GREEN}Building for production...${NC}"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -a -installsuffix cgo -o bin/api-linux cmd/api/main.go
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -a -installsuffix cgo -o bin/worker-linux cmd/worker/main.go
	@echo "${GREEN}✓ Production build complete!${NC}"

deploy: prod-build ## Deploy to production
	@echo "${GREEN}Deploying to production...${NC}"
	@kubectl apply -f deployments/kubernetes/
	@echo "${GREEN}✓ Deployed!${NC}"

# Default target
.DEFAULT_GOAL := help
