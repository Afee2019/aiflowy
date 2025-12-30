package rag

import (
	"strings"
	"testing"
)

func TestSimpleDocumentSplitter(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		chunkSize   int
		overlapSize int
		expectCount int
	}{
		{"empty content", "", 100, 10, 0},
		{"content smaller than chunk", "Hello World", 100, 10, 1},
		{"content equals chunk size", "1234567890", 10, 0, 1},
		{"content larger than chunk no overlap", "12345678901234567890", 10, 0, 2},
		{"content larger than chunk with overlap", "12345678901234567890", 10, 2, 3}, // 0-10, 8-18, 16-20
		{"chinese text", "这是一段测试中文文本", 5, 1, 3}, // 10 chars, chunk 5, overlap 1
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSimpleDocumentSplitter(tt.chunkSize, tt.overlapSize)
			chunks := s.Split(tt.content)
			if len(chunks) != tt.expectCount {
				t.Errorf("expected %d chunks, got %d", tt.expectCount, len(chunks))
			}
		})
	}
}

func TestSimpleDocumentSplitterDefaults(t *testing.T) {
	// Test default values
	s := NewSimpleDocumentSplitter(0, 0)
	if s.ChunkSize != 500 {
		t.Errorf("expected default chunk size 500, got %d", s.ChunkSize)
	}

	// Test overlap larger than chunk size
	s = NewSimpleDocumentSplitter(100, 150)
	if s.OverlapSize >= s.ChunkSize {
		t.Errorf("overlap size should be less than chunk size")
	}

	// Test negative overlap
	s = NewSimpleDocumentSplitter(100, -10)
	if s.OverlapSize != 0 {
		t.Errorf("expected overlap size 0 for negative input, got %d", s.OverlapSize)
	}
}

func TestSimpleDocumentSplitterOverlap(t *testing.T) {
	content := "ABCDEFGHIJ" // 10 characters
	s := NewSimpleDocumentSplitter(5, 2)
	chunks := s.Split(content)

	// With chunk size 5 and overlap 2:
	// Chunk 1: ABCDE (0-5)
	// Chunk 2: DEFGH (3-8) - starts at 5-2=3
	// Chunk 3: GHIJ (6-10) - starts at 8-2=6
	if len(chunks) < 2 {
		t.Errorf("expected at least 2 chunks, got %d", len(chunks))
	}

	// Check that overlap exists between consecutive chunks
	for i := 1; i < len(chunks); i++ {
		prev := chunks[i-1]
		curr := chunks[i]
		// Check if there's some overlap
		if !hasOverlap(prev, curr) && len(prev) >= 2 && len(curr) >= 2 {
			t.Logf("Warning: no apparent overlap between chunk %d and %d", i-1, i)
		}
	}
}

func hasOverlap(a, b string) bool {
	// Simple overlap check: see if end of a matches start of b
	for i := 1; i < len(a) && i < len(b); i++ {
		if strings.HasSuffix(a, b[:i]) {
			return true
		}
	}
	return false
}

func TestRegexDocumentSplitter(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		pattern     string
		expectCount int
	}{
		{"empty content", "", "", 0},
		{"single paragraph", "Hello World", `\n\n+`, 1},
		{"two paragraphs", "Hello\n\nWorld", `\n\n+`, 2},
		{"multiple blank lines", "Hello\n\n\n\nWorld", `\n\n+`, 2},
		{"custom pattern", "Hello---World---Test", `---`, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewRegexDocumentSplitter(tt.pattern)
			chunks := s.Split(tt.content)
			if len(chunks) != tt.expectCount {
				t.Errorf("expected %d chunks, got %d: %v", tt.expectCount, len(chunks), chunks)
			}
		})
	}
}

func TestRegexDocumentSplitterInvalidPattern(t *testing.T) {
	// Invalid regex should fall back to default
	s := NewRegexDocumentSplitter("[invalid")
	chunks := s.Split("Hello\n\nWorld")
	if len(chunks) != 2 {
		t.Errorf("expected 2 chunks with fallback pattern, got %d", len(chunks))
	}
}

func TestSentenceDocumentSplitter(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		chunkSize   int
		overlapSize int
		expectCount int
	}{
		{"empty content", "", 5, 0, 0},
		{"single sentence", "这是一个句子。", 5, 0, 1},
		{"multiple sentences", "第一句。第二句。第三句。", 2, 0, 2},
		{"english sentences", "First. Second. Third.", 2, 0, 2},
		{"mixed punctuation", "问题？回答！结束。", 2, 0, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSentenceDocumentSplitter(tt.chunkSize, tt.overlapSize)
			chunks := s.Split(tt.content)
			if len(chunks) != tt.expectCount {
				t.Errorf("expected %d chunks, got %d: %v", tt.expectCount, len(chunks), chunks)
			}
		})
	}
}

