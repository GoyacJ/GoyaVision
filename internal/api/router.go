package api

import (
	apiErrors "goyavision/internal/api/errors"
	"goyavision/internal/api/handler"

	"github.com/labstack/echo/v4"
)

func RegisterRouter(e *echo.Echo, d Deps) {
	e.HTTPErrorHandler = apiErrors.ErrorHandler
	e.Use(middleware.Logger, middleware.Recover)

	g := e.Group("/api/v1")

	handler.RegisterStream(g, d)
	handler.RegisterAlgorithm(g, d)
	handler.RegisterAlgorithmBinding(g, d)
	handler.RegisterRecord(g, d)
	handler.RegisterInference(g, d)
	handler.RegisterPreview(g, d)

	e.Static("/live", "./data/hls")

	RegisterStatic(e)
}
