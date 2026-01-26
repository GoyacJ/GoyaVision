package handler

import (
	"goyavision/internal/api"
	"goyavision/internal/api/dto"
	"goyavision/internal/app"
	"goyavision/pkg/ffmpeg"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterRecord(g *echo.Group, d api.Deps) {
	pool := ffmpeg.NewPool(d.Cfg.FFmpeg.Bin, d.Cfg.FFmpeg.MaxRecord, d.Cfg.FFmpeg.MaxFrame)
	manager := ffmpeg.NewManager(pool, d.Cfg.Record.BasePath)
	svc := app.NewRecordService(d.Repo, manager, d.Cfg.Record.BasePath, d.Cfg.Record.SegmentSec)
	h := recordHandler{
		d:   d,
		svc: svc,
	}
	g.POST("/streams/:id/record/start", h.Start)
	g.POST("/streams/:id/record/stop", h.Stop)
	g.GET("/streams/:id/record/sessions", h.ListSessions)
}

type recordHandler struct {
	d   api.Deps
	svc *app.RecordService
}

func (h *recordHandler) Start(c echo.Context) error {
	streamIDStr := c.Param("id")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		return c.JSON(400, api.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid stream id",
		})
	}

	session, err := h.svc.Start(c.Request().Context(), streamID)
	if err != nil {
		return err
	}

	return c.JSON(201, dto.RecordStartResponse{
		SessionID: session.ID,
	})
}

func (h *recordHandler) Stop(c echo.Context) error {
	streamIDStr := c.Param("id")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		return c.JSON(400, api.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid stream id",
		})
	}

	if err := h.svc.Stop(c.Request().Context(), streamID); err != nil {
		return err
	}

	return c.NoContent(204)
}

func (h *recordHandler) ListSessions(c echo.Context) error {
	streamIDStr := c.Param("id")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		return c.JSON(400, api.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid stream id",
		})
	}

	sessions, err := h.svc.ListSessions(c.Request().Context(), streamID)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.RecordSessionsToResponse(sessions))
}
