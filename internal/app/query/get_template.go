package query

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/operator"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type GetTemplateHandler struct {
	uow port.UnitOfWork
}

func NewGetTemplateHandler(uow port.UnitOfWork) *GetTemplateHandler {
	return &GetTemplateHandler{uow: uow}
}

func (h *GetTemplateHandler) Handle(ctx context.Context, q dto.GetTemplateQuery) (*operator.OperatorTemplate, error) {
	var tpl *operator.OperatorTemplate
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		item, err := repos.OperatorTemplates.Get(ctx, q.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator template", q.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get template")
		}
		tpl = item
		return nil
	})

	return tpl, err
}
