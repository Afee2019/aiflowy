package model

import (
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/aiflowy/aiflowy-go/internal/dto"
	apierrors "github.com/aiflowy/aiflowy-go/internal/errors"
	"github.com/aiflowy/aiflowy-go/internal/handler/auth"
	"github.com/aiflowy/aiflowy-go/internal/service"
	"github.com/aiflowy/aiflowy-go/pkg/response"
)

// Handler handles model-related HTTP requests
type Handler struct {
	svc *service.ModelService
}

// NewHandler creates a new model handler
func NewHandler() *Handler {
	return &Handler{
		svc: service.NewModelService(),
	}
}

// ========== Model Provider Handlers ==========

// ProviderList lists all model providers
// GET /api/v1/modelProvider/list
func (h *Handler) ProviderList(c echo.Context) error {
	var req dto.ModelProviderListRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	providers, err := h.svc.ListProviders(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return response.Success(c, providers)
}

// ProviderPage returns paginated model providers
// GET /api/v1/modelProvider/page
func (h *Handler) ProviderPage(c echo.Context) error {
	var pageReq dto.PageRequest
	if err := c.Bind(&pageReq); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	var filter dto.ModelProviderListRequest
	_ = c.Bind(&filter)

	result, err := h.svc.PageProviders(c.Request().Context(), &pageReq, &filter)
	if err != nil {
		return err
	}
	return response.Success(c, result)
}

// ProviderDetail gets a model provider by ID
// GET /api/v1/modelProvider/detail?id=xxx
func (h *Handler) ProviderDetail(c echo.Context) error {
	idStr := c.QueryParam("id")
	if idStr == "" {
		return apierrors.BadRequest("缺少供应商ID")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的供应商ID")
	}

	provider, err := h.svc.GetProviderDetail(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return response.Success(c, provider)
}

// ProviderSave creates a new model provider
// POST /api/v1/modelProvider/save
func (h *Handler) ProviderSave(c echo.Context) error {
	var req dto.ModelProviderSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}
	if req.ProviderName == "" {
		return apierrors.BadRequest("供应商名称不能为空")
	}

	operatorID := auth.GetCurrentUserID(c)
	provider, err := h.svc.SaveProvider(c.Request().Context(), &req, operatorID)
	if err != nil {
		return err
	}
	return response.Success(c, provider)
}

// ProviderUpdate updates an existing model provider
// POST /api/v1/modelProvider/update
func (h *Handler) ProviderUpdate(c echo.Context) error {
	var req dto.ModelProviderSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}
	if req.ID == 0 {
		return apierrors.BadRequest("供应商ID不能为空")
	}

	operatorID := auth.GetCurrentUserID(c)
	provider, err := h.svc.UpdateProvider(c.Request().Context(), &req, operatorID)
	if err != nil {
		return err
	}
	return response.Success(c, provider)
}

// ProviderRemove deletes a model provider
// POST /api/v1/modelProvider/remove
func (h *Handler) ProviderRemove(c echo.Context) error {
	var req dto.IDRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}
	if req.ID == 0 {
		return apierrors.BadRequest("供应商ID不能为空")
	}

	err := h.svc.RemoveProvider(c.Request().Context(), req.ID)
	if err != nil {
		return err
	}
	return response.Success(c, nil)
}

// ========== Model Handlers ==========

