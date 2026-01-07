package service

import "github.com/raphael-guer1n/AREA/ServiceService/internal/config"

type WebhookProviderConfigService struct {
	providers map[string]config.WebhookProviderConfig
}

func NewWebhookProviderConfigService(providersDir string) (*WebhookProviderConfigService, error) {
	providers, err := config.LoadWebhookProviderConfigs(providersDir)
	if err != nil {
		return nil, err
	}
	return &WebhookProviderConfigService{providers: providers}, nil
}

// GetAllProviderNames returns a list of all available webhook provider names.
func (s *WebhookProviderConfigService) GetAllProviderNames() []string {
	names := make([]string, 0, len(s.providers))
	for name := range s.providers {
		names = append(names, name)
	}
	return names
}

// GetProviderConfig returns the full webhook provider configuration.
func (s *WebhookProviderConfigService) GetProviderConfig(name string) (*config.WebhookProviderConfig, bool) {
	provider, exists := s.providers[name]
	if !exists {
		return nil, false
	}
	return &provider, true
}
