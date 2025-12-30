package rag

import (
	"context"
	"fmt"
	"sync"

	"github.com/cloudwego/eino/components/embedding"

	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/internal/repository"
)

// Retriever RAG 检索器
type Retriever struct {
	collectionID  int64
	embedder      embedding.Embedder
	vectorStore   VectorStore
	embeddingRepo *repository.DocumentRepository
}

// RetrieverConfig 检索器配置
type RetrieverConfig struct {
	TopK           int     // 返回的最大结果数
	ScoreThreshold float64 // 相似度阈值
}

// DefaultRetrieverConfig 默认检索器配置
func DefaultRetrieverConfig() *RetrieverConfig {
	return &RetrieverConfig{
		TopK:           5,
		ScoreThreshold: 0.5,
	}
}

// NewRetriever 创建检索器
func NewRetriever(collectionID int64, embedder embedding.Embedder, vectorStore VectorStore) *Retriever {
	return &Retriever{
		collectionID:  collectionID,
		embedder:      embedder,
		vectorStore:   vectorStore,
		embeddingRepo: repository.NewDocumentRepository(),
	}
}

// Retrieve 检索相关文档
func (r *Retriever) Retrieve(ctx context.Context, query string, config *RetrieverConfig) ([]*VectorDocument, error) {
	if config == nil {
		config = DefaultRetrieverConfig()
	}

	// 向量化查询文本
	queryVector, err := r.embedder.EmbedStrings(ctx, []string{query})
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	if len(queryVector) == 0 || len(queryVector[0]) == 0 {
		return nil, fmt.Errorf("empty query embedding")
	}

	// 向量检索
	docs, err := r.vectorStore.Search(ctx, queryVector[0], config.TopK, config.ScoreThreshold)
	if err != nil {
		return nil, fmt.Errorf("failed to search vector store: %w", err)
	}

	return docs, nil
}

// RAGService RAG 服务
type RAGService struct {
	embeddingService *EmbeddingService
	modelRepo        *repository.ModelRepository
	docRepo          *repository.DocumentRepository
	collectionRepo   *repository.DocumentCollectionRepository
	mu               sync.RWMutex
	retrievers       map[int64]*Retriever
}

var (
	ragService     *RAGService
	ragServiceOnce sync.Once
)

// GetRAGService 获取 RAG 服务单例
func GetRAGService() *RAGService {
	ragServiceOnce.Do(func() {
		ragService = &RAGService{
			embeddingService: NewEmbeddingService(),
			modelRepo:        repository.NewModelRepository(repository.GetDB()),
			docRepo:          repository.NewDocumentRepository(),
			collectionRepo:   repository.NewDocumentCollectionRepository(),
			retrievers:       make(map[int64]*Retriever),
		}
	})
	return ragService
}

// IndexDocumentChunks 索引文档分块 (向量化并存储)
func (s *RAGService) IndexDocumentChunks(ctx context.Context, collection *entity.DocumentCollection, chunks []*entity.DocumentChunk) error {
	if collection == nil || !collection.VectorStoreEnable {
		return nil
	}

	if collection.VectorEmbedModelID == nil || *collection.VectorEmbedModelID == 0 {
		return fmt.Errorf("embedding model not configured for collection %d", collection.ID)
	}

	// 获取 embedding 模型
	model, err := s.modelRepo.GetModelByID(ctx, *collection.VectorEmbedModelID)
	if err != nil {
		return fmt.Errorf("failed to get embedding model: %w", err)
	}
	if model == nil {
		return fmt.Errorf("embedding model %d not found", *collection.VectorEmbedModelID)
	}

	// 加载模型提供商
	model.ModelProvider, _ = s.modelRepo.GetProviderByID(ctx, model.ProviderID)

	// 创建 embedder
	embedder, err := s.embeddingService.CreateEmbedder(ctx, model)
	if err != nil {
		return fmt.Errorf("failed to create embedder: %w", err)
	}

	// 准备文本
	var texts []string
	var chunkIDs []int64
	for _, chunk := range chunks {
		if chunk.Content == "" {
			continue
		}
		texts = append(texts, chunk.Content)
		chunkIDs = append(chunkIDs, chunk.ID)
	}

	if len(texts) == 0 {
		return nil
	}

	// 批量向量化
	vectors, err := embedder.EmbedStrings(ctx, texts)
	if err != nil {
		return fmt.Errorf("failed to embed texts: %w", err)
	}

	// 构造向量文档
	var vectorDocs []*VectorDocument
	for i, vector := range vectors {
		if i >= len(chunkIDs) {
			break
		}
		vectorDocs = append(vectorDocs, &VectorDocument{
			ID:      chunkIDs[i],
			Content: texts[i],
			Vector:  vector,
			Metadata: map[string]interface{}{
				"collection_id": collection.ID,
			},
		})
	}

	// 存储到向量库
	store := GetVectorStoreManager().GetStore(collection.ID)
	return store.Store(ctx, vectorDocs)
}

