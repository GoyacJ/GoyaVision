package handler

import (
	"goyavision/internal/api"

	"github.com/labstack/echo/v4"
)

func RegisterRecord(g *echo.Group, d api.Deps) {
	h := recordHandler{d: d}
	g.POST("/streams/:id/record/start", h.Start)
	g.POST("/streams/:id/record/stop", h.Stop)
	g.GET("/streams/:id/record/sessions", h.ListSessions)
}

type recordHandler struct{ d api.Deps }

func (h *recordHandler) Start(c echo.Context) error {
	return c.JSON(201, map[string]string{"session_id": ""})
}

func (h *recordHandler) Stop(c echo.Context) error {
	return c.NoContent(204)
}

func (h *recordHandler) ListSessions(c echo.Context) error {
	return c.JSON(200, []any{})
}
