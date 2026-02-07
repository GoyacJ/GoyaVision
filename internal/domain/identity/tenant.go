package identity

import (
	"time"

	"github.com/google/uuid"
)

const (
	TenantStatusDisabled = 0
	TenantStatusEnabled  = 1
)

type Tenant struct {
	ID        uuid.UUID
	Name      string
	Code      string
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTenant(name, code string) *Tenant {
	return &Tenant{
		ID:        uuid.New(),
		Name:      name,
		Code:      code,
		Status:    TenantStatusEnabled,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (t *Tenant) IsEnabled() bool {
	return t.Status == TenantStatusEnabled
}
