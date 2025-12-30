package bot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cloudwego/eino/schema"
	"github.com/labstack/echo/v4"

	"github.com/aiflowy/aiflowy-go/internal/dto"
	apierrors "github.com/aiflowy/aiflowy-go/internal/errors"
	"github.com/aiflowy/aiflowy-go/internal/handler/auth"
	"github.com/aiflowy/aiflowy-go/internal/repository"
	"github.com/aiflowy/aiflowy-go/internal/service"
	"github.com/aiflowy/aiflowy-go/internal/service/llm"
	"github.com/aiflowy/aiflowy-go/pkg/protocol"
	"github.com/aiflowy/aiflowy-go/pkg/response"
)

// Handler handles bot-related HTTP requests
type Handler struct {
	botSvc     *service.BotService
	botChatSvc *service.BotChatService
}

// NewHandler creates a new bot handler
func NewHandler() *Handler {
	return &Handler{
		botSvc:     service.NewBotService(),
		botChatSvc: service.NewBotChatService(),
	}
}

// helper to get user context from JWT claims
func getUserContext(c echo.Context) (userID, tenantID, deptID int64) {
	claims := auth.GetClaims(c)
	if claims == nil {
		return 0, 0, 0
	}
	userID, _ = strconv.ParseInt(claims.UserID, 10, 64)
	tenantID, _ = strconv.ParseInt(claims.TenantID, 10, 64)
	deptID, _ = strconv.ParseInt(claims.DeptID, 10, 64)
	return
}

// ========== Bot Endpoints ==========

