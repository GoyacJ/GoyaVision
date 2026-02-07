package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"goyavision/config"
	"goyavision/internal/app/port"
)

type WechatProvider struct {
	cfg config.OAuthConfig
}

func NewWechatProvider(cfg config.OAuthConfig) *WechatProvider {
	return &WechatProvider{cfg: cfg}
}

func (p *WechatProvider) Name() string {
	return "wechat"
}

func (p *WechatProvider) GetLoginURL(state string) string {
	return fmt.Sprintf("https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect",
		p.cfg.ClientID, url.QueryEscape(p.cfg.RedirectURI), state)
}

func (p *WechatProvider) VerifyCode(ctx context.Context, code string) (*port.AuthUserInfo, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	// 1. Get Access Token
	tokenURL := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		p.cfg.ClientID, p.cfg.ClientSecret, code)

	resp, err := client.Get(tokenURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		OpenID      string `json:"openid"`
		UnionID     string `json:"unionid"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}
	if tokenResp.ErrCode != 0 {
		return nil, fmt.Errorf("wechat token error: %s", tokenResp.ErrMsg)
	}

	// 2. Get User Info
	userURL := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s", tokenResp.AccessToken, tokenResp.OpenID)
	resp, err = client.Get(userURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userResp struct {
		OpenID   string `json:"openid"`
		Nickname string `json:"nickname"`
		HeadImg  string `json:"headimgurl"`
		UnionID  string `json:"unionid"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, err
	}

	id := userResp.UnionID
	if id == "" {
		id = userResp.OpenID
	}

	return &port.AuthUserInfo{
		ID:     id,
		Name:   userResp.Nickname,
		Email:  "", // WeChat doesn't provide email
		Avatar: userResp.HeadImg,
		Raw:    map[string]interface{}{"openid": userResp.OpenID},
	}, nil
}
