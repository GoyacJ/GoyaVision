package handler

import (
	"goyavision/internal/api"
	"goyavision/internal/api/dto"
	"goyavision/internal/app"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterInference(g *echo.Group, d api.Deps) {
	svc := app.NewInferenceService(d.Repo)
	h := inferenceHandler{
		d:   d,
		svc: svc,
	}
	g.GET("/inference_results", h.List)
}

type inferenceHandler struct {
	d   api.Deps
	svc *app.InferenceService
}

func (h *inferenceHandler) List(c echo.Context) error {
	var query dto.InferenceResultListQuery
	if err := c.Bind(&query); err != nil {
		return c.JSON(400, api.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid query parameters",
		})
	}

	var streamID, bindingID *uuid.UUID
	if query.StreamID != nil {
		id, err := uuid.Parse(*query.StreamID)
		if err != nil {
			return c.JSON(400, api.ErrorResponse{
				Error:   "Bad Request",
				Message: "invalid stream_id",
			})
		}
		streamID = &id
	}
	if query.BindingID != nil {
		id, err := uuid.Parse(*query.BindingID)
		if err != nil {
			return c.JSON(400, api.ErrorResponse{
				Error:   "Bad Request",
				Message: "invalid binding_id",
			})
		}
		bindingID = &id
	}

	limit := query.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := query.Offset
	if offset < 0 {
		offset = 0
	}

	results, total, err := h.svc.ListResults(c.Request().Context(), streamID, bindingID, query.From, query.To, limit, offset)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.InferenceResultListResponse{
		Items: dto.InferenceResultsToResponse(results),
		Total: total,
	})
}
