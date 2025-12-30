package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/aiflowy/aiflowy-go/internal/dto"
	"github.com/aiflowy/aiflowy-go/internal/repository"
	"github.com/aiflowy/aiflowy-go/internal/service/llm"
	aitool "github.com/aiflowy/aiflowy-go/internal/service/tool"
)

// ========================== StartNodeExecutor ==========================

// StartNodeExecutor 开始节点执行器
type StartNodeExecutor struct{}

// Execute 执行开始节点
func (e *StartNodeExecutor) Execute(ctx context.Context, state *repository.ChainState, node *dto.WorkflowNode) (map[string]interface{}, error) {
	// 开始节点只是传递输入变量
	return state.Variables, nil
}

// ========================== EndNodeExecutor ==========================

// EndNodeExecutor 结束节点执行器
type EndNodeExecutor struct{}

// Execute 执行结束节点
func (e *EndNodeExecutor) Execute(ctx context.Context, state *repository.ChainState, node *dto.WorkflowNode) (map[string]interface{}, error) {
	// 结束节点返回最终结果
	result := make(map[string]interface{})

	// 从节点数据中提取输出映射
	if node.Data != nil {
		if outputs, ok := node.Data["outputs"].([]interface{}); ok {
			for _, output := range outputs {
				if outputMap, ok := output.(map[string]interface{}); ok {
					key := getStringFromMap(outputMap, "name")
					value := getStringFromMap(outputMap, "value")
					if key != "" && value != "" {
						// 尝试从变量中解析值
						result[key] = resolveVariable(value, state.Variables)
					}
				}
			}
		}
	}

	// 如果没有配置输出，返回所有变量
	if len(result) == 0 {
		return state.Variables, nil
	}

	return result, nil
}

// ========================== LLMNodeExecutor ==========================

// LLMNodeExecutor LLM 节点执行器
type LLMNodeExecutor struct {
	chatService *llm.ChatService
}

// NewLLMNodeExecutor 创建 LLM 节点执行器
func NewLLMNodeExecutor() *LLMNodeExecutor {
	return &LLMNodeExecutor{
		chatService: llm.NewChatService(),
	}
}

// Execute 执行 LLM 节点
func (e *LLMNodeExecutor) Execute(ctx context.Context, state *repository.ChainState, node *dto.WorkflowNode) (map[string]interface{}, error) {
	if node.Data == nil {
		return nil, fmt.Errorf("LLM 节点缺少配置")
	}

	// 获取模型 ID
	modelIDStr := getStringFromMap(node.Data, "modelId")
	if modelIDStr == "" {
		modelIDStr = getStringFromMap(node.Data, "llmId")
	}
	if modelIDStr == "" {
		return nil, fmt.Errorf("LLM 节点未配置模型")
	}

	modelID, err := strconv.ParseInt(modelIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的模型 ID: %s", modelIDStr)
	}

	// 获取提示词
	systemPrompt := getStringFromMap(node.Data, "prompt")
	if systemPrompt == "" {
		systemPrompt = getStringFromMap(node.Data, "systemPrompt")
	}

	// 获取用户消息模板
	userTemplate := getStringFromMap(node.Data, "userPrompt")
	if userTemplate == "" {
		userTemplate = getStringFromMap(node.Data, "message")
	}

	// 替换变量
	userMessage := resolveTemplateString(userTemplate, state.Variables)
	systemPrompt = resolveTemplateString(systemPrompt, state.Variables)

	// 构建消息
	messages := []llm.Message{}
	if systemPrompt != "" {
		messages = append(messages, llm.Message{Role: "system", Content: systemPrompt})
	}
	if userMessage != "" {
		messages = append(messages, llm.Message{Role: "user", Content: userMessage})
	}

	// 调用 LLM
	req := &llm.ChatRequest{
		ModelID:  modelID,
		Messages: messages,
	}

	response, err := e.chatService.Chat(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("LLM 调用失败: %w", err)
	}

	// 获取输出变量名
	outputVar := getStringFromMap(node.Data, "outputVariable")
	if outputVar == "" {
		outputVar = "llmOutput"
	}

	return map[string]interface{}{
		outputVar: response.Content,
	}, nil
}

