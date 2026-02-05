package command

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/pkg/apperr"
)

type DeleteSourceHandler struct {
	uow     port.UnitOfWork
	gateway port.MediaGateway
}

func NewDeleteSourceHandler(uow port.UnitOfWork, gw port.MediaGateway) *DeleteSourceHandler {
	return &DeleteSourceHandler{uow: uow, gateway: gw}
}

func (h *DeleteSourceHandler) Handle(ctx context.Context, cmd dto.DeleteSourceCommand) error {
	var pathName string
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		assets, err := repos.Assets.ListBySource(ctx, cmd.ID)
		if err != nil {
			return apperr.Internal("list assets by source", err)
		}
		if len(assets) > 0 {
			return apperr.HasRelation("存在关联的流媒体资产，请先删除或解除关联")
		}

		src, err := repos.Sources.Get(ctx, cmd.ID)
		if err != nil {
			if apperr.IsNotFound(err) {
				return apperr.NotFound("media source", cmd.ID)
			}
			return apperr.Internal("get source", err)
		}

		pathName = src.PathName

		if err := h.gateway.DeletePath(ctx, pathName); err != nil {
			return apperr.Wrap(err, apperr.CodeServiceUnavailable, "mediamtx delete path failed")
		}

		return repos.Sources.Delete(ctx, cmd.ID)
	})

	return err
}
