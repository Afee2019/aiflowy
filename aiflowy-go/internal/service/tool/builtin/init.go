package builtin

import (
	aitool "github.com/aiflowy/aiflowy-go/internal/service/tool"
)

// RegisterAll registers all builtin tools to the global registry
func RegisterAll() error {
	registry := aitool.GetRegistry()

	tools := []aitool.Tool{
		NewTimeTool(),
		NewCalculatorTool(),
		NewRandomTool(),
	}

	for _, t := range tools {
		if err := registry.Register(t); err != nil {
			return err
		}
	}

	return nil
}

// GetBuiltinTools returns all builtin tool instances
func GetBuiltinTools() []aitool.Tool {
	return []aitool.Tool{
		NewTimeTool(),
		NewCalculatorTool(),
		NewRandomTool(),
	}
}
