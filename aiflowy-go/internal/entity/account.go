package entity

import (
	"time"
)

// SysAccount represents the system user account entity
type SysAccount struct {
	ID          int64     `json:"id" db:"id"`
	DeptID      int64     `json:"deptId" db:"dept_id"`
	TenantID    int64     `json:"tenantId" db:"tenant_id"`
	LoginName   string    `json:"loginName" db:"login_name"`
	Password    string    `json:"-" db:"password"` // Never expose password in JSON
	AccountType int       `json:"accountType" db:"account_type"`
	Nickname    string    `json:"nickname" db:"nickname"`
	Mobile      string    `json:"mobile" db:"mobile"`
	Email       string    `json:"email" db:"email"`
	Avatar      string    `json:"avatar" db:"avatar"`
	Status      int       `json:"status" db:"status"`
	Remark      string    `json:"remark" db:"remark"`
	Created     time.Time `json:"created" db:"created"`
	CreatedBy   int64     `json:"createdBy" db:"created_by"`
	Modified    time.Time `json:"modified" db:"modified"`
	ModifiedBy  int64     `json:"modifiedBy" db:"modified_by"`
}

// AccountStatus constants
const (
	AccountStatusDisabled = 0 // Disabled account
	AccountStatusEnabled  = 1 // Enabled account
)

// AccountType constants
const (
	AccountTypeNormal = 0 // Normal user
	AccountTypeAdmin  = 1 // Admin user
)

// IsEnabled returns true if the account is enabled
func (a *SysAccount) IsEnabled() bool {
	return a.Status == AccountStatusEnabled
}

// LoginAccount represents the logged-in user info stored in context/session
type LoginAccount struct {
	ID          int64  `json:"id"`
	DeptID      int64  `json:"deptId"`
	TenantID    int64  `json:"tenantId"`
	LoginName   string `json:"loginName"`
	AccountType int    `json:"accountType"`
	Nickname    string `json:"nickname"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
	DataScope   int    `json:"dataScope"`
	DeptIDList  string `json:"deptIdList"`
	Remark      string `json:"remark"`
}

// ToLoginAccount converts SysAccount to LoginAccount
func (a *SysAccount) ToLoginAccount() *LoginAccount {
	return &LoginAccount{
		ID:          a.ID,
		DeptID:      a.DeptID,
		TenantID:    a.TenantID,
		LoginName:   a.LoginName,
		AccountType: a.AccountType,
		Nickname:    a.Nickname,
		Mobile:      a.Mobile,
		Email:       a.Email,
		Avatar:      a.Avatar,
		Remark:      a.Remark,
	}
}
