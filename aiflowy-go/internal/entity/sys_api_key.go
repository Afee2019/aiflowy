package entity

import (
	"time"
)

// SysApiKey 系统 API 密钥
type SysApiKey struct {
	ID            int64      `json:"id,string"`
	ApiKey        string     `json:"apiKey"`
	Status        int        `json:"status"` // 0: 禁用, 1: 启用
	ExpiredAt     *time.Time `json:"expiredAt,omitempty"`
	Created       *time.Time `json:"created,omitempty"`
	CreatedBy     *int64     `json:"createdBy,string,omitempty"`
	DeptID        *int64     `json:"deptId,string,omitempty"`
	TenantID      *int64     `json:"tenantId,string,omitempty"`
	PermissionIds []int64    `json:"permissionIds,omitempty"` // 关联的资源权限ID列表
}

// SysApiKeyResource API 密钥可访问的资源
type SysApiKeyResource struct {
	ID               int64  `json:"id,string"`
	RequestInterface string `json:"requestInterface"` // 请求接口路径
	Title            string `json:"title"`            // 标题
}

// SysApiKeyResourceMapping API 密钥与资源的映射
type SysApiKeyResourceMapping struct {
	ID               int64 `json:"id,string"`
	ApiKeyID         int64 `json:"apiKeyId,string"`
	ApiKeyResourceID int64 `json:"apiKeyResourceId,string"`
}
