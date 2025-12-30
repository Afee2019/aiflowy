package entity

import "time"

// SysMenu represents the system menu entity
type SysMenu struct {
	ID            int64     `json:"id" db:"id"`
	ParentID      int64     `json:"parentId" db:"parent_id"`
	MenuType      int       `json:"menuType" db:"menu_type"`
	MenuTitle     string    `json:"menuTitle" db:"menu_title"`
	MenuUrl       string    `json:"menuUrl" db:"menu_url"`
	Component     string    `json:"component" db:"component"`
	MenuIcon      string    `json:"menuIcon" db:"menu_icon"`
	IsShow        int       `json:"isShow" db:"is_show"`
	PermissionTag string    `json:"permissionTag" db:"permission_tag"`
	SortNo        int       `json:"sortNo" db:"sort_no"`
	Status        int       `json:"status" db:"status"`
	Created       time.Time `json:"created" db:"created"`
	CreatedBy     int64     `json:"createdBy" db:"created_by"`
	Modified      time.Time `json:"modified" db:"modified"`
	ModifiedBy    int64     `json:"modifiedBy" db:"modified_by"`
	Remark        string    `json:"remark" db:"remark"`

	// Non-database fields for tree structure
	Children []*SysMenu `json:"children,omitempty" db:"-"`
}

// Menu type constants
const (
	MenuTypeDirectory = 0 // Directory (folder)
	MenuTypeMenu      = 1 // Menu (page)
	MenuTypeButton    = 2 // Button (action)
)

// MenuVo represents the menu view object for frontend routing
type MenuVo struct {
	ID        int64     `json:"id"`
	ParentID  int64     `json:"parentId"`
	Path      string    `json:"path"`
	Component string    `json:"component"`
	Name      string    `json:"name"`
	Meta      MenuMeta  `json:"meta"`
	Children  []*MenuVo `json:"children,omitempty"`
}

// MenuMeta represents menu metadata
type MenuMeta struct {
	Title string `json:"title"`
	Icon  string `json:"icon"`
	Order int    `json:"order"`
}

// ToMenuVo converts SysMenu to MenuVo for frontend
func (m *SysMenu) ToMenuVo() *MenuVo {
	vo := &MenuVo{
		ID:        m.ID,
		ParentID:  m.ParentID,
		Path:      m.MenuUrl,
		Component: m.Component,
		Name:      m.MenuTitle,
		Meta: MenuMeta{
			Title: m.MenuTitle,
			Icon:  m.MenuIcon,
			Order: m.SortNo,
		},
	}

	if len(m.Children) > 0 {
		vo.Children = make([]*MenuVo, len(m.Children))
		for i, child := range m.Children {
			vo.Children[i] = child.ToMenuVo()
		}
	}

	return vo
}
