package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/aiflowy/aiflowy-go/internal/dto"
	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/internal/repository"
)

// ChainExecutor 工作流执行引擎
type ChainExecutor struct {
	workflowRepo *repository.WorkflowRepository
	execRepo     *repository.WorkflowExecRepository
	parser       *WorkflowDSLParser

	// 执行状态缓存 (内存)
	states sync.Map // map[string]*repository.ChainState

	// 节点执行器
	nodeExecutors map[string]NodeExecutor
}

// NodeExecutor 节点执行器接口
type NodeExecutor interface {
	Execute(ctx context.Context, state *repository.ChainState, node *dto.WorkflowNode) (map[string]interface{}, error)
}

// NewChainExecutor 创建工作流执行引擎
func NewChainExecutor() *ChainExecutor {
	executor := &ChainExecutor{
		workflowRepo:  repository.NewWorkflowRepository(),
		execRepo:      repository.NewWorkflowExecRepository(),
		parser:        NewWorkflowDSLParser(),
		nodeExecutors: make(map[string]NodeExecutor),
	}

	// 注册节点执行器
	executor.registerNodeExecutors()

	return executor
}

// registerNodeExecutors 注册内置节点执行器
func (e *ChainExecutor) registerNodeExecutors() {
	e.nodeExecutors[dto.NodeTypeStart] = &StartNodeExecutor{}
	e.nodeExecutors[dto.NodeTypeEnd] = &EndNodeExecutor{}
	e.nodeExecutors[dto.NodeTypeLLM] = NewLLMNodeExecutor()
	e.nodeExecutors[dto.NodeTypeTool] = NewToolNodeExecutor()
	e.nodeExecutors[dto.NodeTypeCondition] = &ConditionNodeExecutor{}
	e.nodeExecutors[dto.NodeTypeHumanConfirm] = &HumanConfirmNodeExecutor{}
	e.nodeExecutors[dto.NodeTypePlugin] = NewPluginNodeExecutor()
	e.nodeExecutors[dto.NodeTypeCode] = &CodeNodeExecutor{}
	e.nodeExecutors[dto.NodeTypeWorkflow] = NewSubWorkflowNodeExecutor(e)
}

