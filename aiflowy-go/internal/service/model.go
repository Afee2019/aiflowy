package service

import (
	"context"
	"time"

	"github.com/aiflowy/aiflowy-go/internal/dto"
	"github.com/aiflowy/aiflowy-go/internal/entity"
	apierrors "github.com/aiflowy/aiflowy-go/internal/errors"
	"github.com/aiflowy/aiflowy-go/internal/repository"
	"github.com/aiflowy/aiflowy-go/pkg/snowflake"
)

// ModelService handles model-related business logic
type ModelService struct {
	repo *repository.ModelRepository
}

// NewModelService creates a new ModelService
func NewModelService() *ModelService {
	return &ModelService{
		repo: repository.GetModelRepository(),
	}
}

// ========== Model Provider Operations ==========

// GetProviderDetail gets a provider by ID
func (s *ModelService) GetProviderDetail(ctx context.Context, id int64) (*entity.ModelProvider, error) {
	provider, err := s.repo.GetProviderByID(ctx, id)
	if err != nil {
		return nil, apierrors.InternalError("查询供应商失败")
	}
	if provider == nil {
		return nil, apierrors.NotFound("供应商不存在")
	}
	return provider, nil
}

// ListProviders lists all providers
func (s *ModelService) ListProviders(ctx context.Context, req *dto.ModelProviderListRequest) ([]*entity.ModelProvider, error) {
	providers, err := s.repo.ListProviders(ctx, req)
	if err != nil {
		return nil, apierrors.InternalError("查询供应商列表失败")
	}
	return providers, nil
}

// PageProviders returns paginated providers
func (s *ModelService) PageProviders(ctx context.Context, pageReq *dto.PageRequest, filter *dto.ModelProviderListRequest) (*dto.PageResponse, error) {
	providers, total, err := s.repo.PageProviders(ctx, pageReq, filter)
	if err != nil {
		return nil, apierrors.InternalError("查询供应商列表失败")
	}
	return dto.NewPageResponse(pageReq.GetPageNumber(), pageReq.GetPageSize(), total, providers), nil
}

// SaveProvider creates a new provider
func (s *ModelService) SaveProvider(ctx context.Context, req *dto.ModelProviderSaveRequest, operatorID int64) (*entity.ModelProvider, error) {
	now := time.Now()
	provider := &entity.ModelProvider{
		ID:           snowflake.MustGenerateID(),
		ProviderName: req.ProviderName,
		ProviderType: req.ProviderType,
		Icon:         req.Icon,
		APIKey:       req.APIKey,
		Endpoint:     req.Endpoint,
		ChatPath:     req.ChatPath,
		EmbedPath:    req.EmbedPath,
		RerankPath:   req.RerankPath,
		Created:      now,
		CreatedBy:    operatorID,
		Modified:     now,
		ModifiedBy:   operatorID,
	}

	if err := s.repo.CreateProvider(ctx, provider); err != nil {
		return nil, apierrors.InternalError("创建供应商失败")
	}
	return provider, nil
}

// UpdateProvider updates an existing provider
func (s *ModelService) UpdateProvider(ctx context.Context, req *dto.ModelProviderSaveRequest, operatorID int64) (*entity.ModelProvider, error) {
	if req.ID == 0 {
		return nil, apierrors.BadRequest("供应商ID不能为空")
	}

	existing, err := s.repo.GetProviderByID(ctx, req.ID)
	if err != nil {
		return nil, apierrors.InternalError("查询供应商失败")
	}
	if existing == nil {
		return nil, apierrors.NotFound("供应商不存在")
	}

	existing.ProviderName = req.ProviderName
	existing.ProviderType = req.ProviderType
	existing.Icon = req.Icon
	existing.APIKey = req.APIKey
	existing.Endpoint = req.Endpoint
	existing.ChatPath = req.ChatPath
	existing.EmbedPath = req.EmbedPath
	existing.RerankPath = req.RerankPath
	existing.Modified = time.Now()
	existing.ModifiedBy = operatorID

	if err := s.repo.UpdateProvider(ctx, existing); err != nil {
		return nil, apierrors.InternalError("更新供应商失败")
	}
	return existing, nil
}

