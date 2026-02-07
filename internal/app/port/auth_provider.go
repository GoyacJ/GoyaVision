package port

import "context"

type AuthUserInfo struct {
	ID       string
	Name     string
	Email    string
	Avatar   string
	Raw      map[string]interface{}
}

type AuthProvider interface {
	Name() string
	GetLoginURL(state string) string
	VerifyCode(ctx context.Context, code string) (*AuthUserInfo, error)
}

type AuthProviderFactory interface {
	Get(provider string) (AuthProvider, error)
}
