package entity

import "time"

// ModelProvider represents the AI model provider entity
type ModelProvider struct {
	ID           int64     `json:"id" db:"id"`
	ProviderName string    `json:"providerName" db:"provider_name"`
	ProviderType string    `json:"providerType" db:"provider_type"`
	Icon         string    `json:"icon" db:"icon"`
	APIKey       string    `json:"apiKey,omitempty" db:"api_key"`
	Endpoint     string    `json:"endpoint" db:"endpoint"`
	ChatPath     string    `json:"chatPath" db:"chat_path"`
	EmbedPath    string    `json:"embedPath" db:"embed_path"`
	RerankPath   string    `json:"rerankPath" db:"rerank_path"`
	Created      time.Time `json:"created" db:"created"`
	CreatedBy    int64     `json:"createdBy" db:"created_by"`
	Modified     time.Time `json:"modified" db:"modified"`
	ModifiedBy   int64     `json:"modifiedBy" db:"modified_by"`
}

// Provider types
const (
	ProviderTypeOpenAI     = "openai"
	ProviderTypeDeepSeek   = "deepseek"
	ProviderTypeOllama     = "ollama"
	ProviderTypeGitee      = "gitee"
	ProviderTypeBaidu      = "baidu"
	ProviderTypeAliyun     = "aliyun"
	ProviderTypeVolcengine = "volcengine"
	ProviderTypeSpark      = "spark"
	ProviderTypeSiliconFlow = "siliconlow"
)
