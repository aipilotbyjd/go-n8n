# 🚀 Complete n8n Clone Architecture - Single Source of Truth

## Table of Contents
1. [Project Overview](#project-overview)
2. [Architecture Principles](#architecture-principles)
3. [Technology Stack](#technology-stack)
4. [Complete Project Structure](#complete-project-structure)
5. [Database Architecture](#database-architecture)
6. [API Architecture](#api-architecture)
7. [Workflow Execution Engine](#workflow-execution-engine)
8. [Node System Architecture](#node-system-architecture)
9. [Frontend Architecture](#frontend-architecture)
10. [Security Architecture](#security-architecture)
11. [Deployment Architecture](#deployment-architecture)
12. [Implementation Roadmap](#implementation-roadmap)

---

## 🎯 Project Overview

### What We're Building
A **production-ready n8n clone** in Go with React frontend that matches n8n's capabilities:
- **200+ Node Types**: All major integrations (Slack, GitHub, OpenAI, databases)
- **Visual Workflow Builder**: Drag-and-drop interface with ReactFlow
- **Execution Engine**: Parallel execution, retries, error handling
- **Real-time Updates**: WebSocket-based live execution status
- **Enterprise Features**: Teams, RBAC, audit logs, API keys
- **Scalability**: Handle 100K+ workflows, 1M+ executions/day

### System Requirements
- **Performance**: < 100ms API response, < 500ms workflow start
- **Scale**: 10K concurrent users, 100K workflows, 1M executions/day
- **Reliability**: 99.9% uptime, zero data loss
- **Security**: JWT auth, OAuth2, encryption at rest/transit

---

## 🏗️ Architecture Principles

1. **Domain-Driven Design (DDD)**: Business logic separated by domains
2. **Hexagonal Architecture**: Core logic independent of infrastructure
3. **CQRS Pattern**: Separate read/write operations
4. **Event-Driven**: Loose coupling via events
5. **Plugin Architecture**: Extensible node system
6. **Microservices-Ready**: Can split when needed
7. **Repository Pattern**: Abstract data access

---

## 🛠️ Technology Stack

### Backend (Go)
```yaml
Core:
  - Language: Go 1.21+
  - Framework: Gin (HTTP)
  - Database: PostgreSQL 15+
  - Cache: Redis 7+
  - Queue: Asynq (Redis-based)
  - WebSocket: Gorilla WebSocket

Libraries:
  - ORM: GORM / sqlx
  - Validation: go-playground/validator
  - JWT: golang-jwt/jwt/v5
  - UUID: google/uuid
  - Config: spf13/viper
  - Logger: uber-go/zap
  - Migration: golang-migrate

Monitoring:
  - Metrics: Prometheus
  - Tracing: OpenTelemetry
  - Logging: Zap + ELK Stack
```

### Frontend (React)
```yaml
Core:
  - Framework: React 18.2+
  - Language: TypeScript 5+
  - Build: Vite 5+
  - Routing: React Router v6

State & Data:
  - State: Redux Toolkit
  - API: TanStack Query
  - WebSocket: Socket.io-client
  - Forms: React Hook Form

UI & Workflow:
  - Workflow Canvas: ReactFlow 11+
  - UI Framework: Material-UI 5+
  - Styling: Tailwind CSS
  - Code Editor: Monaco Editor
  - Icons: MUI Icons
```

---

## 📁 Complete Project Structure

### Backend Structure (Go)

```
go-n8n/
├── cmd/                                    # Entry points for different services
│   ├── api/                               # REST API server
│   │   ├── main.go                        # API server entry point
│   │   ├── server.go                      # HTTP server configuration
│   │   ├── wire.go                        # Dependency injection
│   │   └── config.go                      # API configuration
│   │
│   ├── worker/                            # Background job processor
│   │   ├── main.go                        # Worker entry point
│   │   ├── executor_worker.go             # Workflow execution worker
│   │   ├── webhook_worker.go              # Webhook processing
│   │   └── email_worker.go                # Email sending worker
│   │
│   ├── scheduler/                         # Cron job scheduler
│   │   ├── main.go                        # Scheduler entry point
│   │   ├── cron_parser.go                 # Cron expression parser
│   │   └── timezone_handler.go            # Timezone management
│   │
│   ├── websocket/                         # WebSocket server
│   │   ├── main.go                        # WebSocket server entry
│   │   ├── hub.go                         # Connection management
│   │   └── handlers.go                    # Message handlers
│   │
│   └── migrate/                           # Database tools
│       ├── main.go                        # Migration runner
│       └── seed.go                        # Data seeding
│
├── internal/                              # Private application code
│   ├── domain/                           # Core business logic
│   │   ├── workflow/                     # Workflow domain
│   │   │   ├── entity.go                 # Workflow entity
│   │   │   ├── value_objects.go          # WorkflowStatus, NodePosition
│   │   │   ├── repository.go             # Repository interface
│   │   │   ├── service.go                # Domain service
│   │   │   ├── events.go                 # Domain events
│   │   │   ├── validator.go              # Validation rules
│   │   │   └── version.go                # Versioning logic
│   │   │
│   │   ├── execution/                    # Execution domain
│   │   │   ├── entity.go
│   │   │   ├── state_machine.go          # Execution states
│   │   │   ├── context.go                # Execution context
│   │   │   ├── data_flow.go              # Data between nodes
│   │   │   └── retry_policy.go           # Retry strategies
│   │   │
│   │   ├── node/                         # Node domain
│   │   │   ├── entity.go
│   │   │   ├── types.go                  # Node type definitions
│   │   │   ├── registry.go               # Node registry
│   │   │   ├── connection.go             # Node connections
│   │   │   └── parameters.go             # Node parameters
│   │   │
│   │   ├── credential/                   # Credentials domain
│   │   │   ├── entity.go
│   │   │   ├── encryption.go             # Credential encryption
│   │   │   ├── provider.go               # OAuth providers
│   │   │   └── types.go                  # Credential types
│   │   │
│   │   ├── user/                         # User domain
│   │   │   ├── entity.go
│   │   │   ├── permissions.go            # RBAC
│   │   │   ├── session.go                # User sessions
│   │   │   └── api_key.go                # API key management
│   │   │
│   │   └── webhook/                      # Webhook domain
│   │       ├── entity.go
│   │       ├── registry.go               # Webhook registry
│   │       └── validator.go              # Path validation
│   │
│   ├── application/                      # Use cases (CQRS)
│   │   ├── workflow/
│   │   │   ├── commands/                 # Write operations
│   │   │   │   ├── create_workflow.go
│   │   │   │   ├── update_workflow.go
│   │   │   │   ├── delete_workflow.go
│   │   │   │   ├── activate_workflow.go
│   │   │   │   ├── execute_workflow.go
│   │   │   │   └── duplicate_workflow.go
│   │   │   └── queries/                  # Read operations
│   │   │       ├── get_workflow.go
│   │   │       ├── list_workflows.go
│   │   │       └── search_workflows.go
│   │   │
│   │   └── execution/
│   │       ├── commands/
│   │       │   ├── start_execution.go
│   │       │   ├── stop_execution.go
│   │       │   └── retry_execution.go
│   │       └── queries/
│   │           └── get_execution_status.go
│   │
│   ├── infrastructure/                   # External implementations
│   │   ├── persistence/
│   │   │   ├── postgres/                 # PostgreSQL repositories
│   │   │   │   ├── workflow_repository.go
│   │   │   │   ├── execution_repository.go
│   │   │   │   ├── migrations/           # SQL migrations
│   │   │   │   └── connection.go
│   │   │   └── redis/                    # Redis cache
│   │   │       ├── cache_repository.go
│   │   │       └── session_store.go
│   │   │
│   │   ├── messaging/                    # Message queues
│   │   │   └── asynq/
│   │   │       ├── client.go
│   │   │       └── handlers.go
│   │   │
│   │   └── security/                     # Security implementations
│   │       ├── jwt/
│   │       │   ├── generator.go
│   │       │   └── validator.go
│   │       └── encryption/
│   │           └── aes.go
│   │
│   ├── interfaces/                       # API layer
│   │   ├── http/
│   │   │   ├── rest/
│   │   │   │   └── v1/
│   │   │   │       ├── controllers/      # All controllers
│   │   │   │       │   ├── workflow_controller.go
│   │   │   │       │   ├── execution_controller.go
│   │   │   │       │   ├── node_controller.go
│   │   │   │       │   ├── credential_controller.go
│   │   │   │       │   └── auth_controller.go
│   │   │   │       ├── middleware/       # HTTP middleware
│   │   │   │       │   ├── auth.go
│   │   │   │       │   ├── cors.go
│   │   │   │       │   └── rate_limit.go
│   │   │   │       └── routes.go         # Route definitions
│   │   │   └── websocket/
│   │   │       ├── hub.go                # WebSocket hub
│   │   │       └── handlers.go           # Message handlers
│   │
│   ├── engine/                           # Workflow execution engine
│   │   ├── executor/
│   │   │   ├── workflow_executor.go      # Main executor
│   │   │   ├── node_executor.go          # Node execution
│   │   │   ├── parallel_executor.go      # Parallel branches
│   │   │   └── error_handler.go          # Error handling
│   │   ├── scheduler/
│   │   │   └── cron_scheduler.go         # Cron scheduling
│   │   └── queue/
│   │       └── job_queue.go              # Job queue management
│   │
│   └── nodes/                            # Node implementations
│       ├── registry.go                   # Node registry
│       ├── base_node.go                  # Base node interface
│       ├── core/                         # Built-in nodes
│       │   ├── trigger/
│       │   │   ├── webhook_trigger.go
│       │   │   ├── schedule_trigger.go
│       │   │   └── manual_trigger.go
│       │   ├── action/
│       │   │   ├── http_request.go
│       │   │   ├── email_send.go
│       │   │   └── database_query.go
│       │   └── transform/
│       │       ├── set_node.go
│       │       ├── filter_node.go
│       │       └── merge_node.go
│       └── integrations/                 # Third-party integrations
│           ├── slack/
│           │   └── slack_node.go
│           ├── github/
│           │   └── github_node.go
│           └── openai/
│               └── chatgpt_node.go
│
├── pkg/                                  # Public packages
│   ├── database/                         # Database utilities
│   ├── logger/                           # Logging
│   ├── validator/                        # Validation
│   └── errors/                           # Error handling
│
├── configs/                              # Configuration files
│   └── config.yaml
├── deployments/                          # Deployment configs
│   ├── docker/
│   └── kubernetes/
├── scripts/                              # Build scripts
├── test/                                 # Test files
├── Makefile
├── go.mod
└── README.md
```

### Frontend Structure (React)

```
n8n-frontend/
├── src/
│   ├── app/                             # Application setup
│   │   ├── App.tsx
│   │   ├── store.ts                     # Redux store
│   │   └── router.tsx
│   │
│   ├── features/                        # Feature modules
│   │   ├── editor/                      # Workflow editor
│   │   │   ├── components/
│   │   │   │   ├── Canvas/
│   │   │   │   │   └── WorkflowCanvas.tsx
│   │   │   │   ├── Nodes/
│   │   │   │   │   └── NodeTypes/
│   │   │   │   └── Panels/
│   │   │   │       └── PropertiesPanel.tsx
│   │   │   └── store/
│   │   │       └── editorSlice.ts
│   │   │
│   │   ├── workflows/                   # Workflow management
│   │   │   ├── pages/
│   │   │   └── components/
│   │   │
│   │   └── auth/                        # Authentication
│   │       ├── pages/
│   │       └── hooks/
│   │
│   ├── shared/                          # Shared resources
│   │   ├── components/
│   │   ├── hooks/
│   │   └── services/
│   │
│   └── styles/
│
├── package.json
├── vite.config.ts
└── tsconfig.json
```

---

## 💾 Database Architecture

### Core Tables

```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Workflows table
CREATE TABLE workflows (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    user_id UUID REFERENCES users(id),
    team_id UUID REFERENCES teams(id),
    is_active BOOLEAN DEFAULT false,
    nodes JSONB DEFAULT '[]',
    connections JSONB DEFAULT '[]',
    settings JSONB DEFAULT '{}',
    tags VARCHAR(50)[] DEFAULT '{}',
    version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT workflow_name_user_unique UNIQUE(name, user_id, deleted_at)
);

-- Executions table (partitioned by date)
CREATE TABLE executions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID REFERENCES workflows(id),
    workflow_version INT NOT NULL,
    status VARCHAR(50) NOT NULL, -- waiting, running, success, error, cancelled
    mode VARCHAR(50) NOT NULL, -- manual, trigger, webhook, schedule
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP,
    execution_time_ms INT,
    input_data JSONB DEFAULT '{}',
    output_data JSONB DEFAULT '{}',
    error_message TEXT,
    retry_of UUID REFERENCES executions(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) PARTITION BY RANGE (created_at);

-- Node execution data
CREATE TABLE execution_node_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    execution_id UUID REFERENCES executions(id) ON DELETE CASCADE,
    node_id VARCHAR(255) NOT NULL,
    node_type VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL,
    input_data JSONB DEFAULT '{}',
    output_data JSONB DEFAULT '{}',
    error_message TEXT,
    execution_time_ms INT,
    started_at TIMESTAMP,
    finished_at TIMESTAMP
);

-- Credentials table
CREATE TABLE credentials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL, -- oauth2, api_key, basic_auth
    user_id UUID REFERENCES users(id),
    team_id UUID REFERENCES teams(id),
    node_types VARCHAR(100)[] DEFAULT '{}',
    data BYTEA NOT NULL, -- Encrypted
    iv BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Webhooks table
CREATE TABLE webhooks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID REFERENCES workflows(id) ON DELETE CASCADE,
    node_id VARCHAR(255) NOT NULL,
    path VARCHAR(255) UNIQUE NOT NULL,
    method VARCHAR(10) DEFAULT 'POST',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Scheduled workflows
CREATE TABLE scheduled_workflows (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID REFERENCES workflows(id) ON DELETE CASCADE,
    cron_expression VARCHAR(100) NOT NULL,
    timezone VARCHAR(50) DEFAULT 'UTC',
    is_active BOOLEAN DEFAULT true,
    last_run_at TIMESTAMP,
    next_run_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Variables table
CREATE TABLE variables (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    key VARCHAR(255) UNIQUE NOT NULL,
    value TEXT NOT NULL,
    type VARCHAR(50) DEFAULT 'string',
    is_secret BOOLEAN DEFAULT false,
    user_id UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tags table
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    color VARCHAR(7),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- API Keys table
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    key_hash VARCHAR(255) UNIQUE NOT NULL,
    scopes VARCHAR(50)[] DEFAULT '{}',
    expires_at TIMESTAMP,
    last_used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_workflows_user_active ON workflows(user_id, is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_executions_workflow_status ON executions(workflow_id, status);
CREATE INDEX idx_executions_created_at ON executions(created_at DESC);
CREATE INDEX idx_webhooks_path ON webhooks(path) WHERE is_active = true;
CREATE INDEX idx_scheduled_next_run ON scheduled_workflows(next_run_at) WHERE is_active = true;
```

---

## 🔌 API Architecture

### API Endpoints Structure

```yaml
Authentication:
  POST   /api/v1/auth/register
  POST   /api/v1/auth/login
  POST   /api/v1/auth/refresh
  POST   /api/v1/auth/logout
  GET    /api/v1/auth/me

Workflows:
  GET    /api/v1/workflows                 # List workflows
  POST   /api/v1/workflows                 # Create workflow
  GET    /api/v1/workflows/:id            # Get workflow
  PUT    /api/v1/workflows/:id            # Update workflow
  DELETE /api/v1/workflows/:id            # Delete workflow
  POST   /api/v1/workflows/:id/activate   # Activate
  POST   /api/v1/workflows/:id/deactivate # Deactivate
  POST   /api/v1/workflows/:id/execute    # Execute
  POST   /api/v1/workflows/:id/duplicate  # Duplicate
  GET    /api/v1/workflows/:id/executions # Get executions

Nodes:
  GET    /api/v1/nodes/types              # List node types
  POST   /api/v1/workflows/:id/nodes      # Add node
  PUT    /api/v1/nodes/:id                # Update node
  DELETE /api/v1/nodes/:id                # Delete node
  POST   /api/v1/nodes/:id/test           # Test node

Executions:
  GET    /api/v1/executions                # List executions
  GET    /api/v1/executions/:id           # Get execution
  POST   /api/v1/executions/:id/stop      # Stop execution
  POST   /api/v1/executions/:id/retry     # Retry execution
  DELETE /api/v1/executions/:id           # Delete execution

Credentials:
  GET    /api/v1/credentials               # List credentials
  POST   /api/v1/credentials               # Create credential
  GET    /api/v1/credentials/:id          # Get credential
  PUT    /api/v1/credentials/:id          # Update credential
  DELETE /api/v1/credentials/:id          # Delete credential
  POST   /api/v1/credentials/:id/test     # Test credential

Webhooks:
  GET    /api/v1/webhooks                  # List webhooks
  POST   /api/v1/webhooks                  # Create webhook
  DELETE /api/v1/webhooks/:id             # Delete webhook
  ANY    /webhook/:path                   # Webhook endpoint

Variables:
  GET    /api/v1/variables                 # List variables
  POST   /api/v1/variables                 # Create variable
  PUT    /api/v1/variables/:key           # Update variable
  DELETE /api/v1/variables/:key           # Delete variable

WebSocket:
  WS     /api/v1/ws                       # WebSocket connection
```

### API Response Format

```json
// Success Response
{
  "success": true,
  "data": {
    // Response data
  },
  "meta": {
    "timestamp": "2024-01-01T00:00:00Z",
    "request_id": "uuid"
  }
}

// Error Response
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input",
    "details": {
      "field": "email",
      "reason": "Invalid format"
    }
  },
  "meta": {
    "timestamp": "2024-01-01T00:00:00Z",
    "request_id": "uuid"
  }
}

// Paginated Response
{
  "success": true,
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

---

## ⚙️ Workflow Execution Engine

### Execution Flow Architecture

```
1. Trigger Phase
   ├── Manual Trigger (API call)
   ├── Webhook Trigger (HTTP endpoint)
   ├── Schedule Trigger (Cron job)
   └── Event Trigger (External event)

2. Preparation Phase
   ├── Load Workflow Definition
   ├── Validate Workflow Structure
   ├── Build Execution Graph (DAG)
   ├── Prepare Execution Context
   └── Initialize Data Flow

3. Execution Phase
   ├── Topological Sort Nodes
   ├── Execute Nodes in Order
   │   ├── Parallel Branch Execution
   │   ├── Sequential Node Execution
   │   ├── Conditional Branching
   │   └── Loop Processing
   ├── Data Transformation
   ├── Error Handling & Retries
   └── State Management

4. Completion Phase
   ├── Aggregate Results
   ├── Store Execution Data
   ├── Send Notifications
   └── Trigger Next Workflows
```

### Execution Engine Components

```go
// Workflow Executor
type WorkflowExecutor struct {
    ID           string
    WorkflowID   string
    Status       ExecutionStatus
    Mode         ExecutionMode
    Context      ExecutionContext
    DataFlow     DataFlowManager
    NodeStates   map[string]NodeState
    ErrorHandler ErrorHandler
    RetryPolicy  RetryPolicy
}

// Node Executor
type NodeExecutor interface {
    Execute(ctx context.Context, input NodeInput) (NodeOutput, error)
    Validate(params map[string]interface{}) error
    GetSchema() NodeSchema
}

// Execution Context
type ExecutionContext struct {
    WorkflowID   string
    ExecutionID  string
    Variables    map[string]interface{}
    Credentials  map[string]Credential
    Timezone     string
    StartTime    time.Time
}

// Data Flow Manager
type DataFlowManager struct {
    nodeOutputs map[string]interface{}
    connections []Connection
}
```

### Parallel Execution Strategy

```
                    ┌─────────┐
                    │  Start  │
                    └────┬────┘
                         │
                    ┌────▼────┐
                    │  HTTP   │
                    └────┬────┘
                         │
              ┌──────────┴──────────┐
              │                      │
         ┌────▼────┐            ┌───▼────┐
         │  Slack  │            │ Email  │  ← Parallel
         └────┬────┘            └───┬────┘
              │                      │
              └──────────┬──────────┘
                         │
                    ┌────▼────┐
                    │  Merge  │
                    └────┬────┘
                         │
                    ┌────▼────┐
                    │   End   │
                    └─────────┘
```

---

## 🔧 Node System Architecture

### Node Interface

```go
package nodes

type Node interface {
    // Core methods
    GetType() string
    GetCategory() string
    GetVersion() string
    
    // Execution
    Execute(ctx context.Context, input *NodeInput) (*NodeOutput, error)
    
    // Validation
    Validate(parameters map[string]interface{}) error
    
    // Schema
    GetSchema() *NodeSchema
    GetCredentialTypes() []string
}

type NodeInput struct {
    Data        map[string]interface{}
    Parameters  map[string]interface{}
    Credentials map[string]interface{}
    Context     ExecutionContext
}

type NodeOutput struct {
    Data     map[string]interface{}
    Error    error
    Metadata map[string]interface{}
}

type NodeSchema struct {
    Type        string
    Name        string
    Group       string
    Version     string
    Description string
    Icon        string
    Inputs      []IOSchema
    Outputs     []IOSchema
    Properties  []PropertySchema
    Credentials []CredentialSchema
}
```

### Node Categories

```yaml
Trigger Nodes: (Start workflows)
  - Manual Trigger
  - Webhook Trigger
  - Schedule Trigger (Cron)
  - Email Trigger (IMAP)
  - File Watcher
  - Database Trigger
  - Message Queue Trigger

Action Nodes: (Perform operations)
  - HTTP Request
  - Database Query
  - Send Email
  - Execute Command
  - API Call
  - File Operations

Transform Nodes: (Modify data)
  - Set/Update Data
  - Filter Items
  - Sort Items
  - Merge Data
  - Split Data
  - Aggregate
  - Format Data

Flow Control: (Control execution)
  - IF Condition
  - Switch
  - Loop
  - Wait
  - Stop and Error
  - Merge Branches

Integration Nodes: (Third-party)
  - Slack
  - GitHub
  - Google Sheets
  - OpenAI
  - Stripe
  - Twilio
  - AWS Services
  - 200+ more...
```

---

## 🎨 Frontend Architecture

### Component Structure

```typescript
// Workflow Canvas Component
interface WorkflowCanvasProps {
  workflowId: string;
  nodes: Node[];
  edges: Edge[];
  onNodesChange: (nodes: Node[]) => void;
  onEdgesChange: (edges: Edge[]) => void;
  onExecute: () => void;
}

// Node Component
interface NodeComponentProps {
  id: string;
  type: string;
  data: NodeData;
  selected: boolean;
  onUpdate: (data: NodeData) => void;
  onDelete: () => void;
}

// Properties Panel
interface PropertiesPanelProps {
  node: Node | null;
  onSave: (nodeId: string, data: any) => void;
  onTest: (nodeId: string) => void;
}
```

### State Management (Redux)

```typescript
// Store Structure
interface RootState {
  auth: AuthState;
  workflows: WorkflowsState;
  editor: EditorState;
  executions: ExecutionsState;
  ui: UIState;
}

interface EditorState {
  currentWorkflow: Workflow | null;
  nodes: Node[];
  edges: Edge[];
  selectedNode: string | null;
  isDirty: boolean;
  execution: {
    isRunning: boolean;
    currentNode: string | null;
    results: Record<string, any>;
  };
}
```

### Real-time Updates (WebSocket)

```typescript
// WebSocket Events
enum WSEventType {
  EXECUTION_STARTED = 'execution.started',
  EXECUTION_COMPLETED = 'execution.completed',
  NODE_EXECUTING = 'node.executing',
  NODE_COMPLETED = 'node.completed',
  WORKFLOW_UPDATED = 'workflow.updated'
}

// WebSocket Message
interface WSMessage {
  type: WSEventType;
  data: any;
  timestamp: string;
}
```

---

## 🔐 Security Architecture

### Authentication Flow

```
1. Registration
   ├── Email/Password validation
   ├── Password hashing (bcrypt)
   ├── Email verification
   └── Welcome workflow creation

2. Login
   ├── Credential validation
   ├── JWT token generation (15 min)
   ├── Refresh token (7 days)
   └── Session creation

3. Authorization
   ├── JWT validation
   ├── Role checking (RBAC)
   ├── Resource ownership
   └── Team permissions

4. Security Features
   ├── Rate limiting
   ├── IP whitelisting
   ├── 2FA support
   ├── API key management
   └── Audit logging
```

### Encryption Strategy

```yaml
Data at Rest:
  - Database: AES-256-GCM
  - Files: AES-256-CTR
  - Credentials: Separate encryption key
  
Data in Transit:
  - TLS 1.3 everywhere
  - Certificate pinning
  - Perfect forward secrecy

Secrets Management:
  - Environment variables: Encrypted
  - Credentials: Vault storage
  - API keys: Hashed + salted
  - Rotation: Every 30 days
```

---

## 🚀 Deployment Architecture

### Docker Compose (Development)

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: n8n_db
      POSTGRES_USER: n8n_user
      POSTGRES_PASSWORD: n8n_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  api:
    build: 
      context: .
      dockerfile: Dockerfile.api
    environment:
      DB_HOST: postgres
      REDIS_URL: redis://redis:6379
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis

  worker:
    build:
      context: .
      dockerfile: Dockerfile.worker
    environment:
      DB_HOST: postgres
      REDIS_URL: redis://redis:6379
    depends_on:
      - postgres
      - redis

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    environment:
      VITE_API_URL: http://localhost:8080

volumes:
  postgres_data:
  redis_data:
```

### Kubernetes (Production)

```yaml
# Deployment structure
Namespace: n8n-clone
├── Deployments
│   ├── api (3 replicas)
│   ├── worker (5 replicas)
│   ├── scheduler (1 replica)
│   └── websocket (3 replicas)
├── Services
│   ├── api-service (LoadBalancer)
│   ├── websocket-service
│   └── internal services
├── ConfigMaps
│   └── app-config
├── Secrets
│   └── app-secrets
└── HPA (Horizontal Pod Autoscaler)
    ├── api-hpa (CPU > 70%)
    └── worker-hpa (Queue depth)
```

---

## 📈 Implementation Roadmap

### Phase 1: Foundation (Weeks 1-2)
```
✓ Project setup and structure
✓ Database schema and migrations
✓ Basic domain models
✓ Authentication system
✓ Basic CRUD for workflows
✓ Simple API endpoints
```

### Phase 2: Core Engine (Weeks 3-4)
```
✓ Workflow execution engine
✓ Node system architecture
✓ Basic node types (HTTP, Webhook)
✓ Queue management
✓ Worker implementation
✓ Execution tracking
```

### Phase 3: Frontend (Weeks 5-6)
```
✓ React project setup
✓ Workflow canvas with ReactFlow
✓ Node library and drag-drop
✓ Properties panel
✓ Real-time execution updates
✓ Basic UI/UX
```

### Phase 4: Integrations (Weeks 7-8)
```
✓ 20+ core node types
✓ OAuth2 implementation
✓ Credential management
✓ Popular integrations (Slack, GitHub)
✓ Error handling
✓ Retry mechanisms
```

### Phase 5: Advanced Features (Weeks 9-10)
```
✓ Scheduled workflows
✓ Webhook management
✓ Variables and expressions
✓ Team collaboration
✓ API key management
✓ Workflow templates
```

### Phase 6: Production Ready (Weeks 11-12)
```
✓ Performance optimization
✓ Security hardening
✓ Monitoring and metrics
✓ Documentation
✓ Testing coverage
✓ Deployment setup
```

---

## 🎯 Key Implementation Files

### Critical Backend Files to Implement First

```
1. cmd/api/main.go                        # Entry point
2. internal/domain/workflow/entity.go      # Core domain
3. internal/domain/execution/entity.go     # Execution logic
4. internal/engine/executor/workflow_executor.go # Engine
5. internal/nodes/base_node.go            # Node system
6. internal/nodes/core/action/http_request.go # First node
7. internal/interfaces/http/rest/v1/controllers/workflow_controller.go
8. internal/infrastructure/persistence/postgres/workflow_repository.go
9. pkg/database/connection.go             # Database setup
10. configs/config.yaml                   # Configuration
```

### Critical Frontend Files to Implement First

```
1. src/app/App.tsx                        # Main app
2. src/app/store.ts                       # Redux store
3. src/features/editor/components/Canvas/WorkflowCanvas.tsx
4. src/features/editor/components/Nodes/BaseNode.tsx
5. src/features/editor/store/editorSlice.ts
6. src/shared/services/api.service.ts     # API client
7. src/shared/hooks/useWebSocket.ts       # WebSocket
```

---

## 📊 Success Metrics

### Performance Targets
- API Response: < 100ms
- Workflow Start: < 500ms
- Node Execution: < 1s per node
- WebSocket Latency: < 50ms

### Scale Targets
- Concurrent Users: 10,000+
- Workflows: 100,000+
- Executions/Day: 1,000,000+
- Nodes per Workflow: 500+

### Reliability Targets
- Uptime: 99.9%
- Data Loss: 0%
- Error Rate: < 0.1%
- Recovery Time: < 5 minutes

---

## 🔧 Development Commands

```bash
# Backend
make setup          # Setup development environment
make run-api        # Run API server
make run-worker     # Run worker
make test          # Run tests
make migrate       # Run migrations
make build         # Build binaries

# Frontend
npm install        # Install dependencies
npm run dev        # Development server
npm run build      # Production build
npm test          # Run tests
npm run lint      # Lint code

# Docker
docker-compose up -d     # Start all services
docker-compose logs -f   # View logs
docker-compose down      # Stop services

# Production
kubectl apply -f deployments/kubernetes/  # Deploy to k8s
kubectl get pods -n n8n-clone            # Check pods
kubectl logs -f deployment/api           # View logs
```

---

## 📚 Additional Resources

### Documentation to Create
1. API Documentation (OpenAPI/Swagger)
2. Node Development Guide
3. Deployment Guide
4. Security Best Practices
5. Performance Tuning Guide

### Monitoring Setup
1. Prometheus metrics
2. Grafana dashboards
3. ELK stack for logs
4. Jaeger for tracing
5. Sentry for errors

### Testing Strategy
1. Unit tests (80% coverage)
2. Integration tests
3. E2E tests with Cypress
4. Load testing with k6
5. Security testing

---

## ✅ Checklist for Production

- [ ] All critical features implemented
- [ ] Security measures in place
- [ ] Performance optimized
- [ ] Monitoring configured
- [ ] Documentation complete
- [ ] Tests passing (>80% coverage)
- [ ] CI/CD pipeline setup
- [ ] Backup strategy implemented
- [ ] Disaster recovery plan
- [ ] Load testing completed

---

## 🎉 Final Notes

This architecture provides:
1. **Complete n8n feature parity**
2. **Production-ready from day 1**
3. **Scalable to millions of executions**
4. **Maintainable and extensible**
5. **Enterprise-grade security**

Follow this architecture document as your single source of truth for building the n8n clone. Every component has been carefully designed to work together seamlessly while maintaining flexibility for future growth.

**Start with Phase 1 and progress systematically through each phase for best results!**
