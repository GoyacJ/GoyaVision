package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"goyavision/config"
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

	ctx := context.Background()

	log.Println("\n📊 数据迁移计划:")
	log.Println("1. 更新菜单和权限（V1.0 新功能）")
	log.Println("2. 迁移 streams → media_assets（作为媒体源）")
	log.Println("3. 迁移 algorithms → operators")
	log.Println("4. 清理废弃表（algorithm_bindings、inference_results）")

	if !confirm("\n是否继续？") && !*dryRun {
		log.Println("已取消")
		return
	}

	log.Println("\n开始迁移...")

	if err := updateMenusAndPermissions(ctx, db); err != nil {
		log.Fatalf("更新菜单和权限失败: %v", err)
	}

	if err := migrateStreamsToAssets(ctx, db); err != nil {
		log.Fatalf("迁移 streams 失败: %v", err)
	}

	if err := migrateAlgorithmsToOperators(ctx, db); err != nil {
		log.Fatalf("迁移 algorithms 失败: %v", err)
	}

	if !*dryRun {
		if err := cleanupOldTables(db); err != nil {
			log.Fatalf("清理旧表失败: %v", err)
		}
	}

	log.Println("\n✅ 迁移完成！")
}

// LegacyAlgorithm 旧算法结构（用于迁移）
type LegacyAlgorithm struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Code        string
	Name        string
	Description string
	Type        string
	Endpoint    string
	InputSpec   []byte `gorm:"type:jsonb;column:input_spec"`
	OutputSpec  []byte `gorm:"type:jsonb;column:output_spec"`
	Config      []byte `gorm:"type:jsonb"`
}

func (LegacyAlgorithm) TableName() string { return "algorithms" }

