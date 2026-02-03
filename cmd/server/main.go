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
	"goyavision/internal/adapter/mediamtx"
	"goyavision/internal/adapter/persistence"
	"goyavision/internal/api"
	"goyavision/internal/app"
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
		if err := persistence.InitializeData(db); err != nil {
			log.Printf("init data: %v", err)
		}
	} else {
		log.Print("db.dsn empty, skip database")
	}

	repo := persistence.NewRepository(db)

	mtxCli := mediamtx.NewClient(cfg.MediaMTX.APIAddress)
	if err := mtxCli.Ping(context.Background()); err != nil {
		log.Printf("warning: mediamtx not available: %v", err)
	} else {
		log.Printf("mediamtx connected: %s", cfg.MediaMTX.APIAddress)
	}
	mtxPathSync := mediamtx.NewPathSync(mtxCli)
	mediaSourceService := app.NewMediaSourceService(repo, mtxPathSync)

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

	var workflowScheduler *app.WorkflowScheduler
	if db != nil {
		ctx := context.Background()

		executor := engine.NewHTTPOperatorExecutor()
		workflowEngine := engine.NewSimpleWorkflowEngine(repo, executor)
		workflowScheduler, err = app.NewWorkflowScheduler(repo, workflowEngine)
		if err != nil {
			log.Fatalf("create workflow scheduler: %v", err)
		}

		if err := workflowScheduler.Start(ctx); err != nil {
			log.Fatalf("start workflow scheduler: %v", err)
		}
		log.Print("workflow scheduler started")
		defer workflowScheduler.Stop()
	}

	e := echo.New()

	var webDist fs.FS
	if sub, err := goyavision.GetWebFS(); err == nil {
		webDist = sub
	}

	deps := api.HandlerDeps{
		Repo:               repo,
		Cfg:                cfg,
		MtxCli:             mtxCli,
		MediaSourceService:  mediaSourceService,
		MinIOClient:        minioClient,
		WorkflowScheduler:  workflowScheduler,
	}
	api.RegisterRouter(e, deps, webDist)

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
