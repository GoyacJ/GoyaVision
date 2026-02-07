package query

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/domain/ai_model"
	"goyavision/internal/port"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type GetAIModelHandler struct {
	repo port.Repository
}

func NewGetAIModelHandler(repo port.Repository) *GetAIModelHandler {
	return &GetAIModelHandler{repo: repo}
}

func (h *GetAIModelHandler) Handle(ctx context.Context, query dto.GetAIModelQuery) (*ai_model.AIModel, error) {
	m, err := h.repo.GetAIModel(ctx, query.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("ai model", query.ID.String())
		}
		return nil, apperr.Wrap(err, apperr.CodeDBError, "failed to get ai model")
	}

	return m, nil
}
