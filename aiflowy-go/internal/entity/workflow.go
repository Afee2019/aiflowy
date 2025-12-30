package entity

import "time"

// Workflow 工作流实体
type Workflow struct {
	ID          int64     `db:"id" json:"id,string"`
	Alias       string    `db:"alias" json:"alias,omitempty"`
	DeptID      int64     `db:"dept_id" json:"deptId,string"`
	TenantID    int64     `db:"tenant_id" json:"tenantId,string"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description,omitempty"`
	Icon        string    `db:"icon" json:"icon,omitempty"`
	Content     string    `db:"content" json:"content,omitempty"` // 工作流设计的 JSON 内容
	Created     time.Time `db:"created" json:"created"`
	CreatedBy   int64     `db:"created_by" json:"createdBy,string"`
	Modified    time.Time `db:"modified" json:"modified"`
	ModifiedBy  int64     `db:"modified_by" json:"modifiedBy,string"`
	EnglishName string    `db:"english_name" json:"englishName,omitempty"`
	Status      int       `db:"status" json:"status"`
	CategoryID  int64     `db:"category_id" json:"categoryId,string,omitempty"`

	// 非数据库字段
	CategoryName string `db:"-" json:"categoryName,omitempty"`
}

// TableName 返回表名
func (Workflow) TableName() string {
	return "tb_workflow"
}
