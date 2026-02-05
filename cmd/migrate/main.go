package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"goyavision/config"
	"goyavision/internal/domain/identity"
	"goyavision/internal/domain/media"
	"goyavision/internal/domain/operator"
	"goyavision/internal/infra/persistence/model"

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
	log.Println("0. 创建数据库表结构（如果不存在）")
	log.Println("1. 更新菜单和权限（V1.0 新功能）")
	log.Println("2. 迁移 streams → media_sources（媒体源）")
	log.Println("3. 迁移 streams → media_assets（媒体资产）")
	log.Println("4. 迁移 algorithms → operators（算子）")
	log.Println("5. 清理废弃表（algorithm_bindings、inference_results、streams、record_sessions）")

	if !confirm("\n是否继续？") && !*dryRun {
		log.Println("已取消")
		return
	}

	log.Println("\n开始迁移...")

	if err := createTables(db); err != nil {
		log.Fatalf("创建数据库表失败: %v", err)
	}

	if err := updateMenusAndPermissions(ctx, db); err != nil {
		log.Fatalf("更新菜单和权限失败: %v", err)
	}

	if err := migrateStreamsToSources(ctx, db); err != nil {
		log.Fatalf("迁移 streams → media_sources 失败: %v", err)
	}

	if err := migrateStreamsToAssets(ctx, db); err != nil {
		log.Fatalf("迁移 streams → media_assets 失败: %v", err)
	}

	if err := migrateAlgorithmsToOperators(ctx, db); err != nil {
		log.Fatalf("迁移 algorithms → operators 失败: %v", err)
	}

	if !*dryRun {
		if err := cleanupOldTables(db); err != nil {
			log.Fatalf("清理旧表失败: %v", err)
		}
	}

	log.Println("\n✅ 迁移完成！")
}

// LegacyStream 旧流结构（用于迁移）
type LegacyStream struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	URL       string
	Name      string
	Type      string
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (LegacyStream) TableName() string { return "streams" }

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
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (LegacyAlgorithm) TableName() string { return "algorithms" }

func createTables(db *gorm.DB) error {
	log.Println("\n[0/5] 创建数据库表结构")

	if *dryRun {
		log.Println("  （模拟运行，跳过实际创建）")
		return nil
	}

	log.Println("  创建 V1.0 表结构...")
	if err := db.AutoMigrate(
		&model.UserModel{},
		&model.RoleModel{},
		&model.PermissionModel{},
		&model.MenuModel{},
		&model.MediaSourceModel{},
		&model.MediaAssetModel{},
		&model.OperatorModel{},
		&model.WorkflowModel{},
		&model.WorkflowNodeModel{},
		&model.WorkflowEdgeModel{},
		&model.TaskModel{},
		&model.ArtifactModel{},
		&model.FileModel{},
	); err != nil {
		return fmt.Errorf("AutoMigrate 失败: %w", err)
	}

	log.Println("  ✓ 已创建/更新以下表:")
	log.Println("    - users, roles, permissions, menus")
	log.Println("    - media_sources, media_assets")
	log.Println("    - operators")
	log.Println("    - workflows, workflow_nodes, workflow_edges")
	log.Println("    - tasks, artifacts")
	log.Println("    - files")

	log.Println("✅ 数据库表结构创建完成")
	return nil
}

