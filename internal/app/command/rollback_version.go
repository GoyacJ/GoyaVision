package command

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/operator"
)

type RollbackVersionHandler struct {
	uow port.UnitOfWork
}

func NewRollbackVersionHandler(uow port.UnitOfWork) *RollbackVersionHandler {
	return &RollbackVersionHandler{uow: uow}
}

func (h *RollbackVersionHandler) Handle(ctx context.Context, cmd dto.RollbackVersionCommand) (*operator.Operator, error) {
	activateHandler := NewActivateVersionHandler(h.uow)
	return activateHandler.Handle(ctx, dto.ActivateVersionCommand{
		OperatorID: cmd.OperatorID,
		VersionID:  cmd.VersionID,
	})
}
