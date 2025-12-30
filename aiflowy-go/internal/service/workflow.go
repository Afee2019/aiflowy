package service

import (
	"context"
	"strconv"
	"time"

	"github.com/aiflowy/aiflowy-go/internal/dto"
	"github.com/aiflowy/aiflowy-go/internal/entity"
	apierrors "github.com/aiflowy/aiflowy-go/internal/errors"
	"github.com/aiflowy/aiflowy-go/internal/repository"
	"github.com/aiflowy/aiflowy-go/pkg/snowflake"
)

// WorkflowService 工作流服务
type WorkflowService struct {
	repo *repository.WorkflowRepository
}

// NewWorkflowService 创建 WorkflowService
func NewWorkflowService() *WorkflowService {
	return &WorkflowService{
		repo: repository.NewWorkflowRepository(),
	}
}

// ========================== Workflow ==========================

// GetWorkflow 获取工作流
func (s *WorkflowService) GetWorkflow(ctx context.Context, id string) (*entity.Workflow, error) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, apierrors.BadRequest("无效的工作流 ID")
	}

	workflow, err := s.repo.GetWorkflowByID(ctx, idInt)
	if err != nil {
		return nil, apierrors.InternalError("获取工作流失败")
	}
	if workflow == nil {
		return nil, apierrors.NotFound("工作流不存在")
	}
	return workflow, nil
}

// GetWorkflowByAlias 根据别名获取工作流
func (s *WorkflowService) GetWorkflowByAlias(ctx context.Context, alias string) (*entity.Workflow, error) {
	return s.repo.GetWorkflowByAlias(ctx, alias)
}

// ListWorkflows 获取工作流列表
func (s *WorkflowService) ListWorkflows(ctx context.Context, tenantID int64) ([]*entity.Workflow, error) {
	return s.repo.ListWorkflows(ctx, tenantID)
}

// SaveWorkflow 保存工作流 (创建或更新)
func (s *WorkflowService) SaveWorkflow(ctx context.Context, req *dto.WorkflowSaveRequest, tenantID, userID, deptID int64) (*entity.Workflow, error) {
	var workflow *entity.Workflow
	var isNew bool

	if req.ID != "" {
		// 更新
		idInt, err := strconv.ParseInt(req.ID, 10, 64)
		if err != nil {
			return nil, apierrors.BadRequest("无效的工作流 ID")
		}
		workflow, err = s.repo.GetWorkflowByID(ctx, idInt)
		if err != nil {
			return nil, apierrors.InternalError("获取工作流失败")
		}
		if workflow == nil {
			return nil, apierrors.NotFound("工作流不存在")
		}

		// 检查别名是否重复
		if req.Alias != "" && req.Alias != workflow.Alias {
			existing, _ := s.repo.GetWorkflowByAlias(ctx, req.Alias)
			if existing != nil && existing.ID != workflow.ID {
				return nil, apierrors.BadRequest("别名已存在")
			}
		}

		workflow.Modified = time.Now()
		workflow.ModifiedBy = userID
	} else {
		// 创建
		isNew = true

		// 检查别名是否重复
		if req.Alias != "" {
			existing, _ := s.repo.GetWorkflowByAlias(ctx, req.Alias)
			if existing != nil {
				return nil, apierrors.BadRequest("别名已存在")
			}
		}

		workflow = &entity.Workflow{
			ID:        snowflake.MustGenerateID(),
			Created:   time.Now(),
			TenantID:  tenantID,
			CreatedBy: userID,
			DeptID:    deptID,
		}
	}

	// 更新字段
	workflow.Alias = req.Alias
	workflow.Title = req.Title
	workflow.Description = req.Description
	workflow.Icon = req.Icon
	workflow.Content = req.Content
	workflow.EnglishName = req.EnglishName
	workflow.Status = req.Status

	if req.CategoryID != "" {
		categoryID, _ := strconv.ParseInt(req.CategoryID, 10, 64)
		workflow.CategoryID = categoryID
	}

	if isNew {
		if err := s.repo.CreateWorkflow(ctx, workflow); err != nil {
			return nil, apierrors.InternalError("创建工作流失败")
		}
	} else {
		if err := s.repo.UpdateWorkflow(ctx, workflow); err != nil {
			return nil, apierrors.InternalError("更新工作流失败")
		}
	}

	return workflow, nil
}