// ========================== ToolNodeExecutor ==========================

// ToolNodeExecutor 工具节点执行器
type ToolNodeExecutor struct {
	registry *aitool.Registry
}

// NewToolNodeExecutor 创建工具节点执行器
func NewToolNodeExecutor() *ToolNodeExecutor {
	return &ToolNodeExecutor{
		registry: aitool.GetRegistry(),
	}
}

// Execute 执行工具节点
func (e *ToolNodeExecutor) Execute(ctx context.Context, state *repository.ChainState, node *dto.WorkflowNode) (map[string]interface{}, error) {
	if node.Data == nil {
		return nil, fmt.Errorf("工具节点缺少配置")
	}

	// 获取工具名称
	toolName := getStringFromMap(node.Data, "toolName")
	if toolName == "" {
		toolName = getStringFromMap(node.Data, "name")
	}
	if toolName == "" {
		return nil, fmt.Errorf("工具节点未配置工具名称")
	}

	// 获取工具
	tool, ok := e.registry.Get(toolName)
	if !ok || tool == nil {
		return nil, fmt.Errorf("工具不存在: %s", toolName)
	}

	// 解析参数
	args := make(map[string]interface{})
	if params, ok := node.Data["parameters"].(map[string]interface{}); ok {
		for k, v := range params {
			args[k] = resolveVariable(fmt.Sprintf("%v", v), state.Variables)
		}
	}
	if inputsVal, ok := node.Data["inputs"].([]interface{}); ok {
		for _, input := range inputsVal {
			if inputMap, ok := input.(map[string]interface{}); ok {
				name := getStringFromMap(inputMap, "name")
				value := getStringFromMap(inputMap, "value")
				if name != "" && value != "" {
					args[name] = resolveVariable(value, state.Variables)
				}
			}
		}
	}

	// 执行工具
	result, err := tool.Execute(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("执行工具 %s 失败: %w", toolName, err)
	}

	// 获取输出变量名
	outputVar := getStringFromMap(node.Data, "outputVariable")
	if outputVar == "" {
		outputVar = "toolOutput"
	}

	// 如果结果是 map，直接返回
	if resultMap, ok := result.(map[string]interface{}); ok {
		return resultMap, nil
	}

	return map[string]interface{}{
		outputVar: result,
	}, nil
}

// ========================== ConditionNodeExecutor ==========================

// ConditionNodeExecutor 条件节点执行器
type ConditionNodeExecutor struct{}

// Execute 执行条件节点
func (e *ConditionNodeExecutor) Execute(ctx context.Context, state *repository.ChainState, node *dto.WorkflowNode) (map[string]interface{}, error) {
	if node.Data == nil {
		return map[string]interface{}{"condition": "default"}, nil
	}

	// 获取条件配置
	conditions, ok := node.Data["conditions"].([]interface{})
	if !ok {
		return map[string]interface{}{"condition": "default"}, nil
	}

	// 评估条件
	for _, cond := range conditions {
		condMap, ok := cond.(map[string]interface{})
		if !ok {
			continue
		}

		name := getStringFromMap(condMap, "name")
		expression := getStringFromMap(condMap, "expression")

		if evaluateCondition(expression, state.Variables) {
			return map[string]interface{}{
				"condition": name,
			}, nil
		}
	}

	return map[string]interface{}{"condition": "default"}, nil
}

