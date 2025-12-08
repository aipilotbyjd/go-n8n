# Complete n8n Clone Architecture - Full Production Structure

## 🎯 n8n Clone - Complete Implementation Structure

This is the **COMPLETE file structure** for building a production-ready n8n clone in Go. Every file listed here is necessary for full n8n functionality.

### ✅ What This Structure Includes:
- **200+ Node Types**: All n8n integrations (Slack, GitHub, OpenAI, databases, etc.)
- **Complete Workflow Engine**: Full execution, scheduling, and queue management
- **All API Endpoints**: 200+ REST endpoints matching n8n's API
- **Real-time Features**: WebSocket, live execution updates
- **Enterprise Features**: Teams, RBAC, audit logs, multi-tenancy ready
- **Production Ready**: Monitoring, scaling, security, deployment configs

### 📊 Structure Statistics:
- **Total Files**: ~500+ Go files
- **Node Types**: 200+ integration nodes
- **API Endpoints**: 200+ REST endpoints  
- **Database Tables**: 15+ core tables
- **Services**: 5 separate services (API, Worker, Scheduler, WebSocket, Webhook)

## 🏗️ Complete n8n Clone File Structure

### 📁 Full Project Structure for n8n Clone

```
go-n8n/
├── cmd/                                 # Application entry points
│   ├── api/                            
│   │   ├── main.go                     # REST API server entry
│   │   ├── wire.go                     # Dependency injection setup (Google Wire)
│   │   ├── wire_gen.go                 # Generated wire code
│   │   ├── config.go                   # API-specific configuration
│   │   └── server.go                   # HTTP server setup
│   ├── worker/                         
│   │   ├── main.go                     # Background job processor
│   │   ├── pools.go                    # Worker pool management
│   │   ├── handlers.go                 # Job handlers registration
│   │   ├── executor_worker.go          # Workflow execution worker
│   │   ├── webhook_worker.go           # Webhook processing worker
│   │   └── email_worker.go             # Email sending worker
│   ├── scheduler/                      
│   │   ├── main.go                     # Cron job scheduler
│   │   ├── jobs.go                     # Scheduled job definitions
│   │   ├── manager.go                  # Schedule management
│   │   ├── cron_parser.go              # Cron expression parser
│   │   └── timezone_handler.go         # Timezone management
│   ├── websocket/                      
│   │   ├── main.go                     # WebSocket server
│   │   ├── hub.go                      # Connection hub management
│   │   ├── client.go                   # WebSocket client handler
│   │   ├── broadcast.go                # Message broadcasting
│   │   └── handlers.go                 # Message handlers
│   ├── webhook/                        
│   │   ├── main.go                     # Webhook server (separate for scaling)
│   │   ├── router.go                   # Dynamic webhook routing
│   │   └── processor.go                # Webhook processing
│   └── migrate/                        
│       ├── main.go                     # Database migration tool
│       ├── seed.go                     # Database seeding tool
│       └── rollback.go                 # Migration rollback tool
│
├── internal/
│   ├── domain/                         # Core business logic (Clean Architecture - Entities)
│   │   ├── workflow/
│   │   │   ├── entity.go               # Workflow aggregate root
│   │   │   ├── value_objects.go        # NodePosition, WorkflowStatus, WorkflowTags
│   │   │   ├── repository.go           # Repository interface
│   │   │   ├── service.go              # Domain service for complex logic
│   │   │   ├── events.go               # WorkflowCreated, Updated, Deleted, Activated events
│   │   │   ├── errors.go               # Domain-specific errors
│   │   │   ├── specifications.go       # Business rule specifications
│   │   │   ├── factory.go              # Workflow factory for complex creation
│   │   │   ├── validator.go            # Workflow validation rules
│   │   │   ├── version.go              # Workflow versioning
│   │   │   ├── sharing.go              # Workflow sharing logic
│   │   │   └── template.go             # Workflow template management
│   │   │
│   │   ├── execution/
│   │   │   ├── entity.go               # Execution aggregate
│   │   │   ├── value_objects.go        # ExecutionStatus, ExecutionMode, ExecutionData
│   │   │   ├── repository.go           
│   │   │   ├── service.go              # Execution orchestration
│   │   │   ├── events.go               # ExecutionStarted, NodeExecuted, ExecutionCompleted events
│   │   │   ├── saga.go                 # Distributed transaction handling
│   │   │   ├── state_machine.go        # Execution state transitions
│   │   │   ├── context.go              # Execution context management
│   │   │   ├── data_flow.go            # Data flow between nodes
│   │   │   ├── retry_policy.go         # Retry policies and strategies
│   │   │   └── error_handler.go        # Error handling strategies
│   │   │
│   │   ├── node/
│   │   │   ├── entity.go               # Node entity
│   │   │   ├── types.go                # All node type definitions
│   │   │   ├── repository.go           
│   │   │   ├── registry.go             # Node type registry
│   │   │   ├── validator.go            # Node configuration validator
│   │   │   ├── metadata.go             # Node metadata management
│   │   │   ├── connection.go           # Node connection entity
│   │   │   ├── parameters.go           # Node parameters definition
│   │   │   ├── pin_data.go             # Pinned data for testing
│   │   │   └── categories.go           # Node categorization
│   │   │
│   │   ├── user/
│   │   │   ├── entity.go               # User aggregate root
│   │   │   ├── value_objects.go        # Email, Password, Username value objects
│   │   │   ├── repository.go           
│   │   │   ├── service.go              # User domain service
│   │   │   ├── events.go               # UserRegistered, Activated, PasswordChanged events
│   │   │   ├── permissions.go          # Permission value objects
│   │   │   ├── specifications.go       # User validation rules
│   │   │   ├── preferences.go          # User preferences
│   │   │   ├── api_key.go              # API key management
│   │   │   └── session.go              # User session entity
│   │   │
│   │   ├── credential/
│   │   │   ├── entity.go               # Credential aggregate
│   │   │   ├── value_objects.go        # EncryptedData, CredentialType
│   │   │   ├── repository.go           
│   │   │   ├── service.go              # Credential encryption service
│   │   │   ├── provider.go             # OAuth provider interface
│   │   │   ├── vault.go                # Secure storage interface
│   │   │   ├── types.go                # All credential types (OAuth2, API Key, etc)
│   │   │   ├── sharing.go              # Credential sharing between users
│   │   │   ├── testing.go              # Credential testing logic
│   │   │   └── rotation.go             # Credential rotation policy
│   │   │
│   │   ├── webhook/
│   │   │   ├── entity.go               # Webhook entity
│   │   │   ├── value_objects.go        # WebhookPath, WebhookMethod
│   │   │   ├── repository.go           
│   │   │   ├── service.go              # Webhook management service
│   │   │   ├── validator.go            # Webhook validation
│   │   │   ├── registry.go             # Webhook registry
│   │   │   └── events.go               # WebhookReceived, WebhookProcessed events
│   │   │
│   │   ├── schedule/
│   │   │   ├── entity.go               # Schedule entity
│   │   │   ├── value_objects.go        # CronExpression, Timezone
│   │   │   ├── repository.go           
│   │   │   ├── service.go              # Schedule management service
│   │   │   ├── calculator.go           # Next run calculation
│   │   │   └── events.go               # ScheduleCreated, ScheduleTriggered events
│   │   │
│   │   ├── tag/
│   │   │   ├── entity.go               # Tag entity
│   │   │   ├── repository.go           
│   │   │   └── service.go              # Tag management service
│   │   │
│   │   ├── variable/
│   │   │   ├── entity.go               # Environment variable entity
│   │   │   ├── repository.go           
│   │   │   ├── service.go              # Variable management service
│   │   │   └── resolver.go             # Variable resolution in expressions
│   │   │
│   │   ├── template/
│   │   │   ├── entity.go               # Workflow template entity
│   │   │   ├── repository.go           
│   │   │   ├── service.go              # Template management service
│   │   │   ├── categories.go           # Template categories
│   │   │   └── marketplace.go          # Template marketplace integration
│   │   │
│   │   ├── team/
│   │   │   ├── entity.go               # Team/Organization entity
│   │   │   ├── member.go               # Team member entity
│   │   │   ├── role.go                 # Role-based access control
│   │   │   ├── repository.go           
│   │   │   ├── invitation.go           # Team invitation management
│   │   │   ├── permissions.go          # Team permissions
│   │   │   └── events.go               # TeamCreated, MemberAdded events
│   │   │
│   │   └── common/
│   │       ├── aggregate.go            # Base aggregate root
│   │       ├── entity.go               # Base entity
│   │       ├── event.go                # Base domain event
│   │       ├── specification.go        # Specification pattern base
│   │       ├── result.go               # Result type for error handling
│   │       ├── pagination.go           # Pagination value objects
│   │       └── audit.go                # Audit trail base
│   │
│   ├── application/                    # Application Services (Use Cases)
│   │   ├── workflow/
│   │   │   ├── commands/               # Write operations (CQRS)
│   │   │   │   ├── create_workflow.go
│   │   │   │   ├── update_workflow.go
│   │   │   │   ├── delete_workflow.go
│   │   │   │   ├── activate_workflow.go
│   │   │   │   ├── deactivate_workflow.go
│   │   │   │   ├── duplicate_workflow.go
│   │   │   │   ├── share_workflow.go
│   │   │   │   ├── import_workflow.go
│   │   │   │   ├── export_workflow.go
│   │   │   │   ├── add_node.go
│   │   │   │   ├── update_node.go
│   │   │   │   ├── delete_node.go
│   │   │   │   ├── connect_nodes.go
│   │   │   │   ├── disconnect_nodes.go
│   │   │   │   └── pin_node_data.go
│   │   │   ├── queries/                # Read operations (CQRS)
│   │   │   │   ├── get_workflow.go
│   │   │   │   ├── list_workflows.go
│   │   │   │   ├── search_workflows.go
│   │   │   │   ├── get_workflow_stats.go
│   │   │   │   ├── get_workflow_versions.go
│   │   │   │   ├── get_workflow_nodes.go
│   │   │   │   ├── get_workflow_executions.go
│   │   │   │   └── get_workflow_metrics.go
│   │   │   ├── handlers.go             # Command/Query handlers
│   │   │   ├── validators.go           # Input validation
│   │   │   ├── mappers.go              # DTO to domain mapping
│   │   │   └── saga/                   # Long-running transactions
│   │   │       ├── workflow_deployment.go
│   │   │       └── workflow_migration.go
│   │   │
│   │   ├── execution/
│   │   │   ├── commands/
│   │   │   │   ├── execute_workflow.go
│   │   │   │   ├── stop_execution.go
│   │   │   │   ├── retry_execution.go
│   │   │   │   └── schedule_execution.go
│   │   │   ├── queries/
│   │   │   │   ├── get_execution.go
│   │   │   │   ├── list_executions.go
│   │   │   │   └── get_execution_logs.go
│   │   │   ├── handlers.go
│   │   │   └── orchestrator.go         # Execution orchestration
│   │   │
│   │   ├── auth/
│   │   │   ├── commands/
│   │   │   │   ├── register_user.go
│   │   │   │   ├── login_user.go
│   │   │   │   ├── refresh_token.go
│   │   │   │   ├── reset_password.go
│   │   │   │   └── verify_email.go
│   │   │   ├── queries/
│   │   │   │   └── get_user_profile.go
│   │   │   ├── handlers.go
│   │   │   ├── token_service.go        # JWT token management
│   │   │   └── oauth_service.go        # OAuth flow handling
│   │   │
│   │   └── shared/
│   │       ├── interfaces.go           # Shared interfaces
│   │       ├── transaction.go          # Unit of work pattern
│   │       ├── mediator.go             # CQRS mediator pattern
│   │       ├── event_bus.go            # Application event bus
│   │       └── saga_orchestrator.go    # Saga orchestration
│   │
│   ├── infrastructure/                 # External dependencies (Adapters)
│   │   ├── persistence/
│   │   │   ├── postgres/
│   │   │   │   ├── repositories/       # Repository implementations
│   │   │   │   │   ├── workflow_repository.go
│   │   │   │   │   ├── execution_repository.go
│   │   │   │   │   ├── user_repository.go
│   │   │   │   │   └── credential_repository.go
│   │   │   │   ├── migrations/         # SQL migration files
│   │   │   │   │   ├── 001_initial.up.sql
│   │   │   │   │   ├── 001_initial.down.sql
│   │   │   │   │   └── embed.go        # Embed migrations in binary
│   │   │   │   ├── connection.go       # Connection pool management
│   │   │   │   ├── transaction.go      # Transaction management
│   │   │   │   ├── querybuilder.go     # SQL query builder
│   │   │   │   └── indexes.go          # Index management
│   │   │   │
│   │   │   ├── redis/
│   │   │   │   ├── cache_repository.go # Generic cache operations
│   │   │   │   ├── session_store.go    # User session management
│   │   │   │   ├── rate_limiter.go     # Rate limiting implementation
│   │   │   │   ├── pub_sub.go          # Pub/Sub for events
│   │   │   │   ├── distributed_lock.go # Distributed locking
│   │   │   │   └── connection_pool.go  # Connection pool management
│   │   │   │
│   │   │   ├── elasticsearch/          # For logs and search
│   │   │   │   ├── log_repository.go
│   │   │   │   ├── search_repository.go
│   │   │   │   └── indexer.go
│   │   │   │
│   │   │   └── s3/                     # Object storage
│   │   │       ├── file_repository.go
│   │   │       └── presigned_urls.go
│   │   │
│   │   ├── messaging/
│   │   │   ├── kafka/                  # Event streaming
│   │   │   │   ├── producer.go
│   │   │   │   ├── consumer.go
│   │   │   │   ├── topics.go           # Topic management
│   │   │   │   └── serializer.go       # Message serialization
│   │   │   │
│   │   │   ├── rabbitmq/               # Message queue
│   │   │   │   ├── publisher.go
│   │   │   │   ├── consumer.go
│   │   │   │   ├── exchange.go         # Exchange configuration
│   │   │   │   └── dlq.go              # Dead letter queue
│   │   │   │
│   │   │   └── nats/                   # High-performance messaging
│   │   │       └── client.go
│   │   │
│   │   ├── security/
│   │   │   ├── jwt/
│   │   │   │   ├── generator.go        # JWT token generation
│   │   │   │   ├── validator.go        # JWT validation
│   │   │   │   └── claims.go           # Custom claims
│   │   │   │
│   │   │   ├── oauth/
│   │   │   │   ├── providers/          # OAuth providers
│   │   │   │   │   ├── google.go
│   │   │   │   │   ├── github.go
│   │   │   │   │   └── microsoft.go
│   │   │   │   └── flow.go             # OAuth flow implementation
│   │   │   │
│   │   │   ├── encryption/
│   │   │   │   ├── aes.go              # AES encryption
│   │   │   │   ├── rsa.go              # RSA encryption
│   │   │   │   ├── hasher.go           # Password hashing
│   │   │   │   └── vault_client.go     # HashiCorp Vault integration
│   │   │   │
│   │   │   └── rbac/                   # Role-based access control
│   │   │       ├── enforcer.go         # Policy enforcement
│   │   │       ├── policies.go         # Policy definitions
│   │   │       └── middleware.go       # RBAC middleware
│   │   │
│   │   ├── observability/
│   │   │   ├── tracing/
│   │   │   │   ├── jaeger.go           # Jaeger integration
│   │   │   │   ├── tracer.go           # OpenTelemetry tracer
│   │   │   │   └── middleware.go       # Tracing middleware
│   │   │   │
│   │   │   ├── metrics/
│   │   │   │   ├── prometheus.go       # Prometheus metrics
│   │   │   │   ├── collectors.go       # Custom collectors
│   │   │   │   └── middleware.go       # Metrics middleware
│   │   │   │
│   │   │   ├── logging/
│   │   │   │   ├── zap.go              # Zap logger setup
│   │   │   │   ├── context.go          # Context-aware logging
│   │   │   │   └── middleware.go       # Logging middleware
│   │   │   │
│   │   │   └── health/
│   │   │       ├── checker.go          # Health check implementation
│   │   │       └── probes.go           # Liveness/Readiness probes
│   │   │
│   │   └── external/                   # External service integrations
│   │       ├── email/
│   │       │   ├── smtp.go
│   │       │   ├── sendgrid.go
│   │       │   └── templates.go        # Email templates
│   │       │
│   │       ├── sms/
│   │       │   └── twilio.go
│   │       │
│   │       └── payment/
│   │           └── stripe.go
│   │
│   ├── interfaces/                     # Interface adapters (Ports)
│   │   ├── http/
│   │   │   ├── rest/
│   │   │   │   ├── v1/                 # API Version 1
│   │   │   │   │   ├── controllers/    # HTTP controllers
│   │   │   │   │   │   ├── workflow_controller.go
│   │   │   │   │   │   ├── execution_controller.go
│   │   │   │   │   │   ├── node_controller.go
│   │   │   │   │   │   ├── credential_controller.go
│   │   │   │   │   │   ├── webhook_controller.go
│   │   │   │   │   │   ├── auth_controller.go
│   │   │   │   │   │   ├── user_controller.go
│   │   │   │   │   │   ├── tag_controller.go
│   │   │   │   │   │   ├── variable_controller.go
│   │   │   │   │   │   ├── template_controller.go
│   │   │   │   │   │   ├── schedule_controller.go
│   │   │   │   │   │   ├── api_key_controller.go
│   │   │   │   │   │   ├── team_controller.go
│   │   │   │   │   │   ├── settings_controller.go
│   │   │   │   │   │   ├── metrics_controller.go
│   │   │   │   │   │   ├── health_controller.go
│   │   │   │   │   │   ├── audit_controller.go
│   │   │   │   │   │   ├── search_controller.go
│   │   │   │   │   │   ├── notification_controller.go
│   │   │   │   │   │   └── community_controller.go
│   │   │   │   │   ├── dto/            # Data transfer objects
│   │   │   │   │   │   ├── requests/
│   │   │   │   │   │   │   ├── workflow_requests.go
│   │   │   │   │   │   │   ├── execution_requests.go
│   │   │   │   │   │   │   ├── node_requests.go
│   │   │   │   │   │   │   ├── credential_requests.go
│   │   │   │   │   │   │   ├── auth_requests.go
│   │   │   │   │   │   │   └── pagination_request.go
│   │   │   │   │   │   └── responses/
│   │   │   │   │   │       ├── workflow_responses.go
│   │   │   │   │   │       ├── execution_responses.go
│   │   │   │   │   │       ├── node_responses.go
│   │   │   │   │   │       ├── credential_responses.go
│   │   │   │   │   │       ├── auth_responses.go
│   │   │   │   │   │       └── pagination_response.go
│   │   │   │   │   ├── validators/    # Request validation
│   │   │   │   │   ├── mappers/        # DTO mapping
│   │   │   │   │   ├── router.go       # Main route definitions
│   │   │   │   │   └── routes/         # Route groups
│   │   │   │   │       ├── workflow_routes.go
│   │   │   │   │       ├── execution_routes.go
│   │   │   │   │       ├── node_routes.go
│   │   │   │   │       ├── credential_routes.go
│   │   │   │   │       ├── webhook_routes.go
│   │   │   │   │       └── auth_routes.go
│   │   │   │   │
│   │   │   │   └── v2/                 # API Version 2
│   │   │   │
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go             # Authentication middleware
│   │   │   │   ├── cors.go             # CORS configuration
│   │   │   │   ├── rate_limit.go       # Rate limiting
│   │   │   │   ├── recovery.go         # Panic recovery
│   │   │   │   ├── timeout.go          # Request timeout
│   │   │   │   ├── compression.go      # Response compression
│   │   │   │   ├── cache.go            # Response caching
│   │   │   │   ├── security.go         # Security headers
│   │   │   │   └── request_id.go       # Request ID generation
│   │   │   │
│   │   │   ├── filters/                # Request/Response filters
│   │   │   └── swagger/                # API documentation
│   │   │       ├── docs.go
│   │   │       └── swagger.yaml
│   │   │
│   │   ├── grpc/
│   │   │   ├── proto/                  # Protocol buffer definitions
│   │   │   │   ├── workflow.proto
│   │   │   │   ├── execution.proto
│   │   │   │   └── common.proto
│   │   │   ├── services/               # gRPC service implementations
│   │   │   ├── interceptors/           # gRPC interceptors
│   │   │   └── gateway/                # gRPC-Gateway for REST
│   │   │
│   │   ├── websocket/
│   │   │   ├── hub.go                  # WebSocket connection hub
│   │   │   ├── client.go               # WebSocket client handler
│   │   │   ├── rooms.go                # Room-based broadcasting
│   │   │   ├── handlers/               # Message handlers
│   │   │   │   ├── execution.go
│   │   │   │   └── collaboration.go
│   │   │   └── auth.go                 # WebSocket authentication
│   │   │
│   │   └── graphql/
│   │       ├── schema/
│   │       │   ├── schema.graphql
│   │       │   └── types.go
│   │       ├── resolvers/
│   │       └── dataloaders/            # N+1 query optimization
│   │
│   ├── engine/                         # Workflow execution engine
│   │   ├── executor/
│   │   │   ├── workflow_executor.go    # Main workflow executor
│   │   │   ├── node_executor.go        # Individual node execution
│   │   │   ├── parallel_executor.go    # Parallel branch execution
│   │   │   ├── context.go              # Execution context management
│   │   │   ├── data_flow.go            # Data flow between nodes
│   │   │   ├── error_handler.go        # Error handling strategies
│   │   │   └── retry_policy.go         # Retry logic implementation
│   │   │
│   │   ├── scheduler/
│   │   │   ├── cron_scheduler.go       # Cron-based scheduling
│   │   │   ├── interval_scheduler.go   # Interval-based scheduling
│   │   │   ├── delayed_scheduler.go    # Delayed execution
│   │   │   ├── job_store.go            # Scheduled job storage
│   │   │   └── scheduler_interface.go  # Scheduler interface
│   │   │
│   │   ├── queue/
│   │   │   ├── queue_manager.go        # Queue management
│   │   │   ├── priority_queue.go       # Priority-based queuing
│   │   │   ├── worker_pool.go          # Worker pool management
│   │   │   ├── job.go                  # Job definition
│   │   │   └── dlq_handler.go          # Dead letter queue handling
│   │   │
│   │   ├── runtime/
│   │   │   ├── sandbox.go              # Sandboxed execution environment
│   │   │   ├── expression_evaluator.go # Expression evaluation
│   │   │   ├── variable_resolver.go    # Variable resolution
│   │   │   ├── function_registry.go    # Custom function registry
│   │   │   └── security_context.go     # Security context for execution
│   │   │
│   │   ├── graph/
│   │   │   ├── dag.go                  # Directed acyclic graph
│   │   │   ├── validator.go            # Graph validation
│   │   │   ├── optimizer.go            # Execution optimization
│   │   │   └── analyzer.go             # Static analysis
│   │   │
│   │   └── state/
│   │       ├── state_manager.go        # Execution state management
│   │       ├── checkpoint.go           # Execution checkpointing
│   │       └── recovery.go             # State recovery
│   │
│   └── nodes/                          # Node implementations (Plugin system)
│       ├── registry.go                 # Central node registry
│       ├── base_node.go                # Base node interface and implementation
│       ├── loader.go                   # Dynamic node loader
│       ├── validator.go                # Node validation
│       ├── metadata.go                 # Node metadata management
│       ├── executor.go                 # Node execution wrapper
│       ├── context.go                  # Node execution context
│       ├── categories.go               # Node categorization
│       │
│       ├── core/                       # Core built-in nodes
│       │   ├── trigger/                # Trigger nodes (can start workflows)
│       │   │   ├── start_node.go       # Manual trigger
│       │   │   ├── webhook_trigger.go  # Webhook trigger
│       │   │   ├── schedule_trigger.go # Cron schedule trigger
│       │   │   ├── email_trigger.go    # Email IMAP trigger
│       │   │   ├── interval_trigger.go # Interval trigger
│       │   │   ├── mqtt_trigger.go     # MQTT message trigger
│       │   │   ├── sse_trigger.go      # Server-sent events trigger
│       │   │   └── rss_trigger.go      # RSS feed trigger
│       │   │
│       │   ├── action/                 # Action nodes
│       │   │   ├── http_request.go     # HTTP/REST API calls
│       │   │   ├── webhook_response.go # Webhook response node
│       │   │   ├── graphql_request.go  # GraphQL queries
│       │   │   ├── soap_request.go     # SOAP API calls
│       │   │   ├── email_send.go       # Send emails (SMTP)
│       │   │   ├── execute_command.go  # Execute system commands
│       │   │   ├── ssh_command.go      # SSH command execution
│       │   │   └── wait_node.go        # Wait/delay execution
│       │   │
│       │   ├── function/               # Code execution nodes
│       │   │   ├── javascript_node.go  # JavaScript code execution
│       │   │   ├── python_node.go      # Python code execution
│       │   │   ├── code_node.go        # Generic code node
│       │   │   └── expression_node.go  # Expression evaluation
│       │   │
│       │   ├── transform/              # Data transformation nodes
│       │   │   ├── set_node.go         # Set/modify data
│       │   │   ├── rename_keys.go      # Rename object keys
│       │   │   ├── filter_items.go     # Filter array items
│       │   │   ├── sort_items.go       # Sort array items
│       │   │   ├── limit_items.go      # Limit number of items
│       │   │   ├── aggregate_items.go  # Aggregate data
│       │   │   ├── split_in_batches.go # Split data into batches
│       │   │   ├── merge_items.go      # Merge multiple inputs
│       │   │   ├── remove_duplicates.go # Remove duplicate items
│       │   │   ├── html_extract.go     # Extract data from HTML
│       │   │   ├── xml_parser.go       # Parse XML data
│       │   │   ├── csv_parser.go       # Parse CSV data
│       │   │   ├── json_parser.go      # Parse JSON strings
│       │   │   ├── date_time.go        # Date/time operations
│       │   │   ├── crypto_node.go      # Encryption/hashing
│       │   │   └── compress_node.go    # Compress/decompress data
│       │   │
│       │   ├── flow/                   # Flow control nodes
│       │   │   ├── if_node.go          # Conditional branching
│       │   │   ├── switch_node.go      # Multiple conditions
│       │   │   ├── loop_node.go        # Loop over items
│       │   │   ├── merge_node.go       # Merge branches
│       │   │   ├── split_node.go       # Split execution
│       │   │   ├── error_trigger.go    # Error workflow trigger
│       │   │   ├── stop_and_error.go   # Stop with error
│       │   │   └── no_op.go            # No operation node
│       │   │
│       │   └── files/                  # File operation nodes
│       │       ├── read_file.go        # Read files from disk
│       │       ├── write_file.go       # Write files to disk
│       │       ├── delete_file.go      # Delete files
│       │       ├── move_file.go        # Move/rename files
│       │       ├── ftp_node.go         # FTP operations
│       │       └── binary_data.go      # Binary data handling
│       │
│       ├── integrations/               # Third-party integrations
│       │   ├── communication/
│       │   │   ├── slack/
│       │   │   │   ├── slack_node.go   # Main Slack node
│       │   │   │   ├── slack_trigger.go # Slack event trigger
│       │   │   │   ├── operations.go    # All Slack operations
│       │   │   │   └── auth.go          # Slack OAuth
│       │   │   ├── discord/
│       │   │   │   ├── discord_node.go
│       │   │   │   └── discord_trigger.go
│       │   │   ├── telegram/
│       │   │   │   ├── telegram_node.go
│       │   │   │   └── telegram_trigger.go
│       │   │   ├── microsoft_teams/
│       │   │   │   ├── teams_node.go
│       │   │   │   └── teams_trigger.go
│       │   │   ├── whatsapp/
│       │   │   │   └── whatsapp_business.go
│       │   │   ├── twilio/
│       │   │   │   ├── twilio_sms.go
│       │   │   │   └── twilio_voice.go
│       │   │   └── email/
│       │   │       ├── gmail_node.go
│       │   │       ├── outlook_node.go
│       │   │       └── sendgrid_node.go
│       │   │
│       │   ├── databases/
│       │   │   ├── postgres/
│       │   │   │   ├── postgres_node.go
│       │   │   │   ├── postgres_trigger.go
│       │   │   │   └── operations.go
│       │   │   ├── mysql/
│       │   │   │   ├── mysql_node.go
│       │   │   │   └── mysql_trigger.go
│       │   │   ├── mongodb/
│       │   │   │   ├── mongodb_node.go
│       │   │   │   └── operations.go
│       │   │   ├── redis/
│       │   │   │   └── redis_node.go
│       │   │   ├── elasticsearch/
│       │   │   │   └── elasticsearch_node.go
│       │   │   ├── supabase/
│       │   │   │   └── supabase_node.go
│       │   │   └── airtable/
│       │   │       ├── airtable_node.go
│       │   │       └── airtable_trigger.go
│       │   │
│       │   ├── cloud/
│       │   │   ├── aws/
│       │   │   │   ├── s3_node.go
│       │   │   │   ├── lambda_node.go
│       │   │   │   ├── sqs_node.go
│       │   │   │   ├── sns_node.go
│       │   │   │   ├── dynamodb_node.go
│       │   │   │   └── ec2_node.go
│       │   │   ├── gcp/
│       │   │   │   ├── gcs_node.go      # Google Cloud Storage
│       │   │   │   ├── bigquery_node.go
│       │   │   │   ├── pubsub_node.go
│       │   │   │   └── functions_node.go
│       │   │   └── azure/
│       │   │       ├── blob_storage.go
│       │   │       ├── cosmos_db.go
│       │   │       └── functions.go
│       │   │
│       │   ├── devops/
│       │   │   ├── github/
│       │   │   │   ├── github_node.go
│       │   │   │   ├── github_trigger.go
│       │   │   │   └── operations.go   # Issues, PRs, Actions, etc
│       │   │   ├── gitlab/
│       │   │   │   ├── gitlab_node.go
│       │   │   │   └── gitlab_trigger.go
│       │   │   ├── bitbucket/
│       │   │   │   └── bitbucket_node.go
│       │   │   ├── jenkins/
│       │   │   │   └── jenkins_node.go
│       │   │   ├── docker/
│       │   │   │   └── docker_node.go
│       │   │   └── kubernetes/
│       │   │       └── kubernetes_node.go
│       │   │
│       │   ├── crm/
│       │   │   ├── salesforce/
│       │   │   │   ├── salesforce_node.go
│       │   │   │   └── operations.go
│       │   │   ├── hubspot/
│       │   │   │   ├── hubspot_node.go
│       │   │   │   └── hubspot_trigger.go
│       │   │   ├── pipedrive/
│       │   │   │   └── pipedrive_node.go
│       │   │   └── zoho/
│       │   │       └── zoho_crm.go
│       │   │
│       │   ├── productivity/
│       │   │   ├── google_workspace/
│       │   │   │   ├── google_drive.go
│       │   │   │   ├── google_sheets.go
│       │   │   │   ├── google_docs.go
│       │   │   │   ├── google_calendar.go
│       │   │   │   └── google_forms.go
│       │   │   ├── microsoft_365/
│       │   │   │   ├── onedrive.go
│       │   │   │   ├── excel.go
│       │   │   │   └── sharepoint.go
│       │   │   ├── notion/
│       │   │   │   ├── notion_node.go
│       │   │   │   └── notion_trigger.go
│       │   │   ├── trello/
│       │   │   │   ├── trello_node.go
│       │   │   │   └── trello_trigger.go
│       │   │   ├── asana/
│       │   │   │   └── asana_node.go
│       │   │   └── jira/
│       │   │       ├── jira_node.go
│       │   │       └── jira_trigger.go
│       │   │
│       │   ├── marketing/
│       │   │   ├── mailchimp/
│       │   │   │   └── mailchimp_node.go
│       │   │   ├── activecampaign/
│       │   │   │   └── activecampaign_node.go
│       │   │   ├── facebook/
│       │   │   │   └── facebook_ads.go
│       │   │   ├── google_ads/
│       │   │   │   └── google_ads_node.go
│       │   │   └── linkedin/
│       │   │       └── linkedin_node.go
│       │   │
│       │   ├── payment/
│       │   │   ├── stripe/
│       │   │   │   ├── stripe_node.go
│       │   │   │   └── stripe_trigger.go
│       │   │   ├── paypal/
│       │   │   │   └── paypal_node.go
│       │   │   └── square/
│       │   │       └── square_node.go
│       │   │
│       │   ├── analytics/
│       │   │   ├── google_analytics/
│       │   │   │   └── ga_node.go
│       │   │   ├── mixpanel/
│       │   │   │   └── mixpanel_node.go
│       │   │   └── segment/
│       │   │       └── segment_node.go
│       │   │
│       │   ├── ai/
│       │   │   ├── openai/
│       │   │   │   ├── chatgpt_node.go
│       │   │   │   ├── dalle_node.go
│       │   │   │   └── whisper_node.go
│       │   │   ├── anthropic/
│       │   │   │   └── claude_node.go
│       │   │   ├── huggingface/
│       │   │   │   └── huggingface_node.go
│       │   │   ├── stability/
│       │   │   │   └── stable_diffusion.go
│       │   │   └── langchain/
│       │   │       └── langchain_node.go
│       │   │
│       │   └── other/
│       │       ├── rss/
│       │       │   └── rss_feed.go
│       │       ├── mqtt/
│       │       │   └── mqtt_node.go
│       │       ├── kafka/
│       │       │   └── kafka_node.go
│       │       ├── rabbitmq/
│       │       │   └── rabbitmq_node.go
│       │       ├── redis/
│       │       │   └── redis_pubsub.go
│       │       └── custom_api/
│       │           └── custom_api_node.go
│       │
│       └── custom/                     # User-defined custom nodes
│           ├── loader.go               # Custom node loader
│           ├── validator.go            # Custom node validator
│           └── examples/               # Example custom nodes
│
├── pkg/                                # Public packages
│   ├── errors/
│   │   ├── errors.go                   # Error types and creation
│   │   ├── codes.go                    # Error codes
│   │   ├── handler.go                  # Error handling
│   │   ├── stack.go                    # Stack trace support
│   │   └── formatter.go                # Error formatting
│   │
│   ├── logger/
│   │   ├── logger.go                   # Logger interface
│   │   ├── zap.go                      # Zap implementation
│   │   ├── context.go                  # Context logging
│   │   ├── fields.go                   # Structured fields
│   │   └── hooks.go                    # Log hooks
│   │
│   ├── database/
│   │   ├── connection.go               # Database connection
│   │   ├── transaction.go              # Transaction handling
│   │   ├── migrations.go               # Migration runner
│   │   ├── seeder.go                   # Data seeding
│   │   ├── querybuilder.go             # Query builder
│   │   └── pagination.go               # Pagination helpers
│   │
│   ├── validator/
│   │   ├── validator.go                # Validation interface
│   │   ├── rules.go                    # Validation rules
│   │   ├── custom.go                   # Custom validators
│   │   └── errors.go                   # Validation errors
│   │
│   ├── crypto/
│   │   ├── hash.go                     # Hashing utilities
│   │   ├── encrypt.go                  # Encryption/Decryption
│   │   ├── jwt.go                      # JWT utilities
│   │   ├── random.go                   # Random generation
│   │   └── keys.go                     # Key management
│   │
│   ├── event/
│   │   ├── bus.go                      # Event bus interface
│   │   ├── dispatcher.go               # Event dispatcher
│   │   ├── handler.go                  # Event handler
│   │   ├── store.go                    # Event store
│   │   └── replay.go                   # Event replay
│   │
│   ├── cache/
│   │   ├── cache.go                    # Cache interface
│   │   ├── memory.go                   # In-memory cache
│   │   ├── redis.go                    # Redis cache
│   │   ├── multi_tier.go               # Multi-tier caching
│   │   └── invalidation.go             # Cache invalidation
│   │
│   ├── pubsub/
│   │   ├── publisher.go                # Publisher interface
│   │   ├── subscriber.go               # Subscriber interface
│   │   ├── message.go                  # Message definition
│   │   └── broker.go                   # Message broker
│   │
│   ├── http/
│   │   ├── client.go                   # HTTP client wrapper
│   │   ├── response.go                 # Response helpers
│   │   ├── request.go                  # Request helpers
│   │   ├── retry.go                    # Retry logic
│   │   └── circuit_breaker.go         # Circuit breaker
│   │
│   ├── auth/
│   │   ├── authenticator.go            # Authentication interface
│   │   ├── authorizer.go               # Authorization interface
│   │   ├── token.go                    # Token management
│   │   └── session.go                  # Session management
│   │
│   ├── workflow/
│   │   ├── dsl.go                      # Workflow DSL parser
│   │   ├── validator.go                # Workflow validator
│   │   └── optimizer.go                # Workflow optimizer
│   │
│   └── utils/
│       ├── uuid.go                     # UUID utilities
│       ├── time.go                     # Time utilities
│       ├── json.go                     # JSON utilities
│       ├── strings.go                  # String utilities
│       ├── slice.go                    # Slice utilities
│       ├── map.go                      # Map utilities
│       ├── retry.go                    # Retry utilities
│       └── async.go                    # Async utilities
│
├── configs/                            # Configuration files
│   ├── base/
│   │   ├── app.yaml                   # Base application config
│   │   ├── database.yaml              # Database configuration
│   │   ├── cache.yaml                 # Cache configuration
│   │   ├── security.yaml              # Security settings
│   │   └── features.yaml              # Feature flags
│   ├── environments/
│   │   ├── development.yaml           # Development overrides
│   │   ├── staging.yaml               # Staging overrides
│   │   └── production.yaml            # Production overrides
│   └── embed.go                       # Embed configs in binary
│
├── deployments/                        # Deployment configurations
│   ├── docker/
│   │   ├── Dockerfile.api
│   │   ├── Dockerfile.worker
│   │   ├── Dockerfile.scheduler
│   │   ├── docker-compose.yml
│   │   ├── docker-compose.dev.yml
│   │   └── docker-compose.prod.yml
│   │
│   ├── kubernetes/
│   │   ├── base/
│   │   │   ├── namespace.yaml
│   │   │   ├── configmap.yaml
│   │   │   ├── secret.yaml
│   │   │   ├── deployment.yaml
│   │   │   ├── service.yaml
│   │   │   ├── ingress.yaml
│   │   │   └── hpa.yaml              # Horizontal Pod Autoscaler
│   │   ├── overlays/
│   │   │   ├── development/
│   │   │   ├── staging/
│   │   │   └── production/
│   │   └── kustomization.yaml
│   │
│   ├── helm/
│   │   └── go-n8n/
│   │       ├── Chart.yaml
│   │       ├── values.yaml
│   │       ├── values.dev.yaml
│   │       ├── values.prod.yaml
│   │       └── templates/
│   │
│   ├── terraform/
│   │   ├── modules/
│   │   │   ├── vpc/
│   │   │   ├── eks/
│   │   │   ├── rds/
│   │   │   └── redis/
│   │   ├── environments/
│   │   └── main.tf
│   │
│   └── ansible/
│       ├── playbooks/
│       └── inventory/
│
├── scripts/                           # Utility scripts
│   ├── build/
│   │   ├── build.sh                  # Build script
│   │   └── release.sh                # Release script
│   ├── development/
│   │   ├── setup.sh                  # Development setup
│   │   └── seed.sh                   # Database seeding
│   ├── testing/
│   │   ├── test.sh                   # Test runner
│   │   └── benchmark.sh              # Benchmark runner
│   └── deployment/
│       ├── deploy.sh                 # Deployment script
│       └── rollback.sh               # Rollback script
│
├── test/                             # Test files
│   ├── unit/                        # Unit tests
│   ├── integration/                 # Integration tests
│   ├── e2e/                         # End-to-end tests
│   ├── load/                        # Load tests
│   │   ├── k6/
│   │   └── locust/
│   ├── fixtures/                    # Test fixtures
│   ├── mocks/                       # Mock implementations
│   └── helpers/                     # Test helpers
│
├── docs/                            # Documentation
│   ├── api/                        # API documentation
│   │   ├── openapi.yaml
│   │   └── postman/
│   ├── architecture/               # Architecture documentation
│   │   ├── ADR/                    # Architecture Decision Records
│   │   ├── diagrams/               # Architecture diagrams
│   │   └── patterns.md             # Design patterns used
│   ├── development/                # Development guides
│   │   ├── setup.md
│   │   ├── contributing.md
│   │   └── coding-standards.md
│   ├── deployment/                 # Deployment guides
│   └── user/                       # User documentation
│
├── tools/                          # Development tools
│   ├── codegen/                   # Code generation tools
│   │   ├── node_generator.go
│   │   └── templates/
│   ├── migration/                 # Migration tools
│   └── lint/                      # Linting configuration
│       └── .golangci.yml
│
├── .github/                       # GitHub configuration
│   ├── workflows/                # GitHub Actions
│   │   ├── ci.yml
│   │   ├── cd.yml
│   │   ├── security.yml
│   │   └── codeql.yml
│   ├── ISSUE_TEMPLATE/
│   ├── PULL_REQUEST_TEMPLATE.md
│   └── CODEOWNERS
│
├── Makefile
├── go.mod
├── go.sum
├── .env.example
├── .gitignore
├── .dockerignore
├── .editorconfig
├── README.md
├── CHANGELOG.md
├── LICENSE
└── SECURITY.md
```

