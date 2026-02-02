package handler

import (
	"strings"

	"goyavision/internal/api/dto"
	"goyavision/internal/app"
	"goyavision/internal/domain"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterOperator(g *echo.Group, d Deps) {
	svc := app.NewOperatorService(d.Repo)
	h := operatorHandler{svc: svc}
	g.GET("/operators", h.List)
	g.POST("/operators", h.Create)
	g.GET("/operators/:id", h.Get)
	g.PUT("/operators/:id", h.Update)
	g.DELETE("/operators/:id", h.Delete)
	g.POST("/operators/:id/enable", h.Enable)
	g.POST("/operators/:id/disable", h.Disable)
	g.GET("/operators/category/:category", h.ListByCategory)
}

type operatorHandler struct {
	svc *app.OperatorService
}

func (h *operatorHandler) List(c echo.Context) error {
	var query dto.OperatorListQuery
	if err := c.Bind(&query); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid query parameters",
		})
	}

	req := &app.ListOperatorsRequest{
		Limit:  query.Limit,
		Offset: query.Offset,
	}

	if query.Category != nil {
		cat := domain.OperatorCategory(*query.Category)
		req.Category = &cat
	}

	if query.Type != nil {
		t := domain.OperatorType(*query.Type)
		req.Type = &t
	}

	if query.Status != nil {
		s := domain.OperatorStatus(*query.Status)
		req.Status = &s
	}

	if query.IsBuiltin != nil {
		req.IsBuiltin = query.IsBuiltin
	}

	if query.Tags != nil && *query.Tags != "" {
		req.Tags = strings.Split(*query.Tags, ",")
	}

	if query.Keyword != nil {
		req.Keyword = *query.Keyword
	}

	operators, total, err := h.svc.List(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.OperatorListResponse{
		Items: dto.OperatorsToResponse(operators),
		Total: total,
	})
}

func (h *operatorHandler) Create(c echo.Context) error {
	var req dto.OperatorCreateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	status := domain.OperatorStatusDraft
	if req.Status != "" {
		status = domain.OperatorStatus(req.Status)
	}

	createReq := &app.CreateOperatorRequest{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Category:    domain.OperatorCategory(req.Category),
		Type:        domain.OperatorType(req.Type),
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

	operator, err := h.svc.Create(c.Request().Context(), createReq)
	if err != nil {
		return err
	}

	return c.JSON(201, dto.OperatorToResponse(operator))
}

func (h *operatorHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid operator id",
		})
	}

	operator, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.OperatorToResponse(operator))
}

func (h *operatorHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid operator id",
		})
	}

	var req dto.OperatorUpdateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	updateReq := &app.UpdateOperatorRequest{
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
		s := domain.OperatorStatus(*req.Status)
		updateReq.Status = &s
	}

	operator, err := h.svc.Update(c.Request().Context(), id, updateReq)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.OperatorToResponse(operator))
}

func (h *operatorHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid operator id",
		})
	}

	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		return err
	}

	return c.NoContent(204)
}

func (h *operatorHandler) Enable(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid operator id",
		})
	}

	operator, err := h.svc.Enable(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.OperatorToResponse(operator))
}

func (h *operatorHandler) Disable(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid operator id",
		})
	}

	operator, err := h.svc.Disable(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.OperatorToResponse(operator))
}

func (h *operatorHandler) ListByCategory(c echo.Context) error {
	categoryStr := c.Param("category")
	category := domain.OperatorCategory(categoryStr)

	operators, err := h.svc.ListByCategory(c.Request().Context(), category)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.OperatorsToResponse(operators))
}
