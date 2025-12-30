package plugin

import (
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/aiflowy/aiflowy-go/internal/dto"
	apierrors "github.com/aiflowy/aiflowy-go/internal/errors"
	"github.com/aiflowy/aiflowy-go/internal/handler/auth"
	"github.com/aiflowy/aiflowy-go/internal/service"
	"github.com/aiflowy/aiflowy-go/pkg/response"
)

// Handler 插件 API Handler
type Handler struct {
	service *service.PluginService
}

// NewHandler 创建 Handler
func NewHandler() *Handler {
	return &Handler{
		service: service.NewPluginService(),
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

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(g *echo.Group) {
	// 插件
	plugin := g.Group("/plugin")
	plugin.POST("/getList", h.ListPlugins)
	plugin.GET("/getDetail", h.GetPlugin)
	plugin.POST("/plugin/save", h.SavePlugin)
	plugin.POST("/plugin/remove", h.DeletePlugin)
	plugin.GET("/pageByCategory", h.PageByCategory)

	// 插件工具
	pluginItem := g.Group("/pluginItem")
	pluginItem.POST("/toolsList", h.ListPluginItems)
	pluginItem.POST("/tool/search", h.GetPluginItem)
	pluginItem.POST("/tool/save", h.SavePluginItem)
	pluginItem.POST("/tool/update", h.UpdatePluginItem)
	pluginItem.POST("/remove", h.DeletePluginItem)
	pluginItem.POST("/tool/list", h.ListPluginItemsByBot)
	pluginItem.POST("/test", h.TestPluginItem)
	pluginItem.GET("/getTinyFlowData", h.GetTinyFlowData)

	// Bot-插件关联
	botPlugins := g.Group("/botPlugins")
	botPlugins.POST("/getList", h.GetBotPluginList)
	botPlugins.POST("/getBotPluginToolIds", h.GetBotPluginToolIDs)
	botPlugins.POST("/updateBotPluginToolIds", h.UpdateBotPluginToolIDs)
	botPlugins.POST("/doRemove", h.DeleteBotPlugin)

	// 插件分类
	pluginCategory := g.Group("/pluginCategory")
	pluginCategory.GET("/list", h.ListPluginCategories)
	pluginCategory.POST("/save", h.SavePluginCategory)
	pluginCategory.POST("/remove", h.DeletePluginCategory)
}

// ========================== Plugin ==========================

// ListPlugins 获取插件列表
func (h *Handler) ListPlugins(c echo.Context) error {
	ctx := c.Request().Context()
	_, tenantID, _ := getUserContext(c)

	plugins, err := h.service.ListPluginsWithTools(ctx, tenantID)
	if err != nil {
		return err
	}
	return response.Success(c, plugins)
}

// GetPlugin 获取插件详情
func (h *Handler) GetPlugin(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.QueryParam("id")

	plugin, err := h.service.GetPlugin(ctx, id)
	if err != nil {
		return err
	}
	return response.Success(c, plugin)
}

// SavePlugin 保存插件
func (h *Handler) SavePlugin(c echo.Context) error {
	ctx := c.Request().Context()
	userID, tenantID, deptID := getUserContext(c)

	var req dto.PluginSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	plugin, err := h.service.SavePlugin(ctx, &req, tenantID, userID, deptID)
	if err != nil {
		return err
	}
	return response.Success(c, plugin)
}

// DeletePlugin 删除插件
func (h *Handler) DeletePlugin(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		ID string `json:"id"`
	}
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if err := h.service.DeletePlugin(ctx, req.ID); err != nil {
		return err
	}
	return response.Success(c, true)
}

// PageByCategory 按分类分页获取插件
func (h *Handler) PageByCategory(c echo.Context) error {
	// 暂时返回全部列表
	return h.ListPlugins(c)
}

// ========================== PluginItem ==========================

// ListPluginItems 获取插件下的工具列表
func (h *Handler) ListPluginItems(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		PluginID string `json:"pluginId"`
		BotID    string `json:"botId"`
	}
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	items, err := h.service.ListPluginItems(ctx, req.PluginID, req.BotID)
	if err != nil {
		return err
	}
	return response.Success(c, items)
}