func updateMenusAndPermissions(ctx context.Context, db *gorm.DB) error {
	log.Println("\n[1/5] 更新菜单和权限")

	if *dryRun {
		log.Println("  （模拟运行，跳过实际更新）")
		return nil
	}

	log.Println("  清理旧菜单...")
	oldMenuCodes := []string{"stream", "algorithm", "inference", "legacy", "legacy:stream"}

	for _, code := range oldMenuCodes {
		var menu model.MenuModel
		if err := db.Where("code = ?", code).First(&menu).Error; err != nil {
			continue
		}

		if err := db.Exec("DELETE FROM role_menus WHERE menu_id = ?", menu.ID).Error; err != nil {
			log.Printf("  ⚠️  删除菜单关联失败 %s: %v", code, err)
			continue
		}

		if err := db.Where("id = ?", menu.ID).Delete(&model.MenuModel{}).Error; err != nil {
			log.Printf("  ⚠️  删除旧菜单 %s 失败: %v", code, err)
		} else {
			log.Printf("  ✓ 删除旧菜单: %s", code)
		}
	}

	log.Println("  清理旧权限...")
	oldPermCodes := []string{
		"stream:list", "stream:create", "stream:update", "stream:delete",
		"record:start", "record:stop", "record:list",
		"preview:start", "preview:stop",
		"algorithm:list", "algorithm:create", "algorithm:update", "algorithm:delete",
		"binding:list", "binding:create", "binding:update", "binding:delete",
		"inference:list",
	}

	for _, code := range oldPermCodes {
		var perm model.PermissionModel
		if err := db.Where("code = ?", code).First(&perm).Error; err != nil {
			continue
		}

		if err := db.Exec("DELETE FROM role_permissions WHERE permission_id = ?", perm.ID).Error; err != nil {
			log.Printf("  ⚠️  删除权限关联失败 %s: %v", code, err)
			continue
		}

		if err := db.Where("id = ?", perm.ID).Delete(&model.PermissionModel{}).Error; err != nil {
			log.Printf("  ⚠️  删除旧权限 %s 失败: %v", code, err)
		} else {
			log.Printf("  ✓ 删除旧权限: %s", code)
		}
	}

	log.Println("  添加新菜单...")
	newMenus := []struct {
		ID         uuid.UUID
		ParentID   *uuid.UUID
		Code       string
		Name       string
		Type       int
		Path       string
		Icon       string
		Component  string
		Permission string
		Sort       int
		Visible    bool
	}{
		{uuid.MustParse("00000000-0000-0000-0000-000000000010"), nil, "asset", "媒体资产", 2, "/assets", "Files", "asset/index", "asset:list", 1, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000011"), nil, "source", "媒体源", 2, "/sources", "VideoCamera", "source/index", "source:list", 2, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000020"), nil, "operator", "算子管理", 2, "/operators", "Cpu", "operator/index", "operator:list", 3, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000030"), nil, "workflow", "工作流", 2, "/workflows", "Connection", "workflow/index", "workflow:list", 4, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000040"), nil, "task", "任务管理", 2, "/tasks", "List", "task/index", "task:list", 5, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000001"), nil, "system", "系统管理", 1, "/system", "Setting", "", "", 100, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000002"), ptrUUID("00000000-0000-0000-0000-000000000001"), "system:user", "用户管理", 2, "/system/user", "User", "system/user/index", "user:list", 1, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000003"), ptrUUID("00000000-0000-0000-0000-000000000001"), "system:role", "角色管理", 2, "/system/role", "UserFilled", "system/role/index", "role:list", 2, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000004"), ptrUUID("00000000-0000-0000-0000-000000000001"), "system:menu", "菜单管理", 2, "/system/menu", "Menu", "system/menu/index", "menu:list", 3, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000005"), ptrUUID("00000000-0000-0000-0000-000000000001"), "system:file", "文件管理", 2, "/system/file", "Document", "system/file/index", "file:list", 4, true},
	}

	addedMenus := 0
	for _, m := range newMenus {
		var existing model.MenuModel
		err := db.Where("code = ?", m.Code).First(&existing).Error
		if err == nil {
			log.Printf("  ⊙ 菜单已存在，跳过: %s", m.Name)
			continue
		}

		menu := &model.MenuModel{
			ID:         m.ID,
			ParentID:   m.ParentID,
			Code:       m.Code,
			Name:       m.Name,
			Type:       m.Type,
			Path:       m.Path,
			Icon:       m.Icon,
			Component:  m.Component,
			Permission: m.Permission,
			Sort:       m.Sort,
			Visible:    m.Visible,
			Status:     int(identity.MenuStatusEnabled),
		}
		if err := db.Create(menu).Error; err != nil {
			log.Printf("  ⚠️  创建菜单失败 %s: %v", m.Name, err)
		} else {
			addedMenus++
			log.Printf("  ✓ 创建新菜单: %s", m.Name)
		}
	}
	log.Printf("  ✓ 新增菜单: %d 个", addedMenus)

	log.Println("  添加新权限...")
	newPermissions := []struct {
		Code        string
		Name        string
		Method      string
		Path        string
		Description string
	}{
		{"asset:list", "查看媒体资产列表", "GET", "/api/v1/assets", ""},
		{"asset:create", "创建媒体资产", "POST", "/api/v1/assets", ""},
		{"asset:update", "更新媒体资产", "PUT", "/api/v1/assets/*", ""},
		{"asset:delete", "删除媒体资产", "DELETE", "/api/v1/assets/*", ""},
		{"source:list", "查看媒体源列表", "GET", "/api/v1/sources", ""},
		{"source:create", "创建媒体源", "POST", "/api/v1/sources", ""},
		{"source:update", "更新媒体源", "PUT", "/api/v1/sources/*", ""},
		{"source:delete", "删除媒体源", "DELETE", "/api/v1/sources/*", ""},
		{"operator:list", "查看算子列表", "GET", "/api/v1/operators", ""},
		{"operator:create", "创建算子", "POST", "/api/v1/operators", ""},
		{"operator:update", "更新算子", "PUT", "/api/v1/operators/*", ""},
		{"operator:delete", "删除算子", "DELETE", "/api/v1/operators/*", ""},
		{"operator:enable", "启用算子", "PUT", "/api/v1/operators/*/enable", ""},
		{"operator:disable", "禁用算子", "PUT", "/api/v1/operators/*/disable", ""},
		{"workflow:list", "查看工作流列表", "GET", "/api/v1/workflows", ""},
		{"workflow:create", "创建工作流", "POST", "/api/v1/workflows", ""},
		{"workflow:update", "更新工作流", "PUT", "/api/v1/workflows/*", ""},
		{"workflow:delete", "删除工作流", "DELETE", "/api/v1/workflows/*", ""},
		{"workflow:enable", "启用工作流", "PUT", "/api/v1/workflows/*/enable", ""},
		{"workflow:disable", "禁用工作流", "PUT", "/api/v1/workflows/*/disable", ""},
		{"workflow:trigger", "触发工作流", "POST", "/api/v1/workflows/*/trigger", ""},
		{"task:list", "查看任务列表", "GET", "/api/v1/tasks", ""},
		{"task:create", "创建任务", "POST", "/api/v1/tasks", ""},
		{"task:update", "更新任务", "PUT", "/api/v1/tasks/*", ""},
		{"task:delete", "删除任务", "DELETE", "/api/v1/tasks/*", ""},
		{"task:cancel", "取消任务", "POST", "/api/v1/tasks/*/cancel", ""},
		{"artifact:list", "查看产物列表", "GET", "/api/v1/artifacts", ""},
		{"artifact:delete", "删除产物", "DELETE", "/api/v1/artifacts/*", ""},
		{"user:list", "查看用户列表", "GET", "/api/v1/users", ""},
		{"user:create", "创建用户", "POST", "/api/v1/users", ""},
		{"user:update", "更新用户", "PUT", "/api/v1/users/*", ""},
		{"user:delete", "删除用户", "DELETE", "/api/v1/users/*", ""},
		{"role:list", "查看角色列表", "GET", "/api/v1/roles", ""},
		{"role:create", "创建角色", "POST", "/api/v1/roles", ""},
		{"role:update", "更新角色", "PUT", "/api/v1/roles/*", ""},
		{"role:delete", "删除角色", "DELETE", "/api/v1/roles/*", ""},
		{"menu:list", "查看菜单列表", "GET", "/api/v1/menus", ""},
		{"menu:create", "创建菜单", "POST", "/api/v1/menus", ""},
		{"menu:update", "更新菜单", "PUT", "/api/v1/menus/*", ""},
		{"menu:delete", "删除菜单", "DELETE", "/api/v1/menus/*", ""},
		{"file:list", "查看文件列表", "GET", "/api/v1/files", ""},
		{"file:create", "上传文件", "POST", "/api/v1/files", ""},
		{"file:update", "更新文件", "PUT", "/api/v1/files/*", ""},
		{"file:delete", "删除文件", "DELETE", "/api/v1/files/*", ""},
		{"file:download", "下载文件", "GET", "/api/v1/files/*/download", ""},
	}

	addedPerms := 0
	for _, p := range newPermissions {
		var existing model.PermissionModel
		err := db.Where("code = ?", p.Code).First(&existing).Error
		if err == nil {
			continue
		}

		perm := &model.PermissionModel{
			ID:          uuid.New(),
			Code:        p.Code,
			Name:        p.Name,
			Method:      p.Method,
			Path:        p.Path,
			Description: p.Description,
		}
		if err := db.Create(perm).Error; err != nil {
			log.Printf("  ⚠️  创建权限失败 %s: %v", p.Code, err)
		} else {
			addedPerms++
		}
	}
	log.Printf("  ✓ 新增权限: %d 个", addedPerms)

	log.Println("  更新超级管理员角色权限...")
	var superAdminRole model.RoleModel
	if err := db.Where("code = ?", "super_admin").First(&superAdminRole).Error; err == nil {
		db.Exec("DELETE FROM role_permissions WHERE role_id = ?", superAdminRole.ID)

		var allPermissions []model.PermissionModel
		db.Find(&allPermissions)
		for _, perm := range allPermissions {
			db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
				superAdminRole.ID, perm.ID)
		}

		db.Exec("DELETE FROM role_menus WHERE role_id = ?", superAdminRole.ID)
		var allMenus []model.MenuModel
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

