package dto

// ModelProviderSaveRequest represents request to save model provider
type ModelProviderSaveRequest struct {
	ID           int64  `json:"id"`
	ProviderName string `json:"providerName" validate:"required"`
	ProviderType string `json:"providerType"`
	Icon         string `json:"icon"`
	APIKey       string `json:"apiKey"`
	Endpoint     string `json:"endpoint"`
	ChatPath     string `json:"chatPath"`
	EmbedPath    string `json:"embedPath"`
	RerankPath   string `json:"rerankPath"`
}

// ModelProviderListRequest represents request to list model providers
type ModelProviderListRequest struct {
	ProviderName string `query:"providerName" json:"providerName"`
	ProviderType string `query:"providerType" json:"providerType"`
}

// ModelSaveRequest represents request to save a model
type ModelSaveRequest struct {
	ID                  int64  `json:"id"`
	DeptID              int64  `json:"deptId"`
	TenantID            int64  `json:"tenantId"`
	ProviderID          int64  `json:"providerId"`
	Title               string `json:"title"`
	Icon                string `json:"icon"`
	Description         string `json:"description"`
	Endpoint            string `json:"endpoint"`
	RequestPath         string `json:"requestPath"`
	ModelName           string `json:"modelName" validate:"required"`
	APIKey              string `json:"apiKey"`
	ExtraConfig         string `json:"extraConfig"`
	Options             string `json:"options"`
	GroupName           string `json:"groupName"`
	ModelType           string `json:"modelType" validate:"required"`
	WithUsed            bool   `json:"withUsed"`
	SupportThinking     bool   `json:"supportThinking"`
	SupportTool         bool   `json:"supportTool"`
	SupportImage        bool   `json:"supportImage"`
	SupportImageB64Only bool   `json:"supportImageB64Only"`
	SupportVideo        bool   `json:"supportVideo"`
	SupportAudio        bool   `json:"supportAudio"`
	SupportFree         bool   `json:"supportFree"`
}

// ModelListRequest represents request to list models
type ModelListRequest struct {
	PageRequest
	ProviderID int64  `query:"providerId" json:"providerId"`
	ModelType  string `query:"modelType" json:"modelType"`
	WithUsed   *bool  `query:"withUsed" json:"withUsed"`
	GroupName  string `query:"groupName" json:"groupName"`
	SelectText string `query:"selectText" json:"selectText"`
}

// ModelByProviderRequest represents request to get models grouped by provider
type ModelByProviderRequest struct {
	ProviderID int64  `query:"providerId" json:"providerId"`
	ModelType  string `query:"modelType" json:"modelType"`
	WithUsed   *bool  `query:"withUsed" json:"withUsed"`
}

// AddAllLlmRequest represents request to batch add models
type AddAllLlmRequest struct {
	ProviderID int64             `json:"providerId" validate:"required"`
	GroupName  string            `json:"groupName"`
	ModelType  string            `json:"modelType"`
	Models     []ModelSaveRequest `json:"models"`
}

// UpdateByEntityRequest represents request to update models by conditions
type UpdateByEntityRequest struct {
	ProviderID int64  `json:"providerId"`
	GroupName  string `json:"groupName"`
	ModelType  string `json:"modelType"`
	// Fields to update
	WithUsed        *bool `json:"withUsed"`
	SupportThinking *bool `json:"supportThinking"`
	SupportTool     *bool `json:"supportTool"`
	SupportImage    *bool `json:"supportImage"`
}

// RemoveByEntityRequest represents request to remove models by conditions
type RemoveByEntityRequest struct {
	ProviderID int64  `json:"providerId"`
	GroupName  string `json:"groupName"`
	ModelType  string `json:"modelType"`
}

// RemoveLlmByIdsRequest represents request to remove models by IDs
type RemoveLlmByIdsRequest struct {
	IDs []int64 `json:"ids" validate:"required"`
}
