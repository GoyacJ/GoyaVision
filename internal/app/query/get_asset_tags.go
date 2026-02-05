package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/pkg/apperr"
)

type GetAssetTagsHandler struct {
	uow port.UnitOfWork
}

func NewGetAssetTagsHandler(uow port.UnitOfWork) *GetAssetTagsHandler {
	return &GetAssetTagsHandler{uow: uow}
}

func (h *GetAssetTagsHandler) Handle(ctx context.Context, query dto.GetAssetTagsQuery) (*dto.AssetTagsResult, error) {
	var tags []string
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		var err error
		tags, err = repos.Assets.GetAllTags(ctx)
		return err
	})

	if err != nil {
		return nil, apperr.Internal("get asset tags", err)
	}

	return &dto.AssetTagsResult{
		Tags: tags,
	}, nil
}
