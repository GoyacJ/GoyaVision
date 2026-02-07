package dto

import (
	"time"

	"goyavision/internal/domain/identity"

	"github.com/google/uuid"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginOAuthRequest OAuth登录请求
type LoginOAuthRequest struct {
	Provider string `json:"provider" validate:"required"`
	Code     string `json:"code" validate:"required"`
	State    string `json:"state"`
}

// BindIdentityRequest 绑定身份请求
type BindIdentityRequest struct {
	Provider   string                 `json:"provider" validate:"required"`
	Identifier string                 `json:"identifier" validate:"required"`
	Credential string                 `json:"credential"`
	Meta       map[string]interface{} `json:"meta"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	User         *UserInfo `json:"user,omitempty"`
}

// RefreshTokenRequest 刷新 Token 请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID          uuid.UUID     `json:"id"`
	Username    string        `json:"username"`
	Nickname    string        `json:"nickname"`
	Email       string        `json:"email"`
	Phone       string        `json:"phone"`
	Avatar      string        `json:"avatar"`
	Roles       []string      `json:"roles"`
	Permissions []string      `json:"permissions"`
	Menus       []*MenuSimple `json:"menus"`
}

// MenuSimple 简化的菜单信息
type MenuSimple struct {
	ID         uuid.UUID     `json:"id"`
	ParentID   *uuid.UUID    `json:"parent_id,omitempty"`
	Code       string        `json:"code"`
	Name       string        `json:"name"`
	Type       int           `json:"type"`
	Path       string        `json:"path"`
	Icon       string        `json:"icon"`
	Component  string        `json:"component"`
	Permission string        `json:"permission"`
	Sort       int           `json:"sort"`
	Visible    bool          `json:"visible"`
	Children   []*MenuSimple `json:"children,omitempty"`
}

// AppUserInfo 用于从 app 层转换的结构
type AppUserInfo struct {
	ID          uuid.UUID
	Username    string
	Nickname    string
	Email       string
	Phone       string
	Avatar      string
	Roles       []string
	Permissions []string
	Menus       []*identity.Menu
}

// UserInfoFromApp 从 AppUserInfo 转换
func UserInfoFromApp(u *AppUserInfo) *UserInfo {
	if u == nil {
		return nil
	}
	return &UserInfo{
		ID:          u.ID,
		Username:    u.Username,
		Nickname:    u.Nickname,
		Email:       u.Email,
		Phone:       u.Phone,
		Avatar:      u.Avatar,
		Roles:       u.Roles,
		Permissions: u.Permissions,
		Menus:       menusToSimple(u.Menus),
	}
}

func menusToSimple(menus []*identity.Menu) []*MenuSimple {
	if menus == nil {
		return nil
	}
	result := make([]*MenuSimple, len(menus))
	for i, m := range menus {
		result[i] = menuToSimple(m)
	}
	return result
}

func menuToSimple(m *identity.Menu) *MenuSimple {
	if m == nil {
		return nil
	}
	children := make([]*MenuSimple, len(m.Children))
	for i, c := range m.Children {
		child := c
		children[i] = menuToSimple(&child)
	}
	return &MenuSimple{
		ID:         m.ID,
		ParentID:   m.ParentID,
		Code:       m.Code,
		Name:       m.Name,
		Type:       int(m.Type),
		Path:       m.Path,
		Icon:       m.Icon,
		Component:  m.Component,
		Permission: m.Permission,
		Sort:       m.Sort,
		Visible:    m.Visible,
		Children:   children,
	}
}

// ProfileResponse 用户信息响应
type ProfileResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}
