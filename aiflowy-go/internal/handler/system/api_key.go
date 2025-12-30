package system

import (
	"strconv"

	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/internal/handler/auth"
	"github.com/aiflowy/aiflowy-go/internal/service"
	"github.com/aiflowy/aiflowy-go/pkg/response"
	"github.com/labstack/echo/v4"
)

// SysApiKeyHandler 系统 API 密钥处理器
type SysApiKeyHandler struct {
	svc *service.SysApiKeyService
}

// NewSysApiKeyHandler 创建 SysApiKeyHandler
func NewSysApiKeyHandler() *SysApiKeyHandler {
	return &SysApiKeyHandler{
		svc: service.NewSysApiKeyService(),
	}
}

// Register 注册路由
func (h *SysApiKeyHandler) Register(g *echo.Group) {
	g.POST("/key/save", h.Generate)
	g.POST("/save", h.Update)
	g.POST("/remove", h.Delete)
	g.GET("/page", h.Page)
	g.GET("/resources", h.ListResources)
}

// Generate 生成新的 API 密钥
func (h *SysApiKeyHandler) Generate(c echo.Context) error {
	ctx := c.Request().Context()

	// 获取登录用户信息
	userID, tenantID, deptID := auth.GetUserContext(c)

	key, err := h.svc.Generate(ctx, userID, tenantID, deptID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, map[string]interface{}{
		"id":     key.ID,
		"apiKey": key.ApiKey,
	})
}

// Update 更新 API 密钥
func (h *SysApiKeyHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		ID            int64   `json:"id,string"`
		Status        int     `json:"status"`
		PermissionIds []int64 `json:"permissionIds"`
	}
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "参数错误")
	}

	key := &entity.SysApiKey{
		ID:            req.ID,
		Status:        req.Status,
		PermissionIds: req.PermissionIds,
	}

	if err := h.svc.Update(ctx, key); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil)
}

// Delete 删除 API 密钥
func (h *SysApiKeyHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		ID int64 `json:"id,string"`
	}
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "参数错误")
	}

	if err := h.svc.Delete(ctx, req.ID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil)
}

// Page 分页查询
func (h *SysApiKeyHandler) Page(c echo.Context) error {
	ctx := c.Request().Context()

	pageNum, _ := strconv.Atoi(c.QueryParam("pageNumber"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	list, total, err := h.svc.Page(ctx, pageNum, pageSize)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, map[string]interface{}{
		"list":     list,
		"total":    total,
		"pageNo":   pageNum,
		"pageSize": pageSize,
	})
}

// ListResources 获取所有可授权的资源
func (h *SysApiKeyHandler) ListResources(c echo.Context) error {
	ctx := c.Request().Context()

	resources, err := h.svc.ListResources(ctx)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, resources)
}
