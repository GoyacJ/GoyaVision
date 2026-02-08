package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UserModel struct {
	ID        uuid.UUID   `gorm:"type:uuid;primaryKey"`
	Username  string      `gorm:"uniqueIndex;not null;size:64"`
	Password  string      `gorm:"not null;size:128"`
	Nickname  string      `gorm:"size:64"`
	Email     string      `gorm:"size:128"`
	Phone     string      `gorm:"size:32"`
	Avatar    string      `gorm:"size:256"`
	Status    int         `gorm:"default:1"`
	TenantID  *uuid.UUID  `gorm:"type:uuid;index:idx_users_tenant_id"`
	CreatedAt  time.Time           `gorm:"autoCreateTime"`
	UpdatedAt  time.Time           `gorm:"autoUpdateTime"`
	Roles      []RoleModel         `gorm:"many2many:user_roles"`
	Identities []UserIdentityModel `gorm:"foreignKey:UserID"`
}

func (UserModel) TableName() string { return "users" }

type RoleModel struct {
	ID          uuid.UUID         `gorm:"type:uuid;primaryKey"`
	Code        string            `gorm:"uniqueIndex;not null;size:64"`
	Name             string            `gorm:"not null;size:64"`
	Description      string            `gorm:"size:256"`
	Status           int               `gorm:"default:1"`
	IsDefault        bool              `gorm:"default:false;index:idx_roles_is_default"`
	AutoAssignConfig datatypes.JSON    `gorm:"serializer:json"`
	CreatedAt        time.Time         `gorm:"autoCreateTime"`
	UpdatedAt   time.Time         `gorm:"autoUpdateTime"`
	Permissions []PermissionModel `gorm:"many2many:role_permissions"`
	Menus       []MenuModel       `gorm:"many2many:role_menus"`
}

func (RoleModel) TableName() string { return "roles" }

type PermissionModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Code        string    `gorm:"uniqueIndex;not null;size:64"`
	Name        string    `gorm:"not null;size:64"`
	Method      string    `gorm:"size:16"`
	Path        string    `gorm:"size:256"`
	Description string    `gorm:"size:256"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (PermissionModel) TableName() string { return "permissions" }

type MenuModel struct {
	ID         uuid.UUID   `gorm:"type:uuid;primaryKey"`
	ParentID   *uuid.UUID  `gorm:"type:uuid;index"`
	Code       string      `gorm:"uniqueIndex;not null;size:64"`
	Name       string      `gorm:"not null;size:64"`
	Type       int         `gorm:"not null"`
	Path       string      `gorm:"size:256"`
	Icon       string      `gorm:"size:64"`
	Component  string      `gorm:"size:256"`
	Permission string      `gorm:"size:64"`
	Sort       int         `gorm:"default:0"`
	Visible    bool        `gorm:"default:true"`
	Status     int         `gorm:"default:1"`
	CreatedAt  time.Time   `gorm:"autoCreateTime"`
	UpdatedAt  time.Time   `gorm:"autoUpdateTime"`
	Children   []MenuModel `gorm:"-"`
}

func (MenuModel) TableName() string { return "menus" }
