package mapper

import (
	"goyavision/internal/domain/identity"
	"goyavision/internal/infra/persistence/model"
)

func UserToModel(u *identity.User) *model.UserModel {
	m := &model.UserModel{
		ID:        u.ID,
		Username:  u.Username,
		Password:  u.Password,
		Nickname:  u.Nickname,
		Email:     u.Email,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	for _, r := range u.Roles {
		m.Roles = append(m.Roles, *RoleToModel(&r))
	}
	return m
}

func UserToDomain(m *model.UserModel) *identity.User {
	u := &identity.User{
		ID:        m.ID,
		Username:  m.Username,
		Password:  m.Password,
		Nickname:  m.Nickname,
		Email:     m.Email,
		Phone:     m.Phone,
		Avatar:    m.Avatar,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	for _, r := range m.Roles {
		u.Roles = append(u.Roles, *RoleToDomain(&r))
	}
	return u
}

func RoleToModel(r *identity.Role) *model.RoleModel {
	m := &model.RoleModel{
		ID:          r.ID,
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		Status:      r.Status,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
	for _, p := range r.Permissions {
		m.Permissions = append(m.Permissions, *PermissionToModel(&p))
	}
	for _, menu := range r.Menus {
		m.Menus = append(m.Menus, *MenuToModel(&menu))
	}
	return m
}

func RoleToDomain(m *model.RoleModel) *identity.Role {
	r := &identity.Role{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		Status:      m.Status,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	for _, p := range m.Permissions {
		r.Permissions = append(r.Permissions, *PermissionToDomain(&p))
	}
	for _, menu := range m.Menus {
		r.Menus = append(r.Menus, *MenuToDomain(&menu))
	}
	return r
}

func PermissionToModel(p *identity.Permission) *model.PermissionModel {
	return &model.PermissionModel{
		ID:          p.ID,
		Code:        p.Code,
		Name:        p.Name,
		Method:      p.Method,
		Path:        p.Path,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func PermissionToDomain(m *model.PermissionModel) *identity.Permission {
	return &identity.Permission{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Method:      m.Method,
		Path:        m.Path,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func MenuToModel(menu *identity.Menu) *model.MenuModel {
	return &model.MenuModel{
		ID:         menu.ID,
		ParentID:   menu.ParentID,
		Code:       menu.Code,
		Name:       menu.Name,
		Type:       int(menu.Type),
		Path:       menu.Path,
		Icon:       menu.Icon,
		Component:  menu.Component,
		Permission: menu.Permission,
		Sort:       menu.Sort,
		Visible:    menu.Visible,
		Status:     int(menu.Status),
		CreatedAt:  menu.CreatedAt,
		UpdatedAt:  menu.UpdatedAt,
	}
}

func MenuToDomain(m *model.MenuModel) *identity.Menu {
	return &identity.Menu{
		ID:         m.ID,
		ParentID:   m.ParentID,
		Code:       m.Code,
		Name:       m.Name,
		Type:       identity.MenuType(m.Type),
		Path:       m.Path,
		Icon:       m.Icon,
		Component:  m.Component,
		Permission: m.Permission,
		Sort:       m.Sort,
		Visible:    m.Visible,
		Status:     identity.MenuStatus(m.Status),
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}
