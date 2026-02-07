package command

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/identity"
	"goyavision/pkg/apperr"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BindIdentityHandler struct {
	uow          port.UnitOfWork
	tokenService port.TokenService
}

func NewBindIdentityHandler(uow port.UnitOfWork, ts port.TokenService) *BindIdentityHandler {
	return &BindIdentityHandler{uow: uow, tokenService: ts}
}

func (h *BindIdentityHandler) Handle(ctx context.Context, cmd dto.BindIdentityCommand) error {
	return h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		// 1. Check if identity already exists
		existing, err := repos.UserIdentities.GetByIdentifier(ctx, identity.IdentityType(cmd.Provider), cmd.Identifier)
		if err == nil && existing != nil {
			return apperr.Conflict("identity already bound to a user")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to check identity")
		}

		// 2. Create Identity
		identity := &identity.UserIdentity{
			ID:           uuid.New(),
			UserID:       cmd.UserID,
			IdentityType: identity.IdentityType(cmd.Provider),
			Identifier:   cmd.Identifier,
			Credential:   cmd.Credential,
			Meta:         cmd.Meta,
		}
		if err := repos.UserIdentities.Create(ctx, identity); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to create identity")
		}

		// 3. Assign Context-Aware Roles (trigger="bind")
		roles, err := repos.Roles.List(ctx, nil)
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to list roles")
		}

		var rolesToAdd []uuid.UUID
		for _, role := range roles {
			if role.AutoAssignConfig != nil {
				if role.AutoAssignConfig.Trigger == "bind" {
					if p, ok := role.AutoAssignConfig.Conditions["provider"]; ok && p == cmd.Provider {
						rolesToAdd = append(rolesToAdd, role.ID)
					}
				}
			}
		}

		if len(rolesToAdd) > 0 {
			user, err := repos.Users.GetWithRoles(ctx, cmd.UserID)
			if err != nil {
				return apperr.Wrap(err, apperr.CodeDBError, "failed to get user")
			}

			existingMap := make(map[uuid.UUID]bool)
			for _, r := range user.Roles {
				existingMap[r.ID] = true
			}

			var finalRoles []uuid.UUID
			for _, r := range user.Roles {
				finalRoles = append(finalRoles, r.ID)
			}

			changed := false
			for _, id := range rolesToAdd {
				if !existingMap[id] {
					finalRoles = append(finalRoles, id)
					changed = true
				}
			}

			if changed {
				if err := repos.Users.SetUserRoles(ctx, cmd.UserID, finalRoles); err != nil {
					return apperr.Wrap(err, apperr.CodeDBError, "failed to update user roles")
				}
			}
		}

		return nil
	})
}
