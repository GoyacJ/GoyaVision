package command

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/media"
	"goyavision/pkg/apperr"
)

type CreateSourceHandler struct {
	uow     port.UnitOfWork
	gateway port.MediaGateway
}

func NewCreateSourceHandler(uow port.UnitOfWork, gw port.MediaGateway) *CreateSourceHandler {
	return &CreateSourceHandler{uow: uow, gateway: gw}
}

func (h *CreateSourceHandler) Handle(ctx context.Context, cmd dto.CreateSourceCommand) (*media.Source, error) {
	if cmd.Name == "" {
		return nil, apperr.InvalidInput("name is required")
	}
	if cmd.Type != media.SourceTypePull && cmd.Type != media.SourceTypePush {
		return nil, apperr.InvalidInput("invalid source type")
	}

	source := cmd.URL
	if cmd.Type == media.SourceTypePush {
		source = "publisher"
	} else if cmd.URL == "" {
		return nil, apperr.InvalidInput("url is required for pull source")
	}

	pathName := media.GeneratePathName(cmd.Name)

	if err := h.gateway.AddPath(ctx, pathName, source); err != nil {
		return nil, apperr.Wrap(err, apperr.CodeServiceUnavailable, "mediamtx add path failed")
	}

	var src *media.Source
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		src = &media.Source{
			Name:     cmd.Name,
			PathName: pathName,
			Type:     cmd.Type,
			URL:      cmd.URL,
			Protocol: cmd.Protocol,
			Enabled:  cmd.Enabled,
		}
		return repos.Sources.Create(ctx, src)
	})

	if err != nil {
		_ = h.gateway.DeletePath(ctx, pathName)
		return nil, apperr.Internal("create source in db", err)
	}

	return src, nil
}
