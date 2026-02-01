package handler

import (
	"goyavision/internal/api/dto"
	"goyavision/internal/app"
	"goyavision/internal/domain"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/datatypes"
)

func RegisterAlgorithmBinding(g *echo.Group, d Deps) {
	svc := app.NewAlgorithmBindingService(d.Repo)
	h := algorithmBindingHandler{svc: svc}
	g.GET("/streams/:id/algorithm-bindings", h.List)
	g.POST("/streams/:id/algorithm-bindings", h.Create)
	g.GET("/streams/:id/algorithm-bindings/:bid", h.Get)
	g.PUT("/streams/:id/algorithm-bindings/:bid", h.Update)
	g.DELETE("/streams/:id/algorithm-bindings/:bid", h.Delete)
}

type algorithmBindingHandler struct {
	svc *app.AlgorithmBindingService
}

func (h *algorithmBindingHandler) List(c echo.Context) error {
	streamIDStr := c.Param("id")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid stream id",
		})
	}

	bindings, err := h.svc.ListByStream(c.Request().Context(), streamID)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.AlgorithmBindingsToResponse(bindings))
}

func (h *algorithmBindingHandler) Create(c echo.Context) error {
	streamIDStr := c.Param("id")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid stream id",
		})
	}

	var req dto.AlgorithmBindingCreateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	binding := &domain.AlgorithmBinding{
		AlgorithmID:     req.AlgorithmID,
		IntervalSec:     req.IntervalSec,
		InitialDelaySec: req.InitialDelaySec,
		Enabled:         true,
	}
	if req.Enabled != nil {
		binding.Enabled = *req.Enabled
	}
	if req.Schedule != nil {
		binding.Schedule = datatypes.JSON(req.Schedule)
	}
	if req.Config != nil {
		binding.Config = datatypes.JSON(req.Config)
	}

	result, err := h.svc.Create(c.Request().Context(), streamID, binding)
	if err != nil {
		return err
	}

	return c.JSON(201, dto.AlgorithmBindingToResponse(result))
}

func (h *algorithmBindingHandler) Get(c echo.Context) error {
	bindingIDStr := c.Param("bid")
	bindingID, err := uuid.Parse(bindingIDStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid algorithm binding id",
		})
	}

	binding, err := h.svc.Get(c.Request().Context(), bindingID)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.AlgorithmBindingToResponse(binding))
}

func (h *algorithmBindingHandler) Update(c echo.Context) error {
	bindingIDStr := c.Param("bid")
	bindingID, err := uuid.Parse(bindingIDStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid algorithm binding id",
		})
	}

	var req dto.AlgorithmBindingUpdateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	binding := &domain.AlgorithmBinding{}
	if req.AlgorithmID != nil {
		binding.AlgorithmID = *req.AlgorithmID
	}
	if req.IntervalSec != nil {
		binding.IntervalSec = *req.IntervalSec
	}
	if req.InitialDelaySec != nil {
		binding.InitialDelaySec = *req.InitialDelaySec
	}
	if req.Schedule != nil {
		binding.Schedule = datatypes.JSON(req.Schedule)
	}
	if req.Config != nil {
		binding.Config = datatypes.JSON(req.Config)
	}
	if req.Enabled != nil {
		binding.Enabled = *req.Enabled
	}

	result, err := h.svc.Update(c.Request().Context(), bindingID, binding)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.AlgorithmBindingToResponse(result))
}

func (h *algorithmBindingHandler) Delete(c echo.Context) error {
	bindingIDStr := c.Param("bid")
	bindingID, err := uuid.Parse(bindingIDStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid algorithm binding id",
		})
	}

	if err := h.svc.Delete(c.Request().Context(), bindingID); err != nil {
		return err
	}

	return c.NoContent(204)
}
