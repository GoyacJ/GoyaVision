package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goyavision/config"
	"goyavision/internal/adapter/persistence"
	"goyavision/internal/api"

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
	e := echo.New()
	api.RegisterRouter(e, api.Deps{Repo: repo, Cfg: cfg})

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
