# N8N Clone Backend Plan - Go with Gin Framework

## Project Overview
Building a robust workflow automation backend API that replicates n8n's core functionality using Go and Gin framework. This will serve as the backend for a React frontend application.

## Core Architecture

### 1. Technology Stack
- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: PostgreSQL (primary) + Redis (caching/queues)
- **ORM**: GORM or sqlx
- **Message Queue**: Redis + Asynq (for background jobs)
- **WebSocket**: Gorilla WebSocket or Melody
- **Authentication**: JWT with refresh tokens
- **API Documentation**: Swagger/OpenAPI 3.0
- **Container**: Docker & Docker Compose

### 2. Project Structure
```
go-n8n/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/            # HTTP handlers
│   │   │   ├── workflow.go
│   │   │   ├── execution.go
│   │   │   ├── node.go
│   │   │   ├── credential.go
│   │   │   ├── webhook.go
│   │   │   └── user.go
│   │   ├── middleware/          # Custom middleware
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   ├── ratelimit.go
│   │   │   └── logger.go
│   │   └── routes/              # Route definitions
│   │       └── routes.go
│   ├── core/
│   │   ├── workflow/            # Workflow engine
│   │   │   ├── engine.go
│   │   │   ├── executor.go
│   │   │   ├── scheduler.go
│   │   │   └── parser.go
│   │   ├── nodes/               # Node implementations
│   │   │   ├── base.go
│   │   │   ├── http_request.go
│   │   │   ├── webhook.go
│   │   │   ├── schedule.go
│   │   │   ├── function.go
│   │   │   └── transformer.go
│   │   └── integrations/        # External service integrations
│   │       ├── interface.go
│   │       ├── slack/
│   │       ├── github/
│   │       ├── database/
│   │       └── email/
│   ├── models/                  # Data models
│   │   ├── workflow.go
│   │   ├── node.go
│   │   ├── execution.go
│   │   ├── credential.go
│   │   └── user.go
│   ├── repository/              # Database layer
│   │   ├── interfaces.go
│   │   ├── workflow_repo.go
│   │   ├── execution_repo.go
│   │   └── credential_repo.go
│   ├── services/                # Business logic
│   │   ├── workflow_service.go
│   │   ├── execution_service.go
│   │   ├── credential_service.go
│   │   └── auth_service.go
│   ├── utils/                   # Utility functions
│   │   ├── crypto.go
│   │   ├── validator.go
│   │   └── response.go
│   └── config/                  # Configuration
│       └── config.go
├── pkg/                         # Public packages
│   ├── logger/
│   ├── database/
│   └── queue/
├── migrations/                  # Database migrations
├── scripts/                     # Build/deployment scripts
├── docker/                      # Docker configurations
├── docs/                        # API documentation
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

## Database Schema

### Core Tables

#### 1. Users
```sql
- id (UUID, PK)
- email (VARCHAR, UNIQUE)
- password_hash (VARCHAR)
- name (VARCHAR)
- role (ENUM: admin, user)
- is_active (BOOLEAN)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
```

#### 2. Workflows
```sql
- id (UUID, PK)
- name (VARCHAR)
- description (TEXT)
- user_id (UUID, FK)
- is_active (BOOLEAN)
- settings (JSONB)
- tags (VARCHAR[])
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
- deleted_at (TIMESTAMP, soft delete)
```

#### 3. Nodes
```sql
- id (UUID, PK)
- workflow_id (UUID, FK)
- type (VARCHAR)
- name (VARCHAR)
- position (JSONB) # {x: 100, y: 200}
- parameters (JSONB)
- credentials_id (UUID, FK, nullable)
- retry_on_fail (BOOLEAN)
- max_retries (INT)
- wait_between_retries (INT)
- continue_on_fail (BOOLEAN)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
```

#### 4. Node_Connections
```sql
- id (UUID, PK)
- workflow_id (UUID, FK)
- source_node_id (UUID, FK)
- source_output (VARCHAR)
- target_node_id (UUID, FK)
- target_input (VARCHAR)
- created_at (TIMESTAMP)
```

#### 5. Executions
```sql
- id (UUID, PK)
- workflow_id (UUID, FK)
- status (ENUM: waiting, running, success, error, cancelled)
- mode (ENUM: manual, trigger, webhook, schedule)
- started_at (TIMESTAMP)
- finished_at (TIMESTAMP)
- execution_time_ms (INT)
- data (JSONB) # Input/output data
- error (TEXT)
- retry_of (UUID, FK, nullable)
```

#### 6. Execution_Data
```sql
- id (UUID, PK)
- execution_id (UUID, FK)
- node_id (UUID, FK)
- data (JSONB)
- error (TEXT)
- execution_time_ms (INT)
- status (ENUM: waiting, running, success, error, skipped)
- started_at (TIMESTAMP)
- finished_at (TIMESTAMP)
```

#### 7. Credentials
```sql
- id (UUID, PK)
- name (VARCHAR)
- type (VARCHAR) # api_key, oauth2, basic_auth
- user_id (UUID, FK)
- data (JSONB, encrypted)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
```

#### 8. Webhooks
```sql
- id (UUID, PK)
- workflow_id (UUID, FK)
- node_id (UUID, FK)
- path (VARCHAR, UNIQUE)
- method (VARCHAR)
- is_active (BOOLEAN)
- created_at (TIMESTAMP)
```

#### 9. Scheduled_Workflows
```sql
- id (UUID, PK)
- workflow_id (UUID, FK)
- cron_expression (VARCHAR)
- timezone (VARCHAR)
- is_active (BOOLEAN)
- last_run_at (TIMESTAMP)
- next_run_at (TIMESTAMP)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
```

## API Endpoints

### Authentication
```
POST   /api/auth/register
POST   /api/auth/login
POST   /api/auth/refresh
POST   /api/auth/logout
GET    /api/auth/me
```

### Workflows
```
GET    /api/workflows                 # List all workflows
POST   /api/workflows                 # Create workflow
GET    /api/workflows/:id            # Get workflow details
PUT    /api/workflows/:id            # Update workflow
DELETE /api/workflows/:id            # Delete workflow
POST   /api/workflows/:id/activate   # Activate workflow
POST   /api/workflows/:id/deactivate # Deactivate workflow
POST   /api/workflows/:id/duplicate  # Duplicate workflow
GET    /api/workflows/:id/nodes      # Get workflow nodes
POST   /api/workflows/:id/execute    # Manual execution
GET    /api/workflows/:id/executions # Get execution history
```

### Nodes
```
GET    /api/nodes/types              # List available node types
GET    /api/nodes/:id                # Get node details
POST   /api/workflows/:wf_id/nodes   # Add node to workflow
PUT    /api/nodes/:id                # Update node
DELETE /api/nodes/:id                # Delete node
POST   /api/nodes/:id/test           # Test node execution
```

### Connections
```
POST   /api/workflows/:wf_id/connections    # Create connection
DELETE /api/connections/:id                 # Delete connection
```

### Executions
```
GET    /api/executions                      # List all executions
GET    /api/executions/:id                  # Get execution details
POST   /api/executions/:id/retry            # Retry failed execution
POST   /api/executions/:id/stop             # Stop running execution
DELETE /api/executions/:id                  # Delete execution record
GET    /api/executions/:id/data             # Get execution data
```

### Credentials
```
GET    /api/credentials                     # List credentials
POST   /api/credentials                     # Create credential
GET    /api/credentials/:id                 # Get credential (masked)
PUT    /api/credentials/:id                 # Update credential
DELETE /api/credentials/:id                 # Delete credential
POST   /api/credentials/:id/test            # Test credential
```

### Webhooks
```
GET    /api/webhooks                        # List webhooks
POST   /api/webhooks                        # Create webhook
DELETE /api/webhooks/:id                    # Delete webhook
* ANY  /webhook/:path                       # Webhook endpoint
```

### WebSocket
```
WS     /api/ws                              # Real-time updates
```

## Workflow Execution Engine

### Core Components

#### 1. Workflow Parser
- Parse workflow JSON structure
- Build execution graph (DAG)
- Validate node connections
- Check for circular dependencies

#### 2. Execution Queue
```go
type ExecutionRequest struct {
    WorkflowID   string
    ExecutionID  string
    TriggerData  map[string]interface{}
    Mode         ExecutionMode
    RetryOf      *string
}
```

#### 3. Node Executor
```go
type NodeExecutor interface {
    Execute(ctx context.Context, input NodeInput) (NodeOutput, error)
    Validate(parameters map[string]interface{}) error
    GetSchema() NodeSchema
}
```

#### 4. Execution Flow
1. Receive execution trigger (manual/webhook/schedule)
2. Load workflow configuration
3. Parse and validate workflow
4. Create execution record
5. Queue execution job
6. Execute nodes in topological order
7. Handle branching and merging
8. Process errors and retries
9. Store execution results
10. Send real-time updates via WebSocket

### Node Types

#### Basic Nodes
- **Start Node**: Entry point for workflows
- **Schedule Trigger**: Cron-based triggers
- **Webhook Trigger**: HTTP webhook receiver
- **Manual Trigger**: User-initiated execution

#### Action Nodes
- **HTTP Request**: Make HTTP API calls
- **Database Query**: Execute SQL queries
- **Code Function**: Run custom JavaScript/Python code
- **Wait**: Delay execution
- **Set Variable**: Set workflow variables
- **If/Switch**: Conditional branching
- **Loop**: Iterate over arrays
- **Merge**: Combine multiple branches
- **Split**: Split data into batches

#### Integration Nodes
- **Email**: Send emails
- **Slack**: Send messages, interact with Slack API
- **GitHub**: Issues, PRs, Actions
- **Google Sheets**: Read/write spreadsheets
- **Database**: MySQL, PostgreSQL, MongoDB
- **AWS**: S3, Lambda, SQS
- **Custom Webhook**: Send data to external webhooks

## Authentication & Security

### JWT Implementation
```go
type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int    `json:"expires_in"`
}

