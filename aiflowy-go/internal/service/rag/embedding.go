package rag

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/embedding"
	embeddingOpenAI "github.com/cloudwego/eino-ext/components/embedding/openai"

	"github.com/aiflowy/aiflowy-go/internal/entity"
)

// EmbeddingService Embedding 服务
type EmbeddingService struct{}

// NewEmbeddingService 创建 EmbeddingService
func NewEmbeddingService() *EmbeddingService {
	return &EmbeddingService{}
}

// CreateEmbedder 根据模型配置创建 Embedder
func (s *EmbeddingService) CreateEmbedder(ctx context.Context, m *entity.Model) (embedding.Embedder, error) {
	if m == nil {
		return nil, fmt.Errorf("model is nil")
	}

	providerType := ""
	if m.ModelProvider != nil {
		providerType = m.ModelProvider.ProviderType
	}

	// Get effective endpoint and API key
	endpoint := m.Endpoint
	if endpoint == "" && m.ModelProvider != nil {
		endpoint = m.ModelProvider.Endpoint
		if m.ModelProvider.EmbedPath != "" {
			endpoint = endpoint + m.ModelProvider.EmbedPath
		}
	}

	apiKey := m.APIKey
	if apiKey == "" && m.ModelProvider != nil {
		apiKey = m.ModelProvider.APIKey
	}

	switch providerType {
	case entity.ProviderTypeOpenAI:
		return s.createOpenAIEmbedder(ctx, m.ModelName, endpoint, apiKey)
	case entity.ProviderTypeDeepSeek:
		return s.createDeepSeekEmbedder(ctx, m.ModelName, endpoint, apiKey)
	case entity.ProviderTypeSiliconFlow:
		return s.createSiliconFlowEmbedder(ctx, m.ModelName, endpoint, apiKey)
	case entity.ProviderTypeGitee:
		return s.createGiteeEmbedder(ctx, m.ModelName, endpoint, apiKey)
	case entity.ProviderTypeOllama:
		return s.createOllamaEmbedder(ctx, m.ModelName, endpoint)
	default:
		// Try OpenAI-compatible API as fallback
		return s.createOpenAICompatibleEmbedder(ctx, m.ModelName, endpoint, apiKey)
	}
}

// createOpenAIEmbedder 创建 OpenAI Embedder
func (s *EmbeddingService) createOpenAIEmbedder(ctx context.Context, modelName, endpoint, apiKey string) (embedding.Embedder, error) {
	config := &embeddingOpenAI.EmbeddingConfig{
		Model:  modelName,
		APIKey: apiKey,
	}
	if endpoint != "" {
		config.BaseURL = endpoint
	}
	return embeddingOpenAI.NewEmbedder(ctx, config)
}

// createDeepSeekEmbedder 创建 DeepSeek Embedder (不支持 embedding，返回错误)
func (s *EmbeddingService) createDeepSeekEmbedder(ctx context.Context, modelName, endpoint, apiKey string) (embedding.Embedder, error) {
	// DeepSeek 目前不支持 embedding API，使用 OpenAI-compatible 格式尝试
	baseURL := endpoint
	if baseURL == "" {
		baseURL = "https://api.deepseek.com/v1"
	}
	config := &embeddingOpenAI.EmbeddingConfig{
		Model:   modelName,
		APIKey:  apiKey,
		BaseURL: baseURL,
	}
	return embeddingOpenAI.NewEmbedder(ctx, config)
}

// createSiliconFlowEmbedder 创建 SiliconFlow Embedder
func (s *EmbeddingService) createSiliconFlowEmbedder(ctx context.Context, modelName, endpoint, apiKey string) (embedding.Embedder, error) {
	baseURL := endpoint
	if baseURL == "" {
		baseURL = "https://api.siliconflow.cn/v1"
	}
	config := &embeddingOpenAI.EmbeddingConfig{
		Model:   modelName,
		APIKey:  apiKey,
		BaseURL: baseURL,
	}
	return embeddingOpenAI.NewEmbedder(ctx, config)
}

// createGiteeEmbedder 创建 Gitee AI Embedder
func (s *EmbeddingService) createGiteeEmbedder(ctx context.Context, modelName, endpoint, apiKey string) (embedding.Embedder, error) {
	baseURL := endpoint
	if baseURL == "" {
		baseURL = "https://ai.gitee.com/v1"
	}
	config := &embeddingOpenAI.EmbeddingConfig{
		Model:   modelName,
		APIKey:  apiKey,
		BaseURL: baseURL,
	}
	return embeddingOpenAI.NewEmbedder(ctx, config)
}

// createOllamaEmbedder 创建 Ollama Embedder (使用 OpenAI 兼容 API)
func (s *EmbeddingService) createOllamaEmbedder(ctx context.Context, modelName, endpoint string) (embedding.Embedder, error) {
	baseURL := endpoint
	if baseURL == "" {
		baseURL = "http://localhost:11434/v1"
	}
	config := &embeddingOpenAI.EmbeddingConfig{
		Model:   modelName,
		BaseURL: baseURL,
	}
	return embeddingOpenAI.NewEmbedder(ctx, config)
}

// createOpenAICompatibleEmbedder 创建 OpenAI 兼容 Embedder
func (s *EmbeddingService) createOpenAICompatibleEmbedder(ctx context.Context, modelName, endpoint, apiKey string) (embedding.Embedder, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint is required for OpenAI-compatible providers")
	}
	config := &embeddingOpenAI.EmbeddingConfig{
		Model:   modelName,
		APIKey:  apiKey,
		BaseURL: endpoint,
	}
	return embeddingOpenAI.NewEmbedder(ctx, config)
}

// EmbedText 向量化单个文本
func (s *EmbeddingService) EmbedText(ctx context.Context, embedder embedding.Embedder, text string) ([]float64, error) {
	vectors, err := embedder.EmbedStrings(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	if len(vectors) == 0 {
		return nil, fmt.Errorf("empty embedding result")
	}
	return vectors[0], nil
}

// EmbedTexts 向量化多个文本
func (s *EmbeddingService) EmbedTexts(ctx context.Context, embedder embedding.Embedder, texts []string) ([][]float64, error) {
	return embedder.EmbedStrings(ctx, texts)
}
