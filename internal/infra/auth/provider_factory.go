package auth

import (
	"fmt"

	"goyavision/config"
	"goyavision/internal/app/port"
)

type ProviderFactory struct {
	cfg *config.Config
}

func NewProviderFactory(cfg *config.Config) *ProviderFactory {
	return &ProviderFactory{cfg: cfg}
}

func (f *ProviderFactory) Get(provider string) (port.AuthProvider, error) {
	switch provider {
	case "github":
		return NewGitHubProvider(f.cfg.OAuth.Github), nil
	case "wechat":
		return NewWechatProvider(f.cfg.OAuth.Wechat), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
