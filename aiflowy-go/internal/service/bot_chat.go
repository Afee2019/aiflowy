package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"

	"github.com/aiflowy/aiflowy-go/internal/dto"
	"github.com/aiflowy/aiflowy-go/internal/entity"
	apierrors "github.com/aiflowy/aiflowy-go/internal/errors"
	"github.com/aiflowy/aiflowy-go/internal/repository"
	"github.com/aiflowy/aiflowy-go/internal/service/llm"
	aitool "github.com/aiflowy/aiflowy-go/internal/service/tool"
	"github.com/aiflowy/aiflowy-go/pkg/protocol"
	"github.com/aiflowy/aiflowy-go/pkg/snowflake"
)

// BotChatService handles bot chat operations
type BotChatService struct {
	botRepo   *repository.BotRepository
	modelRepo *repository.ModelRepository
	factory   *llm.ModelFactory
}

// NewBotChatService creates a new BotChatService
func NewBotChatService() *BotChatService {
	return &BotChatService{
		botRepo:   repository.GetBotRepository(),
		modelRepo: repository.GetModelRepository(),
		factory:   llm.NewModelFactory(),
	}
}

// BotChatRequest represents a chat request for a bot
type BotChatRequest struct {
	BotID          int64               `json:"botId,string"`
	ConversationID int64               `json:"conversationId,string"`
	Message        string              `json:"message"`
	Image          string              `json:"image,omitempty"`
	Stream         bool                `json:"stream"`
	Options        *BotChatOptions     `json:"options,omitempty"`
}

// BotChatOptions represents chat options that can override bot defaults
type BotChatOptions struct {
	Temperature      *float64 `json:"temperature,omitempty"`
	TopP             *float64 `json:"topP,omitempty"`
	TopK             *int     `json:"topK,omitempty"`
	MaxTokens        *int     `json:"maxTokens,omitempty"`
	EnableThinking   *bool    `json:"enableThinking,omitempty"`
	ThinkingBudget   *int     `json:"thinkingBudget,omitempty"`
	HistoryCount     *int     `json:"historyCount,omitempty"`
}

// BotChatResponse represents a non-streaming chat response
type BotChatResponse struct {
	ConversationID string `json:"conversationId"`
	MessageID      string `json:"messageId"`
	Content        string `json:"content"`
	Thinking       string `json:"thinking,omitempty"`
	Role           string `json:"role"`
}

// StreamCallback is called for each streaming chunk
type StreamCallback func(envelope *protocol.Envelope) error

// ChatContext holds all context needed for a chat session
type ChatContext struct {
	Bot            *entity.Bot
	BotOptions     *entity.BotModelOptions
	Conversation   *entity.BotConversation
	Model          *entity.Model
	Messages       []*entity.BotMessage
	UserMessage    *entity.BotMessage
	AssistantMsgID int64
	Builder        *protocol.Builder
	StartTime      time.Time
	EnableTools    bool                  // Whether tools are enabled for this chat
	ToolInfos      []*schema.ToolInfo    // Tool infos for LLM binding
	ToolNames      []string              // Tool names for execution
}

