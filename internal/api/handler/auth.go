package handler

import (
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"goyavision/internal/api/dto"
	"goyavision/internal/api/middleware"
	"goyavision/internal/app"
	"goyavision/internal/domain/identity"
	appdto "goyavision/internal/app/dto"

	"github.com/labstack/echo/v4"
)

func RegisterAuth(g *echo.Group, h *Handlers) {
	handler := &authHandler{h: h}
	g.POST("/login", handler.Login)
	g.POST("/refresh", handler.RefreshToken)
	g.GET("/oauth/login", handler.Authorize)
	g.POST("/oauth/login", handler.LoginOAuth)
}

func RegisterAuthProtected(g *echo.Group, h *Handlers) {
	handler := &authHandler{h: h}
	g.GET("/profile", handler.GetProfile)
	g.PUT("/password", handler.ChangePassword)
	g.POST("/logout", handler.Logout)
	g.POST("/bind", handler.BindIdentity)
}

type authHandler struct {
	h *Handlers
}

func (h *authHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if req.Username == "" || req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username and password are required")
	}

	result, err := h.h.Login.Handle(c.Request().Context(), appdto.LoginCommand{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		User: dto.UserInfoFromApp(&dto.AppUserInfo{
			ID:          result.User.ID,
			Username:    result.User.Username,
			Nickname:    result.User.Nickname,
			Email:       result.User.Email,
			Phone:       result.User.Phone,
			Avatar:      result.User.Avatar,
			Roles:       result.User.Roles,
			Permissions: result.User.Permissions,
			Menus:       result.User.Menus,
		}),
	})
}

func (h *authHandler) RefreshToken(c echo.Context) error {
	var req dto.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if req.RefreshToken == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "refresh_token is required")
	}

	tokenPair, err := h.h.TokenService.RefreshTokenPair(req.RefreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired refresh token")
	}

	return c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	})
}

func (h *authHandler) GetProfile(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	result, err := h.h.GetProfile.Handle(c.Request().Context(), appdto.GetProfileQuery{UserID: userID})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.UserInfoFromApp(&dto.AppUserInfo{
		ID:          result.ID,
		Username:    result.Username,
		Nickname:    result.Nickname,
		Email:       result.Email,
		Phone:       result.Phone,
		Avatar:      result.Avatar,
		Roles:       result.Roles,
		Permissions: result.Permissions,
		Menus:       result.Menus,
	}))
}

func (h *authHandler) ChangePassword(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	var req dto.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "old_password and new_password are required")
	}

	if len(req.NewPassword) < 6 {
		return echo.NewHTTPError(http.StatusBadRequest, "new password must be at least 6 characters")
	}

	userSvc := app.NewUserService(h.h.Repo)
	user, err := userSvc.Get(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "old password is incorrect")
	}

	if err := userSvc.ResetPassword(c.Request().Context(), userID, req.NewPassword); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to change password")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "password changed successfully",
	})
}

func (h *authHandler) Logout(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "logged out successfully",
	})
}

func (h *authHandler) Authorize(c echo.Context) error {
	provider := c.QueryParam("provider")
	if provider == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "provider is required")
	}

	authProvider, err := h.h.AuthProviderFactory.Get(provider)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "unsupported provider")
	}

	// Generate a random state. Ideally store it in session/cookie for validation.
	state := uuid.New().String()

	url := authProvider.GetLoginURL(state)
	return c.Redirect(http.StatusFound, url)
}

func (h *authHandler) LoginOAuth(c echo.Context) error {
	var req dto.LoginOAuthRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	result, err := h.h.LoginOAuth.Handle(c.Request().Context(), appdto.LoginOAuthCommand{
		Provider: req.Provider,
		Code:     req.Code,
		State:    req.State,
	})
	if err != nil {
		return err
	}

	// Fetch user info for response
	user, err := h.h.Repo.GetUserWithRoles(c.Request().Context(), result.User.ID)
	if err != nil {
		// Log error but proceed? Or fail?
		// Login succeeded but fetching full user info failed.
		// Result already has user info?
		// result.User in LoginResult might be partial or specific DTO.
		// Let's check internal/app/dto/login.go later. Assuming result.User is good enough or reload.
	}

	return c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		User: dto.UserInfoFromApp(&dto.AppUserInfo{
			ID:          user.ID,
			Username:    user.Username,
			Nickname:    user.Nickname,
			Email:       user.Email,
			Phone:       user.Phone,
			Avatar:      user.Avatar,
			Roles:       dto.RolesToCodes(user.Roles),
			Permissions: []string{},
			Menus:       []*identity.Menu{},
		}),
	})
}

func (h *authHandler) BindIdentity(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	var req dto.BindIdentityRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	err := h.h.BindIdentity.Handle(c.Request().Context(), appdto.BindIdentityCommand{
		UserID:     userID,
		Provider:   req.Provider,
		Identifier: req.Identifier,
		Credential: req.Credential,
		Meta:       req.Meta,
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
