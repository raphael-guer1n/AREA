package oauth2

import (
	"fmt"
	"sync"
)

// Manager manages multiple OAuth2 providers
type Manager struct {
	config    *ConfigOAuth2
	providers map[string]*Provider
	states    *StateStore
}

// StateStore manages OAuth2 state parameters for CSRF protection
type StateStore struct {
	mu     sync.RWMutex
	states map[string]string // state -> provider name
}

// NewStateStore creates a new state store
func NewStateStore() *StateStore {
	return &StateStore{
		states: make(map[string]string),
	}
}

// Set stores a state-provider mapping
func (s *StateStore) Set(state, provider string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[state] = provider
}

// Get retrieves and removes a state-provider mapping
func (s *StateStore) Get(state string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	provider, ok := s.states[state]
	if ok {
		delete(s.states, state) // State is single-use
	}
	return provider, ok
}

// NewManager creates a new OAuth2 manager
func NewManager(config *ConfigOAuth2) *Manager {
	providers := make(map[string]*Provider)
	for name, cfg := range config.Providers {
		providers[name] = NewProvider(cfg)
	}

	return &Manager{
		config:    config,
		providers: providers,
		states:    NewStateStore(),
	}
}

// GetAuthURL generates an OAuth2 authorization URL for a specific provider
func (m *Manager) GetAuthURL(providerName string) (string, error) {
	provider, ok := m.providers[providerName]
	if !ok {
		return "", fmt.Errorf("provider %s not found", providerName)
	}

	// Generate CSRF state token
	state, err := GenerateState()
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	// Store state for validation
	m.states.Set(state, providerName)

	// Generate authorization URL
	authURL := provider.GenerateAuthURL(state)

	return authURL, nil
}

// HandleCallback handles the OAuth2 callback with code and state
func (m *Manager) HandleCallback(state, code string) (*UserInfo, *TokenResponse, string, error) {
	// Validate state and get provider name
	providerName, ok := m.states.Get(state)
	if !ok {
		return nil, nil, "", fmt.Errorf("invalid or expired state parameter")
	}

	provider, ok := m.providers[providerName]
	if !ok {
		return nil, nil, "", fmt.Errorf("provider %s not found", providerName)
	}

	// Exchange code for access token
	tokenResp, err := provider.ExchangeCode(code)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info using access token
	userInfo, err := provider.GetUserInfo(tokenResp.AccessToken)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to get user info: %w", err)
	}

	return userInfo, tokenResp, providerName, nil
}

// ListProviders returns all available provider names
func (m *Manager) ListProviders() []string {
	return m.config.ListProviders()
}

// GetProvider returns a specific provider
func (m *Manager) GetProvider(name string) (*Provider, error) {
	provider, ok := m.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return provider, nil
}
