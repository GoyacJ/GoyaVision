package command

import (
	"context"
	"errors"
	"time"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/agent"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type StopAgentSessionHandler struct {
	uow port.UnitOfWork
}

func NewStopAgentSessionHandler(uow port.UnitOfWork) *StopAgentSessionHandler {
	return &StopAgentSessionHandler{uow: uow}
}

func (h *StopAgentSessionHandler) Handle(ctx context.Context, cmd dto.StopAgentSessionCommand) (*agent.Session, error) {
	status := cmd.Status
	if status == "" {
		status = agent.SessionStatusCancelled
	}

	var result *agent.Session
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if repos.AgentSessions == nil {
			return apperr.Internal("agent session repository is not configured", nil)
		}
		session, err := repos.AgentSessions.Get(ctx, cmd.SessionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("agent_session", cmd.SessionID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get agent session")
		}
		now := time.Now().UTC()
		session.Status = status
		session.EndedAt = &now
		session.UpdatedAt = now
		if err := repos.AgentSessions.Update(ctx, session); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to stop agent session")
		}
		result = session
		return nil
	})
	return result, err
}
