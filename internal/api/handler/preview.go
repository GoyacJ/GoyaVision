package handler

import (
	"goyavision/internal/api/dto"
	"goyavision/internal/app"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterPreview(g *echo.Group, d Deps) {
	svc := app.NewPreviewService(d.Repo, d.MtxCli, d.Cfg.MediaMTX)
	h := previewHandler{svc: svc}
	g.GET("/streams/:id/preview", h.GetURLs)
	g.GET("/streams/:id/preview/start", h.Start)
	g.GET("/streams/:id/preview/ready", h.Ready)
}

type previewHandler struct {
	svc *app.PreviewService
}

func (h *previewHandler) GetURLs(c echo.Context) error {
	streamIDStr := c.Param("id")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid stream id",
		})
	}

	urls, err := h.svc.GetPreviewURLs(c.Request().Context(), streamID)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.PreviewStartResponse{
		HLSURL:    urls.HLS,
		RTSPUrl:   urls.RTSP,
		RTMPUrl:   urls.RTMP,
		WebRTCUrl: urls.WebRTC,
	})
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

	urls, err := h.svc.Start(c.Request().Context(), streamID)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.PreviewStartResponse{
		HLSURL:    urls.HLS,
		RTSPUrl:   urls.RTSP,
		RTMPUrl:   urls.RTMP,
		WebRTCUrl: urls.WebRTC,
	})
}

func (h *previewHandler) Ready(c echo.Context) error {
	streamIDStr := c.Param("id")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid stream id",
		})
	}

	ready, err := h.svc.IsStreamReady(c.Request().Context(), streamID)
	if err != nil {
		return err
	}

	return c.JSON(200, map[string]bool{"ready": ready})
}
