package entity

import "time"

// DocumentCollection 知识库实体
type DocumentCollection struct {
	ID                     int64      `db:"id" json:"id,string"`
	Alias                  string     `db:"alias" json:"alias,omitempty"`
	DeptID                 int64      `db:"dept_id" json:"deptId,string,omitempty"`
	TenantID               int64      `db:"tenant_id" json:"tenantId,string,omitempty"`
	Icon                   string     `db:"icon" json:"icon,omitempty"`
	Title                  string     `db:"title" json:"title,omitempty"`
	Description            string     `db:"description" json:"description,omitempty"`
	Slug                   string     `db:"slug" json:"slug,omitempty"`
	VectorStoreEnable      bool       `db:"vector_store_enable" json:"vectorStoreEnable"`
	VectorStoreType        string     `db:"vector_store_type" json:"vectorStoreType,omitempty"`
	VectorStoreCollection  string     `db:"vector_store_collection" json:"vectorStoreCollection,omitempty"`
	VectorStoreConfig      string     `db:"vector_store_config" json:"vectorStoreConfig,omitempty"`
	VectorEmbedModelID     *int64     `db:"vector_embed_model_id" json:"vectorEmbedModelId,string,omitempty"`
	RerankModelID          *int64     `db:"rerank_model_id" json:"rerankModelId,string,omitempty"`
	SearchEngineEnable     bool       `db:"search_engine_enable" json:"searchEngineEnable"`
	EnglishName            string     `db:"english_name" json:"englishName,omitempty"`
	Options                string     `db:"options" json:"options,omitempty"`
	Created                *time.Time `db:"created" json:"created,omitempty"`
	CreatedBy              *int64     `db:"created_by" json:"createdBy,string,omitempty"`
	Modified               *time.Time `db:"modified" json:"modified,omitempty"`
	ModifiedBy             *int64     `db:"modified_by" json:"modifiedBy,string,omitempty"`

	// 非数据库字段
	DocumentCount          int        `db:"-" json:"documentCount,omitempty"`
	EmbedModel             *Model     `db:"-" json:"embedModel,omitempty"`
	RerankModel            *Model     `db:"-" json:"rerankModel,omitempty"`
}

// BotDocumentCollection Bot-知识库关联实体
type BotDocumentCollection struct {
	ID           int64  `db:"id" json:"id,string"`
	BotID        int64  `db:"bot_id" json:"botId,string"`
	KnowledgeID  int64  `db:"knowledge_id" json:"knowledgeId,string"`

	// 关联数据
	DocumentCollection *DocumentCollection `db:"-" json:"documentCollection,omitempty"`
}

// DocumentCollectionListResponse 知识库列表响应
type DocumentCollectionListResponse struct {
	*DocumentCollection
	DocumentCount int `json:"documentCount"`
}
