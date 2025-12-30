package builtin

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/cloudwego/eino/schema"

	aitool "github.com/aiflowy/aiflowy-go/internal/service/tool"
)

// RandomTool generates random numbers and strings
type RandomTool struct{}

// NewRandomTool creates a new RandomTool
func NewRandomTool() *RandomTool {
	return &RandomTool{}
}

// Name returns the tool name
func (t *RandomTool) Name() string {
	return "random"
}

// Description returns what the tool does
func (t *RandomTool) Description() string {
	return "生成随机数或随机字符串。可以生成指定范围内的随机整数、随机小数、随机字符串或UUID。"
}

// Parameters returns the tool parameters schema
func (t *RandomTool) Parameters() map[string]*schema.ParameterInfo {
	return map[string]*schema.ParameterInfo{
		"type": {
			Type:     schema.String,
			Desc:     "生成类型: integer(整数)、float(小数)、string(字符串)、uuid(UUID)、dice(骰子)、coin(抛硬币)",
			Required: true,
		},
		"min": {
			Type: schema.Integer,
			Desc: "最小值（用于 integer 和 float 类型），默认为 0",
		},
		"max": {
			Type: schema.Integer,
			Desc: "最大值（用于 integer 和 float 类型），默认为 100",
		},
		"length": {
			Type: schema.Integer,
			Desc: "字符串长度（用于 string 类型），默认为 8",
		},
		"sides": {
			Type: schema.Integer,
			Desc: "骰子面数（用于 dice 类型），默认为 6",
		},
	}
}

// Execute runs the tool
func (t *RandomTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	genType, ok := args["type"].(string)
	if !ok {
		return nil, fmt.Errorf("type is required")
	}

	switch genType {
	case "integer":
		return t.randomInteger(args)
	case "float":
		return t.randomFloat(args)
	case "string":
		return t.randomString(args)
	case "uuid":
		return t.randomUUID()
	case "dice":
		return t.rollDice(args)
	case "coin":
		return t.flipCoin()
	default:
		return nil, fmt.Errorf("unknown type: %s", genType)
	}
}

func (t *RandomTool) randomInteger(args map[string]interface{}) (interface{}, error) {
	minVal := int64(0)
	maxVal := int64(100)

	if v, ok := args["min"]; ok {
		if f, ok := v.(float64); ok {
			minVal = int64(f)
		}
	}
	if v, ok := args["max"]; ok {
		if f, ok := v.(float64); ok {
			maxVal = int64(f)
		}
	}

	if minVal >= maxVal {
		return nil, fmt.Errorf("min must be less than max")
	}

	rangeSize := big.NewInt(maxVal - minVal + 1)
	n, err := rand.Int(rand.Reader, rangeSize)
	if err != nil {
		return nil, err
	}

	result := n.Int64() + minVal
	return map[string]interface{}{
		"type":   "integer",
		"result": result,
		"range":  fmt.Sprintf("[%d, %d]", minVal, maxVal),
	}, nil
}

func (t *RandomTool) randomFloat(args map[string]interface{}) (interface{}, error) {
	minVal := float64(0)
	maxVal := float64(100)

	if v, ok := args["min"]; ok {
		if f, ok := v.(float64); ok {
			minVal = f
		}
	}
	if v, ok := args["max"]; ok {
		if f, ok := v.(float64); ok {
			maxVal = f
		}
	}

	if minVal >= maxVal {
		return nil, fmt.Errorf("min must be less than max")
	}

	// Generate a random float between 0 and 1
	precision := int64(1e9)
	n, err := rand.Int(rand.Reader, big.NewInt(precision))
	if err != nil {
		return nil, err
	}

	ratio := float64(n.Int64()) / float64(precision)
	result := minVal + (maxVal-minVal)*ratio

	return map[string]interface{}{
		"type":   "float",
		"result": result,
		"range":  fmt.Sprintf("[%.2f, %.2f]", minVal, maxVal),
	}, nil
}

func (t *RandomTool) randomString(args map[string]interface{}) (interface{}, error) {
	length := 8
	if v, ok := args["length"]; ok {
		if f, ok := v.(float64); ok {
			length = int(f)
		}
	}

	if length <= 0 || length > 256 {
		return nil, fmt.Errorf("length must be between 1 and 256")
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return nil, err
		}
		sb.WriteByte(charset[n.Int64()])
	}

	return map[string]interface{}{
		"type":   "string",
		"result": sb.String(),
		"length": length,
	}, nil
}

func (t *RandomTool) randomUUID() (interface{}, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	// Set version (4) and variant (RFC4122)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])

	return map[string]interface{}{
		"type":   "uuid",
		"result": uuid,
	}, nil
}

func (t *RandomTool) rollDice(args map[string]interface{}) (interface{}, error) {
	sides := int64(6)
	if v, ok := args["sides"]; ok {
		if f, ok := v.(float64); ok {
			sides = int64(f)
		}
	}

	if sides < 2 || sides > 100 {
		return nil, fmt.Errorf("sides must be between 2 and 100")
	}

	n, err := rand.Int(rand.Reader, big.NewInt(sides))
	if err != nil {
		return nil, err
	}

	result := n.Int64() + 1

	return map[string]interface{}{
		"type":   "dice",
		"sides":  sides,
		"result": result,
	}, nil
}

func (t *RandomTool) flipCoin() (interface{}, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(2))
	if err != nil {
		return nil, err
	}

	result := "正面"
	if n.Int64() == 0 {
		result = "反面"
	}

	return map[string]interface{}{
		"type":   "coin",
		"result": result,
	}, nil
}

// Ensure RandomTool implements Tool interface
var _ aitool.Tool = (*RandomTool)(nil)
