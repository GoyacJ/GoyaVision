package command

import (
	"context"

	"golang.org/x/crypto/bcrypt"
	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/identity"
	"goyavision/pkg/apperr"

	"github.com/google/uuid"
)

type LoginHandler struct {
	uow          port.UnitOfWork
	tokenService port.TokenService
}

func NewLoginHandler(uow port.UnitOfWork, ts port.TokenService) *LoginHandler {
	return &LoginHandler{uow: uow, tokenService: ts}
}

func (h *LoginHandler) Handle(ctx context.Context, cmd dto.LoginCommand) (*dto.LoginResult, error) {
	if cmd.Username == "" {
		return nil, apperr.InvalidInput("username is required")
	}
	if cmd.Password == "" {
		return nil, apperr.InvalidInput("password is required")
	}

	var user *identity.User
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		var err error
		user, err = repos.Users.GetByUsername(ctx, cmd.Username)
		return err
	})

	if err != nil {
		if apperr.IsNotFound(err) {
			return nil, apperr.Wrap(err, apperr.CodeLoginFailed, "invalid username or password")
		}
		return nil, apperr.Internal("get user by username", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cmd.Password)); err != nil {
		return nil, apperr.New(apperr.CodeLoginFailed, "invalid username or password")
	}

	if !user.IsEnabled() {
		return nil, apperr.New(apperr.CodeForbidden, "user is disabled")
	}

	tenantID := uuid.Nil
	if user.TenantID != nil {
		tenantID = *user.TenantID
	}

	tokens, err := h.tokenService.GenerateTokenPair(user.ID, tenantID, user.Username)
	if err != nil {
		return nil, apperr.Internal("generate token pair", err)
	}

	userInfo, err := h.getUserInfo(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResult{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		User:         userInfo,
	}, nil
}

func (h *LoginHandler) getUserInfo(ctx context.Context, userID uuid.UUID) (*dto.UserInfo, error) {
	var user *identity.User
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		var err error
		user, err = repos.Users.GetWithRoles(ctx, userID)
		return err
	})

	if err != nil {
		return nil, apperr.Internal("get user with roles", err)
	}

	var roleIDs []uuid.UUID
	var roleCodes []string
	isSuperAdmin := false

	for _, role := range user.Roles {
		if role.IsEnabled() {
			roleIDs = append(roleIDs, role.ID)
			roleCodes = append(roleCodes, role.Code)
			if role.Code == "super_admin" {
				isSuperAdmin = true
			}
		}
	}

	var permissions []string
	var menus []*identity.Menu

	if isSuperAdmin {
		permissions = []string{"*"}
		var allMenus []*identity.Menu
		err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
			var err error
			enabled := 1
			allMenus, err = repos.Menus.List(ctx, &enabled)
			return err
		})
		if err != nil {
			return nil, apperr.Internal("list all menus", err)
		}
		menus = buildMenuTree(allMenus)
	} else {
		var perms []*identity.Permission
		err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
			var err error
			perms, err = repos.Permissions.GetByRoleIDs(ctx, roleIDs)
			return err
		})
		if err != nil {
			return nil, apperr.Internal("get permissions by role ids", err)
		}
		for _, p := range perms {
			permissions = append(permissions, p.Code)
		}

		var menuList []*identity.Menu
		err = h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
			var err error
			menuList, err = repos.Menus.GetByRoleIDs(ctx, roleIDs)
			return err
		})
		if err != nil {
			return nil, apperr.Internal("get menus by role ids", err)
		}
		menus = buildMenuTree(menuList)
	}

	return &dto.UserInfo{
		ID:          user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Email:       user.Email,
		Phone:       user.Phone,
		Avatar:      user.Avatar,
		Roles:       roleCodes,
		Permissions: permissions,
		Menus:       menus,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func buildMenuTree(menus []*identity.Menu) []*identity.Menu {
	menuMap := make(map[uuid.UUID]*identity.Menu)
	var roots []*identity.Menu

	for _, m := range menus {
		menuCopy := *m
		menuCopy.Children = []identity.Menu{}
		menuMap[m.ID] = &menuCopy
	}

	for _, m := range menus {
		menu := menuMap[m.ID]
		if m.ParentID == nil {
			roots = append(roots, menu)
		} else {
			if parent, ok := menuMap[*m.ParentID]; ok {
				parent.Children = append(parent.Children, *menu)
			} else {
				roots = append(roots, menu)
			}
		}
	}

	return roots
}
