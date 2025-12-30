package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/pkg/snowflake"
)

// DocumentRepository 文档数据访问层
type DocumentRepository struct {
	db *sql.DB
}

// NewDocumentRepository 创建 DocumentRepository
func NewDocumentRepository() *DocumentRepository {
	return &DocumentRepository{
		db: GetDB(),
	}
}

// Create 创建文档
func (r *DocumentRepository) Create(ctx context.Context, doc *entity.Document) error {
	if doc.ID == 0 {
		doc.ID, _ = snowflake.GenerateID()
	}
	now := time.Now()
	doc.Created = &now
	doc.Modified = &now

	query := `
		INSERT INTO tb_document
		(id, collection_id, document_type, document_path, title, content, content_type, slug, order_no, options, created, created_by, modified, modified_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		doc.ID, doc.CollectionID, doc.DocumentType, doc.DocumentPath, doc.Title,
		doc.Content, doc.ContentType, doc.Slug, doc.OrderNo, doc.Options,
		doc.Created, doc.CreatedBy, doc.Modified, doc.ModifiedBy,
	)
	return err
}

// Update 更新文档
func (r *DocumentRepository) Update(ctx context.Context, doc *entity.Document) error {
	now := time.Now()
	doc.Modified = &now

	query := `
		UPDATE tb_document SET
			document_type = ?, document_path = ?, title = ?, content = ?, content_type = ?,
			slug = ?, order_no = ?, options = ?, modified = ?, modified_by = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		doc.DocumentType, doc.DocumentPath, doc.Title, doc.Content, doc.ContentType,
		doc.Slug, doc.OrderNo, doc.Options, doc.Modified, doc.ModifiedBy, doc.ID,
	)
	return err
}

