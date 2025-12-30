package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/schema"

	"github.com/aiflowy/aiflowy-go/internal/dto"
	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/internal/repository"
	aitool "github.com/aiflowy/aiflowy-go/internal/service/tool"
)

// WorkflowToolService 工作流工具服务
type WorkflowToolService struct {
	repo *repository.WorkflowRepository
}

// NewWorkflowToolService 创建 WorkflowToolService
func NewWorkflowToolService() *WorkflowToolService {
	return &WorkflowToolService{
		repo: repository.NewWorkflowRepository(),
	}
}

// ========================== WorkflowTool (Eino Tool 实现) ==========================

// WorkflowTool 实现 Eino 的 Tool 接口
type WorkflowTool struct {
	workflow   *entity.Workflow
	parameters []*dto.WorkflowParameter
	service    *WorkflowToolService
}

// NewWorkflowTool 创建 WorkflowTool
func NewWorkflowTool(workflow *entity.Workflow) *WorkflowTool {
	wt := &WorkflowTool{
		workflow: workflow,
		service:  NewWorkflowToolService(),
	}

	// 解析工作流参数
	if workflow.Content != "" {
		parser := NewWorkflowDSLParser()
		definition, err := parser.Parse(workflow.Content)
		if err == nil {
			wt.parameters = parser.GetStartParameters(definition)
		}
	}

	return wt
}

// Name 返回工具名称
func (t *WorkflowTool) Name() string {
	if t.workflow.EnglishName != "" {
		return t.workflow.EnglishName
	}
	return t.workflow.Title
}

// Description 返回工具描述
func (t *WorkflowTool) Description() string {
	return t.workflow.Description
}

// Parameters 返回参数定义
func (t *WorkflowTool) Parameters() map[string]*schema.ParameterInfo {
	params := make(map[string]*schema.ParameterInfo)

	for _, p := range t.parameters {
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

// Execute 执行工具 (工作流)
func (t *WorkflowTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// TODO: Stage 11 将实现工作流执行引擎
	// 目前返回占位信息
	return map[string]interface{}{
		"status":  "pending",
		"message": "工作流执行功能将在 Stage 11 实现",
		"workflow": map[string]interface{}{
			"id":    fmt.Sprintf("%d", t.workflow.ID),
			"title": t.workflow.Title,
		},
		"input": args,
	}, nil
}

// ========================== 加载 Bot 工作流工具 ==========================

// LoadBotWorkflowTools 加载 Bot 关联的工作流工具并注册到 Registry
func (s *WorkflowToolService) LoadBotWorkflowTools(ctx context.Context, botID int64) ([]*schema.ToolInfo, error) {
	// 获取 Bot 关联的工作流
	workflows, err := s.repo.ListWorkflowsByBotID(ctx, botID)
	if err != nil {
		return nil, err
	}

	if len(workflows) == 0 {
		return nil, nil
	}

	var toolInfos []*schema.ToolInfo

	for _, workflow := range workflows {
		// 创建 WorkflowTool
		workflowTool := NewWorkflowTool(workflow)

		// 转换为 Eino ToolInfo
		toolInfo := &schema.ToolInfo{
			Name: workflowTool.Name(),
			Desc: workflowTool.Description(),
			ParamsOneOf: schema.NewParamsOneOfByParams(
				convertWorkflowParams(workflowTool.Parameters()),
			),
		}
		toolInfos = append(toolInfos, toolInfo)

		// 注册到全局 Registry (用于执行)
		registry := aitool.GetRegistry()
		registry.Register(workflowTool)
	}

	return toolInfos, nil
}

// convertWorkflowParams 转换参数格式
func convertWorkflowParams(params map[string]*schema.ParameterInfo) map[string]*schema.ParameterInfo {
	return params
}