## 🔐 Security Architecture

### Authentication & Authorization Flow

```
1. Multi-Factor Authentication
   ├── Password + TOTP
   ├── OAuth2 + Email verification
   └── API Key + IP Whitelist

2. Token Management
   ├── JWT with short expiry (15 min)
   ├── Refresh tokens (7 days)
   ├── Token rotation on refresh
   └── Blacklist for revoked tokens

3. Authorization Levels
   ├── Resource-based (owns workflow)
   ├── Role-based (admin, user, viewer)
   ├── Team-based (organization access)
   └── Feature-based (feature flags)

4. Security Headers
   ├── CSP (Content Security Policy)
   ├── HSTS (HTTP Strict Transport Security)
   ├── X-Frame-Options
   └── X-Content-Type-Options
```

### Encryption Strategy

```go
// internal/infrastructure/security/encryption/strategy.go

type EncryptionStrategy struct {
    // Data at rest
    DatabaseEncryption: AES256_GCM
    FileEncryption:     AES256_CTR
    
    // Data in transit
    TLSVersion:         TLS1.3
    CipherSuites:       []Modern_Ciphers
    
    // Credentials
    CredentialVault:    HashiCorp_Vault
    KeyRotation:        Monthly
    
    // Secrets
    EnvironmentVars:    Encrypted_with_KMS
    ConfigFiles:        Sealed_Secrets
}
```