func migrateStreamsToSources(ctx context.Context, db *gorm.DB) error {
	log.Println("\n[2/5] 迁移 Streams → MediaSources")

	var streams []LegacyStream
	if err := db.Find(&streams).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("  未找到 streams 表或数据，跳过")
			return nil
		}
		return err
	}

	log.Printf("找到 %d 个流", len(streams))

	if *dryRun {
		log.Println("  （模拟运行，跳过实际迁移）")
		return nil
	}

	migrated := 0
	for _, stream := range streams {
		sourceType := media.SourceTypePull
		if stream.Type == "push" {
			sourceType = media.SourceTypePush
		}

		protocol := "rtsp"
		if stream.URL != "" {
			if len(stream.URL) > 4 {
				prefix := stream.URL[:4]
				if prefix == "rtmp" {
					protocol = "rtmp"
				} else if prefix == "http" {
					protocol = "hls"
				}
			}
		}

		source := &model.MediaSourceModel{
			ID:            stream.ID,
			Name:          stream.Name,
			PathName:      media.GeneratePathName(stream.Name),
			Type:          string(sourceType),
			URL:           stream.URL,
			Protocol:      protocol,
			Enabled:       stream.Enabled,
			RecordEnabled: false,
		}

		if err := db.WithContext(ctx).Create(source).Error; err != nil {
			log.Printf("  ⚠️  跳过流 %s: %v", stream.Name, err)
			continue
		}

		migrated++
		log.Printf("  ✓ 迁移流: %s → 媒体源 ID: %s", stream.Name, source.ID)
	}

	log.Printf("✅ 成功迁移 %d/%d 个流到媒体源", migrated, len(streams))
	return nil
}

