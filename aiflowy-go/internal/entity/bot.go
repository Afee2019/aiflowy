package entity

import "time"

// Bot represents a bot entity
type Bot struct {
	ID           int64     `json:"id,string"`
	Alias        string    `json:"alias"`
	DeptID       int64     `json:"deptId,string"`
	TenantID     int64     `json:"tenantId,string"`
	CategoryID   int64     `json:"categoryId,string"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Icon         string    `json:"icon"`
	ModelID      int64     `json:"modelId,string"`
	ModelOptions string    `json:"modelOptions"` // JSON string
	Status       int       `json:"status"`
	Options      string    `json:"options"` // JSON string
	Created      time.Time `json:"created"`
	CreatedBy    int64     `json:"createdBy,string"`
	Modified     time.Time `json:"modified"`
	ModifiedBy   int64     `json:"modifiedBy,string"`

	// Relations (populated by join queries)
	Category *BotCategory `json:"category,omitempty"`
	Model    *Model       `json:"model,omitempty"`
}

// BotModelOptions represents the model configuration for a bot
type BotModelOptions struct {
	SystemPrompt  string  `json:"systemPrompt,omitempty"`
	Temperature   float64 `json:"temperature,omitempty"`
	TopP          float64 `json:"topP,omitempty"`
	TopK          int     `json:"topK,omitempty"`
	MaxTokens     int     `json:"maxTokens,omitempty"`
	PresencePenalty  float64 `json:"presencePenalty,omitempty"`
	FrequencyPenalty float64 `json:"frequencyPenalty,omitempty"`
	EnableThinking   bool    `json:"enableThinking,omitempty"`
	ThinkingBudget   int     `json:"thinkingBudget,omitempty"`
}

// BotOptions represents general options for a bot
type BotOptions struct {
	EnableHistory    bool   `json:"enableHistory,omitempty"`
	HistoryCount     int    `json:"historyCount,omitempty"`
	WelcomeMessage   string `json:"welcomeMessage,omitempty"`
	SuggestedQuestions []string `json:"suggestedQuestions,omitempty"`
}