// evaluateCondition 评估条件表达式
func evaluateCondition(expression string, variables map[string]interface{}) bool {
	if expression == "" {
		return false
	}

	// 简单的条件评估
	// 支持格式: ${variable} == value, ${variable} != value, ${variable} > value 等

	// 替换变量
	resolved := resolveTemplateString(expression, variables)

	// 简单比较
	parts := strings.SplitN(resolved, "==", 2)
	if len(parts) == 2 {
		left := strings.TrimSpace(parts[0])
		right := strings.TrimSpace(parts[1])
		return left == right
	}

	parts = strings.SplitN(resolved, "!=", 2)
	if len(parts) == 2 {
		left := strings.TrimSpace(parts[0])
		right := strings.TrimSpace(parts[1])
		return left != right
	}

	// 如果表达式本身是 true/false
	return resolved == "true" || resolved == "1"
}

// ========================== HumanConfirmNodeExecutor ==========================

// HumanConfirmNodeExecutor 人工确认节点执行器
type HumanConfirmNodeExecutor struct{}

// Execute 执行人工确认节点
func (e *HumanConfirmNodeExecutor) Execute(ctx context.Context, state *repository.ChainState, node *dto.WorkflowNode) (map[string]interface{}, error) {
	// 检查是否已有确认参数
	if node.Data != nil {
		if params, ok := node.Data["confirmParameters"].([]interface{}); ok {
			// 检查是否所有确认参数都已提供
			allProvided := true
			var suspendParams []*repository.SuspendedParam

			for _, param := range params {
				if paramMap, ok := param.(map[string]interface{}); ok {
					name := getStringFromMap(paramMap, "name")
					if _, exists := state.Variables[name]; !exists {
						allProvided = false
						suspendParams = append(suspendParams, &repository.SuspendedParam{
							Name:        name,
							Type:        getStringFromMap(paramMap, "type"),
							Description: getStringFromMap(paramMap, "description"),
							Required:    getBoolFromMap(paramMap, "required"),
						})
					}
				}
			}

			if !allProvided {
				return nil, &SuspendError{Params: suspendParams}
			}
		}
	}

	// 所有确认参数已提供，继续执行
	return state.Variables, nil
}

// ========================== PluginNodeExecutor ==========================

// PluginNodeExecutor 插件节点执行器
type PluginNodeExecutor struct {
	pluginToolService *PluginToolService
}

// NewPluginNodeExecutor 创建插件节点执行器
func NewPluginNodeExecutor() *PluginNodeExecutor {
	return &PluginNodeExecutor{
		pluginToolService: NewPluginToolService(),
	}
}

// Execute 执行插件节点
func (e *PluginNodeExecutor) Execute(ctx context.Context, state *repository.ChainState, node *dto.WorkflowNode) (map[string]interface{}, error) {
	if node.Data == nil {
		return nil, fmt.Errorf("插件节点缺少配置")
	}

	// 获取插件工具 ID
	pluginToolIDStr := getStringFromMap(node.Data, "pluginToolId")
	if pluginToolIDStr == "" {
		pluginToolIDStr = getStringFromMap(node.Data, "pluginId")
	}
	if pluginToolIDStr == "" {
		return nil, fmt.Errorf("插件节点未配置插件工具 ID")
	}

	pluginToolID, err := strconv.ParseInt(pluginToolIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的插件工具 ID: %s", pluginToolIDStr)
	}

	// 解析参数
	args := make(map[string]interface{})
	if params, ok := node.Data["parameters"].(map[string]interface{}); ok {
		for k, v := range params {
			args[k] = resolveVariable(fmt.Sprintf("%v", v), state.Variables)
		}
	}

	// 执行插件工具
	result, err := e.pluginToolService.ExecutePluginToolWithArgs(ctx, pluginToolID, args)
	if err != nil {
		return nil, fmt.Errorf("执行插件工具失败: %w", err)
	}

	// 获取输出变量名
	outputVar := getStringFromMap(node.Data, "outputVariable")
	if outputVar == "" {
		outputVar = "pluginOutput"
	}

	// 如果结果是 map，直接返回
	if resultMap, ok := result.(map[string]interface{}); ok {
		return resultMap, nil
	}

	return map[string]interface{}{
		outputVar: result,
	}, nil
}

