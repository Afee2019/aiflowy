package system

import (
	"strconv"
	"time"

	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/internal/handler/auth"
	"github.com/aiflowy/aiflowy-go/internal/service"
	"github.com/aiflowy/aiflowy-go/pkg/response"
	"github.com/labstack/echo/v4"
)

// SysJobHandler 定时任务处理器
type SysJobHandler struct {
	svc *service.SysJobService
}

// NewSysJobHandler 创建 SysJobHandler
func NewSysJobHandler() *SysJobHandler {
	return &SysJobHandler{
		svc: service.NewSysJobService(),
	}
}

// Register 注册路由
func (h *SysJobHandler) Register(g *echo.Group) {
	g.GET("/page", h.Page)
	g.POST("/page", h.Page)
	g.POST("/save", h.Save)
	g.POST("/remove", h.Delete)
	g.GET("/start", h.Start)
	g.GET("/stop", h.Stop)
	g.GET("/getNextTimes", h.GetNextTimes)
}

// Page 分页查询
func (h *SysJobHandler) Page(c echo.Context) error {
	ctx := c.Request().Context()

	pageNum, _ := strconv.Atoi(c.QueryParam("pageNumber"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	jobName := c.QueryParam("jobName")

	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	list, total, err := h.svc.Page(ctx, pageNum, pageSize, jobName)
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

// Save 保存任务
func (h *SysJobHandler) Save(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		ID              int64  `json:"id,string"`
		JobName         string `json:"jobName"`
		JobType         int    `json:"jobType"`
		JobParams       string `json:"jobParams,omitempty"`
		CronExpression  string `json:"cronExpression"`
		AllowConcurrent int    `json:"allowConcurrent"`
		MisfirePolicy   int    `json:"misfirePolicy"`
		Options         string `json:"options,omitempty"`
		Status          int    `json:"status"`
		Remark          string `json:"remark,omitempty"`
	}
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "参数错误")
	}

	userID, tenantID, deptID := auth.GetUserContext(c)

	job := &entity.SysJob{
		ID:              req.ID,
		JobName:         req.JobName,
		JobType:         req.JobType,
		CronExpression:  req.CronExpression,
		AllowConcurrent: req.AllowConcurrent,
		MisfirePolicy:   req.MisfirePolicy,
		Status:          req.Status,
		TenantID:        tenantID,
		DeptID:          deptID,
	}

	if req.JobParams != "" {
		job.JobParams = &req.JobParams
	}
	if req.Options != "" {
		job.Options = &req.Options
	}
	if req.Remark != "" {
		job.Remark = &req.Remark
	}

	if req.ID == 0 {
		// 创建
		job.CreatedBy = &userID
		now := time.Now()
		job.Created = &now
		job.Modified = &now
		job.ModifiedBy = &userID
		if err := h.svc.Create(ctx, job); err != nil {
			return response.BadRequest(c, err.Error())
		}
	} else {
		// 更新
		job.ModifiedBy = &userID
		if err := h.svc.Update(ctx, job); err != nil {
			return response.BadRequest(c, err.Error())
		}
	}

	return response.Success(c, map[string]interface{}{"id": job.ID})
}

// Delete 删除任务
func (h *SysJobHandler) Delete(c echo.Context) error {
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

// Start 启动任务
func (h *SysJobHandler) Start(c echo.Context) error {
	ctx := c.Request().Context()

	idStr := c.QueryParam("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.svc.Start(ctx, id); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil)
}

// Stop 停止任务
func (h *SysJobHandler) Stop(c echo.Context) error {
	ctx := c.Request().Context()

	idStr := c.QueryParam("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.svc.Stop(ctx, id); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil)
}

// GetNextTimes 获取下次执行时间
func (h *SysJobHandler) GetNextTimes(c echo.Context) error {
	cronExpression := c.QueryParam("cronExpression")
	if cronExpression == "" {
		return response.BadRequest(c, "cron 表达式不能为空")
	}

	times, err := h.svc.GetNextTimes(cronExpression, 5)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, times)
}
