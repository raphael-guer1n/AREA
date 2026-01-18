package service

import (
	"errors"
	"testing"
	"time"

	"github.com/raphael-guer1n/AREA/CronService/internal/domain"
	"github.com/raphael-guer1n/AREA/CronService/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockActionRepository is a mock implementation of ActionRepositoryInterface
type MockActionRepository struct {
	mock.Mock
}

// Ensure MockActionRepository implements the interface
var _ repository.ActionRepositoryInterface = (*MockActionRepository)(nil)

func (m *MockActionRepository) Create(action *domain.Action) error {
	args := m.Called(action)
	return args.Error(0)
}

func (m *MockActionRepository) GetByActionID(actionID int) (*domain.Action, error) {
	args := m.Called(actionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Action), args.Error(1)
}

func (m *MockActionRepository) GetAll() ([]*domain.Action, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Action), args.Error(1)
}

func (m *MockActionRepository) Update(action *domain.Action) error {
	args := m.Called(action)
	return args.Error(0)
}

func (m *MockActionRepository) Delete(actionID int) error {
	args := m.Called(actionID)
	return args.Error(0)
}

func TestNewCronService(t *testing.T) {
	mockRepo := new(MockActionRepository)
	areaServiceURL := "http://localhost:8080"
	internalSecret := "test-secret"

	svc := NewCronService(mockRepo, areaServiceURL, internalSecret)

	assert.NotNil(t, svc)
	assert.NotNil(t, svc.cron)
	assert.NotNil(t, svc.jobs)
	assert.Equal(t, areaServiceURL, svc.areaServiceURL)
	assert.Equal(t, internalSecret, svc.internalSecret)
}

