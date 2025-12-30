package entity

import "time"

// BotConversation represents a conversation with a bot
type BotConversation struct {
	ID         int64     `json:"id,string"`
	Title      string    `json:"title"`
	BotID      int64     `json:"botId,string"`
	AccountID  int64     `json:"accountId,string"`
	Created    time.Time `json:"created"`
	CreatedBy  int64     `json:"createdBy,string"`
	Modified   time.Time `json:"modified"`
	ModifiedBy int64     `json:"modifiedBy,string"`

	// Relations
	Bot      *Bot          `json:"bot,omitempty"`
	Messages []BotMessage  `json:"messages,omitempty"`
}
