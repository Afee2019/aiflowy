package entity

import "time"

// SysDict represents the system dictionary entity
type SysDict struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Code        string    `json:"code" db:"code"`
	Description string    `json:"description" db:"description"`
	DictType    int       `json:"dictType" db:"dict_type"`
	SortNo      int       `json:"sortNo" db:"sort_no"`
	Status      int       `json:"status" db:"status"`
	Options     string    `json:"options" db:"options"` // JSON string
	Created     time.Time `json:"created" db:"created"`
	Modified    time.Time `json:"modified" db:"modified"`

	// Non-database fields
	Items []*SysDictItem `json:"items,omitempty" db:"-"`
}

// Dict type constants
const (
	DictTypeCustom    = 1 // Custom dictionary
	DictTypeTable     = 2 // Database table dictionary
	DictTypeEnum      = 3 // Enum class dictionary
	DictTypeSystem    = 4 // System dictionary
)

// SysDictItem represents dictionary item entity
type SysDictItem struct {
	ID          int64     `json:"id" db:"id"`
	DictID      int64     `json:"dictId" db:"dict_id"`
	Text        string    `json:"text" db:"text"`
	Value       string    `json:"value" db:"value"`
	Description string    `json:"description" db:"description"`
	SortNo      int       `json:"sortNo" db:"sort_no"`
	CssContent  string    `json:"cssContent" db:"css_content"`
	CssClass    string    `json:"cssClass" db:"css_class"`
	Remark      string    `json:"remark" db:"remark"`
	Status      int       `json:"status" db:"status"`
	Created     time.Time `json:"created" db:"created"`
	Modified    time.Time `json:"modified" db:"modified"`
}
