package entity

import "time"

// SysDept represents the system department entity
type SysDept struct {
	ID         int64     `json:"id" db:"id"`
	TenantID   int64     `json:"tenantId" db:"tenant_id"`
	ParentID   int64     `json:"parentId" db:"parent_id"`
	Ancestors  string    `json:"ancestors" db:"ancestors"`
	DeptName   string    `json:"deptName" db:"dept_name"`
	DeptCode   string    `json:"deptCode" db:"dept_code"`
	SortNo     int       `json:"sortNo" db:"sort_no"`
	Status     int       `json:"status" db:"status"`
	Created    time.Time `json:"created" db:"created"`
	CreatedBy  int64     `json:"createdBy" db:"created_by"`
	Modified   time.Time `json:"modified" db:"modified"`
	ModifiedBy int64     `json:"modifiedBy" db:"modified_by"`
	Remark     string    `json:"remark" db:"remark"`

	// Non-database fields for tree structure
	Children []*SysDept `json:"children,omitempty" db:"-"`
}

// IsRoot returns true if this is the root department
func (d *SysDept) IsRoot() bool {
	return d.DeptCode == "root" || d.ParentID == 0
}
