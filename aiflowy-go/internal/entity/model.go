package entity

// Model represents the AI model entity
type Model struct {
	ID                  int64  `json:"id" db:"id"`
	DeptID              int64  `json:"deptId" db:"dept_id"`
	TenantID            int64  `json:"tenantId" db:"tenant_id"`
	ProviderID          int64  `json:"providerId" db:"provider_id"`
	Title               string `json:"title" db:"title"`
	Icon                string `json:"icon" db:"icon"`
	Description         string `json:"description" db:"description"`
	Endpoint            string `json:"endpoint" db:"endpoint"`
	RequestPath         string `json:"requestPath" db:"request_path"`
	ModelName           string `json:"modelName" db:"model_name"`
	APIKey              string `json:"apiKey,omitempty" db:"api_key"`
	ExtraConfig         string `json:"extraConfig" db:"extra_config"`
	Options             string `json:"options" db:"options"`
	GroupName           string `json:"groupName" db:"group_name"`
	ModelType           string `json:"modelType" db:"model_type"`
	WithUsed            bool   `json:"withUsed" db:"with_used"`
	SupportThinking     bool   `json:"supportThinking" db:"support_thinking"`
	SupportTool         bool   `json:"supportTool" db:"support_tool"`
	SupportImage        bool   `json:"supportImage" db:"support_image"`
	SupportImageB64Only bool   `json:"supportImageB64Only" db:"support_image_b64_only"`
	SupportVideo        bool   `json:"supportVideo" db:"support_video"`
	SupportAudio        bool   `json:"supportAudio" db:"support_audio"`
	SupportFree         bool   `json:"supportFree" db:"support_free"`

	// Non-database fields
	ModelProvider *ModelProvider `json:"modelProvider,omitempty" db:"-"`
}

// Model types
const (
	ModelTypeChatModel      = "chatModel"
	ModelTypeEmbeddingModel = "embeddingModel"
	ModelTypeRerankModel    = "rerankModel"
)

// ModelWithProvider represents model with its provider info
type ModelWithProvider struct {
	Model
	ProviderName string `json:"providerName" db:"provider_name"`
	ProviderType string `json:"providerType" db:"provider_type"`
}