// ExecuteAsync 异步执行工作流
func (e *ChainExecutor) ExecuteAsync(ctx context.Context, workflowID string, variables map[string]interface{}, userID, createdBy string) (string, error) {
	// 生成执行 ID
	executeID := uuid.New().String()

	// 解析工作流 ID
	wfID, err := strconv.ParseInt(workflowID, 10, 64)
	if err != nil {
		return "", fmt.Errorf("无效的工作流 ID: %s", workflowID)
	}

	// 加载工作流
	workflow, err := e.workflowRepo.GetWorkflowByID(ctx, wfID)
	if err != nil {
		return "", fmt.Errorf("加载工作流失败: %w", err)
	}
	if workflow == nil {
		return "", fmt.Errorf("工作流不存在: %s", workflowID)
	}

	// 解析工作流定义
	definition, err := e.parser.Parse(workflow.Content)
	if err != nil {
		return "", fmt.Errorf("解析工作流定义失败: %w", err)
	}

	// 验证工作流
	if err := e.parser.Validate(definition); err != nil {
		return "", fmt.Errorf("工作流定义无效: %w", err)
	}

	// 确保 variables 不为 nil
	if variables == nil {
		variables = make(map[string]interface{})
	}

	// 初始化执行状态
	state := &repository.ChainState{
		ExecuteID:  executeID,
		WorkflowID: workflow.ID,
		Status:     entity.ExecStatusRunning,
		Variables:  variables,
		NodeStates: make(map[string]*repository.NodeState),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 存储状态到内存
	e.states.Store(executeID, state)

	// 创建执行记录
	inputJSON, _ := json.Marshal(variables)
	execResult := &entity.WorkflowExecResult{
		ExecKey:      executeID,
		WorkflowID:   workflow.ID,
		Title:        workflow.Title,
		Description:  workflow.Description,
		Input:        string(inputJSON),
		WorkflowJSON: workflow.Content,
		StartTime:    time.Now(),
		Status:       entity.ExecStatusRunning,
		CreatedKey:   userID,
		CreatedBy:    createdBy,
	}

	if err := e.execRepo.CreateExecResult(ctx, execResult); err != nil {
		return "", fmt.Errorf("创建执行记录失败: %w", err)
	}

	state.RecordID = execResult.ID

	// 异步执行
	go e.executeWorkflow(context.Background(), state, definition)

	return executeID, nil
}

// executeWorkflow 执行工作流
func (e *ChainExecutor) executeWorkflow(ctx context.Context, state *repository.ChainState, definition *dto.WorkflowDefinition) {
	defer func() {
		if r := recover(); r != nil {
			e.handleError(ctx, state, fmt.Errorf("执行异常: %v", r))
		}
	}()

	// 找到开始节点
	startNode := e.parser.GetStartNode(definition)
	if startNode == nil {
		e.handleError(ctx, state, fmt.Errorf("工作流没有开始节点"))
		return
	}

	// 从开始节点执行
	e.executeNode(ctx, state, definition, startNode)
}

// executeNode 执行单个节点
func (e *ChainExecutor) executeNode(ctx context.Context, state *repository.ChainState, definition *dto.WorkflowDefinition, node *dto.WorkflowNode) {
	// 检查状态是否已终止
	if state.Status == entity.ExecStatusFailed || state.Status == entity.ExecStatusCompleted {
		return
	}

	// 创建节点状态
	nodeName := node.Name
	if nodeName == "" {
		nodeName = node.Type
	}

	nodeState := &repository.NodeState{
		NodeID:    node.ID,
		NodeName:  nodeName,
		Status:    entity.ExecStatusRunning,
		Input:     state.Variables,
		StartTime: time.Now(),
	}
	state.NodeStates[node.ID] = nodeState
	state.UpdatedAt = time.Now()

	// 记录步骤开始
	nodeDataJSON, _ := json.Marshal(node.Data)
	inputJSON, _ := json.Marshal(state.Variables)
	step := &entity.WorkflowExecStep{
		RecordID:  state.RecordID,
		ExecKey:   state.ExecuteID,
		NodeID:    node.ID,
		NodeName:  nodeName,
		Input:     string(inputJSON),
		NodeData:  string(nodeDataJSON),
		StartTime: time.Now(),
		Status:    entity.ExecStatusRunning,
	}
	e.execRepo.CreateExecStep(ctx, step)

	// 获取节点执行器
	executor, ok := e.nodeExecutors[node.Type]
	if !ok {
		e.handleNodeError(ctx, state, nodeState, step, fmt.Errorf("未知的节点类型: %s", node.Type))
		return
	}

	// 执行节点
	result, err := executor.Execute(ctx, state, node)

	if err != nil {
		// 检查是否是暂停错误
		if suspendErr, ok := err.(*SuspendError); ok {
			e.handleSuspend(ctx, state, nodeState, step, node.ID, suspendErr.Params)
			return
		}
		e.handleNodeError(ctx, state, nodeState, step, err)
		return
	}

	// 更新节点状态
	now := time.Now()
	nodeState.Status = entity.ExecStatusCompleted
	nodeState.Output = result
	nodeState.EndTime = &now

	// 合并输出到变量
	if state.Variables == nil {
		state.Variables = make(map[string]interface{})
	}
	for k, v := range result {
		state.Variables[k] = v
	}

	// 更新步骤记录
	outputJSON, _ := json.Marshal(result)
	step.Output = string(outputJSON)
	step.EndTime = &now
	step.Status = entity.ExecStatusCompleted
	e.execRepo.UpdateExecStep(ctx, step)

	// 检查是否是结束节点
	if node.Type == dto.NodeTypeEnd {
		e.handleComplete(ctx, state, result)
		return
	}

	// 获取下一个节点
	nextNodes := e.parser.GetNextNodes(definition, node.ID)
	if len(nextNodes) == 0 {
		// 没有下一个节点，工作流完成
		e.handleComplete(ctx, state, result)
		return
	}

	// 条件节点：根据条件选择下一个节点
	if node.Type == dto.NodeTypeCondition {
		selectedNode := e.selectNextNodeByCondition(definition, node, result)
		if selectedNode != nil {
			e.executeNode(ctx, state, definition, selectedNode)
		} else {
			e.handleComplete(ctx, state, result)
		}
		return
	}

	// 执行下一个节点
	for _, nextNode := range nextNodes {
		e.executeNode(ctx, state, definition, nextNode)
	}
}

// selectNextNodeByCondition 根据条件选择下一个节点
func (e *ChainExecutor) selectNextNodeByCondition(definition *dto.WorkflowDefinition, node *dto.WorkflowNode, result map[string]interface{}) *dto.WorkflowNode {
	// 获取条件结果
	conditionResult, ok := result["condition"].(string)
	if !ok {
		conditionResult = "default"
	}

	// 遍历边，找到匹配的条件
	for _, edge := range definition.Edges {
		if edge.Source == node.ID {
			// 检查条件是否匹配
			if edge.Condition == conditionResult || edge.Condition == "" || edge.SourcePort == conditionResult {
				return e.parser.GetNodeByID(definition, edge.Target)
			}
		}
	}

	// 找不到匹配的条件，返回第一个下一节点
	for _, edge := range definition.Edges {
		if edge.Source == node.ID {
			return e.parser.GetNodeByID(definition, edge.Target)
		}
	}

	return nil
}

// handleSuspend 处理工作流暂停
func (e *ChainExecutor) handleSuspend(ctx context.Context, state *repository.ChainState, nodeState *repository.NodeState, step *entity.WorkflowExecStep, nodeID string, params []*repository.SuspendedParam) {
	state.Status = entity.ExecStatusSuspended
	state.SuspendedNodeID = nodeID
	state.SuspendedParams = params
	state.UpdatedAt = time.Now()

	nodeState.Status = entity.ExecStatusSuspended

	step.Status = entity.ExecStatusSuspended
	now := time.Now()
	step.EndTime = &now
	e.execRepo.UpdateExecStep(ctx, step)

	// 更新执行记录
	execResult, _ := e.execRepo.GetExecResultByExecKey(ctx, state.ExecuteID)
	if execResult != nil {
		execResult.Status = entity.ExecStatusSuspended
		execResult.EndTime = &now
		e.execRepo.UpdateExecResult(ctx, execResult)
	}
}

// handleError 处理执行错误
func (e *ChainExecutor) handleError(ctx context.Context, state *repository.ChainState, err error) {
	state.Status = entity.ExecStatusFailed
	state.Error = err
	state.UpdatedAt = time.Now()

	// 更新执行记录
	execResult, _ := e.execRepo.GetExecResultByExecKey(ctx, state.ExecuteID)
	if execResult != nil {
		now := time.Now()
		execResult.Status = entity.ExecStatusFailed
		execResult.ErrorInfo = err.Error()
		execResult.EndTime = &now
		e.execRepo.UpdateExecResult(ctx, execResult)
	}
}

// handleNodeError 处理节点执行错误
func (e *ChainExecutor) handleNodeError(ctx context.Context, state *repository.ChainState, nodeState *repository.NodeState, step *entity.WorkflowExecStep, err error) {
	now := time.Now()

	nodeState.Status = entity.ExecStatusFailed
	nodeState.Error = err
	nodeState.EndTime = &now

	step.Status = entity.ExecStatusFailed
	step.ErrorInfo = err.Error()
	step.EndTime = &now
	e.execRepo.UpdateExecStep(ctx, step)

	e.handleError(ctx, state, err)
}

// handleComplete 处理执行完成
func (e *ChainExecutor) handleComplete(ctx context.Context, state *repository.ChainState, result map[string]interface{}) {
	state.Status = entity.ExecStatusCompleted
	state.Result = result
	state.UpdatedAt = time.Now()

	// 更新执行记录
	execResult, _ := e.execRepo.GetExecResultByExecKey(ctx, state.ExecuteID)
	if execResult != nil {
		now := time.Now()
		execResult.Status = entity.ExecStatusCompleted
		outputJSON, _ := json.Marshal(result)
		execResult.Output = string(outputJSON)
		execResult.EndTime = &now
		e.execRepo.UpdateExecResult(ctx, execResult)
	}
}

// GetChainStatus 获取执行状态
func (e *ChainExecutor) GetChainStatus(ctx context.Context, executeID string, requestNodes []*dto.NodeInfo) (*dto.ChainInfo, error) {
	// 尝试从内存获取状态
	if stateVal, ok := e.states.Load(executeID); ok {
		state := stateVal.(*repository.ChainState)
		return e.buildChainInfo(state, requestNodes), nil
	}

	// 从数据库加载
	execResult, err := e.execRepo.GetExecResultByExecKey(ctx, executeID)
	if err != nil {
		return nil, fmt.Errorf("加载执行记录失败: %w", err)
	}
	if execResult == nil {
		return nil, fmt.Errorf("执行记录不存在: %s", executeID)
	}

	// 构建 ChainInfo
	chainInfo := &dto.ChainInfo{
		ExecuteID: executeID,
		Status:    int(execResult.Status),
		Nodes:     make(map[string]*dto.NodeInfo),
	}

	if execResult.ErrorInfo != "" {
		chainInfo.Message = execResult.ErrorInfo
	}

	if execResult.Output != "" {
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(execResult.Output), &result); err == nil {
			chainInfo.Result = result
		}
	}

	// 加载步骤状态
	steps, _ := e.execRepo.GetExecStepsByExecKey(ctx, executeID)
	for _, step := range steps {
		nodeInfo := &dto.NodeInfo{
			NodeID:   step.NodeID,
			NodeName: step.NodeName,
			Status:   int(step.Status),
		}
		if step.ErrorInfo != "" {
			nodeInfo.Message = step.ErrorInfo
		}
		if step.Output != "" {
			var result map[string]interface{}
			if err := json.Unmarshal([]byte(step.Output), &result); err == nil {
				nodeInfo.Result = result
			}
		}
		chainInfo.Nodes[step.NodeID] = nodeInfo
	}

	return chainInfo, nil
}