func TestCronService_CreateAction_Success(t *testing.T) {
	mockRepo := new(MockActionRepository)
	svc := NewCronService(mockRepo, "http://localhost:8080", "secret")

	action := &domain.Action{
		ActionID: 1,
		Active:   true,
		Type:     "cron",
		Provider: "timer",
		Service:  "timer",
		Title:    "delay_action",
		Input: []domain.InputField{
			{Name: "delay", Value: "10"},
		},
	}

	mockRepo.On("Create", action).Return(nil)

	err := svc.CreateAction(action)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCronService_CreateAction_RepositoryError(t *testing.T) {
	mockRepo := new(MockActionRepository)
	svc := NewCronService(mockRepo, "http://localhost:8080", "secret")

	action := &domain.Action{
		ActionID: 1,
		Title:    "delay_action",
		Input: []domain.InputField{
			{Name: "delay", Value: "10"},
		},
	}

	dbError := errors.New("database error")
	mockRepo.On("Create", action).Return(dbError)

	err := svc.CreateAction(action)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create action")
	mockRepo.AssertExpectations(t)
}

func TestCronService_ActivateAction_Success(t *testing.T) {
	mockRepo := new(MockActionRepository)
	svc := NewCronService(mockRepo, "http://localhost:8080", "secret")

	action := &domain.Action{
		ActionID: 1,
		Active:   false,
		Title:    "delay_action",
		Input: []domain.InputField{
			{Name: "delay", Value: "10"},
		},
	}

	mockRepo.On("GetByActionID", 1).Return(action, nil)
	mockRepo.On("Update", mock.MatchedBy(func(a *domain.Action) bool {
		return a.ActionID == 1 && a.Active == true
	})).Return(nil)

	err := svc.ActivateAction(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCronService_ActivateAction_AlreadyActive(t *testing.T) {
	mockRepo := new(MockActionRepository)
	svc := NewCronService(mockRepo, "http://localhost:8080", "secret")

	action := &domain.Action{
		ActionID: 1,
		Active:   true,
		Title:    "delay_action",
		Input: []domain.InputField{
			{Name: "delay", Value: "10"},
		},
	}

	mockRepo.On("GetByActionID", 1).Return(action, nil)

	err := svc.ActivateAction(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	// Update should not be called since it's already active
	mockRepo.AssertNotCalled(t, "Update")
}

func TestCronService_DeactivateAction_Success(t *testing.T) {
	mockRepo := new(MockActionRepository)
	svc := NewCronService(mockRepo, "http://localhost:8080", "secret")

	action := &domain.Action{
		ActionID: 1,
		Active:   true,
		Title:    "delay_action",
		Input: []domain.InputField{
			{Name: "delay", Value: "10"},
		},
	}

	mockRepo.On("GetByActionID", 1).Return(action, nil)
	mockRepo.On("Update", mock.MatchedBy(func(a *domain.Action) bool {
		return a.ActionID == 1 && a.Active == false
	})).Return(nil)

	err := svc.DeactivateAction(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCronService_DeactivateAction_AlreadyInactive(t *testing.T) {
	mockRepo := new(MockActionRepository)
	svc := NewCronService(mockRepo, "http://localhost:8080", "secret")

	action := &domain.Action{
		ActionID: 1,
		Active:   false,
	}

	mockRepo.On("GetByActionID", 1).Return(action, nil)

	err := svc.DeactivateAction(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestCronService_DeleteAction_Success(t *testing.T) {
	mockRepo := new(MockActionRepository)
	svc := NewCronService(mockRepo, "http://localhost:8080", "secret")

	mockRepo.On("Delete", 1).Return(nil)

	err := svc.DeleteAction(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCronService_DeleteAction_Error(t *testing.T) {
	mockRepo := new(MockActionRepository)
	svc := NewCronService(mockRepo, "http://localhost:8080", "secret")

	dbError := errors.New("database error")
	mockRepo.On("Delete", 1).Return(dbError)

	err := svc.DeleteAction(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete action")
	mockRepo.AssertExpectations(t)
}

func TestCronService_BuildCronExpression_DelayAction(t *testing.T) {
	mockRepo := new(MockActionRepository)
	svc := NewCronService(mockRepo, "http://localhost:8080", "secret")

	testCases := []struct {
		name     string
		delay    string
		expected string
		hasError bool
	}{
		{"10 seconds", "10", "@every 10s", false},
		{"60 seconds", "60", "@every 60s", false},
		{"1 second", "1", "@every 1s", false},
		{"invalid - not a number", "abc", "", true},
		{"invalid - negative", "-5", "", true},
		{"invalid - zero", "0", "", true},
		{"missing delay", "", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			action := &domain.Action{
				Title: "delay_action",
				Input: []domain.InputField{},
			}

			if tc.delay != "" {
				action.Input = append(action.Input, domain.InputField{
					Name:  "delay",
					Value: tc.delay,
				})
			}

			expr, err := svc.buildCronExpression(action)

			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, expr)
			}
		})
	}
}

func TestCronService_BuildCronExpression_DailyAction(t *testing.T) {
	mockRepo := new(MockActionRepository)
	svc := NewCronService(mockRepo, "http://localhost:8080", "secret")

	testCases := []struct {
		name     string
		hour     string
		minute   string
		expected string
		hasError bool
	}{
		{"valid - 9:30 AM", "9", "30", "30 9 * * *", false},
		{"valid - midnight", "0", "0", "0 0 * * *", false},
		{"valid - 23:59", "23", "59", "59 23 * * *", false},
		{"invalid - hour too high", "24", "0", "", true},
		{"invalid - hour negative", "-1", "0", "", true},
		{"invalid - minute too high", "12", "60", "", true},
		{"invalid - minute negative", "12", "-1", "", true},
		{"invalid - non-numeric hour", "abc", "30", "", true},
		{"invalid - non-numeric minute", "12", "xyz", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			action := &domain.Action{
				Title: "daily_action",
				Input: []domain.InputField{
					{Name: "hour", Value: tc.hour},
					{Name: "minute", Value: tc.minute},
				},
			}

			expr, err := svc.buildCronExpression(action)

			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, expr)
			}
		})
	}
}

func TestCronService_BuildCronExpression_WeeklyAction(t *testing.T) {
	mockRepo := new(MockActionRepository)
	svc := NewCronService(mockRepo, "http://localhost:8080", "secret")

	testCases := []struct {
		name      string
		dayOfWeek string
		hour      string
		minute    string
		expected  string
		hasError  bool
	}{
		{"Monday 9:00", "1", "9", "0", "0 9 * * 1", false},
		{"Sunday midnight", "0", "0", "0", "0 0 * * 0", false},
		{"Saturday 23:59", "6", "23", "59", "59 23 * * 6", false},
		{"invalid day - too high", "7", "12", "0", "", true},
		{"invalid day - negative", "-1", "12", "0", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			action := &domain.Action{
				Title: "weekly_action",
				Input: []domain.InputField{
					{Name: "day_of_week", Value: tc.dayOfWeek},
					{Name: "hour", Value: tc.hour},
					{Name: "minute", Value: tc.minute},
				},
			}

			expr, err := svc.buildCronExpression(action)

			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, expr)
			}
		})
	}
}

func TestCronService_BuildCronExpression_MonthlyAction(t *testing.T) {
	mockRepo := new(MockActionRepository)
	svc := NewCronService(mockRepo, "http://localhost:8080", "secret")

	testCases := []struct {
		name       string
		dayOfMonth string
		hour       string
		minute     string
		expected   string
		hasError   bool
	}{
		{"1st day 9:00", "1", "9", "0", "0 9 1 * *", false},
		{"15th day midnight", "15", "0", "0", "0 0 15 * *", false},
		{"31st day 23:59", "31", "23", "59", "59 23 31 * *", false},
		{"invalid day - zero", "0", "12", "0", "", true},
		{"invalid day - too high", "32", "12", "0", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			action := &domain.Action{
				Title: "monthly_action",
				Input: []domain.InputField{
					{Name: "day_of_month", Value: tc.dayOfMonth},
					{Name: "hour", Value: tc.hour},
					{Name: "minute", Value: tc.minute},
				},
			}

			expr, err := svc.buildCronExpression(action)

			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, expr)
			}
		})
	}
}

