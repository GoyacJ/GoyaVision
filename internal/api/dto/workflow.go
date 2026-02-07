package dto

import (
	"encoding/json"
	"time"

	"goyavision/internal/domain/workflow"

	"github.com/google/uuid"
)

// WorkflowListQuery 列出工作流查询参数
type WorkflowListQuery struct {
	Status      *string `query:"status"`
	TriggerType *string `query:"trigger_type"`
	Tags        *string `query:"tags"`
	Keyword     *string `query:"keyword"`
	Limit       int     `query:"limit"`
	Offset      int     `query:"offset"`
}

// WorkflowCreateReq 创建工作流请求
type WorkflowCreateReq struct {
	Code        string                 `json:"code" validate:"required"`
	Name        string                 `json:"name" validate:"required"`
	Description string                 `json:"description,omitempty"`
	Version     string                 `json:"version,omitempty"`
	TriggerType string                 `json:"trigger_type" validate:"required"`
	TriggerConf map[string]interface{} `json:"trigger_conf,omitempty"`
	Status         string                 `json:"status,omitempty"`
	Tags           []string               `json:"tags,omitempty"`
	Nodes          []WorkflowNodeInput    `json:"nodes,omitempty"`
	Edges          []WorkflowEdgeInput    `json:"edges,omitempty"`
	Visibility     *int                   `json:"visibility,omitempty"`
	VisibleRoleIDs []string               `json:"visible_role_ids,omitempty"`
}

// WorkflowUpdateReq 更新工作流请求
type WorkflowUpdateReq struct {
	Name           *string                `json:"name,omitempty"`
	Description    *string                `json:"description,omitempty"`
	TriggerConf    map[string]interface{} `json:"trigger_conf,omitempty"`
	Status         *string                `json:"status,omitempty"`
	Tags           []string               `json:"tags,omitempty"`
	Nodes          []WorkflowNodeInput    `json:"nodes,omitempty"`
	Edges          []WorkflowEdgeInput    `json:"edges,omitempty"`
	Visibility     *int                   `json:"visibility,omitempty"`
	VisibleRoleIDs []string               `json:"visible_role_ids,omitempty"`
}

// WorkflowNodeInput 工作流节点输入
type WorkflowNodeInput struct {
	NodeKey    string                 `json:"node_key" validate:"required"`
	NodeType   string                 `json:"node_type" validate:"required"`
	OperatorID *uuid.UUID             `json:"operator_id,omitempty"`
	Config     map[string]interface{} `json:"config,omitempty"`
	Position   map[string]interface{} `json:"position,omitempty"`
}

// WorkflowEdgeInput 工作流边输入
type WorkflowEdgeInput struct {
	SourceKey string                 `json:"source_key" validate:"required"`
	TargetKey string                 `json:"target_key" validate:"required"`
	Condition map[string]interface{} `json:"condition,omitempty"`
}

