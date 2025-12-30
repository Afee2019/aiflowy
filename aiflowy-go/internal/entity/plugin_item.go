package entity

import "time"

// PluginItem 插件工具实体
type PluginItem struct {
	ID            int64     `db:"id" json:"id,string"`
	PluginID      int64     `db:"plugin_id" json:"pluginId,string"`
	Name          string    `db:"name" json:"name"`
	Description   string    `db:"description" json:"description"`
	BasePath      string    `db:"base_path" json:"basePath"`
	Created       time.Time `db:"created" json:"created"`
	Status        int       `db:"status" json:"status"`               // 是否启用
	InputData     string    `db:"input_data" json:"inputData"`        // 输入参数 JSON
	OutputData    string    `db:"output_data" json:"outputData"`      // 输出参数 JSON
	RequestMethod string    `db:"request_method" json:"requestMethod"` // GET/POST/PUT/DELETE
	ServiceStatus int       `db:"service_status" json:"serviceStatus"` // 0=下线, 1=上线
	DebugStatus   int       `db:"debug_status" json:"debugStatus"`     // 0=失败, 1=成功
	EnglishName   string    `db:"english_name" json:"englishName"`

	// 关联数据 (非数据库字段)
	JoinBot bool `db:"-" json:"joinBot,omitempty"` // 是否已关联到 Bot
}
