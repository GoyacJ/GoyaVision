package handler

import (
	"goyavision/internal/api"

	"github.com/labstack/echo/v4"
)

func RegisterPreview(g *echo.Group, d api.Deps) {
	h := previewHandler{d: d}
	g.GET("/streams/:id/preview/start", h.Start)
	g.POST("/streams/:id/preview/stop", h.Stop)
}

type previewHandler struct{ d api.Deps }

func (h *previewHandler) Start(c echo.Context) error {
	return c.JSON(200, map[string]string{"hls_url": ""})
}

func (h *previewHandler) Stop(c echo.Context) error {
	return c.NoContent(204)
}
