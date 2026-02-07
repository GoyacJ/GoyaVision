package app

import (
	"context"
	"errors"

	"goyavision/internal/domain/identity"
	"goyavision/internal/port"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUsernameExists   = errors.New("username already exists")
	ErrRoleNotFound     = errors.New("role not found")
	ErrRoleCodeExists   = errors.New("role code already exists")
	ErrMenuNotFound     = errors.New("menu not found")
	ErrMenuCodeExists   = errors.New("menu code already exists")
	ErrPermissionExists = errors.New("permission code already exists")
)

// UserService 用户管理服务
type UserService struct {
	repo port.Repository
}

func NewUserService(repo port.Repository) *UserService {
	return &UserService{repo: repo}
}

type CreateUserRequest struct {
	Username string
	Password string
	Nickname string
	Email    string
	Phone    string
	Avatar   string
	Status   int
	RoleIDs  []uuid.UUID
}

func (s *UserService) Create(ctx context.Context, req *CreateUserRequest) (*identity.User, error) {
	existing, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err == nil && existing != nil {
		return nil, ErrUsernameExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &identity.User{
		ID:       uuid.New(),
		Username: req.Username,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Status:   req.Status,
	}

	if user.Status == 0 {
		user.Status = identity.UserStatusEnabled
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	roleIDs := req.RoleIDs
	if len(roleIDs) == 0 {
		// 分配默认角色
		defaultRoles, err := s.repo.GetDefaultRoles(ctx)
		if err == nil && len(defaultRoles) > 0 {
			for _, r := range defaultRoles {
				roleIDs = append(roleIDs, r.ID)
			}
		}
	}

	if len(roleIDs) > 0 {
		if err := s.repo.SetUserRoles(ctx, user.ID, roleIDs); err != nil {
			return nil, err
		}
	}

	return s.repo.GetUserWithRoles(ctx, user.ID)
}

func (s *UserService) Get(ctx context.Context, id uuid.UUID) (*identity.User, error) {
	return s.repo.GetUserWithRoles(ctx, id)
}

func (s *UserService) List(ctx context.Context, status *int, limit, offset int) ([]*identity.User, int64, error) {
	return s.repo.ListUsers(ctx, status, limit, offset)
}

type UpdateUserRequest struct {
	Nickname *string
	Email    *string
	Phone    *string
	Avatar   *string
	Status   *int
	Password *string
	RoleIDs  []uuid.UUID
}

func (s *UserService) Update(ctx context.Context, id uuid.UUID, req *UpdateUserRequest) (*identity.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if req.Nickname != nil {
		user.Nickname = *req.Nickname
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	if req.Password != nil && *req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	if req.RoleIDs != nil {
		if err := s.repo.SetUserRoles(ctx, id, req.RoleIDs); err != nil {
			return nil, err
		}
	}

	return s.repo.GetUserWithRoles(ctx, id)
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return s.repo.DeleteUser(ctx, id)
}

func (s *UserService) ResetPassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.repo.UpdateUser(ctx, user)
}

// RoleService 角色管理服务
type RoleService struct {
	repo port.Repository
}

func NewRoleService(repo port.Repository) *RoleService {
	return &RoleService{repo: repo}
}

type CreateRoleRequest struct {
	Code             string
	Name             string
	Description      string
	Status           int
	IsDefault        bool
	AutoAssignConfig *identity.AutoAssignConfig
	PermissionIDs    []uuid.UUID
	MenuIDs          []uuid.UUID
}

func (s *RoleService) Create(ctx context.Context, req *CreateRoleRequest) (*identity.Role, error) {
	existing, err := s.repo.GetRoleByCode(ctx, req.Code)
	if err == nil && existing != nil {
		return nil, ErrRoleCodeExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	role := &identity.Role{
		ID:          uuid.New(),
		Code:        req.Code,
		Name:             req.Name,
		Description:      req.Description,
		Status:           req.Status,
		IsDefault:        req.IsDefault,
		AutoAssignConfig: req.AutoAssignConfig,
	}

	if role.Status == 0 {
		role.Status = identity.RoleStatusEnabled
	}

	if err := s.repo.CreateRole(ctx, role); err != nil {
		return nil, err
	}

	if len(req.PermissionIDs) > 0 {
		if err := s.repo.SetRolePermissions(ctx, role.ID, req.PermissionIDs); err != nil {
			return nil, err
		}
	}

	if len(req.MenuIDs) > 0 {
		if err := s.repo.SetRoleMenus(ctx, role.ID, req.MenuIDs); err != nil {
			return nil, err
		}
	}

	return s.repo.GetRoleWithPermissions(ctx, role.ID)
}

func (s *RoleService) Get(ctx context.Context, id uuid.UUID) (*identity.Role, error) {
	return s.repo.GetRoleWithPermissions(ctx, id)
}

func (s *RoleService) List(ctx context.Context, status *int) ([]*identity.Role, error) {
	return s.repo.ListRoles(ctx, status)
}

type UpdateRoleRequest struct {
	Name             *string
	Description      *string
	Status           *int
	IsDefault        *bool
	AutoAssignConfig *identity.AutoAssignConfig
	PermissionIDs    []uuid.UUID
	MenuIDs          []uuid.UUID
}

func (s *RoleService) Update(ctx context.Context, id uuid.UUID, req *UpdateRoleRequest) (*identity.Role, error) {
	role, err := s.repo.GetRole(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}

	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Description != nil {
		role.Description = *req.Description
	}
	if req.Status != nil {
		role.Status = *req.Status
	}
	if req.IsDefault != nil {
		role.IsDefault = *req.IsDefault
	}
	if req.AutoAssignConfig != nil {
		role.AutoAssignConfig = req.AutoAssignConfig
	}

	if err := s.repo.UpdateRole(ctx, role); err != nil {
		return nil, err
	}

	if req.PermissionIDs != nil {
		if err := s.repo.SetRolePermissions(ctx, id, req.PermissionIDs); err != nil {
			return nil, err
		}
	}

	if req.MenuIDs != nil {
		if err := s.repo.SetRoleMenus(ctx, id, req.MenuIDs); err != nil {
			return nil, err
		}
	}

	return s.repo.GetRoleWithPermissions(ctx, id)
}

func (s *RoleService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetRole(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
		return err
	}
	return s.repo.DeleteRole(ctx, id)
}

// MenuService 菜单管理服务
type MenuService struct {
	repo port.Repository
}

func NewMenuService(repo port.Repository) *MenuService {
	return &MenuService{repo: repo}
}

type CreateMenuRequest struct {
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
	Status     int
}

func (s *MenuService) Create(ctx context.Context, req *CreateMenuRequest) (*identity.Menu, error) {
	existing, err := s.repo.GetMenuByCode(ctx, req.Code)
	if err == nil && existing != nil {
		return nil, ErrMenuCodeExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	menu := &identity.Menu{
		ID:         uuid.New(),
		ParentID:   req.ParentID,
		Code:       req.Code,
		Name:       req.Name,
		Type:       identity.MenuType(req.Type),
		Path:       req.Path,
		Icon:       req.Icon,
		Component:  req.Component,
		Permission: req.Permission,
		Sort:       req.Sort,
		Visible:    req.Visible,
		Status:     identity.MenuStatus(req.Status),
	}

	if menu.Status == 0 {
		menu.Status = identity.MenuStatusEnabled
	}

	if err := s.repo.CreateMenu(ctx, menu); err != nil {
		return nil, err
	}

	return menu, nil
}

func (s *MenuService) Get(ctx context.Context, id uuid.UUID) (*identity.Menu, error) {
	return s.repo.GetMenu(ctx, id)
}

func (s *MenuService) List(ctx context.Context, status *int) ([]*identity.Menu, error) {
	return s.repo.ListMenus(ctx, status)
}

func (s *MenuService) ListTree(ctx context.Context, status *int) ([]*identity.Menu, error) {
	menus, err := s.repo.ListMenus(ctx, status)
	if err != nil {
		return nil, err
	}
	return buildMenuTree(menus), nil
}

type UpdateMenuRequest struct {
	ParentID   *uuid.UUID
	Name       *string
	Type       *int
	Path       *string
	Icon       *string
	Component  *string
	Permission *string
	Sort       *int
	Visible    *bool
	Status     *int
}

func (s *MenuService) Update(ctx context.Context, id uuid.UUID, req *UpdateMenuRequest) (*identity.Menu, error) {
	menu, err := s.repo.GetMenu(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMenuNotFound
		}
		return nil, err
	}

	if req.ParentID != nil {
		menu.ParentID = req.ParentID
	}
	if req.Name != nil {
		menu.Name = *req.Name
	}
	if req.Type != nil {
		menu.Type = identity.MenuType(*req.Type)
	}
	if req.Path != nil {
		menu.Path = *req.Path
	}
	if req.Icon != nil {
		menu.Icon = *req.Icon
	}
	if req.Component != nil {
		menu.Component = *req.Component
	}
	if req.Permission != nil {
		menu.Permission = *req.Permission
	}
	if req.Sort != nil {
		menu.Sort = *req.Sort
	}
	if req.Visible != nil {
		menu.Visible = *req.Visible
	}
	if req.Status != nil {
		menu.Status = identity.MenuStatus(*req.Status)
	}

	if err := s.repo.UpdateMenu(ctx, menu); err != nil {
		return nil, err
	}

	return menu, nil
}

func (s *MenuService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetMenu(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMenuNotFound
		}
		return err
	}
	return s.repo.DeleteMenu(ctx, id)
}

func buildMenuTree(menus []*identity.Menu) []*identity.Menu {
	menuMap := make(map[uuid.UUID]*identity.Menu)
	var roots []*identity.Menu

	for _, m := range menus {
		menuCopy := *m
		menuCopy.Children = []identity.Menu{}
		menuMap[m.ID] = &menuCopy
	}

	for _, m := range menus {
		menu := menuMap[m.ID]
		if m.ParentID == nil {
			roots = append(roots, menu)
		} else {
			if parent, ok := menuMap[*m.ParentID]; ok {
				parent.Children = append(parent.Children, *menu)
			} else {
				roots = append(roots, menu)
			}
		}
	}

	return roots
}

// PermissionService 权限管理服务
type PermissionService struct {
	repo port.Repository
}

func NewPermissionService(repo port.Repository) *PermissionService {
	return &PermissionService{repo: repo}
}

func (s *PermissionService) List(ctx context.Context) ([]*identity.Permission, error) {
	return s.repo.ListPermissions(ctx)
}

func (s *PermissionService) Get(ctx context.Context, id uuid.UUID) (*identity.Permission, error) {
	return s.repo.GetPermission(ctx, id)
}

type CreatePermissionRequest struct {
	Code        string
	Name        string
	Method      string
	Path        string
	Description string
}

func (s *PermissionService) Create(ctx context.Context, req *CreatePermissionRequest) (*identity.Permission, error) {
	existing, err := s.repo.GetPermissionByCode(ctx, req.Code)
	if err == nil && existing != nil {
		return nil, ErrPermissionExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	permission := &identity.Permission{
		ID:          uuid.New(),
		Code:        req.Code,
		Name:        req.Name,
		Method:      req.Method,
		Path:        req.Path,
		Description: req.Description,
	}

	if err := s.repo.CreatePermission(ctx, permission); err != nil {
		return nil, err
	}

	return permission, nil
}

func (s *PermissionService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeletePermission(ctx, id)
}