// GetPluginItem 获取插件工具详情
func (h *Handler) GetPluginItem(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		AIPluginToolID string `json:"aiPluginToolId"`
	}
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	item, err := h.service.GetPluginItem(ctx, req.AIPluginToolID)
	if err != nil {
		return err
	}
	return response.Success(c, item)
}

// SavePluginItem 保存插件工具
func (h *Handler) SavePluginItem(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.PluginItemSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	item, err := h.service.SavePluginItem(ctx, &req)
	if err != nil {
		return err
	}
	return response.Success(c, item)
}

// UpdatePluginItem 更新插件工具
func (h *Handler) UpdatePluginItem(c echo.Context) error {
	return h.SavePluginItem(c)
}

// DeletePluginItem 删除插件工具
func (h *Handler) DeletePluginItem(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	for _, id := range req.IDs {
		if err := h.service.DeletePluginItem(ctx, id); err != nil {
			return err
		}
	}
	return response.Success(c, true)
}

// ListPluginItemsByBot 获取 Bot 关联的插件工具列表
func (h *Handler) ListPluginItemsByBot(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		BotID string `json:"botId"`
	}
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	ids, err := h.service.GetBotPluginIDs(ctx, req.BotID)
	if err != nil {
		return err
	}
	return response.Success(c, ids)
}

// TestPluginItem 测试插件工具
func (h *Handler) TestPluginItem(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.PluginItemTestRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	// 使用 PluginToolService 执行测试
	toolService := service.NewPluginToolService()
	result, err := toolService.ExecutePluginTool(ctx, req.PluginToolID, req.InputData)
	if err != nil {
		return err
	}
	return response.Success(c, result)
}

// GetTinyFlowData 获取工作流节点数据
func (h *Handler) GetTinyFlowData(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.QueryParam("id")

	item, err := h.service.GetPluginItem(ctx, id)
	if err != nil {
		return err
	}

	// 构建工作流节点数据格式
	nodeData := map[string]interface{}{
		"pluginId":   item.ID,
		"pluginName": item.Name,
		"parameters": item.InputData,
		"outputDefs": item.OutputData,
	}
	return response.Success(c, nodeData)
}

// ========================== BotPlugin ==========================

// GetBotPluginList 获取 Bot 关联的插件列表
func (h *Handler) GetBotPluginList(c echo.Context) error {
	// 返回插件列表而非关联 ID
	return h.ListPlugins(c)
}

// GetBotPluginToolIDs 获取 Bot 关联的插件工具 ID 列表
func (h *Handler) GetBotPluginToolIDs(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		BotID string `json:"botId"`
	}
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	ids, err := h.service.GetBotPluginIDs(ctx, req.BotID)
	if err != nil {
		return err
	}
	return response.Success(c, ids)
}

// UpdateBotPluginToolIDs 更新 Bot-插件关联
func (h *Handler) UpdateBotPluginToolIDs(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.BotPluginUpdateRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if err := h.service.UpdateBotPlugins(ctx, &req); err != nil {
		return err
	}
	return response.Success(c, nil)
}

// DeleteBotPlugin 删除单个 Bot-插件关联
func (h *Handler) DeleteBotPlugin(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		BotID        string `json:"botId"`
		PluginToolID string `json:"pluginToolId"`
	}
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if err := h.service.DeleteBotPlugin(ctx, req.BotID, req.PluginToolID); err != nil {
		return err
	}
	return response.Success(c, true)
}

// ========================== PluginCategory ==========================

// ListPluginCategories 获取插件分类列表
func (h *Handler) ListPluginCategories(c echo.Context) error {
	ctx := c.Request().Context()

	categories, err := h.service.ListPluginCategories(ctx)
	if err != nil {
		return err
	}
	return response.Success(c, categories)
}

// SavePluginCategory 保存插件分类
func (h *Handler) SavePluginCategory(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.PluginCategorySaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	category, err := h.service.SavePluginCategory(ctx, &req)
	if err != nil {
		return err
	}
	return response.Success(c, category)
}

// DeletePluginCategory 删除插件分类
func (h *Handler) DeletePluginCategory(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	for _, id := range req.IDs {
		if err := h.service.DeletePluginCategory(ctx, id); err != nil {
			return err
		}
	}
	return response.Success(c, true)
}