// buildChainInfo 构建 ChainInfo
func (e *ChainExecutor) buildChainInfo(state *repository.ChainState, requestNodes []*dto.NodeInfo) *dto.ChainInfo {
	chainInfo := &dto.ChainInfo{
		ExecuteID: state.ExecuteID,
		Status:    int(state.Status),
		Result:    state.Result,
		Nodes:     make(map[string]*dto.NodeInfo),
	}

	if state.Error != nil {
		chainInfo.Message = state.Error.Error()
	}

	// 填充节点状态
	for nodeID, nodeState := range state.NodeStates {
		nodeInfo := &dto.NodeInfo{
			NodeID:   nodeID,
			NodeName: nodeState.NodeName,
			Status:   int(nodeState.Status),
			Result:   nodeState.Output,
		}
		if nodeState.Error != nil {
			nodeInfo.Message = nodeState.Error.Error()
		}
		chainInfo.Nodes[nodeID] = nodeInfo
	}

	// 处理暂停参数
	if state.Status == entity.ExecStatusSuspended && len(state.SuspendedParams) > 0 {
		if nodeInfo, ok := chainInfo.Nodes[state.SuspendedNodeID]; ok {
			nodeInfo.SuspendForParameters = make([]*dto.WorkflowParameter, len(state.SuspendedParams))
			for i, param := range state.SuspendedParams {
				nodeInfo.SuspendForParameters[i] = &dto.WorkflowParameter{
					Name:        param.Name,
					Type:        param.Type,
					Description: param.Description,
					Required:    param.Required,
				}
			}
		}
	}

	return chainInfo
}