// ========================== CodeNodeExecutor ==========================

// CodeNodeExecutor 代码节点执行器
type CodeNodeExecutor struct{}

// Execute 执行代码节点 (简单实现 - 仅支持 JSON 转换)
func (e *CodeNodeExecutor) Execute(ctx context.Context, state *repository.ChainState, node *dto.WorkflowNode) (map[string]interface{}, error) {
	if node.Data == nil {
		return nil, fmt.Errorf("代码节点缺少配置")
	}

	// 获取代码类型
	codeType := getStringFromMap(node.Data, "codeType")
	code := getStringFromMap(node.Data, "code")

	result := make(map[string]interface{})

	switch codeType {
	case "json":
		// JSON 转换
		if code != "" {
			// 替换变量
			resolved := resolveTemplateString(code, state.Variables)
			if err := json.Unmarshal([]byte(resolved), &result); err != nil {
				return nil, fmt.Errorf("JSON 解析失败: %w", err)
			}
		}
	case "template":
		// 模板替换
		if code != "" {
			resolved := resolveTemplateString(code, state.Variables)
			result["output"] = resolved
		}
	default:
		// 简单变量传递
		result = state.Variables
	}

	return result, nil
}

// ========================== SubWorkflowNodeExecutor ==========================

// SubWorkflowNodeExecutor 子工作流节点执行器
type SubWorkflowNodeExecutor struct {
	executor *ChainExecutor
}

// NewSubWorkflowNodeExecutor 创建子工作流节点执行器
func NewSubWorkflowNodeExecutor(executor *ChainExecutor) *SubWorkflowNodeExecutor {
	return &SubWorkflowNodeExecutor{
		executor: executor,
	}
}

// Execute 执行子工作流节点
func (e *SubWorkflowNodeExecutor) Execute(ctx context.Context, state *repository.ChainState, node *dto.WorkflowNode) (map[string]interface{}, error) {
	if node.Data == nil {
		return nil, fmt.Errorf("子工作流节点缺少配置")
	}

	// 获取子工作流 ID
	workflowIDStr := getStringFromMap(node.Data, "workflowId")
	if workflowIDStr == "" {
		return nil, fmt.Errorf("子工作流节点未配置工作流 ID")
	}

	// 解析参数
	variables := make(map[string]interface{})
	if params, ok := node.Data["parameters"].(map[string]interface{}); ok {
		for k, v := range params {
			variables[k] = resolveVariable(fmt.Sprintf("%v", v), state.Variables)
		}
	}

	// 同步执行子工作流 (执行开始节点)
	result, err := e.executor.ExecuteNode(ctx, workflowIDStr, "", variables)
	if err != nil {
		return nil, fmt.Errorf("执行子工作流失败: %w", err)
	}

	return result, nil
}

// ========================== 辅助函数 ==========================

// getStringFromMap 从 map 中获取字符串
func getStringFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
		if s, ok := v.(float64); ok {
			return fmt.Sprintf("%.0f", s)
		}
	}
	return ""
}

// getBoolFromMap 从 map 中获取布尔值
func getBoolFromMap(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// resolveVariable 解析变量 (支持 ${var} 格式)
func resolveVariable(value string, variables map[string]interface{}) interface{} {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		varName := strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")
		if resolved, ok := variables[varName]; ok {
			return resolved
		}
	}
	return value
}

// resolveTemplateString 解析模板字符串 (替换所有 ${var})
func resolveTemplateString(template string, variables map[string]interface{}) string {
	result := template
	for key, value := range variables {
		placeholder := "${" + key + "}"
		valueStr := fmt.Sprintf("%v", value)
		result = strings.ReplaceAll(result, placeholder, valueStr)
	}
	return result
}
