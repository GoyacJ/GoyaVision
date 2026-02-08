package dto

import (
	"time"

	"goyavision/internal/domain/media"
	"goyavision/internal/domain/workflow"

	"github.com/google/uuid"
)

// TaskListQuery 列出任务查询参数
type TaskListQuery struct {
	WorkflowID *uuid.UUID `query:"workflow_id"`
	AssetID    *uuid.UUID `query:"asset_id"`
	Status     *string    `query:"status"`
	From       *int64     `query:"from"`
	To         *int64     `query:"to"`
	Limit      int        `query:"limit"`
	Offset     int        `query:"offset"`
}

// TaskCreateReq 创建任务请求
type TaskCreateReq struct {
	WorkflowID  uuid.UUID              `json:"workflow_id" validate:"required"`
	AssetID     *uuid.UUID             `json:"asset_id,omitempty"`
	InputParams map[string]interface{} `json:"input_params,omitempty"`
}

// TaskUpdateReq 更新任务请求
type TaskUpdateReq struct {
	Status      *string `json:"status,omitempty"`
	Progress    *int    `json:"progress,omitempty"`
	CurrentNode *string `json:"current_node,omitempty"`
	Error       *string `json:"error,omitempty"`
}

// NodeExecutionDTO 节点执行状态 DTO
type NodeExecutionDTO struct {
	NodeKey     string      `json:"node_key"`
	Status      string      `json:"status"`
	Error       string      `json:"error,omitempty"`
	StartedAt   *time.Time  `json:"started_at,omitempty"`
	CompletedAt *time.Time  `json:"completed_at,omitempty"`
	ArtifactIDs []uuid.UUID `json:"artifact_ids,omitempty"`
}

// TaskResponse 任务响应
type TaskResponse struct {
	ID             uuid.UUID              `json:"id"`
	WorkflowID     uuid.UUID              `json:"workflow_id"`
	AssetID        *uuid.UUID             `json:"asset_id,omitempty"`
	Status         string                 `json:"status"`
	Progress       int                    `json:"progress"`
	CurrentNode    string                 `json:"current_node,omitempty"`
	InputParams    map[string]interface{} `json:"input_params,omitempty"`
	Error          string                 `json:"error,omitempty"`
	NodeExecutions []NodeExecutionDTO     `json:"node_executions,omitempty"`
	Duration       float64                `json:"duration"`
	StartedAt      *time.Time             `json:"started_at,omitempty"`
	CompletedAt    *time.Time             `json:"completed_at,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// TaskWithRelationsResponse 任务及关联数据响应
type TaskWithRelationsResponse struct {
	ID             uuid.UUID              `json:"id"`
	WorkflowID     uuid.UUID              `json:"workflow_id"`
	Workflow       *WorkflowResponse      `json:"workflow,omitempty"`
	AssetID        *uuid.UUID             `json:"asset_id,omitempty"`
	Asset          *AssetResponse         `json:"asset,omitempty"`
	Status         string                 `json:"status"`
	Progress       int                    `json:"progress"`
	CurrentNode    string                 `json:"current_node,omitempty"`
	InputParams    map[string]interface{} `json:"input_params,omitempty"`
	Error          string                 `json:"error,omitempty"`
	NodeExecutions []NodeExecutionDTO     `json:"node_executions,omitempty"`
	Duration       float64                `json:"duration"`
	StartedAt      *time.Time             `json:"started_at,omitempty"`
	CompletedAt    *time.Time             `json:"completed_at,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// TaskListResponse 任务列表响应
type TaskListResponse struct {
	Items []*TaskResponse `json:"items"`
	Total int64           `json:"total"`
}

// TaskStatsResponse 任务统计响应
type TaskStatsResponse struct {
	Total     int64 `json:"total"`
	Pending   int64 `json:"pending"`
	Running   int64 `json:"running"`
	Success   int64 `json:"success"`
	Failed    int64 `json:"failed"`
	Cancelled int64 `json:"cancelled"`
}

func nodeExecutionsToDTOs(execs []workflow.NodeExecution) []NodeExecutionDTO {
	if len(execs) == 0 {
		return nil
	}
	dtos := make([]NodeExecutionDTO, len(execs))
	for i, e := range execs {
		dtos[i] = NodeExecutionDTO{
			NodeKey:     e.NodeKey,
			Status:      string(e.Status),
			Error:       e.Error,
			StartedAt:   e.StartedAt,
			CompletedAt: e.CompletedAt,
			ArtifactIDs: e.ArtifactIDs,
		}
	}
	return dtos
}

// TaskToResponse 转换为响应
func TaskToResponse(t *workflow.Task) *TaskResponse {
	if t == nil {
		return nil
	}

	inputParams := t.InputParams
	if inputParams == nil {
		inputParams = make(map[string]interface{})
	}

	return &TaskResponse{
		ID:             t.ID,
		WorkflowID:     t.WorkflowID,
		AssetID:        t.AssetID,
		Status:         string(t.Status),
		Progress:       t.Progress,
		CurrentNode:    t.CurrentNode,
		InputParams:    inputParams,
		Error:          t.Error,
		NodeExecutions: nodeExecutionsToDTOs(t.NodeExecutions),
		Duration:       t.Duration(),
		StartedAt:      t.StartedAt,
		CompletedAt:    t.CompletedAt,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}
}

// TaskToResponseWithRelations 转换为包含关联数据的响应
func TaskToResponseWithRelations(t *workflow.Task, workflow *workflow.Workflow, asset *media.Asset, minioEndpoint, minioBucket, minioPublicBase string, minioUseSSL bool) *TaskWithRelationsResponse {
	if t == nil {
		return nil
	}

	inputParams := t.InputParams
	if inputParams == nil {
		inputParams = make(map[string]interface{})
	}

	resp := &TaskWithRelationsResponse{
		ID:             t.ID,
		WorkflowID:     t.WorkflowID,
		AssetID:        t.AssetID,
		Status:         string(t.Status),
		Progress:       t.Progress,
		CurrentNode:    t.CurrentNode,
		InputParams:    inputParams,
		Error:          t.Error,
		NodeExecutions: nodeExecutionsToDTOs(t.NodeExecutions),
		Duration:       t.Duration(),
		StartedAt:      t.StartedAt,
		CompletedAt:    t.CompletedAt,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}

	if workflow != nil {
		resp.Workflow = WorkflowToResponse(workflow)
	}

	if asset != nil {
		resp.Asset = AssetToResponse(asset, minioEndpoint, minioBucket, minioPublicBase, minioUseSSL)
	}

	return resp
}

// TasksToResponse 转换为响应列表
func TasksToResponse(tasks []*workflow.Task) []*TaskResponse {
	result := make([]*TaskResponse, len(tasks))
	for i, t := range tasks {
		result[i] = TaskToResponse(t)
	}
	return result
}

// TaskStatsToResponse 转换统计为响应
func TaskStatsToResponse(stats *workflow.TaskStats) *TaskStatsResponse {
	if stats == nil {
		return nil
	}
	return &TaskStatsResponse{
		Total:     stats.Total,
		Pending:   stats.Pending,
		Running:   stats.Running,
		Success:   stats.Success,
		Failed:    stats.Failed,
		Cancelled: stats.Cancelled,
	}
}
