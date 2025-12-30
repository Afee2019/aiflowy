package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"

	"github.com/aiflowy/aiflowy-go/internal/dto"
	"github.com/aiflowy/aiflowy-go/internal/entity"
	apierrors "github.com/aiflowy/aiflowy-go/internal/errors"
	"github.com/aiflowy/aiflowy-go/internal/repository"
	aitool "github.com/aiflowy/aiflowy-go/internal/service/tool"
)

// PluginToolService 插件工具服务
type PluginToolService struct {
	repo *repository.PluginRepository
}

// NewPluginToolService 创建 PluginToolService
func NewPluginToolService() *PluginToolService {
	return &PluginToolService{
		repo: repository.NewPluginRepository(),
	}
}

// ExecutePluginTool 执行插件工具
func (s *PluginToolService) ExecutePluginTool(ctx context.Context, pluginToolID string, inputData string) (interface{}, error) {
	idInt, err := strconv.ParseInt(pluginToolID, 10, 64)
	if err != nil {
		return nil, apierrors.BadRequest("无效的插件工具 ID")
	}

	// 获取插件工具
	item, err := s.repo.GetPluginItemByID(ctx, idInt)
	if err != nil {
		return nil, apierrors.InternalError("获取插件工具失败")
	}
	if item == nil {
		return nil, apierrors.NotFound("插件工具不存在")
	}

	// 获取插件
	plugin, err := s.repo.GetPluginByID(ctx, item.PluginID)
	if err != nil {
		return nil, apierrors.InternalError("获取插件失败")
	}
	if plugin == nil {
		return nil, apierrors.NotFound("插件不存在")
	}

	// 解析输入参数
	var argsMap map[string]interface{}
	if inputData != "" {
		// 输入可能是参数数组格式或直接的参数对象
		var params []*dto.PluginParam
		if err := json.Unmarshal([]byte(inputData), &params); err == nil {
			// 转换为 map
			argsMap = make(map[string]interface{})
			for _, p := range params {
				if p.Enabled && p.DefaultValue != nil {
					argsMap[p.Name] = p.DefaultValue
				}
			}
		} else {
			// 尝试直接解析为 map
			json.Unmarshal([]byte(inputData), &argsMap)
		}
	}

	// 执行 HTTP 请求
	return s.executeHTTPRequest(ctx, plugin, item, argsMap)
}

// ExecutePluginToolWithArgs 使用 LLM 参数执行插件工具
func (s *PluginToolService) ExecutePluginToolWithArgs(ctx context.Context, pluginToolID int64, args map[string]interface{}) (interface{}, error) {
	// 获取插件工具
	item, err := s.repo.GetPluginItemByID(ctx, pluginToolID)
	if err != nil {
		return nil, fmt.Errorf("获取插件工具失败: %w", err)
	}
	if item == nil {
		return nil, fmt.Errorf("插件工具不存在")
	}

	// 获取插件
	plugin, err := s.repo.GetPluginByID(ctx, item.PluginID)
	if err != nil {
		return nil, fmt.Errorf("获取插件失败: %w", err)
	}
	if plugin == nil {
		return nil, fmt.Errorf("插件不存在")
	}

	// 执行 HTTP 请求
	return s.executeHTTPRequest(ctx, plugin, item, args)
}

