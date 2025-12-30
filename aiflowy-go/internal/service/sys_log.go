package service

import (
	"context"

	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/internal/repository"
)

// SysLogService 操作日志服务
type SysLogService struct {
	repo *repository.SysLogRepository
}

// NewSysLogService 创建 SysLogService
func NewSysLogService() *SysLogService {
	return &SysLogService{
		repo: repository.NewSysLogRepository(),
	}
}

// Create 创建日志
func (s *SysLogService) Create(ctx context.Context, log *entity.SysLog) error {
	return s.repo.Create(ctx, log)
}

// Page 分页查询
func (s *SysLogService) Page(ctx context.Context, pageNum, pageSize int, actionName, actionType string) ([]*entity.SysLog, int64, error) {
	return s.repo.Page(ctx, pageNum, pageSize, actionName, actionType)
}

// Delete 删除日志
func (s *SysLogService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// RecordAction 记录操作日志
func (s *SysLogService) RecordAction(ctx context.Context, accountID *int64, actionName, actionType, actionURL, actionIP, actionParams string) {
	log := &entity.SysLog{
		AccountID:    accountID,
		ActionName:   actionName,
		ActionType:   actionType,
		ActionURL:    actionURL,
		ActionIP:     actionIP,
		ActionParams: actionParams,
		Status:       1,
	}
	_ = s.repo.Create(ctx, log)
}
