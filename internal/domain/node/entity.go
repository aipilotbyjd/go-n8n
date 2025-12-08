package node

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

// NodeInterface defines the interface all nodes must implement
type NodeInterface interface {
	// Core methods
	GetType() string
	GetName() string
	GetCategory() Category
	GetVersion() string
	GetDescription() string
	GetIcon() string
	
	// Execution
	Execute(ctx context.Context, input *NodeInput) (*NodeOutput, error)
	
	// Validation
	Validate(parameters map[string]interface{}) error
	
	// Schema
	GetSchema() *NodeSchema
	GetCredentialTypes() []string
	
	// Configuration
	GetDefaultParameters() map[string]interface{}
}

// Category represents node category
type Category string

const (
	CategoryTrigger      Category = "trigger"
	CategoryAction       Category = "action"
	CategoryTransform    Category = "transform"
	CategoryFlow         Category = "flow"
	CategoryIntegration  Category = "integration"
	CategoryUtility      Category = "utility"
)

// NodeInput represents input data for node execution
type NodeInput struct {
	Data        []Item                 `json:"data"`
	Parameters  map[string]interface{} `json:"parameters"`
	Credentials map[string]interface{} `json:"credentials"`
	Context     *ExecutionContext      `json:"context"`
}

// NodeOutput represents output data from node execution
type NodeOutput struct {
	Data     []Item                 `json:"data"`
	Error    error                  `json:"error,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Item represents a single data item
type Item struct {
	JSON   map[string]interface{} `json:"json"`
	Binary map[string]Binary      `json:"binary,omitempty"`
}

// Binary represents binary data
type Binary struct {
	Data      []byte `json:"data,omitempty"`
	MimeType  string `json:"mime_type"`
	FileName  string `json:"file_name"`
	FileSize  int64  `json:"file_size"`
	DataURI   string `json:"data_uri,omitempty"`
	ID        string `json:"id,omitempty"`
}

// ExecutionContext provides context for node execution
type ExecutionContext struct {
	WorkflowID    string                 `json:"workflow_id"`
	ExecutionID   string                 `json:"execution_id"`
	NodeID        string                 `json:"node_id"`
	RunIndex      int                    `json:"run_index"`
	ItemIndex     int                    `json:"item_index"`
	ActiveNode    string                 `json:"active_node"`
	Variables     map[string]interface{} `json:"variables"`
	Mode          string                 `json:"mode"`
	Timezone      string                 `json:"timezone"`
	RetryCount    int                    `json:"retry_count"`
	MaxRetries    int                    `json:"max_retries"`
}

// NodeSchema defines the structure and properties of a node
type NodeSchema struct {
	Type        string             `json:"type"`
	Name        string             `json:"name"`
	Group       []string           `json:"group"`
	Version     float64            `json:"version"`
	Description string             `json:"description"`
	Icon        string             `json:"icon"`
	IconColor   string             `json:"icon_color,omitempty"`
	Defaults    NodeDefaults       `json:"defaults"`
	Inputs      []IOSchema         `json:"inputs"`
	Outputs     []IOSchema         `json:"outputs"`
	Properties  []PropertySchema   `json:"properties"`
	Credentials []CredentialSchema `json:"credentials"`
}

// NodeDefaults contains default values for a node
type NodeDefaults struct {
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
}

// IOSchema defines input/output schema
type IOSchema struct {
	Type     string   `json:"type"`
	Required bool     `json:"required"`
	Label    string   `json:"label,omitempty"`
	Options  []string `json:"options,omitempty"`
}

// PropertySchema defines a node property
type PropertySchema struct {
	Name         string                 `json:"name"`
	DisplayName  string                 `json:"display_name"`
	Type         PropertyType           `json:"type"`
	Default      interface{}            `json:"default,omitempty"`
	Required     bool                   `json:"required"`
	Description  string                 `json:"description,omitempty"`
	Hint         string                 `json:"hint,omitempty"`
	Options      []PropertyOption       `json:"options,omitempty"`
	DisplayOptions *DisplayOptions      `json:"display_options,omitempty"`
	Validation   *PropertyValidation    `json:"validation,omitempty"`
	TypeOptions  map[string]interface{} `json:"type_options,omitempty"`
}

// PropertyType represents the type of a property
type PropertyType string

const (
	PropertyTypeString     PropertyType = "string"
	PropertyTypeNumber     PropertyType = "number"
	PropertyTypeBoolean    PropertyType = "boolean"
	PropertyTypeOptions    PropertyType = "options"
	PropertyTypeMultiOptions PropertyType = "multi_options"
	PropertyTypeJSON       PropertyType = "json"
	PropertyTypeCode       PropertyType = "code"
	PropertyTypeDateTime   PropertyType = "datetime"
	PropertyTypeCollection PropertyType = "collection"
	PropertyTypeFixed      PropertyType = "fixed_collection"
	PropertyTypeColor      PropertyType = "color"
	PropertyTypeFile       PropertyType = "file"
	PropertyTypeHidden     PropertyType = "hidden"
)

// PropertyOption represents an option for select properties
type PropertyOption struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

// DisplayOptions controls when a property is displayed
type DisplayOptions struct {
	Show map[string][]interface{} `json:"show,omitempty"`
	Hide map[string][]interface{} `json:"hide,omitempty"`
}

// PropertyValidation defines validation rules for a property
type PropertyValidation struct {
	Min       *float64 `json:"min,omitempty"`
	Max       *float64 `json:"max,omitempty"`
	MinLength *int     `json:"min_length,omitempty"`
	MaxLength *int     `json:"max_length,omitempty"`
	Pattern   string   `json:"pattern,omitempty"`
	Custom    string   `json:"custom,omitempty"`
}

// CredentialSchema defines required credentials for a node
type CredentialSchema struct {
	Name     string   `json:"name"`
	Required bool     `json:"required"`
	Types    []string `json:"types,omitempty"`
}

// NodeRegistration holds node registration information
type NodeRegistration struct {
	Type        string
	Category    Category
	Constructor func() NodeInterface
}

// NodeRegistry manages all registered nodes
type NodeRegistry struct {
	nodes map[string]NodeRegistration
}

// NewNodeRegistry creates a new node registry
func NewNodeRegistry() *NodeRegistry {
	return &NodeRegistry{
		nodes: make(map[string]NodeRegistration),
	}
}

// Register registers a new node type
func (r *NodeRegistry) Register(nodeType string, category Category, constructor func() NodeInterface) error {
	if _, exists := r.nodes[nodeType]; exists {
		return errors.New("node type already registered: " + nodeType)
	}
	
	r.nodes[nodeType] = NodeRegistration{
		Type:        nodeType,
		Category:    category,
		Constructor: constructor,
	}
	
	return nil
}

// Get retrieves a node constructor by type
func (r *NodeRegistry) Get(nodeType string) (func() NodeInterface, error) {
	registration, exists := r.nodes[nodeType]
	if !exists {
		return nil, errors.New("node type not found: " + nodeType)
	}
	return registration.Constructor, nil
}

// List returns all registered node types
func (r *NodeRegistry) List() []NodeRegistration {
	list := make([]NodeRegistration, 0, len(r.nodes))
	for _, reg := range r.nodes {
		list = append(list, reg)
	}
	return list
}

// ListByCategory returns nodes filtered by category
func (r *NodeRegistry) ListByCategory(category Category) []NodeRegistration {
	var list []NodeRegistration
	for _, reg := range r.nodes {
		if reg.Category == category {
			list = append(list, reg)
		}
	}
	return list
}

// BaseNode provides common functionality for all nodes
type BaseNode struct {
	Type        string
	Name        string
	Category    Category
	Version     string
	Description string
	Icon        string
}

// GetType returns the node type
func (n *BaseNode) GetType() string {
	return n.Type
}

// GetName returns the node name
func (n *BaseNode) GetName() string {
	return n.Name
}

// GetCategory returns the node category
func (n *BaseNode) GetCategory() Category {
	return n.Category
}

// GetVersion returns the node version
func (n *BaseNode) GetVersion() string {
	return n.Version
}

// GetDescription returns the node description
func (n *BaseNode) GetDescription() string {
	return n.Description
}

// GetIcon returns the node icon
func (n *BaseNode) GetIcon() string {
	return n.Icon
}

// MarshalJSON marshals an Item to JSON
func (i Item) MarshalJSON() ([]byte, error) {
	type Alias Item
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&i),
	})
}

// UnmarshalJSON unmarshals JSON to an Item
func (i *Item) UnmarshalJSON(data []byte) error {
	type Alias Item
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(i),
	}
	return json.Unmarshal(data, &aux)
}

// NodeExecutionData holds data about node execution
type NodeExecutionData struct {
	NodeID          string        `json:"node_id"`
	NodeType        string        `json:"node_type"`
	StartTime       time.Time     `json:"start_time"`
	EndTime         time.Time     `json:"end_time"`
	ExecutionTimeMs int64         `json:"execution_time_ms"`
	Status          string        `json:"status"`
	Error           string        `json:"error,omitempty"`
	InputItems      int           `json:"input_items"`
	OutputItems     int           `json:"output_items"`
}