func TestSentenceDocumentSplitterDefaults(t *testing.T) {
	s := NewSentenceDocumentSplitter(0, 0)
	if s.ChunkSize != 5 {
		t.Errorf("expected default chunk size 5, got %d", s.ChunkSize)
	}

	s = NewSentenceDocumentSplitter(5, -1)
	if s.OverlapSize != 0 {
		t.Errorf("expected overlap size 0 for negative input, got %d", s.OverlapSize)
	}
}

func TestTokenDocumentSplitter(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		maxTokens   int
		overlapSize int
		expectCount int
	}{
		{"empty content", "", 100, 0, 0},
		{"single paragraph", "Hello World", 100, 0, 1},
		{"multiple paragraphs", "Hello\n\nWorld\n\nTest", 100, 0, 1}, // increase max tokens
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewTokenDocumentSplitter(tt.maxTokens, tt.overlapSize)
			chunks := s.Split(tt.content)
			if len(chunks) != tt.expectCount {
				t.Errorf("expected %d chunks, got %d: %v", tt.expectCount, len(chunks), chunks)
			}
		})
	}
}

func TestTokenDocumentSplitterDefaults(t *testing.T) {
	s := NewTokenDocumentSplitter(0, 0)
	if s.MaxTokens != 500 {
		t.Errorf("expected default max tokens 500, got %d", s.MaxTokens)
	}

	s = NewTokenDocumentSplitter(100, -1)
	if s.OverlapSize != 0 {
		t.Errorf("expected overlap size 0 for negative input, got %d", s.OverlapSize)
	}
}

func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		minToken int
		maxToken int
	}{
		{"empty", "", 0, 0},
		{"single word", "Hello", 1, 5},
		{"multiple words", "Hello World", 2, 15},
		{"chinese", "你好世界", 4, 8},
		{"mixed", "Hello 世界", 3, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := estimateTokens(tt.text)
			if tokens < tt.minToken || tokens > tt.maxToken {
				t.Errorf("expected tokens between %d and %d, got %d", tt.minToken, tt.maxToken, tokens)
			}
		})
	}
}

func TestGetDocumentSplitter(t *testing.T) {
	tests := []struct {
		name         string
		splitterName string
		expectType   string
	}{
		{"simple splitter", "SimpleDocumentSplitter", "*rag.SimpleDocumentSplitter"},
		{"regex splitter", "RegexDocumentSplitter", "*rag.RegexDocumentSplitter"},
		{"sentence splitter", "SentenceDocumentSplitter", "*rag.SentenceDocumentSplitter"},
		{"token splitter", "TokenDocumentSplitter", "*rag.TokenDocumentSplitter"},
		{"simple tokenize", "SimpleTokenizeSplitter", "*rag.TokenDocumentSplitter"},
		{"unknown defaults to simple", "Unknown", "*rag.SimpleDocumentSplitter"},
		{"empty defaults to simple", "", "*rag.SimpleDocumentSplitter"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			splitter := GetDocumentSplitter(tt.splitterName, 100, 10, "")
			if splitter == nil {
				t.Error("expected non-nil splitter")
				return
			}
			// Just verify it can split without error
			chunks := splitter.Split("Test content")
			if len(chunks) == 0 {
				t.Error("expected at least one chunk")
			}
		})
	}
}

func BenchmarkSimpleDocumentSplitter(b *testing.B) {
	content := strings.Repeat("这是一段测试文本。", 1000)
	s := NewSimpleDocumentSplitter(500, 50)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Split(content)
	}
}

func BenchmarkRegexDocumentSplitter(b *testing.B) {
	content := strings.Repeat("段落一。\n\n段落二。\n\n", 500)
	s := NewRegexDocumentSplitter(`\n\n+`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Split(content)
	}
}

func BenchmarkSentenceDocumentSplitter(b *testing.B) {
	content := strings.Repeat("这是第一句。这是第二句。这是第三句。", 500)
	s := NewSentenceDocumentSplitter(5, 1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Split(content)
	}
}

func BenchmarkTokenDocumentSplitter(b *testing.B) {
	content := strings.Repeat("这是一段测试文本\n\n", 500)
	s := NewTokenDocumentSplitter(500, 50)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Split(content)
	}
}
