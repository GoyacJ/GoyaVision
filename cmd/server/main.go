package main

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goyavision"
	"goyavision/config"
	"goyavision/internal/adapter/engine"
	mcpadapter "goyavision/internal/adapter/mcp"
	"goyavision/internal/adapter/mediamtx"
	"goyavision/internal/adapter/persistence"
	"goyavision/internal/adapter/schema"
	"goyavision/internal/api"
	"goyavision/internal/app"
	infraauth "goyavision/internal/infra/auth"
	infraengine "goyavision/internal/infra/engine"
	inframediamtx "goyavision/internal/infra/mediamtx"
	infrapersistence "goyavision/internal/infra/persistence"
	"goyavision/internal/port"
	"goyavision/pkg/storage"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	var db *gorm.DB
	if cfg.DB.DSN != "" {
		db, err = gorm.Open(postgres.Open(cfg.DB.DSN), &gorm.Config{})
		if err != nil {
			log.Fatalf("open db: %v", err)
		}
		sqlDB, _ := db.DB()
		defer sqlDB.Close()
		if err := persistence.AutoMigrate(db); err != nil {
			log.Fatalf("migrate: %v", err)
		}
	} else {
		log.Print("db.dsn empty, skip database")
	}

	repo := persistence.NewRepository(db)
	uow := infrapersistence.NewUnitOfWork(db)
	mediaGateway := inframediamtx.NewGateway(
		cfg.MediaMTX.APIAddress,
		cfg.MediaMTX.Username,
		cfg.MediaMTX.Password,
		cfg.MediaMTX.RecordPath,
		cfg.MediaMTX.RecordFormat,
		cfg.MediaMTX.SegmentDuration,
	)
	tokenService, err := infraauth.NewJWTService(&cfg.JWT)
	if err != nil {
		log.Fatalf("create jwt service: %v", err)
	}
	schemaValidator := schema.NewJSONSchemaValidator()

	mtxCli := mediamtx.NewClient(cfg.MediaMTX.APIAddress, cfg.MediaMTX.Username, cfg.MediaMTX.Password)
	if err := mtxCli.Ping(context.Background()); err != nil {
		log.Printf("warning: mediamtx not available: %v", err)
	} else {
		log.Printf("mediamtx connected: %s", cfg.MediaMTX.APIAddress)
	}

	minioClient, err := storage.NewMinIOClient(
		cfg.MinIO.Endpoint,
		cfg.MinIO.AccessKey,
		cfg.MinIO.SecretKey,
		cfg.MinIO.BucketName,
		cfg.MinIO.UseSSL,
	)
	if err != nil {
		log.Fatalf("create minio client: %v", err)
	}
	log.Printf("minio connected: %s/%s", cfg.MinIO.Endpoint, cfg.MinIO.BucketName)

	mcpClient := mcpadapter.NewStaticClientWithoutDefaults()
	for i := range cfg.MCP.Servers {
		serverCfg := cfg.MCP.Servers[i]
		tools := make([]port.MCPTool, 0, len(serverCfg.Tools))
		for j := range serverCfg.Tools {
			toolCfg := serverCfg.Tools[j]
			tools = append(tools, port.MCPTool{
				Name:         toolCfg.Name,
				Description:  toolCfg.Description,
				Version:      toolCfg.Version,
				InputSchema:  toolCfg.InputSchema,
				OutputSchema: toolCfg.OutputSchema,
			})
		}
		mcpClient.RegisterServerWithConfig(port.MCPServer{
			ID:          serverCfg.ID,
			Name:        serverCfg.Name,
			Description: serverCfg.Description,
			Status:      serverCfg.Status,
		}, tools, serverCfg.Endpoint, serverCfg.APIToken, serverCfg.TimeoutSec)
	}

	var workflowScheduler *app.WorkflowScheduler
	if db != nil {
		ctx := context.Background()

		httpExecutor := engine.NewHTTPOperatorExecutor()
		cliExecutor := engine.NewCLIOperatorExecutor()
		mcpExecutor := engine.NewMCPOperatorExecutor(mcpClient)
		registry := engine.NewExecutorRegistry()
		registry.Register(httpExecutor.Mode(), httpExecutor)
		registry.Register(cliExecutor.Mode(), cliExecutor)
		registry.Register(mcpExecutor.Mode(), mcpExecutor)
		routingExecutor := engine.NewRoutingOperatorExecutor(registry)

		workflowEngine := infraengine.NewDAGWorkflowEngine(uow, routingExecutor, schemaValidator)
		workflowScheduler, err = app.NewWorkflowScheduler(repo, workflowEngine)
		if err != nil {
			log.Fatalf("create workflow scheduler: %v", err)
		}

		if err := workflowScheduler.Start(ctx); err != nil {
			log.Fatalf("start workflow scheduler: %v", err)
		}
		log.Print("workflow scheduler started (DAG engine)")
		defer workflowScheduler.Stop()
	}

	e := echo.New()

	var webDist fs.FS
	if sub, err := goyavision.GetWebFS(); err == nil {
		webDist = sub
	}

	handlers := api.NewHandlers(
		uow,
		schemaValidator,
		mcpClient,
		mcpClient,
		mediaGateway,
		tokenService,
		cfg,
		mtxCli,
		minioClient,
		workflowScheduler,
		repo,
	)
	api.RegisterRouter(e, handlers, webDist)

	srv := &http.Server{Addr: cfg.Server.Addr(), Handler: e}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
	log.Print("server stopped")
}