// DeleteDocumentChunks 删除文档分块的向量
func (s *RAGService) DeleteDocumentChunks(ctx context.Context, collectionID int64, chunkIDs []int64) error {
	store := GetVectorStoreManager().GetStore(collectionID)
	return store.Delete(ctx, chunkIDs)
}

// Search 搜索相关文档
func (s *RAGService) Search(ctx context.Context, collectionID int64, query string, topK int) ([]*VectorDocument, error) {
	// 获取知识库配置
	collection, err := s.collectionRepo.GetByID(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}
	if collection == nil {
		return nil, fmt.Errorf("collection %d not found", collectionID)
	}

	if !collection.VectorStoreEnable {
		// 向量存储未启用，回退到全文搜索
		return s.fullTextSearch(ctx, collectionID, query, topK)
	}

	if collection.VectorEmbedModelID == nil || *collection.VectorEmbedModelID == 0 {
		return s.fullTextSearch(ctx, collectionID, query, topK)
	}

	// 获取 retriever
	retriever, err := s.getOrCreateRetriever(ctx, collection)
	if err != nil {
		return nil, err
	}

	config := &RetrieverConfig{
		TopK:           topK,
		ScoreThreshold: 0.3, // 较低的阈值以获取更多结果
	}

	return retriever.Retrieve(ctx, query, config)
}

// getOrCreateRetriever 获取或创建检索器
func (s *RAGService) getOrCreateRetriever(ctx context.Context, collection *entity.DocumentCollection) (*Retriever, error) {
	s.mu.RLock()
	retriever, ok := s.retrievers[collection.ID]
	s.mu.RUnlock()

	if ok && retriever != nil {
		return retriever, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Double check
	if retriever, ok = s.retrievers[collection.ID]; ok && retriever != nil {
		return retriever, nil
	}

	// 获取 embedding 模型
	model, err := s.modelRepo.GetModelByID(ctx, *collection.VectorEmbedModelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get embedding model: %w", err)
	}
	if model == nil {
		return nil, fmt.Errorf("embedding model %d not found", *collection.VectorEmbedModelID)
	}

	// 加载模型提供商
	model.ModelProvider, _ = s.modelRepo.GetProviderByID(ctx, model.ProviderID)

	// 创建 embedder
	embedder, err := s.embeddingService.CreateEmbedder(ctx, model)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedder: %w", err)
	}

	// 获取向量存储
	store := GetVectorStoreManager().GetStore(collection.ID)

	// 创建检索器
	retriever = NewRetriever(collection.ID, embedder, store)
	s.retrievers[collection.ID] = retriever

	return retriever, nil
}

// fullTextSearch 全文搜索 (简单的关键词匹配)
func (s *RAGService) fullTextSearch(ctx context.Context, collectionID int64, query string, topK int) ([]*VectorDocument, error) {
	// 从数据库获取所有分块
	chunks, err := s.docRepo.ListChunksByCollectionID(ctx, collectionID)
	if err != nil {
		return nil, err
	}

	// 简单的关键词匹配 (实际应用中可以使用更复杂的算法)
	var results []*VectorDocument
	for _, chunk := range chunks {
		// 简单匹配：检查内容是否包含查询关键词
		if containsKeywords(chunk.Content, query) {
			results = append(results, &VectorDocument{
				ID:      chunk.ID,
				Content: chunk.Content,
				Score:   0.5, // 默认分数
			})
		}
	}

	if len(results) > topK {
		results = results[:topK]
	}

	return results, nil
}

// containsKeywords 简单关键词匹配
func containsKeywords(content, query string) bool {
	// 简单实现：检查内容是否包含查询词
	return len(query) > 0 && len(content) > 0 &&
		(contains(content, query) || partialMatch(content, query))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func partialMatch(content, query string) bool {
	// 分词匹配 (简化版)
	queryRunes := []rune(query)
	if len(queryRunes) < 2 {
		return false
	}

	contentRunes := []rune(content)
	matchCount := 0
	for i := 0; i < len(queryRunes)-1; i++ {
		bigram := string(queryRunes[i : i+2])
		if findSubstring(string(contentRunes), bigram) {
			matchCount++
		}
	}

	return matchCount > len(queryRunes)/3
}

// InvalidateRetriever 使检索器失效 (当知识库配置变更时调用)
func (s *RAGService) InvalidateRetriever(collectionID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.retrievers, collectionID)
}

// ClearCollectionIndex 清空知识库索引
func (s *RAGService) ClearCollectionIndex(collectionID int64) {
	GetVectorStoreManager().DeleteStore(collectionID)
	s.InvalidateRetriever(collectionID)
}
