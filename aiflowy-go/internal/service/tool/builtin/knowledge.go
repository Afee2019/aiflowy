package builtin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/schema"

	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/internal/service/rag"
)

// KnowledgeTool 知识库检索工具
type KnowledgeTool struct {
	CollectionID   int64
	ToolName       string
	ToolDesc       string
	EnglishName    string
	useEnglishName bool
}

// NewKnowledgeTool 创建知识库工具
func NewKnowledgeTool(collection *entity.DocumentCollection, useEnglishName bool) *KnowledgeTool {
	name := collection.Title
	if useEnglishName && collection.EnglishName != "" {
		name = collection.EnglishName
	}

	desc := collection.Description
	if desc == "" {
		desc = fmt.Sprintf("搜索 %s 知识库中的相关信息", collection.Title)
	}

	return &KnowledgeTool{
		CollectionID:   collection.ID,
		ToolName:       name,
		ToolDesc:       desc,
		EnglishName:    collection.EnglishName,
		useEnglishName: useEnglishName,
	}
}

// Name 获取工具名称
func (t *KnowledgeTool) Name() string {
	return t.ToolName
}

// Description 获取工具描述
func (t *KnowledgeTool) Description() string {
	return t.ToolDesc
}

// Parameters 获取工具参数
func (t *KnowledgeTool) Parameters() map[string]*schema.ParameterInfo {
	return map[string]*schema.ParameterInfo{
		"input": {
			Type: schema.String,
			Desc: "要在知识库中搜索的关键词或问题",
		},
	}
}

// Execute 执行工具
func (t *KnowledgeTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// 获取查询参数
	input, ok := args["input"].(string)
	if !ok || input == "" {
		return nil, fmt.Errorf("input parameter is required")
	}

	// 调用 RAG 服务进行检索
	ragService := rag.GetRAGService()
	docs, err := ragService.Search(ctx, t.CollectionID, input, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to search knowledge base: %w", err)
	}

	if len(docs) == 0 {
		return "未找到相关信息", nil
	}

	// 构造返回结果
	var results []string
	for i, doc := range docs {
		result := fmt.Sprintf("[%d] %s", i+1, strings.TrimSpace(doc.Content))
		if doc.Score > 0 {
			result = fmt.Sprintf("[%d] (相关度: %.2f) %s", i+1, doc.Score, strings.TrimSpace(doc.Content))
		}
		results = append(results, result)
	}

	return strings.Join(results, "\n\n"), nil
}

// KnowledgeToolResult 知识库工具结果
type KnowledgeToolResult struct {
	Query   string              `json:"query"`
	Results []KnowledgeDocument `json:"results"`
}

// KnowledgeDocument 知识库文档
type KnowledgeDocument struct {
	ID      int64   `json:"id,string"`
	Content string  `json:"content"`
	Score   float64 `json:"score,omitempty"`
}

// ExecuteAndGetStructured 执行工具并返回结构化结果
func (t *KnowledgeTool) ExecuteAndGetStructured(ctx context.Context, query string) (*KnowledgeToolResult, error) {
	ragService := rag.GetRAGService()
	docs, err := ragService.Search(ctx, t.CollectionID, query, 5)
	if err != nil {
		return nil, err
	}

	result := &KnowledgeToolResult{
		Query:   query,
		Results: make([]KnowledgeDocument, 0, len(docs)),
	}

	for _, doc := range docs {
		result.Results = append(result.Results, KnowledgeDocument{
			ID:      doc.ID,
			Content: doc.Content,
			Score:   doc.Score,
		})
	}

	return result, nil
}

// ToJSON 转换为 JSON 字符串
func (r *KnowledgeToolResult) ToJSON() string {
	data, _ := json.Marshal(r)
	return string(data)
}

func init() {
	// 注册知识库工具类型
	// 注意：具体的知识库工具是动态创建的，不在这里注册
}
