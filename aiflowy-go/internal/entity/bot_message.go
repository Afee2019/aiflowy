package entity

import "time"

// BotMessage represents a message in a bot conversation
type BotMessage struct {
	ID             int64     `json:"id,string"`
	BotID          int64     `json:"botId,string"`
	AccountID      int64     `json:"accountId,string"`
	ConversationID int64     `json:"conversationId,string"`
	Role           string    `json:"role"` // user, assistant, system
	Content        string    `json:"content"`
	Image          string    `json:"image,omitempty"`
	Options        string    `json:"options,omitempty"` // JSON string for additional options
	Created        time.Time `json:"created"`
	Modified       time.Time `json:"modified"`
}

// MessageRole constants
const (
	MessageRoleUser      = "user"
	MessageRoleAssistant = "assistant"
	MessageRoleSystem    = "system"
)

// Role aliases for convenience
const (
	RoleUser      = MessageRoleUser
	RoleAssistant = MessageRoleAssistant
	RoleSystem    = MessageRoleSystem
)

// BotMessageOptions represents additional options for a message
type BotMessageOptions struct {
	TokenUsage    *TokenUsage `json:"tokenUsage,omitempty"`
	ModelName     string      `json:"modelName,omitempty"`
	FinishReason  string      `json:"finishReason,omitempty"`
	ThinkingContent string    `json:"thinkingContent,omitempty"`
}

// TokenUsage represents token usage statistics
type TokenUsage struct {
	PromptTokens     int `json:"promptTokens"`
	CompletionTokens int `json:"completionTokens"`
	TotalTokens      int `json:"totalTokens"`
}
