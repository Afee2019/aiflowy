package system

import (
	"strconv"

	"github.com/aiflowy/aiflowy-go/internal/service"
	"github.com/aiflowy/aiflowy-go/pkg/response"
	"github.com/labstack/echo/v4"
)

// SysLogHandler 操作日志处理器
type SysLogHandler struct {
	svc *service.SysLogService
}

// NewSysLogHandler 创建 SysLogHandler
func NewSysLogHandler() *SysLogHandler {
	return &SysLogHandler{
		svc: service.NewSysLogService(),
	}
}

// Register 注册路由
func (h *SysLogHandler) Register(g *echo.Group) {
	g.GET("/page", h.Page)
	g.POST("/page", h.Page)
	g.POST("/remove", h.Delete)
}

// Page 分页查询
func (h *SysLogHandler) Page(c echo.Context) error {
	ctx := c.Request().Context()

	pageNum, _ := strconv.Atoi(c.QueryParam("pageNumber"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	actionName := c.QueryParam("actionName")
	actionType := c.QueryParam("actionType")

	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	list, total, err := h.svc.Page(ctx, pageNum, pageSize, actionName, actionType)
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

// Delete 删除日志
func (h *SysLogHandler) Delete(c echo.Context) error {
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
