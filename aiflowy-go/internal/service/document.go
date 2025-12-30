package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aiflowy/aiflowy-go/internal/config"
	"github.com/aiflowy/aiflowy-go/internal/dto"
	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/internal/repository"
	"github.com/aiflowy/aiflowy-go/internal/service/rag"
)

// DocumentService 文档业务逻辑
type DocumentService struct {
	repo           *repository.DocumentRepository
	collectionRepo *repository.DocumentCollectionRepository
}

// NewDocumentService 创建 DocumentService
func NewDocumentService() *DocumentService {
	return &DocumentService{
		repo:           repository.NewDocumentRepository(),
		collectionRepo: repository.NewDocumentCollectionRepository(),
	}
}

// List 获取文档列表
func (s *DocumentService) List(ctx context.Context, collectionID string) ([]*entity.Document, error) {
	collectionIDInt := parseID(collectionID)
	return s.repo.ListByCollectionID(ctx, collectionIDInt)
}

// ListPaged 分页获取文档列表
func (s *DocumentService) ListPaged(ctx context.Context, collectionID, title string, pageNumber, pageSize int) (*dto.DocumentListResponse, error) {
	collectionIDInt := parseID(collectionID)

	// 默认分页参数
	if pageNumber <= 0 {
		pageNumber = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	docs, total, err := s.repo.ListByCollectionIDPaged(ctx, collectionIDInt, title, pageNumber, pageSize)
	if err != nil {
		return nil, err
	}

	return &dto.DocumentListResponse{
		Total:    total,
		PageNo:   pageNumber,
		PageSize: pageSize,
		List:     docs,
	}, nil
}

// GetByID 根据 ID 获取文档
func (s *DocumentService) GetByID(ctx context.Context, id string) (*entity.Document, error) {
	idInt := parseID(id)
	return s.repo.GetByID(ctx, idInt)
}

// Save 保存文档
func (s *DocumentService) Save(ctx context.Context, req *dto.DocumentSaveRequest, userID int64) (*entity.Document, error) {
	// 验证知识库是否存在
	collectionID := parseID(req.CollectionID)
	collection, err := s.collectionRepo.GetByID(ctx, collectionID)
	if err != nil {
		return nil, err
	}
	if collection == nil {
		return nil, fmt.Errorf("知识库不存在")
	}

	var doc *entity.Document
	if req.ID != "" {
		// 更新
		idInt := parseID(req.ID)
		doc, err = s.repo.GetByID(ctx, idInt)
		if err != nil {
			return nil, err
		}
		if doc == nil {
			return nil, fmt.Errorf("文档不存在")
		}

		// 保存历史记录
		history := &entity.DocumentHistory{
			DocumentID:      doc.ID,
			OldTitle:        doc.Title,
			NewTitle:        req.Title,
			OldContent:      doc.Content,
			NewContent:      req.Content,
			OldDocumentType: doc.DocumentType,
			NewDocumentType: req.DocumentType,
			CreatedBy:       &userID,
		}
		s.repo.CreateHistory(ctx, history)

		s.fillEntity(doc, req)
		doc.ModifiedBy = &userID
		if err := s.repo.Update(ctx, doc); err != nil {
			return nil, err
		}
	} else {
		// 创建
		doc = &entity.Document{
			CollectionID: collectionID,
			CreatedBy:    &userID,
		}
		s.fillEntity(doc, req)
		if err := s.repo.Create(ctx, doc); err != nil {
			return nil, err
		}
	}

	return doc, nil
}

// Delete 删除文档
func (s *DocumentService) Delete(ctx context.Context, id string) error {
	idInt := parseID(id)

	// 删除文档分块
	s.repo.DeleteChunksByDocumentID(ctx, idInt)
	// 删除文档
	return s.repo.Delete(ctx, idInt)
}

// UpdatePosition 更新文档排序
func (s *DocumentService) UpdatePosition(ctx context.Context, id string, orderNo int, userID int64) error {
	idInt := parseID(id)

	// 获取文档
	doc, err := s.repo.GetByID(ctx, idInt)
	if err != nil {
		return err
	}
	if doc == nil {
		return fmt.Errorf("文档不存在")
	}

	// 获取同知识库下的所有文档
	docs, err := s.repo.ListByCollectionID(ctx, doc.CollectionID)
	if err != nil {
		return err
	}

	// 移除当前文档
	var filteredDocs []*entity.Document
	for _, d := range docs {
		if d.ID != doc.ID {
			filteredDocs = append(filteredDocs, d)
		}
	}

	// 调整排序位置
	if orderNo < 0 {
		orderNo = 0
	}
	if orderNo >= len(filteredDocs) {
		filteredDocs = append(filteredDocs, doc)
	} else {
		// 在指定位置插入
		newDocs := make([]*entity.Document, 0, len(filteredDocs)+1)
		newDocs = append(newDocs, filteredDocs[:orderNo]...)
		newDocs = append(newDocs, doc)
		newDocs = append(newDocs, filteredDocs[orderNo:]...)
		filteredDocs = newDocs
	}

	// 更新所有文档的排序
	for i, d := range filteredDocs {
		if err := s.repo.UpdateOrderNo(ctx, d.ID, i); err != nil {
			return err
		}
	}

	return nil
}

// fillEntity 填充实体字段
func (s *DocumentService) fillEntity(doc *entity.Document, req *dto.DocumentSaveRequest) {
	if req.CollectionID != "" {
		doc.CollectionID = parseID(req.CollectionID)
	}
	doc.DocumentType = req.DocumentType
	doc.DocumentPath = req.DocumentPath
	doc.Title = req.Title
	doc.Content = req.Content
	doc.ContentType = req.ContentType
	doc.Slug = req.Slug
	doc.OrderNo = req.OrderNo
	if req.Options != "" {
		doc.Options = req.Options
	}
}

// ========================== 文档分块 ==========================

// SaveChunks 保存文档分块
func (s *DocumentService) SaveChunks(ctx context.Context, documentID, collectionID int64, chunks []string) error {
	// 删除旧分块
	if err := s.repo.DeleteChunksByDocumentID(ctx, documentID); err != nil {
		return err
	}

	// 创建新分块
	for i, content := range chunks {
		chunk := &entity.DocumentChunk{
			DocumentID:           documentID,
			DocumentCollectionID: collectionID,
			Content:              content,
			Sorting:              i,
		}
		if err := s.repo.CreateChunk(ctx, chunk); err != nil {
			return err
		}
	}

	return nil
}

// GetChunks 获取文档分块
func (s *DocumentService) GetChunks(ctx context.Context, documentID string) ([]*entity.DocumentChunk, error) {
	documentIDInt := parseID(documentID)
	return s.repo.ListChunksByDocumentID(ctx, documentIDInt)
}

// TextSplit 文本拆分
func (s *DocumentService) TextSplit(ctx context.Context, req *dto.TextSplitRequest, userID int64) (interface{}, error) {
	// 获取知识库
	collectionID := parseID(req.KnowledgeID)
	collection, err := s.collectionRepo.GetByID(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("获取知识库失败: %w", err)
	}
	if collection == nil {
		return nil, fmt.Errorf("知识库不存在")
	}

	// 读取文件内容
	content, err := s.readFileContent(req.FilePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 设置默认分块参数
	chunkSize := req.ChunkSize
	if chunkSize <= 0 {
		chunkSize = 500
	}
	overlapSize := req.OverlapSize
	if overlapSize < 0 {
		overlapSize = 0
	}

	// 获取分块器
	splitter := rag.GetDocumentSplitter(req.SplitterName, chunkSize, overlapSize, req.Regex)

	// 分块
	chunks := splitter.Split(content)

	// 根据操作类型处理
	if req.Operation == "saveText" {
		// 保存文档和分块
		return s.saveTextResult(ctx, req, collection, content, chunks, userID)
	}

	// 预览模式：分页返回分块
	pageNumber := req.PageNumber
	if pageNumber <= 0 {
		pageNumber = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 分页
	total := len(chunks)
	startIdx := (pageNumber - 1) * pageSize
	endIdx := startIdx + pageSize
	if startIdx >= total {
		return map[string]interface{}{
			"total":       total,
			"previewData": []string{},
		}, nil
	}
	if endIdx > total {
		endIdx = total
	}

	previewChunks := chunks[startIdx:endIdx]

	// 构造预览数据
	var previewData []map[string]interface{}
	for i, chunk := range previewChunks {
		previewData = append(previewData, map[string]interface{}{
			"id":      startIdx + i + 1,
			"content": chunk,
			"sorting": startIdx + i + 1,
		})
	}

	return map[string]interface{}{
		"total":       total,
		"previewData": previewData,
		"aiDocumentData": map[string]interface{}{
			"title":        req.FileOriginName,
			"documentType": getFileExtension(req.FilePath),
			"documentPath": req.FilePath,
			"content":      content,
			"chunkSize":    chunkSize,
			"overlapSize":  overlapSize,
		},
	}, nil
}

// saveTextResult 保存文本分割结果
func (s *DocumentService) saveTextResult(ctx context.Context, req *dto.TextSplitRequest, collection *entity.DocumentCollection, content string, chunks []string, userID int64) (interface{}, error) {
	collectionID := collection.ID

	// 创建文档
	doc := &entity.Document{
		CollectionID: collectionID,
		Title:        req.FileOriginName,
		DocumentType: getFileExtension(req.FilePath),
		DocumentPath: req.FilePath,
		Content:      content,
		CreatedBy:    &userID,
	}

	if err := s.repo.Create(ctx, doc); err != nil {
		return nil, fmt.Errorf("创建文档失败: %w", err)
	}

	// 保存分块到数据库
	var chunkEntities []*entity.DocumentChunk
	for i, chunkContent := range chunks {
		chunk := &entity.DocumentChunk{
			DocumentID:           doc.ID,
			DocumentCollectionID: collectionID,
			Content:              chunkContent,
			Sorting:              i + 1,
		}
		if err := s.repo.CreateChunk(ctx, chunk); err != nil {
			return nil, fmt.Errorf("创建分块失败: %w", err)
		}
		chunkEntities = append(chunkEntities, chunk)
	}

	// 如果启用了向量存储，进行向量化
	if collection.VectorStoreEnable && collection.VectorEmbedModelID != nil {
		ragService := rag.GetRAGService()
		if err := ragService.IndexDocumentChunks(ctx, collection, chunkEntities); err != nil {
			// 向量化失败不影响保存，只记录日志
			fmt.Printf("Warning: failed to index document chunks: %v\n", err)
		}
	}

	return map[string]interface{}{
		"id":         doc.ID,
		"title":      doc.Title,
		"chunkCount": len(chunks),
	}, nil
}

// readFileContent 读取文件内容
func (s *DocumentService) readFileContent(filePath string) (string, error) {
	// 获取完整路径
	cfg := config.GetConfig()
	rootPath := "./uploads"
	if cfg != nil && cfg.Storage.LocalRoot != "" {
		rootPath = cfg.Storage.LocalRoot
	}
	fullPath := filepath.Join(rootPath, filePath)

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("文件不存在: %s", filePath)
	}

	// 读取文件
	file, err := os.Open(fullPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// getFileExtension 获取文件扩展名
func getFileExtension(filePath string) string {
	ext := filepath.Ext(filePath)
	if len(ext) > 0 && ext[0] == '.' {
		return ext[1:]
	}
	return ext
}
