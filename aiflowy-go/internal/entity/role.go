package entity

import "time"

// SysRole represents the system role entity
type SysRole struct {
	ID               int64     `json:"id" db:"id"`
	TenantID         int64     `json:"tenantId" db:"tenant_id"`
	RoleName         string    `json:"roleName" db:"role_name"`
	RoleKey          string    `json:"roleKey" db:"role_key"`
	Status           int       `json:"status" db:"status"`
	DataScope        int       `json:"dataScope" db:"data_scope"`
	MenuCheckStrictly bool     `json:"menuCheckStrictly" db:"menu_check_strictly"`
	DeptCheckStrictly bool     `json:"deptCheckStrictly" db:"dept_check_strictly"`
	Created          time.Time `json:"created" db:"created"`
	CreatedBy        int64     `json:"createdBy" db:"created_by"`
	Modified         time.Time `json:"modified" db:"modified"`
	ModifiedBy       int64     `json:"modifiedBy" db:"modified_by"`
	Remark           string    `json:"remark" db:"remark"`

	// Non-database fields for associations
	MenuIds []int64 `json:"menuIds,omitempty" db:"-"`
	DeptIds []int64 `json:"deptIds,omitempty" db:"-"`
}

// Role key constants
const (
	RoleKeySuperAdmin  = "super_admin"
	RoleKeyTenantAdmin = "tenant_admin"
)

// IsSuperAdmin returns true if this is the super admin role
func (r *SysRole) IsSuperAdmin() bool {
	return r.RoleKey == RoleKeySuperAdmin
}

// IsTenantAdmin returns true if this is the tenant admin role
func (r *SysRole) IsTenantAdmin() bool {
	return r.RoleKey == RoleKeyTenantAdmin
}