func migrateStreamsToAssets(ctx context.Context, db *gorm.DB) error {
	log.Println("\n[3/5] 迁移 Streams → MediaAssets")

	var streams []LegacyStream
	if err := db.Find(&streams).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("  未找到 streams 表或数据，跳过")
			return nil
		}
		return err
	}

	log.Printf("找到 %d 个流", len(streams))

	if *dryRun {
		log.Println("  （模拟运行，跳过实际迁移）")
		return nil
	}

	migrated := 0
	for _, stream := range streams {
		assetType := media.AssetTypeStream
		sourceID := stream.ID

		asset := &model.MediaAssetModel{
			ID:         uuid.New(),
			Type:       string(assetType),
			SourceType: string(media.AssetSourceLive),
			SourceID:   &sourceID,
			Name:       stream.Name,
			Path:       stream.URL,
			Format:     "rtsp",
			Status:     string(media.AssetStatusPending),
		}

		if stream.Enabled {
			asset.Status = string(media.AssetStatusReady)
		}

		if err := db.WithContext(ctx).Create(asset).Error; err != nil {
			log.Printf("  ⚠️  跳过流 %s: %v", stream.Name, err)
			continue
		}

		migrated++
		log.Printf("  ✓ 迁移流: %s → 资产 ID: %s", stream.Name, asset.ID)
	}

	log.Printf("✅ 成功迁移 %d/%d 个流到媒体资产", migrated, len(streams))
	return nil
}

