package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// Tool represents an AIFlowy tool that can be invoked by LLM
type Tool interface {
	// Name returns the unique name of the tool
	Name() string

	// Description returns a description of what the tool does
	Description() string

	// Parameters returns the JSON schema for the tool parameters
	Parameters() map[string]*schema.ParameterInfo

	// Execute runs the tool with the given arguments
	Execute(ctx context.Context, args map[string]interface{}) (interface{}, error)
}

// ToolWrapper wraps an AIFlowy Tool to implement Eino's InvokableTool interface
type ToolWrapper struct {
	tool Tool
}

// NewToolWrapper creates a new ToolWrapper
func NewToolWrapper(t Tool) *ToolWrapper {
	return &ToolWrapper{tool: t}
}

// Info returns the tool information for Eino
func (w *ToolWrapper) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name:       w.tool.Name(),
		Desc:       w.tool.Description(),
		ParamsOneOf: schema.NewParamsOneOfByParams(w.tool.Parameters()),
	}, nil
}

// InvokableRun executes the tool with JSON arguments
func (w *ToolWrapper) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// Parse JSON arguments
	var args map[string]interface{}
	if argumentsInJSON != "" {
		if err := json.Unmarshal([]byte(argumentsInJSON), &args); err != nil {
			return "", fmt.Errorf("failed to parse tool arguments: %w", err)
		}
	} else {
		args = make(map[string]interface{})
	}

	// Execute the tool
	result, err := w.tool.Execute(ctx, args)
	if err != nil {
		return "", err
	}

	// Convert result to JSON string
	if s, ok := result.(string); ok {
		return s, nil
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal tool result: %w", err)
	}

	return string(resultJSON), nil
}

// Ensure ToolWrapper implements InvokableTool
var _ tool.InvokableTool = (*ToolWrapper)(nil)

// Registry manages tool registration and lookup
type Registry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

// Global registry instance
var globalRegistry = &Registry{
	tools: make(map[string]Tool),
}

// GetRegistry returns the global tool registry
func GetRegistry() *Registry {
	return globalRegistry
}

// Register adds a tool to the registry
func (r *Registry) Register(t Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := t.Name()
	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool already registered: %s", name)
	}

	r.tools[name] = t
	return nil
}

// Get retrieves a tool by name
func (r *Registry) Get(name string) (Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.tools[name]
	return t, ok
}

// GetAll returns all registered tools
func (r *Registry) GetAll() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]Tool, 0, len(r.tools))
	for _, t := range r.tools {
		tools = append(tools, t)
	}
	return tools
}

// GetToolInfos returns all tool infos for LLM binding
func (r *Registry) GetToolInfos(ctx context.Context) ([]*schema.ToolInfo, error) {
	tools := r.GetAll()
	infos := make([]*schema.ToolInfo, 0, len(tools))

	for _, t := range tools {
		wrapper := NewToolWrapper(t)
		info, err := wrapper.Info(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get tool info for %s: %w", t.Name(), err)
		}
		infos = append(infos, info)
	}

	return infos, nil
}

// GetInvokableTools returns all tools as InvokableTool for ToolsNode
func (r *Registry) GetInvokableTools() []tool.InvokableTool {
	tools := r.GetAll()
	invokable := make([]tool.InvokableTool, 0, len(tools))

	for _, t := range tools {
		invokable = append(invokable, NewToolWrapper(t))
	}

	return invokable
}

// Execute runs a tool by name with the given arguments
func (r *Registry) Execute(ctx context.Context, name string, argsJSON string) (string, error) {
	t, ok := r.Get(name)
	if !ok {
		return "", fmt.Errorf("tool not found: %s", name)
	}

	wrapper := NewToolWrapper(t)
	return wrapper.InvokableRun(ctx, argsJSON)
}

// GetToolsByNames returns tools by their names
func (r *Registry) GetToolsByNames(names []string) ([]Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]Tool, 0, len(names))
	for _, name := range names {
		t, ok := r.tools[name]
		if !ok {
			return nil, fmt.Errorf("tool not found: %s", name)
		}
		tools = append(tools, t)
	}
	return tools, nil
}

// GetToolInfosByNames returns tool infos for specific tool names
func (r *Registry) GetToolInfosByNames(ctx context.Context, names []string) ([]*schema.ToolInfo, error) {
	tools, err := r.GetToolsByNames(names)
	if err != nil {
		return nil, err
	}

	infos := make([]*schema.ToolInfo, 0, len(tools))
	for _, t := range tools {
		wrapper := NewToolWrapper(t)
		info, err := wrapper.Info(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get tool info for %s: %w", t.Name(), err)
		}
		infos = append(infos, info)
	}

	return infos, nil
}

// Clear removes all tools from the registry (mainly for testing)
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tools = make(map[string]Tool)
}
