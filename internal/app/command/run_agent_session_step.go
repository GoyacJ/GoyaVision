package command

import (
	"context"
	"errors"

	agentruntime "goyavision/internal/app/agent"
	"goyavision/internal/app/dto"
	"goyavision/internal/domain/agent"
	"goyavision/pkg/apperr"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RunAgentSessionStepHandler struct {
	runtime *agentruntime.RunLoop
}

func NewRunAgentSessionStepHandler(runtime *agentruntime.RunLoop) *RunAgentSessionStepHandler {
	return &RunAgentSessionStepHandler{runtime: runtime}
}

func (h *RunAgentSessionStepHandler) Handle(ctx context.Context, cmd dto.RunAgentSessionStepCommand) (*agent.Session, error) {
	if cmd.SessionID == uuid.Nil {
		return nil, apperr.InvalidInput("session_id is required")
	}
	if h.runtime == nil {
		return nil, apperr.Internal("agent runtime is not configured", nil)
	}

	session, err := h.runtime.RunStep(ctx, cmd.SessionID, cmd.MaxActions)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("agent_session", cmd.SessionID.String())
		}
		return nil, apperr.Wrap(err, apperr.CodeInternal, "failed to run agent session step")
	}
	return session, nil
}