func migrateAlgorithmsToOperators(ctx context.Context, db *gorm.DB) error {
	log.Println("\n[4/5] 迁移 Algorithms → Operators")

	var algorithms []LegacyAlgorithm
	if err := db.Find(&algorithms).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("  未找到 algorithms 表或数据，跳过")
			return nil
		}
		return err
	}

	log.Printf("找到 %d 个算法", len(algorithms))

	if *dryRun {
		log.Println("  （模拟运行，跳过实际迁移）")
		return nil
	}

	migrated := 0
	for _, alg := range algorithms {
		category := operator.CategoryAnalysis
		opType := operator.Type("legacy_" + alg.Type)
		if opType == "" {
			opType = operator.TypeObjectDetection
		}

		var inputSchema map[string]interface{}
		if len(alg.InputSpec) > 0 {
			if err := json.Unmarshal(alg.InputSpec, &inputSchema); err != nil {
				log.Printf("  ⚠️  解析 InputSpec 失败 %s: %v", alg.Name, err)
				inputSchema = make(map[string]interface{})
			}
		}

		var outputSpec map[string]interface{}
		if len(alg.OutputSpec) > 0 {
			if err := json.Unmarshal(alg.OutputSpec, &outputSpec); err != nil {
				log.Printf("  ⚠️  解析 OutputSpec 失败 %s: %v", alg.Name, err)
				outputSpec = make(map[string]interface{})
			}
		}

		var config map[string]interface{}
		if len(alg.Config) > 0 {
			if err := json.Unmarshal(alg.Config, &config); err != nil {
				log.Printf("  ⚠️  解析 Config 失败 %s: %v", alg.Name, err)
				config = make(map[string]interface{})
			}
		}

		inputSchemaJSON, _ := json.Marshal(inputSchema)
		outputSpecJSON, _ := json.Marshal(outputSpec)
		configJSON, _ := json.Marshal(config)

		operator := &model.OperatorModel{
			ID:          uuid.New(),
			Code:        alg.Code,
			Name:        alg.Name,
			Description: alg.Description,
			Category:    string(category),
			Type:        string(opType),
			Version:     "1.0.0",
			Endpoint:    alg.Endpoint,
			Method:      "POST",
			InputSchema: inputSchemaJSON,
			OutputSpec:  outputSpecJSON,
			Config:      configJSON,
			Status:      string(operator.StatusEnabled),
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
	log.Println("\n[5/5] 清理废弃表")

	tables := []string{
		"algorithm_bindings",
		"inference_results",
		"streams",
		"record_sessions",
	}

	for _, table := range tables {
		if !db.Migrator().HasTable(table) {
			log.Printf("  ⊙ 表不存在，跳过: %s", table)
			continue
		}

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
