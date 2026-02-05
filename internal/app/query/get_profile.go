package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/identity"
	"goyavision/pkg/apperr"

	"github.com/google/uuid"
)

type GetProfileHandler struct {
	uow port.UnitOfWork
}

func NewGetProfileHandler(uow port.UnitOfWork) *GetProfileHandler {
	return &GetProfileHandler{uow: uow}
}

func (h *GetProfileHandler) Handle(ctx context.Context, query dto.GetProfileQuery) (*dto.UserInfo, error) {
	var user *identity.User
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		var err error
		user, err = repos.Users.GetWithRoles(ctx, query.UserID)
		return err
	})

	if err != nil {
		if apperr.IsNotFound(err) {
			return nil, apperr.NotFound("user", query.UserID)
		}
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
			allMenus, err = repos.Menus.List(ctx, nil)
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
