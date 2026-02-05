package dto

import (
	"time"

	"goyavision/internal/domain/identity"

	"github.com/google/uuid"
)

// MenuCreateRequest 创建菜单请求
type MenuCreateRequest struct {
	ParentID   *uuid.UUID `json:"parent_id"`
	Code       string     `json:"code" validate:"required"`
	Name       string     `json:"name" validate:"required"`
	Type       int        `json:"type" validate:"required"`
	Path       string     `json:"path"`
	Icon       string     `json:"icon"`
	Component  string     `json:"component"`
	Permission string     `json:"permission"`
	Sort       int        `json:"sort"`
	Visible    bool       `json:"visible"`
	Status     int        `json:"status"`
}

// MenuUpdateRequest 更新菜单请求
type MenuUpdateRequest struct {
	ParentID   *uuid.UUID `json:"parent_id"`
	Name       *string    `json:"name"`
	Type       *int       `json:"type"`
	Path       *string    `json:"path"`
	Icon       *string    `json:"icon"`
	Component  *string    `json:"component"`
	Permission *string    `json:"permission"`
	Sort       *int       `json:"sort"`
	Visible    *bool      `json:"visible"`
	Status     *int       `json:"status"`
}

// MenuListQuery 菜单列表查询参数
type MenuListQuery struct {
	Status *int `query:"status"`
}

// MenuResponse 菜单响应
type MenuResponse struct {
	ID         uuid.UUID       `json:"id"`
	ParentID   *uuid.UUID      `json:"parent_id,omitempty"`
	Code       string          `json:"code"`
	Name       string          `json:"name"`
	Type       int             `json:"type"`
	Path       string          `json:"path"`
	Icon       string          `json:"icon"`
	Component  string          `json:"component"`
	Permission string          `json:"permission"`
	Sort       int             `json:"sort"`
	Visible    bool            `json:"visible"`
	Status     int             `json:"status"`
	Children   []*MenuResponse `json:"children,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

// MenuToResponse 将 identity.Menu 转换为 MenuResponse
func MenuToResponse(m *identity.Menu) *MenuResponse {
	if m == nil {
		return nil
	}
	children := make([]*MenuResponse, len(m.Children))
	for i, c := range m.Children {
		child := c
		children[i] = MenuToResponse(&child)
	}
	return &MenuResponse{
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
		Status:     int(m.Status),
		Children:   children,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

// MenusToResponse 批量转换
func MenusToResponse(menus []*identity.Menu) []*MenuResponse {
	result := make([]*MenuResponse, len(menus))
	for i, m := range menus {
		result[i] = MenuToResponse(m)
	}
	return result
}

// PermissionResponse 权限响应
type PermissionResponse struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PermissionToResponse 将 identity.Permission 转换为 PermissionResponse
func PermissionToResponse(p *identity.Permission) *PermissionResponse {
	if p == nil {
		return nil
	}
	return &PermissionResponse{
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

// PermissionsToResponse 批量转换
func PermissionsToResponse(permissions []*identity.Permission) []*PermissionResponse {
	result := make([]*PermissionResponse, len(permissions))
	for i, p := range permissions {
		result[i] = PermissionToResponse(p)
	}
	return result
}
