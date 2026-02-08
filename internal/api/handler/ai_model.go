package handler

import (
	"net/http"

	"goyavision/internal/api/dto"
	authmiddleware "goyavision/internal/api/middleware"
	appdto "goyavision/internal/app/dto"
	"goyavision/internal/domain/ai_model"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterAIModelRoutes(public *echo.Group, protected *echo.Group, h *Handlers) {
	handler := &aiModelHandler{h: h}
	// Public
	public.GET("/ai-models", handler.List)
	public.GET("/ai-models/:id", handler.Get)

	// Protected
	protected.POST("/ai-models", handler.Create)
	protected.PUT("/ai-models/:id", handler.Update)
	protected.DELETE("/ai-models/:id", handler.Delete)
	protected.POST("/ai-models/:id/test-connection", handler.TestConnection)
}

type aiModelHandler struct {
	h *Handlers
}

func (h *aiModelHandler) List(c echo.Context) error {
	var query dto.AIModelListQuery
	if err := c.Bind(&query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters")
	}

	q := appdto.ListAIModelsQuery{
		Keyword: query.Keyword,
		Pagination: appdto.Pagination{
			Limit:  query.Limit,
			Offset: query.Offset,
		},
	}
	if query.Provider != "" {
		q.Provider = &query.Provider
	}
	if query.Status != "" {
		q.Status = &query.Status
	}

	result, err := h.h.ListAIModels.Handle(c.Request().Context(), q)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.AIModelListResponse{
		Items: dto.AIModelsToResponse(result.Items),
		Total: result.Total,
	})
}

func (h *aiModelHandler) Create(c echo.Context) error {
	var req dto.AIModelCreateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	visibility := ai_model.VisibilityPrivate
	if req.Visibility != nil {
		visibility = ai_model.Visibility(*req.Visibility)
	}

	visibleRoleIDs := req.VisibleRoleIDs
	if visibility == ai_model.VisibilityRole && len(visibleRoleIDs) == 0 {
		if ids, ok := c.Get(authmiddleware.ContextKeyRoleIDs).([]uuid.UUID); ok {
			for _, id := range ids {
				visibleRoleIDs = append(visibleRoleIDs, id.String())
			}
		}
	}

	cmd := appdto.CreateAIModelCommand{
		Name:           req.Name,
		Description:    req.Description,
		Provider:       req.Provider,
		Endpoint:       req.Endpoint,
		APIKey:         req.APIKey,
		ModelName:      req.ModelName,
		Config:         req.Config,
		Visibility:     visibility,
		VisibleRoleIDs: visibleRoleIDs,
	}

	model, err := h.h.CreateAIModel.Handle(c.Request().Context(), cmd)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, dto.AIModelToResponse(model))
}

func (h *aiModelHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid ai model id")
	}

	model, err := h.h.GetAIModel.Handle(c.Request().Context(), appdto.GetAIModelQuery{ID: id})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.AIModelToResponse(model))
}

func (h *aiModelHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid ai model id")
	}

	var req dto.AIModelUpdateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	cmd := appdto.UpdateAIModelCommand{
		ID:             id,
		Name:           req.Name,
		Description:    req.Description,
		Provider:       req.Provider,
		Endpoint:       req.Endpoint,
		APIKey:         req.APIKey,
		ModelName:      req.ModelName,
		Config:         req.Config,
		Status:         req.Status,
		VisibleRoleIDs: req.VisibleRoleIDs,
	}

	if req.Visibility != nil {
		v := ai_model.Visibility(*req.Visibility)
		cmd.Visibility = &v
	}

	model, err := h.h.UpdateAIModel.Handle(c.Request().Context(), cmd)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.AIModelToResponse(model))
}

func (h *aiModelHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid ai model id")
	}

	err = h.h.DeleteAIModel.Handle(c.Request().Context(), appdto.DeleteAIModelCommand{ID: id})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *aiModelHandler) TestConnection(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid ai model id")
	}

	result, err := h.h.TestAIModel.Handle(c.Request().Context(), appdto.TestAIModelCommand{ID: id})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.TestAIModelResponse{
		Success: result.Success,
		Message: result.Message,
	})
}
