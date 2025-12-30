package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/pkg/snowflake"
)

// SysApiKeyRepository 系统 API 密钥数据访问层
type SysApiKeyRepository struct {
	db *sql.DB
}

// NewSysApiKeyRepository 创建 SysApiKeyRepository
func NewSysApiKeyRepository() *SysApiKeyRepository {
	return &SysApiKeyRepository{db: GetDB()}
}

// Create 创建 API 密钥
func (r *SysApiKeyRepository) Create(ctx context.Context, apiKey *entity.SysApiKey) error {
	if apiKey.ID == 0 {
		id, _ := snowflake.GenerateID()
		apiKey.ID = id
	}
	now := time.Now()
	apiKey.Created = &now

	query := `
		INSERT INTO tb_sys_api_key (id, api_key, status, expired_at, created, created_by, dept_id, tenant_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		apiKey.ID, apiKey.ApiKey, apiKey.Status, apiKey.ExpiredAt,
		apiKey.Created, apiKey.CreatedBy, apiKey.DeptID, apiKey.TenantID,
	)
	return err
}

// Update 更新 API 密钥
func (r *SysApiKeyRepository) Update(ctx context.Context, apiKey *entity.SysApiKey) error {
	query := `
		UPDATE tb_sys_api_key
		SET api_key = ?, status = ?, expired_at = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		apiKey.ApiKey, apiKey.Status, apiKey.ExpiredAt, apiKey.ID,
	)
	return err
}

// Delete 删除 API 密钥
func (r *SysApiKeyRepository) Delete(ctx context.Context, id int64) error {
	// 先删除关联
	_, _ = r.db.ExecContext(ctx, "DELETE FROM tb_sys_api_key_resource_mapping WHERE api_key_id = ?", id)
	// 删除主记录
	_, err := r.db.ExecContext(ctx, "DELETE FROM tb_sys_api_key WHERE id = ?", id)
	return err
}

// GetByID 根据 ID 获取
func (r *SysApiKeyRepository) GetByID(ctx context.Context, id int64) (*entity.SysApiKey, error) {
	query := `
		SELECT id, api_key, status, expired_at, created, created_by, dept_id, tenant_id
		FROM tb_sys_api_key WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var apiKey entity.SysApiKey
	var expiredAt, created sql.NullTime
	var createdBy, deptID, tenantID sql.NullInt64

	err := row.Scan(
		&apiKey.ID, &apiKey.ApiKey, &apiKey.Status, &expiredAt,
		&created, &createdBy, &deptID, &tenantID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if expiredAt.Valid {
		apiKey.ExpiredAt = &expiredAt.Time
	}
	if created.Valid {
		apiKey.Created = &created.Time
	}
	if createdBy.Valid {
		apiKey.CreatedBy = &createdBy.Int64
	}
	if deptID.Valid {
		apiKey.DeptID = &deptID.Int64
	}
	if tenantID.Valid {
		apiKey.TenantID = &tenantID.Int64
	}

	return &apiKey, nil
}

// GetByApiKey 根据 API Key 获取
func (r *SysApiKeyRepository) GetByApiKey(ctx context.Context, apiKey string) (*entity.SysApiKey, error) {
	query := `
		SELECT id, api_key, status, expired_at, created, created_by, dept_id, tenant_id
		FROM tb_sys_api_key WHERE api_key = ?
	`
	row := r.db.QueryRowContext(ctx, query, apiKey)

	var key entity.SysApiKey
	var expiredAt, created sql.NullTime
	var createdBy, deptID, tenantID sql.NullInt64

	err := row.Scan(
		&key.ID, &key.ApiKey, &key.Status, &expiredAt,
		&created, &createdBy, &deptID, &tenantID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if expiredAt.Valid {
		key.ExpiredAt = &expiredAt.Time
	}
	if created.Valid {
		key.Created = &created.Time
	}
	if createdBy.Valid {
		key.CreatedBy = &createdBy.Int64
	}
	if deptID.Valid {
		key.DeptID = &deptID.Int64
	}
	if tenantID.Valid {
		key.TenantID = &tenantID.Int64
	}

	return &key, nil
}

// Page 分页查询
func (r *SysApiKeyRepository) Page(ctx context.Context, pageNum, pageSize int) ([]*entity.SysApiKey, int64, error) {
	// 查询总数
	var total int64
	countQuery := `SELECT COUNT(*) FROM tb_sys_api_key`
	r.db.QueryRowContext(ctx, countQuery).Scan(&total)

	// 查询数据
	query := `
		SELECT id, api_key, status, expired_at, created, created_by, dept_id, tenant_id
		FROM tb_sys_api_key
		ORDER BY created DESC
		LIMIT ? OFFSET ?
	`
	offset := (pageNum - 1) * pageSize
	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*entity.SysApiKey
	for rows.Next() {
		var apiKey entity.SysApiKey
		var expiredAt, created sql.NullTime
		var createdBy, deptID, tenantID sql.NullInt64

		err := rows.Scan(
			&apiKey.ID, &apiKey.ApiKey, &apiKey.Status, &expiredAt,
			&created, &createdBy, &deptID, &tenantID,
		)
		if err != nil {
			continue
		}

		if expiredAt.Valid {
			apiKey.ExpiredAt = &expiredAt.Time
		}
		if created.Valid {
			apiKey.Created = &created.Time
		}
		if createdBy.Valid {
			apiKey.CreatedBy = &createdBy.Int64
		}
		if deptID.Valid {
			apiKey.DeptID = &deptID.Int64
		}
		if tenantID.Valid {
			apiKey.TenantID = &tenantID.Int64
		}

		list = append(list, &apiKey)
	}

	return list, total, nil
}

// GetPermissionIDs 获取 API 密钥关联的资源 ID 列表
func (r *SysApiKeyRepository) GetPermissionIDs(ctx context.Context, apiKeyID int64) ([]int64, error) {
	query := `SELECT api_key_resource_id FROM tb_sys_api_key_resource_mapping WHERE api_key_id = ?`
	rows, err := r.db.QueryContext(ctx, query, apiKeyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// UpdatePermissions 更新 API 密钥的资源权限
func (r *SysApiKeyRepository) UpdatePermissions(ctx context.Context, apiKeyID int64, resourceIDs []int64) error {
	// 删除旧的关联
	_, err := r.db.ExecContext(ctx, "DELETE FROM tb_sys_api_key_resource_mapping WHERE api_key_id = ?", apiKeyID)
	if err != nil {
		return err
	}

	// 添加新的关联
	if len(resourceIDs) > 0 {
		query := "INSERT INTO tb_sys_api_key_resource_mapping (id, api_key_id, api_key_resource_id) VALUES (?, ?, ?)"
		for _, resID := range resourceIDs {
			id, _ := snowflake.GenerateID()
			_, err := r.db.ExecContext(ctx, query, id, apiKeyID, resID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ListResources 获取所有资源
func (r *SysApiKeyRepository) ListResources(ctx context.Context) ([]*entity.SysApiKeyResource, error) {
	query := `SELECT id, request_interface, title FROM tb_sys_api_key_resource`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.SysApiKeyResource
	for rows.Next() {
		var res entity.SysApiKeyResource
		if err := rows.Scan(&res.ID, &res.RequestInterface, &res.Title); err != nil {
			continue
		}
		list = append(list, &res)
	}
	return list, nil
}
