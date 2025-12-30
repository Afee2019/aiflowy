package system

import (
	"github.com/aiflowy/aiflowy-go/internal/handler/auth"
	"github.com/aiflowy/aiflowy-go/internal/service"
	"github.com/aiflowy/aiflowy-go/pkg/response"
	"github.com/labstack/echo/v4"
)

// SysOptionHandler 系统配置处理器
type SysOptionHandler struct {
	svc *service.SysOptionService
}

// NewSysOptionHandler 创建 SysOptionHandler
func NewSysOptionHandler() *SysOptionHandler {
	return &SysOptionHandler{
		svc: service.NewSysOptionService(),
	}
}

// Register 注册路由
func (h *SysOptionHandler) Register(g *echo.Group) {
	g.GET("/list", h.List)
	g.POST("/save", h.Save)
}

// List 获取配置列表
func (h *SysOptionHandler) List(c echo.Context) error {
	ctx := c.Request().Context()

	keys := c.QueryParams()["keys"]
	if len(keys) == 0 {
		keys = c.QueryParams()["keys[]"]
	}

	if len(keys) == 0 {
		return response.Success(c, map[string]interface{}{})
	}

	// 获取租户 ID
	tenantID := auth.GetCurrentTenantID(c)

	result, err := h.svc.GetMultiple(ctx, tenantID, keys)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, result)
}

// Save 保存配置
func (h *SysOptionHandler) Save(c echo.Context) error {
	ctx := c.Request().Context()

	var req map[string]string
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "参数错误")
	}

	if len(req) == 0 {
		return response.Success(c, nil)
	}

	// 获取租户 ID
	tenantID := auth.GetCurrentTenantID(c)

	if err := h.svc.SetMultiple(ctx, tenantID, req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil)
}
