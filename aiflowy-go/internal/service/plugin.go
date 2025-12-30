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

// PluginService 插件服务
type PluginService struct {
	repo *repository.PluginRepository
}

// NewPluginService 创建 PluginService
func NewPluginService() *PluginService {
	return &PluginService{
		repo: repository.NewPluginRepository(),
	}
}

// ========================== Plugin ==========================

// GetPlugin 获取插件
func (s *PluginService) GetPlugin(ctx context.Context, id string) (*entity.Plugin, error) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, apierrors.BadRequest("无效的插件 ID")
	}

	plugin, err := s.repo.GetPluginByID(ctx, idInt)
	if err != nil {
		return nil, apierrors.InternalError("获取插件失败")
	}
	if plugin == nil {
		return nil, apierrors.NotFound("插件不存在")
	}
	return plugin, nil
}

// ListPlugins 获取插件列表
func (s *PluginService) ListPlugins(ctx context.Context, tenantID int64) ([]*entity.Plugin, error) {
	return s.repo.ListPlugins(ctx, tenantID)
}

// ListPluginsWithTools 获取插件列表(包含工具)
func (s *PluginService) ListPluginsWithTools(ctx context.Context, tenantID int64) ([]*entity.Plugin, error) {
	return s.repo.ListPluginsWithTools(ctx, tenantID)
}

// SavePlugin 保存插件 (创建或更新)
func (s *PluginService) SavePlugin(ctx context.Context, req *dto.PluginSaveRequest, tenantID, userID, deptID int64) (*entity.Plugin, error) {
	var plugin *entity.Plugin
	var isNew bool

	if req.ID != "" {
		// 更新
		idInt, err := strconv.ParseInt(req.ID, 10, 64)
		if err != nil {
			return nil, apierrors.BadRequest("无效的插件 ID")
		}
		plugin, err = s.repo.GetPluginByID(ctx, idInt)
		if err != nil {
			return nil, apierrors.InternalError("获取插件失败")
		}
		if plugin == nil {
			return nil, apierrors.NotFound("插件不存在")
		}
	} else {
		// 创建
		isNew = true
		plugin = &entity.Plugin{
			ID:        snowflake.MustGenerateID(),
			Created:   time.Now(),
			TenantID:  tenantID,
			CreatedBy: userID,
			DeptID:    deptID,
		}
	}

	// 更新字段
	plugin.Alias = req.Alias
	plugin.Name = req.Name
	plugin.Description = req.Description
	plugin.Type = req.Type
	plugin.BaseURL = req.BaseURL
	plugin.AuthType = req.AuthType
	plugin.Icon = req.Icon
	plugin.Position = req.Position
	plugin.Headers = req.Headers
	plugin.TokenKey = req.TokenKey
	plugin.TokenValue = req.TokenValue

	if isNew {
		if err := s.repo.CreatePlugin(ctx, plugin); err != nil {
			return nil, apierrors.InternalError("创建插件失败")
		}
	} else {
		if err := s.repo.UpdatePlugin(ctx, plugin); err != nil {
			return nil, apierrors.InternalError("更新插件失败")
		}
	}

	return plugin, nil
}

// DeletePlugin 删除插件
func (s *PluginService) DeletePlugin(ctx context.Context, id string) error {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的插件 ID")
	}

	// 检查是否有关联的工具
	items, err := s.repo.ListPluginItemsByPluginID(ctx, idInt)
	if err != nil {
		return apierrors.InternalError("检查插件工具失败")
	}
	if len(items) > 0 {
		return apierrors.BadRequest("请先删除插件下的工具")
	}

	return s.repo.DeletePlugin(ctx, idInt)
}

// ========================== PluginItem ==========================

// GetPluginItem 获取插件工具
func (s *PluginService) GetPluginItem(ctx context.Context, id string) (*entity.PluginItem, error) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, apierrors.BadRequest("无效的插件工具 ID")
	}

	item, err := s.repo.GetPluginItemByID(ctx, idInt)
	if err != nil {
		return nil, apierrors.InternalError("获取插件工具失败")
	}
	if item == nil {
		return nil, apierrors.NotFound("插件工具不存在")
	}
	return item, nil
}

// ListPluginItems 获取插件下的工具列表
func (s *PluginService) ListPluginItems(ctx context.Context, pluginID string, botID string) ([]*entity.PluginItem, error) {
	pluginIDInt, err := strconv.ParseInt(pluginID, 10, 64)
	if err != nil {
		return nil, apierrors.BadRequest("无效的插件 ID")
	}

	items, err := s.repo.ListPluginItemsByPluginID(ctx, pluginIDInt)
	if err != nil {
		return nil, apierrors.InternalError("获取插件工具列表失败")
	}

	// 如果提供了 botID，标记已关联的工具
	if botID != "" {
		botIDInt, err := strconv.ParseInt(botID, 10, 64)
		if err == nil {
			botPluginIDs, _ := s.repo.GetBotPluginIDs(ctx, botIDInt)
			idSet := make(map[int64]bool)
			for _, id := range botPluginIDs {
				idSet[id] = true
			}
			for _, item := range items {
				item.JoinBot = idSet[item.ID]
			}
		}
	}

	return items, nil
}

