package handler

import (
	"goyavision/internal/api"

	"github.com/labstack/echo/v4"
)

func RegisterAlgorithm(g *echo.Group, d api.Deps) {
	h := algorithmHandler{d: d}
	g.GET("/algorithms", h.List)
	g.POST("/algorithms", h.Create)
	g.GET("/algorithms/:id", h.Get)
	g.PUT("/algorithms/:id", h.Update)
	g.DELETE("/algorithms/:id", h.Delete)
}

type algorithmHandler struct{ d api.Deps }

func (h *algorithmHandler) List(c echo.Context) error {
	return c.JSON(200, []any{})
}

func (h *algorithmHandler) Create(c echo.Context) error {
	return c.JSON(201, map[string]string{"id": ""})
}

func (h *algorithmHandler) Get(c echo.Context) error {
	return c.JSON(200, map[string]any{})
}

func (h *algorithmHandler) Update(c echo.Context) error {
	return c.JSON(200, map[string]any{})
}

func (h *algorithmHandler) Delete(c echo.Context) error {
	return c.NoContent(204)
}
