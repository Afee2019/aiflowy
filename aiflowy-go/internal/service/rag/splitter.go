package rag

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// DocumentSplitter 文档分块器接口
type DocumentSplitter interface {
	Split(content string) []string
}

// SimpleDocumentSplitter 简单分块器 (按字符数分块，支持重叠)
type SimpleDocumentSplitter struct {
	ChunkSize   int // 分块大小
	OverlapSize int // 重叠大小
}

// NewSimpleDocumentSplitter 创建简单分块器
func NewSimpleDocumentSplitter(chunkSize, overlapSize int) *SimpleDocumentSplitter {
	if chunkSize <= 0 {
		chunkSize = 500
	}
	if overlapSize < 0 {
		overlapSize = 0
	}
	if overlapSize >= chunkSize {
		overlapSize = chunkSize / 4
	}
	return &SimpleDocumentSplitter{
		ChunkSize:   chunkSize,
		OverlapSize: overlapSize,
	}
}

// Split 分块
func (s *SimpleDocumentSplitter) Split(content string) []string {
	if content == "" {
		return nil
	}

	runes := []rune(content)
	length := len(runes)

	if length <= s.ChunkSize {
		return []string{content}
	}

	var chunks []string
	start := 0

	for start < length {
		end := start + s.ChunkSize
		if end > length {
			end = length
		}

		chunk := string(runes[start:end])
		chunks = append(chunks, strings.TrimSpace(chunk))

		// 如果已经处理到末尾，退出循环
		if end >= length {
			break
		}

		// 下一个块的起始位置 (考虑重叠)
		nextStart := end - s.OverlapSize
		if nextStart <= start || s.OverlapSize == 0 {
			nextStart = end
		}
		start = nextStart
	}

	return chunks
}

// RegexDocumentSplitter 正则分块器
type RegexDocumentSplitter struct {
	Pattern string
	regex   *regexp.Regexp
}

// NewRegexDocumentSplitter 创建正则分块器
func NewRegexDocumentSplitter(pattern string) *RegexDocumentSplitter {
	if pattern == "" {
		pattern = `\n\n+` // 默认按段落分割
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		re = regexp.MustCompile(`\n\n+`)
	}
	return &RegexDocumentSplitter{
		Pattern: pattern,
		regex:   re,
	}
}

// Split 分块
func (s *RegexDocumentSplitter) Split(content string) []string {
	if content == "" {
		return nil
	}

	parts := s.regex.Split(content, -1)
	var chunks []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			chunks = append(chunks, trimmed)
		}
	}
	return chunks
}

// SentenceDocumentSplitter 按句子分块器
type SentenceDocumentSplitter struct {
	ChunkSize   int // 每个块的最大句子数
	OverlapSize int // 重叠句子数
}

// NewSentenceDocumentSplitter 创建句子分块器
func NewSentenceDocumentSplitter(chunkSize, overlapSize int) *SentenceDocumentSplitter {
	if chunkSize <= 0 {
		chunkSize = 5
	}
	if overlapSize < 0 {
		overlapSize = 0
	}
	return &SentenceDocumentSplitter{
		ChunkSize:   chunkSize,
		OverlapSize: overlapSize,
	}
}

// Split 分块
func (s *SentenceDocumentSplitter) Split(content string) []string {
	if content == "" {
		return nil
	}

	// 按句子分割 (中英文句号、问号、感叹号)
	sentencePattern := regexp.MustCompile(`[。！？.!?]+`)
	sentences := sentencePattern.Split(content, -1)

	// 过滤空句子
	var validSentences []string
	for _, sent := range sentences {
		trimmed := strings.TrimSpace(sent)
		if trimmed != "" {
			validSentences = append(validSentences, trimmed)
		}
	}

	if len(validSentences) == 0 {
		return nil
	}

	if len(validSentences) <= s.ChunkSize {
		return []string{strings.Join(validSentences, "。")}
	}

	var chunks []string
	start := 0

	for start < len(validSentences) {
		end := start + s.ChunkSize
		if end > len(validSentences) {
			end = len(validSentences)
		}

		chunk := strings.Join(validSentences[start:end], "。")
		chunks = append(chunks, chunk)

		// 如果已经处理到末尾，退出循环
		if end >= len(validSentences) {
			break
		}

		nextStart := end - s.OverlapSize
		if nextStart <= start || s.OverlapSize == 0 {
			nextStart = end
		}
		start = nextStart
	}

	return chunks
}

// TokenDocumentSplitter 按 Token 分块器 (简化版，按字符估算)
type TokenDocumentSplitter struct {
	MaxTokens   int // 每块最大 token 数
	OverlapSize int // 重叠 token 数
}

// NewTokenDocumentSplitter 创建 Token 分块器
func NewTokenDocumentSplitter(maxTokens, overlapSize int) *TokenDocumentSplitter {
	if maxTokens <= 0 {
		maxTokens = 500
	}
	if overlapSize < 0 {
		overlapSize = 0
	}
	return &TokenDocumentSplitter{
		MaxTokens:   maxTokens,
		OverlapSize: overlapSize,
	}
}

// estimateTokens 估算 token 数 (简化：中文按字符，英文按空格分词)
func estimateTokens(text string) int {
	// 简化估算：每个中文字符约 1 token，每个英文单词约 1 token
	tokens := 0
	words := strings.Fields(text)
	for _, word := range words {
		runeCount := utf8.RuneCountInString(word)
		if runeCount > 1 {
			// 可能包含中文
			tokens += runeCount
		} else {
			tokens++
		}
	}
	return tokens
}

// Split 分块
func (s *TokenDocumentSplitter) Split(content string) []string {
	if content == "" {
		return nil
	}

	// 按段落和句子分割，然后合并到 token 限制内
	paragraphs := regexp.MustCompile(`\n+`).Split(content, -1)

	var chunks []string
	var currentChunk strings.Builder
	currentTokens := 0

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		paraTokens := estimateTokens(para)

		if currentTokens+paraTokens > s.MaxTokens && currentTokens > 0 {
			// 当前块已满，保存并开始新块
			chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
			currentChunk.Reset()
			currentTokens = 0
		}

		if currentChunk.Len() > 0 {
			currentChunk.WriteString("\n")
		}
		currentChunk.WriteString(para)
		currentTokens += paraTokens
	}

	// 保存最后一块
	if currentChunk.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
	}

	return chunks
}

// GetDocumentSplitter 根据名称获取分块器
func GetDocumentSplitter(splitterName string, chunkSize, overlapSize int, regex string) DocumentSplitter {
	switch splitterName {
	case "SimpleDocumentSplitter":
		return NewSimpleDocumentSplitter(chunkSize, overlapSize)
	case "RegexDocumentSplitter":
		return NewRegexDocumentSplitter(regex)
	case "SentenceDocumentSplitter":
		return NewSentenceDocumentSplitter(chunkSize, overlapSize)
	case "SimpleTokenizeSplitter", "TokenDocumentSplitter":
		return NewTokenDocumentSplitter(chunkSize, overlapSize)
	default:
		// 默认使用简单分块器
		return NewSimpleDocumentSplitter(chunkSize, overlapSize)
	}
}
