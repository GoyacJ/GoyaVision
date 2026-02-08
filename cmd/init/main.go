package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"goyavision/config"
	"goyavision/internal/adapter/persistence"
	"goyavision/internal/domain/identity"
	"goyavision/internal/infra/persistence/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dryRun = flag.Bool("dry-run", false, "只显示将要执行的操作，不实际执行")
	force  = flag.Bool("force", false, "强制重新初始化（删除现有数据后重新创建）")
)

func main() {
	flag.Parse()

	log.Println("GoyaVision 数据库初始化工具 v1.0")
	log.Println("====================================")

	if *dryRun {
		log.Println("⚠️  模拟运行模式（不会修改数据库）")
	}

	if *force {
		log.Println("⚠️  强制模式：将删除现有数据后重新创建")
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

	log.Println("\n📊 初始化计划:")
	log.Println("1. 创建数据库表结构")
	log.Println("2. 初始化租户数据")
	log.Println("3. 初始化权限数据")
	log.Println("4. 初始化菜单数据")
	log.Println("5. 初始化角色数据（超级管理员）")
	log.Println("6. 初始化管理员用户")

	if !confirm("\n是否继续？") && !*dryRun {
		log.Println("已取消")
		return
	}

	log.Println("\n开始初始化...")

	if err := createTables(db); err != nil {
		log.Fatalf("创建数据库表失败: %v", err)
	}

	if err := initTenants(ctx, db); err != nil {
		log.Fatalf("初始化租户失败: %v", err)
	}

	if err := initPermissions(ctx, db); err != nil {
		log.Fatalf("初始化权限失败: %v", err)
	}

	if err := initMenus(ctx, db); err != nil {
		log.Fatalf("初始化菜单失败: %v", err)
	}

	if err := initRoles(ctx, db); err != nil {
		log.Fatalf("初始化角色失败: %v", err)
	}

	if err := initAdminUser(ctx, db); err != nil {
		log.Fatalf("初始化管理员用户失败: %v", err)
	}

	if err := initSystemConfig(ctx, db); err != nil {
		log.Fatalf("初始化系统配置失败: %v", err)
	}

	log.Println("\n✅ 数据库初始化完成！")
	log.Println("\n默认管理员账号:")
	log.Println("  用户名: admin")
	log.Println("  密码: admin123")
	log.Println("  ⚠️  生产环境请立即修改密码！")
}

func createTables(db *gorm.DB) error {
	log.Println("\n[1/5] 创建数据库表结构")

	if *dryRun {
		log.Println("  （模拟运行，跳过实际创建）")
		return nil
	}

	log.Println("  创建 V1.0 表结构...")
	// 显式先迁移租户表，确保后续初始化可用
	if err := db.AutoMigrate(&model.TenantModel{}); err != nil {
		return fmt.Errorf("迁移租户表失败: %w", err)
	}
	// 再次验证表是否存在，处理某些环境下异步或缓存问题
	if !db.Migrator().HasTable("tenants") {
		log.Println("  ⚠️  警告：AutoMigrate 未能创建 tenants 表，尝试显式创建...")
		if err := db.Migrator().CreateTable(&model.TenantModel{}); err != nil {
			return fmt.Errorf("显式创建租户表失败: %w", err)
		}
	}
	if err := persistence.AutoMigrate(db); err != nil {
		return fmt.Errorf("AutoMigrate 失败: %w", err)
	}

	log.Println("  ✓ 已创建/更新以下表:")
	log.Println("    - users, roles, permissions, menus")
	log.Println("    - media_sources, media_assets")
	log.Println("    - operators, operator_versions, operator_templates, operator_dependencies")
	log.Println("    - ai_models")
	log.Println("    - workflows, workflow_nodes, workflow_edges")
	log.Println("    - tasks, artifacts")
	log.Println("    - files")
	log.Println("    - user_identities")
	log.Println("    - system_configs")
	log.Println("    - user_balances, user_subscriptions, transaction_records, point_records, usage_stats")

	// 兼容性处理：删除旧版本的 legacy 字段
	// 由于 operator 重设计去除了这些字段，如果数据库中残留会导致 GORM 插入失败（因旧字段可能为 NOT NULL）
	legacyColumns := []string{"version", "endpoint", "method", "input_schema", "output_spec", "config", "is_builtin"}
	migrator := db.Migrator()
	if migrator.HasTable(&model.OperatorModel{}) {
		log.Println("  检查并清理 Operator 旧兼容字段...")
		for _, col := range legacyColumns {
			if migrator.HasColumn(&model.OperatorModel{}, col) {
				log.Printf("    - 删除旧字段: %s", col)
				if err := migrator.DropColumn(&model.OperatorModel{}, col); err != nil {
					log.Printf("      ⚠️ 删除失败: %v", err)
				}
			}
		}

		// AI 模型重构：ai_model_id 已移入 ExecConfig.AIModel，需删除旧字段
		if migrator.HasColumn(&model.OperatorModel{}, "ai_model_id") {
			log.Println("    - 删除旧字段: ai_model_id (已迁移至 exec_config)")
			if err := migrator.DropColumn(&model.OperatorModel{}, "ai_model_id"); err != nil {
				log.Printf("      ⚠️ 删除失败: %v", err)
			}
		}
	}

	if migrator.HasTable(&model.FileModel{}) && !migrator.HasColumn(&model.FileModel{}, "Visibility") {
		log.Println("  检查并补充 files 表 visibility 列...")
		if err := db.Exec("ALTER TABLE files ADD COLUMN IF NOT EXISTS visibility int NOT NULL DEFAULT 0").Error; err != nil {
			log.Printf("    ⚠️ 添加 visibility 列失败: %v", err)
		} else {
			log.Println("    ✓ 已添加 files.visibility 列")
		}
	}

	log.Println("✅ 数据库表结构创建完成")
	return nil
}

func initTenants(ctx context.Context, db *gorm.DB) error {
	log.Println("\n[2/6] 初始化租户数据")

	if *dryRun {
		log.Println("  （模拟运行，跳过实际初始化）")
		return nil
	}

	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	var existing model.TenantModel
	err := db.WithContext(ctx).Where("code = ?", "default").First(&existing).Error

	if err == nil {
		if *force {
			log.Println("  更新默认租户...")
			existing.Name = "默认租户"
			existing.Status = 1
			if err := db.WithContext(ctx).Save(&existing).Error; err != nil {
				return fmt.Errorf("更新租户失败: %w", err)
			}
			log.Println("  ✓ 已更新默认租户")
		} else {
			log.Println("  ⊙ 默认租户已存在，跳过创建")
		}
	} else {
		log.Println("  创建默认租户...")
		tenant := &model.TenantModel{
			ID:     tenantID,
			Name:   "默认租户",
			Code:   "default",
			Status: 1,
		}
		if err := db.WithContext(ctx).Create(tenant).Error; err != nil {
			return fmt.Errorf("创建租户失败: %w", err)
		}
		log.Println("  ✓ 已创建默认租户")
	}

	log.Println("✅ 租户数据初始化完成")
	return nil
}

func initPermissions(ctx context.Context, db *gorm.DB) error {
	log.Println("\n[3/6] 初始化权限数据")
	log.Println("\n[2/5] 初始化权限数据")

	if *dryRun {
		log.Println("  （模拟运行，跳过实际初始化）")
		return nil
	}

	if *force {
		log.Println("  清理现有权限...")
		if err := db.Exec("DELETE FROM role_permissions").Error; err != nil {
			log.Printf("  ⚠️  清理权限关联失败: %v", err)
		}
		if err := db.Exec("DELETE FROM permissions").Error; err != nil {
			log.Printf("  ⚠️  清理权限失败: %v", err)
		}
	}

	permissions := []struct {
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
		{"operator:publish", "发布算子", "POST", "/api/v1/operators/*/publish", ""},
		{"operator:deprecate", "弃用算子", "POST", "/api/v1/operators/*/deprecate", ""},
		{"operator:test", "测试算子", "POST", "/api/v1/operators/*/test", ""},
		{"operator:version:list", "查看版本列表", "GET", "/api/v1/operators/*/versions", ""},
		{"operator:version:create", "创建版本", "POST", "/api/v1/operators/*/versions", ""},
		{"operator:version:activate", "激活版本", "POST", "/api/v1/operators/*/versions/activate", ""},
		{"operator:version:rollback", "回滚版本", "POST", "/api/v1/operators/*/versions/rollback", ""},
		{"operator:version:archive", "归档版本", "POST", "/api/v1/operators/*/versions/archive", ""},
		{"operator:template:list", "查看模板列表", "GET", "/api/v1/operators/templates", ""},
		{"operator:template:install", "安装模板", "POST", "/api/v1/operators/templates/install", ""},
		{"operator:dependency:list", "查看依赖列表", "GET", "/api/v1/operators/*/dependencies", ""},
		{"operator:dependency:update", "更新依赖", "PUT", "/api/v1/operators/*/dependencies", ""},
		{"operator:mcp:list", "查看MCP服务", "GET", "/api/v1/operators/mcp/servers", ""},
		{"operator:mcp:install", "安装MCP算子", "POST", "/api/v1/operators/mcp/install", ""},
		{"operator:mcp:sync", "同步MCP模板", "POST", "/api/v1/operators/mcp/sync-templates", ""},
		{"ai-model:list", "查看AI模型列表", "GET", "/api/v1/ai-models", ""},
		{"ai-model:create", "创建AI模型", "POST", "/api/v1/ai-models", ""},
		{"ai-model:update", "更新AI模型", "PUT", "/api/v1/ai-models/*", ""},
		{"ai-model:delete", "删除AI模型", "DELETE", "/api/v1/ai-models/*", ""},
		{"ai-model:test", "测试AI模型连接", "POST", "/api/v1/ai-models/*/test-connection", ""},
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
		{"tenant:list", "查看租户列表", "GET", "/api/v1/tenants", ""},
		{"tenant:create", "创建租户", "POST", "/api/v1/tenants", ""},
		{"tenant:update", "更新租户", "PUT", "/api/v1/tenants/*", ""},
		{"tenant:delete", "删除租户", "DELETE", "/api/v1/tenants/*", ""},
		{"system:config:view", "查看系统配置", "GET", "/api/v1/public/config", "允许查看系统配置"},
		{"system:config:update", "修改系统配置", "PUT", "/api/v1/system/config", "允许修改系统配置"},
	}

	addedPerms := 0
	skippedPerms := 0
	for _, p := range permissions {
		var existing model.PermissionModel
		err := db.WithContext(ctx).Where("code = ?", p.Code).First(&existing).Error
		if err == nil {
			if *force {
				existing.Name = p.Name
				existing.Method = p.Method
				existing.Path = p.Path
				existing.Description = p.Description
				if err := db.WithContext(ctx).Save(&existing).Error; err != nil {
					log.Printf("  ⚠️  更新权限失败 %s: %v", p.Code, err)
				} else {
					log.Printf("  ✓ 更新权限: %s", p.Name)
				}
			} else {
				skippedPerms++
			}
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
		if err := db.WithContext(ctx).Create(perm).Error; err != nil {
			log.Printf("  ⚠️  创建权限失败 %s: %v", p.Code, err)
		} else {
			addedPerms++
			log.Printf("  ✓ 创建权限: %s", p.Name)
		}
	}

	if skippedPerms > 0 {
		log.Printf("  ⊙ 跳过已存在权限: %d 个", skippedPerms)
	}
	log.Printf("  ✓ 新增权限: %d 个", addedPerms)
	log.Println("✅ 权限数据初始化完成")
	return nil
}

func initMenus(ctx context.Context, db *gorm.DB) error {
	log.Println("\n[4/6] 初始化菜单数据")

	if *dryRun {
		log.Println("  （模拟运行，跳过实际初始化）")
		return nil
	}

	if *force {
		log.Println("  清理现有菜单...")
		if err := db.Exec("DELETE FROM role_menus").Error; err != nil {
			log.Printf("  ⚠️  清理菜单关联失败: %v", err)
		}
		if err := db.Exec("DELETE FROM menus").Error; err != nil {
			log.Printf("  ⚠️  清理菜单失败: %v", err)
		}
	}

	menus := []struct {
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
		{uuid.MustParse("00000000-0000-0000-0000-000000000020"), nil, "operator", "算子中心", 1, "/operator", "Cpu", "", "", 3, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000022"), ptrUUID("00000000-0000-0000-0000-000000000020"), "operator-list", "我的算子", 2, "/operator/list", "List", "operator/index", "operator:list", 1, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000023"), ptrUUID("00000000-0000-0000-0000-000000000020"), "ai-model", "AI模型", 2, "/operator/ai-model", "Connection", "operator/ai-model/index", "ai-model:list", 2, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000021"), ptrUUID("00000000-0000-0000-0000-000000000020"), "mcp-market", "MCP市场", 2, "/operator/mcp-market", "Shop", "operator-marketplace/index", "operator:template:list", 3, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000030"), nil, "workflow", "工作流", 2, "/workflows", "Connection", "workflow/index", "workflow:list", 4, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000040"), nil, "task", "任务管理", 2, "/tasks", "List", "task/index", "task:list", 6, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000001"), nil, "system", "系统管理", 1, "/system", "Setting", "", "", 100, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000002"), ptrUUID("00000000-0000-0000-0000-000000000001"), "system:user", "用户管理", 2, "/system/user", "User", "system/user/index", "user:list", 1, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000003"), ptrUUID("00000000-0000-0000-0000-000000000001"), "system:role", "角色管理", 2, "/system/role", "UserFilled", "system/role/index", "role:list", 2, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000004"), ptrUUID("00000000-0000-0000-0000-000000000001"), "system:menu", "菜单管理", 2, "/system/menu", "Menu", "system/menu/index", "menu:list", 3, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000005"), ptrUUID("00000000-0000-0000-0000-000000000001"), "system:file", "文件管理", 2, "/system/file", "Document", "system/file/index", "file:list", 4, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000006"), ptrUUID("00000000-0000-0000-0000-000000000001"), "system:tenant", "租户管理", 2, "/system/tenant", "OfficeBuilding", "system/tenant/index", "tenant:list", 5, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000007"), ptrUUID("00000000-0000-0000-0000-000000000001"), "system:config", "系统配置", 2, "/system/config", "Tools", "system/config/index", "system:config:view", 6, true},
	}

	addedMenus := 0
	skippedMenus := 0
	for _, m := range menus {
		var existing model.MenuModel
		err := db.WithContext(ctx).Where("code = ?", m.Code).First(&existing).Error
		if err == nil {
			if *force {
				existing.Name = m.Name
				existing.Type = m.Type
				existing.Path = m.Path
				existing.Icon = m.Icon
				existing.Component = m.Component
				existing.Permission = m.Permission
				existing.Sort = m.Sort
				existing.Visible = m.Visible
				existing.Status = int(identity.MenuStatusEnabled)
				if err := db.WithContext(ctx).Save(&existing).Error; err != nil {
					log.Printf("  ⚠️  更新菜单失败 %s: %v", m.Code, err)
				} else {
					log.Printf("  ✓ 更新菜单: %s", m.Name)
				}
			} else {
				skippedMenus++
			}
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
		if err := db.WithContext(ctx).Create(menu).Error; err != nil {
			log.Printf("  ⚠️  创建菜单失败 %s: %v", m.Code, err)
		} else {
			addedMenus++
			log.Printf("  ✓ 创建菜单: %s", m.Name)
		}
	}

	if skippedMenus > 0 {
		log.Printf("  ⊙ 跳过已存在菜单: %d 个", skippedMenus)
	}
	log.Printf("  ✓ 新增菜单: %d 个", addedMenus)
	log.Println("✅ 菜单数据初始化完成")
	return nil
}

