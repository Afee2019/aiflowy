package dto

// AccountSaveRequest represents the request to create/update an account
type AccountSaveRequest struct {
	ID          int64   `json:"id"`
	DeptID      int64   `json:"deptId"`
	LoginName   string  `json:"loginName"`
	Password    string  `json:"password"`
	AccountType int     `json:"accountType"`
	Nickname    string  `json:"nickname"`
	Mobile      string  `json:"mobile"`
	Email       string  `json:"email"`
	Avatar      string  `json:"avatar"`
	Status      int     `json:"status"`
	Remark      string  `json:"remark"`
	RoleIds     []int64 `json:"roleIds"`
	PositionIds []int64 `json:"positionIds"`
}

// AccountQueryRequest represents account query parameters
type AccountQueryRequest struct {
	PageRequest
	LoginName   string `query:"loginName"`
	Nickname    string `query:"nickname"`
	Mobile      string `query:"mobile"`
	Email       string `query:"email"`
	Status      *int   `query:"status"`
	DeptID      *int64 `query:"deptId"`
	AccountType *int   `query:"accountType"`
}

// UpdatePasswordRequest represents password update request
type UpdatePasswordRequest struct {
	Password        string `json:"password"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

// RoleSaveRequest represents the request to create/update a role
type RoleSaveRequest struct {
	ID                int64   `json:"id"`
	RoleName          string  `json:"roleName"`
	RoleKey           string  `json:"roleKey"`
	Status            int     `json:"status"`
	DataScope         int     `json:"dataScope"`
	MenuCheckStrictly bool    `json:"menuCheckStrictly"`
	DeptCheckStrictly bool    `json:"deptCheckStrictly"`
	Remark            string  `json:"remark"`
	MenuIds           []int64 `json:"menuIds"`
	DeptIds           []int64 `json:"deptIds"`
}

// RoleQueryRequest represents role query parameters
type RoleQueryRequest struct {
	PageRequest
	RoleName string `query:"roleName"`
	RoleKey  string `query:"roleKey"`
	Status   *int   `query:"status"`
}

// MenuSaveRequest represents the request to create/update a menu
type MenuSaveRequest struct {
	ID            int64  `json:"id"`
	ParentID      int64  `json:"parentId"`
	MenuType      int    `json:"menuType"`
	MenuTitle     string `json:"menuTitle"`
	MenuUrl       string `json:"menuUrl"`
	Component     string `json:"component"`
	MenuIcon      string `json:"menuIcon"`
	IsShow        int    `json:"isShow"`
	PermissionTag string `json:"permissionTag"`
	SortNo        int    `json:"sortNo"`
	Status        int    `json:"status"`
	Remark        string `json:"remark"`
}

// MenuQueryRequest represents menu query parameters
type MenuQueryRequest struct {
	PageRequest
	MenuTitle string `query:"menuTitle"`
	Status    *int   `query:"status"`
	MenuType  *int   `query:"menuType"`
	AsTree    bool   `query:"asTree"`
}

// DeptSaveRequest represents the request to create/update a department
type DeptSaveRequest struct {
	ID       int64  `json:"id"`
	ParentID int64  `json:"parentId"`
	DeptName string `json:"deptName"`
	DeptCode string `json:"deptCode"`
	SortNo   int    `json:"sortNo"`
	Status   int    `json:"status"`
	Remark   string `json:"remark"`
}

// DeptQueryRequest represents department query parameters
type DeptQueryRequest struct {
	PageRequest
	DeptName string `query:"deptName"`
	DeptCode string `query:"deptCode"`
	Status   *int   `query:"status"`
	AsTree   bool   `query:"asTree"`
}

// DictSaveRequest represents the request to create/update a dictionary
type DictSaveRequest struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	DictType    int    `json:"dictType"`
	SortNo      int    `json:"sortNo"`
	Status      int    `json:"status"`
	Options     string `json:"options"`
}

// DictQueryRequest represents dictionary query parameters
type DictQueryRequest struct {
	PageRequest
	Name     string `query:"name"`
	Code     string `query:"code"`
	DictType *int   `query:"dictType"`
	Status   *int   `query:"status"`
}

// DictItemSaveRequest represents the request to create/update a dictionary item
type DictItemSaveRequest struct {
	ID          int64  `json:"id"`
	DictID      int64  `json:"dictId"`
	Text        string `json:"text"`
	Value       string `json:"value"`
	Description string `json:"description"`
	SortNo      int    `json:"sortNo"`
	CssContent  string `json:"cssContent"`
	CssClass    string `json:"cssClass"`
	Remark      string `json:"remark"`
	Status      int    `json:"status"`
}

// DictItemQueryRequest represents dictionary item query parameters
type DictItemQueryRequest struct {
	PageRequest
	DictID int64  `query:"dictId"`
	Text   string `query:"text"`
	Value  string `query:"value"`
	Status *int   `query:"status"`
}
