package handler

import (
	"goyavision/internal/api/dto"
	"goyavision/internal/app"
	"goyavision/internal/domain"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterStream(g *echo.Group, d Deps) {
	svc := app.NewStreamService(d.Repo, d.MtxCli, d.Cfg.MediaMTX)
	h := streamHandler{svc: svc}
	g.GET("/streams", h.List)
	g.POST("/streams", h.Create)
	g.GET("/streams/:id", h.Get)
	g.PUT("/streams/:id", h.Update)
	g.DELETE("/streams/:id", h.Delete)
	g.GET("/streams/:id/status", h.Status)
	g.POST("/streams/:id/enable", h.Enable)
	g.POST("/streams/:id/disable", h.Disable)
}

type streamHandler struct {
	svc *app.StreamService
}

func (h *streamHandler) List(c echo.Context) error {
	var query dto.StreamListQuery
	if err := c.Bind(&query); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid query parameters",
		})
	}

	if query.WithStatus {
		streams, err := h.svc.ListWithStatus(c.Request().Context(), query.Enabled)
		if err != nil {
			return err
		}
		return c.JSON(200, dto.StreamsWithStatusToResponse(streams))
	}

	streams, err := h.svc.List(c.Request().Context(), query.Enabled)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.StreamsToResponse(streams))
}

func (h *streamHandler) Create(c echo.Context) error {
	var req dto.StreamCreateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	streamType := domain.StreamTypePull
	if req.Type == "push" {
		streamType = domain.StreamTypePush
	}

	createReq := &app.CreateStreamRequest{
		URL:     req.URL,
		Name:    req.Name,
		Type:    streamType,
		Enabled: req.Enabled,
	}

	stream, err := h.svc.Create(c.Request().Context(), createReq)
	if err != nil {
		return err
	}

	return c.JSON(201, dto.StreamToResponse(stream))
}

func (h *streamHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid stream id",
		})
	}

	withStatus := c.QueryParam("with_status") == "true"
	if withStatus {
		stream, err := h.svc.GetWithStatus(c.Request().Context(), id)
		if err != nil {
			return err
		}
		return c.JSON(200, dto.StreamWithStatusToResponse(stream))
	}

	stream, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.StreamToResponse(stream))
}

func (h *streamHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid stream id",
		})
	}

	var req dto.StreamUpdateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	updateReq := &app.UpdateStreamRequest{
		URL:     req.URL,
		Name:    req.Name,
		Enabled: req.Enabled,
	}

	stream, err := h.svc.Update(c.Request().Context(), id, updateReq)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.StreamToResponse(stream))
}

func (h *streamHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid stream id",
		})
	}

	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		return err
	}

	return c.NoContent(204)
}

func (h *streamHandler) Status(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid stream id",
		})
	}

	status, err := h.svc.GetStatus(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.StreamStatusToResponse(status))
}

func (h *streamHandler) Enable(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid stream id",
		})
	}

	if err := h.svc.Enable(c.Request().Context(), id); err != nil {
		return err
	}

	return c.JSON(200, map[string]string{"message": "stream enabled"})
}

func (h *streamHandler) Disable(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid stream id",
		})
	}

	if err := h.svc.Disable(c.Request().Context(), id); err != nil {
		return err
	}

	return c.JSON(200, map[string]string{"message": "stream disabled"})
}
