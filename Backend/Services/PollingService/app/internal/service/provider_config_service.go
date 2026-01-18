package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/raphael-guer1n/AREA/PollingService/internal/config"
)

var ErrProviderConfigNotFound = errors.New("provider not found")

type ProviderConfigService struct {
	baseURL        string
	internalSecret string
	client         *http.Client
	cache          map[string]config.PollingProviderConfig
}

func NewProviderConfigService(baseURL string, internalSecret string) *ProviderConfigService {
	return &ProviderConfigService{
		baseURL:        strings.TrimRight(baseURL, "/"),
		internalSecret: strings.TrimSpace(internalSecret),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		cache: make(map[string]config.PollingProviderConfig),
	}
}

func (s *ProviderConfigService) GetProviderConfig(name string) (*config.PollingProviderConfig, error) {
	if cfg, ok := s.cache[name]; ok {
		return &cfg, nil
	}

	cfg, err := s.fetchProviderConfig(name)
	if err != nil {
		return nil, err
	}
	s.cache[name] = *cfg
	return cfg, nil
}

func (s *ProviderConfigService) ListProviders() ([]string, error) {
	endpoint := s.baseURL + "/polling/providers"

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	if s.internalSecret != "" {
		req.Header.Set("X-Internal-Secret", s.internalSecret)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch providers: %w", err)
	}
	defer resp.Body.Close()

	var body struct {
		Success bool `json:"success"`
		Data    struct {
			Providers []string `json:"providers"`
		} `json:"data"`
		Error json.RawMessage `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode providers: %w", err)
	}

	if resp.StatusCode != http.StatusOK || !body.Success {
		message, _ := parseRemoteError(body.Error)
		if message == "" {
			message = "failed to fetch providers"
		}
		return nil, fmt.Errorf(message)
	}

	return body.Data.Providers, nil
}

func (s *ProviderConfigService) fetchProviderConfig(name string) (*config.PollingProviderConfig, error) {
	endpoint := s.baseURL + "/polling/providers/config"
	params := url.Values{}
	params.Set("provider", name)
	endpoint = endpoint + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	if s.internalSecret != "" {
		req.Header.Set("X-Internal-Secret", s.internalSecret)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch provider config: %w", err)
	}
	defer resp.Body.Close()

	var body struct {
		Success bool                         `json:"success"`
		Data    config.PollingProviderConfig `json:"data"`
		Error   json.RawMessage              `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode provider config: %w", err)
	}

	if resp.StatusCode != http.StatusOK || !body.Success {
		message, _ := parseRemoteError(body.Error)
		if resp.StatusCode == http.StatusNotFound && strings.EqualFold(message, "provider not found") {
			return nil, ErrProviderConfigNotFound
		}
		if message == "" {
			message = "failed to fetch provider config"
		}
		return nil, fmt.Errorf(message)
	}

	return &body.Data, nil
}