// DeleteWorkflow 删除工作流
func (s *WorkflowService) DeleteWorkflow(ctx context.Context, id string) error {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的工作流 ID")
	}

	// 检查是否有 Bot 关联
	exists, err := s.repo.ExistsBotWorkflow(ctx, idInt)
	if err != nil {
		return apierrors.InternalError("检查关联失败")
	}
	if exists {
		return apierrors.BadRequest("此工作流还关联有 Bot，请先取消关联后再删除")
	}

	return s.repo.DeleteWorkflow(ctx, idInt)
}

// CopyWorkflow 复制工作流
func (s *WorkflowService) CopyWorkflow(ctx context.Context, id string, tenantID, userID, deptID int64) (*entity.Workflow, error) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, apierrors.BadRequest("无效的工作流 ID")
	}

	// 获取原工作流
	original, err := s.repo.GetWorkflowByID(ctx, idInt)
	if err != nil {
		return nil, apierrors.InternalError("获取工作流失败")
	}
	if original == nil {
		return nil, apierrors.NotFound("工作流不存在")
	}

	// 创建副本
	newWorkflow := &entity.Workflow{
		ID:          snowflake.MustGenerateID(),
		Alias:       snowflake.MustGenerateIDString(), // 使用新的 UUID 作为别名
		DeptID:      deptID,
		TenantID:    tenantID,
		Title:       original.Title + " (副本)",
		Description: original.Description,
		Icon:        original.Icon,
		Content:     original.Content,
		Created:     time.Now(),
		CreatedBy:   userID,
		EnglishName: original.EnglishName,
		Status:      original.Status,
		CategoryID:  original.CategoryID,
	}

	if err := s.repo.CreateWorkflow(ctx, newWorkflow); err != nil {
		return nil, apierrors.InternalError("复制工作流失败")
	}

	return newWorkflow, nil
}

// GetRunningParameters 获取工作流运行参数
func (s *WorkflowService) GetRunningParameters(ctx context.Context, id string) (*dto.RunningParametersResponse, error) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, apierrors.BadRequest("无效的工作流 ID")
	}

	workflow, err := s.repo.GetWorkflowByID(ctx, idInt)
	if err != nil {
		return nil, apierrors.InternalError("获取工作流失败")
	}
	if workflow == nil {
		return nil, apierrors.NotFound("工作流不存在")
	}

	// 解析工作流 DSL 获取参数
	parser := NewWorkflowDSLParser()
	definition, err := parser.Parse(workflow.Content)
	if err != nil {
		return nil, apierrors.BadRequest("工作流配置解析失败: " + err.Error())
	}

	// 获取开始节点的参数
	parameters := parser.GetStartParameters(definition)

	return &dto.RunningParametersResponse{
		Parameters:  parameters,
		Title:       workflow.Title,
		Description: workflow.Description,
		Icon:        workflow.Icon,
	}, nil
}

// ========================== WorkflowCategory ==========================

// GetWorkflowCategory 获取工作流分类
func (s *WorkflowService) GetWorkflowCategory(ctx context.Context, id string) (*entity.WorkflowCategory, error) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, apierrors.BadRequest("无效的分类 ID")
	}

	category, err := s.repo.GetWorkflowCategoryByID(ctx, idInt)
	if err != nil {
		return nil, apierrors.InternalError("获取分类失败")
	}
	if category == nil {
		return nil, apierrors.NotFound("分类不存在")
	}
	return category, nil
}

// ListWorkflowCategories 获取工作流分类列表
func (s *WorkflowService) ListWorkflowCategories(ctx context.Context) ([]*entity.WorkflowCategory, error) {
	return s.repo.ListWorkflowCategories(ctx)
}