type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.StandardClaims
}
```

### Security Measures
1. **Password Hashing**: bcrypt with cost 12
2. **Credential Encryption**: AES-256-GCM
3. **Rate Limiting**: Per-IP and per-user limits
4. **CORS Configuration**: Whitelist frontend domain
5. **Input Validation**: Strict validation on all inputs
6. **SQL Injection Prevention**: Parameterized queries
7. **XSS Prevention**: Sanitize all outputs
8. **Audit Logging**: Track all critical operations

## Background Jobs & Queue

### Job Types
```go
const (
    JobTypeWorkflowExecution = "workflow_execution"
    JobTypeScheduledTrigger  = "scheduled_trigger"
    JobTypeWebhookProcessing = "webhook_processing"
    JobTypeEmailNotification = "email_notification"
    JobTypeCleanup          = "cleanup"
)
```

### Queue Implementation with Asynq
```go
type WorkflowExecutionPayload struct {
    WorkflowID   string
    ExecutionID  string
    TriggerData  map[string]interface{}
}

func HandleWorkflowExecution(ctx context.Context, t *asynq.Task) error {
    var payload WorkflowExecutionPayload
    if err := json.Unmarshal(t.Payload(), &payload); err != nil {
        return err
    }
    // Execute workflow
    return executeWorkflow(ctx, payload)
}
```

## Real-time Updates

### WebSocket Events
```go
type WSMessage struct {
    Type    string      `json:"type"`
    Data    interface{} `json:"data"`
    EventID string      `json:"event_id"`
}

