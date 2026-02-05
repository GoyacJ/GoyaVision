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
	Name        string
	Description string
	Status      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Permissions []Permission
	Menus       []Menu
}

func (r *Role) IsEnabled() bool {
	return r.Status == RoleStatusEnabled
}
