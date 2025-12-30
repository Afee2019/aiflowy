package entity

import "time"

// WorkflowCategory 工作流分类实体
type WorkflowCategory struct {
	ID           int64     `db:"id" json:"id,string"`
	CategoryName string    `db:"category_name" json:"categoryName"`
	SortNo       int       `db:"sort_no" json:"sortNo"`
	Created      time.Time `db:"created" json:"created"`
	CreatedBy    int64     `db:"created_by" json:"createdBy,string"`
	Modified     time.Time `db:"modified" json:"modified"`
	ModifiedBy   int64     `db:"modified_by" json:"modifiedBy,string"`
	Status       int       `db:"status" json:"status"`
}

// TableName 返回表名
func (WorkflowCategory) TableName() string {
	return "tb_workflow_category"
}
