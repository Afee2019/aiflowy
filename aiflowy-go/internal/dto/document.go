package dto

// DocumentCollectionSaveRequest 知识库保存请求
type DocumentCollectionSaveRequest struct {
	ID                    string  `json:"id,omitempty"`
	Alias                 string  `json:"alias,omitempty"`
	Icon                  string  `json:"icon,omitempty"`
	Title                 string  `json:"title"`
	Description           string  `json:"description,omitempty"`
	Slug                  string  `json:"slug,omitempty"`
	VectorStoreEnable     *bool   `json:"vectorStoreEnable,omitempty"`
	VectorStoreType       string  `json:"vectorStoreType,omitempty"`
	VectorStoreCollection string  `json:"vectorStoreCollection,omitempty"`
	VectorStoreConfig     string  `json:"vectorStoreConfig,omitempty"`
	VectorEmbedModelID    string  `json:"vectorEmbedModelId,omitempty"`
	RerankModelID         string  `json:"rerankModelId,omitempty"`
	SearchEngineEnable    *bool   `json:"searchEngineEnable,omitempty"`
	EnglishName           string  `json:"englishName,omitempty"`
	Options               string  `json:"options,omitempty"`
}

// DocumentSaveRequest 文档保存请求
type DocumentSaveRequest struct {
	ID           string `json:"id,omitempty"`
	CollectionID string `json:"collectionId"`
	DocumentType string `json:"documentType,omitempty"`
	DocumentPath string `json:"documentPath,omitempty"`
	Title        string `json:"title"`
	Content      string `json:"content,omitempty"`
	ContentType  string `json:"contentType,omitempty"`
	Slug         string `json:"slug,omitempty"`
	OrderNo      *int   `json:"orderNo,omitempty"`
	Options      string `json:"options,omitempty"`
}

// DocumentListRequest 文档列表请求
type DocumentListRequest struct {
	ID         string `json:"id" query:"id"`           // 知识库 ID
	Title      string `json:"title" query:"title"`     // 文档标题 (模糊搜索)
	PageNumber int    `json:"pageNumber" query:"pageNumber"`
	PageSize   int    `json:"pageSize" query:"pageSize"`
}

// DocumentListResponse 文档列表响应
type DocumentListResponse struct {
	Total    int64       `json:"total"`
	PageNo   int         `json:"pageNo"`
	PageSize int         `json:"pageSize"`
	List     interface{} `json:"list"`
}

// BotDocumentCollectionUpdateRequest Bot-知识库关联更新请求
type BotDocumentCollectionUpdateRequest struct {
	BotID        string   `json:"botId"`
	KnowledgeIDs []string `json:"knowledgeIds"`
}

// TextSplitRequest 文本拆分请求
type TextSplitRequest struct {
	Operation      string `json:"operation" form:"operation"`                // textSplit / saveText
	FilePath       string `json:"filePath" form:"filePath"`
	FileOriginName string `json:"fileOriginName" form:"fileOriginName"`
	KnowledgeID    string `json:"knowledgeId" form:"knowledgeId"`
	SplitterName   string `json:"splitterName" form:"splitterName"`
	ChunkSize      int    `json:"chunkSize" form:"chunkSize"`
	OverlapSize    int    `json:"overlapSize" form:"overlapSize"`
	Regex          string `json:"regex" form:"regex"`
	RowsPerChunk   int    `json:"rowsPerChunk" form:"rowsPerChunk"`
	PageNumber     int    `json:"pageNumber" form:"pageNumber"`
	PageSize       int    `json:"pageSize" form:"pageSize"`
}

// TextSplitResponse 文本拆分响应
type TextSplitResponse struct {
	Total  int64    `json:"total"`
	Chunks []string `json:"chunks"`
}

// UploadResponse 上传响应
type UploadResponse struct {
	Path string `json:"path"`
}
