package rag

import (
	"context"
	"math"
	"sync"
	"testing"
)

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name   string
		a      []float64
		b      []float64
		expect float64
		delta  float64
	}{
		{"identical vectors", []float64{1, 0, 0}, []float64{1, 0, 0}, 1.0, 0.001},
		{"orthogonal vectors", []float64{1, 0, 0}, []float64{0, 1, 0}, 0.0, 0.001},
		{"opposite vectors", []float64{1, 0, 0}, []float64{-1, 0, 0}, -1.0, 0.001},
		{"similar vectors", []float64{1, 1, 0}, []float64{1, 0, 0}, 0.707, 0.01},
		{"empty vectors", []float64{}, []float64{}, 0.0, 0.001},
		{"different lengths", []float64{1, 2}, []float64{1, 2, 3}, 0.0, 0.001},
		{"zero vector", []float64{0, 0, 0}, []float64{1, 1, 1}, 0.0, 0.001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cosineSimilarity(tt.a, tt.b)
			if math.Abs(result-tt.expect) > tt.delta {
				t.Errorf("expected %.4f ± %.4f, got %.4f", tt.expect, tt.delta, result)
			}
		})
	}
}

func TestMemoryVectorStore_Store(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryVectorStore()

	docs := []*VectorDocument{
		{ID: 1, Content: "Hello", Vector: []float64{1, 0, 0}},
		{ID: 2, Content: "World", Vector: []float64{0, 1, 0}},
		{ID: 3, Content: "Test", Vector: []float64{0, 0, 1}},
	}

	err := store.Store(ctx, docs)
	if err != nil {
		t.Fatalf("failed to store documents: %v", err)
	}

	if store.Count() != 3 {
		t.Errorf("expected 3 documents, got %d", store.Count())
	}
}

func TestMemoryVectorStore_StoreSkipsZeroID(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryVectorStore()

	docs := []*VectorDocument{
		{ID: 0, Content: "Skip me", Vector: []float64{1, 0, 0}},
		{ID: 1, Content: "Keep me", Vector: []float64{0, 1, 0}},
	}

	err := store.Store(ctx, docs)
	if err != nil {
		t.Fatalf("failed to store documents: %v", err)
	}

	if store.Count() != 1 {
		t.Errorf("expected 1 document, got %d", store.Count())
	}
}

func TestMemoryVectorStore_Delete(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryVectorStore()

	docs := []*VectorDocument{
		{ID: 1, Content: "Hello", Vector: []float64{1, 0, 0}},
		{ID: 2, Content: "World", Vector: []float64{0, 1, 0}},
	}

	store.Store(ctx, docs)

	err := store.Delete(ctx, []int64{1})
	if err != nil {
		t.Fatalf("failed to delete document: %v", err)
	}

	if store.Count() != 1 {
		t.Errorf("expected 1 document, got %d", store.Count())
	}
}

