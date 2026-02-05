package command

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/pkg/apperr"
)

type DeleteAssetHandler struct {
	uow port.UnitOfWork
}

func NewDeleteAssetHandler(uow port.UnitOfWork) *DeleteAssetHandler {
	return &DeleteAssetHandler{uow: uow}
}

func (h *DeleteAssetHandler) Handle(ctx context.Context, cmd dto.DeleteAssetCommand) error {
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		asset, err := repos.Assets.Get(ctx, cmd.ID)
		if err != nil {
			if apperr.IsNotFound(err) {
				return apperr.NotFound("media asset", cmd.ID)
			}
			return apperr.Internal("get asset", err)
		}

		children, err := repos.Assets.ListByParent(ctx, asset.ID)
		if err != nil {
			return apperr.Internal("list child assets", err)
		}
		if len(children) > 0 {
			return apperr.HasRelation("cannot delete asset with children")
		}

		return repos.Assets.Delete(ctx, cmd.ID)
	})

	return err
}
