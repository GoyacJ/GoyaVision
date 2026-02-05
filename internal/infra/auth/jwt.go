package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"goyavision/config"
	"goyavision/internal/app/port"
	"goyavision/pkg/apperr"
)

// JWTService JWT 令牌服务实现
type JWTService struct {
	secret           []byte
	issuer           string
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
}

// customClaims JWT 自定义声明
type customClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Type     string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// NewJWTService 创建 JWT 服务实例
func NewJWTService(cfg *config.JWT) (*JWTService, error) {
	if cfg.Secret == "" {
		return nil, apperr.InvalidInput("JWT secret is required")
	}
	if cfg.Expire <= 0 {
		cfg.Expire = 2 * time.Hour
	}
	if cfg.RefreshExp <= 0 {
		cfg.RefreshExp = 7 * 24 * time.Hour
	}
	if cfg.Issuer == "" {
		cfg.Issuer = "goyavision"
	}

	return &JWTService{
		secret:          []byte(cfg.Secret),
		issuer:          cfg.Issuer,
		accessTokenTTL:  cfg.Expire,
		refreshTokenTTL: cfg.RefreshExp,
	}, nil
}

// GenerateTokenPair 生成 Token 对（Access + Refresh）
func (s *JWTService) GenerateTokenPair(userID uuid.UUID, username string) (*port.TokenPair, error) {
	if userID == uuid.Nil {
		return nil, apperr.InvalidInput("user ID is required")
	}
	if username == "" {
		return nil, apperr.InvalidInput("username is required")
	}

	now := time.Now()
	accessExpiresAt := now.Add(s.accessTokenTTL)
	refreshExpiresAt := now.Add(s.refreshTokenTTL)

	accessToken, err := s.generateToken(userID, username, "access", accessExpiresAt)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeInternal, "failed to generate access token")
	}

	refreshToken, err := s.generateToken(userID, username, "refresh", refreshExpiresAt)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeInternal, "failed to generate refresh token")
	}

	return &port.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.accessTokenTTL.Seconds()),
		ExpiresAt:    accessExpiresAt,
	}, nil
}

// ValidateAccessToken 验证 Access Token
func (s *JWTService) ValidateAccessToken(tokenString string) (*port.TokenClaims, bool, error) {
	return s.validateToken(tokenString, "access")
}

// ValidateRefreshToken 验证 Refresh Token
func (s *JWTService) ValidateRefreshToken(tokenString string) (*port.TokenClaims, bool, error) {
	return s.validateToken(tokenString, "refresh")
}

// RefreshTokenPair 使用 Refresh Token 刷新 Token 对
func (s *JWTService) RefreshTokenPair(refreshToken string) (*port.TokenPair, error) {
	if refreshToken == "" {
		return nil, apperr.Unauthorized("refresh token is required")
	}

	claims, isExpired, err := s.ValidateRefreshToken(refreshToken)
	if err != nil {
		if isExpired {
			return nil, apperr.New(apperr.CodeTokenExpired, "refresh token has expired")
		}
		return nil, apperr.New(apperr.CodeTokenInvalid, "invalid refresh token")
	}

	return s.GenerateTokenPair(claims.UserID, claims.Username)
}

// generateToken 生成 JWT token
func (s *JWTService) generateToken(userID uuid.UUID, username, tokenType string, expiresAt time.Time) (string, error) {
	now := time.Now()

	claims := &customClaims{
		UserID:   userID.String(),
		Username: username,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// validateToken 验证 token 并返回 claims
func (s *JWTService) validateToken(tokenString, expectedType string) (*port.TokenClaims, bool, error) {
	if tokenString == "" {
		return nil, false, apperr.Unauthorized("token is required")
	}

	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		if err == jwt.ErrTokenExpired {
			return nil, true, apperr.New(apperr.CodeTokenExpired, "token has expired")
		}
		return nil, false, apperr.Wrap(err, apperr.CodeTokenInvalid, "invalid token")
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok || !token.Valid {
		return nil, false, apperr.New(apperr.CodeTokenInvalid, "invalid token claims")
	}

	if claims.Type != expectedType {
		return nil, false, apperr.New(apperr.CodeTokenInvalid, fmt.Sprintf("expected %s token, got %s", expectedType, claims.Type))
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, false, apperr.Wrap(err, apperr.CodeTokenInvalid, "invalid user ID in token")
	}

	isExpired := false
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		isExpired = true
	}

	return &port.TokenClaims{
		UserID:    userID,
		Username:  claims.Username,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
	}, isExpired, nil
}
