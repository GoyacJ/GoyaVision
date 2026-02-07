package command

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/media"
	"goyavision/pkg/apperr"
)

type UpdateAssetHandler struct {
	uow port.UnitOfWork
}

func NewUpdateAssetHandler(uow port.UnitOfWork) *UpdateAssetHandler {
	return &UpdateAssetHandler{uow: uow}
}

func (h *UpdateAssetHandler) Handle(ctx context.Context, cmd dto.UpdateAssetCommand) (*media.Asset, error) {
	var asset *media.Asset
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		var err error
		asset, err = repos.Assets.Get(ctx, cmd.ID)
		if err != nil {
			if apperr.IsNotFound(err) {
				return apperr.NotFound("media asset", cmd.ID)
			}
			return apperr.Internal("get asset", err)
		}

		if cmd.Name != nil {
			asset.Name = *cmd.Name
		}
		if cmd.Status != nil {
			asset.Status = *cmd.Status
		}
		if cmd.Metadata != nil {
			asset.Metadata = cmd.Metadata
		}
		if cmd.Tags != nil {
			asset.Tags = *cmd.Tags
		}
		if cmd.Visibility != nil {
			asset.Visibility = *cmd.Visibility
		}

		return repos.Assets.Update(ctx, asset)
	})

	if err != nil {
		return nil, err
	}

	return asset, nil
}