## 🚀 Performance Optimization

### Caching Strategy

```yaml
# Multi-tier caching architecture
L1_Cache:
  type: in-memory
  size: 100MB
  ttl: 60s
  strategy: LRU
  
L2_Cache:
  type: Redis
  size: 1GB
  ttl: 300s
  strategy: LFU
  
L3_Cache:
  type: CDN
  providers: [CloudFlare, Fastly]
  ttl: 3600s
  
Cache_Patterns:
  - Cache-aside (lazy loading)
  - Write-through (immediate write)
  - Write-behind (async write)
  - Refresh-ahead (proactive refresh)
```

### Database Optimization

```sql
-- Partitioning strategy
CREATE TABLE executions (
    id UUID,
    workflow_id UUID,
    created_at TIMESTAMP
) PARTITION BY RANGE (created_at);

-- Create monthly partitions
CREATE TABLE executions_2024_01 PARTITION OF executions
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

-- Indexes for performance
CREATE INDEX CONCURRENTLY idx_workflows_user_active 
    ON workflows(user_id, is_active) 
    WHERE deleted_at IS NULL;

CREATE INDEX idx_executions_workflow_status 
    ON executions(workflow_id, status) 
    INCLUDE (created_at, finished_at);

-- Materialized views for analytics
CREATE MATERIALIZED VIEW workflow_statistics AS
SELECT 
    workflow_id,
    COUNT(*) as total_executions,
    AVG(execution_time_ms) as avg_duration,
    COUNT(*) FILTER (WHERE status = 'success') as successful
FROM executions
GROUP BY workflow_id;

-- Automatic refresh
CREATE INDEX ON workflow_statistics(workflow_id);
REFRESH MATERIALIZED VIEW CONCURRENTLY workflow_statistics;
```

