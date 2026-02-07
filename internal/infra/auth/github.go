package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"goyavision/config"
	"goyavision/internal/app/port"
)

type GitHubProvider struct {
	cfg config.OAuthConfig
}

func NewGitHubProvider(cfg config.OAuthConfig) *GitHubProvider {
	return &GitHubProvider{cfg: cfg}
}

func (p *GitHubProvider) Name() string {
	return "github"
}

func (p *GitHubProvider) GetLoginURL(state string) string {
	return fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s&scope=user:email",
		p.cfg.ClientID, p.cfg.RedirectURI, state)
}

func (p *GitHubProvider) VerifyCode(ctx context.Context, code string) (*port.AuthUserInfo, error) {
	// Exchange code for token
	tokenReq := map[string]string{
		"client_id":     p.cfg.ClientID,
		"client_secret": p.cfg.ClientSecret,
		"code":          code,
	}
	tokenBody, _ := json.Marshal(tokenReq)

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(tokenBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}
	if tokenResp.Error != "" {
		return nil, fmt.Errorf("github error: %s", tokenResp.Error)
	}

	// Get User Info
	req, _ = http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	resp, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var userResp struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	rawBody, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(rawBody, &userResp); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	name := userResp.Name
	if name == "" {
		name = userResp.Login
	}

	return &port.AuthUserInfo{
		ID:     fmt.Sprintf("%d", userResp.ID),
		Name:   name,
		Email:  userResp.Email,
		Avatar: userResp.AvatarURL,
		Raw:    map[string]interface{}{"login": userResp.Login, "raw": string(rawBody)},
	}, nil
}