// RemoveProvider deletes a provider
func (s *ModelService) RemoveProvider(ctx context.Context, id int64) error {
	// Check if provider exists
	existing, err := s.repo.GetProviderByID(ctx, id)
	if err != nil {
		return apierrors.InternalError("查询供应商失败")
	}
	if existing == nil {
		return apierrors.NotFound("供应商不存在")
	}

	// Check if provider has associated models
	count, err := s.repo.CountModelsByProvider(ctx, id)
	if err != nil {
		return apierrors.InternalError("查询关联模型失败")
	}
	if count > 0 {
		return apierrors.BadRequest("该供应商下存在模型，无法删除")
	}

	if err := s.repo.DeleteProvider(ctx, id); err != nil {
		return apierrors.InternalError("删除供应商失败")
	}
	return nil
}

// ========== Model Operations ==========

// GetModelDetail gets a model by ID
func (s *ModelService) GetModelDetail(ctx context.Context, id int64) (*entity.ModelWithProvider, error) {
	model, err := s.repo.GetModelWithProvider(ctx, id)
	if err != nil {
		return nil, apierrors.InternalError("查询模型失败")
	}
	if model == nil {
		return nil, apierrors.NotFound("模型不存在")
	}
	return model, nil
}

// ListModels lists models with filters
func (s *ModelService) ListModels(ctx context.Context, req *dto.ModelListRequest) ([]*entity.Model, error) {
	models, err := s.repo.ListModels(ctx, req)
	if err != nil {
		return nil, apierrors.InternalError("查询模型列表失败")
	}
	return models, nil
}

// PageModels returns paginated models
func (s *ModelService) PageModels(ctx context.Context, req *dto.ModelListRequest) (*dto.PageResponse, error) {
	models, total, err := s.repo.PageModels(ctx, req)
	if err != nil {
		return nil, apierrors.InternalError("查询模型列表失败")
	}
	return dto.NewPageResponse(req.GetPageNumber(), req.GetPageSize(), total, models), nil
}

// GetList returns models grouped by type and group name
func (s *ModelService) GetList(ctx context.Context, req *dto.ModelByProviderRequest) (map[string]map[string][]*entity.Model, error) {
	result, err := s.repo.GetModelsGroupedByType(ctx, req)
	if err != nil {
		return nil, apierrors.InternalError("查询分组模型失败")
	}
	return result, nil
}

// SaveModel creates a new model
func (s *ModelService) SaveModel(ctx context.Context, req *dto.ModelSaveRequest) (*entity.Model, error) {
	tenantID, deptID, _ := s.repo.GetDefaultTenantAndDept(ctx)
	if req.TenantID == 0 {
		req.TenantID = tenantID
	}
	if req.DeptID == 0 {
		req.DeptID = deptID
	}

	model := &entity.Model{
		ID:                  snowflake.MustGenerateID(),
		DeptID:              req.DeptID,
		TenantID:            req.TenantID,
		ProviderID:          req.ProviderID,
		Title:               req.Title,
		Icon:                req.Icon,
		Description:         req.Description,
		Endpoint:            req.Endpoint,
		RequestPath:         req.RequestPath,
		ModelName:           req.ModelName,
		APIKey:              req.APIKey,
		ExtraConfig:         req.ExtraConfig,
		Options:             req.Options,
		GroupName:           req.GroupName,
		ModelType:           req.ModelType,
		WithUsed:            req.WithUsed,
		SupportThinking:     req.SupportThinking,
		SupportTool:         req.SupportTool,
		SupportImage:        req.SupportImage,
		SupportImageB64Only: req.SupportImageB64Only,
		SupportVideo:        req.SupportVideo,
		SupportAudio:        req.SupportAudio,
		SupportFree:         req.SupportFree,
	}

	if err := s.repo.CreateModel(ctx, model); err != nil {
		return nil, apierrors.InternalError("创建模型失败")
	}
	return model, nil
}