// ListPluginItemsByBotID 获取 Bot 关联的插件工具列表
func (s *PluginService) ListPluginItemsByBotID(ctx context.Context, botID int64) ([]*entity.PluginItem, error) {
	return s.repo.ListPluginItemsByBotID(ctx, botID)
}

// SavePluginItem 保存插件工具
func (s *PluginService) SavePluginItem(ctx context.Context, req *dto.PluginItemSaveRequest) (*entity.PluginItem, error) {
	pluginIDInt, err := strconv.ParseInt(req.PluginID, 10, 64)
	if err != nil {
		return nil, apierrors.BadRequest("无效的插件 ID")
	}

	var item *entity.PluginItem
	var isNew bool

	if req.ID != "" {
		// 更新
		idInt, err := strconv.ParseInt(req.ID, 10, 64)
		if err != nil {
			return nil, apierrors.BadRequest("无效的插件工具 ID")
		}
		item, err = s.repo.GetPluginItemByID(ctx, idInt)
		if err != nil {
			return nil, apierrors.InternalError("获取插件工具失败")
		}
		if item == nil {
			return nil, apierrors.NotFound("插件工具不存在")
		}
	} else {
		// 创建
		isNew = true
		item = &entity.PluginItem{
			ID:       snowflake.MustGenerateID(),
			PluginID: pluginIDInt,
			Created:  time.Now(),
		}
	}

	// 更新字段
	item.Name = req.Name
	item.Description = req.Description
	item.BasePath = req.BasePath
	item.Status = req.Status
	item.InputData = req.InputData
	item.OutputData = req.OutputData
	item.RequestMethod = req.RequestMethod
	item.ServiceStatus = req.ServiceStatus
	item.DebugStatus = req.DebugStatus
	item.EnglishName = req.EnglishName

	if isNew {
		if err := s.repo.CreatePluginItem(ctx, item); err != nil {
			return nil, apierrors.InternalError("创建插件工具失败")
		}
	} else {
		if err := s.repo.UpdatePluginItem(ctx, item); err != nil {
			return nil, apierrors.InternalError("更新插件工具失败")
		}
	}

	return item, nil
}

// DeletePluginItem 删除插件工具
func (s *PluginService) DeletePluginItem(ctx context.Context, id string) error {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的插件工具 ID")
	}

	// 检查是否有 Bot 关联
	exists, err := s.repo.ExistsBotPlugin(ctx, idInt)
	if err != nil {
		return apierrors.InternalError("检查关联失败")
	}
	if exists {
		return apierrors.BadRequest("此工具还关联着 Bot，请先取消关联")
	}

	return s.repo.DeletePluginItem(ctx, idInt)
}

// ========================== BotPlugin ==========================

// GetBotPluginIDs 获取 Bot 关联的插件工具 ID 列表
func (s *PluginService) GetBotPluginIDs(ctx context.Context, botID string) ([]string, error) {
	botIDInt, err := strconv.ParseInt(botID, 10, 64)
	if err != nil {
		return nil, apierrors.BadRequest("无效的 Bot ID")
	}

	ids, err := s.repo.GetBotPluginIDs(ctx, botIDInt)
	if err != nil {
		return nil, apierrors.InternalError("获取关联失败")
	}

	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = strconv.FormatInt(id, 10)
	}
	return result, nil
}

// UpdateBotPlugins 更新 Bot-插件关联
func (s *PluginService) UpdateBotPlugins(ctx context.Context, req *dto.BotPluginUpdateRequest) error {
	botIDInt, err := strconv.ParseInt(req.BotID, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的 Bot ID")
	}

	var pluginItemIDs []int64
	for _, idStr := range req.PluginToolIDs {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			continue
		}
		pluginItemIDs = append(pluginItemIDs, id)
	}

	return s.repo.SaveBotPlugins(ctx, botIDInt, pluginItemIDs)
}

// DeleteBotPlugin 删除单个 Bot-插件关联
func (s *PluginService) DeleteBotPlugin(ctx context.Context, botID, pluginItemID string) error {
	botIDInt, err := strconv.ParseInt(botID, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的 Bot ID")
	}

	pluginItemIDInt, err := strconv.ParseInt(pluginItemID, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的插件工具 ID")
	}

	return s.repo.DeleteBotPlugin(ctx, botIDInt, pluginItemIDInt)
}

// ========================== PluginCategory ==========================

// ListPluginCategories 获取插件分类列表
func (s *PluginService) ListPluginCategories(ctx context.Context) ([]*entity.PluginCategory, error) {
	return s.repo.ListPluginCategories(ctx)
}

// SavePluginCategory 保存插件分类
func (s *PluginService) SavePluginCategory(ctx context.Context, req *dto.PluginCategorySaveRequest) (*entity.PluginCategory, error) {
	category := &entity.PluginCategory{
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	if req.ID != "" {
		idInt, err := strconv.ParseInt(req.ID, 10, 64)
		if err != nil {
			return nil, apierrors.BadRequest("无效的分类 ID")
		}
		category.ID = idInt
	}

	if err := s.repo.CreatePluginCategory(ctx, category); err != nil {
		return nil, apierrors.InternalError("保存分类失败")
	}

	return category, nil
}

// DeletePluginCategory 删除插件分类
func (s *PluginService) DeletePluginCategory(ctx context.Context, id string) error {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的分类 ID")
	}

	return s.repo.DeletePluginCategory(ctx, idInt)
}
