package handler

import (
	"goyavision/internal/api/dto"
	"goyavision/internal/app"
	"goyavision/pkg/ffmpeg"
	"goyavision/pkg/preview"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterPreview(g *echo.Group, d Deps) {
	ffmpegPool := ffmpeg.NewPool(d.Cfg.FFmpeg.Bin, d.Cfg.FFmpeg.MaxRecord, d.Cfg.FFmpeg.MaxFrame)
	previewPool := preview.NewPool(d.Cfg.Preview.MaxPreview)
	manager := preview.NewManager(
		d.Cfg.Preview.Provider,
		d.Cfg.Preview.MediamtxBin,
		ffmpegPool,
		previewPool,
		d.Cfg.Preview.HLSBase,
		"./data/hls",
	)
	svc := app.NewPreviewService(d.Repo, manager)
	h := previewHandler{svc: svc}
	g.GET("/streams/:id/preview/start", h.Start)
	g.POST("/streams/:id/preview/stop", h.Stop)
}

type previewHandler struct {
	svc *app.PreviewService
}

func (h *previewHandler) Start(c echo.Context) error {
	streamIDStr := c.Param("id")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid stream id",
		})
	}

	hlsURL, err := h.svc.Start(c.Request().Context(), streamID)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.PreviewStartResponse{
		HLSURL: hlsURL,
	})
}

func (h *previewHandler) Stop(c echo.Context) error {
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
