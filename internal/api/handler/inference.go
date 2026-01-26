package handler

import (
	"goyavision/internal/api"

	"github.com/labstack/echo/v4"
)

func RegisterInference(g *echo.Group, d api.Deps) {
	h := inferenceHandler{d: d}
	g.GET("/inference_results", h.List)
}

type inferenceHandler struct{ d api.Deps }

func (h *inferenceHandler) List(c echo.Context) error {
	return c.JSON(200, map[string]any{"items": []any{}, "total": 0})
}
