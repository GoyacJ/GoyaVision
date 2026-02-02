package app

import (
	"context"
	"errors"

	"goyavision/internal/domain"
	"goyavision/internal/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateWorkflowRequest 创建工作流请求
type CreateWorkflowRequest struct {
	Code        string                 `json:"code"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Version     string                 `json:"version,omitempty"`
	TriggerType domain.TriggerType     `json:"trigger_type"`
	TriggerConf map[string]interface{} `json:"trigger_conf,omitempty"`
	Status      domain.WorkflowStatus  `json:"status,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Nodes       []WorkflowNodeInput    `json:"nodes,omitempty"`
	Edges       []WorkflowEdgeInput    `json:"edges,omitempty"`
}

// UpdateWorkflowRequest 更新工作流请求
type UpdateWorkflowRequest struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	TriggerConf map[string]interface{} `json:"trigger_conf,omitempty"`
	Status      *domain.WorkflowStatus `json:"status,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Nodes       []WorkflowNodeInput    `json:"nodes,omitempty"`
	Edges       []WorkflowEdgeInput    `json:"edges,omitempty"`
}

// ListWorkflowsRequest 列出工作流请求
type ListWorkflowsRequest struct {
	Status      *domain.WorkflowStatus
	TriggerType *domain.TriggerType
	Tags        []string
	Keyword     string
	Limit       int
	Offset      int
}

// WorkflowNodeInput 工作流节点输入
type WorkflowNodeInput struct {
	NodeKey    string                 `json:"node_key"`
	NodeType   string                 `json:"node_type"`
	OperatorID *uuid.UUID             `json:"operator_id,omitempty"`
	Config     map[string]interface{} `json:"config,omitempty"`
	Position   map[string]interface{} `json:"position,omitempty"`
}

// WorkflowEdgeInput 工作流边输入
type WorkflowEdgeInput struct {
	SourceKey string                 `json:"source_key"`
	TargetKey string                 `json:"target_key"`
	Condition map[string]interface{} `json:"condition,omitempty"`
}

type WorkflowService struct {
	repo port.Repository
}

func NewWorkflowService(repo port.Repository) *WorkflowService {
	return &WorkflowService{
		repo: repo,
	}
}

// Create 创建工作流
func (s *WorkflowService) Create(ctx context.Context, req *CreateWorkflowRequest) (*domain.Workflow, error) {
	if req.Code == "" {
		return nil, errors.New("code is required")
	}
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.TriggerType == "" {
		return nil, errors.New("trigger_type is required")
	}

	if req.TriggerType != domain.TriggerTypeManual &&
		req.TriggerType != domain.TriggerTypeSchedule &&
		req.TriggerType != domain.TriggerTypeEvent &&
		req.TriggerType != domain.TriggerTypeAssetNew &&
		req.TriggerType != domain.TriggerTypeAssetDone {
		return nil, errors.New("invalid trigger type")
	}

	if _, err := s.repo.GetWorkflowByCode(ctx, req.Code); err == nil {
		return nil, errors.New("workflow code already exists")
	}

	version := "1.0.0"
	if req.Version != "" {
		version = req.Version
	}

	status := domain.WorkflowStatusDraft
	if req.Status != "" {
		status = req.Status
	}

	workflow := &domain.Workflow{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Version:     version,
		TriggerType: req.TriggerType,
		Status:      status,
	}

	if err := s.repo.CreateWorkflow(ctx, workflow); err != nil {
		return nil, err
	}

	if len(req.Nodes) > 0 {
		for _, nodeInput := range req.Nodes {
			if nodeInput.OperatorID != nil {
				if _, err := s.repo.GetOperator(ctx, *nodeInput.OperatorID); err != nil {
					return nil, errors.New("operator not found: " + nodeInput.NodeKey)
				}
			}

			node := &domain.WorkflowNode{
				WorkflowID: workflow.ID,
				NodeKey:    nodeInput.NodeKey,
				NodeType:   nodeInput.NodeType,
				OperatorID: nodeInput.OperatorID,
			}
			if err := s.repo.CreateWorkflowNode(ctx, node); err != nil {
				return nil, err
			}
		}
	}

	if len(req.Edges) > 0 {
		for _, edgeInput := range req.Edges {
			edge := &domain.WorkflowEdge{
				WorkflowID: workflow.ID,
				SourceKey:  edgeInput.SourceKey,
				TargetKey:  edgeInput.TargetKey,
			}
			if err := s.repo.CreateWorkflowEdge(ctx, edge); err != nil {
				return nil, err
			}
		}
	}

	return s.repo.GetWorkflowWithNodes(ctx, workflow.ID)
}

// Get 获取工作流
func (s *WorkflowService) Get(ctx context.Context, id uuid.UUID) (*domain.Workflow, error) {
	workflow, err := s.repo.GetWorkflow(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("workflow not found")
		}
		return nil, err
	}
	return workflow, nil
}

// GetWithNodes 获取工作流及其节点和边
func (s *WorkflowService) GetWithNodes(ctx context.Context, id uuid.UUID) (*domain.Workflow, error) {
	workflow, err := s.repo.GetWorkflowWithNodes(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("workflow not found")
		}
		return nil, err
	}
	return workflow, nil
}

// GetByCode 根据代码获取工作流
func (s *WorkflowService) GetByCode(ctx context.Context, code string) (*domain.Workflow, error) {
	workflow, err := s.repo.GetWorkflowByCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("workflow not found")
		}
		return nil, err
	}
	return workflow, nil
}

// List 列出工作流
func (s *WorkflowService) List(ctx context.Context, req *ListWorkflowsRequest) ([]*domain.Workflow, int64, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 1000 {
		req.Limit = 1000
	}

	filter := domain.WorkflowFilter{
		Status:      req.Status,
		TriggerType: req.TriggerType,
		Tags:        req.Tags,
		Keyword:     req.Keyword,
		Limit:       req.Limit,
		Offset:      req.Offset,
	}

	return s.repo.ListWorkflows(ctx, filter)
}

// Update 更新工作流
func (s *WorkflowService) Update(ctx context.Context, id uuid.UUID, req *UpdateWorkflowRequest) (*domain.Workflow, error) {
	workflow, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		workflow.Name = *req.Name
	}
	if req.Description != nil {
		workflow.Description = *req.Description
	}
	if req.Status != nil {
		workflow.Status = *req.Status
	}

	if len(req.Nodes) > 0 {
		if err := s.repo.DeleteWorkflowNodes(ctx, workflow.ID); err != nil {
			return nil, err
		}
		if err := s.repo.DeleteWorkflowEdges(ctx, workflow.ID); err != nil {
			return nil, err
		}

		for _, nodeInput := range req.Nodes {
			if nodeInput.OperatorID != nil {
				if _, err := s.repo.GetOperator(ctx, *nodeInput.OperatorID); err != nil {
					return nil, errors.New("operator not found: " + nodeInput.NodeKey)
				}
			}

			node := &domain.WorkflowNode{
				WorkflowID: workflow.ID,
				NodeKey:    nodeInput.NodeKey,
				NodeType:   nodeInput.NodeType,
				OperatorID: nodeInput.OperatorID,
			}
			if err := s.repo.CreateWorkflowNode(ctx, node); err != nil {
				return nil, err
			}
		}

		if len(req.Edges) > 0 {
			for _, edgeInput := range req.Edges {
				edge := &domain.WorkflowEdge{
					WorkflowID: workflow.ID,
					SourceKey:  edgeInput.SourceKey,
					TargetKey:  edgeInput.TargetKey,
				}
				if err := s.repo.CreateWorkflowEdge(ctx, edge); err != nil {
					return nil, err
				}
			}
		}
	}

	if err := s.repo.UpdateWorkflow(ctx, workflow); err != nil {
		return nil, err
	}

	return s.repo.GetWorkflowWithNodes(ctx, workflow.ID)
}

// Delete 删除工作流
func (s *WorkflowService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}

	return s.repo.DeleteWorkflow(ctx, id)
}

// Enable 启用工作流
func (s *WorkflowService) Enable(ctx context.Context, id uuid.UUID) (*domain.Workflow, error) {
	workflow, err := s.GetWithNodes(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(workflow.Nodes) == 0 {
		return nil, errors.New("workflow must have at least one node")
	}

	workflow.Status = domain.WorkflowStatusEnabled
	if err := s.repo.UpdateWorkflow(ctx, workflow); err != nil {
		return nil, err
	}

	return workflow, nil
}

// Disable 禁用工作流
func (s *WorkflowService) Disable(ctx context.Context, id uuid.UUID) (*domain.Workflow, error) {
	workflow, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	workflow.Status = domain.WorkflowStatusDisabled
	if err := s.repo.UpdateWorkflow(ctx, workflow); err != nil {
		return nil, err
	}

	return workflow, nil
}

// ListEnabled 列出所有启用的工作流
func (s *WorkflowService) ListEnabled(ctx context.Context) ([]*domain.Workflow, error) {
	return s.repo.ListEnabledWorkflows(ctx)
}
