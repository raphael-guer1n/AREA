package service

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/raphael-guer1n/AREA/PollingService/internal/domain"
)

var (
	ErrProviderNotSupported = errors.New("provider not supported")
	ErrInvalidConfig        = errors.New("invalid subscription config")
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrActionExists         = errors.New("action already has a subscription")
	ErrUnauthorizedAction   = errors.New("action does not belong to user")
)

type SubscriptionService struct {
	repo           domain.SubscriptionRepository
	providerConfig ProviderConfigServiceInterface
	requestSvc     RequestServiceInterface
}

func NewSubscriptionService(repo domain.SubscriptionRepository, providerConfig ProviderConfigServiceInterface, requestSvc RequestServiceInterface) *SubscriptionService {
	return &SubscriptionService{
		repo:           repo,
		providerConfig: providerConfig,
		requestSvc:     requestSvc,
	}
}

func (s *SubscriptionService) CreateSubscription(userID, actionID int, provider, service string, cfg json.RawMessage, active bool) (*domain.Subscription, error) {
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
	if providerConfig.IntervalSeconds <= 0 || strings.TrimSpace(providerConfig.Request.Method) == "" || strings.TrimSpace(providerConfig.Request.URLTemplate) == "" {
		return nil, ErrInvalidConfig
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

	cfg, err = json.Marshal(cfgMap)
	if err != nil {
		return nil, ErrInvalidConfig
	}

	var nextRunAt *time.Time
	if active {
		now := time.Now().UTC()
		nextRunAt = &now
	}

	sub := &domain.Subscription{
		UserID:          userID,
		ActionID:        actionID,
		Provider:        providerName,
		Service:         serviceName,
		Active:          active,
		Config:          cfg,
		IntervalSeconds: providerConfig.IntervalSeconds,
		NextRunAt:       nextRunAt,
	}

	created, err := s.repo.Create(sub)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrActionExists
		}
		return nil, err
	}
	return created, nil
}

func (s *SubscriptionService) GetSubscriptionByActionID(actionID int) (*domain.Subscription, error) {
	return s.repo.FindByActionID(actionID)
}

func (s *SubscriptionService) UpdateSubscription(userID, actionID int, provider, service string, cfg json.RawMessage, active bool) (*domain.Subscription, error) {
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

	providerConfig, err := s.providerConfig.GetProviderConfig(serviceName)
	if err != nil {
		if errors.Is(err, ErrProviderConfigNotFound) {
			return nil, ErrProviderNotSupported
		}
		return nil, err
	}
	if providerConfig.IntervalSeconds <= 0 || strings.TrimSpace(providerConfig.Request.Method) == "" || strings.TrimSpace(providerConfig.Request.URLTemplate) == "" {
		return nil, ErrInvalidConfig
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

	newConfig, err := json.Marshal(cfgMap)
	if err != nil {
		return nil, ErrInvalidConfig
	}

	var nextRunAt *time.Time
	if active {
		now := time.Now().UTC()
		nextRunAt = &now
	}

	updatedSub := *subscription
	updatedSub.Provider = providerName
	updatedSub.Service = serviceName
	updatedSub.Config = newConfig
	updatedSub.Active = active
	updatedSub.IntervalSeconds = providerConfig.IntervalSeconds
	updatedSub.LastItemID = ""
	updatedSub.LastPolledAt = nil
	updatedSub.LastError = ""
	updatedSub.NextRunAt = nextRunAt

	saved, err := s.repo.UpdateByActionID(&updatedSub)
	if err != nil {
		return nil, err
	}
	if saved == nil {
		return nil, ErrSubscriptionNotFound
	}

	return saved, nil
}

func (s *SubscriptionService) DeleteSubscription(actionID int) error {
	if err := s.repo.DeleteByActionID(actionID); err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) ActivateSubscription(userID, actionID int) (*domain.Subscription, error) {
	return s.setSubscriptionActive(userID, actionID, true)
}

func (s *SubscriptionService) DeactivateSubscription(userID, actionID int) (*domain.Subscription, error) {
	return s.setSubscriptionActive(userID, actionID, false)
}

func (s *SubscriptionService) setSubscriptionActive(userID, actionID int, active bool) (*domain.Subscription, error) {
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

	updatedSub := *subscription
	updatedSub.Active = active
	if active {
		now := time.Now().UTC()
		updatedSub.NextRunAt = &now
	} else {
		updatedSub.NextRunAt = nil
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

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return false
}
