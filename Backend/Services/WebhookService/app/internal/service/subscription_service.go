package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/raphael-guer1n/AREA/WebhookService/internal/domain"
	"github.com/raphael-guer1n/AREA/WebhookService/internal/utils"
)

var (
	ErrProviderNotSupported = errors.New("provider not supported")
	ErrInvalidConfig        = errors.New("invalid subscription config")
	ErrMissingSecret        = errors.New("missing signature secret")
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrProviderHookMissing  = errors.New("provider hook id is missing")
)

type SubscriptionService struct {
	repo            domain.SubscriptionRepository
	providerConfig  *ProviderConfigService
	webhookSetupSvc *WebhookSetupService
}

func NewSubscriptionService(repo domain.SubscriptionRepository, providerConfig *ProviderConfigService, webhookSetupSvc *WebhookSetupService) *SubscriptionService {
	return &SubscriptionService{
		repo:            repo,
		providerConfig:  providerConfig,
		webhookSetupSvc: webhookSetupSvc,
	}
}

func (s *SubscriptionService) CreateSubscription(userID, areaID int, provider string, cfg json.RawMessage, webhookBaseURL string) (*domain.Subscription, error) {
	providerConfig, err := s.providerConfig.GetProviderConfig(provider)
	if err != nil {
		if errors.Is(err, ErrProviderConfigNotFound) {
			return nil, ErrProviderNotSupported
		}
		return nil, err
	}

	var cfgPayload any = map[string]any{}
	if len(cfg) > 0 {
		if err := json.Unmarshal(cfg, &cfgPayload); err != nil {
			return nil, ErrInvalidConfig
		}
	}

	cfgMap, ok := cfgPayload.(map[string]any)
	if !ok {
		return nil, ErrInvalidConfig
	}

	cfgMap, err = s.applyPrepareSteps(userID, providerConfig, cfgMap)
	if err != nil {
		return nil, err
	}

	if providerConfig.Signature != nil {
		secretValue, ok := utils.ExtractJSONPath(cfgMap, providerConfig.Signature.SecretJSONPath)
		if !ok || fmt.Sprint(secretValue) == "" {
			return nil, ErrMissingSecret
		}
	}

	cfg, err = json.Marshal(cfgMap)
	if err != nil {
		return nil, ErrInvalidConfig
	}

	sub := &domain.Subscription{
		UserID:   userID,
		AreaID:   areaID,
		Provider: provider,
		Config:   cfg,
	}

	for i := 0; i < 3; i++ {
		sub.HookID = generateHookID()
		created, err := s.repo.Create(sub)
		if err == nil {
			if s.webhookSetupSvc != nil && providerConfig.Setup != nil {
				webhookURL := buildWebhookURL(webhookBaseURL, provider, created.HookID)
				providerHookID, err := s.webhookSetupSvc.RegisterWebhook(providerConfig, created, webhookURL, cfgMap)
				if err != nil {
					_ = s.repo.DeleteByHookID(created.HookID)
					return nil, err
				}
				if providerHookID != "" {
					if err := s.repo.UpdateProviderHookID(created.HookID, providerHookID); err != nil {
						return nil, err
					}
					created.ProviderHookID = providerHookID
				}
			}
			return created, nil
		}
		if isUniqueViolation(err) {
			continue
		}
		return nil, err
	}

	return nil, errors.New("failed to generate unique hook id")
}

func (s *SubscriptionService) GetSubscriptionByHookID(hookID string) (*domain.Subscription, error) {
	return s.repo.FindByHookID(hookID)
}

func (s *SubscriptionService) ListSubscriptionsByUserID(userID int) ([]domain.Subscription, error) {
	return s.repo.ListByUserID(userID)
}

func (s *SubscriptionService) DeleteSubscription(hookID, webhookBaseURL string) error {
	subscription, err := s.repo.FindByHookID(hookID)
	if err != nil {
		return err
	}
	if subscription == nil {
		return ErrSubscriptionNotFound
	}

	providerConfig, err := s.providerConfig.GetProviderConfig(subscription.Provider)
	if err != nil {
		if errors.Is(err, ErrProviderConfigNotFound) {
			return ErrProviderNotSupported
		}
		return err
	}

	var cfgPayload any = map[string]any{}
	if len(subscription.Config) > 0 {
		if err := json.Unmarshal(subscription.Config, &cfgPayload); err != nil {
			return ErrInvalidConfig
		}
	}
	if _, ok := cfgPayload.(map[string]any); !ok {
		return ErrInvalidConfig
	}

	if s.webhookSetupSvc != nil && providerConfig.Teardown != nil {
		webhookURL := buildWebhookURL(webhookBaseURL, subscription.Provider, subscription.HookID)
		if err := s.webhookSetupSvc.DeleteWebhook(providerConfig, subscription, webhookURL, cfgPayload); err != nil {
			return err
		}
	}

	return s.repo.DeleteByHookID(hookID)
}

func generateHookID() string {
	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}

func isUniqueViolation(err error) bool {
	pqErr, ok := err.(*pq.Error)
	if !ok {
		return false
	}
	return pqErr.Code == "23505"
}

func buildWebhookURL(baseURL, provider, hookID string) string {
	base := strings.TrimRight(baseURL, "/")
	if base == "" {
		return "/webhooks/" + provider + "/" + hookID
	}
	return base + "/webhooks/" + provider + "/" + hookID
}
