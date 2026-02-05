package port

import (
	"time"

	"github.com/google/uuid"
)

// TokenService JWT 令牌服务接口
//
// 职责：
//  1. 生成 Access Token 和 Refresh Token
//  2. 验证和解析 Token
//  3. Token 刷新
//
// 实现：
//  - infra/auth/jwt.go (golang-jwt/jwt 实现)
type TokenService interface {
	// GenerateTokenPair 生成 Token 对（Access + Refresh）
	GenerateTokenPair(userID uuid.UUID, username string) (*TokenPair, error)

	// ValidateAccessToken 验证 Access Token
	// 返回：Claims、是否过期、错误
	ValidateAccessToken(token string) (*TokenClaims, bool, error)

	// ValidateRefreshToken 验证 Refresh Token
	ValidateRefreshToken(token string) (*TokenClaims, bool, error)

	// RefreshTokenPair 使用 Refresh Token 刷新 Token 对
	RefreshTokenPair(refreshToken string) (*TokenPair, error)
}

// TokenPair Token 对
type TokenPair struct {
	AccessToken  string    // 访问令牌（短期有效，如 2 小时）
	RefreshToken string    // 刷新令牌（长期有效，如 7 天）
	ExpiresIn    int64     // Access Token 过期时间（秒）
	ExpiresAt    time.Time // Access Token 过期时刻
}

// TokenClaims Token 声明
type TokenClaims struct {
	UserID    uuid.UUID
	Username  string
	IssuedAt  time.Time
	ExpiresAt time.Time
}
