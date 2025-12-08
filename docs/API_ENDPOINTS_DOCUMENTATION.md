# Complete API Endpoints Documentation for n8n Clone

## Base URL
```
Production: https://api.your-n8n-clone.com/api/v1
Development: http://localhost:8080/api/v1
```

## Authentication Headers
```
Authorization: Bearer <jwt_token>
X-API-Key: <api_key> (for programmatic access)
Content-Type: application/json
```

## API Endpoints

### 1. Authentication & Authorization

#### 1.1 User Registration
```http
POST /auth/register
```
**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "name": "John Doe",
  "inviteCode": "optional-invite-code"
}
```
**Response:** `201 Created`
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe",
    "role": "user"
  },
  "tokens": {
    "accessToken": "jwt_token",
    "refreshToken": "refresh_token",
    "expiresIn": 900
  }
}
```

#### 1.2 User Login
```http
POST /auth/login
```
**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password"
}
```

#### 1.3 Refresh Token
```http
POST /auth/refresh
```
**Request Body:**
```json
{
  "refreshToken": "refresh_token"
}
```

#### 1.4 Logout
```http
POST /auth/logout
```

#### 1.5 Get Current User
```http
GET /auth/me
```

#### 1.6 Forgot Password
```http
POST /auth/forgot-password
```
**Request Body:**
```json
{
  "email": "user@example.com"
}
```

#### 1.7 Reset Password
```http
POST /auth/reset-password
```
**Request Body:**
```json
{
  "token": "reset_token",
  "newPassword": "NewSecurePassword123!"
}
```

#### 1.8 Verify Email
```http
POST /auth/verify-email
```
**Request Body:**
```json
{
  "token": "verification_token"
}
```

#### 1.9 Change Password
```http
PUT /auth/change-password
```
**Request Body:**
```json
{
  "currentPassword": "current_password",
  "newPassword": "new_password"
}
```

#### 1.10 Enable 2FA
```http
POST /auth/2fa/enable
```

#### 1.11 Disable 2FA
```http
POST /auth/2fa/disable
```

#### 1.12 Verify 2FA
```http
POST /auth/2fa/verify
```
**Request Body:**
```json
{
  "code": "123456"
}
```

### 2. User Management

#### 2.1 List Users (Admin)
```http
GET /users
```
**Query Parameters:**
- `page` (int): Page number
- `limit` (int): Items per page
- `search` (string): Search term
- `role` (string): Filter by role
- `status` (string): active/inactive

#### 2.2 Get User by ID
```http
GET /users/:id
```

#### 2.3 Update User
```http
PUT /users/:id
```
**Request Body:**
```json
{
  "name": "Updated Name",
  "email": "newemail@example.com",
  "role": "admin"
}
```

#### 2.4 Delete User
```http
DELETE /users/:id
```

#### 2.5 Update User Settings
```http
PUT /users/:id/settings
```
**Request Body:**
```json
{
  "timezone": "America/New_York",
  "dateFormat": "MM/DD/YYYY",
  "theme": "dark",
  "notifications": {
    "email": true,
    "push": false
  }
}
```

#### 2.6 Get User Permissions
```http
GET /users/:id/permissions
```

#### 2.7 Update User Permissions
```http
PUT /users/:id/permissions
```

### 3. Workflows

#### 3.1 List Workflows
```http
GET /workflows
```
**Query Parameters:**
- `page` (int): Page number
- `limit` (int): Items per page (default: 20)
- `search` (string): Search in name/description
- `tags[]` (array): Filter by tags
- `active` (boolean): Filter by active status
- `sort` (string): name|created_at|updated_at
- `order` (string): asc|desc

#### 3.2 Create Workflow
```http
POST /workflows
```
**Request Body:**
```json
{
  "name": "My Workflow",
  "description": "Workflow description",
  "nodes": [
    {
      "id": "node1",
      "type": "webhook",
      "name": "Webhook Trigger",
      "position": [100, 100],
      "parameters": {
        "path": "my-webhook",
        "method": "POST"
      }
    }
  ],
  "connections": [
    {
      "sourceNodeId": "node1",
      "sourceOutput": "main",
      "targetNodeId": "node2",
      "targetInput": "main"
    }
  ],
  "settings": {
    "errorWorkflow": "workflow_id",
    "timezone": "UTC",
    "timeout": 3600
  },
  "tags": ["automation", "sales"]
}
```

#### 3.3 Get Workflow
```http
GET /workflows/:id
```

#### 3.4 Update Workflow
```http
PUT /workflows/:id
```

#### 3.5 Delete Workflow
```http
DELETE /workflows/:id
```

#### 3.6 Duplicate Workflow
```http
POST /workflows/:id/duplicate
```

#### 3.7 Activate Workflow
```http
POST /workflows/:id/activate
```

#### 3.8 Deactivate Workflow
```http
POST /workflows/:id/deactivate
```

#### 3.9 Execute Workflow
```http
POST /workflows/:id/execute
```
**Request Body:**
```json
{
  "inputData": {
    "key": "value"
  },
  "runData": {},
  "mode": "manual",
  "startNodes": ["node1"]
}
```

#### 3.10 Test Workflow
```http
POST /workflows/:id/test
```

#### 3.11 Get Workflow Nodes
```http
GET /workflows/:id/nodes
```

#### 3.12 Update Workflow Nodes
```http
PUT /workflows/:id/nodes
```

#### 3.13 Get Workflow Executions
```http
GET /workflows/:id/executions
```

#### 3.14 Export Workflow
```http
GET /workflows/:id/export
```

#### 3.15 Import Workflow
```http
POST /workflows/import
```

#### 3.16 Get Workflow Statistics
```http
GET /workflows/:id/statistics
```

#### 3.17 Get Workflow Versions
```http
GET /workflows/:id/versions
```

#### 3.18 Restore Workflow Version
```http
POST /workflows/:id/versions/:versionId/restore
```

#### 3.19 Share Workflow
```http
POST /workflows/:id/share
```
**Request Body:**
```json
{
  "userIds": ["user1", "user2"],
  "permission": "read|write|execute"
}
```

#### 3.20 Get Workflow Metrics
```http
GET /workflows/:id/metrics
```
**Query Parameters:**
- `startDate` (ISO 8601)
- `endDate` (ISO 8601)
- `granularity` (hour|day|week|month)

### 4. Nodes

#### 4.1 List Available Node Types
```http
GET /nodes/types
```

#### 4.2 Get Node Type Details
```http
GET /nodes/types/:type
```

#### 4.3 Get Node Schema
```http
GET /nodes/types/:type/schema
```

#### 4.4 Create Node
```http
POST /workflows/:workflowId/nodes
```
**Request Body:**
```json
{
  "type": "http_request",
  "name": "HTTP Request",
  "position": [200, 200],
  "parameters": {},
  "credentials": "credential_id"
}
```

#### 4.5 Update Node
```http
PUT /nodes/:id
```

#### 4.6 Delete Node
```http
DELETE /nodes/:id
```

#### 4.7 Test Node
```http
POST /nodes/:id/test
```

#### 4.8 Get Node Execution Data
```http
GET /nodes/:id/executions/:executionId/data
```

#### 4.9 Pin Node Data
```http
POST /nodes/:id/pin
```

#### 4.10 Unpin Node Data
```http
DELETE /nodes/:id/pin
```

### 5. Connections

#### 5.1 Create Connection
```http
POST /workflows/:workflowId/connections
```
**Request Body:**
```json
{
  "sourceNodeId": "node1",
  "sourceOutput": "main",
  "targetNodeId": "node2",
  "targetInput": "main"
}
```

#### 5.2 Delete Connection
```http
DELETE /connections/:id
```

#### 5.3 Update Connection
```http
PUT /connections/:id
```

### 6. Executions

#### 6.1 List All Executions
```http
GET /executions
```
**Query Parameters:**
- `workflowId` (string): Filter by workflow
- `status` (string): waiting|running|success|error|cancelled
- `mode` (string): manual|trigger|webhook|schedule
- `startDate` (ISO 8601)
- `endDate` (ISO 8601)
- `page` (int)
- `limit` (int)

#### 6.2 Get Execution
```http
GET /executions/:id
```

#### 6.3 Get Execution Data
```http
GET /executions/:id/data
```

#### 6.4 Stop Execution
```http
POST /executions/:id/stop
```

#### 6.5 Retry Execution
```http
POST /executions/:id/retry
```

#### 6.6 Delete Execution
```http
DELETE /executions/:id
```

#### 6.7 Delete Multiple Executions
```http
POST /executions/delete
```
**Request Body:**
```json
{
  "ids": ["exec1", "exec2", "exec3"]
}
```

#### 6.8 Get Execution Logs
```http
GET /executions/:id/logs
```

#### 6.9 Get Execution Timeline
```http
GET /executions/:id/timeline
```

### 7. Credentials

#### 7.1 List Credentials
```http
GET /credentials
```

#### 7.2 Create Credential
```http
POST /credentials
```
**Request Body:**
```json
{
  "name": "My API Key",
  "type": "api_key",
  "nodeTypes": ["http_request", "webhook"],
  "data": {
    "apiKey": "secret_key",
    "domain": "api.example.com"
  }
}
```

#### 7.3 Get Credential (Masked)
```http
GET /credentials/:id
```

#### 7.4 Update Credential
```http
PUT /credentials/:id
```

#### 7.5 Delete Credential
```http
DELETE /credentials/:id
```

#### 7.6 Test Credential
```http
POST /credentials/:id/test
```

#### 7.7 Get OAuth2 Redirect URL
```http
GET /credentials/oauth2/:credentialType/auth
```

#### 7.8 OAuth2 Callback
```http
GET /credentials/oauth2/callback
```

#### 7.9 Share Credential
```http
POST /credentials/:id/share
```

### 8. Webhooks

#### 8.1 List Webhooks
```http
GET /webhooks
```

#### 8.2 Create Webhook
```http
POST /webhooks
```
**Request Body:**
```json
{
  "workflowId": "workflow_id",
  "nodeId": "node_id",
  "path": "my-webhook",
  "method": "POST",
  "requiresAuth": false
}
```

#### 8.3 Get Webhook
```http
GET /webhooks/:id
```

#### 8.4 Update Webhook
```http
PUT /webhooks/:id
```

#### 8.5 Delete Webhook
```http
DELETE /webhooks/:id
```

#### 8.6 Test Webhook
```http
POST /webhooks/:id/test
```

#### 8.7 Get Webhook URL
```http
GET /webhooks/:id/url
```

#### 8.8 Webhook Endpoint (Dynamic)
```http
ANY /webhook/:path
ANY /webhook-test/:path
ANY /webhook-waiting/:path
```

### 9. Templates

#### 9.1 List Workflow Templates
```http
GET /templates
```
**Query Parameters:**
- `category` (string): Filter by category
- `search` (string): Search term
- `tags[]` (array): Filter by tags
- `sort` (string): popular|recent|name

#### 9.2 Get Template
```http
GET /templates/:id
```

#### 9.3 Create Template from Workflow
```http
POST /templates
```

#### 9.4 Use Template
```http
POST /templates/:id/use
```

#### 9.5 Update Template
```http
PUT /templates/:id
```

#### 9.6 Delete Template
```http
DELETE /templates/:id
```

#### 9.7 Get Template Categories
```http
GET /templates/categories
```

### 10. Tags

#### 10.1 List Tags
```http
GET /tags
```

#### 10.2 Create Tag
```http
POST /tags
```
**Request Body:**
```json
{
  "name": "automation",
  "color": "#FF5733"
}
```

#### 10.3 Update Tag
```http
PUT /tags/:id
```

#### 10.4 Delete Tag
```http
DELETE /tags/:id
```

#### 10.5 Get Workflows by Tag
```http
GET /tags/:id/workflows
```

### 11. Variables & Environment

#### 11.1 List Variables
```http
GET /variables
```

#### 11.2 Create Variable
```http
POST /variables
```
**Request Body:**
```json
{
  "key": "API_ENDPOINT",
  "value": "https://api.example.com",
  "type": "string",
  "isSecret": false
}
```

#### 11.3 Get Variable
```http
GET /variables/:key
```

#### 11.4 Update Variable
```http
PUT /variables/:key
```

#### 11.5 Delete Variable
```http
DELETE /variables/:key
```

### 12. API Keys

#### 12.1 List API Keys
```http
GET /api-keys
```

#### 12.2 Create API Key
```http
POST /api-keys
```
**Request Body:**
```json
{
  "name": "CI/CD Pipeline",
  "scopes": ["workflows:read", "workflows:execute"],
  "expiresAt": "2024-12-31T23:59:59Z"
}
```

#### 12.3 Get API Key
```http
GET /api-keys/:id
```

#### 12.4 Revoke API Key
```http
DELETE /api-keys/:id
```

### 13. Monitoring & Metrics

#### 13.1 Get System Health
```http
GET /health
```
**Response:**
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime": 86400,
  "services": {
    "database": "healthy",
    "redis": "healthy",
    "queue": "healthy"
  }
}
```

