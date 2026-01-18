package service

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/raphael-guer1n/AREA/PollingService/internal/config"
	"github.com/raphael-guer1n/AREA/PollingService/internal/domain"
	"github.com/raphael-guer1n/AREA/PollingService/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSubscriptionRepository is a mock implementation of SubscriptionRepository
type MockSubscriptionRepository struct {
	mock.Mock
}

func (m *MockSubscriptionRepository) Create(sub *domain.Subscription) (*domain.Subscription, error) {
	args := m.Called(sub)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) FindByActionID(actionID int) (*domain.Subscription, error) {
	args := m.Called(actionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) ListDue(now time.Time) ([]domain.Subscription, error) {
	args := m.Called(now)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) UpdateByActionID(sub *domain.Subscription) (*domain.Subscription, error) {
	args := m.Called(sub)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) UpdatePollingState(actionID int, lastItemID string, nextRunAt time.Time, lastError string, lastPolledAt time.Time) error {
	args := m.Called(actionID, lastItemID, nextRunAt, lastError, lastPolledAt)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) DeleteByActionID(actionID int) error {
	args := m.Called(actionID)
	return args.Error(0)
}

// MockProviderConfigService is a mock for provider config
type MockProviderConfigService struct {
	mock.Mock
}

func (m *MockProviderConfigService) GetProviderConfig(serviceName string) (*config.PollingProviderConfig, error) {
	args := m.Called(serviceName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*config.PollingProviderConfig), args.Error(1)
}

// MockRequestService is a mock for request service
type MockRequestService struct {
	mock.Mock
}

func (m *MockRequestService) ExecuteRequest(request config.PollingProviderRequestConfig, provider string, userID int, ctx utils.TemplateContext, queryOverrides map[string]string) ([]byte, error) {
	args := m.Called(request, provider, userID, ctx, queryOverrides)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func TestSubscriptionService_CreateSubscription_Success(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProviderConfig := new(MockProviderConfigService)
	mockRequestSvc := new(MockRequestService)

	svc := NewSubscriptionService(mockRepo, mockProviderConfig, mockRequestSvc)

	userID := 1
	actionID := 100
	provider := "github"
	service := "github"
	cfg := json.RawMessage(`{"repo": "test/repo"}`)
	active := true

	providerCfg := &config.PollingProviderConfig{
		IntervalSeconds: 300,
		Request: config.PollingProviderRequestConfig{
			Method:      "GET",
			URLTemplate: "https://api.github.com/repos/{{repo}}/events",
		},
	}

	mockRepo.On("FindByActionID", actionID).Return(nil, nil)
	mockProviderConfig.On("GetProviderConfig", service).Return(providerCfg, nil)
	mockRepo.On("Create", mock.AnythingOfType("*domain.Subscription")).Return(&domain.Subscription{
		ID:              1,
		UserID:          userID,
		ActionID:        actionID,
		Provider:        provider,
		Service:         service,
		Active:          active,
		Config:          cfg,
		IntervalSeconds: 300,
	}, nil)

	sub, err := svc.CreateSubscription(userID, actionID, provider, service, cfg, active)

	assert.NoError(t, err)
	assert.NotNil(t, sub)
	assert.Equal(t, actionID, sub.ActionID)
	assert.Equal(t, userID, sub.UserID)
	mockRepo.AssertExpectations(t)
	mockProviderConfig.AssertExpectations(t)
}

func TestSubscriptionService_CreateSubscription_ActionAlreadyExists(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProviderConfig := new(MockProviderConfigService)
	mockRequestSvc := new(MockRequestService)

	svc := NewSubscriptionService(mockRepo, mockProviderConfig, mockRequestSvc)

	actionID := 100
	existingSub := &domain.Subscription{
		ID:       1,
		ActionID: actionID,
	}

	mockRepo.On("FindByActionID", actionID).Return(existingSub, nil)

	sub, err := svc.CreateSubscription(1, actionID, "github", "github", json.RawMessage(`{}`), true)

	assert.Error(t, err)
	assert.Nil(t, sub)
	assert.ErrorIs(t, err, ErrActionExists)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_CreateSubscription_ProviderNotSupported(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProviderConfig := new(MockProviderConfigService)
	mockRequestSvc := new(MockRequestService)

	svc := NewSubscriptionService(mockRepo, mockProviderConfig, mockRequestSvc)

	mockRepo.On("FindByActionID", 100).Return(nil, nil)
	mockProviderConfig.On("GetProviderConfig", "unknown").Return(nil, ErrProviderConfigNotFound)

	sub, err := svc.CreateSubscription(1, 100, "unknown", "unknown", json.RawMessage(`{}`), true)

	assert.Error(t, err)
	assert.Nil(t, sub)
	assert.ErrorIs(t, err, ErrProviderNotSupported)
	mockRepo.AssertExpectations(t)
	mockProviderConfig.AssertExpectations(t)
}

func TestSubscriptionService_CreateSubscription_InvalidConfig(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProviderConfig := new(MockProviderConfigService)
	mockRequestSvc := new(MockRequestService)

	svc := NewSubscriptionService(mockRepo, mockProviderConfig, mockRequestSvc)

	providerCfg := &config.PollingProviderConfig{
		IntervalSeconds: 300,
		Request: config.PollingProviderRequestConfig{
			Method:      "GET",
			URLTemplate: "https://api.example.com",
		},
	}

	mockRepo.On("FindByActionID", 100).Return(nil, nil)
	mockProviderConfig.On("GetProviderConfig", "test").Return(providerCfg, nil)

	// Invalid JSON
	invalidCfg := json.RawMessage(`{invalid json}`)

	sub, err := svc.CreateSubscription(1, 100, "test", "test", invalidCfg, true)

	assert.Error(t, err)
	assert.Nil(t, sub)
	assert.ErrorIs(t, err, ErrInvalidConfig)
	mockRepo.AssertExpectations(t)
	mockProviderConfig.AssertExpectations(t)
}

func TestSubscriptionService_GetSubscriptionByActionID_Success(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProviderConfig := new(MockProviderConfigService)
	mockRequestSvc := new(MockRequestService)

	svc := NewSubscriptionService(mockRepo, mockProviderConfig, mockRequestSvc)

	actionID := 100
	expectedSub := &domain.Subscription{
		ID:       1,
		ActionID: actionID,
		UserID:   1,
	}

	mockRepo.On("FindByActionID", actionID).Return(expectedSub, nil)

	sub, err := svc.GetSubscriptionByActionID(actionID)

	assert.NoError(t, err)
	assert.Equal(t, expectedSub, sub)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_DeleteSubscription_Success(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProviderConfig := new(MockProviderConfigService)
	mockRequestSvc := new(MockRequestService)

	svc := NewSubscriptionService(mockRepo, mockProviderConfig, mockRequestSvc)

	actionID := 100

	mockRepo.On("DeleteByActionID", actionID).Return(nil)

	err := svc.DeleteSubscription(actionID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_DeleteSubscription_Error(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProviderConfig := new(MockProviderConfigService)
	mockRequestSvc := new(MockRequestService)

	svc := NewSubscriptionService(mockRepo, mockProviderConfig, mockRequestSvc)

	actionID := 100
	dbError := errors.New("database error")

	mockRepo.On("DeleteByActionID", actionID).Return(dbError)

	err := svc.DeleteSubscription(actionID)

	assert.Error(t, err)
	assert.Equal(t, dbError, err)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_ActivateSubscription_Success(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProviderConfig := new(MockProviderConfigService)
	mockRequestSvc := new(MockRequestService)

	svc := NewSubscriptionService(mockRepo, mockProviderConfig, mockRequestSvc)

	userID := 1
	actionID := 100

	existingSub := &domain.Subscription{
		ID:       1,
		UserID:   userID,
		ActionID: actionID,
		Active:   false,
		Service:  "github",
	}

	mockRepo.On("FindByActionID", actionID).Return(existingSub, nil)
	mockRepo.On("UpdateByActionID", mock.MatchedBy(func(sub *domain.Subscription) bool {
		return sub.ActionID == actionID && sub.Active == true
	})).Return(&domain.Subscription{
		ID:       1,
		UserID:   userID,
		ActionID: actionID,
		Active:   true,
	}, nil)

	sub, err := svc.ActivateSubscription(userID, actionID)

	assert.NoError(t, err)
	assert.NotNil(t, sub)
	assert.True(t, sub.Active)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_ActivateSubscription_UnauthorizedAction(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProviderConfig := new(MockProviderConfigService)
	mockRequestSvc := new(MockRequestService)

	svc := NewSubscriptionService(mockRepo, mockProviderConfig, mockRequestSvc)

	userID := 1
	actionID := 100

	existingSub := &domain.Subscription{
		ID:       1,
		UserID:   2, // Different user
		ActionID: actionID,
		Active:   false,
	}

	mockRepo.On("FindByActionID", actionID).Return(existingSub, nil)

	sub, err := svc.ActivateSubscription(userID, actionID)

	assert.Error(t, err)
	assert.Nil(t, sub)
	assert.ErrorIs(t, err, ErrUnauthorizedAction)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_DeactivateSubscription_Success(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProviderConfig := new(MockProviderConfigService)
	mockRequestSvc := new(MockRequestService)

	svc := NewSubscriptionService(mockRepo, mockProviderConfig, mockRequestSvc)

	userID := 1
	actionID := 100

	existingSub := &domain.Subscription{
		ID:       1,
		UserID:   userID,
		ActionID: actionID,
		Active:   true,
		Service:  "github",
	}

	mockRepo.On("FindByActionID", actionID).Return(existingSub, nil)
	mockRepo.On("UpdateByActionID", mock.MatchedBy(func(sub *domain.Subscription) bool {
		return sub.ActionID == actionID && sub.Active == false
	})).Return(&domain.Subscription{
		ID:       1,
		UserID:   userID,
		ActionID: actionID,
		Active:   false,
	}, nil)

	sub, err := svc.DeactivateSubscription(userID, actionID)

	assert.NoError(t, err)
	assert.NotNil(t, sub)
	assert.False(t, sub.Active)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_UpdateSubscription_Success(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProviderConfig := new(MockProviderConfigService)
	mockRequestSvc := new(MockRequestService)

	svc := NewSubscriptionService(mockRepo, mockProviderConfig, mockRequestSvc)

	userID := 1
	actionID := 100
	newConfig := json.RawMessage(`{"new": "config"}`)

	existingSub := &domain.Subscription{
		ID:       1,
		UserID:   userID,
		ActionID: actionID,
		Service:  "github",
		Active:   false,
	}

	providerCfg := &config.PollingProviderConfig{
		IntervalSeconds: 300,
		Request: config.PollingProviderRequestConfig{
			Method:      "GET",
			URLTemplate: "https://api.github.com",
		},
	}

	mockRepo.On("FindByActionID", actionID).Return(existingSub, nil)
	mockProviderConfig.On("GetProviderConfig", "github").Return(providerCfg, nil)
	mockRepo.On("UpdateByActionID", mock.AnythingOfType("*domain.Subscription")).Return(&domain.Subscription{
		ID:       1,
		UserID:   userID,
		ActionID: actionID,
		Active:   true,
		Config:   newConfig,
	}, nil)

	sub, err := svc.UpdateSubscription(userID, actionID, "github", "github", newConfig, true)

	assert.NoError(t, err)
	assert.NotNil(t, sub)
	mockRepo.AssertExpectations(t)
	mockProviderConfig.AssertExpectations(t)
}

func TestSubscriptionService_UpdateSubscription_NotFound(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockProviderConfig := new(MockProviderConfigService)
	mockRequestSvc := new(MockRequestService)

	svc := NewSubscriptionService(mockRepo, mockProviderConfig, mockRequestSvc)

	actionID := 100

	mockRepo.On("FindByActionID", actionID).Return(nil, nil)

	sub, err := svc.UpdateSubscription(1, actionID, "github", "github", json.RawMessage(`{}`), true)

	assert.Error(t, err)
	assert.Nil(t, sub)
	assert.ErrorIs(t, err, ErrSubscriptionNotFound)
	mockRepo.AssertExpectations(t)
}

func TestIsUniqueViolation(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "unique violation error",
			err:      &pq.Error{Code: "23505"},
			expected: true,
		},
		{
			name:     "other pq error",
			err:      &pq.Error{Code: "23503"},
			expected: false,
		},
		{
			name:     "generic error",
			err:      errors.New("some error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isUniqueViolation(tc.err)
			assert.Equal(t, tc.expected, result)
		})
	}
}