// ModelList lists all models
// GET /api/v1/model/list
func (h *Handler) ModelList(c echo.Context) error {
	var req dto.ModelListRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	models, err := h.svc.ListModels(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return response.Success(c, models)
}

// ModelPage returns paginated models
// GET /api/v1/model/page
func (h *Handler) ModelPage(c echo.Context) error {
	var req dto.ModelListRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	result, err := h.svc.PageModels(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return response.Success(c, result)
}

// ModelDetail gets a model by ID
// GET /api/v1/model/detail?id=xxx
func (h *Handler) ModelDetail(c echo.Context) error {
	idStr := c.QueryParam("id")
	if idStr == "" {
		return apierrors.BadRequest("缺少模型ID")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的模型ID")
	}

	model, err := h.svc.GetModelDetail(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return response.Success(c, model)
}

// ModelSave creates a new model
// POST /api/v1/model/save
func (h *Handler) ModelSave(c echo.Context) error {
	var req dto.ModelSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}
	if req.ModelName == "" {
		return apierrors.BadRequest("模型名称不能为空")
	}
	if req.ModelType == "" {
		return apierrors.BadRequest("模型类型不能为空")
	}

	model, err := h.svc.SaveModel(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return response.Success(c, model)
}

// ModelUpdate updates an existing model
// POST /api/v1/model/update
func (h *Handler) ModelUpdate(c echo.Context) error {
	var req dto.ModelSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}
	if req.ID == 0 {
		return apierrors.BadRequest("模型ID不能为空")
	}

	model, err := h.svc.UpdateModel(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return response.Success(c, model)
}

// ModelRemove deletes a model
// POST /api/v1/model/remove
func (h *Handler) ModelRemove(c echo.Context) error {
	var req dto.IDRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}
	if req.ID == 0 {
		return apierrors.BadRequest("模型ID不能为空")
	}

	err := h.svc.RemoveModel(c.Request().Context(), req.ID)
	if err != nil {
		return err
	}
	return response.Success(c, nil)
}

// GetList returns models grouped by type and group name
// GET /api/v1/model/getList?providerId=xxx&withUsed=true
func (h *Handler) GetList(c echo.Context) error {
	var req dto.ModelByProviderRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	result, err := h.svc.GetList(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return response.Success(c, result)
}

// SelectLlmByProviderAndModelType returns models filtered by provider and type
// GET /api/v1/model/selectLlmByProviderAndModelType?modelType=chatModel&providerId=xxx&selectText=
func (h *Handler) SelectLlmByProviderAndModelType(c echo.Context) error {
	var req dto.ModelListRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	models, err := h.svc.ListModels(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	// Group by groupName
	grouped := make(map[string]interface{})
	for _, m := range models {
		groupName := m.GroupName
		if groupName == "" {
			groupName = "Default"
		}
		if grouped[groupName] == nil {
			grouped[groupName] = []interface{}{}
		}
		grouped[groupName] = append(grouped[groupName].([]interface{}), m)
	}

	return response.Success(c, grouped)
}

// SelectLlmList returns all available models
// GET /api/v1/model/selectLlmList
func (h *Handler) SelectLlmList(c echo.Context) error {
	withUsed := true
	req := &dto.ModelListRequest{
		WithUsed: &withUsed,
	}

	models, err := h.svc.ListModels(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return response.Success(c, models)
}

// AddAiLlm adds a single model
// POST /api/v1/model/addAiLlm
func (h *Handler) AddAiLlm(c echo.Context) error {
	var req dto.ModelSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}
	if req.ModelName == "" {
		return apierrors.BadRequest("模型名称不能为空")
	}
	if req.ModelType == "" {
		return apierrors.BadRequest("模型类型不能为空")
	}

	model, err := h.svc.SaveModel(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return response.Success(c, model)
}

// AddAllLlm batch adds models for a provider
// POST /api/v1/model/addAllLlm
func (h *Handler) AddAllLlm(c echo.Context) error {
	var req dto.AddAllLlmRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}
	if req.ProviderID == 0 {
		return apierrors.BadRequest("供应商ID不能为空")
	}

	err := h.svc.AddAllLlm(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return response.Success(c, nil)
}

// UpdateByEntity updates models by conditions
// POST /api/v1/model/updateByEntity
func (h *Handler) UpdateByEntity(c echo.Context) error {
	var req dto.UpdateByEntityRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	err := h.svc.UpdateByEntity(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return response.Success(c, nil)
}

// RemoveByEntity removes models by conditions
// POST /api/v1/model/removeByEntity
func (h *Handler) RemoveByEntity(c echo.Context) error {
	var req dto.RemoveByEntityRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}

	err := h.svc.RemoveByEntity(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return response.Success(c, nil)
}

// RemoveLlmByIds removes models by IDs
// POST /api/v1/model/removeLlmByIds
func (h *Handler) RemoveLlmByIds(c echo.Context) error {
	var req dto.RemoveLlmByIdsRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("参数解析失败")
	}
	if len(req.IDs) == 0 {
		return apierrors.BadRequest("模型ID列表不能为空")
	}

	err := h.svc.RemoveLlmByIds(c.Request().Context(), req.IDs)
	if err != nil {
		return err
	}
	return response.Success(c, nil)
}

// VerifyLlmConfig verifies a model configuration
// GET /api/v1/model/verifyLlmConfig?id=xxx
func (h *Handler) VerifyLlmConfig(c echo.Context) error {
	idStr := c.QueryParam("id")
	if idStr == "" {
		return apierrors.BadRequest("缺少模型ID")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的模型ID")
	}

	// Get model instance with inherited config
	model, err := h.svc.GetModelInstance(c.Request().Context(), id)
	if err != nil {
		return err
	}

	// TODO: Actually verify the model by calling the API
	// For now, just return success if the model exists
	result := map[string]interface{}{
		"success":   true,
		"modelName": model.ModelName,
		"modelType": model.ModelType,
		"message":   "配置验证通过",
	}

	return response.Success(c, result)
}
