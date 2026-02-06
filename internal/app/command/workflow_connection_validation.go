package command

import (
	"context"
	"fmt"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/operator"
	"goyavision/pkg/apperr"
)

func validateWorkflowConnections(
	ctx context.Context,
	repos *port.Repositories,
	validator port.SchemaValidator,
	nodes []dto.WorkflowNodeInput,
	edges []dto.WorkflowEdgeInput,
) error {
	if validator == nil {
		return apperr.ServiceUnavailable("schema validator is not configured")
	}

	nodeOperatorMap := make(map[string]*operator.Operator, len(nodes))
	for i := range nodes {
		n := nodes[i]
		if n.OperatorID == nil {
			continue
		}

		op, err := repos.Operators.GetWithActiveVersion(ctx, *n.OperatorID)
		if err != nil {
			return apperr.NotFound("operator", n.NodeKey)
		}
		nodeOperatorMap[n.NodeKey] = op
	}

	for i := range edges {
		e := edges[i]
		upstream := nodeOperatorMap[e.SourceKey]
		downstream := nodeOperatorMap[e.TargetKey]
		if upstream == nil || downstream == nil {
			continue
		}

		upstreamOutputSpec := getOperatorOutputSpec(upstream)
		if len(upstreamOutputSpec) == 0 {
			return apperr.InvalidInput(fmt.Sprintf("workflow connection invalid: 上游节点 %s 的输出 Schema 缺失", e.SourceKey))
		}

		downstreamInputSchema := getOperatorInputSchema(downstream)
		if len(downstreamInputSchema) == 0 {
			return apperr.InvalidInput(fmt.Sprintf("workflow connection invalid: 下游节点 %s 的输入 Schema 缺失", e.TargetKey))
		}

		if err := validator.ValidateConnection(ctx, upstreamOutputSpec, downstreamInputSchema); err != nil {
			return apperr.Wrap(
				err,
				apperr.CodeInvalidInput,
				fmt.Sprintf("workflow connection invalid: %s -> %s", e.SourceKey, e.TargetKey),
			)
		}
	}

	return nil
}

func getOperatorInputSchema(op *operator.Operator) map[string]interface{} {
	if op == nil {
		return nil
	}
	if op.ActiveVersion == nil {
		return nil
	}
	return op.ActiveVersion.InputSchema
}

func getOperatorOutputSpec(op *operator.Operator) map[string]interface{} {
	if op == nil {
		return nil
	}
	if op.ActiveVersion == nil {
		return nil
	}
	return op.ActiveVersion.OutputSpec
}
