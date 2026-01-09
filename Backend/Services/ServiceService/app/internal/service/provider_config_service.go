package service

import (
	"github.com/raphael-guer1n/AREA/ServiceService/internal/config"
)

type ProviderConfigService struct {
	providers map[string]config.ProviderConfig
	services  map[string]config.ServiceConfig
}

func NewProviderConfigService(providersDir string, servicesDir string) (*ProviderConfigService, error) {
	providers, err := config.LoadProviderConfigs(providersDir)
	if err != nil {
		return nil, err
	}
	services, err := config.LoadServiceConfig(servicesDir)
	if err != nil {
		return nil, err
	}
	return &ProviderConfigService{providers: providers, services: services}, nil
}

// GetAllProvidersNames returns a list of all available service provider names
func (s *ProviderConfigService) GetAllProvidersNames() []string {
	names := make([]string, 0, len(s.providers))
	for name := range s.providers {
		names = append(names, name)
	}
	return names
}

// GetOAuth2Config returns the OAuth2 configuration for a specific service
func (s *ProviderConfigService) GetOAuth2Config(serviceName string) (*config.OAuth2Config, bool) {
	provider, exists := s.providers[serviceName]
	if !exists {
		return nil, false
	}
	return &provider.OAuth2, true
}

// GetProviderConfig returns the full provider configuration including mappings
func (s *ProviderConfigService) GetProviderConfig(serviceName string) (*config.ProviderConfig, bool) {
	provider, exists := s.providers[serviceName]
	if !exists {
		return nil, false
	}
	return &provider, true
}

func (s *ProviderConfigService) GetServiceConfig(serviceName string) (*config.ServiceConfig, bool) {
	service, exists := s.services[serviceName]
	if !exists {
		return nil, false
	}
	return &service, true
}

func (s *ProviderConfigService) GetAllServicesNames() []string {
	names := make([]string, 0, len(s.services))
	for name := range s.services {
		names = append(names, name)
	}
	return names
}