// UpdateModel updates an existing model
func (s *ModelService) UpdateModel(ctx context.Context, req *dto.ModelSaveRequest) (*entity.Model, error) {
	if req.ID == 0 {
		return nil, apierrors.BadRequest("模型ID不能为空")
	}

	existing, err := s.repo.GetModelByID(ctx, req.ID)
	if err != nil {
		return nil, apierrors.InternalError("查询模型失败")
	}
	if existing == nil {
		return nil, apierrors.NotFound("模型不存在")
	}

	existing.DeptID = req.DeptID
	existing.TenantID = req.TenantID
	existing.ProviderID = req.ProviderID
	existing.Title = req.Title
	existing.Icon = req.Icon
	existing.Description = req.Description
	existing.Endpoint = req.Endpoint
	existing.RequestPath = req.RequestPath
	existing.ModelName = req.ModelName
	existing.APIKey = req.APIKey
	existing.ExtraConfig = req.ExtraConfig
	existing.Options = req.Options
	existing.GroupName = req.GroupName
	existing.ModelType = req.ModelType
	existing.WithUsed = req.WithUsed
	existing.SupportThinking = req.SupportThinking
	existing.SupportTool = req.SupportTool
	existing.SupportImage = req.SupportImage
	existing.SupportImageB64Only = req.SupportImageB64Only
	existing.SupportVideo = req.SupportVideo
	existing.SupportAudio = req.SupportAudio
	existing.SupportFree = req.SupportFree

	if err := s.repo.UpdateModel(ctx, existing); err != nil {
		return nil, apierrors.InternalError("更新模型失败")
	}
	return existing, nil
}

// RemoveModel deletes a model
func (s *ModelService) RemoveModel(ctx context.Context, id int64) error {
	existing, err := s.repo.GetModelByID(ctx, id)
	if err != nil {
		return apierrors.InternalError("查询模型失败")
	}
	if existing == nil {
		return apierrors.NotFound("模型不存在")
	}

	if err := s.repo.DeleteModel(ctx, id); err != nil {
		return apierrors.InternalError("删除模型失败")
	}
	return nil
}

// AddAllLlm batch adds models for a provider
func (s *ModelService) AddAllLlm(ctx context.Context, req *dto.AddAllLlmRequest) error {
	tenantID, deptID, _ := s.repo.GetDefaultTenantAndDept(ctx)

	for _, modelReq := range req.Models {
		model := &entity.Model{
			ID:                  snowflake.MustGenerateID(),
			DeptID:              deptID,
			TenantID:            tenantID,
			ProviderID:          req.ProviderID,
			Title:               modelReq.Title,
			Icon:                modelReq.Icon,
			Description:         modelReq.Description,
			Endpoint:            modelReq.Endpoint,
			RequestPath:         modelReq.RequestPath,
			ModelName:           modelReq.ModelName,
			APIKey:              modelReq.APIKey,
			ExtraConfig:         modelReq.ExtraConfig,
			Options:             modelReq.Options,
			GroupName:           req.GroupName,
			ModelType:           req.ModelType,
			WithUsed:            modelReq.WithUsed,
			SupportThinking:     modelReq.SupportThinking,
			SupportTool:         modelReq.SupportTool,
			SupportImage:        modelReq.SupportImage,
			SupportImageB64Only: modelReq.SupportImageB64Only,
			SupportVideo:        modelReq.SupportVideo,
			SupportAudio:        modelReq.SupportAudio,
			SupportFree:         modelReq.SupportFree,
		}
		if err := s.repo.CreateModel(ctx, model); err != nil {
			return apierrors.InternalError("批量创建模型失败")
		}
	}
	return nil
}

// UpdateByEntity updates models by conditions
func (s *ModelService) UpdateByEntity(ctx context.Context, req *dto.UpdateByEntityRequest) error {
	if err := s.repo.UpdateModelsByCondition(ctx, req); err != nil {
		return apierrors.InternalError("批量更新模型失败")
	}
	return nil
}

// RemoveByEntity removes models by conditions
func (s *ModelService) RemoveByEntity(ctx context.Context, req *dto.RemoveByEntityRequest) error {
	if err := s.repo.DeleteModelsByCondition(ctx, req); err != nil {
		return apierrors.InternalError("批量删除模型失败")
	}
	return nil
}

// RemoveLlmByIds removes models by IDs
func (s *ModelService) RemoveLlmByIds(ctx context.Context, ids []int64) error {
	if err := s.repo.DeleteModelsByIDs(ctx, ids); err != nil {
		return apierrors.InternalError("批量删除模型失败")
	}
	return nil
}

// GetModelInstance gets a model with inherited provider config
func (s *ModelService) GetModelInstance(ctx context.Context, id int64) (*entity.Model, error) {
	model, err := s.repo.GetModelInstance(ctx, id)
	if err != nil {
		return nil, apierrors.InternalError("查询模型实例失败")
	}
	if model == nil {
		return nil, apierrors.NotFound("模型不存在")
	}
	return model, nil
}
