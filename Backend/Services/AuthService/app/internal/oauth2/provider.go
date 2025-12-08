package oauth2

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// TokenResponse represents the OAuth2 token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// UserInfo represents generic user information from OAuth2 provider
type UserInfo struct {
	ID       string                 `json:"id"`
	Email    string                 `json:"email"`
	Name     string                 `json:"name"`
	Username string                 `json:"username,omitempty"`
	RawData  map[string]interface{} `json:"raw_data"`
}

// Provider handles OAuth2 authentication flow for a specific provider
type Provider struct {
	config ProviderConfig
}

// NewProvider creates a new OAuth2 provider
func NewProvider(config ProviderConfig) *Provider {
	return &Provider{config: config}
}

// GenerateAuthURL generates the OAuth2 authorization URL with state parameter
func (p *Provider) GenerateAuthURL(state string) string {
	params := url.Values{}
	params.Add("client_id", p.config.ClientID)
	params.Add("redirect_uri", p.config.RedirectURI)
	params.Add("response_type", "code")
	params.Add("state", state)

	if len(p.config.Scopes) > 0 {
		params.Add("scope", strings.Join(p.config.Scopes, " "))
	}

	return fmt.Sprintf("%s?%s", p.config.AuthURL, params.Encode())
}

// ExchangeCode exchanges the authorization code for an access token
func (p *Provider) ExchangeCode(code string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", p.config.RedirectURI)
	data.Set("client_id", p.config.ClientID)
	data.Set("client_secret", p.config.ClientSecret)

	req, err := http.NewRequest("POST", p.config.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}

// GetUserInfo retrieves user information using the access token
func (p *Provider) GetUserInfo(accessToken string) (*UserInfo, error) {
	req, err := http.NewRequest("GET", p.config.UserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user info request failed: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var rawData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawData); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	userInfo := &UserInfo{
		RawData: rawData,
	}

	// Extract common fields (providers may use different field names)
	if id, ok := rawData["id"].(string); ok {
		userInfo.ID = id
	} else if sub, ok := rawData["sub"].(string); ok {
		userInfo.ID = sub
	}

	if email, ok := rawData["email"].(string); ok {
		userInfo.Email = email
	}

	if name, ok := rawData["name"].(string); ok {
		userInfo.Name = name
	}

	if username, ok := rawData["username"].(string); ok {
		userInfo.Username = username
	} else if login, ok := rawData["login"].(string); ok {
		userInfo.Username = login
	}

	return userInfo, nil
}

// GenerateState generates a random state parameter for CSRF protection
func GenerateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random state: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
