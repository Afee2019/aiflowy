package dto

// BotSaveRequest represents request to save a bot
type BotSaveRequest struct {
	ID           int64  `json:"id,string"`
	Alias        string `json:"alias"`
	CategoryID   int64  `json:"categoryId,string"`
	Title        string `json:"title" validate:"required"`
	Description  string `json:"description"`
	Icon         string `json:"icon"`
	ModelID      int64  `json:"modelId,string"`
	ModelOptions string `json:"modelOptions"`
	Status       int    `json:"status"`
	Options      string `json:"options"`
}

// BotListRequest represents request to list bots
type BotListRequest struct {
	PageRequest
	CategoryID int64  `query:"categoryId" json:"categoryId,string"`
	Title      string `query:"title" json:"title"`
	Status     *int   `query:"status" json:"status"`
}

// BotUpdateLlmOptionsRequest represents request to update bot's LLM options
type BotUpdateLlmOptionsRequest struct {
	ID           int64  `json:"id,string" validate:"required"`
	ModelID      int64  `json:"modelId,string"`
	ModelOptions string `json:"modelOptions"`
}

// BotCategorySaveRequest represents request to save bot category
type BotCategorySaveRequest struct {
	ID           int64  `json:"id,string"`
	CategoryName string `json:"categoryName" validate:"required"`
	SortNo       int    `json:"sortNo"`
	Status       int    `json:"status"`
}

// BotCategoryListRequest represents request to list bot categories
type BotCategoryListRequest struct {
	CategoryName string `query:"categoryName" json:"categoryName"`
	Status       *int   `query:"status" json:"status"`
}

// BotConversationSaveRequest represents request to save conversation
type BotConversationSaveRequest struct {
	ID        int64  `json:"id,string"`
	Title     string `json:"title" validate:"required"`
	BotID     int64  `json:"botId,string" validate:"required"`
	AccountID int64  `json:"accountId,string"`
}

// BotConversationListRequest represents request to list conversations
type BotConversationListRequest struct {
	PageRequest
	BotID     int64 `query:"botId" json:"botId,string"`
	AccountID int64 `query:"accountId" json:"accountId,string"`
}

// BotMessageSaveRequest represents request to save message
type BotMessageSaveRequest struct {
	ID             int64  `json:"id,string"`
	BotID          int64  `json:"botId,string" validate:"required"`
	AccountID      int64  `json:"accountId,string"`
	ConversationID int64  `json:"conversationId,string" validate:"required"`
	Role           string `json:"role" validate:"required"`
	Content        string `json:"content"`
	Image          string `json:"image"`
	Options        string `json:"options"`
}

// BotMessageListRequest represents request to list messages
type BotMessageListRequest struct {
	PageRequest
	BotID          int64 `query:"botId" json:"botId,string"`
	ConversationID int64 `query:"conversationId" json:"conversationId,string"`
	AccountID      int64 `query:"accountId" json:"accountId,string"`
}

// BotDetailResponse represents bot detail with related data
type BotDetailResponse struct {
	ID           int64           `json:"id,string"`
	Alias        string          `json:"alias"`
	DeptID       int64           `json:"deptId,string"`
	TenantID     int64           `json:"tenantId,string"`
	CategoryID   int64           `json:"categoryId,string"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Icon         string          `json:"icon"`
	ModelID      int64           `json:"modelId,string"`
	ModelOptions interface{}     `json:"modelOptions"`
	Status       int             `json:"status"`
	Options      interface{}     `json:"options"`
	Category     *CategoryInfo   `json:"category,omitempty"`
	Model        *ModelInfo      `json:"model,omitempty"`
}

// CategoryInfo represents category info in bot detail
type CategoryInfo struct {
	ID           int64  `json:"id,string"`
	CategoryName string `json:"categoryName"`
}

// ModelInfo represents model info in bot detail
type ModelInfo struct {
	ID           int64  `json:"id,string"`
	Title        string `json:"title"`
	ModelName    string `json:"modelName"`
	ProviderType string `json:"providerType"`
}
