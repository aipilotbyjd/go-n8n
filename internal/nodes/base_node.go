package nodes

import (
	"context"
	"errors"
	"fmt"

	"github.com/jaydeep/go-n8n/internal/domain/node"
)

// BaseNode provides common implementation for all nodes
type BaseNode struct {
	Type        string
	Name        string
	Category    node.Category
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
func (n *BaseNode) GetCategory() node.Category {
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

// GetCredentialTypes returns empty list by default
func (n *BaseNode) GetCredentialTypes() []string {
	return []string{}
}

// GetDefaultParameters returns empty map by default
func (n *BaseNode) GetDefaultParameters() map[string]interface{} {
	return make(map[string]interface{})
}

// ValidateRequired validates required parameters
func ValidateRequired(parameters map[string]interface{}, required []string) error {
	for _, key := range required {
		if _, exists := parameters[key]; !exists {
			return fmt.Errorf("required parameter '%s' is missing", key)
		}
	}
	return nil
}

// GetString gets a string parameter with default value
func GetString(parameters map[string]interface{}, key string, defaultValue string) string {
	if val, exists := parameters[key]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

// GetInt gets an integer parameter with default value
func GetInt(parameters map[string]interface{}, key string, defaultValue int) int {
	if val, exists := parameters[key]; exists {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case float32:
			return int(v)
		}
	}
	return defaultValue
}

// GetBool gets a boolean parameter with default value
func GetBool(parameters map[string]interface{}, key string, defaultValue bool) bool {
	if val, exists := parameters[key]; exists {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return defaultValue
}

// GetStringSlice gets a string slice parameter
func GetStringSlice(parameters map[string]interface{}, key string) []string {
	if val, exists := parameters[key]; exists {
		switch v := val.(type) {
		case []string:
			return v
		case []interface{}:
			result := make([]string, len(v))
			for i, item := range v {
				if str, ok := item.(string); ok {
					result[i] = str
				}
			}
			return result
		}
	}
	return []string{}
}

// GetMap gets a map parameter
func GetMap(parameters map[string]interface{}, key string) map[string]interface{} {
	if val, exists := parameters[key]; exists {
		if m, ok := val.(map[string]interface{}); ok {
			return m
		}
	}
	return make(map[string]interface{})
}

// ProcessItems applies a function to each input item
func ProcessItems(ctx context.Context, input *node.NodeInput, fn func(context.Context, node.Item, int) (node.Item, error)) (*node.NodeOutput, error) {
	output := &node.NodeOutput{
		Data:     make([]node.Item, 0, len(input.Data)),
		Metadata: make(map[string]interface{}),
	}

	for i, item := range input.Data {
		select {
		case <-ctx.Done():
			return nil, errors.New("execution cancelled")
		default:
			processedItem, err := fn(ctx, item, i)
			if err != nil {
				output.Error = err
				return output, err
			}
			output.Data = append(output.Data, processedItem)
		}
	}

	return output, nil
}

// MergeItems merges multiple items into one
func MergeItems(items []node.Item) node.Item {
	merged := node.Item{
		JSON:   make(map[string]interface{}),
		Binary: make(map[string]node.Binary),
	}

	for _, item := range items {
		// Merge JSON data
		for k, v := range item.JSON {
			merged.JSON[k] = v
		}
		// Merge binary data
		for k, v := range item.Binary {
			merged.Binary[k] = v
		}
	}

	return merged
}

// SplitItems splits items based on a key
func SplitItems(items []node.Item, key string) map[string][]node.Item {
	groups := make(map[string][]node.Item)
	
	for _, item := range items {
		if val, exists := item.JSON[key]; exists {
			groupKey := fmt.Sprintf("%v", val)
			groups[groupKey] = append(groups[groupKey], item)
		} else {
			groups["_undefined"] = append(groups["_undefined"], item)
		}
	}
	
	return groups
}

// FilterItems filters items based on a condition
func FilterItems(items []node.Item, condition func(node.Item) bool) []node.Item {
	filtered := make([]node.Item, 0)
	for _, item := range items {
		if condition(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// TransformItem transforms a single item
func TransformItem(item node.Item, transform func(map[string]interface{}) map[string]interface{}) node.Item {
	return node.Item{
		JSON:   transform(item.JSON),
		Binary: item.Binary,
	}
}

// CreateSingleItem creates a single item output
func CreateSingleItem(data map[string]interface{}) *node.NodeOutput {
	return &node.NodeOutput{
		Data: []node.Item{
			{
				JSON:   data,
				Binary: make(map[string]node.Binary),
			},
		},
		Metadata: make(map[string]interface{}),
	}
}

// CreateErrorOutput creates an error output
func CreateErrorOutput(err error) *node.NodeOutput {
	return &node.NodeOutput{
		Data:     []node.Item{},
		Error:    err,
		Metadata: map[string]interface{}{"error": err.Error()},
	}
}

// CreateEmptyOutput creates an empty output
func CreateEmptyOutput() *node.NodeOutput {
	return &node.NodeOutput{
		Data:     []node.Item{},
		Metadata: make(map[string]interface{}),
	}
}
