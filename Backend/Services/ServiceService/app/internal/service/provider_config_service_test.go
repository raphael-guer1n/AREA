package service

import (
	"testing"

	"github.com/raphael-guer1n/AREA/ServiceService/internal/config"
	"github.com/stretchr/testify/assert"
)

// Note: These tests assume you have a way to mock or provide test provider configs
// In a real scenario, you would create a test fixtures directory with sample config files

func TestProviderConfigService_GetAllProvidersNames(t *testing.T) {
	// Create a mock service with test data
	svc := &ProviderConfigService{
		providers: map[string]config.ProviderConfig{
			"google":  {},
			"github":  {},
			"discord": {},
		},
		services: map[string]config.ServiceConfig{},
	}

	names := svc.GetAllProvidersNames()

	assert.Len(t, names, 3)
	assert.Contains(t, names, "google")
	assert.Contains(t, names, "github")
	assert.Contains(t, names, "discord")
}

func TestProviderConfigService_GetAllProvidersNames_Empty(t *testing.T) {
	svc := &ProviderConfigService{
		providers: map[string]config.ProviderConfig{},
		services:  map[string]config.ServiceConfig{},
	}

	names := svc.GetAllProvidersNames()

	assert.Empty(t, names)
	assert.NotNil(t, names)
}

func TestProviderConfigService_GetAllProviderSummaries(t *testing.T) {
	svc := &ProviderConfigService{
		providers: map[string]config.ProviderConfig{
			"google": {
				LogoURL: "https://example.com/google.png",
			},
			"github": {
				LogoURL: "https://example.com/github.png",
			},
			"discord": {
				LogoURL: "https://example.com/discord.png",
			},
		},
		services: map[string]config.ServiceConfig{},
	}

	summaries := svc.GetAllProviderSummaries()

	assert.Len(t, summaries, 3)

	// Check that results are sorted by name
	for i := 0; i < len(summaries)-1; i++ {
		assert.True(t, summaries[i].Name < summaries[i+1].Name)
	}

	// Find specific provider
	var discordSummary *ProviderSummary
	for i := range summaries {
		if summaries[i].Name == "discord" {
			discordSummary = &summaries[i]
			break
		}
	}

	assert.NotNil(t, discordSummary)
	assert.Equal(t, "discord", discordSummary.Name)
	assert.Equal(t, "https://example.com/discord.png", discordSummary.LogoURL)
}

func TestProviderConfigService_GetAllProviderSummaries_Empty(t *testing.T) {
	svc := &ProviderConfigService{
		providers: map[string]config.ProviderConfig{},
		services:  map[string]config.ServiceConfig{},
	}

	summaries := svc.GetAllProviderSummaries()

	assert.Empty(t, summaries)
	assert.NotNil(t, summaries)
}

func TestProviderConfigService_GetOAuth2Config_Success(t *testing.T) {
	expectedOAuth := config.OAuth2Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "write"},
	}

	svc := &ProviderConfigService{
		providers: map[string]config.ProviderConfig{
			"google": {
				OAuth2: expectedOAuth,
			},
		},
		services: map[string]config.ServiceConfig{},
	}

	oauth2Config, exists := svc.GetOAuth2Config("google")

	assert.True(t, exists)
	assert.NotNil(t, oauth2Config)
	assert.Equal(t, expectedOAuth.ClientID, oauth2Config.ClientID)
	assert.Equal(t, expectedOAuth.ClientSecret, oauth2Config.ClientSecret)
	assert.Equal(t, expectedOAuth.RedirectURI, oauth2Config.RedirectURI)
	assert.Equal(t, expectedOAuth.Scopes, oauth2Config.Scopes)
}

func TestProviderConfigService_GetOAuth2Config_NotFound(t *testing.T) {
	svc := &ProviderConfigService{
		providers: map[string]config.ProviderConfig{
			"google": {},
		},
		services: map[string]config.ServiceConfig{},
	}

	oauth2Config, exists := svc.GetOAuth2Config("nonexistent")

	assert.False(t, exists)
	assert.Nil(t, oauth2Config)
}

func TestProviderConfigService_GetProviderConfig_Success(t *testing.T) {
	expectedConfig := config.ProviderConfig{
		LogoURL: "https://example.com/logo.png",
		OAuth2: config.OAuth2Config{
			ClientID: "test-id",
		},
	}

	svc := &ProviderConfigService{
		providers: map[string]config.ProviderConfig{
			"github": expectedConfig,
		},
		services: map[string]config.ServiceConfig{},
	}

	providerConfig, exists := svc.GetProviderConfig("github")

	assert.True(t, exists)
	assert.NotNil(t, providerConfig)
	assert.Equal(t, expectedConfig.LogoURL, providerConfig.LogoURL)
	assert.Equal(t, expectedConfig.OAuth2.ClientID, providerConfig.OAuth2.ClientID)
}

func TestProviderConfigService_GetProviderConfig_NotFound(t *testing.T) {
	svc := &ProviderConfigService{
		providers: map[string]config.ProviderConfig{},
		services:  map[string]config.ServiceConfig{},
	}

	providerConfig, exists := svc.GetProviderConfig("nonexistent")

	assert.False(t, exists)
	assert.Nil(t, providerConfig)
}

