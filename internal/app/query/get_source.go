package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/media"
	"goyavision/pkg/apperr"
)

type GetSourceHandler struct {
	uow port.UnitOfWork
}

func NewGetSourceHandler(uow port.UnitOfWork) *GetSourceHandler {
	return &GetSourceHandler{uow: uow}
}

func (h *GetSourceHandler) Handle(ctx context.Context, query dto.GetSourceQuery) (*media.Source, error) {
	var src *media.Source
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		var err error
		src, err = repos.Sources.Get(ctx, query.ID)
		return err
	})

	if err != nil {
		if apperr.IsNotFound(err) {
			return nil, apperr.NotFound("media source", query.ID)
		}
		return nil, apperr.Internal("get source", err)
	}

	return src, nil
}
