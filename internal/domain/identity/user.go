package identity

import (
	"time"

	"github.com/google/uuid"
)

const (
	UserStatusDisabled = 0
	UserStatusEnabled  = 1
)

type User struct {
	ID        uuid.UUID
	TenantID  *uuid.UUID
	Username  string
	Password  string
	Nickname  string
	Email     string
	Phone     string
	Avatar    string
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Roles     []Role
}

func (u *User) IsEnabled() bool {
	return u.Status == UserStatusEnabled
}
