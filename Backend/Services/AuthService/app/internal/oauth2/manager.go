package oauth2

import (
	"fmt"
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
	// Check cache first
	m.providersMu.RLock()
	if provider, exists := m.providers[providerName]; exists {
		m.providersMu.RUnlock()
		return provider, nil
	}
	m.providersMu.RUnlock()

	// Load from service-service
	m.providersMu.Lock()
	defer m.providersMu.Unlock()

	// Double-check after acquiring write lock
	if provider, exists := m.providers[providerName]; exists {
		return provider, nil
	}

	// Fetch config from service-service
	config, err := m.configLoader.GetProvider(providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to load provider config: %w", err)
	}

	// Create and cache provider
	provider := NewProvider(*config)
	m.providers[providerName] = provider

	return provider, nil
}

// GetAuthURL generates an OAuth2 authorization URL for a specific provider with metadata
func (m *Manager) GetAuthURL(providerName string, userID int, callbackURL string, platform string) (string, error) {
	// Load provider on-demand
	provider, err := m.getOrLoadProvider(providerName)
	if err != nil {
		return "", err
	}

	// Generate CSRF state token
	state, err := GenerateState()
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	// Store state with metadata for validation
	m.states.Set(state, &StateData{
		Provider:    providerName,
		UserID:      userID,
		CallbackURL: callbackURL,
		Platform:    platform,
	})

	// Generate authorization URL
	authURL := provider.GenerateAuthURL(state, callbackURL)

	return authURL, nil
}

// HandleCallback handles the OAuth2 callback with code and state, returns StateData
func (m *Manager) HandleCallback(state, code string) (*UserInfo, *TokenResponse, *StateData, error) {
	// Validate state and get state data
	stateData, ok := m.states.Get(state)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid or expired state parameter")
	}

	// Load provider on-demand
	provider, err := m.getOrLoadProvider(stateData.Provider)
	if err != nil {
		return nil, nil, nil, err
	}

	// Exchange code for access token
	tokenResp, err := provider.ExchangeCode(code, stateData.CallbackURL)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info using access token
	userInfo, err := provider.GetUserInfo(tokenResp.AccessToken)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return userInfo, tokenResp, stateData, nil
}

// ListProviders returns all available provider names from service-service
func (m *Manager) ListProviders() ([]string, error) {
	return m.configLoader.ListProviders()
}

// GetProvider returns a specific provider (lazy loaded)
func (m *Manager) GetProvider(name string) (*Provider, error) {
	return m.getOrLoadProvider(name)
}
