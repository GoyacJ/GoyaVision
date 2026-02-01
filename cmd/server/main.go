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
	"goyavision/internal/adapter/ai"
	"goyavision/internal/adapter/persistence"
	"goyavision/internal/api"
	"goyavision/internal/app"
	"goyavision/pkg/ffmpeg"

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

	var scheduler *app.Scheduler
	if db != nil {
		inferenceAdapter := ai.NewInferenceAdapter(cfg.AI.Timeout, cfg.AI.Retry)
		pool := ffmpeg.NewPool(cfg.FFmpeg.Bin, cfg.FFmpeg.MaxRecord, cfg.FFmpeg.MaxFrame)
		manager := ffmpeg.NewManager(pool, cfg.Record.BasePath)
		frameBasePath := "./data/frames"
		scheduler, err = app.NewScheduler(repo, inferenceAdapter, manager, frameBasePath)
		if err != nil {
			log.Fatalf("create scheduler: %v", err)
		}

		ctx := context.Background()
		if err := scheduler.Start(ctx); err != nil {
			log.Fatalf("start scheduler: %v", err)
		}
		log.Print("scheduler started")
		defer scheduler.Stop()
	}

	e := echo.New()

	var webDist fs.FS
	if sub, err := goyavision.GetWebFS(); err == nil {
		webDist = sub
	}
	api.RegisterRouter(e, repo, cfg, webDist)

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