func TestCronService_BuildCronExpression_UnsupportedAction(t *testing.T) {
	mockRepo := new(MockActionRepository)
	svc := NewCronService(mockRepo, "http://localhost:8080", "secret")

	action := &domain.Action{
		Title: "unsupported_action",
		Input: []domain.InputField{},
	}

	expr, err := svc.buildCronExpression(action)

	assert.Error(t, err)
	assert.Empty(t, expr)
	assert.Contains(t, err.Error(), "unsupported action title")
}

func TestCronService_BuildOutputFields_DelayAction(t *testing.T) {
	mockRepo := new(MockActionRepository)
	svc := NewCronService(mockRepo, "http://localhost:8080", "secret")

	action := &domain.Action{
		Title: "delay_action",
		Input: []domain.InputField{
			{Name: "delay", Value: "60"},
		},
	}

	fields := svc.buildOutputFields(action)

	assert.Len(t, fields, 1)
	assert.Equal(t, "delay", fields[0].Name)
	assert.Equal(t, "60", fields[0].Value)
}

func TestCronService_BuildOutputFields_DailyAction(t *testing.T) {
	mockRepo := new(MockActionRepository)
	svc := NewCronService(mockRepo, "http://localhost:8080", "secret")

	action := &domain.Action{
		Title: "daily_action",
		Input: []domain.InputField{
			{Name: "hour", Value: "9"},
			{Name: "minute", Value: "30"},
		},
	}

	beforeTime := time.Now().Add(-1 * time.Second)
	fields := svc.buildOutputFields(action)
	afterTime := time.Now().Add(1 * time.Second)

	assert.Len(t, fields, 1)
	assert.Equal(t, "triggered_at", fields[0].Name)

	// Verify the timestamp is valid and recent
	triggeredAt, err := time.Parse(time.RFC3339, fields[0].Value)
	assert.NoError(t, err)
	assert.True(t, triggeredAt.After(beforeTime) || triggeredAt.Equal(beforeTime))
	assert.True(t, triggeredAt.Before(afterTime) || triggeredAt.Equal(afterTime))
}

func TestCronService_StartAndStop(t *testing.T) {
	mockRepo := new(MockActionRepository)
	svc := NewCronService(mockRepo, "http://localhost:8080", "secret")

	mockRepo.On("GetAll").Return([]*domain.Action{}, nil)

	// Start the service
	svc.Start()

	// Give it a moment to start
	time.Sleep(10 * time.Millisecond)

	// Stop the service
	svc.Stop()

	mockRepo.AssertExpectations(t)
}
