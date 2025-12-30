package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/aiflowy/aiflowy-go/internal/dto"
	"github.com/aiflowy/aiflowy-go/internal/entity"
)

// ModelRepository handles model and model provider database operations
type ModelRepository struct {
	db *sql.DB
}

// NewModelRepository creates a new ModelRepository
func NewModelRepository(db *sql.DB) *ModelRepository {
	return &ModelRepository{db: db}
}

// ========== Model Provider Operations ==========

// GetProviderByID retrieves a model provider by ID
func (r *ModelRepository) GetProviderByID(ctx context.Context, id int64) (*entity.ModelProvider, error) {
	query := `SELECT id, provider_name, COALESCE(provider_type,''), COALESCE(icon,''), COALESCE(api_key,''), COALESCE(endpoint,''),
		COALESCE(chat_path,''), COALESCE(embed_path,''), COALESCE(rerank_path,''), created, created_by, modified, modified_by
		FROM tb_model_provider WHERE id = ?`

	var p entity.ModelProvider
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.ProviderName, &p.ProviderType, &p.Icon, &p.APIKey, &p.Endpoint,
		&p.ChatPath, &p.EmbedPath, &p.RerankPath, &p.Created, &p.CreatedBy, &p.Modified, &p.ModifiedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ListProviders lists all model providers
func (r *ModelRepository) ListProviders(ctx context.Context, req *dto.ModelProviderListRequest) ([]*entity.ModelProvider, error) {
	query := `SELECT id, provider_name, COALESCE(provider_type,''), COALESCE(icon,''), COALESCE(api_key,''), COALESCE(endpoint,''),
		COALESCE(chat_path,''), COALESCE(embed_path,''), COALESCE(rerank_path,''), created, created_by, modified, modified_by
		FROM tb_model_provider WHERE 1=1`
	var args []interface{}

	if req != nil {
		if req.ProviderName != "" {
			query += " AND provider_name LIKE ?"
			args = append(args, "%"+req.ProviderName+"%")
		}
		if req.ProviderType != "" {
			query += " AND provider_type = ?"
			args = append(args, req.ProviderType)
		}
	}

	query += " ORDER BY created DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []*entity.ModelProvider
	for rows.Next() {
		var p entity.ModelProvider
		err := rows.Scan(
			&p.ID, &p.ProviderName, &p.ProviderType, &p.Icon, &p.APIKey, &p.Endpoint,
			&p.ChatPath, &p.EmbedPath, &p.RerankPath, &p.Created, &p.CreatedBy, &p.Modified, &p.ModifiedBy,
		)
		if err != nil {
			return nil, err
		}
		providers = append(providers, &p)
	}
	return providers, nil
}

// PageProviders returns paginated providers
func (r *ModelRepository) PageProviders(ctx context.Context, req *dto.PageRequest, filter *dto.ModelProviderListRequest) ([]*entity.ModelProvider, int64, error) {
	countQuery := "SELECT COUNT(*) FROM tb_model_provider WHERE 1=1"
	query := `SELECT id, provider_name, COALESCE(provider_type,''), COALESCE(icon,''), COALESCE(api_key,''), COALESCE(endpoint,''),
		COALESCE(chat_path,''), COALESCE(embed_path,''), COALESCE(rerank_path,''), created, created_by, modified, modified_by
		FROM tb_model_provider WHERE 1=1`
	var args []interface{}

	if filter != nil {
		if filter.ProviderName != "" {
			countQuery += " AND provider_name LIKE ?"
			query += " AND provider_name LIKE ?"
			args = append(args, "%"+filter.ProviderName+"%")
		}
		if filter.ProviderType != "" {
			countQuery += " AND provider_type = ?"
			query += " AND provider_type = ?"
			args = append(args, filter.ProviderType)
		}
	}

	var total int64
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query += " ORDER BY created DESC LIMIT ? OFFSET ?"
	args = append(args, req.GetPageSize(), req.GetOffset())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var providers []*entity.ModelProvider
	for rows.Next() {
		var p entity.ModelProvider
		err := rows.Scan(
			&p.ID, &p.ProviderName, &p.ProviderType, &p.Icon, &p.APIKey, &p.Endpoint,
			&p.ChatPath, &p.EmbedPath, &p.RerankPath, &p.Created, &p.CreatedBy, &p.Modified, &p.ModifiedBy,
		)
		if err != nil {
			return nil, 0, err
		}
		providers = append(providers, &p)
	}
	return providers, total, nil
}

// CreateProvider creates a new model provider
func (r *ModelRepository) CreateProvider(ctx context.Context, p *entity.ModelProvider) error {
	query := `INSERT INTO tb_model_provider
		(id, provider_name, provider_type, icon, api_key, endpoint, chat_path, embed_path, rerank_path, created, created_by, modified, modified_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.ProviderName, p.ProviderType, p.Icon, p.APIKey, p.Endpoint,
		p.ChatPath, p.EmbedPath, p.RerankPath, p.Created, p.CreatedBy, p.Modified, p.ModifiedBy,
	)
	return err
}

// UpdateProvider updates an existing model provider
func (r *ModelRepository) UpdateProvider(ctx context.Context, p *entity.ModelProvider) error {
	query := `UPDATE tb_model_provider SET
		provider_name = ?, provider_type = ?, icon = ?, api_key = ?, endpoint = ?,
		chat_path = ?, embed_path = ?, rerank_path = ?, modified = ?, modified_by = ?
		WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		p.ProviderName, p.ProviderType, p.Icon, p.APIKey, p.Endpoint,
		p.ChatPath, p.EmbedPath, p.RerankPath, p.Modified, p.ModifiedBy, p.ID,
	)
	return err
}

