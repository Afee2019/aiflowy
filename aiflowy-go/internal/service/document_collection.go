package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aiflowy/aiflowy-go/internal/dto"
	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/internal/repository"
)

// DocumentCollectionService 知识库业务逻辑
type DocumentCollectionService struct {
	repo        *repository.DocumentCollectionRepository
	docRepo     *repository.DocumentRepository
	modelRepo   *repository.ModelRepository
}

// NewDocumentCollectionService 创建 DocumentCollectionService
func NewDocumentCollectionService() *DocumentCollectionService {
	return &DocumentCollectionService{
		repo:      repository.NewDocumentCollectionRepository(),
		docRepo:   repository.NewDocumentRepository(),
		modelRepo: repository.NewModelRepository(repository.GetDB()),
	}
}

// List 获取知识库列表
func (s *DocumentCollectionService) List(ctx context.Context, tenantID int64) ([]*entity.DocumentCollection, error) {
	list, err := s.repo.List(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 加载文档数量
	for _, dc := range list {
		count, _ := s.repo.GetDocumentCount(ctx, dc.ID)
		dc.DocumentCount = count
	}

	return list, nil
}

// GetByID 根据 ID 获取知识库
func (s *DocumentCollectionService) GetByID(ctx context.Context, id string) (*entity.DocumentCollection, error) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的知识库 ID: %s", id)
	}

	dc, err := s.repo.GetByID(ctx, idInt)
	if err != nil {
		return nil, err
	}
	if dc == nil {
		return nil, nil
	}

	// 加载文档数量
	count, _ := s.repo.GetDocumentCount(ctx, dc.ID)
	dc.DocumentCount = count

	// 加载模型信息
	if dc.VectorEmbedModelID != nil {
		model, _ := s.modelRepo.GetModelByID(ctx, *dc.VectorEmbedModelID)
		dc.EmbedModel = model
	}
	if dc.RerankModelID != nil {
		model, _ := s.modelRepo.GetModelByID(ctx, *dc.RerankModelID)
		dc.RerankModel = model
	}

	return dc, nil
}

// GetDetail 获取知识库详情 (带关联数据)
func (s *DocumentCollectionService) GetDetail(ctx context.Context, id string) (*entity.DocumentCollection, error) {
	return s.GetByID(ctx, id)
}

// Save 保存知识库
func (s *DocumentCollectionService) Save(ctx context.Context, req *dto.DocumentCollectionSaveRequest, tenantID, userID, deptID int64) (*entity.DocumentCollection, error) {
	// 检查别名是否唯一
	if req.Alias != "" {
		existing, err := s.repo.GetByAlias(ctx, req.Alias)
		if err != nil {
			return nil, err
		}
		if existing != nil && (req.ID == "" || existing.ID != parseID(req.ID)) {
			return nil, fmt.Errorf("别名已存在")
		}
	}

	var dc *entity.DocumentCollection
	if req.ID != "" {
		// 更新
		idInt := parseID(req.ID)
		dc, _ = s.repo.GetByID(ctx, idInt)
		if dc == nil {
			return nil, fmt.Errorf("知识库不存在")
		}
		s.fillEntity(dc, req)
		dc.ModifiedBy = &userID
		if err := s.repo.Update(ctx, dc); err != nil {
			return nil, err
		}
	} else {
		// 创建
		dc = &entity.DocumentCollection{
			TenantID:  tenantID,
			DeptID:    deptID,
			CreatedBy: &userID,
		}
		s.fillEntity(dc, req)
		// 默认值
		if dc.Options == "" {
			dc.Options = `{"canUpdateEmbedding":true}`
		}
		if err := s.repo.Create(ctx, dc); err != nil {
			return nil, err
		}
	}

	return dc, nil
}

// Delete 删除知识库
func (s *DocumentCollectionService) Delete(ctx context.Context, id string) error {
	idInt := parseID(id)

	// 检查是否存在 Bot 关联
	exists, err := s.repo.ExistsBotRelation(ctx, idInt)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("此知识库还关联着Bot，请先取消关联")
	}

	// 删除文档分块
	s.docRepo.DeleteChunksByCollectionID(ctx, idInt)
	// 删除文档
	s.docRepo.DeleteByCollectionID(ctx, idInt)
	// 删除知识库
	return s.repo.Delete(ctx, idInt)
}

