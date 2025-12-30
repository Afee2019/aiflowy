package entity

import "time"

// PluginCategory 插件分类实体
type PluginCategory struct {
	ID        int64     `db:"id" json:"id,string"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// PluginCategoryMapping 插件分类关联实体
type PluginCategoryMapping struct {
	CategoryID int64 `db:"category_id" json:"categoryId,string"`
	PluginID   int64 `db:"plugin_id" json:"pluginId,string"`
}
