package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/pkg/snowflake"
)

// DocumentCollectionRepository 知识库数据访问层
type DocumentCollectionRepository struct {
	db *sql.DB
}

// NewDocumentCollectionRepository 创建 DocumentCollectionRepository
func NewDocumentCollectionRepository() *DocumentCollectionRepository {
	return &DocumentCollectionRepository{
		db: GetDB(),
	}
}

// Create 创建知识库
func (r *DocumentCollectionRepository) Create(ctx context.Context, dc *entity.DocumentCollection) error {
	if dc.ID == 0 {
		dc.ID, _ = snowflake.GenerateID()
	}
	now := time.Now()
	dc.Created = &now
	dc.Modified = &now

	query := `
		INSERT INTO tb_document_collection
		(id, alias, dept_id, tenant_id, icon, title, description, slug,
		 vector_store_enable, vector_store_type, vector_store_collection, vector_store_config,
		 vector_embed_model_id, rerank_model_id, search_engine_enable, english_name, options,
		 created, created_by, modified, modified_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		dc.ID, dc.Alias, dc.DeptID, dc.TenantID, dc.Icon, dc.Title, dc.Description, dc.Slug,
		dc.VectorStoreEnable, dc.VectorStoreType, dc.VectorStoreCollection, dc.VectorStoreConfig,
		dc.VectorEmbedModelID, dc.RerankModelID, dc.SearchEngineEnable, dc.EnglishName, dc.Options,
		dc.Created, dc.CreatedBy, dc.Modified, dc.ModifiedBy,
	)
	return err
}

// Update 更新知识库
func (r *DocumentCollectionRepository) Update(ctx context.Context, dc *entity.DocumentCollection) error {
	now := time.Now()
	dc.Modified = &now

	query := `
		UPDATE tb_document_collection SET
			alias = ?, icon = ?, title = ?, description = ?, slug = ?,
			vector_store_enable = ?, vector_store_type = ?, vector_store_collection = ?, vector_store_config = ?,
			vector_embed_model_id = ?, rerank_model_id = ?, search_engine_enable = ?, english_name = ?, options = ?,
			modified = ?, modified_by = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		dc.Alias, dc.Icon, dc.Title, dc.Description, dc.Slug,
		dc.VectorStoreEnable, dc.VectorStoreType, dc.VectorStoreCollection, dc.VectorStoreConfig,
		dc.VectorEmbedModelID, dc.RerankModelID, dc.SearchEngineEnable, dc.EnglishName, dc.Options,
		dc.Modified, dc.ModifiedBy, dc.ID,
	)
	return err
}

