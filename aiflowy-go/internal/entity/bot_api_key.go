package entity

import (
	"time"
)

// BotApiKey Bot API 密钥
type BotApiKey struct {
	ID         int64      `json:"id,string"`
	ApiKey     string     `json:"apiKey"`
	BotID      int64      `json:"botId,string"`
	Salt       string     `json:"-"` // 加密用的盐,不返回给前端
	Options    *string    `json:"options,omitempty"`
	Created    *time.Time `json:"created,omitempty"`
	CreatedBy  *int64     `json:"createdBy,string,omitempty"`
	Modified   *time.Time `json:"modified,omitempty"`
	ModifiedBy *int64     `json:"modifiedBy,string,omitempty"`
}
