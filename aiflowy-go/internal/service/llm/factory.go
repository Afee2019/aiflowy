package llm

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"

	"github.com/aiflowy/aiflowy-go/internal/entity"
)

// ModelFactory creates LLM model instances based on configuration
type ModelFactory struct{}

// NewModelFactory creates a new ModelFactory
func NewModelFactory() *ModelFactory {
	return &ModelFactory{}
}

// CreateChatModel creates a ChatModel based on model configuration
func (f *ModelFactory) CreateChatModel(ctx context.Context, m *entity.Model) (model.ChatModel, error) {
	if m == nil {
		return nil, fmt.Errorf("model is nil")
	}

	providerType := ""
	if m.ModelProvider != nil {
		providerType = m.ModelProvider.ProviderType
	}

	// Get effective endpoint and API key (model overrides provider)
	endpoint := m.Endpoint
	if endpoint == "" && m.ModelProvider != nil {
		endpoint = m.ModelProvider.Endpoint
		if m.ModelProvider.ChatPath != "" {
			endpoint = endpoint + m.ModelProvider.ChatPath
		}
	}

	apiKey := m.APIKey
	if apiKey == "" && m.ModelProvider != nil {
		apiKey = m.ModelProvider.APIKey
	}

	switch providerType {
	case entity.ProviderTypeOpenAI:
		return f.createOpenAIChatModel(ctx, m.ModelName, endpoint, apiKey)
	case entity.ProviderTypeDeepSeek:
		return f.createDeepSeekChatModel(ctx, m.ModelName, endpoint, apiKey)
	case entity.ProviderTypeOllama:
		return f.createOllamaChatModel(ctx, m.ModelName, endpoint)
	case entity.ProviderTypeGitee:
		return f.createGiteeChatModel(ctx, m.ModelName, endpoint, apiKey)
	case entity.ProviderTypeSiliconFlow:
		return f.createSiliconFlowChatModel(ctx, m.ModelName, endpoint, apiKey)
	default:
		// Try OpenAI-compatible API as fallback
		return f.createOpenAICompatibleChatModel(ctx, m.ModelName, endpoint, apiKey)
	}
}

// createOpenAIChatModel creates an OpenAI ChatModel
func (f *ModelFactory) createOpenAIChatModel(ctx context.Context, modelName, endpoint, apiKey string) (model.ChatModel, error) {
	config := &openai.ChatModelConfig{
		Model:  modelName,
		APIKey: apiKey,
	}
	if endpoint != "" {
		config.BaseURL = endpoint
	}
	return openai.NewChatModel(ctx, config)
}

// createDeepSeekChatModel creates a DeepSeek ChatModel (OpenAI-compatible)
func (f *ModelFactory) createDeepSeekChatModel(ctx context.Context, modelName, endpoint, apiKey string) (model.ChatModel, error) {
	baseURL := endpoint
	if baseURL == "" {
		baseURL = "https://api.deepseek.com/v1"
	}
	config := &openai.ChatModelConfig{
		Model:   modelName,
		APIKey:  apiKey,
		BaseURL: baseURL,
	}
	return openai.NewChatModel(ctx, config)
}

// createOllamaChatModel creates an Ollama ChatModel
func (f *ModelFactory) createOllamaChatModel(ctx context.Context, modelName, endpoint string) (model.ChatModel, error) {
	baseURL := endpoint
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	config := &ollama.ChatModelConfig{
		BaseURL: baseURL,
		Model:   modelName,
		Timeout: 120 * time.Second,
	}
	return ollama.NewChatModel(ctx, config)
}

// createGiteeChatModel creates a Gitee AI ChatModel (OpenAI-compatible)
func (f *ModelFactory) createGiteeChatModel(ctx context.Context, modelName, endpoint, apiKey string) (model.ChatModel, error) {
	baseURL := endpoint
	if baseURL == "" {
		baseURL = "https://ai.gitee.com/v1"
	}
	config := &openai.ChatModelConfig{
		Model:   modelName,
		APIKey:  apiKey,
		BaseURL: baseURL,
	}
	return openai.NewChatModel(ctx, config)
}

// createSiliconFlowChatModel creates a SiliconFlow ChatModel (OpenAI-compatible)
func (f *ModelFactory) createSiliconFlowChatModel(ctx context.Context, modelName, endpoint, apiKey string) (model.ChatModel, error) {
	baseURL := endpoint
	if baseURL == "" {
		baseURL = "https://api.siliconflow.cn/v1"
	}
	config := &openai.ChatModelConfig{
		Model:   modelName,
		APIKey:  apiKey,
		BaseURL: baseURL,
	}
	return openai.NewChatModel(ctx, config)
}

// createOpenAICompatibleChatModel creates an OpenAI-compatible ChatModel for other providers
func (f *ModelFactory) createOpenAICompatibleChatModel(ctx context.Context, modelName, endpoint, apiKey string) (model.ChatModel, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint is required for OpenAI-compatible providers")
	}
	config := &openai.ChatModelConfig{
		Model:   modelName,
		APIKey:  apiKey,
		BaseURL: endpoint,
	}
	return openai.NewChatModel(ctx, config)
}
