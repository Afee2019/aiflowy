package entity

import "time"

// WorkflowExecStatus 工作流执行状态
type WorkflowExecStatus int

const (
	ExecStatusPending   WorkflowExecStatus = 0 // 待执行
	ExecStatusRunning   WorkflowExecStatus = 1 // 执行中
	ExecStatusCompleted WorkflowExecStatus = 2 // 已完成
	ExecStatusFailed    WorkflowExecStatus = 3 // 执行失败
	ExecStatusSuspended WorkflowExecStatus = 4 // 已暂停 (等待人工确认)
)

// WorkflowExecResult 工作流执行记录实体
type WorkflowExecResult struct {
	ID           int64              `db:"id" json:"id,string"`
	ExecKey      string             `db:"exec_key" json:"execKey"`           // 执行标识 (UUID)
	WorkflowID   int64              `db:"workflow_id" json:"workflowId,string"`
	Title        string             `db:"title" json:"title,omitempty"`
	Description  string             `db:"description" json:"description,omitempty"`
	Input        string             `db:"input" json:"input,omitempty"`         // JSON 格式输入
	Output       string             `db:"output" json:"output,omitempty"`       // JSON 格式输出
	WorkflowJSON string             `db:"workflow_json" json:"workflowJson,omitempty"` // 执行时的工作流配置快照
	StartTime    time.Time          `db:"start_time" json:"startTime"`
	EndTime      *time.Time         `db:"end_time" json:"endTime,omitempty"`
	Tokens       int64              `db:"tokens" json:"tokens,omitempty"`
	Status       WorkflowExecStatus `db:"status" json:"status"`
	CreatedKey   string             `db:"created_key" json:"createdKey,omitempty"` // 执行人标识
	CreatedBy    string             `db:"created_by" json:"createdBy,omitempty"`   // 执行人
	ErrorInfo    string             `db:"error_info" json:"errorInfo,omitempty"`
}

// WorkflowExecStep 工作流执行步骤实体
type WorkflowExecStep struct {
	ID        int64              `db:"id" json:"id,string"`
	RecordID  int64              `db:"record_id" json:"recordId,string"` // 关联的执行记录 ID
	ExecKey   string             `db:"exec_key" json:"execKey"`          // 执行标识
	NodeID    string             `db:"node_id" json:"nodeId"`
	NodeName  string             `db:"node_name" json:"nodeName"`
	Input     string             `db:"input" json:"input,omitempty"`     // JSON 格式输入
	Output    string             `db:"output" json:"output,omitempty"`   // JSON 格式输出
	NodeData  string             `db:"node_data" json:"nodeData,omitempty"` // 节点配置 JSON
	StartTime time.Time          `db:"start_time" json:"startTime"`
	EndTime   *time.Time         `db:"end_time" json:"endTime,omitempty"`
	Tokens    int64              `db:"tokens" json:"tokens,omitempty"`
	Status    WorkflowExecStatus `db:"status" json:"status"`
	ErrorInfo string             `db:"error_info" json:"errorInfo,omitempty"`
}
