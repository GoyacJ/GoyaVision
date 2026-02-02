package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"goyavision/config"
	"goyavision/internal/adapter/persistence"
	"goyavision/internal/domain"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dryRun = flag.Bool("dry-run", false, "只显示将要执行的操作，不实际执行")
)

func main() {
	flag.Parse()

	log.Println("GoyaVision 数据迁移工具 v1.0")
	log.Println("================================")

	if *dryRun {
		log.Println("⚠️  模拟运行模式（不会修改数据库）")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	if cfg.DB.DSN == "" {
		log.Fatal("数据库 DSN 未配置")
	}

	db, err := gorm.Open(postgres.Open(cfg.DB.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	repo := persistence.NewRepository(db)
	ctx := context.Background()

	log.Println("\n📊 数据迁移计划:")
	log.Println("1. 迁移 streams → media_assets（作为媒体源）")
	log.Println("2. 迁移 algorithms → operators")
	log.Println("3. 清理废弃表（algorithm_bindings、inference_results）")

	if !confirm("\n是否继续？") && !*dryRun {
		log.Println("已取消")
		return
	}

	log.Println("\n开始迁移...")

	if err := migrateStreamsToAssets(ctx, db, repo); err != nil {
		log.Fatalf("迁移 streams 失败: %v", err)
	}

	if err := migrateAlgorithmsToOperators(ctx, db, repo); err != nil {
		log.Fatalf("迁移 algorithms 失败: %v", err)
	}

	if !*dryRun {
		if err := cleanupOldTables(db); err != nil {
			log.Fatalf("清理旧表失败: %v", err)
		}
	}

	log.Println("\n✅ 迁移完成！")
}

func migrateStreamsToAssets(ctx context.Context, db *gorm.DB, repo *persistence.Repository) error {
	log.Println("\n[1/3] 迁移 Streams → MediaAssets")

	var streams []domain.Stream
	if err := db.Find(&streams).Error; err != nil {
		return err
	}

	log.Printf("找到 %d 个流", len(streams))

	if *dryRun {
		log.Println("  （模拟运行，跳过实际迁移）")
		return nil
	}

	migrated := 0
	for _, stream := range streams {
		assetType := domain.AssetTypeVideo

		asset := &domain.MediaAsset{
			ID:         uuid.New(),
			Type:       assetType,
			SourceType: domain.AssetSourceTypeStreamCapture,
			SourceID:   &stream.ID,
			Name:       stream.Name,
			Path:       stream.URL,
			Format:     "rtsp",
			Status:     domain.AssetStatusReady,
		}

		if stream.Enabled {
			asset.Status = domain.AssetStatusReady
		} else {
			asset.Status = domain.AssetStatusPending
		}

		if err := repo.CreateMediaAsset(ctx, asset); err != nil {
			log.Printf("  ⚠️  跳过流 %s: %v", stream.Name, err)
			continue
		}

		migrated++
		log.Printf("  ✓ 迁移流: %s → 资产 ID: %s", stream.Name, asset.ID)
	}

	log.Printf("✅ 成功迁移 %d/%d 个流", migrated, len(streams))
	return nil
}

func migrateAlgorithmsToOperators(ctx context.Context, db *gorm.DB, repo *persistence.Repository) error {
	log.Println("\n[2/3] 迁移 Algorithms → Operators")

	var algorithms []domain.Algorithm
	if err := db.Find(&algorithms).Error; err != nil {
		return err
	}

	log.Printf("找到 %d 个算法", len(algorithms))

	if *dryRun {
		log.Println("  （模拟运行，跳过实际迁移）")
		return nil
	}

	migrated := 0
	for _, alg := range algorithms {
		category := domain.OperatorCategoryAnalysis
		opType := domain.OperatorType("legacy_" + alg.Type)

		operator := &domain.Operator{
			ID:          uuid.New(),
			Code:        alg.Code,
			Name:        alg.Name,
			Description: alg.Description,
			Category:    category,
			Type:        opType,
			Version:     "1.0.0",
			Endpoint:    alg.Endpoint,
			Method:      "POST",
			InputSchema: alg.InputSpec,
			OutputSpec:  alg.OutputSpec,
			Config:      alg.Config,
			Status:      domain.OperatorStatusPublished,
			IsBuiltin:   false,
		}

		if err := repo.CreateOperator(ctx, operator); err != nil {
			log.Printf("  ⚠️  跳过算法 %s: %v", alg.Name, err)
			continue
		}

		migrated++
		log.Printf("  ✓ 迁移算法: %s → 算子 ID: %s", alg.Name, operator.ID)
	}

	log.Printf("✅ 成功迁移 %d/%d 个算法", migrated, len(algorithms))
	return nil
}

func cleanupOldTables(db *gorm.DB) error {
	log.Println("\n[3/3] 清理废弃表")

	tables := []string{
		"algorithm_bindings",
		"inference_results",
	}

	for _, table := range tables {
		log.Printf("  删除表: %s", table)
		if err := db.Migrator().DropTable(table); err != nil {
			log.Printf("  ⚠️  删除表失败: %v", err)
			continue
		}
		log.Printf("  ✓ 已删除: %s", table)
	}

	log.Println("✅ 清理完成")
	return nil
}

func confirm(msg string) bool {
	if *dryRun {
		return true
	}

	log.Print(msg + " [y/N]: ")
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false
	}
	return response == "y" || response == "Y"
}
