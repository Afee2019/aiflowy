package tool

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/schema"
)

// MockTool is a simple tool for testing
type MockTool struct {
	name        string
	description string
	params      map[string]*schema.ParameterInfo
	execFn      func(ctx context.Context, args map[string]interface{}) (interface{}, error)
}

func (m *MockTool) Name() string {
	return m.name
}

func (m *MockTool) Description() string {
	return m.description
}

func (m *MockTool) Parameters() map[string]*schema.ParameterInfo {
	return m.params
}

func (m *MockTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	if m.execFn != nil {
		return m.execFn(ctx, args)
	}
	return "mock result", nil
}

func NewMockTool(name, desc string) *MockTool {
	return &MockTool{
		name:        name,
		description: desc,
		params: map[string]*schema.ParameterInfo{
			"input": {
				Type: "string",
				Desc: "Input parameter",
			},
		},
	}
}

func TestToolWrapper_Info(t *testing.T) {
	ctx := context.Background()
	mockTool := NewMockTool("test_tool", "A test tool")
	wrapper := NewToolWrapper(mockTool)

	info, err := wrapper.Info(ctx)
	if err != nil {
		t.Fatalf("failed to get info: %v", err)
	}

	if info.Name != "test_tool" {
		t.Errorf("expected name 'test_tool', got '%s'", info.Name)
	}
	if info.Desc != "A test tool" {
		t.Errorf("expected desc 'A test tool', got '%s'", info.Desc)
	}
}

func TestToolWrapper_InvokableRun(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		argsJSON      string
		execFn        func(ctx context.Context, args map[string]interface{}) (interface{}, error)
		expectResult  string
		expectError   bool
	}{
		{
			name:         "empty args",
			argsJSON:     "",
			expectResult: "mock result",
			expectError:  false,
		},
		{
			name:         "valid JSON args",
			argsJSON:     `{"input": "test"}`,
			expectResult: "mock result",
			expectError:  false,
		},
		{
			name:        "invalid JSON args",
			argsJSON:    `{invalid`,
			expectError: true,
		},
		{
			name:     "returns map result",
			argsJSON: "",
			execFn: func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
				return map[string]string{"key": "value"}, nil
			},
			expectResult: `{"key":"value"}`,
			expectError:  false,
		},
		{
			name:     "returns string result",
			argsJSON: "",
			execFn: func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
				return "direct string", nil
			},
			expectResult: "direct string",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTool := NewMockTool("test", "test")
			mockTool.execFn = tt.execFn
			wrapper := NewToolWrapper(mockTool)

			result, err := wrapper.InvokableRun(ctx, tt.argsJSON)
			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			if !tt.expectError && result != tt.expectResult {
				t.Errorf("expected result '%s', got '%s'", tt.expectResult, result)
			}
		})
	}
}

func TestRegistry_Register(t *testing.T) {
	registry := &Registry{tools: make(map[string]Tool)}

	tool1 := NewMockTool("tool1", "First tool")
	tool2 := NewMockTool("tool1", "Duplicate tool")

	// First registration should succeed
	err := registry.Register(tool1)
	if err != nil {
		t.Errorf("first registration should succeed: %v", err)
	}

	// Duplicate registration should fail
	err = registry.Register(tool2)
	if err == nil {
		t.Error("duplicate registration should fail")
	}
}

func TestRegistry_Get(t *testing.T) {
	registry := &Registry{tools: make(map[string]Tool)}

	tool := NewMockTool("test_tool", "Test tool")
	registry.Register(tool)

	// Get existing tool
	found, ok := registry.Get("test_tool")
	if !ok {
		t.Error("expected to find registered tool")
	}
	if found.Name() != "test_tool" {
		t.Errorf("expected name 'test_tool', got '%s'", found.Name())
	}

	// Get non-existing tool
	_, ok = registry.Get("non_existing")
	if ok {
		t.Error("expected not to find non-existing tool")
	}
}

