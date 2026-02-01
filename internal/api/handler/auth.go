package handler

import (
	"net/http"

	"goyavision/internal/api/dto"
	"goyavision/internal/api/middleware"
	"goyavision/internal/app"

	"github.com/labstack/echo/v4"
)

func RegisterAuth(g *echo.Group, d Deps) {
	svc := app.NewAuthService(d.Repo, d.Cfg.JWT)
	h := authHandler{svc: svc, deps: d}

	g.POST("/login", h.Login)
	g.POST("/refresh", h.RefreshToken)
}

func RegisterAuthProtected(g *echo.Group, d Deps) {
	svc := app.NewAuthService(d.Repo, d.Cfg.JWT)
	h := authHandler{svc: svc, deps: d}

	g.GET("/profile", h.GetProfile)
	g.PUT("/password", h.ChangePassword)
	g.POST("/logout", h.Logout)
}

type authHandler struct {
	svc  *app.AuthService
	deps Deps
}

func (h *authHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "username and password are required",
		})
	}

	result, err := h.svc.Login(c.Request().Context(), &app.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		switch err {
		case app.ErrInvalidCredentials:
			return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "Unauthorized",
				Message: "invalid username or password",
			})
		case app.ErrUserDisabled:
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "Forbidden",
				Message: "user is disabled",
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Internal Server Error",
				Message: err.Error(),
			})
		}
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
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	if req.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "refresh_token is required",
		})
	}

	result, err := h.svc.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "invalid or expired refresh token",
		})
	}

	return c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	})
}

func (h *authHandler) GetProfile(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "user not authenticated",
		})
	}

	userInfo, err := h.svc.GetUserInfo(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.UserInfoFromApp(&dto.AppUserInfo{
		ID:          userInfo.ID,
		Username:    userInfo.Username,
		Nickname:    userInfo.Nickname,
		Email:       userInfo.Email,
		Phone:       userInfo.Phone,
		Avatar:      userInfo.Avatar,
		Roles:       userInfo.Roles,
		Permissions: userInfo.Permissions,
		Menus:       userInfo.Menus,
	}))
}

func (h *authHandler) ChangePassword(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "user not authenticated",
		})
	}

	var req dto.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "old_password and new_password are required",
		})
	}

	if len(req.NewPassword) < 6 {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "new password must be at least 6 characters",
		})
	}

	err := h.svc.ChangePassword(c.Request().Context(), userID, &app.ChangePasswordRequest{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		if err == app.ErrPasswordMismatch {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Bad Request",
				Message: "old password is incorrect",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
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
