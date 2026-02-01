package persistence

import (
	"context"
	"log"

	"goyavision/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// InitializeData 初始化默认数据
func InitializeData(db *gorm.DB) error {
	ctx := context.Background()
	repo := NewRepository(db)

	if err := initPermissions(ctx, repo); err != nil {
		return err
	}

	if err := initMenus(ctx, repo); err != nil {
		return err
	}

	if err := initRoles(ctx, repo); err != nil {
		return err
	}

	if err := initAdminUser(ctx, repo); err != nil {
		return err
	}

	log.Println("初始化数据完成")
	return nil
}

func initPermissions(ctx context.Context, repo *repository) error {
	permissions := []struct {
		Code        string
		Name        string
		Method      string
		Path        string
		Description string
	}{
		{"stream:list", "查看视频流列表", "GET", "/api/v1/streams", ""},
		{"stream:create", "创建视频流", "POST", "/api/v1/streams", ""},
		{"stream:update", "更新视频流", "PUT", "/api/v1/streams/*", ""},
		{"stream:delete", "删除视频流", "DELETE", "/api/v1/streams/*", ""},
		{"algorithm:list", "查看算法列表", "GET", "/api/v1/algorithms", ""},
		{"algorithm:create", "创建算法", "POST", "/api/v1/algorithms", ""},
		{"algorithm:update", "更新算法", "PUT", "/api/v1/algorithms/*", ""},
		{"algorithm:delete", "删除算法", "DELETE", "/api/v1/algorithms/*", ""},
		{"binding:list", "查看算法绑定列表", "GET", "/api/v1/streams/*/algorithm-bindings", ""},
		{"binding:create", "创建算法绑定", "POST", "/api/v1/streams/*/algorithm-bindings", ""},
		{"binding:update", "更新算法绑定", "PUT", "/api/v1/streams/*/algorithm-bindings/*", ""},
		{"binding:delete", "删除算法绑定", "DELETE", "/api/v1/streams/*/algorithm-bindings/*", ""},
		{"record:start", "启动录制", "POST", "/api/v1/streams/*/record/start", ""},
		{"record:stop", "停止录制", "POST", "/api/v1/streams/*/record/stop", ""},
		{"record:list", "查看录制会话", "GET", "/api/v1/streams/*/record/sessions", ""},
		{"inference:list", "查看推理结果", "GET", "/api/v1/inference_results", ""},
		{"preview:start", "启动预览", "GET", "/api/v1/streams/*/preview/start", ""},
		{"preview:stop", "停止预览", "POST", "/api/v1/streams/*/preview/stop", ""},
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
	}

	for _, p := range permissions {
		existing, _ := repo.GetPermissionByCode(ctx, p.Code)
		if existing != nil {
			continue
		}
		perm := &domain.Permission{
			ID:          uuid.New(),
			Code:        p.Code,
			Name:        p.Name,
			Method:      p.Method,
			Path:        p.Path,
			Description: p.Description,
		}
		if err := repo.CreatePermission(ctx, perm); err != nil {
			log.Printf("创建权限失败 %s: %v", p.Code, err)
		}
	}

	return nil
}

func initMenus(ctx context.Context, repo *repository) error {
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
		{uuid.MustParse("00000000-0000-0000-0000-000000000010"), nil, "stream", "视频流管理", 2, "/streams", "Monitor", "StreamList", "stream:list", 1, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000020"), nil, "algorithm", "算法管理", 2, "/algorithms", "Cpu", "AlgorithmList", "algorithm:list", 2, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000030"), nil, "inference", "推理结果", 2, "/inference-results", "DataAnalysis", "InferenceResultList", "inference:list", 3, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000001"), nil, "system", "系统管理", 1, "/system", "Setting", "", "", 100, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000002"), ptrUUID("00000000-0000-0000-0000-000000000001"), "system:user", "用户管理", 2, "/system/user", "User", "system/user/index", "user:list", 1, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000003"), ptrUUID("00000000-0000-0000-0000-000000000001"), "system:role", "角色管理", 2, "/system/role", "UserFilled", "system/role/index", "role:list", 2, true},
		{uuid.MustParse("00000000-0000-0000-0000-000000000004"), ptrUUID("00000000-0000-0000-0000-000000000001"), "system:menu", "菜单管理", 2, "/system/menu", "Menu", "system/menu/index", "menu:list", 3, true},
	}

	for _, m := range menus {
		existing, _ := repo.GetMenuByCode(ctx, m.Code)
		if existing != nil {
			continue
		}
		menu := &domain.Menu{
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
			Status:     domain.MenuStatusEnabled,
		}
		if err := repo.CreateMenu(ctx, menu); err != nil {
			log.Printf("创建菜单失败 %s: %v", m.Code, err)
		}
	}

	return nil
}

func initRoles(ctx context.Context, repo *repository) error {
	existing, _ := repo.GetRoleByCode(ctx, "super_admin")
	if existing != nil {
		return nil
	}

	role := &domain.Role{
		ID:          uuid.MustParse("00000000-0000-0000-0000-000000000100"),
		Code:        "super_admin",
		Name:        "超级管理员",
		Description: "拥有所有权限",
		Status:      domain.RoleStatusEnabled,
	}

	if err := repo.CreateRole(ctx, role); err != nil {
		return err
	}

	permissions, _ := repo.ListPermissions(ctx)
	var permIDs []uuid.UUID
	for _, p := range permissions {
		permIDs = append(permIDs, p.ID)
	}
	if len(permIDs) > 0 {
		repo.SetRolePermissions(ctx, role.ID, permIDs)
	}

	menus, _ := repo.ListMenus(ctx, nil)
	var menuIDs []uuid.UUID
	for _, m := range menus {
		menuIDs = append(menuIDs, m.ID)
	}
	if len(menuIDs) > 0 {
		repo.SetRoleMenus(ctx, role.ID, menuIDs)
	}

	return nil
}

func initAdminUser(ctx context.Context, repo *repository) error {
	existing, _ := repo.GetUserByUsername(ctx, "admin")
	if existing != nil {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		ID:       uuid.MustParse("00000000-0000-0000-0000-000000000200"),
		Username: "admin",
		Password: string(hashedPassword),
		Nickname: "管理员",
		Status:   domain.UserStatusEnabled,
	}

	if err := repo.CreateUser(ctx, user); err != nil {
		return err
	}

	roleID := uuid.MustParse("00000000-0000-0000-0000-000000000100")
	return repo.SetUserRoles(ctx, user.ID, []uuid.UUID{roleID})
}

func ptrUUID(s string) *uuid.UUID {
	id := uuid.MustParse(s)
	return &id
}
