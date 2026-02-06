package handler

import (
	"net/http"
	"strings"

	"goyavision/internal/api/dto"
	appdto "goyavision/internal/app/dto"
	"goyavision/internal/domain/operator"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterOperator(g *echo.Group, h *Handlers) {
	handler := &operatorHandler{h: h}
	g.GET("/operators", handler.List)
	g.POST("/operators", handler.Create)
	g.GET("/operators/:id", handler.Get)
	g.PUT("/operators/:id", handler.Update)
	g.DELETE("/operators/:id", handler.Delete)
	g.POST("/operators/:id/enable", handler.Enable)
	g.POST("/operators/:id/disable", handler.Disable)
	g.GET("/operators/category/:category", handler.ListByCategory)
}

type operatorHandler struct {
	h *Handlers
}

func (h *operatorHandler) List(c echo.Context) error {
	var query dto.OperatorListQuery
	if err := c.Bind(&query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters")
	}

	q := appdto.ListOperatorsQuery{
		Pagination: appdto.Pagination{
			Limit:  query.Limit,
			Offset: query.Offset,
		},
	}

	if query.Category != nil {
		cat := operator.Category(*query.Category)
		q.Category = &cat
	}

	if query.Type != nil {
		t := operator.Type(*query.Type)
		q.Type = &t
	}

	if query.Status != nil {
		s := operator.Status(*query.Status)
		q.Status = &s
	}

	if query.IsBuiltin != nil {
		q.IsBuiltin = query.IsBuiltin
	}

	if query.Tags != nil && *query.Tags != "" {
		q.Tags = strings.Split(*query.Tags, ",")
	}

	if query.Keyword != nil {
		q.Keyword = *query.Keyword
	}

	result, err := h.h.ListOperators.Handle(c.Request().Context(), q)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.OperatorListResponse{
		Items: dto.OperatorsToResponse(result.Items),
		Total: result.Total,
	})
}

func (h *operatorHandler) Create(c echo.Context) error {
	var req dto.OperatorCreateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	status := operator.StatusDraft
	if req.Status != "" {
		status = operator.Status(req.Status)
	}

	cmd := appdto.CreateOperatorCommand{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Category:    operator.Category(req.Category),
		Type:        operator.Type(req.Type),
		Version:     req.Version,
		Endpoint:    req.Endpoint,
		Method:      req.Method,
		InputSchema: req.InputSchema,
		OutputSpec:  req.OutputSpec,
		Config:      req.Config,
		Status:      status,
		IsBuiltin:   req.IsBuiltin,
		Tags:        req.Tags,
	}

	op, err := h.h.CreateOperator.Handle(c.Request().Context(), cmd)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, dto.OperatorToResponse(op))
}

func (h *operatorHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid operator id")
	}

	op, err := h.h.GetOperator.Handle(c.Request().Context(), appdto.GetOperatorQuery{ID: id})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.OperatorToResponse(op))
}

func (h *operatorHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid operator id")
	}

	var req dto.OperatorUpdateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	cmd := appdto.UpdateOperatorCommand{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Endpoint:    req.Endpoint,
		Method:      req.Method,
		InputSchema: req.InputSchema,
		OutputSpec:  req.OutputSpec,
		Config:      req.Config,
		Tags:        req.Tags,
	}

	if req.Status != nil {
		s := operator.Status(*req.Status)
		cmd.Status = &s
	}

	op, err := h.h.UpdateOperator.Handle(c.Request().Context(), cmd)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.OperatorToResponse(op))
}

func (h *operatorHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid operator id")
	}

	err = h.h.DeleteOperator.Handle(c.Request().Context(), appdto.DeleteOperatorCommand{ID: id})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *operatorHandler) Enable(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid operator id")
	}

	op, err := h.h.EnableOperator.Handle(c.Request().Context(), appdto.EnableOperatorCommand{ID: id, Enabled: true})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.OperatorToResponse(op))
}

func (h *operatorHandler) Disable(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid operator id")
	}

	op, err := h.h.EnableOperator.Handle(c.Request().Context(), appdto.EnableOperatorCommand{ID: id, Enabled: false})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.OperatorToResponse(op))
}

func (h *operatorHandler) ListByCategory(c echo.Context) error {
	categoryStr := c.Param("category")
	category := operator.Category(categoryStr)

	q := appdto.ListOperatorsQuery{
		Category: &category,
		Pagination: appdto.Pagination{
			Limit:  100,
			Offset: 0,
		},
	}

	result, err := h.h.ListOperators.Handle(c.Request().Context(), q)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.OperatorsToResponse(result.Items))
}
