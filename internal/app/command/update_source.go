package command

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/media"
	"goyavision/pkg/apperr"
)

type UpdateSourceHandler struct {
	uow     port.UnitOfWork
	gateway port.MediaGateway
}

func NewUpdateSourceHandler(uow port.UnitOfWork, gw port.MediaGateway) *UpdateSourceHandler {
	return &UpdateSourceHandler{uow: uow, gateway: gw}
}

func (h *UpdateSourceHandler) Handle(ctx context.Context, cmd dto.UpdateSourceCommand) (*media.Source, error) {
	var src *media.Source
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		var err error
		src, err = repos.Sources.Get(ctx, cmd.ID)
		if err != nil {
			if apperr.IsNotFound(err) {
				return apperr.NotFound("media source", cmd.ID)
			}
			return apperr.Internal("get source", err)
		}

		if cmd.Name != nil {
			src.Name = *cmd.Name
		}
		if cmd.URL != nil {
			src.URL = *cmd.URL
		}
		if cmd.Protocol != nil {
			src.Protocol = *cmd.Protocol
		}
		if cmd.Enabled != nil {
			src.Enabled = *cmd.Enabled
		}

		source := src.URL
		if src.Type == media.SourceTypePush {
			source = "publisher"
		}

		if err := h.gateway.PatchPath(ctx, src.PathName, source); err != nil {
			return apperr.Wrap(err, apperr.CodeServiceUnavailable, "mediamtx patch path failed")
		}

		return repos.Sources.Update(ctx, src)
	})

	if err != nil {
		return nil, err
	}

	return src, nil
}
