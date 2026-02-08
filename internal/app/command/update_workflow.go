package command

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/workflow"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type UpdateWorkflowHandler struct {
	uow             port.UnitOfWork
	schemaValidator port.SchemaValidator
}

func NewUpdateWorkflowHandler(uow port.UnitOfWork, schemaValidator port.SchemaValidator) *UpdateWorkflowHandler {
	return &UpdateWorkflowHandler{uow: uow, schemaValidator: schemaValidator}
}

func (h *UpdateWorkflowHandler) Handle(ctx context.Context, cmd dto.UpdateWorkflowCommand) (*workflow.Workflow, error) {
	var result *workflow.Workflow
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		wf, err := repos.Workflows.Get(ctx, cmd.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("workflow", cmd.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get workflow")
		}

		if cmd.Name != nil {
			wf.Name = *cmd.Name
		}
		if cmd.Description != nil {
			wf.Description = *cmd.Description
		}
		if cmd.TriggerConf != nil {
			triggerConf, err := buildTriggerConfig(cmd.TriggerConf)
			if err != nil {
				return apperr.InvalidInput(err.Error())
			}
			wf.TriggerConf = triggerConf
		}
		if cmd.ContextSpec != nil {
			contextSpec, err := parseContextSpec(cmd.ContextSpec)
			if err != nil {
				return apperr.InvalidInput(err.Error())
			}
			wf.ContextSpec = contextSpec
		}
		if cmd.Status != nil {
			wf.Status = *cmd.Status
		}
		if len(cmd.Tags) > 0 {
			wf.Tags = cmd.Tags
		}
		if cmd.Visibility != nil {
			wf.Visibility = *cmd.Visibility
		}
		if cmd.VisibleRoleIDs != nil {
			wf.VisibleRoleIDs = cmd.VisibleRoleIDs
		}

		if len(cmd.Nodes) > 0 {
			if err := validateWorkflowConnections(ctx, repos, h.schemaValidator, cmd.Nodes, cmd.Edges); err != nil {
				return err
			}

			if err := repos.Workflows.DeleteNodes(ctx, wf.ID); err != nil {
				return apperr.Wrap(err, apperr.CodeDBError, "failed to delete old nodes")
			}
			if err := repos.Workflows.DeleteEdges(ctx, wf.ID); err != nil {
				return apperr.Wrap(err, apperr.CodeDBError, "failed to delete old edges")
			}

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
					Config:     parseNodeConfig(nodeInput.Config),
					Position:   parseNodePosition(nodeInput.Position),
				}
				if err := repos.Workflows.CreateNode(ctx, node); err != nil {
					return apperr.Wrap(err, apperr.CodeDBError, "failed to create workflow node")
				}
			}

			if len(cmd.Edges) > 0 {
				for _, edgeInput := range cmd.Edges {
					edge := &workflow.Edge{
						WorkflowID: wf.ID,
						SourceKey:  edgeInput.SourceKey,
						TargetKey:  edgeInput.TargetKey,
						Condition:  parseEdgeCondition(edgeInput.Condition),
					}
					if err := repos.Workflows.CreateEdge(ctx, edge); err != nil {
						return apperr.Wrap(err, apperr.CodeDBError, "failed to create workflow edge")
					}
				}
			}
		}

		if err := repos.Workflows.Update(ctx, wf); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to update workflow")
		}

		wfWithNodes, err := repos.Workflows.GetWithNodes(ctx, wf.ID)
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get workflow with nodes")
		}
		if _, err := persistAndActivateWorkflowRevision(ctx, repos, wfWithNodes, nextWorkflowRevision(wf.CurrentRevision)); err != nil {
			return err
		}
		wfWithNodes, err = repos.Workflows.GetWithNodes(ctx, wf.ID)
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to reload workflow with revision")
		}
		result = wfWithNodes
		return nil
	})

	return result, err
}