#### 13.2 Get System Metrics
```http
GET /metrics
```

#### 13.3 Get Queue Status
```http
GET /metrics/queue
```

#### 13.4 Get Execution Statistics
```http
GET /metrics/executions
```
**Query Parameters:**
- `period` (string): hour|day|week|month
- `groupBy` (string): workflow|status|mode

#### 13.5 Get Worker Status
```http
GET /metrics/workers
```

#### 13.6 Get Performance Metrics
```http
GET /metrics/performance
```

### 14. Audit Logs

#### 14.1 List Audit Logs
```http
GET /audit-logs
```
**Query Parameters:**
- `userId` (string): Filter by user
- `action` (string): Filter by action
- `entityType` (string): workflow|credential|user
- `startDate` (ISO 8601)
- `endDate` (ISO 8601)

#### 14.2 Get Audit Log Entry
```http
GET /audit-logs/:id
```

### 15. Settings & Configuration

#### 15.1 Get Instance Settings
```http
GET /settings
```

#### 15.2 Update Instance Settings (Admin)
```http
PUT /settings
```
**Request Body:**
```json
{
  "instanceName": "My n8n Instance",
  "allowSignup": true,
  "requireEmailVerification": true,
  "defaultRole": "user",
  "maxWorkflowsPerUser": 100,
  "executionTimeout": 3600
}
```