// Chat performs a bot chat (non-streaming)
func (s *BotChatService) Chat(ctx context.Context, req *BotChatRequest, userID int64) (*BotChatResponse, error) {
	chatCtx, err := s.prepareChat(ctx, req, userID)
	if err != nil {
		return nil, err
	}

	// Build messages for LLM
	llmMessages := s.buildLLMMessages(chatCtx)

	// Create chat model
	baseChatModel, err := s.factory.CreateChatModel(ctx, chatCtx.Model)
	if err != nil {
		return nil, apierrors.InternalError(fmt.Sprintf("创建模型实例失败: %v", err))
	}

	// Bind tools if enabled
	var chatModel model.BaseChatModel = baseChatModel
	if chatCtx.EnableTools && len(chatCtx.ToolInfos) > 0 {
		// Check if model supports tool calling
		if tcm, ok := baseChatModel.(model.ToolCallingChatModel); ok {
			toolModel, err := tcm.WithTools(chatCtx.ToolInfos)
			if err != nil {
				fmt.Printf("Warning: Failed to bind tools: %v\n", err)
			} else {
				chatModel = toolModel
			}
		}
	}

	// Tool call loop - max 5 iterations to prevent infinite loops
	const maxToolIterations = 5
	var finalContent string

	for i := 0; i < maxToolIterations; i++ {
		// Generate response
		result, err := chatModel.Generate(ctx, llmMessages)
		if err != nil {
			return nil, apierrors.InternalError(fmt.Sprintf("生成回复失败: %v", err))
		}

		// Check if LLM wants to call tools
		if len(result.ToolCalls) > 0 {
			// Add assistant message with tool calls to history
			llmMessages = append(llmMessages, result)

			// Execute tools
			for _, tc := range result.ToolCalls {
				toolResult, err := s.executeTool(ctx, tc)
				if err != nil {
					// Add error as tool result
					llmMessages = append(llmMessages, schema.ToolMessage(
						fmt.Sprintf("Error: %v", err),
						tc.ID,
						schema.WithToolName(tc.Function.Name),
					))
				} else {
					llmMessages = append(llmMessages, schema.ToolMessage(
						toolResult,
						tc.ID,
						schema.WithToolName(tc.Function.Name),
					))
				}
			}
			// Continue loop to get final response
			continue
		}

		// No tool calls, we have the final response
		finalContent = result.Content
		break
	}

	// Save assistant message
	assistantMsg := &entity.BotMessage{
		ID:             chatCtx.AssistantMsgID,
		BotID:          req.BotID,
		AccountID:      userID,
		ConversationID: req.ConversationID,
		Role:           entity.RoleAssistant,
		Content:        finalContent,
		Created:        time.Now(),
		Modified:       time.Now(),
	}
	if err := s.botRepo.CreateMessage(ctx, assistantMsg); err != nil {
		// Log but don't fail the response
		fmt.Printf("Failed to save assistant message: %v\n", err)
	}

	return &BotChatResponse{
		ConversationID: strconv.FormatInt(req.ConversationID, 10),
		MessageID:      strconv.FormatInt(chatCtx.AssistantMsgID, 10),
		Content:        finalContent,
		Role:           entity.RoleAssistant,
	}, nil
}

// executeTool executes a tool and returns the result
func (s *BotChatService) executeTool(ctx context.Context, tc schema.ToolCall) (string, error) {
	registry := aitool.GetRegistry()
	return registry.Execute(ctx, tc.Function.Name, tc.Function.Arguments)
}