func TestMemoryVectorStore_Search(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryVectorStore()

	docs := []*VectorDocument{
		{ID: 1, Content: "Exact match", Vector: []float64{1, 0, 0}},
		{ID: 2, Content: "Orthogonal", Vector: []float64{0, 1, 0}},
		{ID: 3, Content: "Similar", Vector: []float64{0.9, 0.1, 0}},
	}

	store.Store(ctx, docs)

	// Search with query vector similar to ID 1
	results, err := store.Search(ctx, []float64{1, 0, 0}, 2, 0)
	if err != nil {
		t.Fatalf("failed to search: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	// First result should be the exact match
	if len(results) > 0 && results[0].ID != 1 {
		t.Errorf("expected first result to be ID 1, got %d", results[0].ID)
	}

	// Scores should be in descending order
	for i := 1; i < len(results); i++ {
		if results[i-1].Score < results[i].Score {
			t.Errorf("results should be sorted by score descending")
		}
	}
}

func TestMemoryVectorStore_SearchWithThreshold(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryVectorStore()

	docs := []*VectorDocument{
		{ID: 1, Content: "High match", Vector: []float64{1, 0, 0}},
		{ID: 2, Content: "Low match", Vector: []float64{0, 1, 0}},
	}

	store.Store(ctx, docs)

	// Search with high threshold
	results, err := store.Search(ctx, []float64{1, 0, 0}, 10, 0.9)
	if err != nil {
		t.Fatalf("failed to search: %v", err)
	}

	// Only the high match should pass threshold
	if len(results) != 1 {
		t.Errorf("expected 1 result with threshold 0.9, got %d", len(results))
	}
}

func TestMemoryVectorStore_SearchEmptyStore(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryVectorStore()

	results, err := store.Search(ctx, []float64{1, 0, 0}, 5, 0)
	if err != nil {
		t.Fatalf("failed to search empty store: %v", err)
	}

	if results != nil {
		t.Errorf("expected nil results for empty store, got %v", results)
	}
}

func TestMemoryVectorStore_Clear(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryVectorStore()

	docs := []*VectorDocument{
		{ID: 1, Content: "Hello", Vector: []float64{1, 0, 0}},
		{ID: 2, Content: "World", Vector: []float64{0, 1, 0}},
	}

	store.Store(ctx, docs)
	store.Clear(ctx)

	if store.Count() != 0 {
		t.Errorf("expected 0 documents after clear, got %d", store.Count())
	}
}

func TestMemoryVectorStore_Concurrent(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryVectorStore()

	var wg sync.WaitGroup
	const numGoroutines = 100

	// Concurrent stores
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			doc := &VectorDocument{
				ID:      id,
				Content: "test",
				Vector:  []float64{float64(id), 0, 0},
			}
			store.Store(ctx, []*VectorDocument{doc})
		}(int64(i + 1))
	}

	wg.Wait()

	if store.Count() != numGoroutines {
		t.Errorf("expected %d documents, got %d", numGoroutines, store.Count())
	}

	// Concurrent searches
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := store.Search(ctx, []float64{1, 0, 0}, 5, 0)
			if err != nil {
				t.Errorf("concurrent search failed: %v", err)
			}
		}()
	}

	wg.Wait()
}

func TestVectorStoreManager_GetStore(t *testing.T) {
	manager := GetVectorStoreManager()

	store1 := manager.GetStore(1)
	store2 := manager.GetStore(1)

	// Same collection should return same store
	if store1 != store2 {
		t.Error("expected same store for same collection ID")
	}

	// Different collection should return different store
	store3 := manager.GetStore(2)
	if store1 == store3 {
		t.Error("expected different store for different collection ID")
	}
}

func TestVectorStoreManager_DeleteStore(t *testing.T) {
	manager := GetVectorStoreManager()

	store1 := manager.GetStore(100)
	manager.DeleteStore(100)
	store2 := manager.GetStore(100)

	// After deletion, should get a new store
	if store1 == store2 {
		t.Error("expected new store after deletion")
	}
}

func TestCreateVectorStore(t *testing.T) {
	tests := []struct {
		name        string
		storeType   VectorStoreType
		expectError bool
	}{
		{"memory", VectorStoreTypeMemory, false},
		{"empty defaults to memory", "", false},
		{"redis not implemented", VectorStoreTypeRedis, true},
		{"milvus not implemented", VectorStoreTypeMilvus, true},
		{"elasticsearch not implemented", VectorStoreTypeElasticsearch, true},
		{"unknown type", VectorStoreType("unknown"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, err := CreateVectorStore(tt.storeType, nil)
			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			if !tt.expectError && store == nil {
				t.Error("expected non-nil store")
			}
		})
	}
}

func BenchmarkCosineSimilarity(b *testing.B) {
	a := make([]float64, 1536) // Common embedding dimension
	bb := make([]float64, 1536)
	for i := range a {
		a[i] = float64(i) / 1536.0
		bb[i] = float64(1536-i) / 1536.0
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cosineSimilarity(a, bb)
	}
}

func BenchmarkMemoryVectorStore_Search(b *testing.B) {
	ctx := context.Background()
	store := NewMemoryVectorStore()

	// Add 1000 documents
	docs := make([]*VectorDocument, 1000)
	for i := range docs {
		vector := make([]float64, 1536)
		for j := range vector {
			vector[j] = float64(i*1536+j) / 1536000.0
		}
		docs[i] = &VectorDocument{
			ID:      int64(i + 1),
			Content: "test",
			Vector:  vector,
		}
	}
	store.Store(ctx, docs)

	queryVector := make([]float64, 1536)
	for i := range queryVector {
		queryVector[i] = 0.5
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.Search(ctx, queryVector, 10, 0)
	}
}
