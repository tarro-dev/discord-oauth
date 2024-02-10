package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ajg/form"
)

const (
	discordTokenAPI  = "https://discord.com/api/oauth2/token"
	discordRevokeAPI = "https://discord.com/api/oauth2/token/revoke"
)

type Discord struct {
	clientID     string
	clientSecret string
	uri          string
}

func NewDiscord(clientID, clientSecret, uri string) *Discord {
	return &Discord{
		clientID:     clientID,
		clientSecret: clientSecret,
		uri:          uri,
	}
}

type Token struct {
	Type    string
	Access  string
	Refresh string
	Expires time.Time
	Scope   string
}

func (d *Discord) RequestToken(ctx context.Context, code string) (Token, error) {
	form := TokenRequest{
		GrantType:    "authorization_code",
		Code:         code,
		RedirectURI:  d.uri,
		ClientID:     d.clientID,
		ClientSecret: d.clientSecret,
	}.Form()

	resp, err := http.PostForm(discordTokenAPI, form)
	if err != nil {
		return Token{}, fmt.Errorf("failed to request token: %w", err)
	}
	defer resp.Body.Close()

	var tresp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tresp); err != nil {
		return Token{}, fmt.Errorf("failed to decode token response: %w", err)
	}

	return Token{
		Type:    tresp.TokenType,
		Access:  tresp.AccessToken,
		Refresh: tresp.RefreshToken,
		Expires: time.Now().Add(time.Duration(tresp.ExpiresIn) * time.Second),
		Scope:   tresp.Scope,
	}, nil
}

func (d *Discord) RefreshToken(ctx context.Context, token string) (Token, error) {
	form := RefreshTokenRequest{
		GrantType:    "refresh_token",
		RefreshToken: token,
		ClientID:     d.clientID,
		ClientSecret: d.clientSecret,
	}.Form()

	resp, err := http.PostForm(discordTokenAPI, form)
	if err != nil {
		return Token{}, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	var tresp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tresp); err != nil {
		return Token{}, fmt.Errorf("failed to decode token response: %w", err)
	}

	return Token{
		Type:    tresp.TokenType,
		Access:  tresp.AccessToken,
		Refresh: tresp.RefreshToken,
		Expires: time.Now().Add(time.Duration(tresp.ExpiresIn) * time.Second),
		Scope:   tresp.Scope,
	}, nil
}

func (d *Discord) RevokeToken(ctx context.Context, token string) error {
	form := RevokeTokenRequest{
		Token: token,
	}.Form()

	resp, err := http.PostForm(discordRevokeAPI, form)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

type TokenRequest struct {
	GrantType    string `form:"grant_type"`
	Code         string `form:"code"`
	RedirectURI  string `form:"redirect_uri"`
	ClientID     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
}

func (tr TokenRequest) Form() url.Values {
	u, _ := form.EncodeToValues(tr)
	return u
}

type RefreshTokenRequest struct {
	GrantType    string `form:"grant_type"`
	RefreshToken string `form:"refresh_token"`
	ClientID     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
}

func (tr RefreshTokenRequest) Form() url.Values {
	u, _ := form.EncodeToValues(tr)
	return u
}

type RevokeTokenRequest struct {
	Token string `form:"token"`
}

func (tr RevokeTokenRequest) Form() url.Values {
	u, _ := form.EncodeToValues(tr)
	return u
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}
