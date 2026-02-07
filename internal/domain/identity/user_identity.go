package identity

import (
	"time"

	"github.com/google/uuid"
)

type IdentityType string

const (
	IdentityTypePassword IdentityType = "password"
	IdentityTypeGithub   IdentityType = "github"
	IdentityTypeWechat   IdentityType = "wechat"
	IdentityTypePhone    IdentityType = "phone"
)

type UserIdentity struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	IdentityType IdentityType
	Identifier   string // username, openid, phone
	Credential   string // password hash, access token (if needed)
	Meta         map[string]interface{}
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