// WorkflowResponse 工作流响应
type WorkflowResponse struct {
	ID          uuid.UUID              `json:"id"`
	Code        string                 `json:"code"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Version     string                 `json:"version"`
	TriggerType string                 `json:"trigger_type"`
	TriggerConf    map[string]interface{} `json:"trigger_conf,omitempty"`
	Status         string                 `json:"status"`
	Tags           []string               `json:"tags,omitempty"`
	Visibility     int                    `json:"visibility"`
	VisibleRoleIDs []string               `json:"visible_role_ids,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// WorkflowWithNodesResponse 工作流及节点响应
type WorkflowWithNodesResponse struct {
	ID          uuid.UUID              `json:"id"`
	Code        string                 `json:"code"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Version     string                 `json:"version"`
	TriggerType string                 `json:"trigger_type"`
	TriggerConf    map[string]interface{} `json:"trigger_conf,omitempty"`
	Status         string                 `json:"status"`
	Tags           []string               `json:"tags,omitempty"`
	Visibility     int                    `json:"visibility"`
	VisibleRoleIDs []string               `json:"visible_role_ids,omitempty"`
	Nodes          []WorkflowNodeResponse `json:"nodes"`
	Edges          []WorkflowEdgeResponse `json:"edges"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// WorkflowNodeResponse 工作流节点响应
type WorkflowNodeResponse struct {
	ID         uuid.UUID              `json:"id"`
	NodeKey    string                 `json:"node_key"`
	NodeType   string                 `json:"node_type"`
	OperatorID *uuid.UUID             `json:"operator_id,omitempty"`
	Operator   *OperatorResponse      `json:"operator,omitempty"`
	Config     map[string]interface{} `json:"config,omitempty"`
	Position   map[string]interface{} `json:"position,omitempty"`
}

// WorkflowEdgeResponse 工作流边响应
type WorkflowEdgeResponse struct {
	ID        uuid.UUID              `json:"id"`
	SourceKey string                 `json:"source_key"`
	TargetKey string                 `json:"target_key"`
	Condition map[string]interface{} `json:"condition,omitempty"`
}

// WorkflowListResponse 工作流列表响应
type WorkflowListResponse struct {
	Items []*WorkflowResponse `json:"items"`
	Total int64               `json:"total"`
}

// WorkflowToResponse 转换为响应
func WorkflowToResponse(w *workflow.Workflow) *WorkflowResponse {
	if w == nil {
		return nil
	}

	var triggerConf map[string]interface{}
	if w.TriggerConf != nil {
		triggerConfBytes, _ := json.Marshal(w.TriggerConf)
		json.Unmarshal(triggerConfBytes, &triggerConf)
	}

	tags := w.Tags
	if tags == nil {
		tags = []string{}
	}

	return &WorkflowResponse{
		ID:          w.ID,
		Code:        w.Code,
		Name:        w.Name,
		Description: w.Description,
		Version:     w.Version,
		TriggerType:    string(w.TriggerType),
		TriggerConf:    triggerConf,
		Status:         string(w.Status),
		Tags:           tags,
		Visibility:     int(w.Visibility),
		VisibleRoleIDs: w.VisibleRoleIDs,
		CreatedAt:      w.CreatedAt,
		UpdatedAt:      w.UpdatedAt,
	}
}

// WorkflowToResponseWithNodes 转换为包含节点的响应
func WorkflowToResponseWithNodes(w *workflow.Workflow) *WorkflowWithNodesResponse {
	if w == nil {
		return nil
	}

	var triggerConf map[string]interface{}
	if w.TriggerConf != nil {
		triggerConfBytes, _ := json.Marshal(w.TriggerConf)
		json.Unmarshal(triggerConfBytes, &triggerConf)
	}

	tags := w.Tags
	if tags == nil {
		tags = []string{}
	}

	nodes := make([]WorkflowNodeResponse, 0, len(w.Nodes))
	for _, n := range w.Nodes {
		var config map[string]interface{}
		if n.Config != nil {
			configBytes, _ := json.Marshal(n.Config)
			json.Unmarshal(configBytes, &config)
		}

		var position map[string]interface{}
		if n.Position != nil {
			positionBytes, _ := json.Marshal(n.Position)
			json.Unmarshal(positionBytes, &position)
		}

		nodeResp := WorkflowNodeResponse{
			ID:         n.ID,
			NodeKey:    n.NodeKey,
			NodeType:   n.NodeType,
			OperatorID: n.OperatorID,
			Config:     config,
			Position:   position,
		}

		nodes = append(nodes, nodeResp)
	}

	edges := make([]WorkflowEdgeResponse, 0, len(w.Edges))
	for _, e := range w.Edges {
		var condition map[string]interface{}
		if e.Condition != nil {
			conditionBytes, _ := json.Marshal(e.Condition)
			json.Unmarshal(conditionBytes, &condition)
		}

		edges = append(edges, WorkflowEdgeResponse{
			ID:        e.ID,
			SourceKey: e.SourceKey,
			TargetKey: e.TargetKey,
			Condition: condition,
		})
	}

	return &WorkflowWithNodesResponse{
		ID:          w.ID,
		Code:        w.Code,
		Name:        w.Name,
		Description: w.Description,
		Version:     w.Version,
		TriggerType:    string(w.TriggerType),
		TriggerConf:    triggerConf,
		Status:         string(w.Status),
		Tags:           tags,
		Visibility:     int(w.Visibility),
		VisibleRoleIDs: w.VisibleRoleIDs,
		Nodes:          nodes,
		Edges:          edges,
		CreatedAt:      w.CreatedAt,
		UpdatedAt:      w.UpdatedAt,
	}
}

// WorkflowsToResponse 转换为响应列表
func WorkflowsToResponse(workflows []*workflow.Workflow) []*WorkflowResponse {
	result := make([]*WorkflowResponse, len(workflows))
	for i, w := range workflows {
		result[i] = WorkflowToResponse(w)
	}
	return result
}
