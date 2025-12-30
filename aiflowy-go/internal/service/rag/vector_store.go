package rag

import (
	"context"
	"fmt"
	"math"
	"sort"
	"sync"
)

// VectorDocument 向量化文档
type VectorDocument struct {
	ID        int64     // 文档分块 ID
	Content   string    // 原始文本内容
	Vector    []float64 // 向量表示
	Metadata  map[string]interface{}
	Score     float64 // 相似度分数 (检索时填充)
}

// VectorStore 向量存储接口
type VectorStore interface {
	// Store 存储向量化文档
	Store(ctx context.Context, docs []*VectorDocument) error
	// Delete 删除指定 ID 的文档
	Delete(ctx context.Context, ids []int64) error
	// Search 向量相似度检索
	Search(ctx context.Context, queryVector []float64, topK int, threshold float64) ([]*VectorDocument, error)
	// Clear 清空所有文档
	Clear(ctx context.Context) error
	// Count 获取文档数量
	Count() int
}

// MemoryVectorStore 内存向量存储
type MemoryVectorStore struct {
	mu   sync.RWMutex
	docs map[int64]*VectorDocument
}

// NewMemoryVectorStore 创建内存向量存储
func NewMemoryVectorStore() *MemoryVectorStore {
	return &MemoryVectorStore{
		docs: make(map[int64]*VectorDocument),
	}
}

// Store 存储向量化文档
func (s *MemoryVectorStore) Store(ctx context.Context, docs []*VectorDocument) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, doc := range docs {
		if doc.ID == 0 {
			continue
		}
		s.docs[doc.ID] = doc
	}
	return nil
}

// Delete 删除指定 ID 的文档
func (s *MemoryVectorStore) Delete(ctx context.Context, ids []int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, id := range ids {
		delete(s.docs, id)
	}
	return nil
}

// Search 向量相似度检索
func (s *MemoryVectorStore) Search(ctx context.Context, queryVector []float64, topK int, threshold float64) ([]*VectorDocument, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.docs) == 0 {
		return nil, nil
	}

	if topK <= 0 {
		topK = 5
	}

	// 计算所有文档的相似度
	type scored struct {
		doc   *VectorDocument
		score float64
	}
	var results []scored

	for _, doc := range s.docs {
		if len(doc.Vector) == 0 {
			continue
		}
		score := cosineSimilarity(queryVector, doc.Vector)
		if threshold > 0 && score < threshold {
			continue
		}
		results = append(results, scored{doc: doc, score: score})
	}

	// 按相似度降序排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	// 取 TopK
	if len(results) > topK {
		results = results[:topK]
	}

	// 构造返回结果
	var docs []*VectorDocument
	for _, r := range results {
		doc := &VectorDocument{
			ID:       r.doc.ID,
			Content:  r.doc.Content,
			Vector:   r.doc.Vector,
			Metadata: r.doc.Metadata,
			Score:    r.score,
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

// Clear 清空所有文档
func (s *MemoryVectorStore) Clear(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.docs = make(map[int64]*VectorDocument)
	return nil
}

// Count 获取文档数量
func (s *MemoryVectorStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.docs)
}

// cosineSimilarity 计算余弦相似度
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// VectorStoreManager 向量存储管理器 (按知识库 ID 管理)
type VectorStoreManager struct {
	mu     sync.RWMutex
	stores map[int64]VectorStore
}

var (
	vectorStoreManager     *VectorStoreManager
	vectorStoreManagerOnce sync.Once
)

// GetVectorStoreManager 获取向量存储管理器单例
func GetVectorStoreManager() *VectorStoreManager {
	vectorStoreManagerOnce.Do(func() {
		vectorStoreManager = &VectorStoreManager{
			stores: make(map[int64]VectorStore),
		}
	})
	return vectorStoreManager
}

// GetStore 获取指定知识库的向量存储
func (m *VectorStoreManager) GetStore(collectionID int64) VectorStore {
	m.mu.RLock()
	store, ok := m.stores[collectionID]
	m.mu.RUnlock()

	if ok {
		return store
	}

	// 创建新的内存存储
	m.mu.Lock()
	defer m.mu.Unlock()

	// Double check
	if store, ok = m.stores[collectionID]; ok {
		return store
	}

	store = NewMemoryVectorStore()
	m.stores[collectionID] = store
	return store
}

// DeleteStore 删除指定知识库的向量存储
func (m *VectorStoreManager) DeleteStore(collectionID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.stores, collectionID)
}

// VectorStoreType 向量存储类型
type VectorStoreType string

const (
	VectorStoreTypeMemory        VectorStoreType = "memory"
	VectorStoreTypeRedis         VectorStoreType = "redis"
	VectorStoreTypeMilvus        VectorStoreType = "milvus"
	VectorStoreTypeElasticsearch VectorStoreType = "elasticsearch"
)

// CreateVectorStore 根据类型创建向量存储
func CreateVectorStore(storeType VectorStoreType, config map[string]interface{}) (VectorStore, error) {
	switch storeType {
	case VectorStoreTypeMemory, "":
		return NewMemoryVectorStore(), nil
	case VectorStoreTypeRedis:
		// TODO: 实现 Redis 向量存储
		return nil, fmt.Errorf("redis vector store not implemented yet")
	case VectorStoreTypeMilvus:
		// TODO: 实现 Milvus 向量存储
		return nil, fmt.Errorf("milvus vector store not implemented yet")
	case VectorStoreTypeElasticsearch:
		// TODO: 实现 Elasticsearch 向量存储
		return nil, fmt.Errorf("elasticsearch vector store not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported vector store type: %s", storeType)
	}
}
