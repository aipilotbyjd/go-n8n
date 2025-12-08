package workflow

import "errors"

var (
	// Workflow errors
	ErrWorkflowNotFound      = errors.New("workflow not found")
	ErrWorkflowNameRequired  = errors.New("workflow name is required")
	ErrWorkflowNodesRequired = errors.New("workflow must have at least one node")
	ErrWorkflowAlreadyActive = errors.New("workflow is already active")
	ErrWorkflowNotActive     = errors.New("workflow is not active")
	ErrWorkflowInvalid       = errors.New("workflow configuration is invalid")
	
	// Node errors
	ErrNodeNotFound      = errors.New("node not found")
	ErrNodeIDRequired    = errors.New("node ID is required")
	ErrNodeTypeRequired  = errors.New("node type is required")
	ErrNodeNameRequired  = errors.New("node name is required")
	ErrNodeTypeInvalid   = errors.New("node type is invalid")
	ErrNodeConfigInvalid = errors.New("node configuration is invalid")
	
	// Connection errors
	ErrConnectionNodesRequired = errors.New("connection source and target nodes are required")
	ErrConnectionSelfLoop      = errors.New("connection cannot create a self-loop")
	ErrConnectionInvalid       = errors.New("connection configuration is invalid")
	ErrConnectionDuplicate     = errors.New("connection already exists")
	
	// Execution errors
	ErrWorkflowExecutionFailed = errors.New("workflow execution failed")
	ErrWorkflowCycleDetected   = errors.New("workflow contains a cycle")
)
