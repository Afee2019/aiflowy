package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/pkg/snowflake"
)

// WorkflowExecRepository 工作流执行记录数据访问层
type WorkflowExecRepository struct {
	db *sql.DB
}

// NewWorkflowExecRepository 创建 WorkflowExecRepository
func NewWorkflowExecRepository() *WorkflowExecRepository {
	return &WorkflowExecRepository{
		db: GetDB(),
	}
}

// ========================== WorkflowExecResult ==========================

// CreateExecResult 创建执行记录
func (r *WorkflowExecRepository) CreateExecResult(ctx context.Context, result *entity.WorkflowExecResult) error {
	if result.ID == 0 {
		result.ID, _ = snowflake.GenerateID()
	}

	query := `
		INSERT INTO tb_workflow_exec_result
		(id, exec_key, workflow_id, title, description, input, output, workflow_json,
		 start_time, end_time, tokens, status, created_key, created_by, error_info)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		result.ID, result.ExecKey, result.WorkflowID, result.Title, result.Description,
		result.Input, result.Output, result.WorkflowJSON, result.StartTime, result.EndTime,
		result.Tokens, result.Status, result.CreatedKey, result.CreatedBy, result.ErrorInfo,
	)
	return err
}

// UpdateExecResult 更新执行记录
func (r *WorkflowExecRepository) UpdateExecResult(ctx context.Context, result *entity.WorkflowExecResult) error {
	query := `
		UPDATE tb_workflow_exec_result
		SET output = ?, end_time = ?, tokens = ?, status = ?, error_info = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		result.Output, result.EndTime, result.Tokens, result.Status, result.ErrorInfo, result.ID,
	)
	return err
}

