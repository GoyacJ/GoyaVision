package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/domain/ai_model"
	"goyavision/internal/port"
	"goyavision/pkg/apperr"
)

type ListAIModelsHandler struct {
	repo port.Repository
}

func NewListAIModelsHandler(repo port.Repository) *ListAIModelsHandler {
	return &ListAIModelsHandler{repo: repo}
}

func (h *ListAIModelsHandler) Handle(ctx context.Context, query dto.ListAIModelsQuery) (*dto.PagedResult[*ai_model.AIModel], error) {
	filter := ai_model.Filter{
		Keyword: query.Keyword,
		Limit:   query.Limit,
		Offset:  query.Offset,
	}
	if query.Provider != nil {
		p := ai_model.Provider(*query.Provider)
		filter.Provider = &p
	}
	if query.Status != nil {
		s := ai_model.Status(*query.Status)
		filter.Status = &s
	}

	items, total, err := h.repo.ListAIModels(ctx, filter)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeDBError, "failed to list ai models")
	}

	return &dto.PagedResult[*ai_model.AIModel]{
		Items:  items,
		Total:  total,
		Limit:  query.Limit,
		Offset: query.Offset,
	}, nil
}
