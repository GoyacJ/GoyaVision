package handler

import (
	"goyavision/internal/api"

	"github.com/labstack/echo/v4"
)

func RegisterAlgorithmBinding(g *echo.Group, d api.Deps) {
	h := algorithmBindingHandler{d: d}
	g.GET("/streams/:id/algorithm-bindings", h.List)
	g.POST("/streams/:id/algorithm-bindings", h.Create)
	g.GET("/streams/:id/algorithm-bindings/:bid", h.Get)
	g.PUT("/streams/:id/algorithm-bindings/:bid", h.Update)
	g.DELETE("/streams/:id/algorithm-bindings/:bid", h.Delete)
}

type algorithmBindingHandler struct{ d api.Deps }

func (h *algorithmBindingHandler) List(c echo.Context) error {
	return c.JSON(200, []any{})
}

func (h *algorithmBindingHandler) Create(c echo.Context) error {
	return c.JSON(201, map[string]string{"id": ""})
}

func (h *algorithmBindingHandler) Get(c echo.Context) error {
	return c.JSON(200, map[string]any{})
}

func (h *algorithmBindingHandler) Update(c echo.Context) error {
	return c.JSON(200, map[string]any{})
}

func (h *algorithmBindingHandler) Delete(c echo.Context) error {
	return c.NoContent(204)
}