### Connection Pooling

```go
// pkg/database/pool_config.go

type PoolConfig struct {
    // API Server Pool
    APIPool: DBPool{
        MaxConnections:     50,
        MaxIdleConnections: 10,
        ConnectionLifetime: 5 * time.Minute,
        HealthCheckPeriod:  30 * time.Second,
    },
    
    // Worker Pool (larger for batch operations)
    WorkerPool: DBPool{
        MaxConnections:     100,
        MaxIdleConnections: 20,
        ConnectionLifetime: 10 * time.Minute,
    },
    
    // Read Replica Pool
    ReadPool: DBPool{
        MaxConnections:     200,
        LoadBalancing:      RoundRobin,
        Replicas:          []string{"read1.db", "read2.db"},
    },
}
```

## 🔄 Event-Driven Architecture

### Event Sourcing Pattern

```go
// internal/domain/common/event_sourcing.go

type EventStore interface {
    Append(aggregateID string, events []Event) error
    Load(aggregateID string) ([]Event, error)
    LoadSnapshot(aggregateID string) (*Snapshot, error)
    SaveSnapshot(aggregateID string, snapshot Snapshot) error
}

type EventStream struct {
    Events    []Event
    Version   int
    Timestamp time.Time
}

// CQRS Implementation
type CommandBus interface {
    Send(ctx context.Context, command Command) error
}

type QueryBus interface {
    Send(ctx context.Context, query Query) (interface{}, error)
}

type EventBus interface {
    Publish(ctx context.Context, events ...Event) error
    Subscribe(eventType string, handler EventHandler) error
}
```

