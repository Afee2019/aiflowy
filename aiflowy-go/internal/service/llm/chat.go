package llm

import (
	"context"
	"fmt"
	"io"

	"github.com/cloudwego/eino/schema"

	apierrors "github.com/aiflowy/aiflowy-go/internal/errors"
	"github.com/aiflowy/aiflowy-go/internal/repository"
)

// ChatService handles LLM chat operations
type ChatService struct {
	factory   *ModelFactory
	modelRepo *repository.ModelRepository
}

// NewChatService creates a new ChatService
func NewChatService() *ChatService {
	return &ChatService{
		factory:   NewModelFactory(),
		modelRepo: repository.GetModelRepository(),
	}
}

// ChatRequest represents a chat request
type ChatRequest struct {
	ModelID  int64    `json:"modelId"`
	Messages []Message `json:"messages"`
	Options  *ChatOptions `json:"options,omitempty"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatOptions represents chat options
type ChatOptions struct {
	Temperature float64 `json:"temperature,omitempty"`
	MaxTokens   int     `json:"maxTokens,omitempty"`
	TopP        float64 `json:"topP,omitempty"`
	TopK        int     `json:"topK,omitempty"`
}

// ChatResponse represents a chat response
type ChatResponse struct {
	Content      string `json:"content"`
	Role         string `json:"role"`
	ModelName    string `json:"modelName"`
	ProviderType string `json:"providerType"`
	FinishReason string `json:"finishReason,omitempty"`
}

// StreamChunk represents a streaming response chunk
type StreamChunk struct {
	Content      string `json:"content"`
	Done         bool   `json:"done"`
	FinishReason string `json:"finishReason,omitempty"`
}

// Chat performs a synchronous chat completion
func (s *ChatService) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// Get model with provider info
	model, err := s.modelRepo.GetModelInstance(ctx, req.ModelID)
	if err != nil {
		return nil, apierrors.InternalError("获取模型失败")
	}
	if model == nil {
		return nil, apierrors.NotFound("模型不存在")
	}

	// Create chat model
	chatModel, err := s.factory.CreateChatModel(ctx, model)
	if err != nil {
		return nil, apierrors.InternalError(fmt.Sprintf("创建模型实例失败: %v", err))
	}

	// Convert messages to Eino schema
	messages := make([]*schema.Message, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = &schema.Message{
			Role:    schema.RoleType(msg.Role),
			Content: msg.Content,
		}
	}

	// Generate response
	result, err := chatModel.Generate(ctx, messages)
	if err != nil {
		return nil, apierrors.InternalError(fmt.Sprintf("生成回复失败: %v", err))
	}

	providerType := ""
	if model.ModelProvider != nil {
		providerType = model.ModelProvider.ProviderType
	}

	return &ChatResponse{
		Content:      result.Content,
		Role:         string(result.Role),
		ModelName:    model.ModelName,
		ProviderType: providerType,
	}, nil
}

// ChatStream performs a streaming chat completion
func (s *ChatService) ChatStream(ctx context.Context, req *ChatRequest, onChunk func(*StreamChunk) error) error {
	// Get model with provider info
	model, err := s.modelRepo.GetModelInstance(ctx, req.ModelID)
	if err != nil {
		return apierrors.InternalError("获取模型失败")
	}
	if model == nil {
		return apierrors.NotFound("模型不存在")
	}

	// Create chat model
	chatModel, err := s.factory.CreateChatModel(ctx, model)
	if err != nil {
		return apierrors.InternalError(fmt.Sprintf("创建模型实例失败: %v", err))
	}

	// Convert messages to Eino schema
	messages := make([]*schema.Message, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = &schema.Message{
			Role:    schema.RoleType(msg.Role),
			Content: msg.Content,
		}
	}

	// Generate streaming response
	streamReader, err := chatModel.Stream(ctx, messages)
	if err != nil {
		return apierrors.InternalError(fmt.Sprintf("开始流式生成失败: %v", err))
	}
	defer streamReader.Close()

	// Read stream chunks
	for {
		chunk, err := streamReader.Recv()
		if err == io.EOF {
			// Send final done chunk
			if err := onChunk(&StreamChunk{Done: true}); err != nil {
				return err
			}
			break
		}
		if err != nil {
			return apierrors.InternalError(fmt.Sprintf("读取流失败: %v", err))
		}

		// Send chunk to callback
		if err := onChunk(&StreamChunk{Content: chunk.Content}); err != nil {
			return err
		}
	}

	return nil
}

// TestModel tests if a model can be called successfully
func (s *ChatService) TestModel(ctx context.Context, modelID int64, prompt string) (*ChatResponse, error) {
	if prompt == "" {
		prompt = "Hello"
	}

	req := &ChatRequest{
		ModelID: modelID,
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	}

	return s.Chat(ctx, req)
}
