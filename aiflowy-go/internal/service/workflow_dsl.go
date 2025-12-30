package service

import (
	"encoding/json"
	"fmt"

	"github.com/aiflowy/aiflowy-go/internal/dto"
)

// WorkflowDSLParser 工作流 DSL 解析器
type WorkflowDSLParser struct{}

// NewWorkflowDSLParser 创建 DSL 解析器
func NewWorkflowDSLParser() *WorkflowDSLParser {
	return &WorkflowDSLParser{}
}

// Parse 解析工作流 JSON 定义
func (p *WorkflowDSLParser) Parse(content string) (*dto.WorkflowDefinition, error) {
	if content == "" {
		return nil, fmt.Errorf("工作流内容为空")
	}

	var definition dto.WorkflowDefinition
	if err := json.Unmarshal([]byte(content), &definition); err != nil {
		return nil, fmt.Errorf("解析工作流 JSON 失败: %w", err)
	}

	return &definition, nil
}

// GetStartParameters 获取开始节点的参数定义
func (p *WorkflowDSLParser) GetStartParameters(definition *dto.WorkflowDefinition) []*dto.WorkflowParameter {
	if definition == nil || len(definition.Nodes) == 0 {
		return nil
	}

	// 找到开始节点
	for _, node := range definition.Nodes {
		if node.Type == dto.NodeTypeStart {
			return p.extractNodeParameters(node)
		}
	}

	return nil
}

// extractNodeParameters 从节点数据中提取参数
func (p *WorkflowDSLParser) extractNodeParameters(node *dto.WorkflowNode) []*dto.WorkflowParameter {
	if node.Parameters != nil {
		return node.Parameters
	}

	// 尝试从 data 字段提取
	if node.Data == nil {
		return nil
	}

	// 查找 parameters 字段
	if params, ok := node.Data["parameters"]; ok {
		if paramList, ok := params.([]interface{}); ok {
			var result []*dto.WorkflowParameter
			for _, item := range paramList {
				if paramMap, ok := item.(map[string]interface{}); ok {
					param := &dto.WorkflowParameter{}
					if name, ok := paramMap["name"].(string); ok {
						param.Name = name
					}
					if typ, ok := paramMap["type"].(string); ok {
						param.Type = typ
					}
					if desc, ok := paramMap["description"].(string); ok {
						param.Description = desc
					}
					if required, ok := paramMap["required"].(bool); ok {
						param.Required = required
					}
					if defVal, ok := paramMap["defaultValue"]; ok {
						param.DefaultValue = defVal
					}
					result = append(result, param)
				}
			}
			return result
		}
	}

	return nil
}

// GetNodeByID 根据 ID 获取节点
func (p *WorkflowDSLParser) GetNodeByID(definition *dto.WorkflowDefinition, nodeID string) *dto.WorkflowNode {
	if definition == nil {
		return nil
	}

	for _, node := range definition.Nodes {
		if node.ID == nodeID {
			return node
		}
	}

	return nil
}

// GetNextNodes 获取指定节点的下一个节点列表
func (p *WorkflowDSLParser) GetNextNodes(definition *dto.WorkflowDefinition, nodeID string) []*dto.WorkflowNode {
	if definition == nil {
		return nil
	}

	var nextNodes []*dto.WorkflowNode
	for _, edge := range definition.Edges {
		if edge.Source == nodeID {
			if node := p.GetNodeByID(definition, edge.Target); node != nil {
				nextNodes = append(nextNodes, node)
			}
		}
	}

	return nextNodes
}

// GetPreviousNodes 获取指定节点的上一个节点列表
func (p *WorkflowDSLParser) GetPreviousNodes(definition *dto.WorkflowDefinition, nodeID string) []*dto.WorkflowNode {
	if definition == nil {
		return nil
	}

	var prevNodes []*dto.WorkflowNode
	for _, edge := range definition.Edges {
		if edge.Target == nodeID {
			if node := p.GetNodeByID(definition, edge.Source); node != nil {
				prevNodes = append(prevNodes, node)
			}
		}
	}

	return prevNodes
}

// GetStartNode 获取开始节点
func (p *WorkflowDSLParser) GetStartNode(definition *dto.WorkflowDefinition) *dto.WorkflowNode {
	if definition == nil {
		return nil
	}

	for _, node := range definition.Nodes {
		if node.Type == dto.NodeTypeStart {
			return node
		}
	}

	return nil
}

// GetEndNodes 获取结束节点列表
func (p *WorkflowDSLParser) GetEndNodes(definition *dto.WorkflowDefinition) []*dto.WorkflowNode {
	if definition == nil {
		return nil
	}

	var endNodes []*dto.WorkflowNode
	for _, node := range definition.Nodes {
		if node.Type == dto.NodeTypeEnd {
			endNodes = append(endNodes, node)
		}
	}

	return endNodes
}

// Validate 验证工作流定义
func (p *WorkflowDSLParser) Validate(definition *dto.WorkflowDefinition) error {
	if definition == nil {
		return fmt.Errorf("工作流定义为空")
	}

	if len(definition.Nodes) == 0 {
		return fmt.Errorf("工作流没有节点")
	}

	// 检查开始节点
	startNode := p.GetStartNode(definition)
	if startNode == nil {
		return fmt.Errorf("工作流没有开始节点")
	}

	// 检查结束节点
	endNodes := p.GetEndNodes(definition)
	if len(endNodes) == 0 {
		return fmt.Errorf("工作流没有结束节点")
	}

	// 检查节点 ID 唯一性
	nodeIDs := make(map[string]bool)
	for _, node := range definition.Nodes {
		if node.ID == "" {
			return fmt.Errorf("节点 ID 不能为空")
		}
		if nodeIDs[node.ID] {
			return fmt.Errorf("节点 ID 重复: %s", node.ID)
		}
		nodeIDs[node.ID] = true
	}

	// 检查边的有效性
	for _, edge := range definition.Edges {
		if !nodeIDs[edge.Source] {
			return fmt.Errorf("边的源节点不存在: %s", edge.Source)
		}
		if !nodeIDs[edge.Target] {
			return fmt.Errorf("边的目标节点不存在: %s", edge.Target)
		}
	}

	return nil
}

// ToJSON 将工作流定义转换为 JSON
func (p *WorkflowDSLParser) ToJSON(definition *dto.WorkflowDefinition) (string, error) {
	data, err := json.Marshal(definition)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