func TestRegistry_GetAll(t *testing.T) {
	registry := &Registry{tools: make(map[string]Tool)}

	registry.Register(NewMockTool("tool1", "Tool 1"))
	registry.Register(NewMockTool("tool2", "Tool 2"))
	registry.Register(NewMockTool("tool3", "Tool 3"))

	all := registry.GetAll()
	if len(all) != 3 {
		t.Errorf("expected 3 tools, got %d", len(all))
	}
}

func TestRegistry_GetToolInfos(t *testing.T) {
	ctx := context.Background()
	registry := &Registry{tools: make(map[string]Tool)}

	registry.Register(NewMockTool("tool1", "Tool 1"))
	registry.Register(NewMockTool("tool2", "Tool 2"))

	infos, err := registry.GetToolInfos(ctx)
	if err != nil {
		t.Fatalf("failed to get tool infos: %v", err)
	}

	if len(infos) != 2 {
		t.Errorf("expected 2 tool infos, got %d", len(infos))
	}
}

func TestRegistry_GetInvokableTools(t *testing.T) {
	registry := &Registry{tools: make(map[string]Tool)}

	registry.Register(NewMockTool("tool1", "Tool 1"))
	registry.Register(NewMockTool("tool2", "Tool 2"))

	invokable := registry.GetInvokableTools()
	if len(invokable) != 2 {
		t.Errorf("expected 2 invokable tools, got %d", len(invokable))
	}
}

func TestRegistry_Execute(t *testing.T) {
	ctx := context.Background()
	registry := &Registry{tools: make(map[string]Tool)}

	tool := NewMockTool("test_tool", "Test tool")
	tool.execFn = func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		return "executed with: " + args["input"].(string), nil
	}
	registry.Register(tool)

	// Execute existing tool
	result, err := registry.Execute(ctx, "test_tool", `{"input": "test"}`)
	if err != nil {
		t.Fatalf("failed to execute tool: %v", err)
	}
	if result != "executed with: test" {
		t.Errorf("unexpected result: %s", result)
	}

	// Execute non-existing tool
	_, err = registry.Execute(ctx, "non_existing", "")
	if err == nil {
		t.Error("expected error for non-existing tool")
	}
}

func TestRegistry_GetToolsByNames(t *testing.T) {
	registry := &Registry{tools: make(map[string]Tool)}

	registry.Register(NewMockTool("tool1", "Tool 1"))
	registry.Register(NewMockTool("tool2", "Tool 2"))
	registry.Register(NewMockTool("tool3", "Tool 3"))

	// Get existing tools
	tools, err := registry.GetToolsByNames([]string{"tool1", "tool3"})
	if err != nil {
		t.Fatalf("failed to get tools: %v", err)
	}
	if len(tools) != 2 {
		t.Errorf("expected 2 tools, got %d", len(tools))
	}

	// Get with non-existing tool
	_, err = registry.GetToolsByNames([]string{"tool1", "non_existing"})
	if err == nil {
		t.Error("expected error for non-existing tool")
	}
}

func TestRegistry_Clear(t *testing.T) {
	registry := &Registry{tools: make(map[string]Tool)}

	registry.Register(NewMockTool("tool1", "Tool 1"))
	registry.Register(NewMockTool("tool2", "Tool 2"))

	registry.Clear()

	all := registry.GetAll()
	if len(all) != 0 {
		t.Errorf("expected 0 tools after clear, got %d", len(all))
	}
}

func TestGetRegistry(t *testing.T) {
	r1 := GetRegistry()
	r2 := GetRegistry()

	if r1 != r2 {
		t.Error("expected same global registry instance")
	}
}

func BenchmarkRegistry_Get(b *testing.B) {
	registry := &Registry{tools: make(map[string]Tool)}
	for i := 0; i < 100; i++ {
		registry.Register(NewMockTool(string(rune('A'+i)), "Tool"))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		registry.Get("M")
	}
}

func BenchmarkToolWrapper_InvokableRun(b *testing.B) {
	ctx := context.Background()
	mockTool := NewMockTool("test", "test")
	wrapper := NewToolWrapper(mockTool)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = wrapper.InvokableRun(ctx, `{"input": "test"}`)
	}
}
