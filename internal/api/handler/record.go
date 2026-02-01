package handler

import (
	"goyavision/internal/api/dto"
	"goyavision/internal/app"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterRecord(g *echo.Group, d Deps) {
	svc := app.NewRecordService(d.Repo, d.MtxCli, d.Cfg.MediaMTX)
	h := recordHandler{svc: svc}
	g.POST("/streams/:id/record/start", h.Start)
	g.POST("/streams/:id/record/stop", h.Stop)
	g.GET("/streams/:id/record/sessions", h.ListSessions)
	g.GET("/streams/:id/record/files", h.GetRecordings)
	g.GET("/streams/:id/record/status", h.Status)
}

type recordHandler struct {
	svc *app.RecordService
}

func (h *recordHandler) Start(c echo.Context) error {
	streamIDStr := c.Param("id")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
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
		return c.JSON(400, dto.ErrorResponse{
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
		return c.JSON(400, dto.ErrorResponse{
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

func (h *recordHandler) GetRecordings(c echo.Context) error {
	streamIDStr := c.Param("id")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid stream id",
		})
	}

	recordings, err := h.svc.GetRecordings(c.Request().Context(), streamID)
	if err != nil {
		return err
	}

	if recordings == nil {
		return c.JSON(200, nil)
	}

	segments := make([]dto.RecordingSegmentResponse, len(recordings.Segments))
	for i, seg := range recordings.Segments {
		segments[i] = dto.RecordingSegmentResponse{Start: seg.Start}
	}

	return c.JSON(200, dto.RecordingsResponse{
		Name:     recordings.Name,
		Segments: segments,
	})
}

func (h *recordHandler) Status(c echo.Context) error {
	streamIDStr := c.Param("id")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid stream id",
		})
	}

	isRecording, err := h.svc.IsRecording(c.Request().Context(), streamID)
	if err != nil {
		return err
	}

	return c.JSON(200, map[string]bool{"recording": isRecording})
}
