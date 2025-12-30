package entity

// BotPlugin Bot-插件关联实体
type BotPlugin struct {
	ID           int64  `db:"id" json:"id,string"`
	BotID        int64  `db:"bot_id" json:"botId,string"`
	PluginItemID int64  `db:"plugin_item_id" json:"pluginItemId,string"`
	Options      string `db:"options" json:"options"` // JSON 格式的配置选项
}
