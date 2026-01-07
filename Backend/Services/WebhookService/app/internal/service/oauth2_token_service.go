package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type OAuth2TokenService struct {
	baseURL string
	client  *http.Client
}

func NewOAuth2TokenService(baseURL string) *OAuth2TokenService {
	return &OAuth2TokenService{
		baseURL: strings.TrimRight(baseURL, "/"),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *OAuth2TokenService) GetProviderToken(userID int, provider string) (string, error) {
	endpoint := s.baseURL + "/oauth2/provider/token/"
	params := url.Values{}
	params.Set("user_id", fmt.Sprintf("%d", userID))
	params.Set("service", provider)
	endpoint = endpoint + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch provider token: %w", err)
	}
	defer resp.Body.Close()

	var body struct {
		Success bool `json:"success"`
		Data    struct {
			ProviderToken string `json:"providerToken"`
		} `json:"data"`
		Error json.RawMessage `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("decode provider token: %w", err)
	}

	if resp.StatusCode != http.StatusOK || !body.Success {
		message, _ := parseRemoteError(body.Error)
		if message == "" {
			message = "failed to fetch provider token"
		}
		return "", fmt.Errorf(message)
	}

	if body.Data.ProviderToken == "" {
		return "", fmt.Errorf("provider token not found")
	}

	return body.Data.ProviderToken, nil
}
