package entity

import (
	"time"
)

// SysJob 系统定时任务
type SysJob struct {
	ID              int64      `json:"id,string"`
	DeptID          int64      `json:"deptId,string"`
	TenantID        int64      `json:"tenantId,string"`
	JobName         string     `json:"jobName"`         // 任务名称
	JobType         int        `json:"jobType"`         // 任务类型
	JobParams       *string    `json:"jobParams,omitempty"` // 任务参数
	CronExpression  string     `json:"cronExpression"`  // cron表达式
	AllowConcurrent int        `json:"allowConcurrent"` // 是否并发执行
	MisfirePolicy   int        `json:"misfirePolicy"`   // 错过策略
	Options         *string    `json:"options,omitempty"` // 其他配置
	Status          int        `json:"status"`          // 状态: 0-停止, 1-运行
	Created         *time.Time `json:"created,omitempty"`
	CreatedBy       *int64     `json:"createdBy,string,omitempty"`
	Modified        *time.Time `json:"modified,omitempty"`
	ModifiedBy      *int64     `json:"modifiedBy,string,omitempty"`
	Remark          *string    `json:"remark,omitempty"` // 备注
}

// SysJobLog 定时任务日志
type SysJobLog struct {
	ID         int64      `json:"id,string"`
	JobID      int64      `json:"jobId,string"`
	JobName    string     `json:"jobName"`
	JobParams  *string    `json:"jobParams,omitempty"`
	Result     *string    `json:"result,omitempty"`   // 执行结果
	Error      *string    `json:"error,omitempty"`    // 错误信息
	Status     int        `json:"status"`             // 状态: 0-失败, 1-成功
	Duration   int64      `json:"duration,omitempty"` // 执行时长(ms)
	Created    *time.Time `json:"created,omitempty"`
}