// Resume 恢复执行
func (e *ChainExecutor) Resume(ctx context.Context, executeID string, confirmParams map[string]interface{}) error {
	// 从内存获取状态
	stateVal, ok := e.states.Load(executeID)
	if !ok {
		return fmt.Errorf("执行记录不存在或已完成: %s", executeID)
	}

	state := stateVal.(*repository.ChainState)

	if state.Status != entity.ExecStatusSuspended {
		return fmt.Errorf("工作流未处于暂停状态")
	}

	// 合并确认参数
	for k, v := range confirmParams {
		state.Variables[k] = v
	}

	// 重新加载工作流定义
	workflow, err := e.workflowRepo.GetWorkflowByID(ctx, state.WorkflowID)
	if err != nil {
		return fmt.Errorf("加载工作流失败: %w", err)
	}

	definition, err := e.parser.Parse(workflow.Content)
	if err != nil {
		return fmt.Errorf("解析工作流定义失败: %w", err)
	}

	// 获取暂停的节点
	suspendedNodeID := state.SuspendedNodeID
	state.Status = entity.ExecStatusRunning
	state.SuspendedNodeID = ""
	state.SuspendedParams = nil
	state.UpdatedAt = time.Now()

	// 更新执行记录状态
	execResult, _ := e.execRepo.GetExecResultByExecKey(ctx, executeID)
	if execResult != nil {
		execResult.Status = entity.ExecStatusRunning
		e.execRepo.UpdateExecResult(ctx, execResult)
	}

	// 获取下一个节点并继续执行
	nextNodes := e.parser.GetNextNodes(definition, suspendedNodeID)
	if len(nextNodes) == 0 {
		e.handleComplete(ctx, state, state.Variables)
		return nil
	}

	// 异步继续执行
	go func() {
		for _, nextNode := range nextNodes {
			e.executeNode(context.Background(), state, definition, nextNode)
		}
	}()

	return nil
}