// ChatStream performs a bot chat with streaming response
func (s *BotChatService) ChatStream(ctx context.Context, req *BotChatRequest, userID int64, callback StreamCallback) error {
	chatCtx, err := s.prepareChat(ctx, req, userID)
	if err != nil {
		return err
	}

	// Send status: running
	if err := callback(chatCtx.Builder.SystemStatus("running")); err != nil {
		return err
	}

	// Build messages for LLM
	llmMessages := s.buildLLMMessages(chatCtx)

	// Create chat model
	baseChatModel, err := s.factory.CreateChatModel(ctx, chatCtx.Model)
	if err != nil {
		errEnv := chatCtx.Builder.SystemError("MODEL_INIT_FAILED", fmt.Sprintf("创建模型实例失败: %v", err), false)
		callback(errEnv)
		return apierrors.InternalError(fmt.Sprintf("创建模型实例失败: %v", err))
	}

	// Bind tools if enabled
	var chatModel model.BaseChatModel = baseChatModel
	if chatCtx.EnableTools && len(chatCtx.ToolInfos) > 0 {
		if tcm, ok := baseChatModel.(model.ToolCallingChatModel); ok {
			toolModel, err := tcm.WithTools(chatCtx.ToolInfos)
			if err != nil {
				fmt.Printf("Warning: Failed to bind tools: %v\n", err)
			} else {
				chatModel = toolModel
			}
		}
	}

	// Tool call loop with streaming
	const maxToolIterations = 5
	var fullContent string
	var fullThinking string

	for iteration := 0; iteration < maxToolIterations; iteration++ {
		// Generate streaming response
		streamReader, err := chatModel.Stream(ctx, llmMessages)
		if err != nil {
			errEnv := chatCtx.Builder.SystemError("STREAM_INIT_FAILED", fmt.Sprintf("开始流式生成失败: %v", err), true)
			callback(errEnv)
			return apierrors.InternalError(fmt.Sprintf("开始流式生成失败: %v", err))
		}

		// Collect full response and tool calls
		var iterationContent string
		var iterationThinking string
		var inThinking bool
		toolCallsMap := make(map[int]*schema.ToolCall) // Use map to merge tool calls by index
		var currentMsg *schema.Message

		// Read stream chunks
		for {
			chunk, err := streamReader.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				streamReader.Close()
				errEnv := chatCtx.Builder.SystemError("STREAM_READ_FAILED", fmt.Sprintf("读取流失败: %v", err), true)
				callback(errEnv)
				return apierrors.InternalError(fmt.Sprintf("读取流失败: %v", err))
			}

			currentMsg = chunk

			// Collect and merge tool calls from chunks (they come in pieces with Index)
			for _, tc := range chunk.ToolCalls {
				index := 0
				if tc.Index != nil {
					index = *tc.Index
				}

				if existing, ok := toolCallsMap[index]; ok {
					// Merge with existing tool call
					if tc.ID != "" {
						existing.ID = tc.ID
					}
					if tc.Type != "" {
						existing.Type = tc.Type
					}
					if tc.Function.Name != "" {
						existing.Function.Name = tc.Function.Name
					}
					existing.Function.Arguments += tc.Function.Arguments
				} else {
					// New tool call
					tcCopy := tc
					toolCallsMap[index] = &tcCopy
				}
			}

			// Handle thinking content
			if chunk.ReasoningContent != "" {
				iterationThinking += chunk.ReasoningContent
				if err := callback(chatCtx.Builder.LLMThinkingDelta(chunk.ReasoningContent)); err != nil {
					streamReader.Close()
					return err
				}
			} else if chunk.Role == "thinking" || (inThinking && chunk.Content != "") {
				inThinking = true
				iterationThinking += chunk.Content
				if err := callback(chatCtx.Builder.LLMThinkingDelta(chunk.Content)); err != nil {
					streamReader.Close()
					return err
				}
			} else {
				if inThinking {
					inThinking = false
				}
				iterationContent += chunk.Content
				if chunk.Content != "" {
					if err := callback(chatCtx.Builder.LLMMessageDelta(chunk.Content)); err != nil {
						streamReader.Close()
						return err
					}
				}
			}
		}
		streamReader.Close()

		// Convert tool calls map to slice, filtering out invalid ones
		var toolCalls []schema.ToolCall
		for _, tc := range toolCallsMap {
			// Skip invalid tool calls (empty name or ID)
			if tc.Function.Name == "" || tc.ID == "" {
				continue
			}
			toolCalls = append(toolCalls, *tc)
		}

		// Check if there are tool calls to execute
		if len(toolCalls) > 0 {
			// Add assistant message with tool calls
			assistantWithTools := &schema.Message{
				Role:      schema.Assistant,
				Content:   iterationContent,
				ToolCalls: toolCalls,
			}
			llmMessages = append(llmMessages, assistantWithTools)

			// Execute each tool and add results
			for _, tc := range toolCalls {
				// Send tool call event
				var args map[string]interface{}
				json.Unmarshal([]byte(tc.Function.Arguments), &args)
				if err := callback(chatCtx.Builder.ToolCall(tc.ID, tc.Function.Name, args)); err != nil {
					return err
				}

				// Execute tool
				toolResult, execErr := s.executeTool(ctx, tc)
				status := "success"
				if execErr != nil {
					status = "error"
					toolResult = execErr.Error()
				}

				// Send tool result event
				if err := callback(chatCtx.Builder.ToolResult(tc.ID, status, toolResult)); err != nil {
					return err
				}

				// Add tool result to messages
				llmMessages = append(llmMessages, schema.ToolMessage(
					toolResult,
					tc.ID,
					schema.WithToolName(tc.Function.Name),
				))
			}

			// Continue loop to get response after tool execution
			continue
		}

		// No tool calls - we have the final response
		fullContent += iterationContent
		fullThinking += iterationThinking

		// If we got content from currentMsg, use it
		if currentMsg != nil && currentMsg.Content != "" && iterationContent == "" {
			fullContent = currentMsg.Content
		}

		break
	}

	// Save assistant message with full content
	assistantMsg := &entity.BotMessage{
		ID:             chatCtx.AssistantMsgID,
		BotID:          req.BotID,
		AccountID:      userID,
		ConversationID: req.ConversationID,
		Role:           entity.RoleAssistant,
		Content:        fullContent,
		Created:        time.Now(),
		Modified:       time.Now(),
	}

	// Store thinking in options if present
	if fullThinking != "" {
		thinkingOpts, _ := json.Marshal(map[string]string{"thinking": fullThinking})
		assistantMsg.Options = string(thinkingOpts)
	}

	if err := s.botRepo.CreateMessage(ctx, assistantMsg); err != nil {
		// Log but don't fail
		fmt.Printf("Failed to save assistant message: %v\n", err)
	}

	// Send done event with metadata
	latency := time.Since(chatCtx.StartTime).Milliseconds()
	meta := &protocol.Meta{
		LatencyMs: latency,
		ModelName: chatCtx.Model.ModelName,
	}
	if err := callback(chatCtx.Builder.SystemDone(meta)); err != nil {
		return err
	}

	return nil
}

