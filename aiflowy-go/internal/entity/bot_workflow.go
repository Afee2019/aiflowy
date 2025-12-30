package entity

// BotWorkflow Bot-工作流关联实体
type BotWorkflow struct {
	ID         int64  `db:"id" json:"id,string"`
	BotID      int64  `db:"bot_id" json:"botId,string"`
	WorkflowID int64  `db:"workflow_id" json:"workflowId,string"`
	Options    string `db:"options" json:"options,omitempty"` // JSON 格式的配置选项

	// 非数据库字段 - 关联的工作流信息
	Workflow *Workflow `db:"-" json:"workflow,omitempty"`
}

// TableName 返回表名
func (BotWorkflow) TableName() string {
	return "tb_bot_workflow"
}