// ExecuteNode 执行单个节点 (用于调试)
func (e *ChainExecutor) ExecuteNode(ctx context.Context, workflowID, nodeID string, variables map[string]interface{}) (map[string]interface{}, error) {
	// 解析工作流 ID
	wfID, err := strconv.ParseInt(workflowID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的工作流 ID: %s", workflowID)
	}

	// 加载工作流
	workflow, err := e.workflowRepo.GetWorkflowByID(ctx, wfID)
	if err != nil {
		return nil, fmt.Errorf("加载工作流失败: %w", err)
	}
	if workflow == nil {
		return nil, fmt.Errorf("工作流不存在: %s", workflowID)
	}

	// 解析工作流定义
	definition, err := e.parser.Parse(workflow.Content)
	if err != nil {
		return nil, fmt.Errorf("解析工作流定义失败: %w", err)
	}

	// 查找节点
	node := e.parser.GetNodeByID(definition, nodeID)
	if node == nil {
		return nil, fmt.Errorf("节点不存在: %s", nodeID)
	}

	// 确保 variables 不为 nil
	if variables == nil {
		variables = make(map[string]interface{})
	}

	// 创建临时状态
	state := &repository.ChainState{
		ExecuteID:  uuid.New().String(),
		WorkflowID: workflow.ID,
		Status:     entity.ExecStatusRunning,
		Variables:  variables,
		NodeStates: make(map[string]*repository.NodeState),
	}

	// 获取节点执行器
	executor, ok := e.nodeExecutors[node.Type]
	if !ok {
		return nil, fmt.Errorf("未知的节点类型: %s", node.Type)
	}

	// 执行节点
	return executor.Execute(ctx, state, node)
}

// ========================== 错误类型 ==========================

// SuspendError 暂停错误 (用于人工确认节点)
type SuspendError struct {
	Params []*repository.SuspendedParam
}

func (e *SuspendError) Error() string {
	return "workflow suspended for confirmation"
}

// ========================== 全局执行器实例 ==========================

var (
	chainExecutorOnce sync.Once
	chainExecutor     *ChainExecutor
)

// GetChainExecutor 获取全局执行器实例
func GetChainExecutor() *ChainExecutor {
	chainExecutorOnce.Do(func() {
		chainExecutor = NewChainExecutor()
	})
	return chainExecutor
}
