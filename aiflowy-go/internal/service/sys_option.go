package service

import (
	"context"

	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/internal/repository"
)

// SysOptionService 系统配置服务
type SysOptionService struct {
	repo *repository.SysOptionRepository
}

// NewSysOptionService 创建 SysOptionService
func NewSysOptionService() *SysOptionService {
	return &SysOptionService{
		repo: repository.NewSysOptionRepository(),
	}
}

// Get 获取配置
func (s *SysOptionService) Get(ctx context.Context, tenantID int64, key string) (string, error) {
	opt, err := s.repo.Get(ctx, tenantID, key)
	if err != nil {
		return "", err
	}
	if opt == nil {
		return "", nil
	}
	return opt.Value, nil
}

// Set 设置配置
func (s *SysOptionService) Set(ctx context.Context, tenantID int64, key, value string) error {
	return s.repo.Set(ctx, tenantID, key, value)
}

// GetMultiple 批量获取配置
func (s *SysOptionService) GetMultiple(ctx context.Context, tenantID int64, keys []string) (map[string]string, error) {
	options, err := s.repo.ListByKeys(ctx, tenantID, keys)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, opt := range options {
		result[opt.Key] = opt.Value
	}
	return result, nil
}

// SetMultiple 批量设置配置
func (s *SysOptionService) SetMultiple(ctx context.Context, tenantID int64, options map[string]string) error {
	for key, value := range options {
		if err := s.repo.Set(ctx, tenantID, key, value); err != nil {
			return err
		}
	}
	return nil
}

// Delete 删除配置
func (s *SysOptionService) Delete(ctx context.Context, tenantID int64, key string) error {
	return s.repo.Delete(ctx, tenantID, key)
}

// ListByKeys 根据 keys 获取配置列表
func (s *SysOptionService) ListByKeys(ctx context.Context, tenantID int64, keys []string) ([]*entity.SysOption, error) {
	return s.repo.ListByKeys(ctx, tenantID, keys)
}
