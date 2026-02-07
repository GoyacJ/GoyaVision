package command

import (
	"context"
	"errors"
	"fmt"
	"time"

	"goyavision/internal/app/dto"
	appport "goyavision/internal/app/port"
	"goyavision/internal/domain/operator"
	coreport "goyavision/internal/port"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type TestOperatorHandler struct {
	uow      appport.UnitOfWork
	registry coreport.ExecutorRegistry
}

func NewTestOperatorHandler(uow appport.UnitOfWork, registry coreport.ExecutorRegistry) *TestOperatorHandler {
	return &TestOperatorHandler{uow: uow, registry: registry}
}

func (h *TestOperatorHandler) Handle(ctx context.Context, cmd dto.TestOperatorCommand) (*dto.TestOperatorResult, error) {
	result := &dto.TestOperatorResult{Success: true, Message: "test passed"}

	err := h.uow.Do(ctx, func(ctx context.Context, repos *appport.Repositories) error {
		if h.registry == nil {
			return apperr.Internal("executor registry is not configured", nil)
		}

		op, err := repos.Operators.GetWithActiveVersion(ctx, cmd.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator", cmd.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get operator")
		}

		if op.ActiveVersion == nil {
			return apperr.InvalidInput("operator has no active version")
		}

		executor, err := h.registry.Get(op.ActiveVersion.ExecMode)
		if err != nil {
			return apperr.Wrap(err, apperr.CodeInvalidInput, "executor not found for active version")
		}

		if err := executor.HealthCheck(ctx, op.ActiveVersion); err != nil {
			return apperr.Wrap(err, apperr.CodeInvalidInput, "operator health check failed")
		}

		start := time.Now()
		input := &operator.Input{
			Params: cmd.Params,
		}
		if cmd.AssetID != nil {
			input.AssetID = *cmd.AssetID
		}

		output, err := executor.Execute(ctx, op.ActiveVersion, input)
		if err != nil {
			// 将具体错误信息包含在消息中，以便前端展示给用户
			return apperr.Wrap(err, apperr.CodeInternal, fmt.Sprintf("operator test execution failed: %v", err))
		}

		result.Diagnostics = map[string]interface{}{
			"operator_id":       op.ID.String(),
			"operator_status":   string(op.Status),
			"version_id":        op.ActiveVersion.ID.String(),
			"version":           op.ActiveVersion.Version,
			"exec_mode":         string(op.ActiveVersion.ExecMode),
			"duration_ms":       time.Since(start).Milliseconds(),
			"output_assets":     len(output.OutputAssets),
			"result_items":      len(output.Results),
			"timeline_events":   len(output.Timeline),
			"has_diagnostics":   output.Diagnostics != nil,
			"health_check_pass": true,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}
