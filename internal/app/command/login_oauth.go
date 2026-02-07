package command

import (
	"context"
	"errors"
	"fmt"

	"goyavision/internal/app"
	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/identity"
	"goyavision/pkg/apperr"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoginOAuthHandler struct {
	uow             port.UnitOfWork
	tokenService    port.TokenService
	providerFactory port.AuthProviderFactory
	userService     *app.UserService
}

func NewLoginOAuthHandler(uow port.UnitOfWork, ts port.TokenService, pf port.AuthProviderFactory, us *app.UserService) *LoginOAuthHandler {
	return &LoginOAuthHandler{uow: uow, tokenService: ts, providerFactory: pf, userService: us}
}

func (h *LoginOAuthHandler) Handle(ctx context.Context, cmd dto.LoginOAuthCommand) (*dto.LoginResult, error) {
	provider, err := h.providerFactory.Get(cmd.Provider)
	if err != nil {
		return nil, apperr.InvalidInput("unsupported auth provider: " + cmd.Provider)
	}

	userInfo, err := provider.VerifyCode(ctx, cmd.Code)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeUnauthorized, "oauth verification failed")
	}

	var user *identity.User
	var isNewUser bool

	err = h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		// 1. Check if identity exists
		userIdentity, err := repos.UserIdentities.GetByIdentifier(ctx, identity.IdentityType(cmd.Provider), userInfo.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to check user identity")
		}

		if userIdentity != nil {
			// Found identity, get user
			user, err = repos.Users.GetWithRoles(ctx, userIdentity.UserID)
			if err != nil {
				return apperr.Wrap(err, apperr.CodeDBError, "failed to get user")
			}
			return nil
		}

		// 2. Identity not found, create user and identity
		isNewUser = true
		// Create User (Auto generate username/password)
		username := fmt.Sprintf("%s_%s", cmd.Provider, userInfo.ID)
		// Check username
		_, err = repos.Users.GetByUsername(ctx, username)
		if err == nil {
			username = fmt.Sprintf("%s_%s_%s", cmd.Provider, userInfo.ID, uuid.New().String()[:8])
		}

		// Use manual creation to ensure transaction consistency (UserService uses its own transaction logic usually, or just repos)
		// But here we are inside a UoW transaction, so we use repos directly.
		
		user = &identity.User{
			ID:       uuid.New(),
			Username: username,
			Password: "", // No password for oauth user initially
			Nickname: userInfo.Name,
			Email:    userInfo.Email,
			Avatar:   userInfo.Avatar,
			Status:   1,
		}
		
		if err := repos.Users.Create(ctx, user); err != nil {
			return err
		}

		// Assign Default Roles
		defaultRoles, err := repos.Roles.GetDefaultRoles(ctx)
		if err == nil && len(defaultRoles) > 0 {
			roleIDs := make([]uuid.UUID, len(defaultRoles))
			for i, r := range defaultRoles {
				roleIDs[i] = r.ID
			}
			if err := repos.Users.SetUserRoles(ctx, user.ID, roleIDs); err != nil {
				return err
			}
		}

		// Create Identity
		userIdentityModel := &identity.UserIdentity{
			ID:           uuid.New(),
			UserID:       user.ID,
			IdentityType: identity.IdentityType(cmd.Provider),
			Identifier:   userInfo.ID,
			Meta:         userInfo.Raw,
		}
		if err := repos.UserIdentities.Create(ctx, userIdentityModel); err != nil {
			return err
		}
		
		// Reload user with roles for token generation context
		user, err = repos.Users.GetWithRoles(ctx, user.ID)
		return err
	})

	if err != nil {
		return nil, err
	}

	if !user.IsEnabled() {
		return nil, apperr.Forbidden("user is disabled")
	}

	// 3. Assign Context-Aware Roles (AutoAssignConfig)
	
	err = h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		roles, err := repos.Roles.List(ctx, nil)
		if err != nil {
			return err
		}

		var rolesToAdd []uuid.UUID
		for _, role := range roles {
			if role.AutoAssignConfig != nil {
				// Check trigger
				trigger := role.AutoAssignConfig.Trigger
				if trigger == "login" || (isNewUser && trigger == "register") || trigger == "bind" {
					// Check conditions
					if p, ok := role.AutoAssignConfig.Conditions["provider"]; ok && p == cmd.Provider {
						rolesToAdd = append(rolesToAdd, role.ID)
					}
				}
			}
		}

		if len(rolesToAdd) > 0 {
			// Reload user roles to avoid duplicates
			currentUser, err := repos.Users.GetWithRoles(ctx, user.ID)
			if err != nil { return err }
			
			existingMap := make(map[uuid.UUID]bool)
			for _, r := range currentUser.Roles {
				existingMap[r.ID] = true
			}
			
			var finalRoles []uuid.UUID
			for _, r := range currentUser.Roles {
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
				if err := repos.Users.SetUserRoles(ctx, user.ID, finalRoles); err != nil {
					return err
				}
			}
		}
		return nil
	})
	
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeInternal, "failed to assign conditional roles")
	}

	// Generate Tokens
	tokens, err := h.tokenService.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		return nil, apperr.Internal("generate token pair", err)
	}

	return &dto.LoginResult{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	}, nil
}