// Event Types
const (
    EventExecutionStarted   = "execution.started"
    EventExecutionCompleted = "execution.completed"
    EventExecutionFailed    = "execution.failed"
    EventNodeExecuting      = "node.executing"
    EventNodeCompleted      = "node.completed"
    EventWorkflowUpdated    = "workflow.updated"
)
```

## Error Handling

### Error Types
```go
type AppError struct {
    Code       string `json:"code"`
    Message    string `json:"message"`
    StatusCode int    `json:"status_code"`
    Details    map[string]interface{} `json:"details,omitempty"`
}

const (
    ErrCodeValidation      = "VALIDATION_ERROR"
    ErrCodeUnauthorized    = "UNAUTHORIZED"
    ErrCodeForbidden       = "FORBIDDEN"
    ErrCodeNotFound        = "NOT_FOUND"
    ErrCodeConflict        = "CONFLICT"
    ErrCodeInternal        = "INTERNAL_ERROR"
    ErrCodeRateLimit       = "RATE_LIMIT"
    ErrCodeWorkflowInvalid = "WORKFLOW_INVALID"
    ErrCodeExecutionFailed = "EXECUTION_FAILED"
)
```

## Performance Optimizations

1. **Connection Pooling**: Database connection pool
2. **Query Optimization**: Indexes on frequently queried fields
3. **Caching Strategy**:
   - Redis for session data
   - In-memory cache for workflow definitions
   - Cache execution results for retry
4. **Pagination**: Limit/offset for list endpoints
5. **Lazy Loading**: Load node details on demand
6. **Batch Processing**: Process multiple nodes in parallel when possible
7. **Resource Limits**: Max execution time, memory limits

## Monitoring & Logging

### Structured Logging
```go
logger.Info("workflow executed",
    zap.String("workflow_id", workflowID),
    zap.String("execution_id", executionID),
    zap.Duration("duration", duration),
    zap.String("status", status),
)
```

### Metrics to Track
- API response times
- Workflow execution duration
- Queue depth
- Error rates
- Active connections
- Database query performance

## Development Phases

### Phase 1: Foundation (Week 1-2)
- [ ] Project setup and structure
- [ ] Database schema and migrations
- [ ] Basic CRUD for workflows
- [ ] Authentication system
- [ ] API documentation setup

### Phase 2: Core Engine (Week 3-4)
- [ ] Workflow parser and validator
- [ ] Basic node executor
- [ ] Execution queue setup
- [ ] Simple node types (HTTP, webhook, manual trigger)
- [ ] Execution history tracking

### Phase 3: Advanced Features (Week 5-6)
- [ ] Complex node types (conditions, loops, merge)
- [ ] Credential management
- [ ] Error handling and retries
- [ ] WebSocket implementation
- [ ] Scheduled workflows

### Phase 4: Integrations (Week 7-8)
- [ ] External service integrations
- [ ] OAuth2 implementation
- [ ] Custom code execution nodes
- [ ] Webhook management

### Phase 5: Optimization & Testing (Week 9-10)
- [ ] Performance optimization
- [ ] Comprehensive testing
- [ ] Load testing
- [ ] Security audit
- [ ] Documentation completion

## Testing Strategy

### Unit Tests
- Repository layer tests
- Service layer tests
- Node executor tests
- Utility function tests

### Integration Tests
- API endpoint tests
- Workflow execution tests
- Database integration tests
- Queue processing tests

### E2E Tests
- Complete workflow scenarios
- Authentication flows
- WebSocket communication
- Error handling scenarios

## Deployment

### Docker Setup
```dockerfile
# Multi-stage build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

