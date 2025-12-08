#!/bin/bash

# Setup script for n8n Clone

set -e

echo "🚀 n8n Clone Setup Script"
echo "========================="
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed${NC}"
    echo "Please install Go 1.21+ from https://golang.org/dl/"
    exit 1
else
    echo -e "${GREEN}✅ Go is installed:${NC} $(go version)"
fi

# Check if Docker is running
if docker info &> /dev/null; then
    echo -e "${GREEN}✅ Docker is running${NC}"
else
    echo -e "${YELLOW}⚠️  Docker is not running${NC}"
    echo "Please start Docker Desktop or run: open -a Docker"
    echo "Waiting for Docker to start..."
    
    # Try to start Docker on macOS
    if [[ "$OSTYPE" == "darwin"* ]]; then
        open -a Docker 2>/dev/null || true
        
        # Wait for Docker to be ready
        counter=0
        while ! docker info &> /dev/null && [ $counter -lt 30 ]; do
            echo -n "."
            sleep 2
            counter=$((counter+1))
        done
        
        if docker info &> /dev/null; then
            echo -e "\n${GREEN}✅ Docker started successfully${NC}"
        else
            echo -e "\n${RED}❌ Docker failed to start. Please start it manually.${NC}"
            exit 1
        fi
    else
        exit 1
    fi
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo -e "${YELLOW}Creating .env file...${NC}"
    cp .env.example .env
    echo -e "${GREEN}✅ .env file created${NC}"
else
    echo -e "${GREEN}✅ .env file exists${NC}"
fi

# Install Go dependencies
echo -e "${YELLOW}Installing Go dependencies...${NC}"
go mod download
echo -e "${GREEN}✅ Go dependencies installed${NC}"

# Start Docker services
echo -e "${YELLOW}Starting PostgreSQL and Redis...${NC}"
docker-compose up -d
echo -e "${GREEN}✅ Database services started${NC}"

# Wait for PostgreSQL to be ready
echo -e "${YELLOW}Waiting for PostgreSQL to be ready...${NC}"
sleep 5

# Run database migrations
echo -e "${YELLOW}Running database migrations...${NC}"
docker-compose exec -T postgres psql -U n8n_user -d n8n_db -f - < internal/infrastructure/persistence/postgres/migrations/001_initial_schema.sql 2>/dev/null || {
    echo -e "${YELLOW}Note: Migrations will be applied when you run 'make migrate-up'${NC}"
}

echo ""
echo -e "${GREEN}🎉 Setup Complete!${NC}"
echo ""
echo "Next steps:"
echo "1. Start the API server:     make run-api"
echo "2. Or use hot reload:        make dev"
echo "3. View all commands:        make help"
echo ""
echo "Services running:"
echo "  - PostgreSQL: localhost:5432"
echo "  - Redis:      localhost:6379"
echo ""
echo "Default credentials:"
echo "  - DB User: n8n_user"
echo "  - DB Pass: n8n_password"
echo "  - DB Name: n8n_db"
echo ""
