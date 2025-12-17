package oauth2

import (
	"fmt"
	"os"
	"sync"
)

// Manager manages multiple OAuth2 providers with lazy loading
type Manager struct {
	configLoader *ConfigLoader
	providers    map[string]*Provider
	providersMu  sync.RWMutex
	states       *StateStore
}

// StateData holds metadata associated with an OAuth2 state
type StateData struct {
	Provider    string
	UserID      int
	CallbackURL string
	Platform    string // web, android, ios
}

// StateStore manages OAuth2 state parameters for CSRF protection
type StateStore struct {
	mu     sync.RWMutex
	states map[string]*StateData // state -> state data
}

// NewStateStore creates a new state store
func NewStateStore() *StateStore {
	return &StateStore{
		states: make(map[string]*StateData),
	}
}

// Set stores a state with associated metadata
func (s *StateStore) Set(state string, data *StateData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[state] = data
}

// Get retrieves and removes a state's data
func (s *StateStore) Get(state string) (*StateData, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, ok := s.states[state]
	if ok {
		delete(s.states, state) // State is single-use
	}
	return data, ok
}

// NewManager creates a new OAuth2 manager with lazy loading from service-service
func NewManager(serviceServiceURL string) *Manager {
	return &Manager{
		configLoader: NewConfigLoader(serviceServiceURL),
		providers:    make(map[string]*Provider),
		states:       NewStateStore(),
	}
}

// getOrLoadProvider gets a provider from cache or loads it from service-service
func (m *Manager) getOrLoadProvider(providerName string) (*Provider, error) {
	m.providersMu.RLock()
	if provider, exists := m.providers[providerName]; exists {
		m.providersMu.RUnlock()
		return provider, nil
	}
	m.providersMu.RUnlock()

	m.providersMu.Lock()
	defer m.providersMu.Unlock()

	if provider, exists := m.providers[providerName]; exists {
		return provider, nil
	}

	config, err := m.configLoader.GetProvider(providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to load provider config: %w", err)
	}

	provider := NewProvider(*config)
	m.providers[providerName] = provider
	return provider, nil
}

// GetAuthURL generates an OAuth2 authorization URL for a specific provider with metadata
func (m *Manager) GetAuthURL(providerName string, userID int, callbackURL string, platform string) (string, error) {
	provider, err := m.getOrLoadProvider(providerName)
	if err != nil {
		return "", err
	}

	state, err := GenerateState()
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	m.states.Set(state, &StateData{
		Provider:    providerName,
		UserID:      userID,
		CallbackURL: callbackURL,
		Platform:    platform,
	})

	authURL := provider.GenerateAuthURL(state, callbackURL)
	return authURL, nil
}

// HandleCallback handles the OAuth2 callback with code and state, returns StateData
func (m *Manager) HandleCallback(state, code string) (*UserInfo, *TokenResponse, *StateData, error) {
	// Validate state and get associated metadata
	stateData, ok := m.states.Get(state)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid or expired state parameter")
	}

	provider, err := m.getOrLoadProvider(stateData.Provider)
	if err != nil {
		return nil, nil, nil, err
	}

	redirectURI := stateData.CallbackURL
	if redirectURI == "" {
		publicURL := os.Getenv("PUBLIC_URL")
		if publicURL == "" {
			publicURL = "http://localhost:8083"
		}
		redirectURI = publicURL + "/oauth2/callback"
	}

	tokenResp, err := provider.ExchangeCodeWithRedirect(code, redirectURI)
	if err != nil {
		return nil, nil, nil,
			fmt.Errorf("failed to exchange code: %w", err)
	}

	userInfo, err := provider.GetUserInfo(tokenResp.AccessToken)
	if err != nil {
		return nil, nil, nil,
			fmt.Errorf("failed to get user info: %w", err)
	}

	return userInfo, tokenResp, stateData, nil
}

// ListProviders returns all available provider names
func (m *Manager) ListProviders() ([]string, error) {
	return m.configLoader.ListProviders()
}

// GetProvider returns a specific provider (lazy loaded)
func (m *Manager) GetProvider(name string) (*Provider, error) {
	return m.getOrLoadProvider(name)
}