// Delete 删除知识库
func (r *DocumentCollectionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tb_document_collection WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetByID 根据 ID 获取知识库
func (r *DocumentCollectionRepository) GetByID(ctx context.Context, id int64) (*entity.DocumentCollection, error) {
	query := `
		SELECT id, alias, dept_id, tenant_id, icon, title, description, slug,
		       vector_store_enable, vector_store_type, vector_store_collection, vector_store_config,
		       vector_embed_model_id, rerank_model_id, search_engine_enable, english_name, options,
		       created, created_by, modified, modified_by
		FROM tb_document_collection
		WHERE id = ?
	`

	return r.scanOne(ctx, query, id)
}

// GetByAlias 根据别名获取知识库
func (r *DocumentCollectionRepository) GetByAlias(ctx context.Context, alias string) (*entity.DocumentCollection, error) {
	query := `
		SELECT id, alias, dept_id, tenant_id, icon, title, description, slug,
		       vector_store_enable, vector_store_type, vector_store_collection, vector_store_config,
		       vector_embed_model_id, rerank_model_id, search_engine_enable, english_name, options,
		       created, created_by, modified, modified_by
		FROM tb_document_collection
		WHERE alias = ?
	`

	return r.scanOne(ctx, query, alias)
}

// GetBySlug 根据 slug 获取知识库
func (r *DocumentCollectionRepository) GetBySlug(ctx context.Context, slug string) (*entity.DocumentCollection, error) {
	query := `
		SELECT id, alias, dept_id, tenant_id, icon, title, description, slug,
		       vector_store_enable, vector_store_type, vector_store_collection, vector_store_config,
		       vector_embed_model_id, rerank_model_id, search_engine_enable, english_name, options,
		       created, created_by, modified, modified_by
		FROM tb_document_collection
		WHERE slug = ?
	`

	return r.scanOne(ctx, query, slug)
}

// List 获取知识库列表
func (r *DocumentCollectionRepository) List(ctx context.Context, tenantID int64) ([]*entity.DocumentCollection, error) {
	query := `
		SELECT id, alias, dept_id, tenant_id, icon, title, description, slug,
		       vector_store_enable, vector_store_type, vector_store_collection, vector_store_config,
		       vector_embed_model_id, rerank_model_id, search_engine_enable, english_name, options,
		       created, created_by, modified, modified_by
		FROM tb_document_collection
		WHERE tenant_id = ?
		ORDER BY created DESC
	`

	return r.scanList(ctx, query, tenantID)
}

// GetDocumentCount 获取知识库的文档数量
func (r *DocumentCollectionRepository) GetDocumentCount(ctx context.Context, collectionID int64) (int, error) {
	query := `SELECT COUNT(*) FROM tb_document WHERE collection_id = ?`
	var count int
	err := r.db.QueryRowContext(ctx, query, collectionID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ExistsBotRelation 检查是否存在 Bot 关联
func (r *DocumentCollectionRepository) ExistsBotRelation(ctx context.Context, collectionID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM tb_bot_document_collection WHERE document_collection_id = ?`
	var count int
	err := r.db.QueryRowContext(ctx, query, collectionID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// scanOne 扫描单条记录
func (r *DocumentCollectionRepository) scanOne(ctx context.Context, query string, args ...interface{}) (*entity.DocumentCollection, error) {
	var dc entity.DocumentCollection
	var alias, icon, title, description, slug, vectorStoreType, vectorStoreCollection, vectorStoreConfig sql.NullString
	var englishName, options sql.NullString
	var vectorEmbedModelID, rerankModelID, createdBy, modifiedBy sql.NullInt64
	var created, modified sql.NullTime

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&dc.ID, &alias, &dc.DeptID, &dc.TenantID, &icon, &title, &description, &slug,
		&dc.VectorStoreEnable, &vectorStoreType, &vectorStoreCollection, &vectorStoreConfig,
		&vectorEmbedModelID, &rerankModelID, &dc.SearchEngineEnable, &englishName, &options,
		&created, &createdBy, &modified, &modifiedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	dc.Alias = alias.String
	dc.Icon = icon.String
	dc.Title = title.String
	dc.Description = description.String
	dc.Slug = slug.String
	dc.VectorStoreType = vectorStoreType.String
	dc.VectorStoreCollection = vectorStoreCollection.String
	dc.VectorStoreConfig = vectorStoreConfig.String
	dc.EnglishName = englishName.String
	dc.Options = options.String
	if vectorEmbedModelID.Valid {
		dc.VectorEmbedModelID = &vectorEmbedModelID.Int64
	}
	if rerankModelID.Valid {
		dc.RerankModelID = &rerankModelID.Int64
	}
	if created.Valid {
		dc.Created = &created.Time
	}
	if createdBy.Valid {
		dc.CreatedBy = &createdBy.Int64
	}
	if modified.Valid {
		dc.Modified = &modified.Time
	}
	if modifiedBy.Valid {
		dc.ModifiedBy = &modifiedBy.Int64
	}

	return &dc, nil
}

// scanList 扫描多条记录
func (r *DocumentCollectionRepository) scanList(ctx context.Context, query string, args ...interface{}) ([]*entity.DocumentCollection, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.DocumentCollection
	for rows.Next() {
		var dc entity.DocumentCollection
		var alias, icon, title, description, slug, vectorStoreType, vectorStoreCollection, vectorStoreConfig sql.NullString
		var englishName, options sql.NullString
		var vectorEmbedModelID, rerankModelID, createdBy, modifiedBy sql.NullInt64
		var created, modified sql.NullTime

		err := rows.Scan(
			&dc.ID, &alias, &dc.DeptID, &dc.TenantID, &icon, &title, &description, &slug,
			&dc.VectorStoreEnable, &vectorStoreType, &vectorStoreCollection, &vectorStoreConfig,
			&vectorEmbedModelID, &rerankModelID, &dc.SearchEngineEnable, &englishName, &options,
			&created, &createdBy, &modified, &modifiedBy,
		)
		if err != nil {
			return nil, err
		}

		dc.Alias = alias.String
		dc.Icon = icon.String
		dc.Title = title.String
		dc.Description = description.String
		dc.Slug = slug.String
		dc.VectorStoreType = vectorStoreType.String
		dc.VectorStoreCollection = vectorStoreCollection.String
		dc.VectorStoreConfig = vectorStoreConfig.String
		dc.EnglishName = englishName.String
		dc.Options = options.String
		if vectorEmbedModelID.Valid {
			dc.VectorEmbedModelID = &vectorEmbedModelID.Int64
		}
		if rerankModelID.Valid {
			dc.RerankModelID = &rerankModelID.Int64
		}
		if created.Valid {
			dc.Created = &created.Time
		}
		if createdBy.Valid {
			dc.CreatedBy = &createdBy.Int64
		}
		if modified.Valid {
			dc.Modified = &modified.Time
		}
		if modifiedBy.Valid {
			dc.ModifiedBy = &modifiedBy.Int64
		}

		list = append(list, &dc)
	}

	return list, nil
}

// ========================== BotDocumentCollection ==========================

// ListByBotID 获取 Bot 关联的知识库列表
func (r *DocumentCollectionRepository) ListByBotID(ctx context.Context, botID int64) ([]*entity.BotDocumentCollection, error) {
	query := `
		SELECT bdc.id, bdc.bot_id, bdc.document_collection_id,
		       dc.id, dc.alias, dc.dept_id, dc.tenant_id, dc.icon, dc.title, dc.description, dc.slug,
		       dc.vector_store_enable, dc.vector_store_type, dc.vector_store_collection, dc.vector_store_config,
		       dc.vector_embed_model_id, dc.rerank_model_id, dc.search_engine_enable, dc.english_name, dc.options,
		       dc.created, dc.created_by, dc.modified, dc.modified_by
		FROM tb_bot_document_collection bdc
		LEFT JOIN tb_document_collection dc ON bdc.document_collection_id = dc.id
		WHERE bdc.bot_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, botID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.BotDocumentCollection
	for rows.Next() {
		var bdc entity.BotDocumentCollection
		var dc entity.DocumentCollection
		var alias, icon, title, description, slug, vectorStoreType, vectorStoreCollection, vectorStoreConfig sql.NullString
		var englishName, options sql.NullString
		var vectorEmbedModelID, rerankModelID, createdBy, modifiedBy sql.NullInt64
		var created, modified sql.NullTime

		err := rows.Scan(
			&bdc.ID, &bdc.BotID, &bdc.KnowledgeID,
			&dc.ID, &alias, &dc.DeptID, &dc.TenantID, &icon, &title, &description, &slug,
			&dc.VectorStoreEnable, &vectorStoreType, &vectorStoreCollection, &vectorStoreConfig,
			&vectorEmbedModelID, &rerankModelID, &dc.SearchEngineEnable, &englishName, &options,
			&created, &createdBy, &modified, &modifiedBy,
		)
		if err != nil {
			return nil, err
		}

		dc.Alias = alias.String
		dc.Icon = icon.String
		dc.Title = title.String
		dc.Description = description.String
		dc.Slug = slug.String
		dc.VectorStoreType = vectorStoreType.String
		dc.VectorStoreCollection = vectorStoreCollection.String
		dc.VectorStoreConfig = vectorStoreConfig.String
		dc.EnglishName = englishName.String
		dc.Options = options.String
		if vectorEmbedModelID.Valid {
			dc.VectorEmbedModelID = &vectorEmbedModelID.Int64
		}
		if rerankModelID.Valid {
			dc.RerankModelID = &rerankModelID.Int64
		}
		if created.Valid {
			dc.Created = &created.Time
		}
		if createdBy.Valid {
			dc.CreatedBy = &createdBy.Int64
		}
		if modified.Valid {
			dc.Modified = &modified.Time
		}
		if modifiedBy.Valid {
			dc.ModifiedBy = &modifiedBy.Int64
		}

		bdc.DocumentCollection = &dc
		list = append(list, &bdc)
	}

	return list, nil
}

// GetBotKnowledgeIDs 获取 Bot 关联的知识库 ID 列表
func (r *DocumentCollectionRepository) GetBotKnowledgeIDs(ctx context.Context, botID int64) ([]int64, error) {
	query := `SELECT document_collection_id FROM tb_bot_document_collection WHERE bot_id = ?`
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

// UpdateBotKnowledges 更新 Bot-知识库关联
func (r *DocumentCollectionRepository) UpdateBotKnowledges(ctx context.Context, botID int64, knowledgeIDs []int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除旧关联
	_, err = tx.ExecContext(ctx, `DELETE FROM tb_bot_document_collection WHERE bot_id = ?`, botID)
	if err != nil {
		return err
	}

	// 插入新关联
	for _, knowledgeID := range knowledgeIDs {
		id, _ := snowflake.GenerateID()
		_, err = tx.ExecContext(ctx,
			`INSERT INTO tb_bot_document_collection (id, bot_id, document_collection_id) VALUES (?, ?, ?)`,
			id, botID, knowledgeID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// DeleteBotKnowledge 删除单个 Bot-知识库关联
func (r *DocumentCollectionRepository) DeleteBotKnowledge(ctx context.Context, botID, knowledgeID int64) error {
	query := `DELETE FROM tb_bot_document_collection WHERE bot_id = ? AND document_collection_id = ?`
	_, err := r.db.ExecContext(ctx, query, botID, knowledgeID)
	return err
}
