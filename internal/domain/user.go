package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserStatus 用户状态
const (
	UserStatusDisabled = 0
	UserStatusEnabled  = 1
)

// User 用户实体
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username  string    `gorm:"uniqueIndex;not null;size:64"`
	Password  string    `gorm:"not null;size:128"`
	Nickname  string    `gorm:"size:64"`
	Email     string    `gorm:"size:128"`
	Phone     string    `gorm:"size:32"`
	Avatar    string    `gorm:"size:256"`
	Status    int       `gorm:"default:1"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Roles     []Role    `gorm:"many2many:user_roles"`
}

func (User) TableName() string {
	return "users"
}

// IsEnabled 检查用户是否启用
func (u *User) IsEnabled() bool {
	return u.Status == UserStatusEnabled
}
