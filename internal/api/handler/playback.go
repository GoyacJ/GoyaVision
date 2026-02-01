package handler

import (
	"time"

	"goyavision/internal/api/dto"
	"goyavision/internal/app"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterPlayback(g *echo.Group, d Deps) {
	svc := app.NewPlaybackService(d.Repo, d.MtxCli, d.Cfg.MediaMTX)
	h := playbackHandler{svc: svc}
	g.GET("/streams/:id/playback", h.GetPlaybackURL)
	g.GET("/streams/:id/playback/segments", h.ListSegments)
}

type playbackHandler struct {
	svc *app.PlaybackService
}

func (h *playbackHandler) GetPlaybackURL(c echo.Context) error {
	streamIDStr := c.Param("id")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid stream id",
		})
	}

	startStr := c.QueryParam("start")
	if startStr == "" {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "start parameter is required",
		})
	}

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid start time format, use RFC3339",
		})
	}

	urls, err := h.svc.GetPlaybackURL(c.Request().Context(), streamID, start)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.PlaybackURLsResponse{
		HLS: urls.HLS,
		MP4: urls.MP4,
	})
}

func (h *playbackHandler) ListSegments(c echo.Context) error {
	streamIDStr := c.Param("id")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid stream id",
		})
	}

	segments, err := h.svc.ListRecordingSegments(c.Request().Context(), streamID)
	if err != nil {
		return err
	}

	result := make([]dto.PlaybackSegmentResponse, len(segments))
	for i, seg := range segments {
		result[i] = dto.PlaybackSegmentResponse{
			Start:       seg.Start,
			PlaybackURL: seg.PlaybackURL,
		}
	}

	return c.JSON(200, result)
}
