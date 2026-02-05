package handler

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
	appdto "goyavision/internal/app/dto"
	"goyavision/internal/api/dto"
	"goyavision/internal/api/middleware"
	"goyavision/internal/app"

	"github.com/labstack/echo/v4"
)

func RegisterAuth(g *echo.Group, h *Handlers) {
	handler := &authHandler{h: h}
	g.POST("/login", handler.Login)
	g.POST("/refresh", handler.RefreshToken)
}

func RegisterAuthProtected(g *echo.Group, h *Handlers) {
	handler := &authHandler{h: h}
	g.GET("/profile", handler.GetProfile)
	g.PUT("/password", handler.ChangePassword)
	g.POST("/logout", handler.Logout)
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
