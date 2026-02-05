package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"goyavision/config"
	"goyavision/internal/api/dto"
	"goyavision/internal/port"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Context keys
const (
	ContextKeyUserID      = "user_id"
	ContextKeyUsername    = "username"
	ContextKeyRoles       = "roles"
	ContextKeyPermissions = "permissions"
)

// Token types
const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// JWTClaims 自定义 JWT Claims
type JWTClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	TokenType string    `json:"token_type"`
	LegacyType string   `json:"type"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(cfg config.JWT, userID uuid.UUID, username string, tokenType string) (string, error) {
	var expiration time.Duration
	if tokenType == TokenTypeRefresh {
		expiration = cfg.RefreshExp
	} else {
		expiration = cfg.Expire
	}

	claims := JWTClaims{
		UserID:    userID,
		Username:  username,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// ParseToken 解析 JWT Token
func ParseToken(cfg config.JWT, tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// JWTAuth JWT 认证中间件
func JWTAuth(cfg config.JWT) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error:   "Unauthorized",
					Message: "missing authorization header",
				})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error:   "Unauthorized",
					Message: "invalid authorization header format",
				})
			}

			claims, err := ParseToken(cfg, parts[1])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error:   "Unauthorized",
					Message: "invalid or expired token",
				})
			}

			tokenType := claims.TokenType
			if tokenType == "" {
				tokenType = claims.LegacyType
			}
			if tokenType != TokenTypeAccess {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error:   "Unauthorized",
					Message: "invalid token type",
				})
			}

			c.Set(ContextKeyUserID, claims.UserID)
			c.Set(ContextKeyUsername, claims.Username)

			return next(c)
		}
	}
}

// RequirePermission 权限校验中间件
func RequirePermission(repo port.Repository, requiredPermissions ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, ok := c.Get(ContextKeyUserID).(uuid.UUID)
			if !ok {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error:   "Unauthorized",
					Message: "user not authenticated",
				})
			}

			user, err := repo.GetUserWithRoles(c.Request().Context(), userID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error:   "Unauthorized",
					Message: "user not found",
				})
			}

			if !user.IsEnabled() {
				return c.JSON(http.StatusForbidden, dto.ErrorResponse{
					Error:   "Forbidden",
					Message: "user is disabled",
				})
			}

			var roleIDs []uuid.UUID
			var roleCodes []string
			for _, role := range user.Roles {
				if role.IsEnabled() {
					roleIDs = append(roleIDs, role.ID)
					roleCodes = append(roleCodes, role.Code)
				}
			}

			for _, code := range roleCodes {
				if code == "super_admin" {
					c.Set(ContextKeyRoles, roleCodes)
					c.Set(ContextKeyPermissions, []string{"*"})
					return next(c)
				}
			}

			permissions, err := repo.GetPermissionsByRoleIDs(c.Request().Context(), roleIDs)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error:   "Internal Server Error",
					Message: "failed to get permissions",
				})
			}

			permissionCodes := make([]string, len(permissions))
			for i, p := range permissions {
				permissionCodes[i] = p.Code
			}

			c.Set(ContextKeyRoles, roleCodes)
			c.Set(ContextKeyPermissions, permissionCodes)

			if len(requiredPermissions) > 0 {
				hasPermission := false
				for _, required := range requiredPermissions {
					for _, code := range permissionCodes {
						if code == required || code == "*" {
							hasPermission = true
							break
						}
					}
					if hasPermission {
						break
					}
				}
				if !hasPermission {
					return c.JSON(http.StatusForbidden, dto.ErrorResponse{
						Error:   "Forbidden",
						Message: "insufficient permissions",
					})
				}
			}

			return next(c)
		}
	}
}

// LoadUserPermissions 加载用户权限中间件（不校验，只加载）
func LoadUserPermissions(repo port.Repository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, ok := c.Get(ContextKeyUserID).(uuid.UUID)
			if !ok {
				return next(c)
			}

			user, err := repo.GetUserWithRoles(c.Request().Context(), userID)
			if err != nil {
				return next(c)
			}

			var roleIDs []uuid.UUID
			var roleCodes []string
			for _, role := range user.Roles {
				if role.IsEnabled() {
					roleIDs = append(roleIDs, role.ID)
					roleCodes = append(roleCodes, role.Code)
				}
			}

			for _, code := range roleCodes {
				if code == "super_admin" {
					c.Set(ContextKeyRoles, roleCodes)
					c.Set(ContextKeyPermissions, []string{"*"})
					return next(c)
				}
			}

			permissions, err := repo.GetPermissionsByRoleIDs(c.Request().Context(), roleIDs)
			if err != nil {
				return next(c)
			}

			permissionCodes := make([]string, len(permissions))
			for i, p := range permissions {
				permissionCodes[i] = p.Code
			}

			c.Set(ContextKeyRoles, roleCodes)
			c.Set(ContextKeyPermissions, permissionCodes)

			return next(c)
		}
	}
}

// GetUserID 从 Context 获取用户 ID
func GetUserID(c echo.Context) (uuid.UUID, bool) {
	userID, ok := c.Get(ContextKeyUserID).(uuid.UUID)
	return userID, ok
}

// GetUsername 从 Context 获取用户名
func GetUsername(c echo.Context) (string, bool) {
	username, ok := c.Get(ContextKeyUsername).(string)
	return username, ok
}

// GetPermissions 从 Context 获取权限列表
func GetPermissions(c echo.Context) []string {
	permissions, ok := c.Get(ContextKeyPermissions).([]string)
	if !ok {
		return []string{}
	}
	return permissions
}

// HasPermission 检查是否有指定权限
func HasPermission(c echo.Context, permission string) bool {
	permissions := GetPermissions(c)
	for _, p := range permissions {
		if p == permission || p == "*" {
			return true
		}
	}
	return false
}
