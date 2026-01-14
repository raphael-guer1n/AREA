package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type AuthService struct {
	baseURL string
	client  *http.Client
}

func NewAuthService(baseURL string) *AuthService {
	return &AuthService{
		baseURL: strings.TrimRight(baseURL, "/"),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *AuthService) GetUserID(authHeader string) (int, error) {
	if strings.TrimSpace(authHeader) == "" {
		return 0, errors.New("missing authorization header")
	}
	endpoint := s.baseURL + "/auth/me"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", authHeader)

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("fetch user: %w", err)
	}
	defer resp.Body.Close()

	var body struct {
		Success bool `json:"success"`
		Data    struct {
			User struct {
				ID int `json:"id"`
			} `json:"user"`
		} `json:"data"`
		Error json.RawMessage `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return 0, fmt.Errorf("decode user: %w", err)
	}

	if resp.StatusCode != http.StatusOK || !body.Success {
		message, _ := parseRemoteError(body.Error)
		if message == "" {
			message = "failed to fetch user"
		}
		return 0, errors.New(message)
	}

	if body.Data.User.ID <= 0 {
		return 0, errors.New("user not found")
	}

	return body.Data.User.ID, nil
}
