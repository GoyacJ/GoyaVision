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

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"github.com/labstack/echo/v4/middleware"
)

func NewHandlers(
	uow port.UnitOfWork,
	schemaValidator port.SchemaValidator,
	mcpClient portrepo.MCPClient,
	mcpRegistry portrepo.MCPRegistry,
	mediaGateway port.MediaGateway,
	tokenService port.TokenService,
	db *gorm.DB,
	cfg *config.Config,
	mtxCli *mediamtx.Client,
	minioClient *storage.MinIOClient,
	workflowScheduler *app.WorkflowScheduler,
	repo portrepo.Repository,
) *handler.Handlers {
	return handler.NewHandlers(
		uow,
		schemaValidator,
		mcpClient,
		mcpRegistry,
		mediaGateway,
		tokenService,
		db,
		cfg,
		mtxCli,
		minioClient,
		workflowScheduler,
		repo,
	)
}

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func RegisterRouter(e *echo.Echo, h *handler.Handlers, webFS fs.FS) {
	e.HTTPErrorHandler = ErrorHandler
	e.Validator = &customValidator{validator: validator.New()}
	e.Use(middleware.Logger(), middleware.Recover())

	authGroup := e.Group("/api/v1/auth")
	handler.RegisterAuth(authGroup, h)

	api := e.Group("/api/v1", authMiddleware.JWTAuth(h.Cfg.JWT))
	optionalApi := e.Group("/api/v1", authMiddleware.OptionalJWTAuth(h.Cfg.JWT))

	authProtected := api.Group("/auth")
	handler.RegisterAuthProtected(authProtected, h)

	api.Use(authMiddleware.LoadUserPermissions(h.Repo))
	optionalApi.Use(authMiddleware.LoadUserPermissions(h.Repo))

	handler.RegisterSystemConfig(optionalApi, api, h)
	handler.RegisterAssetRoutes(optionalApi, api, h)
	handler.RegisterSourceRoutes(optionalApi, api, h)
	handler.RegisterUpload(api, h)
	handler.RegisterFile(api, h)
	handler.RegisterOperatorRoutes(optionalApi, api, h)
	handler.RegisterWorkflowRoutes(optionalApi, api, h)
	handler.RegisterTaskRoutes(optionalApi, api, h)
	handler.RegisterArtifact(api, h)
	handler.RegisterAIModelRoutes(optionalApi, api, h)
	handler.RegisterUserAssetRoutes(api, h)

	admin := api.Group("")
	handler.RegisterUser(admin, h)
	handler.RegisterRole(admin, h)
	handler.RegisterMenu(admin, h)
	handler.RegisterTenant(admin, h)

	RegisterStatic(e, webFS)
}
