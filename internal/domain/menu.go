package domain

import (
	"time"

	"github.com/google/uuid"
)

// MenuType 菜单类型
const (
	MenuTypeDirectory = 1
	MenuTypeMenu      = 2
	MenuTypeButton    = 3
)

// MenuStatus 菜单状态
const (
	MenuStatusDisabled = 0
	MenuStatusEnabled  = 1
)

// Menu 菜单实体
type Menu struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey"`
	ParentID   *uuid.UUID `gorm:"type:uuid;index"`
	Code       string     `gorm:"uniqueIndex;not null;size:64"`
	Name       string     `gorm:"not null;size:64"`
	Type       int        `gorm:"not null"`
	Path       string     `gorm:"size:256"`
	Icon       string     `gorm:"size:64"`
	Component  string     `gorm:"size:256"`
	Permission string     `gorm:"size:64"`
	Sort       int        `gorm:"default:0"`
	Visible    bool       `gorm:"default:true"`
	Status     int        `gorm:"default:1"`
	CreatedAt  time.Time  `gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime"`
	Children   []Menu     `gorm:"-"`
}

func (Menu) TableName() string {
	return "menus"
}

// IsEnabled 检查菜单是否启用
func (m *Menu) IsEnabled() bool {
	return m.Status == MenuStatusEnabled
}

// IsDirectory 是否为目录
func (m *Menu) IsDirectory() bool {
	return m.Type == MenuTypeDirectory
}

// IsMenu 是否为菜单
func (m *Menu) IsMenu() bool {
	return m.Type == MenuTypeMenu
}

// IsButton 是否为按钮
func (m *Menu) IsButton() bool {
	return m.Type == MenuTypeButton
}