// DeleteProvider deletes a model provider
func (r *ModelRepository) DeleteProvider(ctx context.Context, id int64) error {
	query := "DELETE FROM tb_model_provider WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// ========== Model Operations ==========

// GetModelByID retrieves a model by ID
func (r *ModelRepository) GetModelByID(ctx context.Context, id int64) (*entity.Model, error) {
	query := `SELECT id, dept_id, tenant_id, provider_id, COALESCE(title,''), COALESCE(icon,''), COALESCE(description,''), COALESCE(endpoint,''),
		COALESCE(request_path,''), COALESCE(model_name,''), COALESCE(api_key,''), COALESCE(extra_config,''), COALESCE(options,''), COALESCE(group_name,''), COALESCE(model_type,''),
		COALESCE(with_used,false), COALESCE(support_thinking,false), COALESCE(support_tool,false), COALESCE(support_image,false), COALESCE(support_image_b64_only,false),
		COALESCE(support_video,false), COALESCE(support_audio,false), COALESCE(support_free,false)
		FROM tb_model WHERE id = ?`

	var m entity.Model
	var providerID sql.NullInt64
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.DeptID, &m.TenantID, &providerID, &m.Title, &m.Icon, &m.Description, &m.Endpoint,
		&m.RequestPath, &m.ModelName, &m.APIKey, &m.ExtraConfig, &m.Options, &m.GroupName, &m.ModelType,
		&m.WithUsed, &m.SupportThinking, &m.SupportTool, &m.SupportImage, &m.SupportImageB64Only,
		&m.SupportVideo, &m.SupportAudio, &m.SupportFree,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if providerID.Valid {
		m.ProviderID = providerID.Int64
	}
	return &m, nil
}

// GetModelWithProvider retrieves a model with its provider info
func (r *ModelRepository) GetModelWithProvider(ctx context.Context, id int64) (*entity.ModelWithProvider, error) {
	query := `SELECT m.id, m.dept_id, m.tenant_id, m.provider_id, COALESCE(m.title,''), COALESCE(m.icon,''), COALESCE(m.description,''), COALESCE(m.endpoint,''),
		COALESCE(m.request_path,''), COALESCE(m.model_name,''), COALESCE(m.api_key,''), COALESCE(m.extra_config,''), COALESCE(m.options,''), COALESCE(m.group_name,''), COALESCE(m.model_type,''),
		COALESCE(m.with_used,false), COALESCE(m.support_thinking,false), COALESCE(m.support_tool,false), COALESCE(m.support_image,false), COALESCE(m.support_image_b64_only,false),
		COALESCE(m.support_video,false), COALESCE(m.support_audio,false), COALESCE(m.support_free,false),
		COALESCE(p.provider_name, ''), COALESCE(p.provider_type, '')
		FROM tb_model m
		LEFT JOIN tb_model_provider p ON m.provider_id = p.id
		WHERE m.id = ?`

	var m entity.ModelWithProvider
	var providerID sql.NullInt64
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.DeptID, &m.TenantID, &providerID, &m.Title, &m.Icon, &m.Description, &m.Endpoint,
		&m.RequestPath, &m.ModelName, &m.APIKey, &m.ExtraConfig, &m.Options, &m.GroupName, &m.ModelType,
		&m.WithUsed, &m.SupportThinking, &m.SupportTool, &m.SupportImage, &m.SupportImageB64Only,
		&m.SupportVideo, &m.SupportAudio, &m.SupportFree, &m.ProviderName, &m.ProviderType,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if providerID.Valid {
		m.ProviderID = providerID.Int64
	}
	return &m, nil
}

// ListModels lists models with optional filters
func (r *ModelRepository) ListModels(ctx context.Context, req *dto.ModelListRequest) ([]*entity.Model, error) {
	query := `SELECT id, dept_id, tenant_id, provider_id, COALESCE(title,''), COALESCE(icon,''), COALESCE(description,''), COALESCE(endpoint,''),
		COALESCE(request_path,''), COALESCE(model_name,''), COALESCE(api_key,''), COALESCE(extra_config,''), COALESCE(options,''), COALESCE(group_name,''), COALESCE(model_type,''),
		COALESCE(with_used,false), COALESCE(support_thinking,false), COALESCE(support_tool,false), COALESCE(support_image,false), COALESCE(support_image_b64_only,false),
		COALESCE(support_video,false), COALESCE(support_audio,false), COALESCE(support_free,false)
		FROM tb_model WHERE 1=1`
	var args []interface{}

	if req != nil {
		if req.ProviderID > 0 {
			query += " AND provider_id = ?"
			args = append(args, req.ProviderID)
		}
		if req.ModelType != "" {
			query += " AND model_type = ?"
			args = append(args, req.ModelType)
		}
		if req.WithUsed != nil {
			query += " AND with_used = ?"
			args = append(args, *req.WithUsed)
		}
		if req.GroupName != "" {
			query += " AND group_name = ?"
			args = append(args, req.GroupName)
		}
		if req.SelectText != "" {
			query += " AND (title LIKE ? OR model_name LIKE ?)"
			args = append(args, "%"+req.SelectText+"%", "%"+req.SelectText+"%")
		}
	}

	query += " ORDER BY group_name, title"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []*entity.Model
	for rows.Next() {
		var m entity.Model
		var providerID sql.NullInt64
		err := rows.Scan(
			&m.ID, &m.DeptID, &m.TenantID, &providerID, &m.Title, &m.Icon, &m.Description, &m.Endpoint,
			&m.RequestPath, &m.ModelName, &m.APIKey, &m.ExtraConfig, &m.Options, &m.GroupName, &m.ModelType,
			&m.WithUsed, &m.SupportThinking, &m.SupportTool, &m.SupportImage, &m.SupportImageB64Only,
			&m.SupportVideo, &m.SupportAudio, &m.SupportFree,
		)
		if err != nil {
			return nil, err
		}
		if providerID.Valid {
			m.ProviderID = providerID.Int64
		}
		models = append(models, &m)
	}
	return models, nil
}

// PageModels returns paginated models
func (r *ModelRepository) PageModels(ctx context.Context, req *dto.ModelListRequest) ([]*entity.Model, int64, error) {
	countQuery := "SELECT COUNT(*) FROM tb_model WHERE 1=1"
	query := `SELECT id, dept_id, tenant_id, provider_id, COALESCE(title,''), COALESCE(icon,''), COALESCE(description,''), COALESCE(endpoint,''),
		COALESCE(request_path,''), COALESCE(model_name,''), COALESCE(api_key,''), COALESCE(extra_config,''), COALESCE(options,''), COALESCE(group_name,''), COALESCE(model_type,''),
		COALESCE(with_used,false), COALESCE(support_thinking,false), COALESCE(support_tool,false), COALESCE(support_image,false), COALESCE(support_image_b64_only,false),
		COALESCE(support_video,false), COALESCE(support_audio,false), COALESCE(support_free,false)
		FROM tb_model WHERE 1=1`
	var args []interface{}

	if req != nil {
		if req.ProviderID > 0 {
			countQuery += " AND provider_id = ?"
			query += " AND provider_id = ?"
			args = append(args, req.ProviderID)
		}
		if req.ModelType != "" {
			countQuery += " AND model_type = ?"
			query += " AND model_type = ?"
			args = append(args, req.ModelType)
		}
		if req.WithUsed != nil {
			countQuery += " AND with_used = ?"
			query += " AND with_used = ?"
			args = append(args, *req.WithUsed)
		}
		if req.GroupName != "" {
			countQuery += " AND group_name = ?"
			query += " AND group_name = ?"
			args = append(args, req.GroupName)
		}
		if req.SelectText != "" {
			countQuery += " AND (title LIKE ? OR model_name LIKE ?)"
			query += " AND (title LIKE ? OR model_name LIKE ?)"
			args = append(args, "%"+req.SelectText+"%", "%"+req.SelectText+"%")
		}
	}

	var total int64
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query += " ORDER BY group_name, title LIMIT ? OFFSET ?"
	args = append(args, req.GetPageSize(), req.GetOffset())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var models []*entity.Model
	for rows.Next() {
		var m entity.Model
		var providerID sql.NullInt64
		err := rows.Scan(
			&m.ID, &m.DeptID, &m.TenantID, &providerID, &m.Title, &m.Icon, &m.Description, &m.Endpoint,
			&m.RequestPath, &m.ModelName, &m.APIKey, &m.ExtraConfig, &m.Options, &m.GroupName, &m.ModelType,
			&m.WithUsed, &m.SupportThinking, &m.SupportTool, &m.SupportImage, &m.SupportImageB64Only,
			&m.SupportVideo, &m.SupportAudio, &m.SupportFree,
		)
		if err != nil {
			return nil, 0, err
		}
		if providerID.Valid {
			m.ProviderID = providerID.Int64
		}
		models = append(models, &m)
	}
	return models, total, nil
}

// GetModelsGroupedByType returns models grouped by model type and group name
func (r *ModelRepository) GetModelsGroupedByType(ctx context.Context, req *dto.ModelByProviderRequest) (map[string]map[string][]*entity.Model, error) {
	query := `SELECT id, dept_id, tenant_id, provider_id, COALESCE(title,''), COALESCE(icon,''), COALESCE(description,''), COALESCE(endpoint,''),
		COALESCE(request_path,''), COALESCE(model_name,''), COALESCE(api_key,''), COALESCE(extra_config,''), COALESCE(options,''), COALESCE(group_name,''), COALESCE(model_type,''),
		COALESCE(with_used,false), COALESCE(support_thinking,false), COALESCE(support_tool,false), COALESCE(support_image,false), COALESCE(support_image_b64_only,false),
		COALESCE(support_video,false), COALESCE(support_audio,false), COALESCE(support_free,false)
		FROM tb_model WHERE 1=1`
	var args []interface{}

	if req != nil {
		if req.ProviderID > 0 {
			query += " AND provider_id = ?"
			args = append(args, req.ProviderID)
		}
		if req.ModelType != "" {
			query += " AND model_type = ?"
			args = append(args, req.ModelType)
		}
		if req.WithUsed != nil {
			query += " AND with_used = ?"
			args = append(args, *req.WithUsed)
		}
	}

	query += " ORDER BY model_type, group_name, title"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Result: map[modelType]map[groupName][]Model
	result := make(map[string]map[string][]*entity.Model)

	for rows.Next() {
		var m entity.Model
		var providerID sql.NullInt64
		err := rows.Scan(
			&m.ID, &m.DeptID, &m.TenantID, &providerID, &m.Title, &m.Icon, &m.Description, &m.Endpoint,
			&m.RequestPath, &m.ModelName, &m.APIKey, &m.ExtraConfig, &m.Options, &m.GroupName, &m.ModelType,
			&m.WithUsed, &m.SupportThinking, &m.SupportTool, &m.SupportImage, &m.SupportImageB64Only,
			&m.SupportVideo, &m.SupportAudio, &m.SupportFree,
		)
		if err != nil {
			return nil, err
		}
		if providerID.Valid {
			m.ProviderID = providerID.Int64
		}

		modelType := m.ModelType
		groupName := m.GroupName
		if groupName == "" {
			groupName = "Default"
		}

		if result[modelType] == nil {
			result[modelType] = make(map[string][]*entity.Model)
		}
		model := m
		result[modelType][groupName] = append(result[modelType][groupName], &model)
	}

	return result, nil
}

// CreateModel creates a new model
func (r *ModelRepository) CreateModel(ctx context.Context, m *entity.Model) error {
	query := `INSERT INTO tb_model
		(id, dept_id, tenant_id, provider_id, title, icon, description, endpoint, request_path,
		model_name, api_key, extra_config, options, group_name, model_type, with_used,
		support_thinking, support_tool, support_image, support_image_b64_only,
		support_video, support_audio, support_free)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var providerID interface{} = nil
	if m.ProviderID > 0 {
		providerID = m.ProviderID
	}

	_, err := r.db.ExecContext(ctx, query,
		m.ID, m.DeptID, m.TenantID, providerID, m.Title, m.Icon, m.Description, m.Endpoint, m.RequestPath,
		m.ModelName, m.APIKey, m.ExtraConfig, m.Options, m.GroupName, m.ModelType, m.WithUsed,
		m.SupportThinking, m.SupportTool, m.SupportImage, m.SupportImageB64Only,
		m.SupportVideo, m.SupportAudio, m.SupportFree,
	)
	return err
}

// UpdateModel updates an existing model
func (r *ModelRepository) UpdateModel(ctx context.Context, m *entity.Model) error {
	query := `UPDATE tb_model SET
		dept_id = ?, tenant_id = ?, provider_id = ?, title = ?, icon = ?, description = ?,
		endpoint = ?, request_path = ?, model_name = ?, api_key = ?, extra_config = ?,
		options = ?, group_name = ?, model_type = ?, with_used = ?, support_thinking = ?,
		support_tool = ?, support_image = ?, support_image_b64_only = ?,
		support_video = ?, support_audio = ?, support_free = ?
		WHERE id = ?`

	var providerID interface{} = nil
	if m.ProviderID > 0 {
		providerID = m.ProviderID
	}

	_, err := r.db.ExecContext(ctx, query,
		m.DeptID, m.TenantID, providerID, m.Title, m.Icon, m.Description,
		m.Endpoint, m.RequestPath, m.ModelName, m.APIKey, m.ExtraConfig,
		m.Options, m.GroupName, m.ModelType, m.WithUsed, m.SupportThinking,
		m.SupportTool, m.SupportImage, m.SupportImageB64Only,
		m.SupportVideo, m.SupportAudio, m.SupportFree, m.ID,
	)
	return err
}

// UpdateModelsByCondition updates models matching conditions
func (r *ModelRepository) UpdateModelsByCondition(ctx context.Context, req *dto.UpdateByEntityRequest) error {
	var setClauses []string
	var args []interface{}

	if req.WithUsed != nil {
		setClauses = append(setClauses, "with_used = ?")
		args = append(args, *req.WithUsed)
	}
	if req.SupportThinking != nil {
		setClauses = append(setClauses, "support_thinking = ?")
		args = append(args, *req.SupportThinking)
	}
	if req.SupportTool != nil {
		setClauses = append(setClauses, "support_tool = ?")
		args = append(args, *req.SupportTool)
	}
	if req.SupportImage != nil {
		setClauses = append(setClauses, "support_image = ?")
		args = append(args, *req.SupportImage)
	}

	if len(setClauses) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE tb_model SET %s WHERE 1=1", strings.Join(setClauses, ", "))

	if req.ProviderID > 0 {
		query += " AND provider_id = ?"
		args = append(args, req.ProviderID)
	}
	if req.GroupName != "" {
		query += " AND group_name = ?"
		args = append(args, req.GroupName)
	}
	if req.ModelType != "" {
		query += " AND model_type = ?"
		args = append(args, req.ModelType)
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// DeleteModel deletes a model
func (r *ModelRepository) DeleteModel(ctx context.Context, id int64) error {
	query := "DELETE FROM tb_model WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// DeleteModelsByCondition deletes models matching conditions
func (r *ModelRepository) DeleteModelsByCondition(ctx context.Context, req *dto.RemoveByEntityRequest) error {
	query := "DELETE FROM tb_model WHERE 1=1"
	var args []interface{}

	if req.ProviderID > 0 {
		query += " AND provider_id = ?"
		args = append(args, req.ProviderID)
	}
	if req.GroupName != "" {
		query += " AND group_name = ?"
		args = append(args, req.GroupName)
	}
	if req.ModelType != "" {
		query += " AND model_type = ?"
		args = append(args, req.ModelType)
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// DeleteModelsByIDs deletes models by IDs
func (r *ModelRepository) DeleteModelsByIDs(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf("DELETE FROM tb_model WHERE id IN (%s)", strings.Join(placeholders, ","))
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// CountModelsByProvider counts models for a provider
func (r *ModelRepository) CountModelsByProvider(ctx context.Context, providerID int64) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM tb_model WHERE provider_id = ?", providerID).Scan(&count)
	return count, err
}

// GetDefaultTenantAndDept gets default tenant and dept IDs
func (r *ModelRepository) GetDefaultTenantAndDept(ctx context.Context) (tenantID, deptID int64, err error) {
	// Default values from the admin user
	return 1000000, 1, nil
}

// GetModelInstance gets a model with inherited provider config
func (r *ModelRepository) GetModelInstance(ctx context.Context, id int64) (*entity.Model, error) {
	model, err := r.GetModelByID(ctx, id)
	if err != nil || model == nil {
		return model, err
	}

	// Load provider if exists
	if model.ProviderID > 0 {
		provider, err := r.GetProviderByID(ctx, model.ProviderID)
		if err != nil {
			return nil, err
		}
		if provider != nil {
			model.ModelProvider = provider
			// Inherit missing fields from provider
			if model.Endpoint == "" {
				model.Endpoint = provider.Endpoint
			}
			if model.APIKey == "" {
				model.APIKey = provider.APIKey
			}
			if model.RequestPath == "" {
				switch model.ModelType {
				case entity.ModelTypeChatModel:
					model.RequestPath = provider.ChatPath
				case entity.ModelTypeEmbeddingModel:
					model.RequestPath = provider.EmbedPath
				case entity.ModelTypeRerankModel:
					model.RequestPath = provider.RerankPath
				}
			}
		}
	}

	return model, nil
}

var modelRepo *ModelRepository
var modelRepoInit = false

// GetModelRepository returns the singleton ModelRepository
func GetModelRepository() *ModelRepository {
	if !modelRepoInit {
		modelRepo = NewModelRepository(GetDB())
		modelRepoInit = true
	}
	return modelRepo
}

// Helper function
func timeNow() time.Time {
	return time.Now()
}
