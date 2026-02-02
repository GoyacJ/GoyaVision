package api

import (
	"io/fs"

	"goyavision/internal/api/handler"
	authMiddleware "goyavision/internal/api/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// HandlerDeps 处理器依赖
type HandlerDeps = handler.Deps

func RegisterRouter(e *echo.Echo, d HandlerDeps, webFS fs.FS) {
	e.HTTPErrorHandler = ErrorHandler
	e.Use(middleware.Logger(), middleware.Recover())

	authGroup := e.Group("/api/v1/auth")
	handler.RegisterAuth(authGroup, d)

	api := e.Group("/api/v1", authMiddleware.JWTAuth(d.Cfg.JWT))

	authProtected := api.Group("/auth")
	handler.RegisterAuthProtected(authProtected, d)

	api.Use(authMiddleware.LoadUserPermissions(d.Repo))

	handler.RegisterStream(api, d)
	handler.RegisterAsset(api, d)
	handler.RegisterOperator(api, d)
	handler.RegisterWorkflow(api, d)
	handler.RegisterTask(api, d)
	handler.RegisterArtifact(api, d)
	handler.RegisterAlgorithm(api, d)
	handler.RegisterAlgorithmBinding(api, d)
	handler.RegisterRecord(api, d)
	handler.RegisterInference(api, d)
	handler.RegisterPreview(api, d)
	handler.RegisterPlayback(api, d)

	admin := api.Group("")
	handler.RegisterUser(admin, d)
	handler.RegisterRole(admin, d)
	handler.RegisterMenu(admin, d)

	e.Static("/live", "./data/hls")

	RegisterStatic(e, webFS)
}
