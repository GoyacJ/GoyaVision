package auth

import (
	"context"

	"goyavision/internal/app/port"
)

type MockProvider struct{ name string }

func (p *MockProvider) Name() string { return p.name }
func (p *MockProvider) GetLoginURL(state string) string { return "http://mock-login" }
func (p *MockProvider) VerifyCode(ctx context.Context, code string) (*port.AuthUserInfo, error) {
	return &port.AuthUserInfo{
		ID:     "mock-id-" + code,
		Name:   "Mock User " + p.name,
		Email:  "mock@example.com",
		Avatar: "",
		Raw:    map[string]interface{}{"provider": p.name},
	}, nil
}

type ProviderFactory struct{}

func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{}
}

func (f *ProviderFactory) Get(provider string) (port.AuthProvider, error) {
	return &MockProvider{name: provider}, nil
}
