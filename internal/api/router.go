package api

import (
	"io/fs"

	"goyavision/config"
	"goyavision/internal/api/handler"
	"goyavision/internal/port"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterRouter(e *echo.Echo, repo port.Repository, cfg *config.Config, webFS fs.FS) {
	e.HTTPErrorHandler = ErrorHandler
	e.Use(middleware.Logger(), middleware.Recover())

	g := e.Group("/api/v1")

	d := handler.Deps{Repo: repo, Cfg: cfg}
	handler.RegisterStream(g, d)
	handler.RegisterAlgorithm(g, d)
	handler.RegisterAlgorithmBinding(g, d)
	handler.RegisterRecord(g, d)
	handler.RegisterInference(g, d)
	handler.RegisterPreview(g, d)

	e.Static("/live", "./data/hls")

	RegisterStatic(e, webFS)
}
