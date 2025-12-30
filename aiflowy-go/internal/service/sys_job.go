package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/internal/repository"
	"github.com/robfig/cron/v3"
)

// SysJobService 定时任务服务
type SysJobService struct {
	repo      *repository.SysJobRepository
	scheduler *cron.Cron
	entryMap  sync.Map // map[int64]cron.EntryID
}

var (
	jobServiceInstance *SysJobService
	jobServiceOnce     sync.Once
)

// GetSysJobService 获取单例
func GetSysJobService() *SysJobService {
	jobServiceOnce.Do(func() {
		jobServiceInstance = &SysJobService{
			repo:      repository.NewSysJobRepository(),
			scheduler: cron.New(cron.WithSeconds()),
		}
		// 启动调度器
		jobServiceInstance.scheduler.Start()
	})
	return jobServiceInstance
}

// NewSysJobService 创建 SysJobService
func NewSysJobService() *SysJobService {
	return GetSysJobService()
}

// Create 创建任务
func (s *SysJobService) Create(ctx context.Context, job *entity.SysJob) error {
	return s.repo.Create(ctx, job)
}

// Update 更新任务
func (s *SysJobService) Update(ctx context.Context, job *entity.SysJob) error {
	// 如果任务正在运行,先停止
	s.StopJob(job.ID)
	return s.repo.Update(ctx, job)
}

// Delete 删除任务
func (s *SysJobService) Delete(ctx context.Context, id int64) error {
	s.StopJob(id)
	return s.repo.Delete(ctx, id)
}

// GetByID 根据 ID 获取
func (s *SysJobService) GetByID(ctx context.Context, id int64) (*entity.SysJob, error) {
	return s.repo.GetByID(ctx, id)
}

// Page 分页查询
func (s *SysJobService) Page(ctx context.Context, pageNum, pageSize int, jobName string) ([]*entity.SysJob, int64, error) {
	return s.repo.Page(ctx, pageNum, pageSize, jobName)
}

// Start 启动任务
func (s *SysJobService) Start(ctx context.Context, id int64) error {
	job, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if job == nil {
		return fmt.Errorf("任务不存在")
	}

	// 添加到调度器
	entryID, err := s.scheduler.AddFunc(job.CronExpression, func() {
		s.executeJob(context.Background(), job)
	})
	if err != nil {
		return fmt.Errorf("添加任务失败: %w", err)
	}

	// 保存 entryID
	s.entryMap.Store(id, entryID)

	// 更新状态
	return s.repo.UpdateStatus(ctx, id, 1)
}

// Stop 停止任务
func (s *SysJobService) Stop(ctx context.Context, id int64) error {
	s.StopJob(id)
	return s.repo.UpdateStatus(ctx, id, 0)
}

// StopJob 停止任务调度
func (s *SysJobService) StopJob(id int64) {
	if entryID, ok := s.entryMap.Load(id); ok {
		s.scheduler.Remove(entryID.(cron.EntryID))
		s.entryMap.Delete(id)
	}
}

// executeJob 执行任务
func (s *SysJobService) executeJob(ctx context.Context, job *entity.SysJob) {
	startTime := time.Now()
	var result, errMsg string
	status := 1

	// 这里根据任务类型执行不同的逻辑
	// 目前只是记录日志
	result = fmt.Sprintf("任务 %s 执行成功", job.JobName)

	duration := time.Since(startTime).Milliseconds()

	// 记录日志
	log := &entity.SysJobLog{
		JobID:    job.ID,
		JobName:  job.JobName,
		JobParams: job.JobParams,
		Result:   &result,
		Status:   status,
		Duration: duration,
	}
	if errMsg != "" {
		log.Error = &errMsg
		log.Status = 0
	}
	_ = s.repo.CreateLog(ctx, log)
}

// GetNextTimes 获取下次执行时间
func (s *SysJobService) GetNextTimes(cronExpression string, count int) ([]string, error) {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(cronExpression)
	if err != nil {
		return nil, fmt.Errorf("cron 表达式错误: %w", err)
	}

	var times []string
	next := time.Now()
	for i := 0; i < count; i++ {
		next = schedule.Next(next)
		times = append(times, next.Format("2006-01-02 15:04:05"))
	}
	return times, nil
}

// LoadRunningJobs 加载运行中的任务
func (s *SysJobService) LoadRunningJobs(ctx context.Context) error {
	jobs, err := s.repo.ListRunning(ctx)
	if err != nil {
		return err
	}

	for _, job := range jobs {
		entryID, err := s.scheduler.AddFunc(job.CronExpression, func() {
			s.executeJob(context.Background(), job)
		})
		if err != nil {
			continue
		}
		s.entryMap.Store(job.ID, entryID)
	}
	return nil
}

// Shutdown 关闭调度器
func (s *SysJobService) Shutdown() {
	s.scheduler.Stop()
}
