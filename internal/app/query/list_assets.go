package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/media"
	"goyavision/pkg/apperr"
)

type ListAssetsHandler struct {
	uow port.UnitOfWork
}

func NewListAssetsHandler(uow port.UnitOfWork) *ListAssetsHandler {
	return &ListAssetsHandler{uow: uow}
}

func (h *ListAssetsHandler) Handle(ctx context.Context, query dto.ListAssetsQuery) (*dto.PagedResult[*media.Asset], error) {
	query.Pagination.Normalize()

	filter := media.AssetFilter{
		Type:       query.Type,
		SourceType: query.SourceType,
		SourceID:   query.SourceID,
		ParentID:   query.ParentID,
		Status:     query.Status,
		Tags:       query.Tags,
		From:       query.From,
		To:         query.To,
		Limit:      query.Pagination.Limit,
		Offset:     query.Pagination.Offset,
	}

	var assets []*media.Asset
	var total int64
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		var err error
		assets, total, err = repos.Assets.List(ctx, filter)
		return err
	})

	if err != nil {
		return nil, apperr.Internal("list assets", err)
	}

	return &dto.PagedResult[*media.Asset]{
		Items:  assets,
		Total:  total,
		Limit:  query.Pagination.Limit,
		Offset: query.Pagination.Offset,
	}, nil
}
