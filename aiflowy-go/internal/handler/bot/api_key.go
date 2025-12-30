package bot

import (
	"strconv"

	"github.com/aiflowy/aiflowy-go/internal/dto"
	"github.com/aiflowy/aiflowy-go/internal/handler/auth"
	"github.com/aiflowy/aiflowy-go/internal/service"
	"github.com/aiflowy/aiflowy-go/pkg/response"
	"github.com/labstack/echo/v4"
)

// BotApiKeyHandler Bot API 密钥处理器
type BotApiKeyHandler struct {
	svc *service.BotApiKeyService
}

// NewBotApiKeyHandler 创建 BotApiKeyHandler
func NewBotApiKeyHandler() *BotApiKeyHandler {
	return &BotApiKeyHandler{
		svc: service.NewBotApiKeyService(),
	}
}

// Register 注册路由
func (h *BotApiKeyHandler) Register(g *echo.Group) {
	g.POST("/addKey", h.AddKey)
	g.GET("/list", h.List)
	g.POST("/list", h.List)
	g.POST("/remove", h.Delete)
}

// AddKey 生成 Bot API 密钥
func (h *BotApiKeyHandler) AddKey(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		BotID string `json:"botId"`
	}
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "参数错误")
	}

	botID, _ := strconv.ParseInt(req.BotID, 10, 64)
	if botID == 0 {
		return response.BadRequest(c, "botId 不能为空")
	}

	userID := auth.GetCurrentUserID(c)

	apiKey, err := h.svc.GenerateByBotID(ctx, botID, userID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, apiKey)
}

// List 获取 Bot 的 API 密钥列表
func (h *BotApiKeyHandler) List(c echo.Context) error {
	ctx := c.Request().Context()

	botIDStr := c.QueryParam("botId")
	if botIDStr == "" {
		var req struct {
			BotID string `json:"botId"`
		}
		c.Bind(&req)
		botIDStr = req.BotID
	}

	botID, _ := strconv.ParseInt(botIDStr, 10, 64)
	if botID == 0 {
		return response.BadRequest(c, "botId 不能为空")
	}

	list, err := h.svc.ListByBotID(ctx, botID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, list)
}

// Delete 删除 API 密钥
func (h *BotApiKeyHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.IDRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "参数错误")
	}

	if err := h.svc.Delete(ctx, req.ID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil)
}
