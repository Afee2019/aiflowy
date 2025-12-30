package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/internal/repository"
	"github.com/google/uuid"
)

// SysApiKeyService 系统 API 密钥服务
type SysApiKeyService struct {
	repo *repository.SysApiKeyRepository
}

// NewSysApiKeyService 创建 SysApiKeyService
func NewSysApiKeyService() *SysApiKeyService {
	return &SysApiKeyService{
		repo: repository.NewSysApiKeyRepository(),
	}
}

// Generate 生成新的 API 密钥
func (s *SysApiKeyService) Generate(ctx context.Context, userID, tenantID, deptID int64) (*entity.SysApiKey, error) {
	// 生成 UUID 作为 API Key
	apiKey := uuid.New().String()

	// 默认30天有效期
	expiredAt := time.Now().AddDate(0, 0, 30)

	key := &entity.SysApiKey{
		ApiKey:    apiKey,
		Status:    1, // 启用
		ExpiredAt: &expiredAt,
		CreatedBy: &userID,
		TenantID:  &tenantID,
		DeptID:    &deptID,
	}

	if err := s.repo.Create(ctx, key); err != nil {
		return nil, fmt.Errorf("创建 API 密钥失败: %w", err)
	}

	return key, nil
}

// Update 更新 API 密钥
func (s *SysApiKeyService) Update(ctx context.Context, apiKey *entity.SysApiKey) error {
	// 更新权限
	if apiKey.PermissionIds != nil {
		if err := s.repo.UpdatePermissions(ctx, apiKey.ID, apiKey.PermissionIds); err != nil {
			return fmt.Errorf("更新权限失败: %w", err)
		}
	}

	return s.repo.Update(ctx, apiKey)
}

// Delete 删除 API 密钥
func (s *SysApiKeyService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// GetByID 根据 ID 获取
func (s *SysApiKeyService) GetByID(ctx context.Context, id int64) (*entity.SysApiKey, error) {
	key, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if key != nil {
		// 获取关联的权限
		ids, _ := s.repo.GetPermissionIDs(ctx, key.ID)
		key.PermissionIds = ids
	}
	return key, nil
}

// Page 分页查询
func (s *SysApiKeyService) Page(ctx context.Context, pageNum, pageSize int) ([]*entity.SysApiKey, int64, error) {
	list, total, err := s.repo.Page(ctx, pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 获取每个密钥的权限
	for _, key := range list {
		ids, _ := s.repo.GetPermissionIDs(ctx, key.ID)
		key.PermissionIds = ids
	}

	return list, total, nil
}

// CheckApiKey 验证 API 密钥
func (s *SysApiKeyService) CheckApiKey(ctx context.Context, apiKey string) error {
	key, err := s.repo.GetByApiKey(ctx, apiKey)
	if err != nil {
		return fmt.Errorf("查询 API 密钥失败: %w", err)
	}
	if key == nil {
		return fmt.Errorf("API 密钥不存在")
	}
	if key.Status == 0 {
		return fmt.Errorf("API 密钥未启用")
	}
	if key.ExpiredAt != nil && key.ExpiredAt.Before(time.Now()) {
		return fmt.Errorf("API 密钥已过期")
	}
	return nil
}

// CheckApiKeyPermission 验证 API 密钥权限
func (s *SysApiKeyService) CheckApiKeyPermission(ctx context.Context, apiKey, requestURI string) error {
	key, err := s.repo.GetByApiKey(ctx, apiKey)
	if err != nil {
		return fmt.Errorf("查询 API 密钥失败: %w", err)
	}
	if key == nil {
		return fmt.Errorf("API 密钥不存在")
	}
	if key.Status == 0 {
		return fmt.Errorf("API 密钥未启用")
	}
	if key.ExpiredAt != nil && key.ExpiredAt.Before(time.Now()) {
		return fmt.Errorf("API 密钥已过期")
	}

	// 获取密钥的权限列表
	permissionIDs, err := s.repo.GetPermissionIDs(ctx, key.ID)
	if err != nil {
		return fmt.Errorf("获取权限失败: %w", err)
	}

	// 获取所有资源
	resources, err := s.repo.ListResources(ctx)
	if err != nil {
		return fmt.Errorf("获取资源失败: %w", err)
	}

	// 检查权限
	permissionSet := make(map[int64]bool)
	for _, id := range permissionIDs {
		permissionSet[id] = true
	}

	for _, res := range resources {
		if permissionSet[res.ID] && res.RequestInterface == requestURI {
			return nil
		}
	}

	return fmt.Errorf("没有权限访问该资源")
}

// ListResources 获取所有可授权的资源
func (s *SysApiKeyService) ListResources(ctx context.Context) ([]*entity.SysApiKeyResource, error) {
	return s.repo.ListResources(ctx)
}