### SAGA Pattern for Distributed Transactions

```go
// internal/application/shared/saga_orchestrator.go

type SagaStep struct {
    Name            string
    Execute         func(context.Context, interface{}) error
    Compensate      func(context.Context, interface{}) error
    RetryPolicy     RetryPolicy
}

type Saga struct {
    ID              string
    Steps           []SagaStep
    State           SagaState
    CompletedSteps  []string
}

type SagaOrchestrator struct {
    repository SagaRepository
    eventBus   EventBus
    
    Execute(ctx context.Context, saga *Saga) error {
        for _, step := range saga.Steps {
            if err := step.Execute(ctx, saga.State); err != nil {
                // Compensate in reverse order
                return s.compensate(ctx, saga)
            }
            saga.CompletedSteps = append(saga.CompletedSteps, step.Name)
            s.repository.SaveProgress(saga)
        }
        return nil
    }
}
```

## 📊 Monitoring & Observability

### Metrics Collection

```go
// internal/infrastructure/observability/metrics/collectors.go

var (
    WorkflowExecutions = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "workflow_executions_total",
            Help: "Total number of workflow executions",
        },
        []string{"workflow_id", "status"},
    )
    
    NodeExecutionDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "node_execution_duration_seconds",
            Help:    "Node execution duration in seconds",
            Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
        },
        []string{"node_type", "status"},
    )
    
    QueueDepth = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "queue_depth",
            Help: "Current queue depth",
        },
        []string{"queue_name", "priority"},
    )
)
```

