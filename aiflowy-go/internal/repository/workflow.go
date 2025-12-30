package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/pkg/snowflake"
)

// WorkflowRepository 工作流数据访问层
type WorkflowRepository struct {
	db *sql.DB
}

// NewWorkflowRepository 创建 WorkflowRepository
func NewWorkflowRepository() *WorkflowRepository {
	return &WorkflowRepository{
		db: GetDB(),
	}
}

// ========================== Workflow ==========================

// GetWorkflowByID 根据 ID 获取工作流
func (r *WorkflowRepository) GetWorkflowByID(ctx context.Context, id int64) (*entity.Workflow, error) {
	query := `SELECT id, alias, dept_id, tenant_id, title, description, icon, content,
		created, created_by, modified, modified_by, english_name, status, category_id
		FROM tb_workflow WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)
	workflow := &entity.Workflow{}

	var alias, description, icon, content, englishName sql.NullString
	var categoryID sql.NullInt64
	var modified sql.NullTime
	var modifiedBy sql.NullInt64

	err := row.Scan(
		&workflow.ID, &alias, &workflow.DeptID, &workflow.TenantID,
		&workflow.Title, &description, &icon, &content,
		&workflow.Created, &workflow.CreatedBy, &modified, &modifiedBy,
		&englishName, &workflow.Status, &categoryID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	workflow.Alias = alias.String
	workflow.Description = description.String
	workflow.Icon = icon.String
	workflow.Content = content.String
	workflow.EnglishName = englishName.String
	if categoryID.Valid {
		workflow.CategoryID = categoryID.Int64
	}
	if modified.Valid {
		workflow.Modified = modified.Time
	}
	if modifiedBy.Valid {
		workflow.ModifiedBy = modifiedBy.Int64
	}

	return workflow, nil
}

// GetWorkflowByAlias 根据别名获取工作流
func (r *WorkflowRepository) GetWorkflowByAlias(ctx context.Context, alias string) (*entity.Workflow, error) {
	query := `SELECT id, alias, dept_id, tenant_id, title, description, icon, content,
		created, created_by, modified, modified_by, english_name, status, category_id
		FROM tb_workflow WHERE alias = ?`

	row := r.db.QueryRowContext(ctx, query, alias)
	workflow := &entity.Workflow{}

	var aliasVal, description, icon, content, englishName sql.NullString
	var categoryID sql.NullInt64
	var modified sql.NullTime
	var modifiedBy sql.NullInt64

	err := row.Scan(
		&workflow.ID, &aliasVal, &workflow.DeptID, &workflow.TenantID,
		&workflow.Title, &description, &icon, &content,
		&workflow.Created, &workflow.CreatedBy, &modified, &modifiedBy,
		&englishName, &workflow.Status, &categoryID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	workflow.Alias = aliasVal.String
	workflow.Description = description.String
	workflow.Icon = icon.String
	workflow.Content = content.String
	workflow.EnglishName = englishName.String
	if categoryID.Valid {
		workflow.CategoryID = categoryID.Int64
	}
	if modified.Valid {
		workflow.Modified = modified.Time
	}
	if modifiedBy.Valid {
		workflow.ModifiedBy = modifiedBy.Int64
	}

	return workflow, nil
}

// ListWorkflows 获取工作流列表
func (r *WorkflowRepository) ListWorkflows(ctx context.Context, tenantID int64) ([]*entity.Workflow, error) {
	query := `SELECT w.id, w.alias, w.dept_id, w.tenant_id, w.title, w.description, w.icon,
		w.created, w.created_by, w.modified, w.modified_by, w.english_name, w.status, w.category_id,
		c.category_name
		FROM tb_workflow w
		LEFT JOIN tb_workflow_category c ON w.category_id = c.id
		WHERE w.tenant_id = ? AND w.status >= 0
		ORDER BY w.created DESC`

	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []*entity.Workflow
	for rows.Next() {
		workflow := &entity.Workflow{}
		var alias, description, icon, englishName, categoryName sql.NullString
		var categoryID sql.NullInt64
		var modified sql.NullTime
		var modifiedBy sql.NullInt64

		err := rows.Scan(
			&workflow.ID, &alias, &workflow.DeptID, &workflow.TenantID,
			&workflow.Title, &description, &icon,
			&workflow.Created, &workflow.CreatedBy, &modified, &modifiedBy,
			&englishName, &workflow.Status, &categoryID, &categoryName,
		)
		if err != nil {
			return nil, err
		}

		workflow.Alias = alias.String
		workflow.Description = description.String
		workflow.Icon = icon.String
		workflow.EnglishName = englishName.String
		workflow.CategoryName = categoryName.String
		if categoryID.Valid {
			workflow.CategoryID = categoryID.Int64
		}
		if modified.Valid {
			workflow.Modified = modified.Time
		}
		if modifiedBy.Valid {
			workflow.ModifiedBy = modifiedBy.Int64
		}

		workflows = append(workflows, workflow)
	}

	return workflows, nil
}

// CreateWorkflow 创建工作流
func (r *WorkflowRepository) CreateWorkflow(ctx context.Context, workflow *entity.Workflow) error {
	query := `INSERT INTO tb_workflow (id, alias, dept_id, tenant_id, title, description, icon,
		content, created, created_by, english_name, status, category_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		workflow.ID, nullString(workflow.Alias), workflow.DeptID, workflow.TenantID,
		workflow.Title, nullString(workflow.Description), nullString(workflow.Icon),
		nullString(workflow.Content), workflow.Created, workflow.CreatedBy,
		nullString(workflow.EnglishName), workflow.Status, nullInt64(workflow.CategoryID),
	)
	return err
}

