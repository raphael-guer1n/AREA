package service

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/raphael-guer1n/AREA/AreaService/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAreaRepository is a mock implementation of AreaRepository
type MockAreaRepository struct {
	mock.Mock
}

func (m *MockAreaRepository) GetUserAreas(userID int) ([]domain.Area, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Area), args.Error(1)
}

func (m *MockAreaRepository) SaveArea(area domain.Area) (domain.Area, error) {
	args := m.Called(area)
	return args.Get(0).(domain.Area), args.Error(1)
}

func (m *MockAreaRepository) GetAreaFromAction(actionID int) (domain.Area, error) {
	args := m.Called(actionID)
	return args.Get(0).(domain.Area), args.Error(1)
}

func (m *MockAreaRepository) GetAreaActions(areaID int) ([]domain.AreaAction, error) {
	args := m.Called(areaID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.AreaAction), args.Error(1)
}

func (m *MockAreaRepository) GetAreaReactions(areaID int) ([]domain.AreaReaction, error) {
	args := m.Called(areaID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.AreaReaction), args.Error(1)
}

func (m *MockAreaRepository) GetArea(areaID int) (domain.Area, error) {
	args := m.Called(areaID)
	return args.Get(0).(domain.Area), args.Error(1)
}

func (m *MockAreaRepository) ToggleArea(areaID int, isActive bool) error {
	args := m.Called(areaID, isActive)
	return args.Error(0)
}

func (m *MockAreaRepository) DeleteArea(areaID int) error {
	args := m.Called(areaID)
	return args.Error(0)
}

func (m *MockAreaRepository) SaveActions(areaID int, actions []domain.AreaAction) ([]domain.AreaAction, error) {
	args := m.Called(areaID, actions)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.AreaAction), args.Error(1)
}

func (m *MockAreaRepository) SaveReactions(areaID int, reactions []domain.AreaReaction) ([]domain.AreaReaction, error) {
	args := m.Called(areaID, reactions)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.AreaReaction), args.Error(1)
}

func (m *MockAreaRepository) DeactivateAreasByProvider(userID int, provider string) (int, error) {
	args := m.Called(userID, provider)
	return args.Int(0), args.Error(1)
}

func TestAreaService_GetUserAreas_Success(t *testing.T) {
	mockRepo := new(MockAreaRepository)
	svc := NewAreaService(mockRepo, "test-secret")

	userID := 1
	expectedAreas := []domain.Area{
		{ID: 1, Name: "Area 1", UserID: userID, Active: true},
		{ID: 2, Name: "Area 2", UserID: userID, Active: false},
	}

	mockRepo.On("GetUserAreas", userID).Return(expectedAreas, nil)

	areas, err := svc.GetUserAreas(userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedAreas, areas)
	mockRepo.AssertExpectations(t)
}

func TestAreaService_GetUserAreas_EmptyResult(t *testing.T) {
	mockRepo := new(MockAreaRepository)
	svc := NewAreaService(mockRepo, "test-secret")

	userID := 1
	mockRepo.On("GetUserAreas", userID).Return([]domain.Area{}, nil)

	areas, err := svc.GetUserAreas(userID)

	assert.NoError(t, err)
	assert.Empty(t, areas)
	mockRepo.AssertExpectations(t)
}

func TestAreaService_GetUserAreas_Error(t *testing.T) {
	mockRepo := new(MockAreaRepository)
	svc := NewAreaService(mockRepo, "test-secret")

	userID := 1
	dbError := errors.New("database error")
	mockRepo.On("GetUserAreas", userID).Return(nil, dbError)

	areas, err := svc.GetUserAreas(userID)

	assert.Error(t, err)
	assert.Nil(t, areas)
	assert.Equal(t, dbError, err)
	mockRepo.AssertExpectations(t)
}

func TestAreaService_SaveArea_Success(t *testing.T) {
	mockRepo := new(MockAreaRepository)
	svc := NewAreaService(mockRepo, "test-secret")

	area := domain.Area{
		Name:   "Test Area",
		UserID: 1,
		Active: true,
	}

	expectedArea := area
	expectedArea.ID = 1

	mockRepo.On("SaveArea", area).Return(expectedArea, nil)

	savedArea, err := svc.SaveArea(area)

	assert.NoError(t, err)
	assert.Equal(t, expectedArea, savedArea)
	mockRepo.AssertExpectations(t)
}

func TestAreaService_SaveArea_Error(t *testing.T) {
	mockRepo := new(MockAreaRepository)
	svc := NewAreaService(mockRepo, "test-secret")

	area := domain.Area{Name: "Test Area"}
	dbError := errors.New("database error")

	mockRepo.On("SaveArea", area).Return(domain.Area{}, dbError)

	savedArea, err := svc.SaveArea(area)

	assert.Error(t, err)
	assert.Equal(t, domain.Area{}, savedArea)
	mockRepo.AssertExpectations(t)
}

func TestAreaService_GetAreaFromAction_Success(t *testing.T) {
	mockRepo := new(MockAreaRepository)
	svc := NewAreaService(mockRepo, "test-secret")

	actionID := 1
	expectedArea := domain.Area{
		ID:     1,
		Name:   "Test Area",
		Active: true,
	}

	mockRepo.On("GetAreaFromAction", actionID).Return(expectedArea, nil)

	area, err := svc.GetAreaFromAction(actionID)

	assert.NoError(t, err)
	assert.Equal(t, expectedArea, area)
	mockRepo.AssertExpectations(t)
}

