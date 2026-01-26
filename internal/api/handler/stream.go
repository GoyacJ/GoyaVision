package handler

import (
	"goyavision/internal/api"
	"goyavision/internal/api/dto"

	"github.com/labstack/echo/v4"
)

func RegisterStream(g *echo.Group, d api.Deps) {
	h := streamHandler{d: d}
	g.GET("/streams", h.List)
	g.POST("/streams", h.Create)
	g.GET("/streams/:id", h.Get)
	g.PUT("/streams/:id", h.Update)
	g.DELETE("/streams/:id", h.Delete)
}

type streamHandler struct{ d api.Deps }

func (h *streamHandler) List(c echo.Context) error {
	_ = dto.StreamListQuery{}
	return c.JSON(200, []any{})
}

func (h *streamHandler) Create(c echo.Context) error {
	_ = dto.StreamCreateReq{}
	return c.JSON(201, map[string]string{"id": ""})
}

func (h *streamHandler) Get(c echo.Context) error {
	return c.JSON(200, map[string]any{})
}

func (h *streamHandler) Update(c echo.Context) error {
	_ = dto.StreamUpdateReq{}
	return c.JSON(200, map[string]any{})
}

func (h *streamHandler) Delete(c echo.Context) error {
	return c.NoContent(204)
}
