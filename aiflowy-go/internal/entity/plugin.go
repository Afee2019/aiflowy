package entity

import "time"

// Plugin 插件实体
type Plugin struct {
	ID          int64     `db:"id" json:"id,string"`
	Alias       string    `db:"alias" json:"alias"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Type        int       `db:"type" json:"type"`
	BaseURL     string    `db:"base_url" json:"baseUrl"`
	AuthType    string    `db:"auth_type" json:"authType"` // apiKey/none
	Created     time.Time `db:"created" json:"created"`
	Icon        string    `db:"icon" json:"icon"`
	Position    string    `db:"position" json:"position"` // headers/query
	Headers     string    `db:"headers" json:"headers"`   // JSON 格式的请求头
	TokenKey    string    `db:"token_key" json:"tokenKey"`
	TokenValue  string    `db:"token_value" json:"tokenValue"`
	DeptID      int64     `db:"dept_id" json:"deptId,string"`
	TenantID    int64     `db:"tenant_id" json:"tenantId,string"`
	CreatedBy   int64     `db:"created_by" json:"createdBy,string"`

	// 关联数据 (非数据库字段)
	Tools []*PluginItem `db:"-" json:"tools,omitempty"`
}

// Title 返回插件标题 (兼容前端)
func (p *Plugin) Title() string {
	return p.Name
}