// Delete 删除文档
func (r *DocumentRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tb_document WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// DeleteByCollectionID 删除知识库下的所有文档
func (r *DocumentRepository) DeleteByCollectionID(ctx context.Context, collectionID int64) error {
	query := `DELETE FROM tb_document WHERE collection_id = ?`
	_, err := r.db.ExecContext(ctx, query, collectionID)
	return err
}

// GetByID 根据 ID 获取文档
func (r *DocumentRepository) GetByID(ctx context.Context, id int64) (*entity.Document, error) {
	query := `
		SELECT id, collection_id, document_type, document_path, title, content, content_type, slug, order_no, options, created, created_by, modified, modified_by
		FROM tb_document
		WHERE id = ?
	`

	return r.scanOne(ctx, query, id)
}

// ListByCollectionID 获取知识库下的文档列表
func (r *DocumentRepository) ListByCollectionID(ctx context.Context, collectionID int64) ([]*entity.Document, error) {
	query := `
		SELECT id, collection_id, document_type, document_path, title, content, content_type, slug, order_no, options, created, created_by, modified, modified_by
		FROM tb_document
		WHERE collection_id = ?
		ORDER BY order_no ASC, created DESC
	`

	return r.scanList(ctx, query, collectionID)
}

// ListByCollectionIDPaged 分页获取文档列表
func (r *DocumentRepository) ListByCollectionIDPaged(ctx context.Context, collectionID int64, title string, pageNumber, pageSize int) ([]*entity.Document, int64, error) {
	// 计算总数
	countQuery := `SELECT COUNT(*) FROM tb_document WHERE collection_id = ?`
	countArgs := []interface{}{collectionID}
	if title != "" {
		countQuery += ` AND title LIKE ?`
		countArgs = append(countArgs, "%"+title+"%")
	}

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	query := `
		SELECT id, collection_id, document_type, document_path, title, content, content_type, slug, order_no, options, created, created_by, modified, modified_by
		FROM tb_document
		WHERE collection_id = ?
	`
	args := []interface{}{collectionID}
	if title != "" {
		query += ` AND title LIKE ?`
		args = append(args, "%"+title+"%")
	}
	query += ` ORDER BY order_no ASC, created DESC LIMIT ? OFFSET ?`
	args = append(args, pageSize, (pageNumber-1)*pageSize)

	docs, err := r.scanList(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return docs, total, nil
}

// UpdateOrderNo 更新文档排序
func (r *DocumentRepository) UpdateOrderNo(ctx context.Context, id int64, orderNo int) error {
	query := `UPDATE tb_document SET order_no = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, orderNo, id)
	return err
}

// scanOne 扫描单条记录
func (r *DocumentRepository) scanOne(ctx context.Context, query string, args ...interface{}) (*entity.Document, error) {
	var doc entity.Document
	var documentType, documentPath, title, content, contentType, slug, options sql.NullString
	var orderNo sql.NullInt32
	var createdBy, modifiedBy sql.NullInt64
	var created, modified sql.NullTime

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&doc.ID, &doc.CollectionID, &documentType, &documentPath, &title, &content,
		&contentType, &slug, &orderNo, &options, &created, &createdBy, &modified, &modifiedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	doc.DocumentType = documentType.String
	doc.DocumentPath = documentPath.String
	doc.Title = title.String
	doc.Content = content.String
	doc.ContentType = contentType.String
	doc.Slug = slug.String
	doc.Options = options.String
	if orderNo.Valid {
		o := int(orderNo.Int32)
		doc.OrderNo = &o
	}
	if created.Valid {
		doc.Created = &created.Time
	}
	if createdBy.Valid {
		doc.CreatedBy = &createdBy.Int64
	}
	if modified.Valid {
		doc.Modified = &modified.Time
	}
	if modifiedBy.Valid {
		doc.ModifiedBy = &modifiedBy.Int64
	}

	return &doc, nil
}

// scanList 扫描多条记录
func (r *DocumentRepository) scanList(ctx context.Context, query string, args ...interface{}) ([]*entity.Document, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.Document
	for rows.Next() {
		var doc entity.Document
		var documentType, documentPath, title, content, contentType, slug, options sql.NullString
		var orderNo sql.NullInt32
		var createdBy, modifiedBy sql.NullInt64
		var created, modified sql.NullTime

		err := rows.Scan(
			&doc.ID, &doc.CollectionID, &documentType, &documentPath, &title, &content,
			&contentType, &slug, &orderNo, &options, &created, &createdBy, &modified, &modifiedBy,
		)
		if err != nil {
			return nil, err
		}

		doc.DocumentType = documentType.String
		doc.DocumentPath = documentPath.String
		doc.Title = title.String
		doc.Content = content.String
		doc.ContentType = contentType.String
		doc.Slug = slug.String
		doc.Options = options.String
		if orderNo.Valid {
			o := int(orderNo.Int32)
			doc.OrderNo = &o
		}
		if created.Valid {
			doc.Created = &created.Time
		}
		if createdBy.Valid {
			doc.CreatedBy = &createdBy.Int64
		}
		if modified.Valid {
			doc.Modified = &modified.Time
		}
		if modifiedBy.Valid {
			doc.ModifiedBy = &modifiedBy.Int64
		}

		list = append(list, &doc)
	}

	return list, nil
}

// ========================== DocumentChunk ==========================

// CreateChunk 创建文档分块
func (r *DocumentRepository) CreateChunk(ctx context.Context, chunk *entity.DocumentChunk) error {
	if chunk.ID == 0 {
		chunk.ID, _ = snowflake.GenerateID()
	}

	query := `
		INSERT INTO tb_document_chunk (id, document_id, document_collection_id, content, sorting)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		chunk.ID, chunk.DocumentID, chunk.DocumentCollectionID, chunk.Content, chunk.Sorting,
	)
	return err
}

// DeleteChunksByDocumentID 删除文档的所有分块
func (r *DocumentRepository) DeleteChunksByDocumentID(ctx context.Context, documentID int64) error {
	query := `DELETE FROM tb_document_chunk WHERE document_id = ?`
	_, err := r.db.ExecContext(ctx, query, documentID)
	return err
}

// DeleteChunksByCollectionID 删除知识库的所有分块
func (r *DocumentRepository) DeleteChunksByCollectionID(ctx context.Context, collectionID int64) error {
	query := `DELETE FROM tb_document_chunk WHERE document_collection_id = ?`
	_, err := r.db.ExecContext(ctx, query, collectionID)
	return err
}

// ListChunksByDocumentID 获取文档的所有分块
func (r *DocumentRepository) ListChunksByDocumentID(ctx context.Context, documentID int64) ([]*entity.DocumentChunk, error) {
	query := `
		SELECT id, document_id, document_collection_id, content, sorting
		FROM tb_document_chunk
		WHERE document_id = ?
		ORDER BY sorting ASC
	`

	rows, err := r.db.QueryContext(ctx, query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.DocumentChunk
	for rows.Next() {
		var chunk entity.DocumentChunk
		var content sql.NullString
		err := rows.Scan(&chunk.ID, &chunk.DocumentID, &chunk.DocumentCollectionID, &content, &chunk.Sorting)
		if err != nil {
			return nil, err
		}
		chunk.Content = content.String
		list = append(list, &chunk)
	}

	return list, nil
}

// ListChunksByCollectionID 获取知识库的所有分块
func (r *DocumentRepository) ListChunksByCollectionID(ctx context.Context, collectionID int64) ([]*entity.DocumentChunk, error) {
	query := `
		SELECT id, document_id, document_collection_id, content, sorting
		FROM tb_document_chunk
		WHERE document_collection_id = ?
		ORDER BY document_id, sorting ASC
	`

	rows, err := r.db.QueryContext(ctx, query, collectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.DocumentChunk
	for rows.Next() {
		var chunk entity.DocumentChunk
		var content sql.NullString
		err := rows.Scan(&chunk.ID, &chunk.DocumentID, &chunk.DocumentCollectionID, &content, &chunk.Sorting)
		if err != nil {
			return nil, err
		}
		chunk.Content = content.String
		list = append(list, &chunk)
	}

	return list, nil
}

// GetChunkIDs 获取文档的分块 ID 列表
func (r *DocumentRepository) GetChunkIDs(ctx context.Context, documentID int64) ([]int64, error) {
	query := `SELECT id FROM tb_document_chunk WHERE document_id = ?`
	rows, err := r.db.QueryContext(ctx, query, documentID)
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

// ========================== DocumentHistory ==========================

// CreateHistory 创建文档历史记录
func (r *DocumentRepository) CreateHistory(ctx context.Context, history *entity.DocumentHistory) error {
	query := `
		INSERT INTO tb_document_history
		(document_id, old_title, new_title, old_content, new_content, old_document_type, new_document_type, created, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		history.DocumentID, history.OldTitle, history.NewTitle,
		history.OldContent, history.NewContent, history.OldDocumentType, history.NewDocumentType,
		now, history.CreatedBy,
	)
	return err
}
