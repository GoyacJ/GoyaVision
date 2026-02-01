package dto

import (
	"time"

	"goyavision/internal/domain"

	"github.com/google/uuid"
)

// RoleCreateRequest 创建角色请求
type RoleCreateRequest struct {
	Code          string      `json:"code" validate:"required"`
	Name          string      `json:"name" validate:"required"`
	Description   string      `json:"description"`
	Status        int         `json:"status"`
	PermissionIDs []uuid.UUID `json:"permission_ids"`
	MenuIDs       []uuid.UUID `json:"menu_ids"`
}

// RoleUpdateRequest 更新角色请求
type RoleUpdateRequest struct {
	Name          *string     `json:"name"`
	Description   *string     `json:"description"`
	Status        *int        `json:"status"`
	PermissionIDs []uuid.UUID `json:"permission_ids"`
	MenuIDs       []uuid.UUID `json:"menu_ids"`
}

// RoleListQuery 角色列表查询参数
type RoleListQuery struct {
	Status *int `query:"status"`
}

// RoleResponse 角色响应
type RoleResponse struct {
	ID          uuid.UUID          `json:"id"`
	Code        string             `json:"code"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Status      int                `json:"status"`
	Permissions []PermissionSimple `json:"permissions,omitempty"`
	Menus       []MenuSimple       `json:"menus,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// PermissionSimple 简化的权限信息
type PermissionSimple struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

// RoleToResponse 将 domain.Role 转换为 RoleResponse
func RoleToResponse(r *domain.Role) *RoleResponse {
	if r == nil {
		return nil
	}
	permissions := make([]PermissionSimple, len(r.Permissions))
	for i, p := range r.Permissions {
		permissions[i] = PermissionSimple{
			ID:   p.ID,
			Code: p.Code,
			Name: p.Name,
		}
	}
	menus := make([]MenuSimple, len(r.Menus))
	for i, m := range r.Menus {
		menus[i] = MenuSimple{
			ID:   m.ID,
			Code: m.Code,
			Name: m.Name,
		}
	}
	return &RoleResponse{
		ID:          r.ID,
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		Status:      r.Status,
		Permissions: permissions,
		Menus:       menus,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// RolesToResponse 批量转换
func RolesToResponse(roles []*domain.Role) []*RoleResponse {
	result := make([]*RoleResponse, len(roles))
	for i, r := range roles {
		result[i] = RoleToResponse(r)
	}
	return result
}
