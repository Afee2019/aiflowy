package builtin

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/cloudwego/eino/schema"

	aitool "github.com/aiflowy/aiflowy-go/internal/service/tool"
)

// CalculatorTool performs basic mathematical calculations
type CalculatorTool struct{}

// NewCalculatorTool creates a new CalculatorTool
func NewCalculatorTool() *CalculatorTool {
	return &CalculatorTool{}
}

// Name returns the tool name
func (t *CalculatorTool) Name() string {
	return "calculator"
}

// Description returns what the tool does
func (t *CalculatorTool) Description() string {
	return "执行基本的数学计算。支持加法(add)、减法(subtract)、乘法(multiply)、除法(divide)、幂运算(power)、平方根(sqrt)等操作。"
}

// Parameters returns the tool parameters schema
func (t *CalculatorTool) Parameters() map[string]*schema.ParameterInfo {
	return map[string]*schema.ParameterInfo{
		"operation": {
			Type:     schema.String,
			Desc:     "运算类型: add(加)、subtract(减)、multiply(乘)、divide(除)、power(幂)、sqrt(平方根)、percent(百分比)",
			Required: true,
		},
		"a": {
			Type:     schema.Number,
			Desc:     "第一个操作数",
			Required: true,
		},
		"b": {
			Type: schema.Number,
			Desc: "第二个操作数（sqrt 操作时不需要）",
		},
	}
}

// Execute runs the tool
func (t *CalculatorTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	operation, ok := args["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation is required")
	}

	a, err := toFloat64(args["a"])
	if err != nil {
		return nil, fmt.Errorf("invalid value for 'a': %w", err)
	}

	var result float64
	var expression string

	switch operation {
	case "add", "+":
		b, err := toFloat64(args["b"])
		if err != nil {
			return nil, fmt.Errorf("invalid value for 'b': %w", err)
		}
		result = a + b
		expression = fmt.Sprintf("%v + %v", a, b)

	case "subtract", "-":
		b, err := toFloat64(args["b"])
		if err != nil {
			return nil, fmt.Errorf("invalid value for 'b': %w", err)
		}
		result = a - b
		expression = fmt.Sprintf("%v - %v", a, b)

	case "multiply", "*":
		b, err := toFloat64(args["b"])
		if err != nil {
			return nil, fmt.Errorf("invalid value for 'b': %w", err)
		}
		result = a * b
		expression = fmt.Sprintf("%v × %v", a, b)

	case "divide", "/":
		b, err := toFloat64(args["b"])
		if err != nil {
			return nil, fmt.Errorf("invalid value for 'b': %w", err)
		}
		if b == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		result = a / b
		expression = fmt.Sprintf("%v ÷ %v", a, b)

	case "power", "^":
		b, err := toFloat64(args["b"])
		if err != nil {
			return nil, fmt.Errorf("invalid value for 'b': %w", err)
		}
		result = math.Pow(a, b)
		expression = fmt.Sprintf("%v ^ %v", a, b)

	case "sqrt":
		if a < 0 {
			return nil, fmt.Errorf("cannot calculate square root of negative number")
		}
		result = math.Sqrt(a)
		expression = fmt.Sprintf("√%v", a)

	case "percent":
		b, err := toFloat64(args["b"])
		if err != nil {
			return nil, fmt.Errorf("invalid value for 'b': %w", err)
		}
		result = a * b / 100
		expression = fmt.Sprintf("%v%% of %v", b, a)

	default:
		return nil, fmt.Errorf("unknown operation: %s", operation)
	}

	return map[string]interface{}{
		"expression": expression,
		"result":     result,
	}, nil
}

// toFloat64 converts an interface{} to float64
func toFloat64(v interface{}) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	case nil:
		return 0, fmt.Errorf("value is nil")
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}

// Ensure CalculatorTool implements Tool interface
var _ aitool.Tool = (*CalculatorTool)(nil)
