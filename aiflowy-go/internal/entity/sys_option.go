package entity

// SysOption 系统配置
type SysOption struct {
	TenantID int64  `json:"tenantId,string"`
	Key      string `json:"key"`
	Value    string `json:"value,omitempty"`
}