// executeHTTPRequest 执行 HTTP 请求
func (s *PluginToolService) executeHTTPRequest(ctx context.Context, plugin *entity.Plugin, item *entity.PluginItem, args map[string]interface{}) (interface{}, error) {
	// 构建 URL
	baseURL := plugin.BaseURL
	path := item.BasePath
	if path == "" {
		path = "/" + item.Name
	}
	fullURL := strings.TrimRight(baseURL, "/") + path

	// 解析输入参数定义
	var params []*dto.PluginParam
	if item.InputData != "" {
		json.Unmarshal([]byte(item.InputData), &params)
	}

	// 分类参数
	queryParams := url.Values{}
	bodyParams := make(map[string]interface{})
	headers := make(map[string]string)

	// 设置插件请求头
	if plugin.Headers != "" {
		var headerList []dto.PluginHeader
		if err := json.Unmarshal([]byte(plugin.Headers), &headerList); err == nil {
			for _, h := range headerList {
				headers[h.Label] = h.Value
			}
		}
	}

	// 设置认证
	if plugin.AuthType == "apiKey" && plugin.TokenKey != "" && plugin.TokenValue != "" {
		if plugin.Position == "headers" {
			headers[plugin.TokenKey] = plugin.TokenValue
		} else {
			queryParams.Set(plugin.TokenKey, plugin.TokenValue)
		}
	}

	// 处理参数
	for _, p := range params {
		if !p.Enabled {
			continue
		}

		// 获取参数值：优先使用 LLM 传递的值，否则使用默认值
		var value interface{}
		if args != nil {
			if v, ok := args[p.Name]; ok {
				value = v
			}
		}
		if value == nil && p.DefaultValue != nil {
			value = p.DefaultValue
		}
		if value == nil {
			continue
		}

		// 根据参数位置分类
		switch strings.ToLower(p.Method) {
		case "query":
			queryParams.Set(p.Name, fmt.Sprintf("%v", value))
		case "body":
			bodyParams[p.Name] = value
		case "header":
			headers[p.Name] = fmt.Sprintf("%v", value)
		case "path":
			fullURL = strings.ReplaceAll(fullURL, "{"+p.Name+"}", fmt.Sprintf("%v", value))
		}
	}

	// 添加 query 参数到 URL
	if len(queryParams) > 0 {
		if strings.Contains(fullURL, "?") {
			fullURL += "&" + queryParams.Encode()
		} else {
			fullURL += "?" + queryParams.Encode()
		}
	}

	// 构建请求
	method := strings.ToUpper(item.RequestMethod)
	if method == "" {
		method = "GET"
	}

	var body io.Reader
	if method != "GET" && method != "HEAD" && len(bodyParams) > 0 {
		bodyBytes, _ := json.Marshal(bodyParams)
		body = bytes.NewReader(bodyBytes)
		if headers["Content-Type"] == "" {
			headers["Content-Type"] = "application/json"
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 发送请求
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 尝试解析为 JSON
	var result interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		// 返回原始字符串
		return string(respBody), nil
	}
	return result, nil
}

// ========================== PluginTool (Eino Tool 实现) ==========================

// PluginTool 实现 Eino 的 Tool 接口
type PluginTool struct {
	pluginItem *entity.PluginItem
	plugin     *entity.Plugin
	params     []*dto.PluginParam
	service    *PluginToolService
}

// NewPluginTool 创建 PluginTool
func NewPluginTool(plugin *entity.Plugin, item *entity.PluginItem) *PluginTool {
	var params []*dto.PluginParam
	if item.InputData != "" {
		json.Unmarshal([]byte(item.InputData), &params)
	}

	return &PluginTool{
		pluginItem: item,
		plugin:     plugin,
		params:     params,
		service:    NewPluginToolService(),
	}
}

// Name 返回工具名称
func (t *PluginTool) Name() string {
	if t.pluginItem.EnglishName != "" {
		return t.pluginItem.EnglishName
	}
	return t.pluginItem.Name
}

// Description 返回工具描述
func (t *PluginTool) Description() string {
	return t.pluginItem.Description
}

// Parameters 返回参数定义
func (t *PluginTool) Parameters() map[string]*schema.ParameterInfo {
	params := make(map[string]*schema.ParameterInfo)

	for _, p := range t.params {
		if !p.Enabled {
			continue
		}

		var paramType schema.DataType = schema.String
		switch strings.ToLower(p.Type) {
		case "number", "integer", "int":
			paramType = schema.Integer
		case "float", "double":
			paramType = schema.Number
		case "boolean", "bool":
			paramType = schema.Boolean
		case "array":
			paramType = schema.Array
		case "object":
			paramType = schema.Object
		}

		params[p.Name] = &schema.ParameterInfo{
			Type: paramType,
			Desc: p.Description,
		}
		if p.Required {
			params[p.Name].Required = true
		}
	}

	return params
}

// Execute 执行工具
func (t *PluginTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return t.service.executeHTTPRequest(ctx, t.plugin, t.pluginItem, args)
}

// ========================== 加载 Bot 插件工具 ==========================

// LoadBotPluginTools 加载 Bot 关联的插件工具并注册到 Registry
func (s *PluginToolService) LoadBotPluginTools(ctx context.Context, botID int64) ([]*schema.ToolInfo, error) {
	// 获取 Bot 关联的插件工具
	items, err := s.repo.ListPluginItemsByBotID(ctx, botID)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, nil
	}

	var toolInfos []*schema.ToolInfo

	for _, item := range items {
		// 获取插件
		plugin, err := s.repo.GetPluginByID(ctx, item.PluginID)
		if err != nil || plugin == nil {
			continue
		}

		// 创建 PluginTool
		pluginTool := NewPluginTool(plugin, item)

		// 转换为 Eino ToolInfo
		toolInfo := &schema.ToolInfo{
			Name: pluginTool.Name(),
			Desc: pluginTool.Description(),
			ParamsOneOf: schema.NewParamsOneOfByParams(
				convertToSchemaParams(pluginTool.Parameters()),
			),
		}
		toolInfos = append(toolInfos, toolInfo)

		// 注册到全局 Registry (用于执行)
		registry := aitool.GetRegistry()
		registry.Register(pluginTool)
	}

	return toolInfos, nil
}

// convertToSchemaParams 转换参数格式
func convertToSchemaParams(params map[string]*schema.ParameterInfo) map[string]*schema.ParameterInfo {
	return params
}