// BotPage handles GET /api/v1/bot/page
func (h *Handler) BotPage(c echo.Context) error {
	var req dto.BotListRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	bots, total, err := h.botSvc.PageBots(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return response.PageSuccess(c, bots, total, req.GetPage(), req.GetPageSize())
}

// BotList handles GET /api/v1/bot/list
func (h *Handler) BotList(c echo.Context) error {
	var req dto.BotListRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	bots, err := h.botSvc.ListBots(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return response.Success(c, bots)
}

// BotDetail handles GET /api/v1/bot/detail
func (h *Handler) BotDetail(c echo.Context) error {
	id, err := strconv.ParseInt(c.QueryParam("id"), 10, 64)
	if err != nil || id == 0 {
		return apierrors.BadRequest("无效的ID")
	}

	bot, err := h.botSvc.GetBot(c.Request().Context(), id)
	if err != nil {
		return err
	}
	if bot == nil {
		return apierrors.NotFound("机器人不存在")
	}

	return response.Success(c, bot)
}

// GetDetail handles GET /api/v1/bot/getDetail - returns bot with related data
func (h *Handler) GetDetail(c echo.Context) error {
	id, err := strconv.ParseInt(c.QueryParam("id"), 10, 64)
	if err != nil || id == 0 {
		return apierrors.BadRequest("无效的ID")
	}

	detail, err := h.botSvc.GetBotDetail(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return response.Success(c, detail)
}

// BotSave handles POST /api/v1/bot/save
func (h *Handler) BotSave(c echo.Context) error {
	var req dto.BotSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	userID, tenantID, deptID := getUserContext(c)

	bot, err := h.botSvc.SaveBot(c.Request().Context(), &req, userID, tenantID, deptID)
	if err != nil {
		return err
	}

	return response.Success(c, bot)
}

// BotUpdate handles POST /api/v1/bot/update
func (h *Handler) BotUpdate(c echo.Context) error {
	var req dto.BotSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	if req.ID == 0 {
		return apierrors.BadRequest("缺少ID")
	}

	userID, _, _ := getUserContext(c)
	if err := h.botSvc.UpdateBot(c.Request().Context(), &req, userID); err != nil {
		return err
	}

	return response.Success(c, nil)
}

// UpdateLlmOptions handles POST /api/v1/bot/updateLlmOptions
func (h *Handler) UpdateLlmOptions(c echo.Context) error {
	var req dto.BotUpdateLlmOptionsRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	if req.ID == 0 {
		return apierrors.BadRequest("缺少ID")
	}

	userID, _, _ := getUserContext(c)
	if err := h.botSvc.UpdateLlmOptions(c.Request().Context(), &req, userID); err != nil {
		return err
	}

	return response.Success(c, nil)
}

// BotRemove handles POST /api/v1/bot/remove
func (h *Handler) BotRemove(c echo.Context) error {
	var req dto.IDRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	if req.ID == 0 {
		return apierrors.BadRequest("缺少ID")
	}

	if err := h.botSvc.DeleteBot(c.Request().Context(), req.ID); err != nil {
		return err
	}

	return response.Success(c, nil)
}

// GenerateConversationId handles GET /api/v1/bot/generateConversationId
func (h *Handler) GenerateConversationId(c echo.Context) error {
	id := h.botSvc.GenerateConversationID()
	return response.Success(c, map[string]interface{}{
		"conversationId": strconv.FormatInt(id, 10),
	})
}

// ========== Chat Endpoint ==========

// BotChatRequest represents the chat request body
type BotChatRequest struct {
	BotID          int64                   `json:"botId,string"`
	ConversationID int64                   `json:"conversationId,string"`
	Message        string                  `json:"message"`
	Image          string                  `json:"image,omitempty"`
	Stream         bool                    `json:"stream"`
	Options        *service.BotChatOptions `json:"options,omitempty"`
}

// Chat handles POST /api/v1/bot/chat - Bot streaming chat API
func (h *Handler) Chat(c echo.Context) error {
	var req BotChatRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	if req.BotID == 0 {
		return apierrors.BadRequest("缺少机器人ID")
	}
	if req.Message == "" {
		return apierrors.BadRequest("消息内容不能为空")
	}

	userID, _, _ := getUserContext(c)

	chatReq := &service.BotChatRequest{
		BotID:          req.BotID,
		ConversationID: req.ConversationID,
		Message:        req.Message,
		Image:          req.Image,
		Stream:         req.Stream,
		Options:        req.Options,
	}

	// Non-streaming response
	if !req.Stream {
		result, err := h.botChatSvc.Chat(c.Request().Context(), chatReq, userID)
		if err != nil {
			return err
		}
		return response.Success(c, result)
	}

	// Streaming response (SSE)
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("X-Accel-Buffering", "no") // Disable nginx buffering
	c.Response().WriteHeader(200)

	err := h.botChatSvc.ChatStream(c.Request().Context(), chatReq, userID, func(envelope *protocol.Envelope) error {
		sseData, err := envelope.ToSSE()
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(c.Response(), sseData)
		if err != nil {
			return err
		}
		c.Response().Flush()
		return nil
	})

	if err != nil {
		// Send error event if possible
		builder := protocol.NewBuilder("", "")
		errEnv := builder.SystemError("CHAT_ERROR", err.Error(), false)
		sseData, _ := errEnv.ToSSE()
		fmt.Fprint(c.Response(), sseData)
		c.Response().Flush()
	}

	return nil
}

// ========== Category Endpoints ==========

// CategoryList handles GET /api/v1/botCategory/list
func (h *Handler) CategoryList(c echo.Context) error {
	var req dto.BotCategoryListRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	categories, err := h.botSvc.ListCategories(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return response.Success(c, categories)
}

// CategoryDetail handles GET /api/v1/botCategory/detail
func (h *Handler) CategoryDetail(c echo.Context) error {
	id, err := strconv.ParseInt(c.QueryParam("id"), 10, 64)
	if err != nil || id == 0 {
		return apierrors.BadRequest("无效的ID")
	}

	category, err := h.botSvc.GetCategory(c.Request().Context(), id)
	if err != nil {
		return err
	}
	if category == nil {
		return apierrors.NotFound("分类不存在")
	}

	return response.Success(c, category)
}

// CategorySave handles POST /api/v1/botCategory/save
func (h *Handler) CategorySave(c echo.Context) error {
	var req dto.BotCategorySaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	userID, _, _ := getUserContext(c)
	category, err := h.botSvc.SaveCategory(c.Request().Context(), &req, userID)
	if err != nil {
		return err
	}

	return response.Success(c, category)
}

// CategoryUpdate handles POST /api/v1/botCategory/update
func (h *Handler) CategoryUpdate(c echo.Context) error {
	var req dto.BotCategorySaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	if req.ID == 0 {
		return apierrors.BadRequest("缺少ID")
	}

	userID, _, _ := getUserContext(c)
	if err := h.botSvc.UpdateCategory(c.Request().Context(), &req, userID); err != nil {
		return err
	}

	return response.Success(c, nil)
}

// CategoryRemove handles POST /api/v1/botCategory/remove
func (h *Handler) CategoryRemove(c echo.Context) error {
	var req dto.IDRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	if req.ID == 0 {
		return apierrors.BadRequest("缺少ID")
	}

	if err := h.botSvc.DeleteCategory(c.Request().Context(), req.ID); err != nil {
		return err
	}

	return response.Success(c, nil)
}

// ========== Conversation Endpoints ==========

// ConversationPage handles GET /api/v1/botConversation/page
func (h *Handler) ConversationPage(c echo.Context) error {
	var req dto.BotConversationListRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	conversations, total, err := h.botSvc.PageConversations(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return response.PageSuccess(c, conversations, total, req.GetPage(), req.GetPageSize())
}

// ConversationList handles GET /api/v1/botConversation/list
func (h *Handler) ConversationList(c echo.Context) error {
	var req dto.BotConversationListRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	conversations, err := h.botSvc.ListConversations(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return response.Success(c, conversations)
}

// ConversationDetail handles GET /api/v1/botConversation/detail
func (h *Handler) ConversationDetail(c echo.Context) error {
	id, err := strconv.ParseInt(c.QueryParam("id"), 10, 64)
	if err != nil || id == 0 {
		return apierrors.BadRequest("无效的ID")
	}

	conversation, err := h.botSvc.GetConversation(c.Request().Context(), id)
	if err != nil {
		return err
	}
	if conversation == nil {
		return apierrors.NotFound("会话不存在")
	}

	return response.Success(c, conversation)
}

// ConversationSave handles POST /api/v1/botConversation/save
func (h *Handler) ConversationSave(c echo.Context) error {
	var req dto.BotConversationSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	userID, _, _ := getUserContext(c)
	conversation, err := h.botSvc.SaveConversation(c.Request().Context(), &req, userID)
	if err != nil {
		return err
	}

	return response.Success(c, conversation)
}

// ConversationUpdate handles POST /api/v1/botConversation/update
func (h *Handler) ConversationUpdate(c echo.Context) error {
	var req dto.BotConversationSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	if req.ID == 0 {
		return apierrors.BadRequest("缺少ID")
	}

	userID, _, _ := getUserContext(c)
	if err := h.botSvc.UpdateConversation(c.Request().Context(), &req, userID); err != nil {
		return err
	}

	return response.Success(c, nil)
}

// ConversationRemove handles POST /api/v1/botConversation/remove
func (h *Handler) ConversationRemove(c echo.Context) error {
	var req dto.IDRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	if req.ID == 0 {
		return apierrors.BadRequest("缺少ID")
	}

	if err := h.botSvc.DeleteConversation(c.Request().Context(), req.ID); err != nil {
		return err
	}

	return response.Success(c, nil)
}

// ========== Message Endpoints ==========

// MessagePage handles GET /api/v1/botMessage/page
func (h *Handler) MessagePage(c echo.Context) error {
	var req dto.BotMessageListRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	messages, total, err := h.botSvc.PageMessages(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return response.PageSuccess(c, messages, total, req.GetPage(), req.GetPageSize())
}

// MessageList handles GET /api/v1/botMessage/list
func (h *Handler) MessageList(c echo.Context) error {
	var req dto.BotMessageListRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	messages, err := h.botSvc.ListMessages(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return response.Success(c, messages)
}

// MessageDetail handles GET /api/v1/botMessage/detail
func (h *Handler) MessageDetail(c echo.Context) error {
	id, err := strconv.ParseInt(c.QueryParam("id"), 10, 64)
	if err != nil || id == 0 {
		return apierrors.BadRequest("无效的ID")
	}

	message, err := h.botSvc.GetMessage(c.Request().Context(), id)
	if err != nil {
		return err
	}
	if message == nil {
		return apierrors.NotFound("消息不存在")
	}

	return response.Success(c, message)
}

// MessageSave handles POST /api/v1/botMessage/save
func (h *Handler) MessageSave(c echo.Context) error {
	var req dto.BotMessageSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	userID, _, _ := getUserContext(c)
	message, err := h.botSvc.SaveMessage(c.Request().Context(), &req, userID)
	if err != nil {
		return err
	}

	return response.Success(c, message)
}

// MessageUpdate handles POST /api/v1/botMessage/update
func (h *Handler) MessageUpdate(c echo.Context) error {
	var req dto.BotMessageSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	if req.ID == 0 {
		return apierrors.BadRequest("缺少ID")
	}

	userID, _, _ := getUserContext(c)
	if err := h.botSvc.UpdateMessage(c.Request().Context(), &req, userID); err != nil {
		return err
	}

	return response.Success(c, nil)
}

// MessageRemove handles POST /api/v1/botMessage/remove
func (h *Handler) MessageRemove(c echo.Context) error {
	var req dto.IDRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	if req.ID == 0 {
		return apierrors.BadRequest("缺少ID")
	}

	if err := h.botSvc.DeleteMessage(c.Request().Context(), req.ID); err != nil {
		return err
	}

	return response.Success(c, nil)
}

// ========== Voice Input & Prompt Chore ==========

// VoiceInput handles POST /api/v1/bot/voiceInput - 语音识别输入
// 注意: 语音识别需要集成第三方 STT (Speech-to-Text) 服务
// 目前只返回占位响应,实际实现需要集成如 OpenAI Whisper API
func (h *Handler) VoiceInput(c echo.Context) error {
	// 获取上传的音频文件
	file, err := c.FormFile("audio")
	if err != nil {
		return response.BadRequest(c, "请上传音频文件")
	}

	// 检查文件类型
	contentType := file.Header.Get("Content-Type")
	validTypes := map[string]bool{
		"audio/wav":                 true,
		"audio/mpeg":                true,
		"audio/mp3":                 true,
		"audio/webm":                true,
		"audio/ogg":                 true,
		"audio/x-wav":               true,
		"audio/x-m4a":               true,
		"application/octet-stream":  true,
	}
	if !validTypes[contentType] {
		return response.BadRequest(c, "不支持的音频格式")
	}

	// TODO: 集成 STT 服务 (如 OpenAI Whisper API)
	// 目前返回提示信息
	return response.Success(c, "语音识别服务暂未配置,请手动输入")
}

// PromptChoreChat handles POST /api/v1/bot/prompt/chore/chat - 提示词优化
func (h *Handler) PromptChoreChat(c echo.Context) error {
	var req struct {
		Prompt string `json:"prompt"`
		BotID  string `json:"botId"`
	}
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	if req.Prompt == "" {
		return response.BadRequest(c, "提示词不能为空")
	}

	botID, _ := strconv.ParseInt(req.BotID, 10, 64)
	if botID == 0 {
		return response.BadRequest(c, "botId 不能为空")
	}

	ctx := c.Request().Context()

	// 获取 Bot
	bot, err := h.botSvc.GetBot(ctx, botID)
	if err != nil || bot == nil {
		return response.BadRequest(c, "聊天助手不存在")
	}

	// 获取模型
	modelRepo := repository.GetModelRepository()
	model, err := modelRepo.GetModelInstance(ctx, bot.ModelID)
	if err != nil || model == nil {
		return response.BadRequest(c, "模型不存在")
	}

	// 提示词优化的系统提示
	promptChoreSystemPrompt := `# 角色与目标

你是一位专业的提示词工程师（Prompt Engineer）。你的任务是，分析我提供的"用户原始提示词"，并将其优化成一个结构清晰、指令明确、效果最优的"系统提示词（System Prompt）"。

这个优化后的系统提示词将直接用于引导一个AI助手，使其能够精准、高效地完成用户的请求。

# 优化指南 (请严格遵循)

在优化过程中，请遵循以下原则，以确保最终提示词的质量：

1.  **角色定义 (Role Definition)**：
    *   为AI助手明确一个具体、专业的角色。这个角色应该与任务高度相关。
    *   例如："你是一位资深的软件架构师"、"你是一位经验丰富的产品经理"。

2.  **任务与目标 (Task & Goal)**：
    *   清晰、具体地描述AI需要完成的任务。
    *   明确指出期望的最终输出是什么，以及输出的目标和用途。

3.  **上下文与背景 (Context & Background)**：
    *   如果用户的原始提示词中包含背景信息，请保留并整合。
    *   补充必要的上下文，帮助AI更好地理解任务。

4.  **输出格式 (Output Format)**：
    *   明确规定AI回复的格式，如Markdown、列表、代码块等。
    *   如果适用，提供输出的示例。

5.  **约束与限制 (Constraints & Limitations)**：
    *   列出AI应避免的行为或话题。
    *   设定回复的风格、语气或长度。

# 核心输出要求

**只输出优化后的提示词文本本身，不要包含任何解释、说明、标题或额外格式（如"优化后的提示词："这类前缀）。**

---

用户原始提示词：
[` + req.Prompt + `]`

	// 创建 LLM 模型
	factory := llm.NewModelFactory()
	chatModel, err := factory.CreateChatModel(ctx, model)
	if err != nil {
		return response.BadRequest(c, "创建模型失败: "+err.Error())
	}

	// SSE 流式响应
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("X-Accel-Buffering", "no")
	c.Response().WriteHeader(200)

	// 构建消息
	messages := []*schema.Message{
		schema.UserMessage(promptChoreSystemPrompt),
	}

	// 流式生成
	stream, err := chatModel.Stream(ctx, messages)
	if err != nil {
		errData := fmt.Sprintf("event: error\ndata: %s\n\n", err.Error())
		fmt.Fprint(c.Response(), errData)
		c.Response().Flush()
		return nil
	}

	streamContent(ctx, c, stream)
	return nil
}

// streamContent 流式输出内容
func streamContent(ctx context.Context, c echo.Context, stream *schema.StreamReader[*schema.Message]) {
	defer stream.Close()

	for {
		msg, err := stream.Recv()
		if err != nil {
			break
		}

		content := msg.Content
		if content == "" {
			continue
		}

		// 发送 SSE 事件
		sseData := fmt.Sprintf("event: message\ndata: %s\n\n", content)
		fmt.Fprint(c.Response(), sseData)
		c.Response().Flush()
	}

	// 发送结束事件
	fmt.Fprint(c.Response(), "event: done\ndata: [DONE]\n\n")
	c.Response().Flush()
}
