package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/media"
	"goyavision/pkg/apperr"
)

type GetAssetHandler struct {
	uow port.UnitOfWork
}

func NewGetAssetHandler(uow port.UnitOfWork) *GetAssetHandler {
	return &GetAssetHandler{uow: uow}
}

func (h *GetAssetHandler) Handle(ctx context.Context, query dto.GetAssetQuery) (*media.Asset, error) {
	var asset *media.Asset
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		var err error
		asset, err = repos.Assets.Get(ctx, query.ID)
		return err
	})

	if err != nil {
		if apperr.IsNotFound(err) {
			return nil, apperr.NotFound("media asset", query.ID)
		}
		return nil, apperr.Internal("get asset", err)
	}

	return asset, nil
}
