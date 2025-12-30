package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/pkg/snowflake"
)

// PluginRepository 插件数据访问层
type PluginRepository struct{}

// NewPluginRepository 创建 PluginRepository
func NewPluginRepository() *PluginRepository {
	return &PluginRepository{}
}

// ========================== Plugin ==========================

// GetPluginByID 根据 ID 获取插件
func (r *PluginRepository) GetPluginByID(ctx context.Context, id int64) (*entity.Plugin, error) {
	query := `SELECT id, alias, name, description, type, base_url, auth_type, created,
		icon, position, headers, token_key, token_value, dept_id, tenant_id, created_by
		FROM tb_plugin WHERE id = ?`

	var p entity.Plugin
	err := db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Alias, &p.Name, &p.Description, &p.Type, &p.BaseURL, &p.AuthType, &p.Created,
		&p.Icon, &p.Position, &p.Headers, &p.TokenKey, &p.TokenValue, &p.DeptID, &p.TenantID, &p.CreatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ListPlugins 获取插件列表
func (r *PluginRepository) ListPlugins(ctx context.Context, tenantID int64) ([]*entity.Plugin, error) {
	query := `SELECT id, alias, name, description, type, base_url, auth_type, created,
		icon, position, headers, token_key, token_value, dept_id, tenant_id, created_by
		FROM tb_plugin WHERE tenant_id = ? ORDER BY created DESC`

	rows, err := db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plugins []*entity.Plugin
	for rows.Next() {
		var p entity.Plugin
		err := rows.Scan(
			&p.ID, &p.Alias, &p.Name, &p.Description, &p.Type, &p.BaseURL, &p.AuthType, &p.Created,
			&p.Icon, &p.Position, &p.Headers, &p.TokenKey, &p.TokenValue, &p.DeptID, &p.TenantID, &p.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		plugins = append(plugins, &p)
	}
	return plugins, nil
}

// CreatePlugin 创建插件
func (r *PluginRepository) CreatePlugin(ctx context.Context, p *entity.Plugin) error {
	query := `INSERT INTO tb_plugin (id, alias, name, description, type, base_url, auth_type, created,
		icon, position, headers, token_key, token_value, dept_id, tenant_id, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := db.ExecContext(ctx, query,
		p.ID, p.Alias, p.Name, p.Description, p.Type, p.BaseURL, p.AuthType, p.Created,
		p.Icon, p.Position, p.Headers, p.TokenKey, p.TokenValue, p.DeptID, p.TenantID, p.CreatedBy,
	)
	return err
}

// UpdatePlugin 更新插件
func (r *PluginRepository) UpdatePlugin(ctx context.Context, p *entity.Plugin) error {
	query := `UPDATE tb_plugin SET alias=?, name=?, description=?, type=?, base_url=?, auth_type=?,
		icon=?, position=?, headers=?, token_key=?, token_value=? WHERE id=?`

	_, err := db.ExecContext(ctx, query,
		p.Alias, p.Name, p.Description, p.Type, p.BaseURL, p.AuthType,
		p.Icon, p.Position, p.Headers, p.TokenKey, p.TokenValue, p.ID,
	)
	return err
}

// DeletePlugin 删除插件
func (r *PluginRepository) DeletePlugin(ctx context.Context, id int64) error {
	query := `DELETE FROM tb_plugin WHERE id = ?`
	_, err := db.ExecContext(ctx, query, id)
	return err
}

// ========================== PluginItem ==========================

// GetPluginItemByID 根据 ID 获取插件工具
func (r *PluginRepository) GetPluginItemByID(ctx context.Context, id int64) (*entity.PluginItem, error) {
	query := `SELECT id, plugin_id, name, description, base_path, created, status,
		input_data, output_data, request_method, service_status, debug_status, english_name
		FROM tb_plugin_item WHERE id = ?`

	var item entity.PluginItem
	err := db.QueryRowContext(ctx, query, id).Scan(
		&item.ID, &item.PluginID, &item.Name, &item.Description, &item.BasePath, &item.Created, &item.Status,
		&item.InputData, &item.OutputData, &item.RequestMethod, &item.ServiceStatus, &item.DebugStatus, &item.EnglishName,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// ListPluginItemsByPluginID 获取插件下的工具列表
func (r *PluginRepository) ListPluginItemsByPluginID(ctx context.Context, pluginID int64) ([]*entity.PluginItem, error) {
	query := `SELECT id, plugin_id, name, description, base_path, created, status,
		input_data, output_data, request_method, service_status, debug_status, english_name
		FROM tb_plugin_item WHERE plugin_id = ? ORDER BY created DESC`

	rows, err := db.QueryContext(ctx, query, pluginID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*entity.PluginItem
	for rows.Next() {
		var item entity.PluginItem
		err := rows.Scan(
			&item.ID, &item.PluginID, &item.Name, &item.Description, &item.BasePath, &item.Created, &item.Status,
			&item.InputData, &item.OutputData, &item.RequestMethod, &item.ServiceStatus, &item.DebugStatus, &item.EnglishName,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

// ListPluginItemsByBotID 获取 Bot 关联的插件工具列表
func (r *PluginRepository) ListPluginItemsByBotID(ctx context.Context, botID int64) ([]*entity.PluginItem, error) {
	query := `SELECT pi.id, pi.plugin_id, pi.name, pi.description, pi.base_path, pi.created, pi.status,
		pi.input_data, pi.output_data, pi.request_method, pi.service_status, pi.debug_status, pi.english_name
		FROM tb_plugin_item pi
		INNER JOIN tb_bot_plugin bp ON pi.id = bp.plugin_item_id
		WHERE bp.bot_id = ? AND pi.service_status = 1
		ORDER BY pi.created DESC`

	rows, err := db.QueryContext(ctx, query, botID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*entity.PluginItem
	for rows.Next() {
		var item entity.PluginItem
		err := rows.Scan(
			&item.ID, &item.PluginID, &item.Name, &item.Description, &item.BasePath, &item.Created, &item.Status,
			&item.InputData, &item.OutputData, &item.RequestMethod, &item.ServiceStatus, &item.DebugStatus, &item.EnglishName,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

// CreatePluginItem 创建插件工具
func (r *PluginRepository) CreatePluginItem(ctx context.Context, item *entity.PluginItem) error {
	query := `INSERT INTO tb_plugin_item (id, plugin_id, name, description, base_path, created, status,
		input_data, output_data, request_method, service_status, debug_status, english_name)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := db.ExecContext(ctx, query,
		item.ID, item.PluginID, item.Name, item.Description, item.BasePath, item.Created, item.Status,
		item.InputData, item.OutputData, item.RequestMethod, item.ServiceStatus, item.DebugStatus, item.EnglishName,
	)
	return err
}

// UpdatePluginItem 更新插件工具
func (r *PluginRepository) UpdatePluginItem(ctx context.Context, item *entity.PluginItem) error {
	query := `UPDATE tb_plugin_item SET name=?, description=?, base_path=?, status=?,
		input_data=?, output_data=?, request_method=?, service_status=?, debug_status=?, english_name=?
		WHERE id=?`

	_, err := db.ExecContext(ctx, query,
		item.Name, item.Description, item.BasePath, item.Status,
		item.InputData, item.OutputData, item.RequestMethod, item.ServiceStatus, item.DebugStatus, item.EnglishName,
		item.ID,
	)
	return err
}

// DeletePluginItem 删除插件工具
func (r *PluginRepository) DeletePluginItem(ctx context.Context, id int64) error {
	query := `DELETE FROM tb_plugin_item WHERE id = ?`
	_, err := db.ExecContext(ctx, query, id)
	return err
}

// ========================== BotPlugin ==========================

// GetBotPluginIDs 获取 Bot 关联的插件工具 ID 列表
func (r *PluginRepository) GetBotPluginIDs(ctx context.Context, botID int64) ([]int64, error) {
	query := `SELECT plugin_item_id FROM tb_bot_plugin WHERE bot_id = ?`

	rows, err := db.QueryContext(ctx, query, botID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// SaveBotPlugins 保存 Bot-插件关联 (删除旧的，添加新的)
func (r *PluginRepository) SaveBotPlugins(ctx context.Context, botID int64, pluginItemIDs []int64) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除旧的关联
	_, err = tx.ExecContext(ctx, `DELETE FROM tb_bot_plugin WHERE bot_id = ?`, botID)
	if err != nil {
		return err
	}

	// 添加新的关联
	if len(pluginItemIDs) > 0 {
		for _, itemID := range pluginItemIDs {
			id, _ := snowflake.GenerateID()
			_, err = tx.ExecContext(ctx,
				`INSERT INTO tb_bot_plugin (id, bot_id, plugin_item_id) VALUES (?, ?, ?)`,
				id, botID, itemID,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

// DeleteBotPlugin 删除单个 Bot-插件关联
func (r *PluginRepository) DeleteBotPlugin(ctx context.Context, botID, pluginItemID int64) error {
	query := `DELETE FROM tb_bot_plugin WHERE bot_id = ? AND plugin_item_id = ?`
	_, err := db.ExecContext(ctx, query, botID, pluginItemID)
	return err
}

// ExistsBotPlugin 检查 Bot-插件关联是否存在
func (r *PluginRepository) ExistsBotPlugin(ctx context.Context, pluginItemID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM tb_bot_plugin WHERE plugin_item_id = ?`
	var count int
	err := db.QueryRowContext(ctx, query, pluginItemID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ========================== PluginCategory ==========================

// ListPluginCategories 获取插件分类列表
func (r *PluginRepository) ListPluginCategories(ctx context.Context) ([]*entity.PluginCategory, error) {
	query := `SELECT id, name, created_at FROM tb_plugin_category ORDER BY id`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.PluginCategory
	for rows.Next() {
		var c entity.PluginCategory
		var createdAt sql.NullTime
		err := rows.Scan(&c.ID, &c.Name, &createdAt)
		if err != nil {
			return nil, err
		}
		if createdAt.Valid {
			c.CreatedAt = createdAt.Time
		}
		categories = append(categories, &c)
	}
	return categories, nil
}

// CreatePluginCategory 创建插件分类
func (r *PluginRepository) CreatePluginCategory(ctx context.Context, c *entity.PluginCategory) error {
	query := `INSERT INTO tb_plugin_category (name, created_at) VALUES (?, ?)`
	result, err := db.ExecContext(ctx, query, c.Name, time.Now())
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	c.ID = id
	return nil
}

// DeletePluginCategory 删除插件分类
func (r *PluginRepository) DeletePluginCategory(ctx context.Context, id int64) error {
	// 先删除关联
	_, _ = db.ExecContext(ctx, `DELETE FROM tb_plugin_category_mapping WHERE category_id = ?`, id)
	// 再删除分类
	query := `DELETE FROM tb_plugin_category WHERE id = ?`
	_, err := db.ExecContext(ctx, query, id)
	return err
}

// ListPluginsWithTools 获取插件列表(包含工具)
func (r *PluginRepository) ListPluginsWithTools(ctx context.Context, tenantID int64) ([]*entity.Plugin, error) {
	plugins, err := r.ListPlugins(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	for _, p := range plugins {
		tools, err := r.ListPluginItemsByPluginID(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		p.Tools = tools
	}
	return plugins, nil
}
