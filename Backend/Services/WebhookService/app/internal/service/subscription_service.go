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
	ErrActionExists         = errors.New("action already has a webhook")
	ErrUnauthorizedAction   = errors.New("action does not belong to user")
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

func (s *SubscriptionService) CreateSubscription(userID, actionID int, provider, service string, cfg json.RawMessage, authToken string, active bool, webhookBaseURL string) (*domain.Subscription, error) {
	if existing, err := s.repo.FindByActionID(actionID); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, ErrActionExists
	}

	serviceName := strings.TrimSpace(service)
	if serviceName == "" {
		serviceName = strings.TrimSpace(provider)
	}
	if serviceName == "" {
		return nil, ErrProviderNotSupported
	}
	providerName := strings.TrimSpace(provider)
	if providerName == "" {
		providerName = serviceName
	}

	providerConfig, err := s.providerConfig.GetProviderConfig(serviceName)
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
		UserID:    userID,
		ActionID:  actionID,
		Provider:  providerName,
		Service:   serviceName,
		AuthToken: authToken,
		Active:    active,
		Config:    cfg,
	}

	for i := 0; i < 3; i++ {
		sub.HookID = generateHookID()
		created, err := s.repo.Create(sub)
		if err == nil {
			if active && s.webhookSetupSvc != nil && providerConfig.Setup != nil {
				webhookURL := buildWebhookURL(webhookBaseURL, serviceName, created.HookID)
				providerHookID, err := s.webhookSetupSvc.RegisterWebhook(providerConfig, created, webhookURL, cfgMap)
				if err != nil {
					_ = s.repo.DeleteByActionID(created.ActionID)
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

func (s *SubscriptionService) GetSubscriptionByActionID(actionID int) (*domain.Subscription, error) {
	return s.repo.FindByActionID(actionID)
}

func (s *SubscriptionService) UpdateSubscription(userID, actionID int, provider, service string, cfg json.RawMessage, authToken string, active bool, webhookBaseURL string) (*domain.Subscription, error) {
	subscription, err := s.repo.FindByActionID(actionID)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return nil, ErrSubscriptionNotFound
	}
	if subscription.UserID != userID {
		return nil, ErrUnauthorizedAction
	}

	serviceName := strings.TrimSpace(service)
	if serviceName == "" {
		serviceName = strings.TrimSpace(provider)
	}
	if serviceName == "" {
		return nil, ErrProviderNotSupported
	}
	providerName := strings.TrimSpace(provider)
	if providerName == "" {
		providerName = serviceName
	}

	newProviderConfig, err := s.providerConfig.GetProviderConfig(serviceName)
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

	cfgMap, err = s.applyPrepareSteps(userID, newProviderConfig, cfgMap)
	if err != nil {
		return nil, err
	}

	if newProviderConfig.Signature != nil {
		secretValue, ok := utils.ExtractJSONPath(cfgMap, newProviderConfig.Signature.SecretJSONPath)
		if !ok || fmt.Sprint(secretValue) == "" {
			return nil, ErrMissingSecret
		}
	}

	newConfig, err := json.Marshal(cfgMap)
	if err != nil {
		return nil, ErrInvalidConfig
	}

	if subscription.Active && s.webhookSetupSvc != nil {
		oldProviderConfig, err := s.providerConfig.GetProviderConfig(subscription.Service)
		if err == nil && oldProviderConfig.Teardown != nil {
			var oldCfgPayload any = map[string]any{}
			if len(subscription.Config) > 0 {
				if err := json.Unmarshal(subscription.Config, &oldCfgPayload); err == nil {
					webhookURL := buildWebhookURL(webhookBaseURL, subscription.Service, subscription.HookID)
					_ = s.webhookSetupSvc.DeleteWebhook(oldProviderConfig, subscription, webhookURL, oldCfgPayload)
				}
			}
		}
	}

	updatedSub := *subscription
	updatedSub.Provider = providerName
	updatedSub.Service = serviceName
	updatedSub.Config = newConfig
	updatedSub.AuthToken = authToken
	updatedSub.Active = active
	updatedSub.ProviderHookID = subscription.ProviderHookID

	if active && s.webhookSetupSvc != nil && newProviderConfig.Setup != nil {
		webhookURL := buildWebhookURL(webhookBaseURL, serviceName, updatedSub.HookID)
		providerHookID, err := s.webhookSetupSvc.RegisterWebhook(newProviderConfig, &updatedSub, webhookURL, cfgMap)
		if err != nil {
			return nil, err
		}
		updatedSub.ProviderHookID = providerHookID
	} else {
		updatedSub.ProviderHookID = ""
	}

	saved, err := s.repo.UpdateByActionID(&updatedSub)
	if err != nil {
		return nil, err
	}
	if saved == nil {
		return nil, ErrSubscriptionNotFound
	}
	return saved, nil
}

func (s *SubscriptionService) DeleteSubscription(actionID int, webhookBaseURL string) error {
	subscription, err := s.repo.FindByActionID(actionID)
	if err != nil {
		return err
	}
	if subscription == nil {
		return ErrSubscriptionNotFound
	}

	providerConfig, err := s.providerConfig.GetProviderConfig(subscription.Service)
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

	if subscription.Active && s.webhookSetupSvc != nil && providerConfig.Teardown != nil {
		webhookURL := buildWebhookURL(webhookBaseURL, subscription.Service, subscription.HookID)
		if err := s.webhookSetupSvc.DeleteWebhook(providerConfig, subscription, webhookURL, cfgPayload); err != nil {
			return err
		}
	}

	return s.repo.DeleteByActionID(actionID)
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
