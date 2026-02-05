package dto

import (
	"time"

	"goyavision/internal/domain/identity"

	"github.com/google/uuid"
)

// UserCreateRequest 创建用户请求
type UserCreateRequest struct {
	Username string      `json:"username" validate:"required"`
	Password string      `json:"password" validate:"required,min=6"`
	Nickname string      `json:"nickname"`
	Email    string      `json:"email"`
	Phone    string      `json:"phone"`
	Avatar   string      `json:"avatar"`
	Status   int         `json:"status"`
	RoleIDs  []uuid.UUID `json:"role_ids"`
}

// UserUpdateRequest 更新用户请求
type UserUpdateRequest struct {
	Nickname *string     `json:"nickname"`
	Email    *string     `json:"email"`
	Phone    *string     `json:"phone"`
	Avatar   *string     `json:"avatar"`
	Status   *int        `json:"status"`
	Password *string     `json:"password"`
	RoleIDs  []uuid.UUID `json:"role_ids"`
}

// UserListQuery 用户列表查询参数
type UserListQuery struct {
	Status *int `query:"status"`
	Limit  int  `query:"limit"`
	Offset int  `query:"offset"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID        uuid.UUID    `json:"id"`
	Username  string       `json:"username"`
	Nickname  string       `json:"nickname"`
	Email     string       `json:"email"`
	Phone     string       `json:"phone"`
	Avatar    string       `json:"avatar"`
	Status    int          `json:"status"`
	Roles     []RoleSimple `json:"roles,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// RoleSimple 简化的角色信息
type RoleSimple struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Items []*UserResponse `json:"items"`
	Total int64           `json:"total"`
}

// UserToResponse 将 identity.User 转换为 UserResponse
func UserToResponse(u *identity.User) *UserResponse {
	if u == nil {
		return nil
	}
	roles := make([]RoleSimple, len(u.Roles))
	for i, r := range u.Roles {
		roles[i] = RoleSimple{
			ID:   r.ID,
			Code: r.Code,
			Name: r.Name,
		}
	}
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Nickname:  u.Nickname,
		Email:     u.Email,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		Status:    u.Status,
		Roles:     roles,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// UsersToResponse 批量转换
func UsersToResponse(users []*identity.User) []*UserResponse {
	result := make([]*UserResponse, len(users))
	for i, u := range users {
		result[i] = UserToResponse(u)
	}
	return result
}
