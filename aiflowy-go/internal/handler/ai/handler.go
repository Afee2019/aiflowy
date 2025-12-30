package ai

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"

	apierrors "github.com/aiflowy/aiflowy-go/internal/errors"
	"github.com/aiflowy/aiflowy-go/internal/service/llm"
	"github.com/aiflowy/aiflowy-go/pkg/response"
)

// Handler handles AI-related HTTP requests
type Handler struct {
	chatSvc *llm.ChatService
}

// NewHandler creates a new AI handler
func NewHandler() *Handler {
	return &Handler{
		chatSvc: llm.NewChatService(),
	}
}

// TestRequest represents a test API request
type TestRequest struct {
	ModelID int64  `json:"modelId" query:"modelId"`
	Prompt  string `json:"prompt" query:"prompt"`
}

// ChatRequest represents a chat API request
type ChatRequest struct {
	ModelID  int64        `json:"modelId"`
	Messages []llm.Message `json:"messages"`
	Stream   bool         `json:"stream"`
}

// Test tests a model with a simple prompt
// GET /api/v1/ai/test?modelId=xxx&prompt=Hello
// POST /api/v1/ai/test with {"modelId": xxx, "prompt": "Hello"}
func (h *Handler) Test(c echo.Context) error {
	var req TestRequest

	// Try binding from both query params and body
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	// Also try query params directly for GET requests
	if req.ModelID == 0 {
		if idStr := c.QueryParam("modelId"); idStr != "" {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				return apierrors.BadRequest("无效的模型ID")
			}
			req.ModelID = id
		}
	}
	if req.Prompt == "" {
		req.Prompt = c.QueryParam("prompt")
	}

	if req.ModelID == 0 {
		return apierrors.BadRequest("缺少模型ID")
	}
	if req.Prompt == "" {
		req.Prompt = "Hello, please introduce yourself briefly."
	}

	result, err := h.chatSvc.TestModel(c.Request().Context(), req.ModelID, req.Prompt)
	if err != nil {
		return err
	}

	return response.Success(c, result)
}

// Chat performs a chat completion
// POST /api/v1/ai/chat
func (h *Handler) Chat(c echo.Context) error {
	var req ChatRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	if req.ModelID == 0 {
		return apierrors.BadRequest("缺少模型ID")
	}
	if len(req.Messages) == 0 {
		return apierrors.BadRequest("消息列表不能为空")
	}

	// Non-streaming response
	if !req.Stream {
		chatReq := &llm.ChatRequest{
			ModelID:  req.ModelID,
			Messages: req.Messages,
		}
		result, err := h.chatSvc.Chat(c.Request().Context(), chatReq)
		if err != nil {
			return err
		}
		return response.Success(c, result)
	}

	// Streaming response (SSE)
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().WriteHeader(200)

	chatReq := &llm.ChatRequest{
		ModelID:  req.ModelID,
		Messages: req.Messages,
	}

	err := h.chatSvc.ChatStream(c.Request().Context(), chatReq, func(chunk *llm.StreamChunk) error {
		if chunk.Done {
			// Send done event
			_, err := fmt.Fprintf(c.Response(), "event: message\ndata: {\"domain\":\"system\",\"type\":\"done\",\"content\":\"\"}\n\n")
			if err != nil {
				return err
			}
		} else {
			// Send content chunk
			_, err := fmt.Fprintf(c.Response(), "event: message\ndata: {\"domain\":\"llm\",\"type\":\"text\",\"content\":%q}\n\n", chunk.Content)
			if err != nil {
				return err
			}
		}
		c.Response().Flush()
		return nil
	})

	if err != nil {
		// Send error event
		fmt.Fprintf(c.Response(), "event: message\ndata: {\"domain\":\"system\",\"type\":\"error\",\"content\":%q}\n\n", err.Error())
		c.Response().Flush()
	}

	return nil
}