// GetExecResultByExecKey 根据执行 Key 获取执行记录
func (r *WorkflowExecRepository) GetExecResultByExecKey(ctx context.Context, execKey string) (*entity.WorkflowExecResult, error) {
	query := `
		SELECT id, exec_key, workflow_id, title, description, input, output, workflow_json,
		       start_time, end_time, tokens, status, created_key, created_by, error_info
		FROM tb_workflow_exec_result
		WHERE exec_key = ?
	`

	var result entity.WorkflowExecResult
	var endTime sql.NullTime
	var title, description, input, output, workflowJSON, createdKey, createdBy, errorInfo sql.NullString

	err := r.db.QueryRowContext(ctx, query, execKey).Scan(
		&result.ID, &result.ExecKey, &result.WorkflowID, &title, &description,
		&input, &output, &workflowJSON, &result.StartTime, &endTime,
		&result.Tokens, &result.Status, &createdKey, &createdBy, &errorInfo,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	result.Title = title.String
	result.Description = description.String
	result.Input = input.String
	result.Output = output.String
	result.WorkflowJSON = workflowJSON.String
	result.CreatedKey = createdKey.String
	result.CreatedBy = createdBy.String
	result.ErrorInfo = errorInfo.String
	if endTime.Valid {
		result.EndTime = &endTime.Time
	}

	return &result, nil
}

// GetExecResultByID 根据 ID 获取执行记录
func (r *WorkflowExecRepository) GetExecResultByID(ctx context.Context, id int64) (*entity.WorkflowExecResult, error) {
	query := `
		SELECT id, exec_key, workflow_id, title, description, input, output, workflow_json,
		       start_time, end_time, tokens, status, created_key, created_by, error_info
		FROM tb_workflow_exec_result
		WHERE id = ?
	`

	var result entity.WorkflowExecResult
	var endTime sql.NullTime
	var title, description, input, output, workflowJSON, createdKey, createdBy, errorInfo sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&result.ID, &result.ExecKey, &result.WorkflowID, &title, &description,
		&input, &output, &workflowJSON, &result.StartTime, &endTime,
		&result.Tokens, &result.Status, &createdKey, &createdBy, &errorInfo,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	result.Title = title.String
	result.Description = description.String
	result.Input = input.String
	result.Output = output.String
	result.WorkflowJSON = workflowJSON.String
	result.CreatedKey = createdKey.String
	result.CreatedBy = createdBy.String
	result.ErrorInfo = errorInfo.String
	if endTime.Valid {
		result.EndTime = &endTime.Time
	}

	return &result, nil
}

// ListExecResultsByWorkflowID 获取工作流的执行记录列表
func (r *WorkflowExecRepository) ListExecResultsByWorkflowID(ctx context.Context, workflowID int64) ([]*entity.WorkflowExecResult, error) {
	query := `
		SELECT id, exec_key, workflow_id, title, description, input, output, workflow_json,
		       start_time, end_time, tokens, status, created_key, created_by, error_info
		FROM tb_workflow_exec_result
		WHERE workflow_id = ?
		ORDER BY start_time DESC
		LIMIT 100
	`

	rows, err := r.db.QueryContext(ctx, query, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*entity.WorkflowExecResult
	for rows.Next() {
		var result entity.WorkflowExecResult
		var endTime sql.NullTime
		var title, description, input, output, workflowJSON, createdKey, createdBy, errorInfo sql.NullString

		err := rows.Scan(
			&result.ID, &result.ExecKey, &result.WorkflowID, &title, &description,
			&input, &output, &workflowJSON, &result.StartTime, &endTime,
			&result.Tokens, &result.Status, &createdKey, &createdBy, &errorInfo,
		)
		if err != nil {
			return nil, err
		}

		result.Title = title.String
		result.Description = description.String
		result.Input = input.String
		result.Output = output.String
		result.WorkflowJSON = workflowJSON.String
		result.CreatedKey = createdKey.String
		result.CreatedBy = createdBy.String
		result.ErrorInfo = errorInfo.String
		if endTime.Valid {
			result.EndTime = &endTime.Time
		}

		results = append(results, &result)
	}

	return results, nil
}

// ========================== WorkflowExecStep ==========================

// CreateExecStep 创建执行步骤记录
func (r *WorkflowExecRepository) CreateExecStep(ctx context.Context, step *entity.WorkflowExecStep) error {
	if step.ID == 0 {
		step.ID, _ = snowflake.GenerateID()
	}

	query := `
		INSERT INTO tb_workflow_exec_step
		(id, record_id, exec_key, node_id, node_name, input, output, node_data,
		 start_time, end_time, tokens, status, error_info)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		step.ID, step.RecordID, step.ExecKey, step.NodeID, step.NodeName,
		step.Input, step.Output, step.NodeData, step.StartTime, step.EndTime,
		step.Tokens, step.Status, step.ErrorInfo,
	)
	return err
}

// UpdateExecStep 更新执行步骤
func (r *WorkflowExecRepository) UpdateExecStep(ctx context.Context, step *entity.WorkflowExecStep) error {
	query := `
		UPDATE tb_workflow_exec_step
		SET output = ?, end_time = ?, tokens = ?, status = ?, error_info = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		step.Output, step.EndTime, step.Tokens, step.Status, step.ErrorInfo, step.ID,
	)
	return err
}

// GetExecStepsByRecordID 获取执行记录的所有步骤
func (r *WorkflowExecRepository) GetExecStepsByRecordID(ctx context.Context, recordID int64) ([]*entity.WorkflowExecStep, error) {
	query := `
		SELECT id, record_id, exec_key, node_id, node_name, input, output, node_data,
		       start_time, end_time, tokens, status, error_info
		FROM tb_workflow_exec_step
		WHERE record_id = ?
		ORDER BY start_time ASC
	`

	rows, err := r.db.QueryContext(ctx, query, recordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []*entity.WorkflowExecStep
	for rows.Next() {
		var step entity.WorkflowExecStep
		var endTime sql.NullTime
		var input, output, nodeData, errorInfo sql.NullString

		err := rows.Scan(
			&step.ID, &step.RecordID, &step.ExecKey, &step.NodeID, &step.NodeName,
			&input, &output, &nodeData, &step.StartTime, &endTime,
			&step.Tokens, &step.Status, &errorInfo,
		)
		if err != nil {
			return nil, err
		}

		step.Input = input.String
		step.Output = output.String
		step.NodeData = nodeData.String
		step.ErrorInfo = errorInfo.String
		if endTime.Valid {
			step.EndTime = &endTime.Time
		}

		steps = append(steps, &step)
	}

	return steps, nil
}

// GetExecStepsByExecKey 根据执行 Key 获取所有步骤
func (r *WorkflowExecRepository) GetExecStepsByExecKey(ctx context.Context, execKey string) ([]*entity.WorkflowExecStep, error) {
	query := `
		SELECT id, record_id, exec_key, node_id, node_name, input, output, node_data,
		       start_time, end_time, tokens, status, error_info
		FROM tb_workflow_exec_step
		WHERE exec_key = ?
		ORDER BY start_time ASC
	`

	rows, err := r.db.QueryContext(ctx, query, execKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []*entity.WorkflowExecStep
	for rows.Next() {
		var step entity.WorkflowExecStep
		var endTime sql.NullTime
		var input, output, nodeData, errorInfo sql.NullString

		err := rows.Scan(
			&step.ID, &step.RecordID, &step.ExecKey, &step.NodeID, &step.NodeName,
			&input, &output, &nodeData, &step.StartTime, &endTime,
			&step.Tokens, &step.Status, &errorInfo,
		)
		if err != nil {
			return nil, err
		}

		step.Input = input.String
		step.Output = output.String
		step.NodeData = nodeData.String
		step.ErrorInfo = errorInfo.String
		if endTime.Valid {
			step.EndTime = &endTime.Time
		}

		steps = append(steps, &step)
	}

	return steps, nil
}

// ========================== 辅助方法 ==========================

// CreateExecResultWithSteps 创建执行记录和第一个步骤
func (r *WorkflowExecRepository) CreateExecResultWithSteps(ctx context.Context, result *entity.WorkflowExecResult) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 创建执行记录
	if result.ID == 0 {
		result.ID, _ = snowflake.GenerateID()
	}

	query := `
		INSERT INTO tb_workflow_exec_result
		(id, exec_key, workflow_id, title, description, input, output, workflow_json,
		 start_time, end_time, tokens, status, created_key, created_by, error_info)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = tx.ExecContext(ctx, query,
		result.ID, result.ExecKey, result.WorkflowID, result.Title, result.Description,
		result.Input, result.Output, result.WorkflowJSON, result.StartTime, result.EndTime,
		result.Tokens, result.Status, result.CreatedKey, result.CreatedBy, result.ErrorInfo,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// marshalJSON 将对象序列化为 JSON 字符串
func marshalJSON(v interface{}) string {
	if v == nil {
		return ""
	}
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}

// unmarshalJSON 将 JSON 字符串反序列化为对象
func unmarshalJSON(s string, v interface{}) error {
	if s == "" {
		return nil
	}
	return json.Unmarshal([]byte(s), v)
}

// ========================== 执行状态缓存 (内存) ==========================

// ChainState 工作流执行状态 (内存中保持)
type ChainState struct {
	ExecuteID         string
	WorkflowID        int64
	RecordID          int64
	Status            entity.WorkflowExecStatus
	Variables         map[string]interface{}
	NodeStates        map[string]*NodeState
	Result            map[string]interface{}
	Error             error
	SuspendedNodeID   string                   // 当前暂停的节点 ID
	SuspendedParams   []*SuspendedParam        // 暂停时等待的参数
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// NodeState 节点执行状态
type NodeState struct {
	NodeID    string
	NodeName  string
	Status    entity.WorkflowExecStatus
	Input     map[string]interface{}
	Output    map[string]interface{}
	Error     error
	StartTime time.Time
	EndTime   *time.Time
	Tokens    int64
}

// SuspendedParam 暂停时等待的参数
type SuspendedParam struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}