// prepareChat prepares the chat context
func (s *BotChatService) prepareChat(ctx context.Context, req *BotChatRequest, userID int64) (*ChatContext, error) {
	startTime := time.Now()

	// Validate request
	if req.BotID == 0 {
		return nil, apierrors.BadRequest("缺少机器人ID")
	}
	if req.Message == "" {
		return nil, apierrors.BadRequest("消息内容不能为空")
	}

	// Get bot
	bot, err := s.botRepo.GetBotByID(ctx, req.BotID)
	if err != nil {
		return nil, apierrors.InternalError("获取机器人失败")
	}
	if bot == nil {
		return nil, apierrors.NotFound("机器人不存在")
	}

	// Parse bot model options
	var botOptions entity.BotModelOptions
	if bot.ModelOptions != "" {
		json.Unmarshal([]byte(bot.ModelOptions), &botOptions)
	}

	// Apply request options overrides
	if req.Options != nil {
		s.applyOptionsOverride(&botOptions, req.Options)
	}

	// Get model
	modelID := bot.ModelID
	if modelID == 0 {
		return nil, apierrors.BadRequest("机器人未配置模型")
	}

	model, err := s.modelRepo.GetModelInstance(ctx, modelID)
	if err != nil {
		return nil, apierrors.InternalError("获取模型失败")
	}
	if model == nil {
		return nil, apierrors.NotFound("模型不存在")
	}

	// Handle conversation
	var conversation *entity.BotConversation
	if req.ConversationID == 0 {
		// Generate new conversation ID
		req.ConversationID = snowflake.MustGenerateID()
	}

	// Check if conversation exists
	conversation, err = s.botRepo.GetConversationByID(ctx, req.ConversationID)
	if err != nil {
		// Log but continue
		fmt.Printf("Failed to get conversation: %v\n", err)
	}

	// Create conversation if not exists
	if conversation == nil {
		now := time.Now()
		conversation = &entity.BotConversation{
			ID:         req.ConversationID,
			Title:      s.generateConversationTitle(req.Message),
			BotID:      req.BotID,
			AccountID:  userID,
			Created:    now,
			CreatedBy:  userID,
			Modified:   now,
			ModifiedBy: userID,
		}
		if err := s.botRepo.CreateConversation(ctx, conversation); err != nil {
			// Log but continue
			fmt.Printf("Failed to create conversation: %v\n", err)
		}
	}

	// Get historical messages for context
	var historyMessages []*entity.BotMessage
	historyCount := 10 // Default
	if botOptions.EnableThinking && botOptions.ThinkingBudget > 0 {
		// Reduce history when using thinking mode to save tokens
		historyCount = 5
	}
	if req.Options != nil && req.Options.HistoryCount != nil {
		historyCount = *req.Options.HistoryCount
	}

	historyMessages, err = s.botRepo.GetRecentMessages(ctx, req.ConversationID, historyCount)
	if err != nil {
		// Log but continue
		fmt.Printf("Failed to get history messages: %v\n", err)
	}

	// Save user message
	userMsgID := snowflake.MustGenerateID()
	userMsg := &entity.BotMessage{
		ID:             userMsgID,
		BotID:          req.BotID,
		AccountID:      userID,
		ConversationID: req.ConversationID,
		Role:           entity.RoleUser,
		Content:        req.Message,
		Image:          req.Image,
		Created:        time.Now(),
		Modified:       time.Now(),
	}
	if err := s.botRepo.CreateMessage(ctx, userMsg); err != nil {
		// Log but continue
		fmt.Printf("Failed to save user message: %v\n", err)
	}

	// Generate assistant message ID
	assistantMsgID := snowflake.MustGenerateID()

	// Create protocol builder
	builder := protocol.NewBuilder(
		strconv.FormatInt(req.ConversationID, 10),
		strconv.FormatInt(assistantMsgID, 10),
	)

	// Load tools - builtin tools + Bot plugins
	var enableTools bool
	var toolInfos []*schema.ToolInfo
	var toolNames []string

	registry := aitool.GetRegistry()

	// 1. Load builtin tools
	builtinTools := registry.GetAll()
	if len(builtinTools) > 0 {
		enableTools = true
		builtinInfos, _ := registry.GetToolInfos(ctx)
		toolInfos = append(toolInfos, builtinInfos...)
		for _, t := range builtinTools {
			toolNames = append(toolNames, t.Name())
		}
	}

	// 2. Load Bot plugin tools
	pluginToolService := NewPluginToolService()
	pluginToolInfos, err := pluginToolService.LoadBotPluginTools(ctx, req.BotID)
	if err == nil && len(pluginToolInfos) > 0 {
		enableTools = true
		toolInfos = append(toolInfos, pluginToolInfos...)
		for _, ti := range pluginToolInfos {
			toolNames = append(toolNames, ti.Name)
		}
	}

	// 3. Load Bot knowledge base tools (RAG)
	knowledgeToolInfos, err := s.loadBotKnowledgeTools(ctx, req.BotID)
	if err == nil && len(knowledgeToolInfos) > 0 {
		enableTools = true
		toolInfos = append(toolInfos, knowledgeToolInfos...)
		for _, ti := range knowledgeToolInfos {
			toolNames = append(toolNames, ti.Name)
		}
	}

	return &ChatContext{
		Bot:            bot,
		BotOptions:     &botOptions,
		Conversation:   conversation,
		Model:          model,
		Messages:       historyMessages,
		UserMessage:    userMsg,
		AssistantMsgID: assistantMsgID,
		Builder:        builder,
		StartTime:      startTime,
		EnableTools:    enableTools,
		ToolInfos:      toolInfos,
		ToolNames:      toolNames,
	}, nil
}

