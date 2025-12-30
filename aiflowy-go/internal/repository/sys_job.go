package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/pkg/snowflake"
)

// SysJobRepository 定时任务数据访问层
type SysJobRepository struct {
	db *sql.DB
}

// NewSysJobRepository 创建 SysJobRepository
func NewSysJobRepository() *SysJobRepository {
	return &SysJobRepository{db: GetDB()}
}

// Create 创建任务
func (r *SysJobRepository) Create(ctx context.Context, job *entity.SysJob) error {
	if job.ID == 0 {
		id, _ := snowflake.GenerateID()
		job.ID = id
	}
	now := time.Now()
	job.Created = &now
	job.Modified = &now

	query := `
		INSERT INTO tb_sys_job (id, dept_id, tenant_id, job_name, job_type, job_params,
			cron_expression, allow_concurrent, misfire_policy, options, status,
			created, created_by, modified, modified_by, remark)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		job.ID, job.DeptID, job.TenantID, job.JobName, job.JobType, job.JobParams,
		job.CronExpression, job.AllowConcurrent, job.MisfirePolicy, job.Options, job.Status,
		job.Created, job.CreatedBy, job.Modified, job.ModifiedBy, job.Remark,
	)
	return err
}

// Update 更新任务
func (r *SysJobRepository) Update(ctx context.Context, job *entity.SysJob) error {
	now := time.Now()
	job.Modified = &now

	query := `
		UPDATE tb_sys_job
		SET job_name = ?, job_type = ?, job_params = ?, cron_expression = ?,
			allow_concurrent = ?, misfire_policy = ?, options = ?, status = ?,
			modified = ?, modified_by = ?, remark = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		job.JobName, job.JobType, job.JobParams, job.CronExpression,
		job.AllowConcurrent, job.MisfirePolicy, job.Options, job.Status,
		job.Modified, job.ModifiedBy, job.Remark, job.ID,
	)
	return err
}

// UpdateStatus 更新状态
func (r *SysJobRepository) UpdateStatus(ctx context.Context, id int64, status int) error {
	query := `UPDATE tb_sys_job SET status = ?, modified = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

// Delete 删除任务
func (r *SysJobRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM tb_sys_job WHERE id = ?", id)
	return err
}

// GetByID 根据 ID 获取
func (r *SysJobRepository) GetByID(ctx context.Context, id int64) (*entity.SysJob, error) {
	query := `
		SELECT id, dept_id, tenant_id, job_name, job_type, job_params, cron_expression,
			allow_concurrent, misfire_policy, options, status, created, created_by, modified, modified_by, remark
		FROM tb_sys_job WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var job entity.SysJob
	var jobParams, options, remark sql.NullString
	var created, modified sql.NullTime
	var createdBy, modifiedBy sql.NullInt64

	err := row.Scan(
		&job.ID, &job.DeptID, &job.TenantID, &job.JobName, &job.JobType, &jobParams,
		&job.CronExpression, &job.AllowConcurrent, &job.MisfirePolicy, &options, &job.Status,
		&created, &createdBy, &modified, &modifiedBy, &remark,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if jobParams.Valid {
		job.JobParams = &jobParams.String
	}
	if options.Valid {
		job.Options = &options.String
	}
	if remark.Valid {
		job.Remark = &remark.String
	}
	if created.Valid {
		job.Created = &created.Time
	}
	if createdBy.Valid {
		job.CreatedBy = &createdBy.Int64
	}
	if modified.Valid {
		job.Modified = &modified.Time
	}
	if modifiedBy.Valid {
		job.ModifiedBy = &modifiedBy.Int64
	}

	return &job, nil
}

// Page 分页查询
func (r *SysJobRepository) Page(ctx context.Context, pageNum, pageSize int, jobName string) ([]*entity.SysJob, int64, error) {
	// 构建查询条件
	countQuery := `SELECT COUNT(*) FROM tb_sys_job WHERE 1=1`
	query := `
		SELECT id, dept_id, tenant_id, job_name, job_type, job_params, cron_expression,
			allow_concurrent, misfire_policy, options, status, created, created_by, modified, modified_by, remark
		FROM tb_sys_job WHERE 1=1
	`
	args := []interface{}{}
	countArgs := []interface{}{}

	if jobName != "" {
		countQuery += " AND job_name LIKE ?"
		query += " AND job_name LIKE ?"
		args = append(args, "%"+jobName+"%")
		countArgs = append(countArgs, "%"+jobName+"%")
	}

	// 查询总数
	var total int64
	r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)

	// 添加排序和分页
	query += " ORDER BY created DESC LIMIT ? OFFSET ?"
	offset := (pageNum - 1) * pageSize
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*entity.SysJob
	for rows.Next() {
		var job entity.SysJob
		var jobParams, options, remark sql.NullString
		var created, modified sql.NullTime
		var createdBy, modifiedBy sql.NullInt64

		err := rows.Scan(
			&job.ID, &job.DeptID, &job.TenantID, &job.JobName, &job.JobType, &jobParams,
			&job.CronExpression, &job.AllowConcurrent, &job.MisfirePolicy, &options, &job.Status,
			&created, &createdBy, &modified, &modifiedBy, &remark,
		)
		if err != nil {
			continue
		}

		if jobParams.Valid {
			job.JobParams = &jobParams.String
		}
		if options.Valid {
			job.Options = &options.String
		}
		if remark.Valid {
			job.Remark = &remark.String
		}
		if created.Valid {
			job.Created = &created.Time
		}
		if createdBy.Valid {
			job.CreatedBy = &createdBy.Int64
		}
		if modified.Valid {
			job.Modified = &modified.Time
		}
		if modifiedBy.Valid {
			job.ModifiedBy = &modifiedBy.Int64
		}

		list = append(list, &job)
	}

	return list, total, nil
}

// ListRunning 获取所有运行中的任务
func (r *SysJobRepository) ListRunning(ctx context.Context) ([]*entity.SysJob, error) {
	query := `
		SELECT id, dept_id, tenant_id, job_name, job_type, job_params, cron_expression,
			allow_concurrent, misfire_policy, options, status
		FROM tb_sys_job WHERE status = 1
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.SysJob
	for rows.Next() {
		var job entity.SysJob
		var jobParams, options sql.NullString

		err := rows.Scan(
			&job.ID, &job.DeptID, &job.TenantID, &job.JobName, &job.JobType, &jobParams,
			&job.CronExpression, &job.AllowConcurrent, &job.MisfirePolicy, &options, &job.Status,
		)
		if err != nil {
			continue
		}

		if jobParams.Valid {
			job.JobParams = &jobParams.String
		}
		if options.Valid {
			job.Options = &options.String
		}

		list = append(list, &job)
	}

	return list, nil
}

// CreateLog 创建任务日志
func (r *SysJobRepository) CreateLog(ctx context.Context, log *entity.SysJobLog) error {
	if log.ID == 0 {
		id, _ := snowflake.GenerateID()
		log.ID = id
	}
	now := time.Now()
	log.Created = &now

	query := `
		INSERT INTO tb_sys_job_log (id, job_id, job_name, job_params, result, error, status, duration, created)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		log.ID, log.JobID, log.JobName, log.JobParams,
		log.Result, log.Error, log.Status, log.Duration, log.Created,
	)
	return err
}