func TestProviderConfigService_GetServiceConfig_Success(t *testing.T) {
	expectedConfig := config.ServiceConfig{
		Name:     "test-service",
		Label:    "Test Service",
		Provider: "test-provider",
	}

	svc := &ProviderConfigService{
		providers: map[string]config.ProviderConfig{},
		services: map[string]config.ServiceConfig{
			"test-service": expectedConfig,
		},
	}

	serviceConfig, exists := svc.GetServiceConfig("test-service")

	assert.True(t, exists)
	assert.NotNil(t, serviceConfig)
	assert.Equal(t, expectedConfig.Name, serviceConfig.Name)
	assert.Equal(t, expectedConfig.Label, serviceConfig.Label)
}

func TestProviderConfigService_GetServiceConfig_NotFound(t *testing.T) {
	svc := &ProviderConfigService{
		providers: map[string]config.ProviderConfig{},
		services:  map[string]config.ServiceConfig{},
	}

	serviceConfig, exists := svc.GetServiceConfig("nonexistent")

	assert.False(t, exists)
	assert.Nil(t, serviceConfig)
}

func TestProviderConfigService_GetAllServicesNames(t *testing.T) {
	svc := &ProviderConfigService{
		providers: map[string]config.ProviderConfig{},
		services: map[string]config.ServiceConfig{
			"service1": {},
			"service2": {},
			"service3": {},
		},
	}

	names := svc.GetAllServicesNames()

	assert.Len(t, names, 3)
	assert.Contains(t, names, "service1")
	assert.Contains(t, names, "service2")
	assert.Contains(t, names, "service3")
}

func TestProviderConfigService_GetAllServicesNames_Empty(t *testing.T) {
	svc := &ProviderConfigService{
		providers: map[string]config.ProviderConfig{},
		services:  map[string]config.ServiceConfig{},
	}

	names := svc.GetAllServicesNames()

	assert.Empty(t, names)
	assert.NotNil(t, names)
}

func TestProviderConfigService_MultipleProviders(t *testing.T) {
	svc := &ProviderConfigService{
		providers: map[string]config.ProviderConfig{
			"google": {
				LogoURL: "https://example.com/google.png",
				OAuth2: config.OAuth2Config{
					ClientID: "google-client-id",
				},
			},
			"github": {
				LogoURL: "https://example.com/github.png",
				OAuth2: config.OAuth2Config{
					ClientID: "github-client-id",
				},
			},
		},
		services: map[string]config.ServiceConfig{},
	}

	// Test getting multiple providers
	googleConfig, googleExists := svc.GetProviderConfig("google")
	githubConfig, githubExists := svc.GetProviderConfig("github")

	assert.True(t, googleExists)
	assert.True(t, githubExists)
	assert.Equal(t, "google-client-id", googleConfig.OAuth2.ClientID)
	assert.Equal(t, "github-client-id", githubConfig.OAuth2.ClientID)

	// Test getting all names includes both
	names := svc.GetAllProvidersNames()
	assert.Contains(t, names, "google")
	assert.Contains(t, names, "github")
}

func TestProviderConfigService_CaseSensitivity(t *testing.T) {
	svc := &ProviderConfigService{
		providers: map[string]config.ProviderConfig{
			"GitHub": {
				LogoURL: "https://example.com/github.png",
			},
		},
		services: map[string]config.ServiceConfig{},
	}

	// Keys are case-sensitive
	_, existsCorrectCase := svc.GetProviderConfig("GitHub")
	_, existsLowerCase := svc.GetProviderConfig("github")

	assert.True(t, existsCorrectCase)
	assert.False(t, existsLowerCase)
}

func TestProviderSummary_Structure(t *testing.T) {
	summary := ProviderSummary{
		Name:    "test-provider",
		LogoURL: "https://example.com/logo.png",
	}

	assert.Equal(t, "test-provider", summary.Name)
	assert.Equal(t, "https://example.com/logo.png", summary.LogoURL)
}

func TestProviderConfigService_EmptyLogoURL(t *testing.T) {
	svc := &ProviderConfigService{
		providers: map[string]config.ProviderConfig{
			"test": {
				LogoURL: "",
			},
		},
		services: map[string]config.ServiceConfig{},
	}

	summaries := svc.GetAllProviderSummaries()

	assert.Len(t, summaries, 1)
	assert.Equal(t, "test", summaries[0].Name)
	assert.Equal(t, "", summaries[0].LogoURL)
}

func TestProviderConfigService_SummariesSorted(t *testing.T) {
	svc := &ProviderConfigService{
		providers: map[string]config.ProviderConfig{
			"zebra":  {},
			"apple":  {},
			"monkey": {},
			"banana": {},
		},
		services: map[string]config.ServiceConfig{},
	}

	summaries := svc.GetAllProviderSummaries()

	assert.Len(t, summaries, 4)
	assert.Equal(t, "apple", summaries[0].Name)
	assert.Equal(t, "banana", summaries[1].Name)
	assert.Equal(t, "monkey", summaries[2].Name)
	assert.Equal(t, "zebra", summaries[3].Name)
}