// buildLLMMessages builds the message list for LLM
func (s *BotChatService) buildLLMMessages(chatCtx *ChatContext) []*schema.Message {
	var messages []*schema.Message

	// Add system prompt if present
	if chatCtx.BotOptions.SystemPrompt != "" {
		messages = append(messages, &schema.Message{
			Role:    schema.System,
			Content: chatCtx.BotOptions.SystemPrompt,
		})
	}

	// Add historical messages
	for _, msg := range chatCtx.Messages {
		role := schema.RoleType(msg.Role)
		messages = append(messages, &schema.Message{
			Role:    role,
			Content: msg.Content,
		})
	}

	// Add current user message
	messages = append(messages, &schema.Message{
		Role:    schema.User,
		Content: chatCtx.UserMessage.Content,
	})

	return messages
}

// applyOptionsOverride applies request options to bot options
func (s *BotChatService) applyOptionsOverride(botOpts *entity.BotModelOptions, reqOpts *BotChatOptions) {
	if reqOpts.Temperature != nil {
		botOpts.Temperature = *reqOpts.Temperature
	}
	if reqOpts.TopP != nil {
		botOpts.TopP = *reqOpts.TopP
	}
	if reqOpts.TopK != nil {
		botOpts.TopK = *reqOpts.TopK
	}
	if reqOpts.MaxTokens != nil {
		botOpts.MaxTokens = *reqOpts.MaxTokens
	}
	if reqOpts.EnableThinking != nil {
		botOpts.EnableThinking = *reqOpts.EnableThinking
	}
	if reqOpts.ThinkingBudget != nil {
		botOpts.ThinkingBudget = *reqOpts.ThinkingBudget
	}
}