func initRoles(ctx context.Context, db *gorm.DB) error {
	log.Println("\n[5/6] 初始化角色数据")

	if *dryRun {
		log.Println("  （模拟运行，跳过实际初始化）")
		return nil
	}

	roleID := uuid.MustParse("00000000-0000-0000-0000-000000000100")
	var existingRole model.RoleModel
	err := db.WithContext(ctx).Where("code = ?", "super_admin").First(&existingRole).Error

	if err == nil {
		if *force {
			log.Println("  更新超级管理员角色...")
			existingRole.Name = "超级管理员"
			existingRole.Description = "拥有所有权限"
			existingRole.Status = int(identity.RoleStatusEnabled)
			if err := db.WithContext(ctx).Save(&existingRole).Error; err != nil {
				return fmt.Errorf("更新角色失败: %w", err)
			}
			log.Println("  ✓ 已更新超级管理员角色")
		} else {
			log.Println("  ⊙ 超级管理员角色已存在，跳过创建")
		}
	} else {
		log.Println("  创建超级管理员角色...")
		role := &model.RoleModel{
			ID:          roleID,
			Code:        "super_admin",
			Name:        "超级管理员",
			Description: "拥有所有权限",
			Status:      int(identity.RoleStatusEnabled),
		}
		if err := db.WithContext(ctx).Create(role).Error; err != nil {
			return fmt.Errorf("创建角色失败: %w", err)
		}
		log.Println("  ✓ 已创建超级管理员角色")
	}

	log.Println("  检查并分配权限和菜单...")
	var superAdminRole model.RoleModel
	if err := db.WithContext(ctx).Where("code = ?", "super_admin").First(&superAdminRole).Error; err != nil {
		return fmt.Errorf("获取角色失败: %w", err)
	}

	if *force {
		log.Println("  清理现有权限和菜单关联...")
		db.Exec("DELETE FROM role_permissions WHERE role_model_id = ?", superAdminRole.ID)
		db.Exec("DELETE FROM role_menus WHERE role_model_id = ?", superAdminRole.ID)
	}

	// Create default user role
	userRoleID := uuid.MustParse("00000000-0000-0000-0000-000000000101")
	var existingUserRole model.RoleModel
	err = db.WithContext(ctx).Where("code = ?", "user").First(&existingUserRole).Error

	if err == nil {
		if *force {
			log.Println("  更新普通用户角色...")
			existingUserRole.Name = "普通用户"
			existingUserRole.Description = "默认普通用户"
			existingUserRole.Status = int(identity.RoleStatusEnabled)
			existingUserRole.IsDefault = true
			if err := db.WithContext(ctx).Save(&existingUserRole).Error; err != nil {
				return fmt.Errorf("更新角色失败: %w", err)
			}
			log.Println("  ✓ 已更新普通用户角色")
		} else {
			log.Println("  ⊙ 普通用户角色已存在，跳过创建")
		}
	} else {
		log.Println("  创建普通用户角色...")
		role := &model.RoleModel{
			ID:          userRoleID,
			Code:        "user",
			Name:        "普通用户",
			Description: "默认普通用户",
			Status:      int(identity.RoleStatusEnabled),
			IsDefault:   true,
		}
		if err := db.WithContext(ctx).Create(role).Error; err != nil {
			return fmt.Errorf("创建角色失败: %w", err)
		}
		log.Println("  ✓ 已创建普通用户角色")
	}

	// Assign basic permissions to user role
	var userRole model.RoleModel
	if err := db.WithContext(ctx).Where("code = ?", "user").First(&userRole).Error; err == nil {
		basicPermCodes := []string{"asset:list", "asset:create", "operator:list", "workflow:list", "task:list"}
		var basicPerms []model.PermissionModel
		if err := db.WithContext(ctx).Where("code IN ?", basicPermCodes).Find(&basicPerms).Error; err == nil {
			if *force {
				db.Exec("DELETE FROM role_permissions WHERE role_model_id = ?", userRole.ID)
			}
			for _, perm := range basicPerms {
				var count int64
				if err := db.WithContext(ctx).Table("role_permissions").
					Where("role_model_id = ? AND permission_model_id = ?", userRole.ID, perm.ID).
					Count(&count).Error; err == nil && count == 0 {
					db.WithContext(ctx).Exec(
						"INSERT INTO role_permissions (role_model_id, permission_model_id) VALUES (?, ?)",
						userRole.ID, perm.ID,
					)
				}
			}
		}
	}

	var allPermissions []model.PermissionModel
	if err := db.WithContext(ctx).Find(&allPermissions).Error; err != nil {
		return fmt.Errorf("查询权限失败: %w", err)
	}

	addedPerms := 0
	for _, perm := range allPermissions {
		var count int64
		if err := db.WithContext(ctx).Table("role_permissions").
			Where("role_model_id = ? AND permission_model_id = ?", superAdminRole.ID, perm.ID).
			Count(&count).Error; err == nil && count == 0 {
			if err := db.WithContext(ctx).Exec(
				"INSERT INTO role_permissions (role_model_id, permission_model_id) VALUES (?, ?)",
				superAdminRole.ID, perm.ID,
			).Error; err != nil {
				log.Printf("  ⚠️  分配权限失败 %s: %v", perm.Code, err)
			} else {
				addedPerms++
			}
		}
	}
	if addedPerms > 0 {
		log.Printf("  ✓ 新增分配 %d 个权限", addedPerms)
	}
	log.Printf("  ✓ 权限总计: %d 个", len(allPermissions))

	var allMenus []model.MenuModel
	if err := db.WithContext(ctx).Find(&allMenus).Error; err != nil {
		return fmt.Errorf("查询菜单失败: %w", err)
	}

	addedMenus := 0
	for _, menu := range allMenus {
		var count int64
		if err := db.WithContext(ctx).Table("role_menus").
			Where("role_model_id = ? AND menu_model_id = ?", superAdminRole.ID, menu.ID).
			Count(&count).Error; err == nil && count == 0 {
			if err := db.WithContext(ctx).Exec(
				"INSERT INTO role_menus (role_model_id, menu_model_id) VALUES (?, ?)",
				superAdminRole.ID, menu.ID,
			).Error; err != nil {
				log.Printf("  ⚠️  分配菜单失败 %s: %v", menu.Code, err)
			} else {
				addedMenus++
			}
		}
	}
	if addedMenus > 0 {
		log.Printf("  ✓ 新增分配 %d 个菜单", addedMenus)
	}
	log.Printf("  ✓ 菜单总计: %d 个", len(allMenus))

	log.Println("✅ 角色数据初始化完成")
	return nil
}

