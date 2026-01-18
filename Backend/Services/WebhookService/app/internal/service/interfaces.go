package service

import (
	"github.com/raphael-guer1n/AREA/WebhookService/internal/config"
	"github.com/raphael-guer1n/AREA/WebhookService/internal/domain"
	"github.com/raphael-guer1n/AREA/WebhookService/internal/utils"
)

// ProviderConfigServiceInterface defines the interface for provider config service
type ProviderConfigServiceInterface interface {
	GetProviderConfig(name string) (*config.WebhookProviderConfig, error)
}

// WebhookSetupServiceInterface defines the interface for webhook setup service
type WebhookSetupServiceInterface interface {
	RegisterWebhook(providerConfig *config.WebhookProviderConfig, sub *domain.Subscription, webhookURL string, subscriptionConfig any) (string, error)
	DeleteWebhook(providerConfig *config.WebhookProviderConfig, sub *domain.Subscription, webhookURL string, subscriptionConfig any) error
	executeActionOnce(action *config.WebhookProviderSetupConfig, provider string, userID int, ctx utils.TemplateContext, label string, queryOverrides map[string]string) ([]byte, error)
}