// SaveWorkflowCategory 保存工作流分类
func (s *WorkflowService) SaveWorkflowCategory(ctx context.Context, req *dto.WorkflowCategorySaveRequest, userID int64) (*entity.WorkflowCategory, error) {
	var category *entity.WorkflowCategory
	var isNew bool

	if req.ID != "" {
		// 更新
		idInt, err := strconv.ParseInt(req.ID, 10, 64)
		if err != nil {
			return nil, apierrors.BadRequest("无效的分类 ID")
		}
		category, err = s.repo.GetWorkflowCategoryByID(ctx, idInt)
		if err != nil {
			return nil, apierrors.InternalError("获取分类失败")
		}
		if category == nil {
			return nil, apierrors.NotFound("分类不存在")
		}
		category.Modified = time.Now()
		category.ModifiedBy = userID
	} else {
		// 创建
		isNew = true
		category = &entity.WorkflowCategory{
			ID:         snowflake.MustGenerateID(),
			Created:    time.Now(),
			CreatedBy:  userID,
			Modified:   time.Now(),
			ModifiedBy: userID,
		}
	}

	// 更新字段
	category.CategoryName = req.CategoryName
	category.SortNo = req.SortNo
	category.Status = req.Status

	if isNew {
		if err := s.repo.CreateWorkflowCategory(ctx, category); err != nil {
			return nil, apierrors.InternalError("创建分类失败")
		}
	} else {
		if err := s.repo.UpdateWorkflowCategory(ctx, category); err != nil {
			return nil, apierrors.InternalError("更新分类失败")
		}
	}

	return category, nil
}

// DeleteWorkflowCategory 删除工作流分类
func (s *WorkflowService) DeleteWorkflowCategory(ctx context.Context, id string) error {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的分类 ID")
	}
	return s.repo.DeleteWorkflowCategory(ctx, idInt)
}

// ========================== BotWorkflow ==========================

// GetBotWorkflowIDs 获取 Bot 关联的工作流 ID 列表
func (s *WorkflowService) GetBotWorkflowIDs(ctx context.Context, botID string) ([]string, error) {
	botIDInt, err := strconv.ParseInt(botID, 10, 64)
	if err != nil {
		return nil, apierrors.BadRequest("无效的 Bot ID")
	}

	ids, err := s.repo.GetBotWorkflowIDs(ctx, botIDInt)
	if err != nil {
		return nil, apierrors.InternalError("获取关联失败")
	}

	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = strconv.FormatInt(id, 10)
	}
	return result, nil
}

// ListBotWorkflows 获取 Bot 关联的工作流列表
func (s *WorkflowService) ListBotWorkflows(ctx context.Context, botID string) ([]*entity.BotWorkflow, error) {
	botIDInt, err := strconv.ParseInt(botID, 10, 64)
	if err != nil {
		return nil, apierrors.BadRequest("无效的 Bot ID")
	}
	return s.repo.ListBotWorkflows(ctx, botIDInt)
}

// ListWorkflowsByBotID 获取 Bot 关联的工作流
func (s *WorkflowService) ListWorkflowsByBotID(ctx context.Context, botID int64) ([]*entity.Workflow, error) {
	return s.repo.ListWorkflowsByBotID(ctx, botID)
}

// UpdateBotWorkflows 更新 Bot-工作流关联
func (s *WorkflowService) UpdateBotWorkflows(ctx context.Context, req *dto.BotWorkflowUpdateRequest) error {
	botIDInt, err := strconv.ParseInt(req.BotID, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的 Bot ID")
	}

	var workflowIDs []int64
	for _, idStr := range req.WorkflowIDs {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			continue
		}
		workflowIDs = append(workflowIDs, id)
	}

	return s.repo.SaveBotWorkflows(ctx, botIDInt, workflowIDs)
}

// DeleteBotWorkflow 删除单个 Bot-工作流关联
func (s *WorkflowService) DeleteBotWorkflow(ctx context.Context, botID, workflowID string) error {
	botIDInt, err := strconv.ParseInt(botID, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的 Bot ID")
	}

	workflowIDInt, err := strconv.ParseInt(workflowID, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的工作流 ID")
	}

	return s.repo.DeleteBotWorkflow(ctx, botIDInt, workflowIDInt)
}