// UpdateWorkflow 更新工作流
func (r *WorkflowRepository) UpdateWorkflow(ctx context.Context, workflow *entity.Workflow) error {
	query := `UPDATE tb_workflow SET alias = ?, title = ?, description = ?, icon = ?,
		content = ?, modified = ?, modified_by = ?, english_name = ?, status = ?, category_id = ?
		WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		nullString(workflow.Alias), workflow.Title, nullString(workflow.Description),
		nullString(workflow.Icon), nullString(workflow.Content),
		workflow.Modified, workflow.ModifiedBy, nullString(workflow.EnglishName),
		workflow.Status, nullInt64(workflow.CategoryID), workflow.ID,
	)
	return err
}

// DeleteWorkflow 删除工作流 (软删除)
func (r *WorkflowRepository) DeleteWorkflow(ctx context.Context, id int64) error {
	query := `UPDATE tb_workflow SET status = -1 WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// ========================== WorkflowCategory ==========================

// GetWorkflowCategoryByID 根据 ID 获取工作流分类
func (r *WorkflowRepository) GetWorkflowCategoryByID(ctx context.Context, id int64) (*entity.WorkflowCategory, error) {
	query := `SELECT id, category_name, sort_no, created, created_by, modified, modified_by, status
		FROM tb_workflow_category WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)
	category := &entity.WorkflowCategory{}

	var modified sql.NullTime
	var modifiedBy sql.NullInt64

	err := row.Scan(
		&category.ID, &category.CategoryName, &category.SortNo,
		&category.Created, &category.CreatedBy, &modified, &modifiedBy, &category.Status,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if modified.Valid {
		category.Modified = modified.Time
	}
	if modifiedBy.Valid {
		category.ModifiedBy = modifiedBy.Int64
	}

	return category, nil
}

// ListWorkflowCategories 获取工作流分类列表
func (r *WorkflowRepository) ListWorkflowCategories(ctx context.Context) ([]*entity.WorkflowCategory, error) {
	query := `SELECT id, category_name, sort_no, created, created_by, modified, modified_by, status
		FROM tb_workflow_category WHERE status >= 0 ORDER BY sort_no, created`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.WorkflowCategory
	for rows.Next() {
		category := &entity.WorkflowCategory{}
		var modified sql.NullTime
		var modifiedBy sql.NullInt64

		err := rows.Scan(
			&category.ID, &category.CategoryName, &category.SortNo,
			&category.Created, &category.CreatedBy, &modified, &modifiedBy, &category.Status,
		)
		if err != nil {
			return nil, err
		}

		if modified.Valid {
			category.Modified = modified.Time
		}
		if modifiedBy.Valid {
			category.ModifiedBy = modifiedBy.Int64
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// CreateWorkflowCategory 创建工作流分类
func (r *WorkflowRepository) CreateWorkflowCategory(ctx context.Context, category *entity.WorkflowCategory) error {
	query := `INSERT INTO tb_workflow_category (id, category_name, sort_no, created, created_by, modified, modified_by, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		category.ID, category.CategoryName, category.SortNo,
		category.Created, category.CreatedBy, category.Modified, category.ModifiedBy, category.Status,
	)
	return err
}

// UpdateWorkflowCategory 更新工作流分类
func (r *WorkflowRepository) UpdateWorkflowCategory(ctx context.Context, category *entity.WorkflowCategory) error {
	query := `UPDATE tb_workflow_category SET category_name = ?, sort_no = ?, modified = ?, modified_by = ?, status = ?
		WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		category.CategoryName, category.SortNo, category.Modified, category.ModifiedBy, category.Status, category.ID,
	)
	return err
}

// DeleteWorkflowCategory 删除工作流分类 (软删除)
func (r *WorkflowRepository) DeleteWorkflowCategory(ctx context.Context, id int64) error {
	query := `UPDATE tb_workflow_category SET status = -1 WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// ========================== BotWorkflow ==========================

// GetBotWorkflowIDs 获取 Bot 关联的工作流 ID 列表
func (r *WorkflowRepository) GetBotWorkflowIDs(ctx context.Context, botID int64) ([]int64, error) {
	query := `SELECT workflow_id FROM tb_bot_workflow WHERE bot_id = ?`

	rows, err := r.db.QueryContext(ctx, query, botID)
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

// ListBotWorkflows 获取 Bot 关联的工作流列表 (包含工作流详情)
func (r *WorkflowRepository) ListBotWorkflows(ctx context.Context, botID int64) ([]*entity.BotWorkflow, error) {
	query := `SELECT bw.id, bw.bot_id, bw.workflow_id, bw.options,
		w.id, w.alias, w.title, w.description, w.icon, w.english_name, w.status
		FROM tb_bot_workflow bw
		LEFT JOIN tb_workflow w ON bw.workflow_id = w.id
		WHERE bw.bot_id = ? AND w.status >= 0`

	rows, err := r.db.QueryContext(ctx, query, botID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var botWorkflows []*entity.BotWorkflow
	for rows.Next() {
		bw := &entity.BotWorkflow{}
		w := &entity.Workflow{}
		var options, alias, description, icon, englishName sql.NullString

		err := rows.Scan(
			&bw.ID, &bw.BotID, &bw.WorkflowID, &options,
			&w.ID, &alias, &w.Title, &description, &icon, &englishName, &w.Status,
		)
		if err != nil {
			return nil, err
		}

		bw.Options = options.String
		w.Alias = alias.String
		w.Description = description.String
		w.Icon = icon.String
		w.EnglishName = englishName.String
		bw.Workflow = w

		botWorkflows = append(botWorkflows, bw)
	}

	return botWorkflows, nil
}

// ListWorkflowsByBotID 获取 Bot 关联的工作流
func (r *WorkflowRepository) ListWorkflowsByBotID(ctx context.Context, botID int64) ([]*entity.Workflow, error) {
	query := `SELECT w.id, w.alias, w.dept_id, w.tenant_id, w.title, w.description, w.icon, w.content,
		w.created, w.created_by, w.modified, w.modified_by, w.english_name, w.status, w.category_id
		FROM tb_workflow w
		INNER JOIN tb_bot_workflow bw ON w.id = bw.workflow_id
		WHERE bw.bot_id = ? AND w.status >= 0`

	rows, err := r.db.QueryContext(ctx, query, botID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []*entity.Workflow
	for rows.Next() {
		workflow := &entity.Workflow{}
		var alias, description, icon, content, englishName sql.NullString
		var categoryID sql.NullInt64
		var modified sql.NullTime
		var modifiedBy sql.NullInt64

		err := rows.Scan(
			&workflow.ID, &alias, &workflow.DeptID, &workflow.TenantID,
			&workflow.Title, &description, &icon, &content,
			&workflow.Created, &workflow.CreatedBy, &modified, &modifiedBy,
			&englishName, &workflow.Status, &categoryID,
		)
		if err != nil {
			return nil, err
		}

		workflow.Alias = alias.String
		workflow.Description = description.String
		workflow.Icon = icon.String
		workflow.Content = content.String
		workflow.EnglishName = englishName.String
		if categoryID.Valid {
			workflow.CategoryID = categoryID.Int64
		}
		if modified.Valid {
			workflow.Modified = modified.Time
		}
		if modifiedBy.Valid {
			workflow.ModifiedBy = modifiedBy.Int64
		}

		workflows = append(workflows, workflow)
	}

	return workflows, nil
}

// SaveBotWorkflows 保存 Bot-工作流关联 (全量替换)
func (r *WorkflowRepository) SaveBotWorkflows(ctx context.Context, botID int64, workflowIDs []int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除旧的关联
	_, err = tx.ExecContext(ctx, `DELETE FROM tb_bot_workflow WHERE bot_id = ?`, botID)
	if err != nil {
		return err
	}

	// 添加新的关联
	if len(workflowIDs) > 0 {
		for _, workflowID := range workflowIDs {
			id := snowflake.MustGenerateID()
			_, err = tx.ExecContext(ctx,
				`INSERT INTO tb_bot_workflow (id, bot_id, workflow_id) VALUES (?, ?, ?)`,
				id, botID, workflowID,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

// ExistsBotWorkflow 检查工作流是否有 Bot 关联
func (r *WorkflowRepository) ExistsBotWorkflow(ctx context.Context, workflowID int64) (bool, error) {
	query := `SELECT 1 FROM tb_bot_workflow WHERE workflow_id = ? LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, workflowID)
	var exists int
	err := row.Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteBotWorkflow 删除单个 Bot-工作流关联
func (r *WorkflowRepository) DeleteBotWorkflow(ctx context.Context, botID, workflowID int64) error {
	query := `DELETE FROM tb_bot_workflow WHERE bot_id = ? AND workflow_id = ?`
	_, err := r.db.ExecContext(ctx, query, botID, workflowID)
	return err
}

// helper functions
func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullInt64(i int64) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: i, Valid: true}
}

func nullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: t, Valid: true}
}
