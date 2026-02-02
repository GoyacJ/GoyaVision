package dto

import (
	"encoding/json"
	"time"

	"goyavision/internal/domain"

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

// TaskResponse 任务响应
type TaskResponse struct {
	ID          uuid.UUID              `json:"id"`
	WorkflowID  uuid.UUID              `json:"workflow_id"`
	AssetID     *uuid.UUID             `json:"asset_id,omitempty"`
	Status      string                 `json:"status"`
	Progress    int                    `json:"progress"`
	CurrentNode string                 `json:"current_node,omitempty"`
	InputParams map[string]interface{} `json:"input_params,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Duration    float64                `json:"duration"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// TaskWithRelationsResponse 任务及关联数据响应
type TaskWithRelationsResponse struct {
	ID          uuid.UUID              `json:"id"`
	WorkflowID  uuid.UUID              `json:"workflow_id"`
	Workflow    *WorkflowResponse      `json:"workflow,omitempty"`
	AssetID     *uuid.UUID             `json:"asset_id,omitempty"`
	Asset       *AssetResponse         `json:"asset,omitempty"`
	Status      string                 `json:"status"`
	Progress    int                    `json:"progress"`
	CurrentNode string                 `json:"current_node,omitempty"`
	InputParams map[string]interface{} `json:"input_params,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Duration    float64                `json:"duration"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
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

// TaskToResponse 转换为响应
func TaskToResponse(t *domain.Task) *TaskResponse {
	if t == nil {
		return nil
	}

	var inputParams map[string]interface{}
	if t.InputParams != nil && len(t.InputParams) > 0 {
		json.Unmarshal(t.InputParams, &inputParams)
	}

	return &TaskResponse{
		ID:          t.ID,
		WorkflowID:  t.WorkflowID,
		AssetID:     t.AssetID,
		Status:      string(t.Status),
		Progress:    t.Progress,
		CurrentNode: t.CurrentNode,
		InputParams: inputParams,
		Error:       t.Error,
		Duration:    t.Duration(),
		StartedAt:   t.StartedAt,
		CompletedAt: t.CompletedAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// TaskToResponseWithRelations 转换为包含关联数据的响应
func TaskToResponseWithRelations(t *domain.Task) *TaskWithRelationsResponse {
	if t == nil {
		return nil
	}

	var inputParams map[string]interface{}
	if t.InputParams != nil && len(t.InputParams) > 0 {
		json.Unmarshal(t.InputParams, &inputParams)
	}

	resp := &TaskWithRelationsResponse{
		ID:          t.ID,
		WorkflowID:  t.WorkflowID,
		AssetID:     t.AssetID,
		Status:      string(t.Status),
		Progress:    t.Progress,
		CurrentNode: t.CurrentNode,
		InputParams: inputParams,
		Error:       t.Error,
		Duration:    t.Duration(),
		StartedAt:   t.StartedAt,
		CompletedAt: t.CompletedAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}

	if t.Workflow != nil {
		resp.Workflow = WorkflowToResponse(t.Workflow)
	}

	if t.Asset != nil {
		resp.Asset = AssetToResponse(t.Asset)
	}

	return resp
}

// TasksToResponse 转换为响应列表
func TasksToResponse(tasks []*domain.Task) []*TaskResponse {
	result := make([]*TaskResponse, len(tasks))
	for i, t := range tasks {
		result[i] = TaskToResponse(t)
	}
	return result
}

// TaskStatsToResponse 转换统计为响应
func TaskStatsToResponse(stats *domain.TaskStats) *TaskStatsResponse {
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
