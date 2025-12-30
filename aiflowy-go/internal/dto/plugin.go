package dto

// PluginSaveRequest 保存插件请求
type PluginSaveRequest struct {
	ID          string `json:"id,omitempty"`
	Alias       string `json:"alias"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Type        int    `json:"type"`
	BaseURL     string `json:"baseUrl"`
	AuthType    string `json:"authType"`  // apiKey/none
	Icon        string `json:"icon"`
	Position    string `json:"position"` // headers/query
	Headers     string `json:"headers"`
	TokenKey    string `json:"tokenKey"`
	TokenValue  string `json:"tokenValue"`
}

// PluginItemSaveRequest 保存插件工具请求
type PluginItemSaveRequest struct {
	ID            string `json:"id,omitempty"`
	PluginID      string `json:"pluginId" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Description   string `json:"description"`
	BasePath      string `json:"basePath"`
	Status        int    `json:"status"`
	InputData     string `json:"inputData"`
	OutputData    string `json:"outputData"`
	RequestMethod string `json:"requestMethod"`
	ServiceStatus int    `json:"serviceStatus"`
	DebugStatus   int    `json:"debugStatus"`
	EnglishName   string `json:"englishName"`
}

// PluginItemTestRequest 插件工具测试请求
type PluginItemTestRequest struct {
	PluginToolID string `json:"pluginToolId" validate:"required"`
	InputData    string `json:"inputData"`
}

// BotPluginUpdateRequest 更新 Bot-插件关联请求
type BotPluginUpdateRequest struct {
	BotID         string   `json:"botId" validate:"required"`
	PluginToolIDs []string `json:"pluginToolIds"`
}

// PluginCategorySaveRequest 保存插件分类请求
type PluginCategorySaveRequest struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name" validate:"required"`
}

// PluginParam 插件参数定义
type PluginParam struct {
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Type         string         `json:"type"`          // String/Number/Boolean/Array/Object/File
	Required     bool           `json:"required"`
	Enabled      bool           `json:"enabled"`
	Method       string         `json:"method"`        // query/body/header/path
	DefaultValue interface{}    `json:"defaultValue"`
	Children     []*PluginParam `json:"children,omitempty"`
}

// PluginHeader 插件请求头
type PluginHeader struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
