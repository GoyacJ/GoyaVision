package api

import (
	"io/fs"

	"goyavision/config"
	"goyavision/internal/adapter/mediamtx"
	"goyavision/internal/api/handler"
	authMiddleware "goyavision/internal/api/middleware"
	"goyavision/internal/app"
	"goyavision/internal/app/port"
	portrepo "goyavision/internal/port"
	"goyavision/pkg/storage"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewHandlers(
	uow port.UnitOfWork,
	mediaGateway port.MediaGateway,
	tokenService port.TokenService,
	cfg *config.Config,
	mtxCli *mediamtx.Client,
	minioClient *storage.MinIOClient,
	workflowScheduler *app.WorkflowScheduler,
	repo portrepo.Repository,
) *handler.Handlers {
	return handler.NewHandlers(
		uow,
		mediaGateway,
		tokenService,
		cfg,
		mtxCli,
		minioClient,
		workflowScheduler,
		repo,
	)
}

func RegisterRouter(e *echo.Echo, h *handler.Handlers, webFS fs.FS) {
	e.HTTPErrorHandler = ErrorHandler
	e.Use(middleware.Logger(), middleware.Recover())

	authGroup := e.Group("/api/v1/auth")
	handler.RegisterAuth(authGroup, h)

	api := e.Group("/api/v1", authMiddleware.JWTAuth(h.Cfg.JWT))

	authProtected := api.Group("/auth")
	handler.RegisterAuthProtected(authProtected, h)

	api.Use(authMiddleware.LoadUserPermissions(h.Repo))

	handler.RegisterAsset(api, h)
	handler.RegisterSource(api, h)
	handler.RegisterUpload(api, h)
	handler.RegisterFile(api, h)
	handler.RegisterOperator(api, h)
	handler.RegisterWorkflow(api, h)
	handler.RegisterTask(api, h)
	handler.RegisterArtifact(api, h)

	admin := api.Group("")
	handler.RegisterUser(admin, h)
	handler.RegisterRole(admin, h)
	handler.RegisterMenu(admin, h)

	RegisterStatic(e, webFS)
}
