package service

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/raphael-guer1n/AREA/WebhookService/internal/config"
	"github.com/raphael-guer1n/AREA/WebhookService/internal/domain"
)

type SubscriptionRenewalService struct {
	repo            domain.SubscriptionRepository
	providerConfig  *ProviderConfigService
	webhookSetupSvc *WebhookSetupService
	publicBaseURL   string
	pollInterval    time.Duration
}

func NewSubscriptionRenewalService(repo domain.SubscriptionRepository, providerConfig *ProviderConfigService, webhookSetupSvc *WebhookSetupService, publicBaseURL string) *SubscriptionRenewalService {
	return &SubscriptionRenewalService{
		repo:            repo,
		providerConfig:  providerConfig,
		webhookSetupSvc: webhookSetupSvc,
		publicBaseURL:   strings.TrimRight(publicBaseURL, "/"),
		pollInterval:    time.Minute,
	}
}

func (s *SubscriptionRenewalService) Start() {
	if s.webhookSetupSvc == nil || s.providerConfig == nil || s.repo == nil {
		return
	}
	if s.publicBaseURL == "" {
		log.Printf("subscription renewal disabled: public base URL is empty")
		return
	}

	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	for {
		s.runOnce()
		<-ticker.C
	}
}

func (s *SubscriptionRenewalService) runOnce() {
	providers, err := s.providerConfig.ListProviders()
	if err != nil {
		log.Printf("subscription renewal: failed to list providers: %v", err)
		return
	}

	for _, provider := range providers {
		cfg, err := s.providerConfig.GetProviderConfig(provider)
		if err != nil {
			log.Printf("subscription renewal: failed to load provider config %s: %v", provider, err)
			continue
		}
		if !shouldRenew(cfg) {
			continue
		}
		s.renewProvider(provider, cfg)
	}
}

func shouldRenew(cfg *config.WebhookProviderConfig) bool {
	if cfg == nil || cfg.Renewal == nil {
		return false
	}
	return cfg.Renewal.AfterSeconds > 0
}

func (s *SubscriptionRenewalService) renewProvider(provider string, cfg *config.WebhookProviderConfig) {
	subs, err := s.repo.ListByProvider(provider)
	if err != nil {
		log.Printf("subscription renewal: failed to list subscriptions for %s: %v", provider, err)
		return
	}

	interval := time.Duration(cfg.Renewal.AfterSeconds) * time.Second
	for _, sub := range subs {
		last := sub.UpdatedAt
		if last.IsZero() {
			last = sub.CreatedAt
		}
		if time.Since(last) < interval {
			continue
		}

		var cfgPayload any = map[string]any{}
		if len(sub.Config) > 0 {
			if err := json.Unmarshal(sub.Config, &cfgPayload); err != nil {
				log.Printf("subscription renewal: invalid config for hook_id=%s: %v", sub.HookID, err)
				continue
			}
		}

		webhookURL := buildWebhookURL(s.publicBaseURL, provider, sub.HookID)
		providerHookID, err := s.webhookSetupSvc.RegisterWebhook(cfg, &sub, webhookURL, cfgPayload)
		if err != nil {
			log.Printf("subscription renewal: failed for hook_id=%s provider=%s: %v", sub.HookID, provider, err)
			continue
		}
		if providerHookID != "" && providerHookID != sub.ProviderHookID {
			if err := s.repo.UpdateProviderHookID(sub.HookID, providerHookID); err != nil {
				log.Printf("subscription renewal: failed to update provider_hook_id for hook_id=%s: %v", sub.HookID, err)
				continue
			}
			continue
		}
		if err := s.repo.TouchByHookID(sub.HookID); err != nil {
			log.Printf("subscription renewal: failed to update timestamp for hook_id=%s: %v", sub.HookID, err)
		}
	}
}