### Environment Variables
```env
# Server
PORT=8080
ENV=production
DEBUG=false

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=n8n_user
DB_PASSWORD=secret
DB_NAME=n8n_clone

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=your-secret-key
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# Encryption
ENCRYPTION_KEY=32-byte-key-for-aes-256

# CORS
ALLOWED_ORIGINS=http://localhost:3000

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_DURATION=1m
```

## Key Decisions & Rationale

1. **Go + Gin**: High performance, excellent concurrency support, perfect for workflow execution
2. **PostgreSQL**: ACID compliance, JSONB support for flexible node data
3. **Redis**: Fast caching, reliable queue backend
4. **JWT**: Stateless authentication, scalable
5. **WebSocket**: Real-time execution updates essential for UX
6. **Microservices-ready**: Can be split into services later if needed

## Success Metrics

- API response time < 100ms for CRUD operations
- Workflow execution latency < 500ms
- Support 1000+ concurrent workflow executions
- 99.9% uptime
- Zero data loss
- < 1% error rate

## Next Steps

1. Set up development environment
2. Initialize Go module and install dependencies
3. Create database schema and migrations
4. Implement authentication system
5. Build basic workflow CRUD operations
6. Develop workflow execution engine
7. Add node types incrementally
8. Implement real-time updates
9. Add integrations
10. Optimize and test thoroughly

This plan provides a solid foundation for building a production-ready n8n clone backend in Go.
