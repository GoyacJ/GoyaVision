package identity

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	Get(ctx context.Context, id uuid.UUID) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetWithRoles(ctx context.Context, id uuid.UUID) (*User, error)
	List(ctx context.Context, status *int, limit, offset int) ([]*User, int64, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetUserRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error
}

type RoleRepository interface {
	Create(ctx context.Context, r *Role) error
	Get(ctx context.Context, id uuid.UUID) (*Role, error)
	GetByCode(ctx context.Context, code string) (*Role, error)
	GetWithPermissions(ctx context.Context, id uuid.UUID) (*Role, error)
	List(ctx context.Context, status *int) ([]*Role, error)
	Update(ctx context.Context, r *Role) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetPermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
	SetMenus(ctx context.Context, roleID uuid.UUID, menuIDs []uuid.UUID) error
	GetDefaultRoles(ctx context.Context) ([]*Role, error)
}

type UserIdentityRepository interface {
	Create(ctx context.Context, i *UserIdentity) error
	Get(ctx context.Context, id uuid.UUID) (*UserIdentity, error)
	GetByIdentifier(ctx context.Context, identityType IdentityType, identifier string) (*UserIdentity, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*UserIdentity, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type PermissionRepository interface {
	Create(ctx context.Context, p *Permission) error
	Get(ctx context.Context, id uuid.UUID) (*Permission, error)
	GetByCode(ctx context.Context, code string) (*Permission, error)
	List(ctx context.Context) ([]*Permission, error)
	Update(ctx context.Context, p *Permission) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]*Permission, error)
}

type MenuRepository interface {
	Create(ctx context.Context, m *Menu) error
	Get(ctx context.Context, id uuid.UUID) (*Menu, error)
	GetByCode(ctx context.Context, code string) (*Menu, error)
	List(ctx context.Context, status *int) ([]*Menu, error)
	Update(ctx context.Context, m *Menu) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]*Menu, error)
}
