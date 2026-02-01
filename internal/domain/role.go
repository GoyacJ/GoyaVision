package domain

import (
	"time"

	"github.com/google/uuid"
)

// RoleStatus 角色状态
const (
	RoleStatusDisabled = 0
	RoleStatusEnabled  = 1
)

// Role 角色实体
type Role struct {
	ID          uuid.UUID    `gorm:"type:uuid;primaryKey"`
	Code        string       `gorm:"uniqueIndex;not null;size:64"`
	Name        string       `gorm:"not null;size:64"`
	Description string       `gorm:"size:256"`
	Status      int          `gorm:"default:1"`
	CreatedAt   time.Time    `gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime"`
	Permissions []Permission `gorm:"many2many:role_permissions"`
	Menus       []Menu       `gorm:"many2many:role_menus"`
}

func (Role) TableName() string {
	return "roles"
}

// IsEnabled 检查角色是否启用
func (r *Role) IsEnabled() bool {
	return r.Status == RoleStatusEnabled
}
