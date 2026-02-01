package app

import (
	"context"
	"errors"

	"goyavision/config"
	"goyavision/internal/api/middleware"
	"goyavision/internal/domain"
	"goyavision/internal/port"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserDisabled       = errors.New("user is disabled")
	ErrInvalidToken       = errors.New("invalid token")
	ErrPasswordMismatch   = errors.New("old password is incorrect")
)

type AuthService struct {
	repo port.Repository
	cfg  config.JWT
}

func NewAuthService(repo port.Repository, cfg config.JWT) *AuthService {
	return &AuthService{
		repo: repo,
		cfg:  cfg,
	}
}

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	User         *UserInfo
}

type UserInfo struct {
	ID          uuid.UUID
	Username    string
	Nickname    string
	Email       string
	Phone       string
	Avatar      string
	Roles       []string
	Permissions []string
	Menus       []*domain.Menu
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsEnabled() {
		return nil, ErrUserDisabled
	}

	accessToken, err := middleware.GenerateToken(s.cfg, user.ID, user.Username, middleware.TokenTypeAccess)
	if err != nil {
		return nil, err
	}

	refreshToken, err := middleware.GenerateToken(s.cfg, user.ID, user.Username, middleware.TokenTypeRefresh)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.GetUserInfo(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.cfg.Expire.Seconds()),
		User:         userInfo,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	claims, err := middleware.ParseToken(s.cfg, refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims.TokenType != middleware.TokenTypeRefresh {
		return nil, ErrInvalidToken
	}

	user, err := s.repo.GetUser(ctx, claims.UserID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if !user.IsEnabled() {
		return nil, ErrUserDisabled
	}

	accessToken, err := middleware.GenerateToken(s.cfg, user.ID, user.Username, middleware.TokenTypeAccess)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := middleware.GenerateToken(s.cfg, user.ID, user.Username, middleware.TokenTypeRefresh)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.cfg.Expire.Seconds()),
	}, nil
}

func (s *AuthService) GetUserInfo(ctx context.Context, userID uuid.UUID) (*UserInfo, error) {
	user, err := s.repo.GetUserWithRoles(ctx, userID)
	if err != nil {
		return nil, err
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
	var menus []*domain.Menu

	if isSuperAdmin {
		permissions = []string{"*"}
		allMenus, err := s.repo.ListMenus(ctx, nil)
		if err != nil {
			return nil, err
		}
		menus = buildMenuTree(allMenus)
	} else {
		perms, err := s.repo.GetPermissionsByRoleIDs(ctx, roleIDs)
		if err != nil {
			return nil, err
		}
		for _, p := range perms {
			permissions = append(permissions, p.Code)
		}

		menuList, err := s.repo.GetMenusByRoleIDs(ctx, roleIDs)
		if err != nil {
			return nil, err
		}
		menus = buildMenuTree(menuList)
	}

	return &UserInfo{
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

type ChangePasswordRequest struct {
	OldPassword string
	NewPassword string
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, req *ChangePasswordRequest) error {
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return ErrPasswordMismatch
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.repo.UpdateUser(ctx, user)
}

func buildMenuTree(menus []*domain.Menu) []*domain.Menu {
	menuMap := make(map[uuid.UUID]*domain.Menu)
	var roots []*domain.Menu

	for _, m := range menus {
		menuCopy := *m
		menuCopy.Children = []domain.Menu{}
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

// HashPassword 生成密码哈希
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
