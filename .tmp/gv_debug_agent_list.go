package main

import (
	"context"
	"fmt"

	"goyavision/config"
	persistence "goyavision/internal/adapter/persistence"
	"goyavision/internal/api/middleware"
	"goyavision/internal/domain/agent"
	"goyavision/internal/infra/persistence/repo"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	if cfg.DB.Driver == "" {
		cfg.DB.Driver = "postgres"
	}
	db, err := persistence.OpenDB(cfg.DB.Driver, cfg.DB.DSN)
	if err != nil {
		panic(err)
	}
	db = db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Info)})

	cols, err := db.Migrator().ColumnTypes("agent_sessions")
	if err != nil {
		fmt.Println("column types error:", err)
	} else {
		fmt.Println("agent_sessions columns:")
		for _, c := range cols {
			fmt.Println(" -", c.Name())
		}
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.ContextKeyTenantID, uuid.MustParse("00000000-0000-0000-0000-000000000001"))

	r := repo.NewAgentSessionRepo(db)
	_, _, err = r.List(ctx, agent.SessionFilter{Limit: 20, Offset: 0})
	if err != nil {
		fmt.Println("list error:", err)
	} else {
		fmt.Println("list ok")
	}
}
