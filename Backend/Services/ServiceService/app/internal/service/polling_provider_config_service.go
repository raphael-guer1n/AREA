package service

import "github.com/raphael-guer1n/AREA/ServiceService/internal/config"

type PollingProviderConfigService struct {
	providers map[string]config.PollingProviderConfig
}

func NewPollingProviderConfigService(providersDir string) (*PollingProviderConfigService, error) {
	providers, err := config.LoadPollingProviderConfigs(providersDir)
	if err != nil {
		return nil, err
	}
	return &PollingProviderConfigService{providers: providers}, nil
}

func (s *PollingProviderConfigService) GetAllProviderNames() []string {
	names := make([]string, 0, len(s.providers))
	for name := range s.providers {
		names = append(names, name)
	}
	return names
}

func (s *PollingProviderConfigService) GetProviderConfig(name string) (*config.PollingProviderConfig, bool) {
	cfg, ok := s.providers[name]
	if !ok {
		return nil, false
	}
	return &cfg, true
}