### Distributed Tracing

```go
// internal/infrastructure/observability/tracing/middleware.go

func TracingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tracer := otel.Tracer("api")
        
        ctx, span := tracer.Start(
            c.Request.Context(),
            fmt.Sprintf("%s %s", c.Request.Method, c.FullPath()),
            trace.WithAttributes(
                attribute.String("http.method", c.Request.Method),
                attribute.String("http.url", c.Request.URL.String()),
                attribute.String("http.user_agent", c.Request.UserAgent()),
            ),
        )
        defer span.End()
        
        c.Request = c.Request.WithContext(ctx)
        c.Next()
        
        span.SetAttributes(
            attribute.Int("http.status_code", c.Writer.Status()),
            attribute.Int64("http.response_size", int64(c.Writer.Size())),
        )
    }
}
```

### Health Checks

```go
// internal/infrastructure/observability/health/checker.go

type HealthChecker struct {
    checks map[string]Check
}

func (h *HealthChecker) RegisterChecks() {
    h.Register("database", DatabaseCheck{})
    h.Register("redis", RedisCheck{})
    h.Register("kafka", KafkaCheck{})
    h.Register("disk_space", DiskSpaceCheck{MinFreeBytes: 1GB})
    h.Register("memory", MemoryCheck{MaxUsagePercent: 90})
}

// Kubernetes probes
func (h *HealthChecker) Liveness() HealthStatus {
    // Basic checks - is the application alive?
    return h.checkCritical()
}

func (h *HealthChecker) Readiness() HealthStatus {
    // Full checks - is the application ready to serve traffic?
    return h.checkAll()
}

func (h *HealthChecker) Startup() HealthStatus {
    // Initial checks - has the application started successfully?
    return h.checkStartup()
}
```

