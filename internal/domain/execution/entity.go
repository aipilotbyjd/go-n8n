package execution

import (
	"time"

	"github.com/google/uuid"
)

// Execution represents a workflow execution
type Execution struct {
	ID              uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	WorkflowID      uuid.UUID              `json:"workflow_id" gorm:"type:uuid;not null"`
	WorkflowVersion int                    `json:"workflow_version" gorm:"not null"`
	Status          ExecutionStatus        `json:"status" gorm:"not null"`
	Mode            ExecutionMode          `json:"mode" gorm:"not null"`
	StartedAt       time.Time              `json:"started_at"`
	FinishedAt      *time.Time             `json:"finished_at,omitempty"`
	ExecutionTimeMs int                    `json:"execution_time_ms,omitempty"`
	InputData       map[string]interface{} `json:"input_data" gorm:"serializer:json"`
	OutputData      map[string]interface{} `json:"output_data,omitempty" gorm:"serializer:json"`
	ErrorMessage    string                 `json:"error_message,omitempty"`
	ErrorNode       string                 `json:"error_node,omitempty"`
	RetryOf         *uuid.UUID             `json:"retry_of,omitempty" gorm:"type:uuid"`
	RetryCount      int                    `json:"retry_count" gorm:"default:0"`
	CreatedAt       time.Time              `json:"created_at"`
}

// ExecutionStatus represents the status of an execution
type ExecutionStatus string

const (
	ExecutionStatusWaiting   ExecutionStatus = "waiting"
	ExecutionStatusRunning   ExecutionStatus = "running"
	ExecutionStatusSuccess   ExecutionStatus = "success"
	ExecutionStatusError     ExecutionStatus = "error"
	ExecutionStatusCancelled ExecutionStatus = "cancelled"
	ExecutionStatusCrashed   ExecutionStatus = "crashed"
	ExecutionStatusTimeout   ExecutionStatus = "timeout"
)

// ExecutionMode represents how the execution was triggered
type ExecutionMode string

const (
	ExecutionModeManual   ExecutionMode = "manual"
	ExecutionModeTrigger  ExecutionMode = "trigger"
	ExecutionModeWebhook  ExecutionMode = "webhook"
	ExecutionModeSchedule ExecutionMode = "schedule"
	ExecutionModeRetry    ExecutionMode = "retry"
	ExecutionModeTest     ExecutionMode = "test"
)

// NodeExecution represents the execution state of a single node
type NodeExecution struct {
	ID              uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ExecutionID     uuid.UUID              `json:"execution_id" gorm:"type:uuid;not null"`
	NodeID          string                 `json:"node_id" gorm:"not null"`
	NodeType        string                 `json:"node_type" gorm:"not null"`
	NodeName        string                 `json:"node_name"`
	Status          ExecutionStatus        `json:"status" gorm:"not null"`
	InputData       map[string]interface{} `json:"input_data" gorm:"serializer:json"`
	OutputData      map[string]interface{} `json:"output_data,omitempty" gorm:"serializer:json"`
	ErrorMessage    string                 `json:"error_message,omitempty"`
	ExecutionTimeMs int                    `json:"execution_time_ms,omitempty"`
	StartedAt       time.Time              `json:"started_at"`
	FinishedAt      *time.Time             `json:"finished_at,omitempty"`
	RetryCount      int                    `json:"retry_count" gorm:"default:0"`
}

// ExecutionContext holds the runtime context for an execution
type ExecutionContext struct {
	ExecutionID     uuid.UUID              `json:"execution_id"`
	WorkflowID      uuid.UUID              `json:"workflow_id"`
	Mode            ExecutionMode          `json:"mode"`
	Variables       map[string]interface{} `json:"variables"`
	GlobalData      map[string]interface{} `json:"global_data"`
	NodeOutputs     map[string]interface{} `json:"node_outputs"`
	Credentials     map[string]interface{} `json:"credentials"`
	Timezone        string                 `json:"timezone"`
	StartTime       time.Time              `json:"start_time"`
	MaxExecutionTime time.Duration         `json:"max_execution_time"`
	RetryPolicy     RetryPolicy            `json:"retry_policy"`
}

// RetryPolicy defines retry behavior for failed executions
type RetryPolicy struct {
	MaxRetries       int           `json:"max_retries"`
	RetryInterval    time.Duration `json:"retry_interval"`
	BackoffFactor    float64       `json:"backoff_factor"`
	MaxRetryInterval time.Duration `json:"max_retry_interval"`
}

// ExecutionStatistics holds execution metrics
type ExecutionStatistics struct {
	TotalExecutions int                    `json:"total_executions"`
	SuccessCount    int                    `json:"success_count"`
	ErrorCount      int                    `json:"error_count"`
	AverageTimeMs   int                    `json:"average_time_ms"`
	LastExecution   *time.Time             `json:"last_execution,omitempty"`
	NodeStats       map[string]NodeStats   `json:"node_stats"`
}

// NodeStats holds statistics for a specific node
type NodeStats struct {
	ExecutionCount  int    `json:"execution_count"`
	SuccessCount    int    `json:"success_count"`
	ErrorCount      int    `json:"error_count"`
	AverageTimeMs   int    `json:"average_time_ms"`
	LastError       string `json:"last_error,omitempty"`
}

// IsTerminal returns whether the execution is in a terminal state
func (s ExecutionStatus) IsTerminal() bool {
	switch s {
	case ExecutionStatusSuccess, ExecutionStatusError, ExecutionStatusCancelled, ExecutionStatusCrashed, ExecutionStatusTimeout:
		return true
	default:
		return false
	}
}

// IsRunning returns whether the execution is currently running
func (s ExecutionStatus) IsRunning() bool {
	return s == ExecutionStatusRunning || s == ExecutionStatusWaiting
}

// Start marks the execution as started
func (e *Execution) Start() {
	e.Status = ExecutionStatusRunning
	e.StartedAt = time.Now()
}

// Complete marks the execution as completed successfully
func (e *Execution) Complete(outputData map[string]interface{}) {
	e.Status = ExecutionStatusSuccess
	e.OutputData = outputData
	e.finish()
}

// Fail marks the execution as failed
func (e *Execution) Fail(err error, nodeID string) {
	e.Status = ExecutionStatusError
	e.ErrorMessage = err.Error()
	e.ErrorNode = nodeID
	e.finish()
}

// Cancel marks the execution as cancelled
func (e *Execution) Cancel() {
	e.Status = ExecutionStatusCancelled
	e.finish()
}

// Timeout marks the execution as timed out
func (e *Execution) Timeout() {
	e.Status = ExecutionStatusTimeout
	e.finish()
}

// finish sets the finish time and calculates execution time
func (e *Execution) finish() {
	now := time.Now()
	e.FinishedAt = &now
	e.ExecutionTimeMs = int(now.Sub(e.StartedAt).Milliseconds())
}

// CanRetry returns whether the execution can be retried
func (e *Execution) CanRetry(policy RetryPolicy) bool {
	return e.Status == ExecutionStatusError && e.RetryCount < policy.MaxRetries
}

// CreateRetry creates a new execution as a retry of this one
func (e *Execution) CreateRetry() *Execution {
	retry := &Execution{
		ID:              uuid.New(),
		WorkflowID:      e.WorkflowID,
		WorkflowVersion: e.WorkflowVersion,
		Status:          ExecutionStatusWaiting,
		Mode:            ExecutionModeRetry,
		InputData:       e.InputData,
		RetryOf:         &e.ID,
		RetryCount:      e.RetryCount + 1,
		CreatedAt:       time.Now(),
	}
	return retry
}
