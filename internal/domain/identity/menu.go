package identity

import (
	"time"

	"github.com/google/uuid"
)

// MenuType 菜单类型
type MenuType int

const (
	MenuTypeDirectory MenuType = 1 // 目录
	MenuTypeMenu      MenuType = 2 // 菜单
	MenuTypeButton    MenuType = 3 // 按钮
)

// MenuStatus 菜单状态
type MenuStatus int

const (
	MenuStatusDisabled MenuStatus = 0 // 禁用
	MenuStatusEnabled  MenuStatus = 1 // 启用
)

// Menu 菜单实体（纯域模型，无 ORM 依赖）
type Menu struct {
	ID         uuid.UUID
	ParentID   *uuid.UUID
	Code       string
	Name       string
	Type       MenuType
	Path       string
	Icon       string
	Component  string
	Permission string
	Sort       int
	Visible    bool
	Status     MenuStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Children   []Menu // 子菜单（树形结构）
}

// NewMenu 创建新菜单
func NewMenu(code, name string, menuType MenuType) *Menu {
	return &Menu{
		ID:      uuid.New(),
		Code:    code,
		Name:    name,
		Type:    menuType,
		Visible: true,
		Status:  MenuStatusEnabled,
	}
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

// Enable 启用菜单
func (m *Menu) Enable() {
	m.Status = MenuStatusEnabled
}

// Disable 禁用菜单
func (m *Menu) Disable() {
	m.Status = MenuStatusDisabled
}

// SetParent 设置父菜单
func (m *Menu) SetParent(parentID uuid.UUID) {
	m.ParentID = &parentID
}

// IsTopLevel 是否为顶级菜单
func (m *Menu) IsTopLevel() bool {
	return m.ParentID == nil
}
