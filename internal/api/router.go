package api

import (
	"io/fs"

	"goyavision/config"
	"goyavision/internal/api/handler"
	authMiddleware "goyavision/internal/api/middleware"
	"goyavision/internal/port"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterRouter(e *echo.Echo, repo port.Repository, cfg *config.Config, webFS fs.FS) {
	e.HTTPErrorHandler = ErrorHandler
	e.Use(middleware.Logger(), middleware.Recover())

	d := handler.Deps{Repo: repo, Cfg: cfg}

	authGroup := e.Group("/api/v1/auth")
	handler.RegisterAuth(authGroup, d)

	api := e.Group("/api/v1", authMiddleware.JWTAuth(cfg.JWT))

	authProtected := api.Group("/auth")
	handler.RegisterAuthProtected(authProtected, d)

	api.Use(authMiddleware.LoadUserPermissions(repo))

	handler.RegisterStream(api, d)
	handler.RegisterAlgorithm(api, d)
	handler.RegisterAlgorithmBinding(api, d)
	handler.RegisterRecord(api, d)
	handler.RegisterInference(api, d)
	handler.RegisterPreview(api, d)

	admin := api.Group("")
	handler.RegisterUser(admin, d)
	handler.RegisterRole(admin, d)
	handler.RegisterMenu(admin, d)

	e.Static("/live", "./data/hls")

	RegisterStatic(e, webFS)
}