#### 15.3 Get SMTP Settings
```http
GET /settings/smtp
```

#### 15.4 Update SMTP Settings
```http
PUT /settings/smtp
```

#### 15.5 Test SMTP Settings
```http
POST /settings/smtp/test
```

### 16. Community & Sharing

#### 16.1 Get Community Workflows
```http
GET /community/workflows
```

#### 16.2 Publish Workflow to Community
```http
POST /community/workflows
```

#### 16.3 Get Workflow Reviews
```http
GET /community/workflows/:id/reviews
```

#### 16.4 Add Workflow Review
```http
POST /community/workflows/:id/reviews
```

#### 16.5 Report Workflow
```http
POST /community/workflows/:id/report
```

### 17. Integrations

#### 17.1 List Available Integrations
```http
GET /integrations
```

#### 17.2 Get Integration Details
```http
GET /integrations/:name
```

#### 17.3 Install Integration
```http
POST /integrations/:name/install
```

#### 17.4 Uninstall Integration
```http
POST /integrations/:name/uninstall
```

#### 17.5 Update Integration
```http
PUT /integrations/:name
```

### 18. WebSocket Endpoints

#### 18.1 Real-time Updates
```websocket
WS /ws
```
**Message Types:**
```json
{
  "type": "execution.started|execution.completed|execution.failed|node.executing|node.completed|workflow.updated",
  "data": {},
  "eventId": "uuid",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

#### 18.2 Subscribe to Workflow
```json
{
  "action": "subscribe",
  "workflowId": "workflow_id"
}
```

#### 18.3 Unsubscribe from Workflow
```json
{
  "action": "unsubscribe",
  "workflowId": "workflow_id"
}
```

### 19. Import/Export

#### 19.1 Export All Workflows
```http
GET /export/workflows
```

#### 19.2 Export All Credentials
```http
GET /export/credentials
```

#### 19.3 Export All Data
```http
GET /export/all
```

#### 19.4 Import Data
```http
POST /import
```
**Request Body (multipart/form-data):**
- `file`: JSON export file
- `overwrite` (boolean): Overwrite existing data

### 20. Search

#### 20.1 Global Search
```http
GET /search
```
**Query Parameters:**
- `q` (string): Search query
- `type` (string): workflow|credential|execution|node
- `limit` (int): Results limit

#### 20.2 Search Workflows
```http
GET /search/workflows
```

#### 20.3 Search Executions
```http
GET /search/executions
```

### 21. Notifications

#### 21.1 Get User Notifications
```http
GET /notifications
```

#### 21.2 Mark Notification as Read
```http
PUT /notifications/:id/read
```

#### 21.3 Mark All as Read
```http
PUT /notifications/read-all
```

#### 21.4 Delete Notification
```http
DELETE /notifications/:id
```

#### 21.5 Get Notification Settings
```http
GET /notifications/settings
```

#### 21.6 Update Notification Settings
```http
PUT /notifications/settings
```

### 22. Scheduler

#### 22.1 List Scheduled Workflows
```http
GET /schedules
```

#### 22.2 Create Schedule
```http
POST /schedules
```
**Request Body:**
```json
{
  "workflowId": "workflow_id",
  "cronExpression": "0 0 * * *",
  "timezone": "America/New_York",
  "active": true
}
```

#### 22.3 Update Schedule
```http
PUT /schedules/:id
```

#### 22.4 Delete Schedule
```http
DELETE /schedules/:id
```

#### 22.5 Activate Schedule
```http
POST /schedules/:id/activate
```

#### 22.6 Deactivate Schedule
```http
POST /schedules/:id/deactivate
```

### 23. Teams & Organizations

#### 23.1 List Teams
```http
GET /teams
```

#### 23.2 Create Team
```http
POST /teams
```

#### 23.3 Get Team
```http
GET /teams/:id
```

#### 23.4 Update Team
```http
PUT /teams/:id
```

#### 23.5 Delete Team
```http
DELETE /teams/:id
```

#### 23.6 Add Team Member
```http
POST /teams/:id/members
```

#### 23.7 Remove Team Member
```http
DELETE /teams/:id/members/:userId
```

#### 23.8 Update Team Member Role
```http
PUT /teams/:id/members/:userId
```

### 24. Billing & Usage (Enterprise)

#### 24.1 Get Usage Statistics
```http
GET /billing/usage
```

#### 24.2 Get Billing Information
```http
GET /billing/info
```

#### 24.3 Get Invoices
```http
GET /billing/invoices
```

#### 24.4 Get Subscription
```http
GET /billing/subscription
```

#### 24.5 Update Subscription
```http
PUT /billing/subscription
```

## Error Responses

All error responses follow this format:
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {
      "field": "Additional error context"
    },
    "requestId": "uuid",
    "timestamp": "2024-01-01T00:00:00Z"
  }
}
```

