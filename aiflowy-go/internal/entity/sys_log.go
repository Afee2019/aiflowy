package entity

import (
	"time"
)

// SysLog 操作日志
type SysLog struct {
	ID           int64      `json:"id,string"`
	AccountID    *int64     `json:"accountId,string,omitempty"`
	ActionName   string     `json:"actionName,omitempty"`   // 操作名称
	ActionType   string     `json:"actionType,omitempty"`   // 操作类型
	ActionClass  string     `json:"actionClass,omitempty"`  // 操作类
	ActionMethod string     `json:"actionMethod,omitempty"` // 操作方法
	ActionURL    string     `json:"actionUrl,omitempty"`    // 请求URL
	ActionIP     string     `json:"actionIp,omitempty"`     // 用户IP
	ActionParams string     `json:"actionParams,omitempty"` // 请求参数
	ActionBody   string     `json:"actionBody,omitempty"`   // 请求体
	Status       int        `json:"status"`                 // 状态
	Created      *time.Time `json:"created,omitempty"`      // 操作时间

	// 关联
	Account *SysAccount `json:"account,omitempty"` // 操作人
}
