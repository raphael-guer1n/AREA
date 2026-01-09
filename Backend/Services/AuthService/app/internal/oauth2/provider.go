package oauth2

import (
	"bytes"
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

// UserInfo represents generic user information returned by the provider
type UserInfo struct {
	ID       string                 `json:"id"`
	Email    string                 `json:"email"`
	Name     string                 `json:"name"`
	Username string                 `json:"username,omitempty"`
	RawData  map[string]interface{} `json:"raw_data"`
}

// Provider handles OAuth2 flow for a specific provider
type Provider struct {
	config ProviderConfig
}

// NewProvider creates a new OAuth2 provider
func NewProvider(config ProviderConfig) *Provider {
	return &Provider{config: config}
}

// GenerateAuthURL builds the OAuth2 authorization URL with the given state
func (p *Provider) GenerateAuthURL(state string, callbackURI string) string {
	params := url.Values{}
	params.Add("client_id", p.config.ClientID)

	if callbackURI != "" {
		params.Add("redirect_uri", callbackURI)
	} else {
		params.Add("redirect_uri", p.config.RedirectURI)
	}

	params.Add("response_type", "code")
	params.Add("state", state)

	if len(p.config.Scopes) > 0 {
		params.Add("scope", strings.Join(p.config.Scopes, " "))
	}

	return fmt.Sprintf("%s?%s", p.config.AuthURL, params.Encode())
}

// ExchangeCodeWithRedirect exchanges an auth code for tokens using a specific redirect URI.
func (p *Provider) ExchangeCodeWithRedirect(code, redirectURI string) (*TokenResponse, error) {
	redirect := redirectURI
	if redirect == "" {
		redirect = p.config.RedirectURI
	}

	// Notion requires Basic auth + JSON payload for the token exchange endpoint.
	if strings.Contains(p.config.TokenURL, "notion.com") && p.config.Name == "notion" {
		payload := map[string]string{
			"grant_type":   "authorization_code",
			"code":         code,
			"redirect_uri": redirect,
		}
		body, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal token request: %w", err)
		}

		req, err := http.NewRequest("POST", p.config.TokenURL, bytes.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("failed to create token request: %w", err)
		}

		basic := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", p.config.ClientID, p.config.ClientSecret)))
		req.Header.Set("Authorization", "Basic "+basic)
		req.Header.Set("Content-Type", "application/json")
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

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirect)
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

// ExchangeCode exchanges the authorization code for an access token
func (p *Provider) ExchangeCode(code string, callbackUri string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	redirectURI := p.config.RedirectURI
	if callbackUri != "" {
		redirectURI = callbackUri
	}
	data.Set("redirect_uri", redirectURI)
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
	var (
		req *http.Request
		err error
	)

	// Dropbox users/get_current_account requires POST with an empty body.
	if strings.Contains(p.config.UserInfoURL, "dropboxapi.com/2/users/get_current_account") || p.config.Name == "dropbox" {
		req, err = http.NewRequest("POST", p.config.UserInfoURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create user info request: %w", err)
		}
	} else {
		req, err = http.NewRequest("GET", p.config.UserInfoURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create user info request: %w", err)
		}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/json")
	// Notion requires the Notion-Version header on every API call.
	if strings.Contains(p.config.UserInfoURL, "notion.com") || p.config.Name == "notion" {
		req.Header.Set("Notion-Version", "2022-06-28")
	}

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

	user := &UserInfo{RawData: rawData}
	if id, ok := rawData["id"].(string); ok {
		user.ID = id
	} else if sub, ok := rawData["sub"].(string); ok {
		user.ID = sub
	}
	if email, ok := rawData["email"].(string); ok {
		user.Email = email
	}
	if name, ok := rawData["name"].(string); ok {
		user.Name = name
	}
	if username, ok := rawData["username"].(string); ok {
		user.Username = username
	} else if login, ok := rawData["login"].(string); ok {
		user.Username = login
	}

	return user, nil
}

// GenerateState builds a random CSRF protection parameter
func GenerateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random state: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