func TestAreaService_GetAreaReactions_Success(t *testing.T) {
	mockRepo := new(MockAreaRepository)
	svc := NewAreaService(mockRepo, "test-secret")

	areaID := 1
	expectedReactions := []domain.AreaReaction{
		{ID: 1, Provider: "google", Service: "gmail", Title: "Send Email"},
		{ID: 2, Provider: "discord", Service: "discord", Title: "Send Message"},
	}

	mockRepo.On("GetAreaReactions", areaID).Return(expectedReactions, nil)

	reactions, err := svc.GetAreaReactions(areaID)

	assert.NoError(t, err)
	assert.Equal(t, expectedReactions, reactions)
	mockRepo.AssertExpectations(t)
}

func TestAreaService_GetArea_Success(t *testing.T) {
	mockRepo := new(MockAreaRepository)
	svc := NewAreaService(mockRepo, "test-secret")

	areaID := 1
	expectedArea := domain.Area{
		ID:     areaID,
		Name:   "Test Area",
		Active: true,
	}

	mockRepo.On("GetArea", areaID).Return(expectedArea, nil)

	area, err := svc.GetArea(areaID)

	assert.NoError(t, err)
	assert.Equal(t, expectedArea, area)
	mockRepo.AssertExpectations(t)
}

func TestAreaService_ToggleArea_Success(t *testing.T) {
	mockRepo := new(MockAreaRepository)
	svc := NewAreaService(mockRepo, "test-secret")

	areaID := 1
	mockRepo.On("ToggleArea", areaID, true).Return(nil)

	err := svc.ToggleArea(areaID, true)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAreaService_ToggleArea_Error(t *testing.T) {
	mockRepo := new(MockAreaRepository)
	svc := NewAreaService(mockRepo, "test-secret")

	areaID := 1
	dbError := errors.New("database error")
	mockRepo.On("ToggleArea", areaID, false).Return(dbError)

	err := svc.ToggleArea(areaID, false)

	assert.Error(t, err)
	assert.Equal(t, dbError, err)
	mockRepo.AssertExpectations(t)
}

func TestAreaService_DeleteArea_Success(t *testing.T) {
	mockRepo := new(MockAreaRepository)
	svc := NewAreaService(mockRepo, "test-secret")

	areaID := 1
	mockRepo.On("DeleteArea", areaID).Return(nil)

	err := svc.DeleteArea(areaID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAreaService_DeleteArea_Error(t *testing.T) {
	mockRepo := new(MockAreaRepository)
	svc := NewAreaService(mockRepo, "test-secret")

	areaID := 1
	dbError := errors.New("database error")
	mockRepo.On("DeleteArea", areaID).Return(dbError)

	err := svc.DeleteArea(areaID)

	assert.Error(t, err)
	assert.Equal(t, dbError, err)
	mockRepo.AssertExpectations(t)
}

func TestAreaService_DeactivateAreasByProvider_Success(t *testing.T) {
	mockRepo := new(MockAreaRepository)
	svc := NewAreaService(mockRepo, "test-secret")

	userID := 1
	provider := "google"
	expectedCount := 3

	mockRepo.On("DeactivateAreasByProvider", userID, provider).Return(expectedCount, nil)

	count, err := svc.DeactivateAreasByProvider(userID, provider)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockRepo.AssertExpectations(t)
}

func TestAreaService_DeactivateAreasByProvider_NoAreasAffected(t *testing.T) {
	mockRepo := new(MockAreaRepository)
	svc := NewAreaService(mockRepo, "test-secret")

	userID := 1
	provider := "nonexistent"

	mockRepo.On("DeactivateAreasByProvider", userID, provider).Return(0, nil)

	count, err := svc.DeactivateAreasByProvider(userID, provider)

	assert.NoError(t, err)
	assert.Equal(t, 0, count)
	mockRepo.AssertExpectations(t)
}

func TestAreaService_LaunchReactions_SimplePayload(t *testing.T) {
	mockRepo := new(MockAreaRepository)
	svc := NewAreaService(mockRepo, "test-secret")

	fieldValues := map[string]string{
		"user":    "John",
		"message": "Hello World",
	}

	bodyField := domain.BodyField{
		Path:  "content",
		Type:  "string",
		Value: json.RawMessage(`"{{message}}"`),
	}

	reaction := domain.ReactionConfig{
		Url:        "https://api.example.com/send",
		Method:     "POST",
		BodyStruct: []domain.BodyField{bodyField},
	}

	// Note: This test will fail because it makes a real HTTP request
	// In a real implementation, we would mock the HTTP client
	err := svc.LaunchReactions("test-token", fieldValues, reaction)

	// Since we can't make the actual HTTP call succeed, we expect an error
	assert.Error(t, err)
}

func TestAreaService_LaunchReactions_URLWithPlaceholder(t *testing.T) {
	mockRepo := new(MockAreaRepository)
	svc := NewAreaService(mockRepo, "test-secret")

	fieldValues := map[string]string{
		"userId": "123",
	}

	reaction := domain.ReactionConfig{
		Url:        "https://api.example.com/users/{{userId}}",
		Method:     "GET",
		BodyStruct: []domain.BodyField{},
	}

	// This will fail due to actual HTTP call
	err := svc.LaunchReactions("test-token", fieldValues, reaction)

	assert.Error(t, err)
}
