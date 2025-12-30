package dto

// WorkflowSaveRequest 保存工作流请求
type WorkflowSaveRequest struct {
	ID          string `json:"id,omitempty"`
	Alias       string `json:"alias,omitempty"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Content     string `json:"content,omitempty"`
	EnglishName string `json:"englishName,omitempty"`
	CategoryID  string `json:"categoryId,omitempty"`
	Status      int    `json:"status"`
}

// WorkflowCategorySaveRequest 保存工作流分类请求
type WorkflowCategorySaveRequest struct {
	ID           string `json:"id,omitempty"`
	CategoryName string `json:"categoryName" validate:"required"`
	SortNo       int    `json:"sortNo"`
	Status       int    `json:"status"`
}

// BotWorkflowUpdateRequest 更新 Bot-工作流关联请求
type BotWorkflowUpdateRequest struct {
	BotID       string   `json:"botId" validate:"required"`
	WorkflowIDs []string `json:"workflowIds"`
}

// WorkflowRunRequest 运行工作流请求
type WorkflowRunRequest struct {
	ID        string                 `json:"id" validate:"required"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// WorkflowSingleRunRequest 单节点运行请求
type WorkflowSingleRunRequest struct {
	WorkflowID string                 `json:"workflowId" validate:"required"`
	NodeID     string                 `json:"nodeId" validate:"required"`
	Variables  map[string]interface{} `json:"variables,omitempty"`
}

// WorkflowResumeRequest 恢复工作流请求
type WorkflowResumeRequest struct {
	ExecuteID     string                 `json:"executeId" validate:"required"`
	ConfirmParams map[string]interface{} `json:"confirmParams,omitempty"`
}

// ChainStatusRequest 获取工作流状态请求
type ChainStatusRequest struct {
	ExecuteID string      `json:"executeId" validate:"required"`
	Nodes     []*NodeInfo `json:"nodes,omitempty"`
}

// ========================== 工作流 DSL 相关 ==========================

// WorkflowDefinition 工作流定义 (DSL)
type WorkflowDefinition struct {
	ID          string            `json:"id,omitempty"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Nodes       []*WorkflowNode   `json:"nodes,omitempty"`
	Edges       []*WorkflowEdge   `json:"edges,omitempty"`
	Variables   map[string]string `json:"variables,omitempty"`
}

// WorkflowNode 工作流节点
type WorkflowNode struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`               // start, end, llm, tool, condition, human_confirm, workflow, code
	Name       string                 `json:"name,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`     // 节点配置数据
	Position   *NodePosition          `json:"position,omitempty"` // 可视化位置
	Parameters []*WorkflowParameter   `json:"parameters,omitempty"`
}

// WorkflowEdge 工作流边 (连接)
type WorkflowEdge struct {
	ID         string `json:"id,omitempty"`
	Source     string `json:"source"`               // 源节点 ID
	Target     string `json:"target"`               // 目标节点 ID
	SourcePort string `json:"sourcePort,omitempty"` // 源端口
	TargetPort string `json:"targetPort,omitempty"` // 目标端口
	Condition  string `json:"condition,omitempty"`  // 条件表达式
}

// NodePosition 节点位置
type NodePosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// WorkflowParameter 工作流参数
type WorkflowParameter struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"` // string, number, boolean, array, object
	Description  string      `json:"description,omitempty"`
	Required     bool        `json:"required,omitempty"`
	DefaultValue interface{} `json:"defaultValue,omitempty"`
}

// ========================== 执行状态相关 ==========================

// ChainInfo 工作流执行信息
type ChainInfo struct {
	ExecuteID string                 `json:"executeId"`
	Status    int                    `json:"status"` // 0: pending, 1: running, 2: completed, 3: failed, 4: suspended
	Message   string                 `json:"message,omitempty"`
	Result    map[string]interface{} `json:"result,omitempty"`
	Nodes     map[string]*NodeInfo   `json:"nodes,omitempty"`
}

// NodeInfo 节点执行信息
type NodeInfo struct {
	NodeID               string                 `json:"nodeId"`
	NodeName             string                 `json:"nodeName,omitempty"`
	Status               int                    `json:"status"` // 0: pending, 1: running, 2: completed, 3: failed, 4: suspended
	Message              string                 `json:"message,omitempty"`
	Result               map[string]interface{} `json:"result,omitempty"`
	SuspendForParameters []*WorkflowParameter   `json:"suspendForParameters,omitempty"`
}

// 工作流执行状态常量
const (
	ChainStatusPending   = 0 // 等待执行
	ChainStatusRunning   = 1 // 执行中
	ChainStatusCompleted = 2 // 已完成
	ChainStatusFailed    = 3 // 执行失败
	ChainStatusSuspended = 4 // 已挂起 (等待人工确认)
)

// 节点类型常量
const (
	NodeTypeStart        = "start"         // 开始节点
	NodeTypeEnd          = "end"           // 结束节点
	NodeTypeLLM          = "llm"           // LLM 节点
	NodeTypeTool         = "tool"          // 工具节点
	NodeTypeCondition    = "condition"     // 条件节点
	NodeTypeHumanConfirm = "human_confirm" // 人工确认节点
	NodeTypeWorkflow     = "workflow"      // 子工作流节点
	NodeTypeCode         = "code"          // 代码节点
	NodeTypePlugin       = "plugin"        // 插件节点
	NodeTypeDoc          = "doc"           // 文档节点
	NodeTypeSQL          = "sql"           // SQL 节点
)

// RunningParametersResponse 运行参数响应
type RunningParametersResponse struct {
	Parameters  []*WorkflowParameter `json:"parameters,omitempty"`
	Title       string               `json:"title,omitempty"`
	Description string               `json:"description,omitempty"`
	Icon        string               `json:"icon,omitempty"`
}