## 🧪 Testing Strategy

### Test Pyramid

```
         /\
        /  \  E2E Tests (5%)
       /    \  - Full user journeys
      /------\ Integration Tests (15%)
     /        \ - API tests, DB tests
    /----------\ Component Tests (30%)
   /            \ - Service tests
  /--------------\ Unit Tests (50%)
 /                \ - Domain logic, utilities
```

### Test Organization

```go
// test/helpers/database.go
func SetupTestDB(t *testing.T) *sqlx.DB {
    container := testcontainers.PostgresContainer()
    db := container.GetDB()
    t.Cleanup(func() { container.Stop() })
    return db
}

// test/helpers/fixtures.go
func LoadFixtures(t *testing.T, path string) {
    // Load test data from fixtures
}

// test/integration/workflow_test.go
func TestWorkflowExecution(t *testing.T) {
    db := SetupTestDB(t)
    LoadFixtures(t, "workflows.json")
    
    // Test workflow execution
    result := ExecuteWorkflow(ctx, workflowID)
    assert.Equal(t, "success", result.Status)
}
```

## 🚦 CI/CD Pipeline

### GitHub Actions Workflow

```yaml
# .github/workflows/ci.yml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: golangci/golangci-lint-action@v3
        with:
          version: latest

  test:
    strategy:
      matrix:
        go-version: [1.21, 1.22]
        os: [ubuntu-latest, macos-latest]
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go-version }}
      - run: make test-coverage
      - uses: codecov/codecov-action@v3

  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: securego/gosec@master
      - uses: aquasecurity/trivy-action@master

  build:
    needs: [lint, test, security]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: docker/setup-buildx-action@v2
      - uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: |
            ghcr.io/${{ github.repository }}:${{ github.sha }}
            ghcr.io/${{ github.repository }}:latest
          cache-from: type=gha
          cache-to: type=gha,mode=max

  deploy:
    needs: build
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - uses: azure/k8s-deploy@v4
        with:
          manifests: |
            deployments/kubernetes/
          images: |
            ghcr.io/${{ github.repository }}:${{ github.sha }}
```

## 🎯 Development Workflow

### Git Flow

```
main (production)
  ├── develop (staging)
  │   ├── feature/workflow-improvements
  │   ├── feature/new-node-type
  │   └── feature/performance-optimization
  ├── release/v1.2.0
  └── hotfix/critical-bug-fix
```

### Makefile Commands

```makefile
# Complete Makefile
.PHONY: all build test clean

# Variables
BINARY_NAME=go-n8n
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Development
dev:
	@air -c .air.toml

setup:
	@go mod download
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/golang/mock/mockgen@latest
	@docker-compose -f deployments/docker/docker-compose.dev.yml up -d

# Building
build:
	@echo "Building..."
	@CGO_ENABLED=0 go build ${LDFLAGS} -o bin/api cmd/api/main.go
	@CGO_ENABLED=0 go build ${LDFLAGS} -o bin/worker cmd/worker/main.go
	@CGO_ENABLED=0 go build ${LDFLAGS} -o bin/scheduler cmd/scheduler/main.go

build-docker:
	@docker build -f deployments/docker/Dockerfile.api -t ${BINARY_NAME}-api:${VERSION} .
	@docker build -f deployments/docker/Dockerfile.worker -t ${BINARY_NAME}-worker:${VERSION} .

# Testing
test:
	@go test -v -race ./...

test-coverage:
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

test-integration:
	@go test -v -tags=integration ./test/integration/...

test-e2e:
	@go test -v -tags=e2e ./test/e2e/...

benchmark:
	@go test -bench=. -benchmem ./...

# Database
migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

migrate-create:
	@go run cmd/migrate/main.go create $(name)

seed:
	@go run cmd/migrate/seed.go

# Code Quality
lint:
	@golangci-lint run --deadline=5m

fmt:
	@go fmt ./...
	@gofumpt -l -w .

vet:
	@go vet ./...

security:
	@gosec -fmt sarif -out results.sarif ./...

generate:
	@go generate ./...

mock:
	@mockgen -source=internal/domain/workflow/repository.go -destination=test/mocks/workflow_repository_mock.go

# Deployment
deploy-dev:
	@kubectl apply -k deployments/kubernetes/overlays/development

deploy-staging:
	@kubectl apply -k deployments/kubernetes/overlays/staging

deploy-prod:
	@kubectl apply -k deployments/kubernetes/overlays/production

rollback:
	@kubectl rollout undo deployment/api -n go-n8n
	@kubectl rollout undo deployment/worker -n go-n8n

# Monitoring
logs-api:
	@kubectl logs -f deployment/api -n go-n8n

logs-worker:
	@kubectl logs -f deployment/worker -n go-n8n

metrics:
	@open http://localhost:9090/metrics

# Clean
clean:
	@rm -rf bin/ coverage.* vendor/

# Help
help:
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
```

## 🔄 Data Consistency Patterns

### Transactional Outbox Pattern

```sql
-- Outbox table for reliable event publishing
CREATE TABLE outbox_events (
    id UUID PRIMARY KEY,
    aggregate_id UUID NOT NULL,
    aggregate_type VARCHAR(100) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    event_data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP,
    INDEX idx_unpublished (published_at) WHERE published_at IS NULL
);

-- Trigger to insert events
CREATE TRIGGER workflow_events_trigger
AFTER INSERT OR UPDATE OR DELETE ON workflows
FOR EACH ROW EXECUTE FUNCTION publish_workflow_event();
```

### Two-Phase Commit for Distributed Transactions

