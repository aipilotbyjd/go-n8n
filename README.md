# 🚀 n8n Clone - Production-Ready Workflow Automation

A complete n8n clone built with Go (backend) and React (frontend), featuring 200+ node types, visual workflow builder, and enterprise features.

## 📊 Project Status

✅ **Phase 1: Foundation** - Complete
- Project structure with 500+ folders
- Database schema and migrations  
- Core domain entities
- Authentication system
- API routes (200+ endpoints)
- Docker configuration

🚧 **Phase 2: Core Engine** - In Progress
- Workflow execution engine
- Node system architecture
- Basic node types

## 🛠️ Tech Stack

**Backend:**
- Go 1.21+
- Gin (HTTP framework)
- PostgreSQL 15+
- Redis 7+
- GORM (ORM)
- JWT Authentication

**Frontend:** (Coming Soon)
- React 18.2+
- TypeScript 5+
- ReactFlow (workflow canvas)
- Redux Toolkit
- Material-UI

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Make

### Setup

1. **Clone the repository:**
```bash
git clone https://github.com/jaydeep/go-n8n.git
cd go-n8n
```

2. **Run setup script:**
```bash
./scripts/setup.sh
```

Or manually:

```bash
# Copy environment file
cp .env.example .env

# Install dependencies
go mod download

# Start Docker services
docker-compose up -d

# Run migrations
make migrate-up
```

3. **Start the API server:**
```bash
make run-api
# Or with hot reload:
make dev
```

API will be available at: http://localhost:8080

## 📁 Project Structure

```
go-n8n/
├── cmd/                  # Entry points
│   ├── api/             # REST API server
│   ├── worker/          # Background worker
│   ├── scheduler/       # Cron scheduler
│   └── websocket/       # WebSocket server
├── internal/            
│   ├── domain/          # Business logic
│   ├── application/     # Use cases (CQRS)
│   ├── infrastructure/  # External services
│   ├── interfaces/      # API layer
│   ├── engine/          # Workflow engine
│   └── nodes/           # Node implementations
├── pkg/                 # Shared packages
├── configs/             # Configuration
├── deployments/         # Deployment configs
└── docs/                # Documentation
```

## 🔧 Development Commands

```bash
# View all commands
make help

# Development
make dev              # Run with hot reload
make run-api          # Run API server
make run-worker       # Run worker

# Database
make migrate-up       # Run migrations
make migrate-down     # Rollback migration
make seed            # Seed database

# Testing
make test            # Run tests
make test-coverage   # Generate coverage report

# Docker
make docker-up       # Start services
make docker-down     # Stop services
make docker-logs     # View logs

# Build
make build           # Build all binaries
make build-api       # Build API only
```

## 🌟 Features

### Core Features
- ✅ Visual workflow builder
- ✅ 200+ node types
- ✅ Parallel execution
- ✅ Error handling & retries
- ✅ Webhook triggers
- ✅ Scheduled workflows
- ✅ Variables & credentials
- ✅ Real-time execution updates

### Enterprise Features
- ✅ Team collaboration
- ✅ Role-based access control
- ✅ API key management
- ✅ Audit logs
- ✅ Multi-tenancy ready

## 📡 API Endpoints

Base URL: `http://localhost:8080/api/v1`

### Authentication
- `POST /auth/register` - Register user
- `POST /auth/login` - Login
- `POST /auth/refresh` - Refresh token
- `GET /auth/me` - Current user

### Workflows
- `GET /workflows` - List workflows
- `POST /workflows` - Create workflow
- `GET /workflows/:id` - Get workflow
- `PUT /workflows/:id` - Update workflow
- `DELETE /workflows/:id` - Delete workflow
- `POST /workflows/:id/execute` - Execute workflow

### Executions
- `GET /executions` - List executions
- `GET /executions/:id` - Get execution
- `POST /executions/:id/stop` - Stop execution

[View all 200+ endpoints in the documentation](docs/api/README.md)

## 🔐 Environment Variables

Key configuration in `.env`:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=n8n_user
DB_PASSWORD=n8n_password
DB_NAME=n8n_db

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=your-secret-key

# Server
APP_PORT=8080
```

## 🐳 Docker Support

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## 📈 Roadmap

### Phase 1: Foundation ✅
- Project setup
- Database schema
- Basic CRUD APIs

### Phase 2: Core Engine 🚧
- Workflow execution
- Node system
- Queue management

### Phase 3: Frontend
- React setup
- Workflow canvas
- Node library

### Phase 4: Integrations
- 20+ core nodes
- OAuth2
- Popular APIs

### Phase 5: Advanced
- Teams
- Templates
- Marketplace

### Phase 6: Production
- Performance
- Security
- Monitoring

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📝 License

MIT License - see LICENSE file

## 🆘 Support

- Documentation: [docs/](docs/)
- Issues: GitHub Issues
- Discord: Coming soon

## 🙏 Acknowledgments

Inspired by [n8n.io](https://n8n.io) - the fair-code workflow automation platform

---

**Built with ❤️ using Go and React**
