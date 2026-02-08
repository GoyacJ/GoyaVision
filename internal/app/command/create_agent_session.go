package command

import (
	"context"
	"errors"
	"time"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/agent"
	"goyavision/pkg/apperr"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateAgentSessionHandler struct {
	uow port.UnitOfWork
}

func NewCreateAgentSessionHandler(uow port.UnitOfWork) *CreateAgentSessionHandler {
	return &CreateAgentSessionHandler{uow: uow}
}

func (h *CreateAgentSessionHandler) Handle(ctx context.Context, cmd dto.CreateAgentSessionCommand) (*agent.Session, error) {
	if cmd.TaskID == uuid.Nil {
		return nil, apperr.InvalidInput("task_id is required")
	}

	var result *agent.Session
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if repos.AgentSessions == nil {
			return apperr.Internal("agent session repository is not configured", nil)
		}
		if _, err := repos.Tasks.Get(ctx, cmd.TaskID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("task", cmd.TaskID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get task")
		}

		now := time.Now().UTC()
		session := &agent.Session{
			ID:        uuid.New(),
			TaskID:    cmd.TaskID,
			Status:    agent.SessionStatusRunning,
			Budget:    cmd.Budget,
			StepCount: 0,
			StartedAt: now,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := repos.AgentSessions.Create(ctx, session); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to create agent session")
		}
		result = session
		return nil
	})
	return result, err
}