func initAdminUser(ctx context.Context, db *gorm.DB) error {
	log.Println("\n[6/6] 初始化管理员用户")

	if *dryRun {
		log.Println("  （模拟运行，跳过实际初始化）")
		return nil
	}

	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000200")
	var existingUser model.UserModel
	err := db.WithContext(ctx).Where("username = ?", "admin").First(&existingUser).Error

	if err == nil {
		if *force {
			log.Println("  更新管理员用户...")
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
			if err != nil {
				return fmt.Errorf("加密密码失败: %w", err)
			}
			existingUser.Password = string(hashedPassword)
			existingUser.Nickname = "管理员"
			existingUser.Status = int(identity.UserStatusEnabled)
			existingUser.TenantID = &tenantID
			if err := db.WithContext(ctx).Save(&existingUser).Error; err != nil {
				return fmt.Errorf("更新用户失败: %w", err)
			}

			roleID := uuid.MustParse("00000000-0000-0000-0000-000000000100")
			db.Exec("DELETE FROM user_roles WHERE user_model_id = ?", existingUser.ID)
			if err := db.WithContext(ctx).Exec(
				"INSERT INTO user_roles (user_model_id, role_model_id) VALUES (?, ?)",
				existingUser.ID, roleID,
			).Error; err != nil {
				log.Printf("  ⚠️  分配角色失败: %v", err)
			}

			// 初始化或更新管理员余额
			db.WithContext(ctx).Where("user_id = ?", existingUser.ID).FirstOrCreate(&model.UserBalance{
				UserID:  existingUser.ID,
				Balance: 9999.99,
				Points:  100000,
				Level:   "Administrator",
			})

			log.Println("  ✓ 已更新管理员用户")
		} else {
			log.Println("  ⊙ 管理员用户已存在，跳过创建")
		}
		return nil
	}

	log.Println("  创建管理员用户...")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("加密密码失败: %w", err)
	}

	user := &model.UserModel{
		ID:       userID,
		Username: "admin",
		Password: string(hashedPassword),
		Nickname: "管理员",
		Status:   int(identity.UserStatusEnabled),
		TenantID: &tenantID,
	}

	if err := db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	roleID := uuid.MustParse("00000000-0000-0000-0000-000000000100")
	if err := db.WithContext(ctx).Exec(
		"INSERT INTO user_roles (user_model_id, role_model_id) VALUES (?, ?)",
		user.ID, roleID,
	).Error; err != nil {
		return fmt.Errorf("分配角色失败: %w", err)
	}

	// 初始化管理员余额
	if err := db.WithContext(ctx).Create(&model.UserBalance{
		UserID:  user.ID,
		Balance: 9999.99,
		Points:  100000,
		Level:   "Administrator",
	}).Error; err != nil {
		log.Printf("  ⚠️  初始化管理员余额失败: %v", err)
	}

	log.Println("  ✓ 已创建管理员用户")
	log.Println("✅ 管理员用户初始化完成")
	return nil
}

