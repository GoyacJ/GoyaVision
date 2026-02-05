package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/media"
	"goyavision/pkg/apperr"
)

type ListAssetChildrenHandler struct {
	uow port.UnitOfWork
}

func NewListAssetChildrenHandler(uow port.UnitOfWork) *ListAssetChildrenHandler {
	return &ListAssetChildrenHandler{uow: uow}
}

func (h *ListAssetChildrenHandler) Handle(ctx context.Context, query dto.ListAssetChildrenQuery) (*dto.AssetListResult, error) {
	var children []*media.Asset
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		var err error
		children, err = repos.Assets.ListByParent(ctx, query.ParentID)
		return err
	})

	if err != nil {
		return nil, apperr.Internal("list asset children", err)
	}

	return &dto.AssetListResult{
		Assets: children,
		Total:  int64(len(children)),
	}, nil
}