func updateMenusAndPermissions(ctx context.Context, db *gorm.DB) error {
	log.Println("\n[1/4] 更新菜单和权限")

	if *dryRun {
		log.Println("  （模拟运行，跳过实际更新）")
		return nil
	}

	log.Println("  清理旧菜单...")
	oldMenuCodes := []string{"stream", "algorithm", "inference"}
	for _, code := range oldMenuCodes {
		if err := db.Where("code = ?", code).Delete(&domain.Menu{}).Error; err != nil {
			log.Printf("  ⚠️  删除旧菜单 %s 失败: %v", code, err)
		} else {
			log.Printf("  ✓ 删除旧菜单: %s", code)
		}
	}

	log.Println("  添加新菜单...")
	newMenus := []domain.Menu{
		{
			ID:         uuid.MustParse("00000000-0000-0000-0000-000000000010"),
			Code:       "asset",
			Name:       "媒体资产",
			Type:       2,
			Path:       "/assets",
			Icon:       "Files",
			Component:  "asset/index",
			Permission: "asset:list",
			Sort:       1,
			Visible:    true,
			Status:     1,
		},
		{
			ID:         uuid.MustParse("00000000-0000-0000-0000-000000000020"),
			Code:       "operator",
			Name:       "算子管理",
			Type:       2,
			Path:       "/operators",
			Icon:       "Cpu",
			Component:  "operator/index",
			Permission: "operator:list",
			Sort:       2,
			Visible:    true,
			Status:     1,
		},
		{
			ID:         uuid.MustParse("00000000-0000-0000-0000-000000000030"),
			Code:       "workflow",
			Name:       "工作流",
			Type:       2,
			Path:       "/workflows",
			Icon:       "Connection",
			Component:  "workflow/index",
			Permission: "workflow:list",
			Sort:       3,
			Visible:    true,
			Status:     1,
		},
		{
			ID:         uuid.MustParse("00000000-0000-0000-0000-000000000040"),
			Code:       "task",
			Name:       "任务管理",
			Type:       2,
			Path:       "/tasks",
			Icon:       "List",
			Component:  "task/index",
			Permission: "task:list",
			Sort:       4,
			Visible:    true,
			Status:     1,
		},
		{
			ID:         uuid.MustParse("00000000-0000-0000-0000-000000000050"),
			Code:       "legacy",
			Name:       "旧功能",
			Type:       1,
			Path:       "/legacy",
			Icon:       "FolderOpened",
			Component:  "",
			Permission: "",
			Sort:       90,
			Visible:    true,
			Status:     1,
		},
		{
			ID:       uuid.MustParse("00000000-0000-0000-0000-000000000051"),
			ParentID: ptrUUID("00000000-0000-0000-0000-000000000050"),
			Code:     "legacy:stream",
			Name:     "视频流（旧）",
			Type:     2,
			Path:     "/streams",
			Icon:     "Monitor",
			Component: "stream/index",
			Permission: "stream:list",
			Sort:       1,
			Visible:    true,
			Status:     1,
		},
	}

	for _, menu := range newMenus {
		var existing domain.Menu
		err := db.Where("code = ?", menu.Code).First(&existing).Error
		if err == nil {
			log.Printf("  ⊙ 菜单已存在，跳过: %s", menu.Name)
			continue
		}

		if err := db.Create(&menu).Error; err != nil {
			log.Printf("  ⚠️  创建菜单失败 %s: %v", menu.Name, err)
		} else {
			log.Printf("  ✓ 创建新菜单: %s", menu.Name)
		}
	}

	log.Println("  添加新权限...")
	newPermissions := []struct {
		Code   string
		Name   string
		Method string
		Path   string
	}{
		{"asset:list", "查看媒体资产列表", "GET", "/api/v1/assets"},
		{"asset:create", "创建媒体资产", "POST", "/api/v1/assets"},
		{"asset:update", "更新媒体资产", "PUT", "/api/v1/assets/*"},
		{"asset:delete", "删除媒体资产", "DELETE", "/api/v1/assets/*"},
		{"operator:list", "查看算子列表", "GET", "/api/v1/operators"},
		{"operator:create", "创建算子", "POST", "/api/v1/operators"},
		{"operator:update", "更新算子", "PUT", "/api/v1/operators/*"},
		{"operator:delete", "删除算子", "DELETE", "/api/v1/operators/*"},
		{"operator:enable", "启用算子", "PUT", "/api/v1/operators/*/enable"},
		{"operator:disable", "禁用算子", "PUT", "/api/v1/operators/*/disable"},
		{"workflow:list", "查看工作流列表", "GET", "/api/v1/workflows"},
		{"workflow:create", "创建工作流", "POST", "/api/v1/workflows"},
		{"workflow:update", "更新工作流", "PUT", "/api/v1/workflows/*"},
		{"workflow:delete", "删除工作流", "DELETE", "/api/v1/workflows/*"},
		{"workflow:enable", "启用工作流", "PUT", "/api/v1/workflows/*/enable"},
		{"workflow:disable", "禁用工作流", "PUT", "/api/v1/workflows/*/disable"},
		{"workflow:trigger", "触发工作流", "POST", "/api/v1/workflows/*/trigger"},
		{"task:list", "查看任务列表", "GET", "/api/v1/tasks"},
		{"task:create", "创建任务", "POST", "/api/v1/tasks"},
		{"task:update", "更新任务", "PUT", "/api/v1/tasks/*"},
		{"task:delete", "删除任务", "DELETE", "/api/v1/tasks/*"},
		{"task:cancel", "取消任务", "POST", "/api/v1/tasks/*/cancel"},
		{"artifact:list", "查看产物列表", "GET", "/api/v1/artifacts"},
		{"artifact:delete", "删除产物", "DELETE", "/api/v1/artifacts/*"},
	}

	addedPerms := 0
	for _, p := range newPermissions {
		var existing domain.Permission
		err := db.Where("code = ?", p.Code).First(&existing).Error
		if err == nil {
			continue
		}

		perm := &domain.Permission{
			ID:     uuid.New(),
			Code:   p.Code,
			Name:   p.Name,
			Method: p.Method,
			Path:   p.Path,
		}
		if err := db.Create(perm).Error; err != nil {
			log.Printf("  ⚠️  创建权限失败 %s: %v", p.Code, err)
		} else {
			addedPerms++
		}
	}
	log.Printf("  ✓ 新增权限: %d 个", addedPerms)

	log.Println("  更新超级管理员角色权限...")
	var superAdminRole domain.Role
	if err := db.Where("code = ?", "super_admin").First(&superAdminRole).Error; err == nil {
		db.Exec("DELETE FROM role_permissions WHERE role_id = ?", superAdminRole.ID)

		var allPermissions []domain.Permission
		db.Find(&allPermissions)
		for _, perm := range allPermissions {
			db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
				superAdminRole.ID, perm.ID)
		}

		db.Exec("DELETE FROM role_menus WHERE role_id = ?", superAdminRole.ID)
		var allMenus []domain.Menu
		db.Find(&allMenus)
		for _, menu := range allMenus {
			db.Exec("INSERT INTO role_menus (role_id, menu_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
				superAdminRole.ID, menu.ID)
		}
		log.Println("  ✓ 已更新超级管理员权限")
	}

	log.Println("✅ 菜单和权限更新完成")
	return nil
}

func ptrUUID(s string) *uuid.UUID {
	id := uuid.MustParse(s)
	return &id
}

func migrateStreamsToAssets(ctx context.Context, db *gorm.DB) error {
	log.Println("\n[2/4] 迁移 Streams → MediaAssets")

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
			SourceType: domain.AssetSourceLive,
			SourceID:   &stream.ID,
			Name:       stream.Name,
			Path:       stream.URL,
			Format:     "rtsp",
		}

		if stream.Enabled {
			asset.Status = domain.AssetStatusReady
		} else {
			asset.Status = domain.AssetStatusPending
		}

		if err := db.WithContext(ctx).Create(asset).Error; err != nil {
			log.Printf("  ⚠️  跳过流 %s: %v", stream.Name, err)
			continue
		}

		migrated++
		log.Printf("  ✓ 迁移流: %s → 资产 ID: %s", stream.Name, asset.ID)
	}

	log.Printf("✅ 成功迁移 %d/%d 个流", migrated, len(streams))
	return nil
}

func migrateAlgorithmsToOperators(ctx context.Context, db *gorm.DB) error {
	log.Println("\n[3/4] 迁移 Algorithms → Operators")

	var algorithms []LegacyAlgorithm
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
			Status:      domain.OperatorStatusEnabled,
			IsBuiltin:   false,
		}

		if err := db.WithContext(ctx).Create(operator).Error; err != nil {
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
	log.Println("\n[4/4] 清理废弃表")

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
