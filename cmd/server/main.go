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
	adaptercrypto "goyavision/internal/adapter/crypto"
	"goyavision/internal/adapter/engine"
	mcpadapter "goyavision/internal/adapter/mcp"
	"goyavision/internal/adapter/mediamtx"
	"goyavision/internal/adapter/persistence"
	adapterstorage "goyavision/internal/adapter/storage"
	"goyavision/internal/adapter/schema"
	"goyavision/internal/api"
	"goyavision/internal/app"
	infraeventbus "goyavision/internal/infra/eventbus"
	infraauth "goyavision/internal/infra/auth"
	infraengine "goyavision/internal/infra/engine"
	inframediamtx "goyavision/internal/infra/mediamtx"
	infrapersistence "goyavision/internal/infra/persistence"
	"goyavision/internal/port"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	var db *gorm.DB
	if cfg.DB.DSN != "" {
		driver := cfg.DB.Driver
		if driver == "" {
			driver = "postgres"
		}
		db, err = persistence.OpenDB(driver, cfg.DB.DSN)
		if err != nil {
			log.Fatalf("open db: %v", err)
		}
		sqlDB, _ := db.DB()
		defer sqlDB.Close()
		if err := persistence.AutoMigrate(db); err != nil {
			log.Fatalf("migrate: %v", err)
		}
		log.Printf("db connected: driver=%s", driver)
	} else {
		log.Print("db.dsn empty, skip database")
	}

	repo := persistence.NewRepository(db)
	uow := infrapersistence.NewUnitOfWork(db)
	eventBus := infraeventbus.NewLocalEventBus(100)
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

	fileStorage, storageURLConfig, err := adapterstorage.NewFileStorageFromConfig(cfg)
	if err != nil {
		log.Fatalf("create file storage: %v", err)
	}
	log.Printf("storage connected: type=%s", cfg.Storage.Type)

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

	encryptKey := cfg.EncryptKey
	if encryptKey == "" {
		encryptKey = cfg.JWT.Secret
	}
	cryptoService, _ := adaptercrypto.NewAESCryptoService(encryptKey)

	var workflowScheduler *app.WorkflowScheduler
	if db != nil {
		ctx := context.Background()

		httpExecutor := engine.NewHTTPOperatorExecutor()
		cliExecutor := engine.NewCLIOperatorExecutor()
		mcpExecutor := engine.NewMCPOperatorExecutor(mcpClient)
		aiModelExecutor := engine.NewAIModelExecutor(repo, cryptoService)
		registry := engine.NewExecutorRegistry()
		registry.Register(httpExecutor.Mode(), httpExecutor)
		registry.Register(cliExecutor.Mode(), cliExecutor)
		registry.Register(mcpExecutor.Mode(), mcpExecutor)
		registry.Register(aiModelExecutor.Mode(), aiModelExecutor)
		routingExecutor := engine.NewRoutingOperatorExecutor(registry)

		workflowEngine := infraengine.NewDAGWorkflowEngine(uow, routingExecutor, schemaValidator)
		workflowScheduler, err = app.NewWorkflowScheduler(repo, workflowEngine, eventBus)
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
		db,
		cfg,
		mtxCli,
		fileStorage,
		storageURLConfig,
		workflowScheduler,
		repo,
		eventBus,
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
