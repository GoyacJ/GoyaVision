package command

import (
	"context"
	"fmt"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/workflow"
	"goyavision/pkg/apperr"
)

type CreateWorkflowHandler struct {
	uow port.UnitOfWork
}

func NewCreateWorkflowHandler(uow port.UnitOfWork) *CreateWorkflowHandler {
	return &CreateWorkflowHandler{uow: uow}
}

func (h *CreateWorkflowHandler) Handle(ctx context.Context, cmd dto.CreateWorkflowCommand) (*workflow.Workflow, error) {
	if cmd.Code == "" {
		return nil, apperr.InvalidInput("code is required")
	}
	if cmd.Name == "" {
		return nil, apperr.InvalidInput("name is required")
	}
	if cmd.TriggerType == "" {
		return nil, apperr.InvalidInput("trigger_type is required")
	}

	if cmd.TriggerType != workflow.TriggerTypeManual &&
		cmd.TriggerType != workflow.TriggerTypeSchedule &&
		cmd.TriggerType != workflow.TriggerTypeEvent &&
		cmd.TriggerType != workflow.TriggerTypeAssetNew &&
		cmd.TriggerType != workflow.TriggerTypeAssetDone {
		return nil, apperr.InvalidInput("invalid trigger type")
	}

	version := "1.0.0"
	if cmd.Version != "" {
		version = cmd.Version
	}

	status := workflow.StatusDraft
	if cmd.Status != "" {
		status = cmd.Status
	}

	var result *workflow.Workflow
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if _, err := repos.Workflows.GetByCode(ctx, cmd.Code); err == nil {
			return apperr.Conflict(fmt.Sprintf("workflow with code %s already exists", cmd.Code))
		}

		wf := &workflow.Workflow{
			Code:        cmd.Code,
			Name:        cmd.Name,
			Description: cmd.Description,
			Version:     version,
			TriggerType: cmd.TriggerType,
			Status:      status,
			Tags:        cmd.Tags,
		}

		if err := repos.Workflows.Create(ctx, wf); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to create workflow")
		}

		if len(cmd.Nodes) > 0 {
			for _, nodeInput := range cmd.Nodes {
				if nodeInput.OperatorID != nil {
					if _, err := repos.Operators.Get(ctx, *nodeInput.OperatorID); err != nil {
						return apperr.NotFound("operator", nodeInput.NodeKey)
					}
				}

				node := &workflow.Node{
					WorkflowID: wf.ID,
					NodeKey:    nodeInput.NodeKey,
					NodeType:   nodeInput.NodeType,
					OperatorID: nodeInput.OperatorID,
				}
				if err := repos.Workflows.CreateNode(ctx, node); err != nil {
					return apperr.Wrap(err, apperr.CodeDBError, "failed to create workflow node")
				}
			}
		}

		if len(cmd.Edges) > 0 {
			for _, edgeInput := range cmd.Edges {
				edge := &workflow.Edge{
					WorkflowID: wf.ID,
					SourceKey:  edgeInput.SourceKey,
					TargetKey:  edgeInput.TargetKey,
				}
				if err := repos.Workflows.CreateEdge(ctx, edge); err != nil {
					return apperr.Wrap(err, apperr.CodeDBError, "failed to create workflow edge")
				}
			}
		}

		wfWithNodes, err := repos.Workflows.GetWithNodes(ctx, wf.ID)
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get workflow with nodes")
		}
		result = wfWithNodes
		return nil
	})

	return result, err
}
