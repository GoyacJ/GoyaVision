package command

import (
	"context"
	"encoding/json"
	"fmt"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/workflow"
	"goyavision/pkg/apperr"
)

type CreateWorkflowHandler struct {
	uow             port.UnitOfWork
	schemaValidator port.SchemaValidator
}

func NewCreateWorkflowHandler(uow port.UnitOfWork, schemaValidator port.SchemaValidator) *CreateWorkflowHandler {
	return &CreateWorkflowHandler{uow: uow, schemaValidator: schemaValidator}
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

		if err := validateWorkflowConnections(ctx, repos, h.schemaValidator, cmd.Nodes, cmd.Edges); err != nil {
			return err
		}

		triggerConf, err := buildTriggerConfig(cmd.TriggerConf)
		if err != nil {
			return apperr.InvalidInput(err.Error())
		}

		wf := &workflow.Workflow{
			Code:        cmd.Code,
			Name:        cmd.Name,
			Description: cmd.Description,
			Version:     version,
			TriggerType: cmd.TriggerType,
			TriggerConf: triggerConf,
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
					Config:     parseNodeConfig(nodeInput.Config),
					Position:   parseNodePosition(nodeInput.Position),
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
					Condition:  parseEdgeCondition(edgeInput.Condition),
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

func parseNodeConfig(raw map[string]interface{}) *workflow.NodeConfig {
	if len(raw) == 0 {
		return nil
	}
	data, err := json.Marshal(raw)
	if err != nil {
		return nil
	}
	var cfg workflow.NodeConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil
	}
	return &cfg
}

func parseNodePosition(raw map[string]interface{}) *workflow.NodePosition {
	if len(raw) == 0 {
		return nil
	}
	data, err := json.Marshal(raw)
	if err != nil {
		return nil
	}
	var pos workflow.NodePosition
	if err := json.Unmarshal(data, &pos); err != nil {
		return nil
	}
	return &pos
}

func parseEdgeCondition(raw map[string]interface{}) *workflow.EdgeCondition {
	if len(raw) == 0 {
		return nil
	}
	data, err := json.Marshal(raw)
	if err != nil {
		return nil
	}
	var cond workflow.EdgeCondition
	if err := json.Unmarshal(data, &cond); err != nil {
		return nil
	}
	return &cond
}