```go
// internal/application/shared/transaction_coordinator.go

type TwoPhaseCommit struct {
    participants []Participant
}

func (tpc *TwoPhaseCommit) Execute(ctx context.Context, tx Transaction) error {
    // Phase 1: Prepare
    for _, p := range tpc.participants {
        if err := p.Prepare(ctx, tx); err != nil {
            tpc.abort(ctx, tx)
            return err
        }
    }
    
    // Phase 2: Commit
    for _, p := range tpc.participants {
        if err := p.Commit(ctx, tx); err != nil {
            // Log critical error - manual intervention needed
            return err
        }
    }
    
    return nil
}
```

## 🚀 Performance Patterns

### Circuit Breaker Pattern

```go
// pkg/http/circuit_breaker.go

type CircuitBreaker struct {
    maxFailures  int
    resetTimeout time.Duration
    state        State
    failures     int
    lastFailTime time.Time
    mutex        sync.RWMutex
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()
    
    switch cb.state {
    case Closed:
        err := fn()
        if err != nil {
            cb.failures++
            cb.lastFailTime = time.Now()
            if cb.failures >= cb.maxFailures {
                cb.state = Open
            }
            return err
        }
        cb.failures = 0
        return nil
        
    case Open:
        if time.Since(cb.lastFailTime) > cb.resetTimeout {
            cb.state = HalfOpen
            return cb.Call(fn)
        }
        return ErrCircuitOpen
        
    case HalfOpen:
        err := fn()
        if err != nil {
            cb.state = Open
            cb.lastFailTime = time.Now()
            return err
        }
        cb.state = Closed
        cb.failures = 0
        return nil
    }
}
```

### Rate Limiting with Token Bucket

```go
// internal/infrastructure/persistence/redis/rate_limiter.go

type TokenBucket struct {
    redis    *redis.Client
    capacity int
    refillRate time.Duration
}

func (tb *TokenBucket) Allow(ctx context.Context, key string) bool {
    script := `
        local key = KEYS[1]
        local capacity = tonumber(ARGV[1])
        local now = tonumber(ARGV[2])
        local refill_rate = tonumber(ARGV[3])
        
        local bucket = redis.call('HGETALL', key)
        local tokens = capacity
        local last_refill = now
        
        if #bucket > 0 then
            tokens = tonumber(bucket[2])
            last_refill = tonumber(bucket[4])
            
            local elapsed = now - last_refill
            local refill = math.floor(elapsed / refill_rate)
            tokens = math.min(capacity, tokens + refill)
        end
        
        if tokens > 0 then
            redis.call('HMSET', key, 'tokens', tokens - 1, 'last_refill', now)
            redis.call('EXPIRE', key, 3600)
            return 1
        end
        
        return 0
    `
    
    result, err := tb.redis.Eval(ctx, script, []string{key}, 
        tb.capacity, time.Now().Unix(), tb.refillRate.Seconds()).Result()
    
    return result.(int64) == 1
}
```

## 🔒 Advanced Security Patterns

### Zero Trust Architecture

```go
// internal/infrastructure/security/zero_trust.go

type ZeroTrustValidator struct {
    // Verify everything, trust nothing
    
    ValidateRequest(ctx context.Context, req Request) error {
        // 1. Authenticate user
        if err := s.authenticateUser(req); err != nil {
            return err
        }
        
        // 2. Verify device
        if err := s.verifyDevice(req); err != nil {
            return err
        }
        
        // 3. Check location
        if err := s.validateLocation(req); err != nil {
            return err
        }
        
        // 4. Analyze behavior
        if err := s.analyzeBehavior(req); err != nil {
            return err
        }
        
        // 5. Enforce least privilege
        if err := s.enforcePermissions(req); err != nil {
            return err
        }
        
        return nil
    }
}
```

### API Versioning Strategy

```go
// internal/interfaces/http/rest/versioning.go

type APIVersion string

const (
    V1 APIVersion = "v1"
    V2 APIVersion = "v2"
)

func VersionMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        version := c.Param("version")
        
        if version == "" {
            // Check Accept header
            accept := c.GetHeader("Accept")
            if strings.Contains(accept, "application/vnd.api+json;version=2") {
                version = "v2"
            } else {
                version = "v1" // Default
            }
        }
        
        c.Set("api_version", version)
        c.Next()
    }
}

// Deprecation headers
func DeprecationMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.GetString("api_version") == "v1" {
            c.Header("Sunset", "2025-01-01")
            c.Header("Deprecation", "true")
            c.Header("Link", "</api/v2>; rel=\"successor-version\"")
        }
        c.Next()
    }
}
```

## 📈 Scaling Strategies

### Database Read Replicas

```go
// pkg/database/read_write_split.go

type DBManager struct {
    master    *sqlx.DB
    replicas  []*sqlx.DB
    strategy  LoadBalanceStrategy
}

func (m *DBManager) Read() *sqlx.DB {
    return m.strategy.SelectReplica(m.replicas)
}

func (m *DBManager) Write() *sqlx.DB {
    return m.master
}

// Usage in repository
func (r *WorkflowRepository) FindByID(ctx context.Context, id string) (*Workflow, error) {
    db := r.dbManager.Read() // Use read replica
    // Query execution
}

func (r *WorkflowRepository) Save(ctx context.Context, w *Workflow) error {
    db := r.dbManager.Write() // Use master
    // Insert/Update execution
}
```

### Sharding Strategy

```go
// internal/infrastructure/persistence/sharding.go

type ShardManager struct {
    shards map[int]*sqlx.DB
    shardCount int
}

func (sm *ShardManager) GetShard(key string) *sqlx.DB {
    hash := crc32.ChecksumIEEE([]byte(key))
    shardID := int(hash) % sm.shardCount
    return sm.shards[shardID]
}

// Usage
func (r *ExecutionRepository) Save(ctx context.Context, e *Execution) error {
    db := r.shardManager.GetShard(e.WorkflowID) // Shard by workflow ID
    // Save to specific shard
}
```

## 🔄 Migration Strategy

### Blue-Green Deployment

```yaml
# deployments/kubernetes/blue-green.yaml
apiVersion: v1
kind: Service
metadata:
  name: app-service
spec:
  selector:
    app: go-n8n
    version: green  # Switch between blue/green
  ports:
    - port: 80
      targetPort: 8080

---
# Blue deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-blue
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-n8n
      version: blue
  template:
    metadata:
      labels:
        app: go-n8n
        version: blue
    spec:
      containers:
      - name: app
        image: go-n8n:v1.0.0

---
# Green deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-green
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-n8n
      version: green
  template:
    metadata:
      labels:
        app: go-n8n
        version: green
    spec:
      containers:
      - name: app
        image: go-n8n:v2.0.0
```

## 🎯 Production Readiness Checklist

```yaml
Security:
  ✓ TLS 1.3 everywhere
  ✓ Secrets in Vault
  ✓ RBAC configured
  ✓ Network policies
  ✓ Security scanning
  ✓ Vulnerability management

Performance:
  ✓ Database indexes optimized
  ✓ Caching layers configured
  ✓ Connection pools tuned
  ✓ Rate limiting enabled
  ✓ Circuit breakers configured
  ✓ Load testing completed

Reliability:
  ✓ Health checks implemented
  ✓ Graceful shutdown
  ✓ Retry logic with backoff
  ✓ Timeout configurations
  ✓ Error handling comprehensive
  ✓ Disaster recovery plan

Observability:
  ✓ Metrics collection
  ✓ Distributed tracing
  ✓ Centralized logging
  ✓ Alerting rules defined
  ✓ Dashboards created
  ✓ SLOs defined

Operations:
  ✓ CI/CD pipeline
  ✓ Automated testing
  ✓ Rollback procedures
  ✓ Backup strategy
  ✓ Documentation complete
  ✓ Runbooks prepared
```

This enhanced architecture provides everything you need for a production-grade n8n clone that can scale from startup to enterprise!