func ptrUUID(s string) *uuid.UUID {
	id := uuid.MustParse(s)
	return &id
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

func initSystemConfig(ctx context.Context, db *gorm.DB) error {
	log.Println("\n[7/6] 初始化系统配置")

	if *dryRun {
		log.Println("  （模拟运行，跳过实际初始化）")
		return nil
	}

	configs := []model.SystemConfigModel{
		{
			Key:         "system.home_path",
			Value:       datatypes.JSON([]byte(`"/assets"`)),
			Description: "系统默认首页路径",
		},
		{
			Key:         "system.public_menus",
			Value:       datatypes.JSON([]byte(`[]`)),
			Description: "未登录用户可见的菜单ID列表",
		},
	}

	for _, c := range configs {
		var existing model.SystemConfigModel
		err := db.WithContext(ctx).Where("key = ?", c.Key).First(&existing).Error
		if err == nil {
			if *force {
				existing.Value = c.Value
				existing.Description = c.Description
				if err := db.WithContext(ctx).Save(&existing).Error; err != nil {
					return fmt.Errorf("更新配置失败 %s: %w", c.Key, err)
				}
				log.Printf("  ✓ 更新配置: %s", c.Key)
			} else {
				log.Printf("  ⊙ 跳过已存在配置: %s", c.Key)
			}
		} else {
			if err := db.WithContext(ctx).Create(&c).Error; err != nil {
				return fmt.Errorf("创建配置失败 %s: %w", c.Key, err)
			}
			log.Printf("  ✓ 创建配置: %s", c.Key)
		}
	}

	log.Println("✅ 系统配置初始化完成")
	return nil
}
