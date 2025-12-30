package builtin

import (
	"context"
	"time"

	"github.com/cloudwego/eino/schema"

	aitool "github.com/aiflowy/aiflowy-go/internal/service/tool"
)

// TimeTool returns the current date and time
type TimeTool struct{}

// NewTimeTool creates a new TimeTool
func NewTimeTool() *TimeTool {
	return &TimeTool{}
}

// Name returns the tool name
func (t *TimeTool) Name() string {
	return "get_current_time"
}

// Description returns what the tool does
func (t *TimeTool) Description() string {
	return "获取当前的日期和时间。返回格式化的当前时间，包括年月日、时分秒、星期几。"
}

// Parameters returns the tool parameters schema
func (t *TimeTool) Parameters() map[string]*schema.ParameterInfo {
	return map[string]*schema.ParameterInfo{
		"timezone": {
			Type: schema.String,
			Desc: "时区名称，如 'Asia/Shanghai'、'UTC'、'America/New_York'。如果不指定，使用服务器本地时区。",
		},
		"format": {
			Type: schema.String,
			Desc: "时间格式，可选 'full'（完整格式）、'date'（仅日期）、'time'（仅时间）。默认为 'full'。",
		},
	}
}

// Execute runs the tool
func (t *TimeTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Get timezone
	loc := time.Local
	if tz, ok := args["timezone"].(string); ok && tz != "" {
		var err error
		loc, err = time.LoadLocation(tz)
		if err != nil {
			// Fall back to local time
			loc = time.Local
		}
	}

	now := time.Now().In(loc)

	// Get format
	format := "full"
	if f, ok := args["format"].(string); ok && f != "" {
		format = f
	}

	var result string
	weekday := getChineseWeekday(now.Weekday())

	switch format {
	case "date":
		result = now.Format("2006年01月02日") + " " + weekday
	case "time":
		result = now.Format("15:04:05")
	default: // full
		result = now.Format("2006年01月02日 15:04:05") + " " + weekday
	}

	return map[string]interface{}{
		"datetime":  result,
		"timestamp": now.Unix(),
		"timezone":  loc.String(),
	}, nil
}

// getChineseWeekday converts weekday to Chinese
func getChineseWeekday(w time.Weekday) string {
	days := []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}
	return days[w]
}

// Ensure TimeTool implements Tool interface
var _ aitool.Tool = (*TimeTool)(nil)