// generateConversationTitle generates a title from the first message
func (s *BotChatService) generateConversationTitle(message string) string {
	// Take first 50 chars as title
	if len(message) > 50 {
		// Try to find a good break point
		runes := []rune(message)
		if len(runes) > 50 {
			return string(runes[:50]) + "..."
		}
	}
	return message
}

// GetChatDTO converts BotChatRequest from DTO
type BotChatRequestDTO struct {
	BotID          int64           `json:"botId,string"`
	ConversationID int64           `json:"conversationId,string"`
	Message        string          `json:"message"`
	Image          string          `json:"image,omitempty"`
	Stream         bool            `json:"stream"`
	Options        *BotChatOptions `json:"options,omitempty"`
}

// ToBotChatRequest converts DTO to service request
func (d *BotChatRequestDTO) ToBotChatRequest() *BotChatRequest {
	return &BotChatRequest{
		BotID:          d.BotID,
		ConversationID: d.ConversationID,
		Message:        d.Message,
		Image:          d.Image,
		Stream:         d.Stream,
		Options:        d.Options,
	}
}

// BotChatHistoryRequest for getting chat history
type BotChatHistoryRequest struct {
	dto.PageRequest
	ConversationID int64 `json:"conversationId,string" query:"conversationId"`
}

// loadBotKnowledgeTools loads knowledge base tools for a bot
func (s *BotChatService) loadBotKnowledgeTools(ctx context.Context, botID int64) ([]*schema.ToolInfo, error) {
	// Get Bot's associated knowledge bases
	collectionRepo := repository.NewDocumentCollectionRepository()
	collections, err := collectionRepo.ListByBotID(ctx, botID)
	if err != nil {
		return nil, err
	}

	if len(collections) == 0 {
		return nil, nil
	}

	var toolInfos []*schema.ToolInfo

	for _, bdc := range collections {
		if bdc.DocumentCollection == nil {
			continue
		}
		dc := bdc.DocumentCollection

		// Create tool name - must match pattern ^[a-zA-Z0-9_-]+$
		// Prefer EnglishName, fallback to generated name
		toolName := dc.EnglishName
		if toolName == "" {
			// Generate a safe tool name from collection ID
			toolName = fmt.Sprintf("knowledge_%d", dc.ID)
		}

		// Create tool description
		toolDesc := dc.Description
		if toolDesc == "" {
			toolDesc = "搜索 " + dc.Title + " 知识库中的相关信息"
		}

		// Create tool info
		toolInfo := &schema.ToolInfo{
			Name: toolName,
			Desc: toolDesc,
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"input": {
					Type: schema.String,
					Desc: "要在知识库中搜索的关键词或问题",
				},
			}),
		}

		// Register the tool in the registry for execution
		knowledgeTool := &KnowledgeToolWrapper{
			CollectionID: dc.ID,
			ToolName:     toolName,
		}
		aitool.GetRegistry().Register(knowledgeTool)

		toolInfos = append(toolInfos, toolInfo)
	}

	return toolInfos, nil
}

// KnowledgeToolWrapper wraps knowledge tool for registry
type KnowledgeToolWrapper struct {
	CollectionID int64
	ToolName     string
}

// Name returns the tool name
func (t *KnowledgeToolWrapper) Name() string {
	return t.ToolName
}

// Description returns the tool description
func (t *KnowledgeToolWrapper) Description() string {
	return "搜索知识库中的相关信息"
}

// Parameters returns the tool parameters
func (t *KnowledgeToolWrapper) Parameters() map[string]*schema.ParameterInfo {
	return map[string]*schema.ParameterInfo{
		"input": {
			Type: schema.String,
			Desc: "要在知识库中搜索的关键词或问题",
		},
	}
}

// Execute executes the knowledge tool
func (t *KnowledgeToolWrapper) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	input, ok := args["input"].(string)
	if !ok || input == "" {
		return nil, fmt.Errorf("input parameter is required")
	}

	// Call RAG service
	ragService := NewDocumentCollectionService()
	docs := ragService.SearchByCollectionID(ctx, t.CollectionID, input, 5)

	if len(docs) == 0 {
		return "未找到相关信息", nil
	}

	// Build result
	var result string
	for i, doc := range docs {
		result += fmt.Sprintf("[%d] %s\n\n", i+1, doc.Content)
	}

	return result, nil
}
