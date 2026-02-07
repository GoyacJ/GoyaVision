package identity

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoleStatusDisabled = 0
	RoleStatusEnabled  = 1
)

type Role struct {
	ID          uuid.UUID
	Code        string
	Name             string
	Description      string
	Status           int
	IsDefault        bool
	AutoAssignConfig *AutoAssignConfig
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Permissions []Permission
	Menus       []Menu
}

type AutoAssignConfig struct {
	Trigger    string            `json:"trigger"`    // register, bind, login
	Conditions map[string]string `json:"conditions"` // provider: github
}

func (r *Role) IsEnabled() bool {
	return r.Status == RoleStatusEnabled
}
