# N8N Clone - Detailed Implementation Plan

## Table of Contents
1. [Core Architecture Implementation](#core-architecture-implementation)
2. [Detailed Database Design](#detailed-database-design)
3. [Workflow Execution Engine - Deep Dive](#workflow-execution-engine---deep-dive)
4. [Node Implementation Details](#node-implementation-details)
5. [API Implementation with Code Examples](#api-implementation-with-code-examples)
6. [Security Implementation](#security-implementation)
7. [Performance & Scaling](#performance--scaling)
8. [Testing Implementation](#testing-implementation)
9. [Deployment & DevOps](#deployment--devops)
10. [Migration Strategy](#migration-strategy)

## Core Architecture Implementation

### 1. Application Initialization

```go
// cmd/server/main.go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/yourusername/go-n8n/internal/api/routes"
    "github.com/yourusername/go-n8n/internal/config"
    "github.com/yourusername/go-n8n/internal/core/workflow"
    "github.com/yourusername/go-n8n/pkg/database"
    "github.com/yourusername/go-n8n/pkg/logger"
    "github.com/yourusername/go-n8n/pkg/queue"
    "go.uber.org/zap"
)

func main() {
    // Load configuration
    cfg := config.Load()
    
    // Initialize structured logger
    log := logger.New(cfg.LogLevel)
    defer log.Sync()
    
    // Initialize database
    db, err := database.NewConnection(cfg.Database)
    if err != nil {
        log.Fatal("Failed to connect to database", zap.Error(err))
    }
    defer db.Close()
    
    // Run migrations
    if err := database.Migrate(db); err != nil {
        log.Fatal("Failed to run migrations", zap.Error(err))
    }
    
    // Initialize Redis
    redisClient := queue.NewRedisClient(cfg.Redis)
    defer redisClient.Close()
    
    // Initialize queue system
    queueManager := queue.NewManager(redisClient, log)
    
    // Initialize workflow engine
    engine := workflow.NewEngine(db, redisClient, queueManager, log)
    
    // Start background workers
    workerCtx, cancelWorkers := context.WithCancel(context.Background())
    go queueManager.StartWorkers(workerCtx, cfg.WorkerCount)
    
    // Initialize Gin router
    router := setupRouter(cfg, db, redisClient, engine, log)
    
    // Create HTTP server
    srv := &http.Server{
        Addr:         fmt.Sprintf(":%d", cfg.Port),
        Handler:      router,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    
    // Start server in goroutine
    go func() {
        log.Info("Starting server", zap.Int("port", cfg.Port))
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal("Failed to start server", zap.Error(err))
        }
    }()
    
    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Info("Shutting down server...")
    
    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Stop accepting new jobs
    cancelWorkers()
    
    // Shutdown HTTP server
    if err := srv.Shutdown(ctx); err != nil {
        log.Error("Server forced to shutdown", zap.Error(err))
    }
    
    log.Info("Server shutdown complete")
}

func setupRouter(cfg *config.Config, db database.DB, redis *queue.RedisClient, 
    engine *workflow.Engine, log *zap.Logger) *gin.Engine {
    
    gin.SetMode(cfg.GinMode)
    router := gin.New()
    
    // Global middleware
    router.Use(
        middleware.Logger(log),
        middleware.Recovery(),
        middleware.RequestID(),
        middleware.CORS(cfg.CORS),
        middleware.RateLimit(redis, cfg.RateLimit),
    )
    
    // Initialize repositories
    repos := repository.NewRepositories(db)
    
    // Initialize services
    services := service.NewServices(repos, redis, engine, log)
    
    // Initialize handlers
    handlers := handlers.NewHandlers(services, log)
    
    // Setup routes
    routes.Setup(router, handlers, cfg)
    
    return router
}
```

### 2. Configuration Management

```go
// internal/config/config.go
package config

import (
    "time"
    "github.com/kelseyhightower/envconfig"
    "github.com/joho/godotenv"
)

type Config struct {
    // Server
    Port        int    `envconfig:"PORT" default:"8080"`
    GinMode     string `envconfig:"GIN_MODE" default:"debug"`
    Environment string `envconfig:"ENV" default:"development"`
    
    // Database
    Database DatabaseConfig
    
    // Redis
    Redis RedisConfig
    
    // JWT
    JWT JWTConfig
    
    // Encryption
    EncryptionKey string `envconfig:"ENCRYPTION_KEY" required:"true"`
    
    // CORS
    CORS CORSConfig
    
    // Rate Limiting
    RateLimit RateLimitConfig
    
    // Workers
    WorkerCount int `envconfig:"WORKER_COUNT" default:"10"`
    
    // Execution Limits
    MaxExecutionTime   time.Duration `envconfig:"MAX_EXECUTION_TIME" default:"30m"`
    MaxParallelNodes   int          `envconfig:"MAX_PARALLEL_NODES" default:"10"`
    MaxRetries         int          `envconfig:"MAX_RETRIES" default:"3"`
    RetryBackoff       time.Duration `envconfig:"RETRY_BACKOFF" default:"1m"`
    
    // Storage
    Storage StorageConfig
    
    // Webhook
    WebhookTimeout time.Duration `envconfig:"WEBHOOK_TIMEOUT" default:"30s"`
    WebhookMaxSize int64         `envconfig:"WEBHOOK_MAX_SIZE" default:"10485760"` // 10MB
    
    // Logging
    LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
    
    // Metrics
    MetricsEnabled bool   `envconfig:"METRICS_ENABLED" default:"true"`
    MetricsPort    int    `envconfig:"METRICS_PORT" default:"9090"`
}

type DatabaseConfig struct {
    Host            string        `envconfig:"DB_HOST" default:"localhost"`
    Port            int           `envconfig:"DB_PORT" default:"5432"`
    User            string        `envconfig:"DB_USER" required:"true"`
    Password        string        `envconfig:"DB_PASSWORD" required:"true"`
    Name            string        `envconfig:"DB_NAME" required:"true"`
    SSLMode         string        `envconfig:"DB_SSL_MODE" default:"disable"`
    MaxConnections  int           `envconfig:"DB_MAX_CONNECTIONS" default:"25"`
    MaxIdleConns    int           `envconfig:"DB_MAX_IDLE_CONNS" default:"5"`
    ConnMaxLifetime time.Duration `envconfig:"DB_CONN_MAX_LIFETIME" default:"5m"`
}

type RedisConfig struct {
    URL            string        `envconfig:"REDIS_URL" default:"redis://localhost:6379"`
    MaxRetries     int           `envconfig:"REDIS_MAX_RETRIES" default:"3"`
    PoolSize       int           `envconfig:"REDIS_POOL_SIZE" default:"10"`
    MinIdleConns   int           `envconfig:"REDIS_MIN_IDLE_CONNS" default:"2"`
    MaxConnAge     time.Duration `envconfig:"REDIS_MAX_CONN_AGE" default:"0"`
    ReadTimeout    time.Duration `envconfig:"REDIS_READ_TIMEOUT" default:"3s"`
    WriteTimeout   time.Duration `envconfig:"REDIS_WRITE_TIMEOUT" default:"3s"`
    PoolTimeout    time.Duration `envconfig:"REDIS_POOL_TIMEOUT" default:"4s"`
}

type JWTConfig struct {
    Secret         string        `envconfig:"JWT_SECRET" required:"true"`
    AccessExpiry   time.Duration `envconfig:"JWT_ACCESS_EXPIRY" default:"15m"`
    RefreshExpiry  time.Duration `envconfig:"JWT_REFRESH_EXPIRY" default:"168h"` // 7 days
    Issuer         string        `envconfig:"JWT_ISSUER" default:"go-n8n"`
}

type CORSConfig struct {
    AllowedOrigins   []string `envconfig:"CORS_ALLOWED_ORIGINS" default:"http://localhost:3000"`
    AllowedMethods   []string `envconfig:"CORS_ALLOWED_METHODS" default:"GET,POST,PUT,DELETE,OPTIONS"`
    AllowedHeaders   []string `envconfig:"CORS_ALLOWED_HEADERS" default:"Accept,Authorization,Content-Type,X-Request-ID"`
    ExposedHeaders   []string `envconfig:"CORS_EXPOSED_HEADERS" default:"X-Request-ID"`
    AllowCredentials bool     `envconfig:"CORS_ALLOW_CREDENTIALS" default:"true"`
    MaxAge           int      `envconfig:"CORS_MAX_AGE" default:"86400"`
}

type RateLimitConfig struct {
    Enabled       bool          `envconfig:"RATE_LIMIT_ENABLED" default:"true"`
    Requests      int           `envconfig:"RATE_LIMIT_REQUESTS" default:"100"`
    Duration      time.Duration `envconfig:"RATE_LIMIT_DURATION" default:"1m"`
    BurstRequests int           `envconfig:"RATE_LIMIT_BURST" default:"20"`
}

type StorageConfig struct {
    Type      string `envconfig:"STORAGE_TYPE" default:"local"` // local, s3, gcs
    LocalPath string `envconfig:"STORAGE_LOCAL_PATH" default:"./storage"`
    
    // S3 Configuration
    S3Bucket    string `envconfig:"S3_BUCKET"`
    S3Region    string `envconfig:"S3_REGION" default:"us-east-1"`
    S3Endpoint  string `envconfig:"S3_ENDPOINT"`
    S3AccessKey string `envconfig:"S3_ACCESS_KEY"`
    S3SecretKey string `envconfig:"S3_SECRET_KEY"`
}

func Load() *Config {
    _ = godotenv.Load()
    
    var cfg Config
    if err := envconfig.Process("", &cfg); err != nil {
        panic(err)
    }
    
    return &cfg
}

func (c *Config) GetDSN() string {
    return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        c.Database.Host, c.Database.Port, c.Database.User, 
        c.Database.Password, c.Database.Name, c.Database.SSLMode)
}
```

## Detailed Database Design

### Advanced Schema with Indexes and Constraints

```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Custom types
CREATE TYPE user_role AS ENUM ('admin', 'user', 'viewer');
CREATE TYPE execution_status AS ENUM ('waiting', 'running', 'success', 'error', 'cancelled', 'paused');
CREATE TYPE execution_mode AS ENUM ('manual', 'trigger', 'webhook', 'schedule', 'retry', 'test');
CREATE TYPE node_status AS ENUM ('waiting', 'running', 'success', 'error', 'skipped', 'paused');
CREATE TYPE credential_type AS ENUM ('api_key', 'oauth2', 'basic_auth', 'jwt', 'custom');
CREATE TYPE webhook_method AS ENUM ('GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS');

-- Users table with advanced features
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role user_role DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    email_verified BOOLEAN DEFAULT false,
    verification_token VARCHAR(255),
    reset_token VARCHAR(255),
    reset_token_expires TIMESTAMP,
    last_login_at TIMESTAMP,
    login_attempts INT DEFAULT 0,
    locked_until TIMESTAMP,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    CONSTRAINT email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
);

-- Indexes for users
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_verification_token ON users(verification_token) WHERE verification_token IS NOT NULL;
CREATE INDEX idx_users_reset_token ON users(reset_token) WHERE reset_token IS NOT NULL;
CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- Workflows table with versioning
CREATE TABLE workflows (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT false,
    settings JSONB DEFAULT '{}',
    nodes JSONB DEFAULT '[]', -- Store complete node structure
    connections JSONB DEFAULT '[]', -- Store connections
    static_data JSONB DEFAULT '{}', -- Persistent data between executions
    tags VARCHAR(50)[] DEFAULT '{}',
    version INT DEFAULT 1,
    error_workflow_id UUID REFERENCES workflows(id),
    timezone VARCHAR(50) DEFAULT 'UTC',
    timeout_seconds INT DEFAULT 3600,
    max_consecutive_failures INT DEFAULT 5,
    consecutive_failures INT DEFAULT 0,
    last_execution_at TIMESTAMP,
    next_execution_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    CONSTRAINT workflow_name_user_unique UNIQUE(name, user_id, deleted_at),
    CONSTRAINT timeout_positive CHECK (timeout_seconds > 0)
);

-- Indexes for workflows
CREATE INDEX idx_workflows_user_id ON workflows(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_workflows_is_active ON workflows(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_workflows_tags ON workflows USING GIN(tags) WHERE deleted_at IS NULL;
CREATE INDEX idx_workflows_next_execution ON workflows(next_execution_at) WHERE is_active = true AND deleted_at IS NULL;
CREATE INDEX idx_workflows_updated_at ON workflows(updated_at DESC);

-- Workflow versions for history
CREATE TABLE workflow_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID REFERENCES workflows(id) ON DELETE CASCADE,
    version INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    nodes JSONB NOT NULL,
    connections JSONB NOT NULL,
    settings JSONB DEFAULT '{}',
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT workflow_version_unique UNIQUE(workflow_id, version)
);

CREATE INDEX idx_workflow_versions_workflow_id ON workflow_versions(workflow_id, version DESC);

-- Executions table with partitioning support
CREATE TABLE executions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID REFERENCES workflows(id) ON DELETE CASCADE,
    workflow_version INT NOT NULL,
    status execution_status NOT NULL DEFAULT 'waiting',
    mode execution_mode NOT NULL DEFAULT 'manual',
    started_by UUID REFERENCES users(id),
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP,
    execution_time_ms INT,
    wait_till TIMESTAMP, -- For delayed executions
    retry_of UUID REFERENCES executions(id),
    retry_count INT DEFAULT 0,
    input_data JSONB DEFAULT '{}',
    output_data JSONB DEFAULT '{}',
    error_message TEXT,
    error_node_id VARCHAR(255),
    metadata JSONB DEFAULT '{}',
    webhook_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) PARTITION BY RANGE (created_at);

-- Create monthly partitions for executions
CREATE TABLE executions_2024_01 PARTITION OF executions
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');
-- Add more partitions as needed

-- Indexes for executions
CREATE INDEX idx_executions_workflow_id ON executions(workflow_id, created_at DESC);
CREATE INDEX idx_executions_status ON executions(status) WHERE status IN ('waiting', 'running');
CREATE INDEX idx_executions_mode ON executions(mode);
CREATE INDEX idx_executions_started_by ON executions(started_by);
CREATE INDEX idx_executions_created_at ON executions(created_at DESC);
CREATE INDEX idx_executions_wait_till ON executions(wait_till) WHERE wait_till IS NOT NULL AND status = 'waiting';

-- Node execution data with detailed metrics
CREATE TABLE execution_node_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    execution_id UUID REFERENCES executions(id) ON DELETE CASCADE,
    node_id VARCHAR(255) NOT NULL,
    node_type VARCHAR(100) NOT NULL,
    node_name VARCHAR(255),
    status node_status NOT NULL DEFAULT 'waiting',
    retry_count INT DEFAULT 0,
    input_data JSONB DEFAULT '{}',
    output_data JSONB DEFAULT '{}',
    error_message TEXT,
    error_details JSONB,
    execution_time_ms INT,
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    metadata JSONB DEFAULT '{}'
);

-- Indexes for node execution data
CREATE INDEX idx_execution_node_data_execution ON execution_node_data(execution_id);
CREATE INDEX idx_execution_node_data_node ON execution_node_data(node_id);
CREATE INDEX idx_execution_node_data_status ON execution_node_data(status);

-- Credentials with encryption
CREATE TABLE credentials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    type credential_type NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    node_types VARCHAR(100)[] DEFAULT '{}', -- Which node types can use this credential
    data BYTEA NOT NULL, -- Encrypted credential data
    iv BYTEA NOT NULL, -- Initialization vector for encryption
    is_shared BOOLEAN DEFAULT false,
    shared_with UUID[] DEFAULT '{}', -- User IDs this credential is shared with
    last_used_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    CONSTRAINT credential_name_user_unique UNIQUE(name, user_id, deleted_at)
);

-- Indexes for credentials
CREATE INDEX idx_credentials_user_id ON credentials(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_credentials_type ON credentials(type) WHERE deleted_at IS NULL;
CREATE INDEX idx_credentials_node_types ON credentials USING GIN(node_types) WHERE deleted_at IS NULL;
CREATE INDEX idx_credentials_shared_with ON credentials USING GIN(shared_with) WHERE is_shared = true AND deleted_at IS NULL;

-- Webhooks with path management
CREATE TABLE webhooks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID REFERENCES workflows(id) ON DELETE CASCADE,
    node_id VARCHAR(255) NOT NULL,
    path VARCHAR(255) UNIQUE NOT NULL,
    method webhook_method DEFAULT 'POST',
    is_active BOOLEAN DEFAULT true,
    requires_auth BOOLEAN DEFAULT false,
    auth_type VARCHAR(50), -- basic, bearer, api_key
    auth_value VARCHAR(255), -- Encrypted
    ip_whitelist INET[] DEFAULT '{}',
    headers_filter JSONB DEFAULT '{}',
    response_mode VARCHAR(50) DEFAULT 'immediately', -- immediately, last_node, specific_node
    response_node_id VARCHAR(255),
    response_data JSONB DEFAULT '{}',
    response_code INT DEFAULT 200,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT webhook_path_format CHECK (path ~ '^[a-zA-Z0-9-_/]+$')
);

-- Indexes for webhooks
CREATE INDEX idx_webhooks_workflow_id ON webhooks(workflow_id);
CREATE INDEX idx_webhooks_path ON webhooks(path) WHERE is_active = true;
CREATE INDEX idx_webhooks_node_id ON webhooks(node_id);

-- Scheduled workflows with cron
CREATE TABLE scheduled_workflows (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID REFERENCES workflows(id) ON DELETE CASCADE,
    cron_expression VARCHAR(100) NOT NULL,
    timezone VARCHAR(50) DEFAULT 'UTC',
    is_active BOOLEAN DEFAULT true,
    last_run_at TIMESTAMP,
    next_run_at TIMESTAMP,
    consecutive_failures INT DEFAULT 0,
    max_consecutive_failures INT DEFAULT 5,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT scheduled_workflow_unique UNIQUE(workflow_id)
);

-- Indexes for scheduled workflows
CREATE INDEX idx_scheduled_workflows_next_run ON scheduled_workflows(next_run_at) WHERE is_active = true;
CREATE INDEX idx_scheduled_workflows_workflow_id ON scheduled_workflows(workflow_id);

-- Workflow templates
CREATE TABLE workflow_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    icon VARCHAR(255),
    nodes JSONB NOT NULL,
    connections JSONB NOT NULL,
    credentials_template JSONB DEFAULT '{}',
    tags VARCHAR(50)[] DEFAULT '{}',
    is_public BOOLEAN DEFAULT false,
    created_by UUID REFERENCES users(id),
    usage_count INT DEFAULT 0,
    rating DECIMAL(2,1),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for templates
CREATE INDEX idx_workflow_templates_category ON workflow_templates(category) WHERE is_public = true;
CREATE INDEX idx_workflow_templates_tags ON workflow_templates USING GIN(tags) WHERE is_public = true;
CREATE INDEX idx_workflow_templates_rating ON workflow_templates(rating DESC) WHERE is_public = true;

-- Audit logs
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    old_data JSONB,
    new_data JSONB,
    ip_address INET,
    user_agent TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) PARTITION BY RANGE (created_at);

-- Create monthly partitions for audit logs
CREATE TABLE audit_logs_2024_01 PARTITION OF audit_logs
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

-- Indexes for audit logs
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id, created_at DESC);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);

-- API Keys for programmatic access
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    key_hash VARCHAR(255) UNIQUE NOT NULL,
    prefix VARCHAR(10) NOT NULL, -- For identification (first 10 chars)
    scopes VARCHAR(50)[] DEFAULT '{}',
    expires_at TIMESTAMP,
    last_used_at TIMESTAMP,
    usage_count INT DEFAULT 0,
    ip_whitelist INET[] DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT api_key_name_user_unique UNIQUE(name, user_id)
);

-- Indexes for API keys
CREATE INDEX idx_api_keys_user_id ON api_keys(user_id) WHERE is_active = true;
CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash) WHERE is_active = true;
CREATE INDEX idx_api_keys_prefix ON api_keys(prefix);

-- Trigger functions for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply updated_at triggers
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    
CREATE TRIGGER update_workflows_updated_at BEFORE UPDATE ON workflows
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    
CREATE TRIGGER update_credentials_updated_at BEFORE UPDATE ON credentials
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to clean up old execution data
CREATE OR REPLACE FUNCTION cleanup_old_executions(days_to_keep INT DEFAULT 30)
RETURNS INT AS $$
DECLARE
    deleted_count INT;
BEGIN
    DELETE FROM executions 
    WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '1 day' * days_to_keep
    AND status IN ('success', 'error', 'cancelled');
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;
```

## Workflow Execution Engine - Deep Dive

### Workflow Engine Core Implementation

```go
// internal/core/workflow/engine.go
package workflow

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
    "time"

    "github.com/google/uuid"
    "go.uber.org/zap"
)

type Engine struct {
    db           database.DB
    redis        *redis.Client
    queueManager *queue.Manager
    nodeRegistry *NodeRegistry
    executors    map[string]*Executor
    mu           sync.RWMutex
    logger       *zap.Logger
    metrics      *Metrics
    
    // Configuration
    maxParallelNodes   int
    maxExecutionTime   time.Duration
    defaultTimeout     time.Duration
    maxRetries         int
}

type Executor struct {
    engine       *Engine
    workflow     *models.Workflow
    execution    *models.Execution
    nodeStates   map[string]*NodeState
    dataFlow     *DataFlow
    context      context.Context
    cancel       context.CancelFunc
    logger       *zap.Logger
    
    // Execution control
    mu           sync.RWMutex
    paused       bool
    stopped      bool
    
    // Channels
    nodeChan     chan *NodeExecution
    resultChan   chan *NodeResult
    errorChan    chan error
}

type NodeState struct {
    NodeID      string
    Status      NodeStatus
    InputData   map[string]interface{}
    OutputData  map[string]interface{}
    Error       error
    RetryCount  int
    StartedAt   time.Time
    FinishedAt  time.Time
}

type DataFlow struct {
    data map[string]map[string]interface{}
    mu   sync.RWMutex
}

// NewEngine creates a new workflow execution engine
func NewEngine(db database.DB, redis *redis.Client, queueManager *queue.Manager, logger *zap.Logger) *Engine {
    return &Engine{
        db:               db,
        redis:            redis,
        queueManager:     queueManager,
        nodeRegistry:     NewNodeRegistry(),
        executors:        make(map[string]*Executor),
        logger:           logger,
        metrics:          NewMetrics(),
        maxParallelNodes: 10,
        maxExecutionTime: 30 * time.Minute,
        defaultTimeout:   5 * time.Minute,
        maxRetries:       3,
    }
}

// ExecuteWorkflow starts workflow execution
func (e *Engine) ExecuteWorkflow(ctx context.Context, workflowID string, inputData map[string]interface{}, mode ExecutionMode) (*models.Execution, error) {
    // Load workflow
    workflow, err := e.loadWorkflow(workflowID)
    if err != nil {
        return nil, fmt.Errorf("failed to load workflow: %w", err)
    }
    
    // Validate workflow
    if err := e.validateWorkflow(workflow); err != nil {
        return nil, fmt.Errorf("workflow validation failed: %w", err)
    }
    
    // Create execution record
    execution := &models.Execution{
        ID:              uuid.New().String(),
        WorkflowID:      workflowID,
        WorkflowVersion: workflow.Version,
        Status:          ExecutionStatusRunning,
        Mode:            mode,
        InputData:       inputData,
        StartedAt:       time.Now(),
    }
    
    if err := e.db.CreateExecution(execution); err != nil {
        return nil, fmt.Errorf("failed to create execution record: %w", err)
    }
    
    // Create executor
    execCtx, cancel := context.WithTimeout(ctx, e.maxExecutionTime)
    executor := &Executor{
        engine:      e,
        workflow:    workflow,
        execution:   execution,
        nodeStates:  make(map[string]*NodeState),
        dataFlow:    NewDataFlow(),
        context:     execCtx,
        cancel:      cancel,
        logger:      e.logger.With(zap.String("execution_id", execution.ID)),
        nodeChan:    make(chan *NodeExecution, e.maxParallelNodes),
        resultChan:  make(chan *NodeResult, e.maxParallelNodes),
        errorChan:   make(chan error, 1),
    }
    
    // Store executor
    e.mu.Lock()
    e.executors[execution.ID] = executor
    e.mu.Unlock()
    
    // Start execution in background
    go executor.run()
    
    // Send real-time update
    e.sendExecutionUpdate(execution, "started")
    
    return execution, nil
}

// Executor.run manages the workflow execution
func (ex *Executor) run() {
    defer func() {
        ex.cancel()
        ex.cleanup()
    }()
    
    ex.logger.Info("Starting workflow execution",
        zap.String("workflow_id", ex.workflow.ID),
        zap.String("workflow_name", ex.workflow.Name),
    )
    
    // Build execution graph
    graph, err := ex.buildExecutionGraph()
    if err != nil {
        ex.handleExecutionError(fmt.Errorf("failed to build execution graph: %w", err))
        return
    }
    
    // Initialize trigger node data
    if err := ex.initializeTriggerData(); err != nil {
        ex.handleExecutionError(fmt.Errorf("failed to initialize trigger data: %w", err))
        return
    }
    
    // Start worker pool
    var wg sync.WaitGroup
    for i := 0; i < ex.engine.maxParallelNodes; i++ {
        wg.Add(1)
        go ex.nodeWorker(&wg)
    }
    
    // Start result collector
    go ex.resultCollector()
    
    // Execute nodes in topological order
    for {
        // Get ready nodes (all dependencies satisfied)
        readyNodes := ex.getReadyNodes(graph)
        
        if len(readyNodes) == 0 {
            // Check if all nodes are completed
            if ex.allNodesCompleted() {
                ex.handleExecutionSuccess()
                break
            }
            
            // Check for deadlock
            if ex.hasDeadlock() {
                ex.handleExecutionError(fmt.Errorf("execution deadlock detected"))
                break
            }
            
            // Wait for nodes to complete
            time.Sleep(100 * time.Millisecond)
            continue
        }
        
        // Queue ready nodes for execution
        for _, node := range readyNodes {
            select {
            case <-ex.context.Done():
                ex.handleExecutionError(ex.context.Err())
                return
            case ex.nodeChan <- node:
                ex.updateNodeState(node.ID, NodeStatusRunning)
            }
        }
    }
    
    // Close channels and wait for workers
    close(ex.nodeChan)
    wg.Wait()
}

// nodeWorker processes nodes from the work queue
func (ex *Executor) nodeWorker(wg *sync.WaitGroup) {
    defer wg.Done()
    
    for nodeExec := range ex.nodeChan {
        ex.executeNode(nodeExec)
    }
}

// executeNode executes a single node
func (ex *Executor) executeNode(nodeExec *NodeExecution) {
    nodeState := ex.getNodeState(nodeExec.ID)
    nodeState.StartedAt = time.Now()
    
    ex.logger.Info("Executing node",
        zap.String("node_id", nodeExec.ID),
        zap.String("node_type", nodeExec.Type),
        zap.String("node_name", nodeExec.Name),
    )
    
    // Send real-time update
    ex.engine.sendNodeUpdate(ex.execution.ID, nodeExec.ID, "executing")
    
    // Get node implementation
    nodeImpl, err := ex.engine.nodeRegistry.GetNode(nodeExec.Type)
    if err != nil {
        ex.handleNodeError(nodeExec, fmt.Errorf("node type not found: %w", err))
        return
    }
    
    // Prepare input data
    inputData, err := ex.prepareNodeInput(nodeExec)
    if err != nil {
        ex.handleNodeError(nodeExec, fmt.Errorf("failed to prepare input: %w", err))
        return
    }
    
    // Load credentials if needed
    credentials, err := ex.loadNodeCredentials(nodeExec)
    if err != nil {
        ex.handleNodeError(nodeExec, fmt.Errorf("failed to load credentials: %w", err))
        return
    }
    
    // Create node context with timeout
    nodeTimeout := ex.getNodeTimeout(nodeExec)
    nodeCtx, cancel := context.WithTimeout(ex.context, nodeTimeout)
    defer cancel()
    
    // Execute node with retry logic
    var output *NodeOutput
    var lastError error
    
    for attempt := 0; attempt <= nodeExec.MaxRetries; attempt++ {
        if attempt > 0 {
            ex.logger.Info("Retrying node execution",
                zap.String("node_id", nodeExec.ID),
                zap.Int("attempt", attempt),
            )
            
            // Wait before retry with exponential backoff
            backoff := time.Duration(attempt) * time.Second * 2
            if backoff > 30*time.Second {
                backoff = 30 * time.Second
            }
            time.Sleep(backoff)
        }
        
        // Execute node
        output, lastError = nodeImpl.Execute(nodeCtx, &NodeInput{
            NodeID:      nodeExec.ID,
            Data:        inputData,
            Credentials: credentials,
            Parameters:  nodeExec.Parameters,
            Workflow:    ex.workflow,
            Execution:   ex.execution,
        })
        
        if lastError == nil {
            break
        }
        
        // Check if error is retryable
        if !isRetryableError(lastError) {
            break
        }
        
        nodeState.RetryCount = attempt + 1
    }
    
    // Handle execution result
    if lastError != nil {
        if nodeExec.ContinueOnFail {
            ex.logger.Warn("Node failed but continuing",
                zap.String("node_id", nodeExec.ID),
                zap.Error(lastError),
            )
            
            // Set empty output and continue
            output = &NodeOutput{
                Data: map[string]interface{}{
                    "error": lastError.Error(),
                },
            }
            nodeState.Status = NodeStatusError
        } else {
            ex.handleNodeError(nodeExec, lastError)
            return
        }
    } else {
        nodeState.Status = NodeStatusSuccess
    }
    
    // Store output
    nodeState.OutputData = output.Data
    nodeState.FinishedAt = time.Now()
    
    // Update data flow
    ex.dataFlow.SetNodeData(nodeExec.ID, output.Data)
    
    // Store node execution data
    if err := ex.storeNodeExecutionData(nodeExec, nodeState); err != nil {
        ex.logger.Error("Failed to store node execution data",
            zap.String("node_id", nodeExec.ID),
            zap.Error(err),
        )
    }
    
    // Send result
    ex.resultChan <- &NodeResult{
        NodeID: nodeExec.ID,
        Status: nodeState.Status,
        Data:   output.Data,
        Error:  lastError,
    }
    
    // Send real-time update
    ex.engine.sendNodeUpdate(ex.execution.ID, nodeExec.ID, "completed")
}

// buildExecutionGraph creates a DAG from workflow definition
func (ex *Executor) buildExecutionGraph() (*ExecutionGraph, error) {
    graph := NewExecutionGraph()
    
    // Parse nodes
    var nodes []models.Node
    if err := json.Unmarshal(ex.workflow.Nodes, &nodes); err != nil {
        return nil, fmt.Errorf("failed to parse nodes: %w", err)
    }
    
    // Add nodes to graph
    for _, node := range nodes {
        graph.AddNode(&node)
    }
    
    // Parse connections
    var connections []models.Connection
    if err := json.Unmarshal(ex.workflow.Connections, &connections); err != nil {
        return nil, fmt.Errorf("failed to parse connections: %w", err)
    }
    
    // Add connections to graph
    for _, conn := range connections {
        if err := graph.AddEdge(conn.SourceNodeID, conn.TargetNodeID); err != nil {
            return nil, fmt.Errorf("invalid connection: %w", err)
        }
    }
    
    // Validate graph
    if err := graph.Validate(); err != nil {
        return nil, fmt.Errorf("graph validation failed: %w", err)
    }
    
    // Check for cycles
    if graph.HasCycles() {
        return nil, fmt.Errorf("workflow contains cycles")
    }
    
    return graph, nil
}
```

### Execution Graph and Topological Sorting

```go
// internal/core/workflow/graph.go
package workflow

import (
    "fmt"
    "sync"
)

type ExecutionGraph struct {
    nodes     map[string]*GraphNode
    edges     map[string][]string  // adjacency list
    inDegree  map[string]int
    mu        sync.RWMutex
}

type GraphNode struct {
    ID         string
    Type       string
    Name       string
    Parameters map[string]interface{}
    Status     NodeStatus
}

func NewExecutionGraph() *ExecutionGraph {
    return &ExecutionGraph{
        nodes:    make(map[string]*GraphNode),
        edges:    make(map[string][]string),
        inDegree: make(map[string]int),
    }
}

func (g *ExecutionGraph) AddNode(node *models.Node) {
    g.mu.Lock()
    defer g.mu.Unlock()
    
    g.nodes[node.ID] = &GraphNode{
        ID:         node.ID,
        Type:       node.Type,
        Name:       node.Name,
        Parameters: node.Parameters,
        Status:     NodeStatusWaiting,
    }
    
    if _, exists := g.inDegree[node.ID]; !exists {
        g.inDegree[node.ID] = 0
    }
}

func (g *ExecutionGraph) AddEdge(from, to string) error {
    g.mu.Lock()
    defer g.mu.Unlock()
    
    if _, exists := g.nodes[from]; !exists {
        return fmt.Errorf("source node %s not found", from)
    }
    
    if _, exists := g.nodes[to]; !exists {
        return fmt.Errorf("target node %s not found", to)
    }
    
    g.edges[from] = append(g.edges[from], to)
    g.inDegree[to]++
    
    return nil
}

// TopologicalSort returns nodes in execution order
func (g *ExecutionGraph) TopologicalSort() ([]*GraphNode, error) {
    g.mu.RLock()
    defer g.mu.RUnlock()
    
    // Create a copy of inDegree
    inDegreeCopy := make(map[string]int)
    for k, v := range g.inDegree {
        inDegreeCopy[k] = v
    }
    
    var result []*GraphNode
    var queue []string
    
    // Find all nodes with no incoming edges
    for nodeID, degree := range inDegreeCopy {
        if degree == 0 {
            queue = append(queue, nodeID)
        }
    }
    
    // Process queue
    for len(queue) > 0 {
        // Dequeue
        nodeID := queue[0]
        queue = queue[1:]
        
        result = append(result, g.nodes[nodeID])
        
        // Reduce in-degree for neighbors
        for _, neighbor := range g.edges[nodeID] {
            inDegreeCopy[neighbor]--
            if inDegreeCopy[neighbor] == 0 {
                queue = append(queue, neighbor)
            }
        }
    }
    
    // Check if all nodes were processed
    if len(result) != len(g.nodes) {
        return nil, fmt.Errorf("graph contains cycles")
    }
    
    return result, nil
}

// HasCycles detects if the graph contains cycles
func (g *ExecutionGraph) HasCycles() bool {
    g.mu.RLock()
    defer g.mu.RUnlock()
    
    visited := make(map[string]bool)
    recStack := make(map[string]bool)
    
    for nodeID := range g.nodes {
        if !visited[nodeID] {
            if g.hasCyclesDFS(nodeID, visited, recStack) {
                return true
            }
        }
    }
    
    return false
}

func (g *ExecutionGraph) hasCyclesDFS(nodeID string, visited, recStack map[string]bool) bool {
    visited[nodeID] = true
    recStack[nodeID] = true
    
    for _, neighbor := range g.edges[nodeID] {
        if !visited[neighbor] {
            if g.hasCyclesDFS(neighbor, visited, recStack) {
                return true
            }
        } else if recStack[neighbor] {
            return true
        }
    }
    
    recStack[nodeID] = false
    return false
}

// GetReadyNodes returns nodes that are ready to execute
func (g *ExecutionGraph) GetReadyNodes(completed map[string]bool) []*GraphNode {
    g.mu.RLock()
    defer g.mu.RUnlock()
    
    var ready []*GraphNode
    
    for nodeID, node := range g.nodes {
        // Skip if already completed
        if completed[nodeID] {
            continue
        }
        
        // Check if all dependencies are satisfied
        allDepsCompleted := true
        for depNodeID := range g.nodes {
            // Check if depNodeID has an edge to nodeID
            for _, target := range g.edges[depNodeID] {
                if target == nodeID && !completed[depNodeID] {
                    allDepsCompleted = false
                    break
                }
            }
            if !allDepsCompleted {
                break
            }
        }
        
        if allDepsCompleted {
            ready = append(ready, node)
        }
    }
    
    return ready
}

// ParallelExecutionPlan creates groups of nodes that can execute in parallel
func (g *ExecutionGraph) ParallelExecutionPlan() ([][]string, error) {
    sorted, err := g.TopologicalSort()
    if err != nil {
        return nil, err
    }
    
    var plan [][]string
    processed := make(map[string]bool)
    
    for len(processed) < len(sorted) {
        var group []string
        
        for _, node := range sorted {
            if processed[node.ID] {
                continue
            }
            
            // Check if all dependencies are processed
            canExecute := true
            for depNodeID := range g.nodes {
                for _, target := range g.edges[depNodeID] {
                    if target == node.ID && !processed[depNodeID] {
                        canExecute = false
                        break
                    }
                }
                if !canExecute {
                    break
                }
            }
            
            if canExecute {
                group = append(group, node.ID)
            }
        }
        
        // Mark group as processed
        for _, nodeID := range group {
            processed[nodeID] = true
        }
        
        if len(group) > 0 {
            plan = append(plan, group)
        }
    }
    
    return plan, nil
}
```

## Node Implementation Details

### Base Node Interface and Common Implementations

```go
// internal/core/nodes/base.go
package nodes

import (
    "context"
    "encoding/json"
    "fmt"
)

type NodeType string

const (
    NodeTypeStart         NodeType = "start"
    NodeTypeHTTPRequest   NodeType = "http_request"
    NodeTypeWebhook       NodeType = "webhook"
    NodeTypeSchedule      NodeType = "schedule"
    NodeTypeFunction      NodeType = "function"
    NodeTypeCondition     NodeType = "condition"
    NodeTypeLoop          NodeType = "loop"
    NodeTypeMerge         NodeType = "merge"
    NodeTypeSplit         NodeType = "split"
    NodeTypeWait          NodeType = "wait"
    NodeTypeSetVariable   NodeType = "set_variable"
    NodeTypeDatabase      NodeType = "database"
    NodeTypeEmail         NodeType = "email"
    NodeTypeSlack         NodeType = "slack"
    NodeTypeTransform     NodeType = "transform"
)

type Node interface {
    Execute(ctx context.Context, input *NodeInput) (*NodeOutput, error)
    Validate(parameters map[string]interface{}) error
    GetSchema() *NodeSchema
    GetType() NodeType
}

type NodeInput struct {
    NodeID      string
    Data        map[string]interface{}
    Credentials map[string]interface{}
    Parameters  map[string]interface{}
    Workflow    *models.Workflow
    Execution   *models.Execution
}

type NodeOutput struct {
    Data     map[string]interface{}
    Metadata map[string]interface{}
}

type NodeSchema struct {
    Type        NodeType                 `json:"type"`
    Name        string                   `json:"name"`
    Group       string                   `json:"group"`
    Version     string                   `json:"version"`
    Description string                   `json:"description"`
    Inputs      []NodeIOSchema           `json:"inputs"`
    Outputs     []NodeIOSchema           `json:"outputs"`
    Properties  map[string]PropertySchema `json:"properties"`
    Credentials []CredentialRequirement  `json:"credentials"`
}

type NodeIOSchema struct {
    Name        string `json:"name"`
    Type        string `json:"type"`
    Description string `json:"description"`
    Required    bool   `json:"required"`
}

type PropertySchema struct {
    Type        string      `json:"type"`
    Description string      `json:"description"`
    Default     interface{} `json:"default"`
    Required    bool        `json:"required"`
    Options     []Option    `json:"options,omitempty"`
    Min         *float64    `json:"min,omitempty"`
    Max         *float64    `json:"max,omitempty"`
    Pattern     string      `json:"pattern,omitempty"`
}

type Option struct {
    Name  string      `json:"name"`
    Value interface{} `json:"value"`
}

type CredentialRequirement struct {
    Name     string   `json:"name"`
    Type     string   `json:"type"`
    Required bool     `json:"required"`
}

// BaseNode provides common functionality
type BaseNode struct {
    Type NodeType
}

func (n *BaseNode) GetType() NodeType {
    return n.Type
}

// Registry for all node types
type NodeRegistry struct {
    nodes map[NodeType]func() Node
    mu    sync.RWMutex
}

func NewNodeRegistry() *NodeRegistry {
    r := &NodeRegistry{
        nodes: make(map[NodeType]func() Node),
    }
    
    // Register built-in nodes
    r.Register(NodeTypeHTTPRequest, func() Node { return NewHTTPRequestNode() })
    r.Register(NodeTypeWebhook, func() Node { return NewWebhookNode() })
    r.Register(NodeTypeFunction, func() Node { return NewFunctionNode() })
    r.Register(NodeTypeCondition, func() Node { return NewConditionNode() })
    r.Register(NodeTypeLoop, func() Node { return NewLoopNode() })
    r.Register(NodeTypeMerge, func() Node { return NewMergeNode() })
    r.Register(NodeTypeDatabase, func() Node { return NewDatabaseNode() })
    r.Register(NodeTypeTransform, func() Node { return NewTransformNode() })
    
    return r
}

func (r *NodeRegistry) Register(nodeType NodeType, factory func() Node) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.nodes[nodeType] = factory
}

func (r *NodeRegistry) GetNode(nodeType string) (Node, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    factory, exists := r.nodes[NodeType(nodeType)]
    if !exists {
        return nil, fmt.Errorf("node type %s not registered", nodeType)
    }
    
    return factory(), nil
}
```

### HTTP Request Node Implementation

```go
// internal/core/nodes/http_request.go
package nodes

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
)

type HTTPRequestNode struct {
    BaseNode
    client *http.Client
}

func NewHTTPRequestNode() Node {
    return &HTTPRequestNode{
        BaseNode: BaseNode{Type: NodeTypeHTTPRequest},
        client: &http.Client{
            Timeout: 30 * time.Second,
            Transport: &http.Transport{
                MaxIdleConns:        100,
                MaxIdleConnsPerHost: 10,
                IdleConnTimeout:     90 * time.Second,
            },
        },
    }
}

func (n *HTTPRequestNode) Execute(ctx context.Context, input *NodeInput) (*NodeOutput, error) {
    // Extract parameters
    method := n.getStringParam(input.Parameters, "method", "GET")
    url := n.getStringParam(input.Parameters, "url", "")
    headers := n.getMapParam(input.Parameters, "headers")
    queryParams := n.getMapParam(input.Parameters, "queryParams")
    bodyType := n.getStringParam(input.Parameters, "bodyType", "json")
    body := input.Parameters["body"]
    
    // Authentication
    authType := n.getStringParam(input.Parameters, "authenticationType", "none")
    
    // Response handling
    responseType := n.getStringParam(input.Parameters, "responseType", "json")
    includeHeaders := n.getBoolParam(input.Parameters, "includeResponseHeaders", false)
    
    // Validate URL
    if url == "" {
        return nil, fmt.Errorf("URL is required")
    }
    
    // Replace variables in URL
    url = n.replaceVariables(url, input.Data)
    
    // Build query string
    if len(queryParams) > 0 {
        queryString := n.buildQueryString(queryParams, input.Data)
        if strings.Contains(url, "?") {
            url += "&" + queryString
        } else {
            url += "?" + queryString
        }
    }
    
    // Prepare request body
    var requestBody io.Reader
    if body != nil && method != "GET" && method != "HEAD" {
        bodyData, err := n.prepareBody(bodyType, body, input.Data)
        if err != nil {
            return nil, fmt.Errorf("failed to prepare request body: %w", err)
        }
        requestBody = bytes.NewReader(bodyData)
    }
    
    // Create request
    req, err := http.NewRequestWithContext(ctx, method, url, requestBody)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    // Set headers
    n.setHeaders(req, headers, bodyType, input.Data)
    
    // Apply authentication
    if err := n.applyAuthentication(req, authType, input.Credentials); err != nil {
        return nil, fmt.Errorf("failed to apply authentication: %w", err)
    }
    
    // Execute request with retry logic
    var resp *http.Response
    maxRetries := 3
    
    for attempt := 0; attempt <= maxRetries; attempt++ {
        if attempt > 0 {
            // Exponential backoff
            time.Sleep(time.Duration(attempt*2) * time.Second)
        }
        
        resp, err = n.client.Do(req)
        if err == nil {
            break
        }
        
        // Check if error is retryable
        if !n.isRetryableError(err) {
            break
        }
    }
    
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    // Read response body
    bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024)) // 10MB limit
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }
    
    // Prepare output
    output := &NodeOutput{
        Data: make(map[string]interface{}),
        Metadata: map[string]interface{}{
            "statusCode": resp.StatusCode,
            "statusText": resp.Status,
        },
    }
    
    // Parse response based on type
    switch responseType {
    case "json":
        var jsonData interface{}
        if err := json.Unmarshal(bodyBytes, &jsonData); err != nil {
            // If JSON parsing fails, return as string
            output.Data["body"] = string(bodyBytes)
        } else {
            output.Data["body"] = jsonData
        }
    case "text":
        output.Data["body"] = string(bodyBytes)
    case "binary":
        output.Data["body"] = bodyBytes
    default:
        output.Data["body"] = string(bodyBytes)
    }
    
    // Include headers if requested
    if includeHeaders {
        headers := make(map[string]string)
        for key, values := range resp.Header {
            headers[key] = strings.Join(values, ", ")
        }
        output.Data["headers"] = headers
    }
    
    // Check for HTTP errors
    if resp.StatusCode >= 400 {
        return output, fmt.Errorf("HTTP error: %s", resp.Status)
    }
    
    return output, nil
}

func (n *HTTPRequestNode) Validate(parameters map[string]interface{}) error {
    // Validate required parameters
    if url, ok := parameters["url"].(string); !ok || url == "" {
        return fmt.Errorf("URL is required")
    }
    
    // Validate method
    method := n.getStringParam(parameters, "method", "GET")
    validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
    valid := false
    for _, m := range validMethods {
        if method == m {
            valid = true
            break
        }
    }
    if !valid {
        return fmt.Errorf("invalid HTTP method: %s", method)
    }
    
    return nil
}

func (n *HTTPRequestNode) GetSchema() *NodeSchema {
    return &NodeSchema{
        Type:        NodeTypeHTTPRequest,
        Name:        "HTTP Request",
        Group:       "Core",
        Version:     "1.0.0",
        Description: "Make HTTP requests to any API or website",
        Inputs: []NodeIOSchema{
            {
                Name:        "main",
                Type:        "main",
                Description: "Main input",
                Required:    true,
            },
        },
        Outputs: []NodeIOSchema{
            {
                Name:        "main",
                Type:        "main",
                Description: "Response data",
                Required:    true,
            },
        },
        Properties: map[string]PropertySchema{
            "method": {
                Type:        "options",
                Description: "HTTP Method",
                Default:     "GET",
                Required:    true,
                Options: []Option{
                    {Name: "GET", Value: "GET"},
                    {Name: "POST", Value: "POST"},
                    {Name: "PUT", Value: "PUT"},
                    {Name: "DELETE", Value: "DELETE"},
                    {Name: "PATCH", Value: "PATCH"},
                    {Name: "HEAD", Value: "HEAD"},
                    {Name: "OPTIONS", Value: "OPTIONS"},
                },
            },
            "url": {
                Type:        "string",
                Description: "URL to request",
                Required:    true,
            },
            "authenticationType": {
                Type:        "options",
                Description: "Authentication type",
                Default:     "none",
                Options: []Option{
                    {Name: "None", Value: "none"},
                    {Name: "Basic Auth", Value: "basicAuth"},
                    {Name: "Bearer Token", Value: "bearerToken"},
                    {Name: "API Key", Value: "apiKey"},
                    {Name: "OAuth2", Value: "oauth2"},
                },
            },
            "headers": {
                Type:        "keyValue",
                Description: "Request headers",
            },
            "queryParams": {
                Type:        "keyValue",
                Description: "Query parameters",
            },
            "bodyType": {
                Type:        "options",
                Description: "Body content type",
                Default:     "json",
                Options: []Option{
                    {Name: "JSON", Value: "json"},
                    {Name: "Form Data", Value: "form"},
                    {Name: "Raw", Value: "raw"},
                    {Name: "Binary", Value: "binary"},
                },
            },
            "body": {
                Type:        "json",
                Description: "Request body",
            },
            "responseType": {
                Type:        "options",
                Description: "Response format",
                Default:     "json",
                Options: []Option{
                    {Name: "JSON", Value: "json"},
                    {Name: "Text", Value: "text"},
                    {Name: "Binary", Value: "binary"},
                },
            },
            "includeResponseHeaders": {
                Type:        "boolean",
                Description: "Include response headers in output",
                Default:     false,
            },
        },
        Credentials: []CredentialRequirement{
            {
                Name:     "httpBasicAuth",
                Type:     "basic_auth",
                Required: false,
            },
            {
                Name:     "httpBearerToken",
                Type:     "bearer_token",
                Required: false,
            },
            {
                Name:     "httpApiKey",
                Type:     "api_key",
                Required: false,
            },
        },
    }
}

// Helper methods
func (n *HTTPRequestNode) getStringParam(params map[string]interface{}, key, defaultValue string) string {
    if val, ok := params[key].(string); ok {
        return val
    }
    return defaultValue
}

func (n *HTTPRequestNode) getBoolParam(params map[string]interface{}, key string, defaultValue bool) bool {
    if val, ok := params[key].(bool); ok {
        return val
    }
    return defaultValue
}

func (n *HTTPRequestNode) getMapParam(params map[string]interface{}, key string) map[string]string {
    result := make(map[string]string)
    
    if val, ok := params[key].(map[string]interface{}); ok {
        for k, v := range val {
            if str, ok := v.(string); ok {
                result[k] = str
            }
        }
    }
    
    return result
}

func (n *HTTPRequestNode) replaceVariables(str string, data map[string]interface{}) string {
    // Simple variable replacement
    // In production, use a proper template engine
    for key, value := range data {
        placeholder := fmt.Sprintf("{{%s}}", key)
        str = strings.ReplaceAll(str, placeholder, fmt.Sprintf("%v", value))
    }
    return str
}

func (n *HTTPRequestNode) isRetryableError(err error) bool {
    // Check if error is retryable (network errors, timeouts, etc.)
    return strings.Contains(err.Error(), "timeout") ||
           strings.Contains(err.Error(), "connection refused") ||
           strings.Contains(err.Error(), "temporary failure")
}
```

## API Implementation with Code Examples

### Authentication Handler

```go
// internal/api/handlers/auth.go
package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "go.uber.org/zap"
    "golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
    authService    services.AuthService
    userService    services.UserService
    logger         *zap.Logger
    jwtSecret      []byte
    accessExpiry   time.Duration
    refreshExpiry  time.Duration
}

func NewAuthHandler(authService services.AuthService, userService services.UserService, 
    logger *zap.Logger, config *config.JWTConfig) *AuthHandler {
    return &AuthHandler{
        authService:   authService,
        userService:   userService,
        logger:        logger,
        jwtSecret:     []byte(config.Secret),
        accessExpiry:  config.AccessExpiry,
        refreshExpiry: config.RefreshExpiry,
    }
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.sendError(c, http.StatusBadRequest, "Invalid request", err)
        return
    }
    
    // Validate request
    if err := h.validateRegisterRequest(&req); err != nil {
        h.sendError(c, http.StatusBadRequest, "Validation failed", err)
        return
    }
    
    // Check if user exists
    existingUser, _ := h.userService.GetByEmail(c.Request.Context(), req.Email)
    if existingUser != nil {
        h.sendError(c, http.StatusConflict, "User already exists", nil)
        return
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
    if err != nil {
        h.sendError(c, http.StatusInternalServerError, "Failed to hash password", err)
        return
    }
    
    // Create user
    user := &models.User{
        ID:               uuid.New().String(),
        Email:            req.Email,
        PasswordHash:     string(hashedPassword),
        Name:             req.Name,
        Role:             models.RoleUser,
        IsActive:         true,
        EmailVerified:    false,
        VerificationToken: h.generateToken(),
    }
    
    if err := h.userService.Create(c.Request.Context(), user); err != nil {
        h.sendError(c, http.StatusInternalServerError, "Failed to create user", err)
        return
    }
    
    // Send verification email
    go h.authService.SendVerificationEmail(user)
    
    // Generate tokens
    tokenPair, err := h.generateTokenPair(user)
    if err != nil {
        h.sendError(c, http.StatusInternalServerError, "Failed to generate tokens", err)
        return
    }
    
    // Log successful registration
    h.logger.Info("User registered",
        zap.String("user_id", user.ID),
        zap.String("email", user.Email),
    )
    
    c.JSON(http.StatusCreated, gin.H{
        "message": "Registration successful",
        "user": gin.H{
            "id":    user.ID,
            "email": user.Email,
            "name":  user.Name,
            "role":  user.Role,
        },
        "tokens": tokenPair,
    })
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.sendError(c, http.StatusBadRequest, "Invalid request", err)
        return
    }
    
    // Get user
    user, err := h.userService.GetByEmail(c.Request.Context(), req.Email)
    if err != nil {
        h.sendError(c, http.StatusUnauthorized, "Invalid credentials", nil)
        return
    }
    
    // Check if user is locked
    if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
        h.sendError(c, http.StatusUnauthorized, "Account is locked", nil)
        return
    }
    
    // Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        // Increment login attempts
        h.userService.IncrementLoginAttempts(c.Request.Context(), user.ID)
        
        // Lock account after 5 failed attempts
        if user.LoginAttempts >= 4 {
            lockUntil := time.Now().Add(30 * time.Minute)
            h.userService.LockAccount(c.Request.Context(), user.ID, lockUntil)
            h.sendError(c, http.StatusUnauthorized, "Account locked due to too many failed attempts", nil)
            return
        }
        
        h.sendError(c, http.StatusUnauthorized, "Invalid credentials", nil)
        return
    }
    
    // Check if email is verified
    if !user.EmailVerified && h.requireEmailVerification() {
        h.sendError(c, http.StatusUnauthorized, "Email not verified", nil)
        return
    }
    
    // Reset login attempts
    h.userService.ResetLoginAttempts(c.Request.Context(), user.ID)
    
    // Update last login
    h.userService.UpdateLastLogin(c.Request.Context(), user.ID)
    
    // Generate tokens
    tokenPair, err := h.generateTokenPair(user)
    if err != nil {
        h.sendError(c, http.StatusInternalServerError, "Failed to generate tokens", err)
        return
    }
    
    // Create session
    session := &models.Session{
        ID:           uuid.New().String(),
        UserID:       user.ID,
        RefreshToken: tokenPair.RefreshToken,
        UserAgent:    c.Request.UserAgent(),
        IPAddress:    c.ClientIP(),
        ExpiresAt:    time.Now().Add(h.refreshExpiry),
    }
    
    if err := h.authService.CreateSession(c.Request.Context(), session); err != nil {
        h.logger.Error("Failed to create session",
            zap.String("user_id", user.ID),
            zap.Error(err),
        )
    }
    
    h.logger.Info("User logged in",
        zap.String("user_id", user.ID),
        zap.String("email", user.Email),
        zap.String("ip", c.ClientIP()),
    )
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Login successful",
        "user": gin.H{
            "id":    user.ID,
            "email": user.Email,
            "name":  user.Name,
            "role":  user.Role,
        },
        "tokens": tokenPair,
    })
}

// Refresh handles token refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
    var req RefreshRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.sendError(c, http.StatusBadRequest, "Invalid request", err)
        return
    }
    
    // Validate refresh token
    claims, err := h.validateRefreshToken(req.RefreshToken)
    if err != nil {
        h.sendError(c, http.StatusUnauthorized, "Invalid refresh token", err)
        return
    }
    
    // Get session
    session, err := h.authService.GetSessionByToken(c.Request.Context(), req.RefreshToken)
    if err != nil {
        h.sendError(c, http.StatusUnauthorized, "Session not found", err)
        return
    }
    
    // Check if session is expired
    if session.ExpiresAt.Before(time.Now()) {
        h.sendError(c, http.StatusUnauthorized, "Session expired", nil)
        return
    }
    
    // Get user
    user, err := h.userService.GetByID(c.Request.Context(), claims.UserID)
    if err != nil {
        h.sendError(c, http.StatusUnauthorized, "User not found", err)
        return
    }
    
    // Generate new access token
    accessToken, err := h.generateAccessToken(user)
    if err != nil {
        h.sendError(c, http.StatusInternalServerError, "Failed to generate access token", err)
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "access_token": accessToken,
        "expires_in":   int(h.accessExpiry.Seconds()),
    })
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
    // Get token from header
    token := h.getTokenFromHeader(c)
    if token == "" {
        h.sendError(c, http.StatusUnauthorized, "No token provided", nil)
        return
    }
    
    // Validate token
    claims, err := h.validateAccessToken(token)
    if err != nil {
        h.sendError(c, http.StatusUnauthorized, "Invalid token", err)
        return
    }
    
    // Revoke all sessions for user
    if err := h.authService.RevokeAllSessions(c.Request.Context(), claims.UserID); err != nil {
        h.logger.Error("Failed to revoke sessions",
            zap.String("user_id", claims.UserID),
            zap.Error(err),
        )
    }
    
    // Add token to blacklist (with expiry)
    h.authService.BlacklistToken(c.Request.Context(), token, time.Until(time.Unix(claims.ExpiresAt, 0)))
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Logout successful",
    })
}

// GetMe returns current user information
func (h *AuthHandler) GetMe(c *gin.Context) {
    // Get user from context (set by auth middleware)
    userID, exists := c.Get("user_id")
    if !exists {
        h.sendError(c, http.StatusUnauthorized, "User not found in context", nil)
        return
    }
    
    user, err := h.userService.GetByID(c.Request.Context(), userID.(string))
    if err != nil {
        h.sendError(c, http.StatusNotFound, "User not found", err)
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "user": gin.H{
            "id":             user.ID,
            "email":          user.Email,
            "name":           user.Name,
            "role":           user.Role,
            "email_verified": user.EmailVerified,
            "created_at":     user.CreatedAt,
        },
    })
}
```

## Security Implementation

### Encryption Service

```go
// internal/services/encryption.go
package services

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "io"
)

type EncryptionService struct {
    key []byte
}

func NewEncryptionService(key string) *EncryptionService {
    // Use SHA-256 to ensure key is correct length
    hash := sha256.Sum256([]byte(key))
    return &EncryptionService{
        key: hash[:],
    }
}

// Encrypt encrypts data using AES-256-GCM
func (s *EncryptionService) Encrypt(plaintext []byte) ([]byte, []byte, error) {
    block, err := aes.NewCipher(s.key)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to create cipher: %w", err)
    }
    
    // Create GCM cipher
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to create GCM: %w", err)
    }
    
    // Create nonce
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, nil, fmt.Errorf("failed to generate nonce: %w", err)
    }
    
    // Encrypt data
    ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
    
    return ciphertext, nonce, nil
}

// Decrypt decrypts data using AES-256-GCM
func (s *EncryptionService) Decrypt(ciphertext, nonce []byte) ([]byte, error) {
    block, err := aes.NewCipher(s.key)
    if err != nil {
        return nil, fmt.Errorf("failed to create cipher: %w", err)
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM: %w", err)
    }
    
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt: %w", err)
    }
    
    return plaintext, nil
}

// EncryptString encrypts a string and returns base64 encoded result
func (s *EncryptionService) EncryptString(plaintext string) (string, string, error) {
    ciphertext, nonce, err := s.Encrypt([]byte(plaintext))
    if err != nil {
        return "", "", err
    }
    
    return base64.StdEncoding.EncodeToString(ciphertext),
           base64.StdEncoding.EncodeToString(nonce), nil
}

// DecryptString decrypts a base64 encoded string
func (s *EncryptionService) DecryptString(ciphertext, nonce string) (string, error) {
    ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", fmt.Errorf("failed to decode ciphertext: %w", err)
    }
    
    nonceBytes, err := base64.StdEncoding.DecodeString(nonce)
    if err != nil {
        return "", fmt.Errorf("failed to decode nonce: %w", err)
    }
    
    plaintext, err := s.Decrypt(ciphertextBytes, nonceBytes)
    if err != nil {
        return "", err
    }
    
    return string(plaintext), nil
}

// EncryptCredentials encrypts credential data
func (s *EncryptionService) EncryptCredentials(credentials map[string]interface{}) ([]byte, []byte, error) {
    // Convert to JSON
    jsonData, err := json.Marshal(credentials)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to marshal credentials: %w", err)
    }
    
    return s.Encrypt(jsonData)
}

// DecryptCredentials decrypts credential data
func (s *EncryptionService) DecryptCredentials(ciphertext, nonce []byte) (map[string]interface{}, error) {
    plaintext, err := s.Decrypt(ciphertext, nonce)
    if err != nil {
        return nil, err
    }
    
    var credentials map[string]interface{}
    if err := json.Unmarshal(plaintext, &credentials); err != nil {
        return nil, fmt.Errorf("failed to unmarshal credentials: %w", err)
    }
    
    return credentials, nil
}
```

## Performance & Scaling

### Caching Strategy

```go
// internal/services/cache.go
package services

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/go-redis/redis/v8"
    "go.uber.org/zap"
)

type CacheService struct {
    redis  *redis.Client
    logger *zap.Logger
}

func NewCacheService(redis *redis.Client, logger *zap.Logger) *CacheService {
    return &CacheService{
        redis:  redis,
        logger: logger,
    }
}

// Workflow caching
func (s *CacheService) GetWorkflow(ctx context.Context, workflowID string) (*models.Workflow, error) {
    key := fmt.Sprintf("workflow:%s", workflowID)
    
    data, err := s.redis.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    var workflow models.Workflow
    if err := json.Unmarshal([]byte(data), &workflow); err != nil {
        return nil, err
    }
    
    return &workflow, nil
}

func (s *CacheService) SetWorkflow(ctx context.Context, workflow *models.Workflow, ttl time.Duration) error {
    key := fmt.Sprintf("workflow:%s", workflow.ID)
    
    data, err := json.Marshal(workflow)
    if err != nil {
        return err
    }
    
    return s.redis.Set(ctx, key, data, ttl).Err()
}

func (s *CacheService) InvalidateWorkflow(ctx context.Context, workflowID string) error {
    key := fmt.Sprintf("workflow:%s", workflowID)
    return s.redis.Del(ctx, key).Err()
}

// Execution caching for retry scenarios
func (s *CacheService) GetExecutionState(ctx context.Context, executionID string) (map[string]interface{}, error) {
    key := fmt.Sprintf("execution:state:%s", executionID)
    
    data, err := s.redis.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    var state map[string]interface{}
    if err := json.Unmarshal([]byte(data), &state); err != nil {
        return nil, err
    }
    
    return state, nil
}

func (s *CacheService) SetExecutionState(ctx context.Context, executionID string, 
    state map[string]interface{}, ttl time.Duration) error {
    
    key := fmt.Sprintf("execution:state:%s", executionID)
    
    data, err := json.Marshal(state)
    if err != nil {
        return err
    }
    
    return s.redis.Set(ctx, key, data, ttl).Err()
}

// Rate limiting
func (s *CacheService) CheckRateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
    pipe := s.redis.Pipeline()
    
    now := time.Now().Unix()
    windowStart := now - int64(window.Seconds())
    
    // Remove old entries
    pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))
    
    // Count current entries
    count := pipe.ZCard(ctx, key)
    
    // Add current request
    pipe.ZAdd(ctx, key, &redis.Z{
        Score:  float64(now),
        Member: fmt.Sprintf("%d", now),
    })
    
    // Set expiry
    pipe.Expire(ctx, key, window)
    
    _, err := pipe.Exec(ctx)
    if err != nil {
        return false, err
    }
    
    return count.Val() < int64(limit), nil
}

// Distributed locking
func (s *CacheService) AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
    return s.redis.SetNX(ctx, fmt.Sprintf("lock:%s", key), "1", ttl).Result()
}

func (s *CacheService) ReleaseLock(ctx context.Context, key string) error {
    return s.redis.Del(ctx, fmt.Sprintf("lock:%s", key)).Err()
}
```

### Load Balancing and Horizontal Scaling

```go
// internal/services/coordinator.go
package services

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/hashicorp/consul/api"
    "go.uber.org/zap"
)

type WorkerCoordinator struct {
    consul       *api.Client
    instanceID   string
    serviceName  string
    logger       *zap.Logger
    
    // Worker pool management
    workers      map[string]*Worker
    mu           sync.RWMutex
    
    // Load balancing
    loadBalancer LoadBalancer
}

type Worker struct {
    ID           string
    Address      string
    Port         int
    Load         int
    MaxCapacity  int
    LastHeartbeat time.Time
    Healthy      bool
}

type LoadBalancer interface {
    SelectWorker(workers []*Worker) *Worker
}

// LeastLoadedBalancer selects worker with least load
type LeastLoadedBalancer struct{}

func (b *LeastLoadedBalancer) SelectWorker(workers []*Worker) *Worker {
    if len(workers) == 0 {
        return nil
    }
    
    var selected *Worker
    minLoad := float64(1.0)
    
    for _, worker := range workers {
        if !worker.Healthy {
            continue
        }
        
        loadRatio := float64(worker.Load) / float64(worker.MaxCapacity)
        if loadRatio < minLoad {
            minLoad = loadRatio
            selected = worker
        }
    }
    
    return selected
}

func NewWorkerCoordinator(consulAddr, instanceID, serviceName string, logger *zap.Logger) (*WorkerCoordinator, error) {
    config := api.DefaultConfig()
    config.Address = consulAddr
    
    consul, err := api.NewClient(config)
    if err != nil {
        return nil, fmt.Errorf("failed to create consul client: %w", err)
    }
    
    return &WorkerCoordinator{
        consul:       consul,
        instanceID:   instanceID,
        serviceName:  serviceName,
        logger:       logger,
        workers:      make(map[string]*Worker),
        loadBalancer: &LeastLoadedBalancer{},
    }, nil
}

// RegisterWorker registers this instance as a worker
func (c *WorkerCoordinator) RegisterWorker(ctx context.Context, address string, port int) error {
    registration := &api.AgentServiceRegistration{
        ID:      c.instanceID,
        Name:    c.serviceName,
        Address: address,
        Port:    port,
        Check: &api.AgentServiceCheck{
            HTTP:                           fmt.Sprintf("http://%s:%d/health", address, port),
            Interval:                       "10s",
            Timeout:                        "5s",
            DeregisterCriticalServiceAfter: "30s",
        },
        Meta: map[string]string{
            "version":      "1.0.0",
            "max_capacity": "100",
        },
    }
    
    if err := c.consul.Agent().ServiceRegister(registration); err != nil {
        return fmt.Errorf("failed to register service: %w", err)
    }
    
    // Start heartbeat
    go c.heartbeat(ctx)
    
    return nil
}

// DiscoverWorkers finds all available workers
func (c *WorkerCoordinator) DiscoverWorkers(ctx context.Context) error {
    services, _, err := c.consul.Health().Service(c.serviceName, "", true, nil)
    if err != nil {
        return fmt.Errorf("failed to discover services: %w", err)
    }
    
    c.mu.Lock()
    defer c.mu.Unlock()
    
    // Update worker list
    newWorkers := make(map[string]*Worker)
    
    for _, service := range services {
        worker := &Worker{
            ID:           service.Service.ID,
            Address:      service.Service.Address,
            Port:         service.Service.Port,
            MaxCapacity:  100, // Default, can be from metadata
            LastHeartbeat: time.Now(),
            Healthy:      true,
        }
        
        // Preserve load information if worker exists
        if existing, ok := c.workers[worker.ID]; ok {
            worker.Load = existing.Load
        }
        
        newWorkers[worker.ID] = worker
    }
    
    c.workers = newWorkers
    
    c.logger.Info("Discovered workers",
        zap.Int("count", len(newWorkers)),
    )
    
    return nil
}

// AssignWork assigns work to least loaded worker
func (c *WorkerCoordinator) AssignWork(ctx context.Context, workItem interface{}) (*Worker, error) {
    c.mu.RLock()
    
    // Get healthy workers
    var healthyWorkers []*Worker
    for _, worker := range c.workers {
        if worker.Healthy {
            healthyWorkers = append(healthyWorkers, worker)
        }
    }
    c.mu.RUnlock()
    
    if len(healthyWorkers) == 0 {
        return nil, fmt.Errorf("no healthy workers available")
    }
    
    // Select worker using load balancer
    worker := c.loadBalancer.SelectWorker(healthyWorkers)
    if worker == nil {
        return nil, fmt.Errorf("failed to select worker")
    }
    
    // Update load
    c.mu.Lock()
    worker.Load++
    c.mu.Unlock()
    
    return worker, nil
}

// ReportCompletion reports work completion
func (c *WorkerCoordinator) ReportCompletion(workerID string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if worker, ok := c.workers[workerID]; ok {
        if worker.Load > 0 {
            worker.Load--
        }
    }
}

// heartbeat sends periodic heartbeats
func (c *WorkerCoordinator) heartbeat(ctx context.Context) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            // Update service health
            if err := c.consul.Agent().UpdateTTL(
                "service:"+c.instanceID,
                "alive",
                api.HealthPassing,
            ); err != nil {
                c.logger.Error("Failed to send heartbeat", zap.Error(err))
            }
        }
    }
}
```

## Testing Implementation

### Unit Tests Example

```go
// internal/services/workflow_service_test.go
package services

import (
    "context"
    "testing"
    "time"

    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestWorkflowService_Create(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockWorkflowRepository(ctrl)
    mockCache := mocks.NewMockCacheService(ctrl)
    mockValidator := mocks.NewMockValidator(ctrl)
    
    service := NewWorkflowService(mockRepo, mockCache, mockValidator)
    
    t.Run("successful creation", func(t *testing.T) {
        ctx := context.Background()
        workflow := &models.Workflow{
            Name:        "Test Workflow",
            Description: "Test Description",
            UserID:      "user-123",
            Nodes:       json.RawMessage(`[{"id":"node1","type":"start"}]`),
            Connections: json.RawMessage(`[]`),
        }
        
        mockValidator.EXPECT().
            ValidateWorkflow(workflow).
            Return(nil)
        
        mockRepo.EXPECT().
            Create(ctx, workflow).
            DoAndReturn(func(ctx context.Context, w *models.Workflow) error {
                w.ID = "workflow-123"
                w.CreatedAt = time.Now()
                return nil
            })
        
        mockCache.EXPECT().
            SetWorkflow(ctx, gomock.Any(), 5*time.Minute).
            Return(nil)
        
        err := service.Create(ctx, workflow)
        
        require.NoError(t, err)
        assert.Equal(t, "workflow-123", workflow.ID)
        assert.NotZero(t, workflow.CreatedAt)
    })
    
    t.Run("validation failure", func(t *testing.T) {
        ctx := context.Background()
        workflow := &models.Workflow{
            Name: "", // Empty name should fail validation
        }
        
        mockValidator.EXPECT().
            ValidateWorkflow(workflow).
            Return(ErrInvalidWorkflow)
        
        err := service.Create(ctx, workflow)
        
        assert.ErrorIs(t, err, ErrInvalidWorkflow)
    })
    
    t.Run("duplicate name", func(t *testing.T) {
        ctx := context.Background()
        workflow := &models.Workflow{
            Name:   "Existing Workflow",
            UserID: "user-123",
        }
        
        mockValidator.EXPECT().
            ValidateWorkflow(workflow).
            Return(nil)
        
        mockRepo.EXPECT().
            GetByName(ctx, workflow.Name, workflow.UserID).
            Return(&models.Workflow{}, nil)
        
        err := service.Create(ctx, workflow)
        
        assert.ErrorIs(t, err, ErrWorkflowExists)
    })
}

func TestWorkflowService_Execute(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockWorkflowRepository(ctrl)
    mockEngine := mocks.NewMockEngine(ctrl)
    mockQueue := mocks.NewMockQueueManager(ctrl)
    
    service := NewWorkflowService(mockRepo, nil, nil)
    service.SetEngine(mockEngine)
    service.SetQueue(mockQueue)
    
    t.Run("manual execution", func(t *testing.T) {
        ctx := context.Background()
        workflowID := "workflow-123"
        inputData := map[string]interface{}{
            "key": "value",
        }
        
        workflow := &models.Workflow{
            ID:       workflowID,
            IsActive: true,
        }
        
        execution := &models.Execution{
            ID:         "exec-123",
            WorkflowID: workflowID,
            Status:     ExecutionStatusRunning,
        }
        
        mockRepo.EXPECT().
            GetByID(ctx, workflowID).
            Return(workflow, nil)
        
        mockEngine.EXPECT().
            ExecuteWorkflow(ctx, workflowID, inputData, ExecutionModeManual).
            Return(execution, nil)
        
        result, err := service.Execute(ctx, workflowID, inputData, ExecutionModeManual)
        
        require.NoError(t, err)
        assert.Equal(t, execution.ID, result.ID)
    })
    
    t.Run("workflow not active", func(t *testing.T) {
        ctx := context.Background()
        workflowID := "workflow-123"
        
        workflow := &models.Workflow{
            ID:       workflowID,
            IsActive: false,
        }
        
        mockRepo.EXPECT().
            GetByID(ctx, workflowID).
            Return(workflow, nil)
        
        _, err := service.Execute(ctx, workflowID, nil, ExecutionModeManual)
        
        assert.ErrorIs(t, err, ErrWorkflowNotActive)
    })
}

func TestWorkflowService_ConcurrentExecutions(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockWorkflowRepository(ctrl)
    mockEngine := mocks.NewMockEngine(ctrl)
    
    service := NewWorkflowService(mockRepo, nil, nil)
    service.SetEngine(mockEngine)
    
    ctx := context.Background()
    workflowID := "workflow-123"
    
    workflow := &models.Workflow{
        ID:       workflowID,
        IsActive: true,
    }
    
    mockRepo.EXPECT().
        GetByID(ctx, workflowID).
        Return(workflow, nil).
        Times(10)
    
    mockEngine.EXPECT().
        ExecuteWorkflow(ctx, workflowID, gomock.Any(), ExecutionModeManual).
        Return(&models.Execution{
            ID: "exec-" + time.Now().String(),
        }, nil).
        Times(10)
    
    // Execute 10 workflows concurrently
    done := make(chan bool, 10)
    
    for i := 0; i < 10; i++ {
        go func(idx int) {
            _, err := service.Execute(ctx, workflowID, map[string]interface{}{
                "index": idx,
            }, ExecutionModeManual)
            
            assert.NoError(t, err)
            done <- true
        }(i)
    }
    
    // Wait for all executions to complete
    for i := 0; i < 10; i++ {
        <-done
    }
}
```

### Integration Tests

```go
// internal/api/handlers/workflow_handler_test.go
package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestWorkflowHandler_Integration(t *testing.T) {
    // Setup test database
    db := setupTestDatabase(t)
    defer db.Close()
    
    // Setup test services
    services := setupTestServices(db)
    
    // Create handler
    handler := NewWorkflowHandler(services.WorkflowService, services.ExecutionService)
    
    // Setup router
    router := gin.New()
    setupWorkflowRoutes(router, handler)
    
    t.Run("Create and Execute Workflow", func(t *testing.T) {
        // Create workflow
        createReq := CreateWorkflowRequest{
            Name:        "Test Workflow",
            Description: "Integration test workflow",
            Nodes: []NodeDefinition{
                {
                    ID:   "start",
                    Type: "start",
                    Name: "Start Node",
                },
                {
                    ID:   "http",
                    Type: "http_request",
                    Name: "HTTP Request",
                    Parameters: map[string]interface{}{
                        "method": "GET",
                        "url":    "https://api.example.com/test",
                    },
                },
            },
            Connections: []ConnectionDefinition{
                {
                    Source: "start",
                    Target: "http",
                },
            },
        }
        
        body, _ := json.Marshal(createReq)
        req := httptest.NewRequest(http.MethodPost, "/workflows", bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", "Bearer "+getTestToken())
        
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
        
        require.Equal(t, http.StatusCreated, w.Code)
        
        var createResp CreateWorkflowResponse
        err := json.Unmarshal(w.Body.Bytes(), &createResp)
        require.NoError(t, err)
        assert.NotEmpty(t, createResp.WorkflowID)
        
        workflowID := createResp.WorkflowID
        
        // Activate workflow
        activateReq := httptest.NewRequest(http.MethodPost, 
            "/workflows/"+workflowID+"/activate", nil)
        activateReq.Header.Set("Authorization", "Bearer "+getTestToken())
        
        w = httptest.NewRecorder()
        router.ServeHTTP(w, activateReq)
        
        require.Equal(t, http.StatusOK, w.Code)
        
        // Execute workflow
        executeReq := ExecuteWorkflowRequest{
            InputData: map[string]interface{}{
                "test": "data",
            },
        }
        
        body, _ = json.Marshal(executeReq)
        req = httptest.NewRequest(http.MethodPost, 
            "/workflows/"+workflowID+"/execute", bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", "Bearer "+getTestToken())
        
        w = httptest.NewRecorder()
        router.ServeHTTP(w, req)
        
        require.Equal(t, http.StatusOK, w.Code)
        
        var executeResp ExecuteWorkflowResponse
        err = json.Unmarshal(w.Body.Bytes(), &executeResp)
        require.NoError(t, err)
        assert.NotEmpty(t, executeResp.ExecutionID)
        
        // Get execution status
        statusReq := httptest.NewRequest(http.MethodGet, 
            "/executions/"+executeResp.ExecutionID, nil)
        statusReq.Header.Set("Authorization", "Bearer "+getTestToken())
        
        w = httptest.NewRecorder()
        router.ServeHTTP(w, statusReq)
        
        require.Equal(t, http.StatusOK, w.Code)
        
        var statusResp ExecutionStatusResponse
        err = json.Unmarshal(w.Body.Bytes(), &statusResp)
        require.NoError(t, err)
        assert.Contains(t, []string{"running", "success", "error"}, statusResp.Status)
    })
}
```

## Deployment & DevOps

### Docker Compose Setup

```yaml
# docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: go-n8n-postgres
    environment:
      POSTGRES_USER: n8n_user
      POSTGRES_PASSWORD: n8n_password
      POSTGRES_DB: n8n_db
      POSTGRES_INITDB_ARGS: "--encoding=UTF8 --lc-collate=C --lc-ctype=C"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations/init.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "5432:5432"
    networks:
      - n8n_network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U n8n_user -d n8n_db"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: go-n8n-redis
    command: redis-server --appendonly yes --maxmemory 512mb --maxmemory-policy allkeys-lru
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"
    networks:
      - n8n_network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  app:
    build:
      context: .
      dockerfile: Dockerfile
      args:
        - BUILD_VERSION=${BUILD_VERSION:-dev}
    container_name: go-n8n-app
    environment:
      - ENV=production
      - PORT=8080
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=n8n_user
      - DB_PASSWORD=n8n_password
      - DB_NAME=n8n_db
      - REDIS_URL=redis://redis:6379
      - JWT_SECRET=${JWT_SECRET}
      - ENCRYPTION_KEY=${ENCRYPTION_KEY}
      - CORS_ALLOWED_ORIGINS=${CORS_ALLOWED_ORIGINS:-http://localhost:3000}
    volumes:
      - ./storage:/app/storage
      - ./logs:/app/logs
    ports:
      - "8080:8080"
    networks:
      - n8n_network
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    container_name: go-n8n-nginx
    volumes:
      - ./docker/nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./docker/nginx/sites:/etc/nginx/sites-available
      - nginx_logs:/var/log/nginx
    ports:
      - "80:80"
      - "443:443"
    networks:
      - n8n_network
    depends_on:
      - app
    restart: unless-stopped

  prometheus:
    image: prom/prometheus:latest
    container_name: go-n8n-prometheus
    volumes:
      - ./docker/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--web.enable-lifecycle'
    ports:
      - "9090:9090"
    networks:
      - n8n_network
    restart: unless-stopped

  grafana:
    image: grafana/grafana:latest
    container_name: go-n8n-grafana
    volumes:
      - grafana_data:/var/lib/grafana
      - ./docker/grafana/provisioning:/etc/grafana/provisioning
    environment:
      - GF_SECURITY_ADMIN_USER=${GRAFANA_USER:-admin}
      - GF_SECURITY_ADMIN_PASSWORD=${GRAFANA_PASSWORD:-admin}
      - GF_INSTALL_PLUGINS=redis-datasource
    ports:
      - "3001:3000"
    networks:
      - n8n_network
    depends_on:
      - prometheus
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
  nginx_logs:
  prometheus_data:
  grafana_data:

networks:
  n8n_network:
    driver: bridge
```

### Production Dockerfile

```dockerfile
# Dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-w -s -X main.Version=${BUILD_VERSION} -X main.BuildTime=$(date -u '+%Y-%m-%d_%H:%M:%S')" \
    -a -installsuffix cgo \
    -o main cmd/server/main.go

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata curl

# Create non-root user
RUN addgroup -g 1001 -S n8n && \
    adduser -u 1001 -S n8n -G n8n

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy migration files
COPY --from=builder /app/migrations ./migrations

# Create necessary directories
RUN mkdir -p /app/storage /app/logs && \
    chown -R n8n:n8n /app

# Switch to non-root user
USER n8n

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run the binary
ENTRYPOINT ["./main"]
```

### Kubernetes Deployment

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-n8n
  namespace: default
  labels:
    app: go-n8n
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-n8n
  template:
    metadata:
      labels:
        app: go-n8n
    spec:
      containers:
      - name: go-n8n
        image: yourusername/go-n8n:latest
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 9090
          name: metrics
        env:
        - name: ENV
          value: "production"
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: go-n8n-secrets
              key: db-host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: go-n8n-secrets
              key: db-password
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: go-n8n-secrets
              key: jwt-secret
        - name: ENCRYPTION_KEY
          valueFrom:
            secretKeyRef:
              name: go-n8n-secrets
              key: encryption-key
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "1Gi"
            cpu: "1000m"
        livenessProbe:
          httpGet:
            path: /health
            port: http
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: http
          initialDelaySeconds: 5
          periodSeconds: 5
        volumeMounts:
        - name: storage
          mountPath: /app/storage
      volumes:
      - name: storage
        persistentVolumeClaim:
          claimName: go-n8n-storage

---
apiVersion: v1
kind: Service
metadata:
  name: go-n8n-service
  namespace: default
spec:
  type: LoadBalancer
  selector:
    app: go-n8n
  ports:
  - port: 80
    targetPort: 8080
    protocol: TCP
    name: http
  - port: 9090
    targetPort: 9090
    protocol: TCP
    name: metrics

---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: go-n8n-hpa
  namespace: default
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: go-n8n
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

## Migration Strategy

### Makefile for Development

```makefile
# Makefile
.PHONY: help build run test clean docker-up docker-down migrate

# Variables
APP_NAME := go-n8n
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Default target
help:
	@echo "Available targets:"
	@echo "  build       - Build the application"
	@echo "  run         - Run the application locally"
	@echo "  test        - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  lint        - Run linter"
	@echo "  fmt         - Format code"
	@echo "  clean       - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-up   - Start all services with docker-compose"
	@echo "  docker-down - Stop all services"
	@echo "  migrate-up  - Run database migrations"
	@echo "  migrate-down - Rollback database migrations"
	@echo "  migrate-create - Create new migration"
	@echo "  seed        - Seed database with test data"

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build $(LDFLAGS) -o $(GOBIN)/$(APP_NAME) cmd/server/main.go

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	@$(GOBIN)/$(APP_NAME)

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run --timeout=5m

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(GOBIN)
	@rm -f coverage.out coverage.html

# Docker targets
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):$(VERSION) --build-arg BUILD_VERSION=$(VERSION) .

docker-up:
	@echo "Starting services..."
	@docker-compose up -d

docker-down:
	@echo "Stopping services..."
	@docker-compose down

docker-logs:
	@docker-compose logs -f app

# Database migrations
migrate-up:
	@echo "Running migrations..."
	@migrate -path ./migrations -database "postgres://n8n_user:n8n_password@localhost:5432/n8n_db?sslmode=disable" up

migrate-down:
	@echo "Rolling back migrations..."
	@migrate -path ./migrations -database "postgres://n8n_user:n8n_password@localhost:5432/n8n_db?sslmode=disable" down 1

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ./migrations -seq $$name

# Seed database
seed:
	@echo "Seeding database..."
	@go run scripts/seed.go

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	@go mod download
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@cp .env.example .env
	@echo "Development environment ready!"

# Generate mocks
mocks:
	@echo "Generating mocks..."
	@mockgen -source=internal/repository/interfaces.go -destination=internal/mocks/repository_mocks.go -package=mocks
	@mockgen -source=internal/services/interfaces.go -destination=internal/mocks/service_mocks.go -package=mocks
	@echo "Mocks generated!"

# Run specific test
test-specific:
	@read -p "Enter test name pattern: " pattern; \
	go test -v -race -run $$pattern ./...

# Benchmark
benchmark:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Security check
security:
	@echo "Running security checks..."
	@gosec ./...
	@nancy go.sum

# API documentation
docs:
	@echo "Generating API documentation..."
	@swag init -g cmd/server/main.go

# Performance profiling
profile:
	@echo "Running profiling..."
	@go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=. ./internal/core/workflow
	@echo "Use 'go tool pprof cpu.prof' and 'go tool pprof mem.prof' to analyze"
```

This comprehensive plan provides:

1. **Complete architecture** with clean separation of concerns
2. **Detailed database schema** with indexes, constraints, and partitioning
3. **Full workflow execution engine** with DAG processing, parallel execution, and error handling
4. **Node system** with extensible architecture for custom nodes
5. **Security implementation** with encryption, JWT, and rate limiting
6. **Performance optimizations** including caching, load balancing, and horizontal scaling
7. **Complete testing strategy** with unit, integration, and E2E tests
8. **Production-ready deployment** with Docker, Kubernetes, and monitoring
9. **Development workflow** with Makefile for common tasks

The implementation focuses on scalability, maintainability, and production readiness, providing a solid foundation for building an n8n clone that can handle enterprise workloads.