### Common Error Codes:
- `VALIDATION_ERROR`: Invalid input data
- `UNAUTHORIZED`: Missing or invalid authentication
- `FORBIDDEN`: Insufficient permissions
- `NOT_FOUND`: Resource not found
- `CONFLICT`: Resource already exists
- `RATE_LIMIT`: Too many requests
- `INTERNAL_ERROR`: Server error
- `WORKFLOW_INACTIVE`: Workflow is not active
- `EXECUTION_FAILED`: Execution failed
- `INVALID_CREDENTIALS`: Invalid credentials
- `QUOTA_EXCEEDED`: Usage limit exceeded

## Rate Limiting

API rate limits:
- **Anonymous**: 10 requests/minute
- **Authenticated**: 100 requests/minute
- **Pro**: 500 requests/minute
- **Enterprise**: Unlimited

Rate limit headers:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640995200
```

## Pagination

List endpoints support pagination:
```
GET /workflows?page=1&limit=20
```

Response includes pagination metadata:
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "totalPages": 5,
    "hasNext": true,
    "hasPrev": false
  }
}
```

## Filtering & Sorting

Most list endpoints support:
- `sort`: Field to sort by
- `order`: asc|desc
- `search`: Text search
- `filter[field]`: Field-specific filters

Example:
```
GET /workflows?sort=created_at&order=desc&filter[active]=true&search=sales
```

## Batch Operations

Some endpoints support batch operations:
```http
POST /workflows/batch
```
**Request Body:**
```json
{
  "operation": "activate|deactivate|delete",
  "ids": ["id1", "id2", "id3"]
}
```

This comprehensive API documentation covers all endpoints needed for a complete n8n clone backend implementation.
