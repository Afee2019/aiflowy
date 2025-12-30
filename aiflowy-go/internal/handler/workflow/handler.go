package workflow

import (
	"io"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/aiflowy/aiflowy-go/internal/dto"
	apierrors "github.com/aiflowy/aiflowy-go/internal/errors"
	"github.com/aiflowy/aiflowy-go/internal/handler/auth"
	"github.com/aiflowy/aiflowy-go/internal/service"
	"github.com/aiflowy/aiflowy-go/pkg/response"
)

// Handler 工作流 API Handler
type Handler struct {
	service  *service.WorkflowService
	executor *service.ChainExecutor
}

// NewHandler 创建 Handler
func NewHandler() *Handler {
	return &Handler{
		service:  service.NewWorkflowService(),
		executor: service.GetChainExecutor(),
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
	// 工作流
	workflow := g.Group("/workflow")
	workflow.POST("/list", h.ListWorkflows)
	workflow.GET("/list", h.ListWorkflows)
	workflow.GET("/getDetail", h.GetWorkflow)
	workflow.POST("/save", h.SaveWorkflow)
	workflow.POST("/remove", h.DeleteWorkflow)
	workflow.GET("/copy", h.CopyWorkflow)
	workflow.GET("/getRunningParameters", h.GetRunningParameters)
	workflow.POST("/importWorkFlow", h.ImportWorkflow)
	workflow.GET("/exportWorkFlow", h.ExportWorkflow)

	// 工作流执行 (Stage 11 将实现完整功能)
	workflow.POST("/runAsync", h.RunAsync)
	workflow.POST("/getChainStatus", h.GetChainStatus)
	workflow.POST("/resume", h.Resume)
	workflow.POST("/singleRun", h.SingleRun)

	// 工作流分类
	workflowCategory := g.Group("/workflowCategory")
	workflowCategory.GET("/list", h.ListWorkflowCategories)
	workflowCategory.POST("/list", h.ListWorkflowCategories)
	workflowCategory.POST("/save", h.SaveWorkflowCategory)
	workflowCategory.POST("/remove", h.DeleteWorkflowCategory)

	// Bot-工作流关联
	botWorkflow := g.Group("/botWorkflow")
	botWorkflow.GET("/list", h.ListBotWorkflows)
	botWorkflow.POST("/list", h.ListBotWorkflows)
	botWorkflow.POST("/updateBotWorkflowIds", h.UpdateBotWorkflowIDs)
	botWorkflow.POST("/getBotWorkflowIds", h.GetBotWorkflowIDs)
	botWorkflow.POST("/remove", h.DeleteBotWorkflow)
}

// ========================== Workflow ==========================

// ListWorkflows 获取工作流列表
func (h *Handler) ListWorkflows(c echo.Context) error {
	ctx := c.Request().Context()
	_, tenantID, _ := getUserContext(c)

	workflows, err := h.service.ListWorkflows(ctx, tenantID)
	if err != nil {
		return err
	}
	return response.Success(c, workflows)
}

// GetWorkflow 获取工作流详情
func (h *Handler) GetWorkflow(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.QueryParam("id")

	workflow, err := h.service.GetWorkflow(ctx, id)
	if err != nil {
		return err
	}
	return response.Success(c, workflow)
}

// SaveWorkflow 保存工作流
func (h *Handler) SaveWorkflow(c echo.Context) error {
	ctx := c.Request().Context()
	userID, tenantID, deptID := getUserContext(c)

	var req dto.WorkflowSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	workflow, err := h.service.SaveWorkflow(ctx, &req, tenantID, userID, deptID)
	if err != nil {
		return err
	}
	return response.Success(c, workflow)
}

// DeleteWorkflow 删除工作流
func (h *Handler) DeleteWorkflow(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	for _, id := range req.IDs {
		if err := h.service.DeleteWorkflow(ctx, id); err != nil {
			return err
		}
	}
	return response.Success(c, true)
}

// CopyWorkflow 复制工作流
func (h *Handler) CopyWorkflow(c echo.Context) error {
	ctx := c.Request().Context()
	userID, tenantID, deptID := getUserContext(c)
	id := c.QueryParam("id")

	_, err := h.service.CopyWorkflow(ctx, id, tenantID, userID, deptID)
	if err != nil {
		return err
	}
	return response.Success(c, nil)
}

// GetRunningParameters 获取工作流运行参数
func (h *Handler) GetRunningParameters(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.QueryParam("id")

	result, err := h.service.GetRunningParameters(ctx, id)
	if err != nil {
		return err
	}
	return response.Success(c, result)
}

// ImportWorkflow 导入工作流
func (h *Handler) ImportWorkflow(c echo.Context) error {
	ctx := c.Request().Context()
	userID, tenantID, deptID := getUserContext(c)

	// 获取上传的 JSON 文件
	file, err := c.FormFile("jsonFile")
	if err != nil {
		return apierrors.BadRequest("请上传工作流 JSON 文件")
	}

	src, err := file.Open()
	if err != nil {
		return apierrors.InternalError("读取文件失败")
	}
	defer src.Close()

	content, err := io.ReadAll(src)
	if err != nil {
		return apierrors.InternalError("读取文件内容失败")
	}

	// 获取其他表单参数
	req := &dto.WorkflowSaveRequest{
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
		Icon:        c.FormValue("icon"),
		Content:     string(content),
	}

	workflow, err := h.service.SaveWorkflow(ctx, req, tenantID, userID, deptID)
	if err != nil {
		return err
	}
	return response.Success(c, workflow)
}

// ExportWorkflow 导出工作流
func (h *Handler) ExportWorkflow(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.QueryParam("id")

	workflow, err := h.service.GetWorkflow(ctx, id)
	if err != nil {
		return err
	}

	// 直接返回工作流内容
	return response.Success(c, workflow.Content)
}

// ========================== 工作流执行 ==========================

// RunAsync 异步运行工作流
func (h *Handler) RunAsync(c echo.Context) error {
	ctx := c.Request().Context()
	userID, _, _ := getUserContext(c)

	var req dto.WorkflowRunRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if req.ID == "" {
		return apierrors.BadRequest("工作流 ID 不能为空")
	}

	// 获取用户信息
	claims := auth.GetClaims(c)
	createdBy := ""
	if claims != nil {
		createdBy = claims.Nickname
	}

	// 执行工作流
	executeID, err := h.executor.ExecuteAsync(ctx, req.ID, req.Variables, strconv.FormatInt(userID, 10), createdBy)
	if err != nil {
		return apierrors.InternalError(err.Error())
	}

	return response.Success(c, executeID)
}

// GetChainStatus 获取工作流执行状态
func (h *Handler) GetChainStatus(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.ChainStatusRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if req.ExecuteID == "" {
		return apierrors.BadRequest("执行 ID 不能为空")
	}

	chainInfo, err := h.executor.GetChainStatus(ctx, req.ExecuteID, req.Nodes)
	if err != nil {
		return apierrors.InternalError(err.Error())
	}

	return response.Success(c, chainInfo)
}

// Resume 恢复工作流执行
func (h *Handler) Resume(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.WorkflowResumeRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if req.ExecuteID == "" {
		return apierrors.BadRequest("执行 ID 不能为空")
	}

	if err := h.executor.Resume(ctx, req.ExecuteID, req.ConfirmParams); err != nil {
		return apierrors.InternalError(err.Error())
	}

	return response.Success(c, nil)
}

// SingleRun 单节点运行
func (h *Handler) SingleRun(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.WorkflowSingleRunRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if req.WorkflowID == "" {
		return apierrors.BadRequest("工作流 ID 不能为空")
	}
	if req.NodeID == "" {
		return apierrors.BadRequest("节点 ID 不能为空")
	}

	result, err := h.executor.ExecuteNode(ctx, req.WorkflowID, req.NodeID, req.Variables)
	if err != nil {
		return apierrors.InternalError(err.Error())
	}

	return response.Success(c, result)
}

// ========================== WorkflowCategory ==========================

// ListWorkflowCategories 获取工作流分类列表
func (h *Handler) ListWorkflowCategories(c echo.Context) error {
	ctx := c.Request().Context()

	categories, err := h.service.ListWorkflowCategories(ctx)
	if err != nil {
		return err
	}
	return response.Success(c, categories)
}

// SaveWorkflowCategory 保存工作流分类
func (h *Handler) SaveWorkflowCategory(c echo.Context) error {
	ctx := c.Request().Context()
	userID, _, _ := getUserContext(c)

	var req dto.WorkflowCategorySaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	category, err := h.service.SaveWorkflowCategory(ctx, &req, userID)
	if err != nil {
		return err
	}
	return response.Success(c, category)
}

// DeleteWorkflowCategory 删除工作流分类
func (h *Handler) DeleteWorkflowCategory(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	for _, id := range req.IDs {
		if err := h.service.DeleteWorkflowCategory(ctx, id); err != nil {
			return err
		}
	}
	return response.Success(c, true)
}

// ========================== BotWorkflow ==========================

// ListBotWorkflows 获取 Bot 关联的工作流列表
func (h *Handler) ListBotWorkflows(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		BotID string `json:"botId" query:"botId"`
	}
	if c.Request().Method == "POST" {
		if err := c.Bind(&req); err != nil {
			return apierrors.BadRequest("无效的请求参数")
		}
	} else {
		req.BotID = c.QueryParam("botId")
	}

	botWorkflows, err := h.service.ListBotWorkflows(ctx, req.BotID)
	if err != nil {
		return err
	}
	return response.Success(c, botWorkflows)
}

// GetBotWorkflowIDs 获取 Bot 关联的工作流 ID 列表
func (h *Handler) GetBotWorkflowIDs(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		BotID string `json:"botId"`
	}
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	ids, err := h.service.GetBotWorkflowIDs(ctx, req.BotID)
	if err != nil {
		return err
	}
	return response.Success(c, ids)
}

// UpdateBotWorkflowIDs 更新 Bot-工作流关联
func (h *Handler) UpdateBotWorkflowIDs(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.BotWorkflowUpdateRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if err := h.service.UpdateBotWorkflows(ctx, &req); err != nil {
		return err
	}
	return response.Success(c, nil)
}

// DeleteBotWorkflow 删除单个 Bot-工作流关联
func (h *Handler) DeleteBotWorkflow(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		BotID      string `json:"botId"`
		WorkflowID string `json:"workflowId"`
	}
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if err := h.service.DeleteBotWorkflow(ctx, req.BotID, req.WorkflowID); err != nil {
		return err
	}
	return response.Success(c, true)
}