// fillEntity 填充实体字段
func (s *DocumentCollectionService) fillEntity(dc *entity.DocumentCollection, req *dto.DocumentCollectionSaveRequest) {
	dc.Alias = req.Alias
	dc.Icon = req.Icon
	dc.Title = req.Title
	dc.Description = req.Description
	dc.Slug = req.Slug
	if req.VectorStoreEnable != nil {
		dc.VectorStoreEnable = *req.VectorStoreEnable
	}
	dc.VectorStoreType = req.VectorStoreType
	dc.VectorStoreCollection = req.VectorStoreCollection
	dc.VectorStoreConfig = req.VectorStoreConfig
	if req.VectorEmbedModelID != "" {
		id := parseID(req.VectorEmbedModelID)
		dc.VectorEmbedModelID = &id
	} else {
		dc.VectorEmbedModelID = nil
	}
	if req.RerankModelID != "" {
		id := parseID(req.RerankModelID)
		dc.RerankModelID = &id
	} else {
		dc.RerankModelID = nil
	}
	if req.SearchEngineEnable != nil {
		dc.SearchEngineEnable = *req.SearchEngineEnable
	}
	dc.EnglishName = req.EnglishName
	if req.Options != "" {
		dc.Options = req.Options
	}
}

// parseID 解析 ID 字符串
func parseID(id string) int64 {
	idInt, _ := strconv.ParseInt(id, 10, 64)
	return idInt
}

// ========================== Bot-知识库关联 ==========================

// ListByBotID 获取 Bot 关联的知识库列表
func (s *DocumentCollectionService) ListByBotID(ctx context.Context, botID string) ([]*entity.BotDocumentCollection, error) {
	botIDInt := parseID(botID)
	return s.repo.ListByBotID(ctx, botIDInt)
}

// GetBotKnowledgeIDs 获取 Bot 关联的知识库 ID 列表
func (s *DocumentCollectionService) GetBotKnowledgeIDs(ctx context.Context, botID string) ([]string, error) {
	botIDInt := parseID(botID)
	ids, err := s.repo.GetBotKnowledgeIDs(ctx, botIDInt)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = strconv.FormatInt(id, 10)
	}
	return result, nil
}

// UpdateBotKnowledges 更新 Bot-知识库关联
func (s *DocumentCollectionService) UpdateBotKnowledges(ctx context.Context, req *dto.BotDocumentCollectionUpdateRequest) error {
	botIDInt := parseID(req.BotID)
	knowledgeIDs := make([]int64, len(req.KnowledgeIDs))
	for i, id := range req.KnowledgeIDs {
		knowledgeIDs[i] = parseID(id)
	}
	return s.repo.UpdateBotKnowledges(ctx, botIDInt, knowledgeIDs)
}

// DeleteBotKnowledge 删除单个 Bot-知识库关联
func (s *DocumentCollectionService) DeleteBotKnowledge(ctx context.Context, botID, knowledgeID string) error {
	return s.repo.DeleteBotKnowledge(ctx, parseID(botID), parseID(knowledgeID))
}

// ========================== 知识库检索 (RAG) ==========================

// SearchResult 检索结果
type SearchResult struct {
	ID      int64   `json:"id,string"`
	Content string  `json:"content"`
	Score   float64 `json:"score,omitempty"`
}

// SearchByCollectionID 在知识库中检索
func (s *DocumentCollectionService) SearchByCollectionID(ctx context.Context, collectionID int64, query string, topK int) []*SearchResult {
	if query == "" || collectionID == 0 {
		return nil
	}

	if topK <= 0 {
		topK = 5
	}

	// 获取知识库的所有文档分块
	chunks, err := s.docRepo.ListChunksByCollectionID(ctx, collectionID)
	if err != nil || len(chunks) == 0 {
		return nil
	}

	// 简单的关键词匹配搜索
	var results []*SearchResult
	queryRunes := []rune(query)

	for _, chunk := range chunks {
		if chunk.Content == "" {
			continue
		}

		// 简单匹配：检查内容是否包含查询词
		if containsQuery(chunk.Content, query, queryRunes) {
			results = append(results, &SearchResult{
				ID:      chunk.ID,
				Content: chunk.Content,
				Score:   0.5,
			})
		}

		if len(results) >= topK {
			break
		}
	}

	return results
}

// containsQuery 检查内容是否包含查询词
func containsQuery(content, query string, queryRunes []rune) bool {
	contentRunes := []rune(content)

	// 1. 直接包含
	if len(queryRunes) <= len(contentRunes) {
		for i := 0; i <= len(contentRunes)-len(queryRunes); i++ {
			match := true
			for j := 0; j < len(queryRunes); j++ {
				if contentRunes[i+j] != queryRunes[j] {
					match = false
					break
				}
			}
			if match {
				return true
			}
		}
	}

	// 2. 部分匹配 (bigram)
	if len(queryRunes) >= 2 {
		matchCount := 0
		for i := 0; i < len(queryRunes)-1; i++ {
			bigram := string(queryRunes[i : i+2])
			for j := 0; j < len(contentRunes)-1; j++ {
				if string(contentRunes[j:j+2]) == bigram {
					matchCount++
					break
				}
			}
		}
		if matchCount > len(queryRunes)/3 {
			return true
		}
	}

	return false
}
