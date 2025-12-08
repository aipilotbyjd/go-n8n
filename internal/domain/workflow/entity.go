package workflow

import (
	"time"

	"github.com/google/uuid"
)

// Workflow represents a workflow entity
type Workflow struct {
	ID          uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string                 `json:"name" gorm:"not null"`
	Description string                 `json:"description"`
	UserID      uuid.UUID              `json:"user_id" gorm:"type:uuid;not null"`
	TeamID      *uuid.UUID             `json:"team_id,omitempty" gorm:"type:uuid"`
	IsActive    bool                   `json:"is_active" gorm:"default:false"`
	Nodes       []Node                 `json:"nodes" gorm:"serializer:json"`
	Connections []Connection           `json:"connections" gorm:"serializer:json"`
	Settings    WorkflowSettings       `json:"settings" gorm:"serializer:json"`
	Tags        []string               `json:"tags" gorm:"type:text[]"`
	Version     int                    `json:"version" gorm:"default:1"`
	Variables   map[string]interface{} `json:"variables" gorm:"serializer:json"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DeletedAt   *time.Time             `json:"deleted_at,omitempty" gorm:"index"`
}

// Node represents a node in a workflow
type Node struct {
	ID             string                 `json:"id"`
	Type           string                 `json:"type"`
	Name           string                 `json:"name"`
	Position       NodePosition           `json:"position"`
	Parameters     map[string]interface{} `json:"parameters"`
	CredentialID   *uuid.UUID             `json:"credential_id,omitempty"`
	Disabled       bool                   `json:"disabled"`
	Notes          string                 `json:"notes,omitempty"`
	RetryOnFail    bool                   `json:"retry_on_fail"`
	MaxRetries     int                    `json:"max_retries"`
	WaitBetweenTries int                  `json:"wait_between_tries"` // milliseconds
	ContinueOnFail bool                   `json:"continue_on_fail"`
	ExecuteOnce    bool                   `json:"execute_once"`
}

// NodePosition represents the position of a node on the canvas
type NodePosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Connection represents a connection between nodes
type Connection struct {
	Source     ConnectionPoint `json:"source"`
	Target     ConnectionPoint `json:"target"`
	Data       ConnectionData  `json:"data,omitempty"`
}

// ConnectionPoint represents an endpoint of a connection
type ConnectionPoint struct {
	NodeID string `json:"node_id"`
	Type   string `json:"type"` // "main" or custom output
	Index  int    `json:"index"`
}

// ConnectionData contains additional connection metadata
type ConnectionData struct {
	Disabled bool   `json:"disabled,omitempty"`
	Label    string `json:"label,omitempty"`
}

// WorkflowSettings contains workflow configuration
type WorkflowSettings struct {
	ExecutionOrder    string                 `json:"execution_order"` // "sequential" or "parallel"
	SaveExecutions    bool                   `json:"save_executions"`
	SaveDataOnError   bool                   `json:"save_data_on_error"`
	SaveDataOnSuccess bool                   `json:"save_data_on_success"`
	SaveDataManual    bool                   `json:"save_data_manual"`
	Timezone          string                 `json:"timezone"`
	ErrorWorkflow     *uuid.UUID             `json:"error_workflow,omitempty"`
	MaxExecutionTime  int                    `json:"max_execution_time"` // seconds
	Timeout           int                    `json:"timeout"`             // seconds
	CustomData        map[string]interface{} `json:"custom_data,omitempty"`
}

// WorkflowStatus represents the status of a workflow
type WorkflowStatus string

const (
	WorkflowStatusActive   WorkflowStatus = "active"
	WorkflowStatusInactive WorkflowStatus = "inactive"
	WorkflowStatusError    WorkflowStatus = "error"
)

// Validate validates the workflow entity
func (w *Workflow) Validate() error {
	if w.Name == "" {
		return ErrWorkflowNameRequired
	}
	
	if len(w.Nodes) == 0 {
		return ErrWorkflowNodesRequired
	}
	
	// Validate each node
	for _, node := range w.Nodes {
		if err := node.Validate(); err != nil {
			return err
		}
	}
	
	// Validate connections
	for _, conn := range w.Connections {
		if err := conn.Validate(); err != nil {
			return err
		}
	}
	
	return nil
}

// Validate validates a node
func (n *Node) Validate() error {
	if n.ID == "" {
		return ErrNodeIDRequired
	}
	
	if n.Type == "" {
		return ErrNodeTypeRequired
	}
	
	if n.Name == "" {
		return ErrNodeNameRequired
	}
	
	return nil
}

// Validate validates a connection
func (c *Connection) Validate() error {
	if c.Source.NodeID == "" || c.Target.NodeID == "" {
		return ErrConnectionNodesRequired
	}
	
	if c.Source.NodeID == c.Target.NodeID {
		return ErrConnectionSelfLoop
	}
	
	return nil
}

// Activate activates the workflow
func (w *Workflow) Activate() error {
	if w.IsActive {
		return ErrWorkflowAlreadyActive
	}
	
	if err := w.Validate(); err != nil {
		return err
	}
	
	w.IsActive = true
	w.UpdatedAt = time.Now()
	return nil
}

// Deactivate deactivates the workflow
func (w *Workflow) Deactivate() {
	w.IsActive = false
	w.UpdatedAt = time.Now()
}

// IncrementVersion increments the workflow version
func (w *Workflow) IncrementVersion() {
	w.Version++
	w.UpdatedAt = time.Now()
}

// Clone creates a copy of the workflow
func (w *Workflow) Clone() *Workflow {
	clone := &Workflow{
		ID:          uuid.New(),
		Name:        w.Name + " (Copy)",
		Description: w.Description,
		UserID:      w.UserID,
		TeamID:      w.TeamID,
		IsActive:    false,
		Nodes:       make([]Node, len(w.Nodes)),
		Connections: make([]Connection, len(w.Connections)),
		Settings:    w.Settings,
		Tags:        make([]string, len(w.Tags)),
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	copy(clone.Nodes, w.Nodes)
	copy(clone.Connections, w.Connections)
	copy(clone.Tags, w.Tags)
	
	return clone
}
