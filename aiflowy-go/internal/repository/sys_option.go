package repository

import (
	"context"
	"database/sql"

	"github.com/aiflowy/aiflowy-go/internal/entity"
)

// SysOptionRepository 系统配置数据访问层
type SysOptionRepository struct {
	db *sql.DB
}

// NewSysOptionRepository 创建 SysOptionRepository
func NewSysOptionRepository() *SysOptionRepository {
	return &SysOptionRepository{db: GetDB()}
}

// Get 获取配置
func (r *SysOptionRepository) Get(ctx context.Context, tenantID int64, key string) (*entity.SysOption, error) {
	query := `SELECT tenant_id, ` + "`key`" + `, value FROM tb_sys_option WHERE tenant_id = ? AND ` + "`key`" + ` = ?`
	row := r.db.QueryRowContext(ctx, query, tenantID, key)

	var opt entity.SysOption
	var value sql.NullString
	err := row.Scan(&opt.TenantID, &opt.Key, &value)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if value.Valid {
		opt.Value = value.String
	}
	return &opt, nil
}

// Set 设置配置 (存在则更新,不存在则创建)
func (r *SysOptionRepository) Set(ctx context.Context, tenantID int64, key, value string) error {
	// 先尝试更新
	updateQuery := `UPDATE tb_sys_option SET value = ? WHERE tenant_id = ? AND ` + "`key`" + ` = ?`
	result, err := r.db.ExecContext(ctx, updateQuery, value, tenantID, key)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		return nil
	}

	// 不存在则插入
	insertQuery := `INSERT INTO tb_sys_option (tenant_id, ` + "`key`" + `, value) VALUES (?, ?, ?)`
	_, err = r.db.ExecContext(ctx, insertQuery, tenantID, key, value)
	return err
}

// ListByKeys 根据 keys 列表获取配置
func (r *SysOptionRepository) ListByKeys(ctx context.Context, tenantID int64, keys []string) ([]*entity.SysOption, error) {
	if len(keys) == 0 {
		return nil, nil
	}

	query := `SELECT tenant_id, ` + "`key`" + `, value FROM tb_sys_option WHERE tenant_id = ? AND ` + "`key`" + ` IN (`
	args := []interface{}{tenantID}
	for i, key := range keys {
		if i > 0 {
			query += ","
		}
		query += "?"
		args = append(args, key)
	}
	query += ")"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.SysOption
	for rows.Next() {
		var opt entity.SysOption
		var value sql.NullString
		if err := rows.Scan(&opt.TenantID, &opt.Key, &value); err != nil {
			continue
		}
		if value.Valid {
			opt.Value = value.String
		}
		list = append(list, &opt)
	}
	return list, nil
}

// Delete 删除配置
func (r *SysOptionRepository) Delete(ctx context.Context, tenantID int64, key string) error {
	query := `DELETE FROM tb_sys_option WHERE tenant_id = ? AND ` + "`key`" + ` = ?`
	_, err := r.db.ExecContext(ctx, query, tenantID, key)
	return err
}
