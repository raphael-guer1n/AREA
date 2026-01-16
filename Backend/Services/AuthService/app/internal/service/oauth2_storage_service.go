package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/raphael-guer1n/AREA/AuthService/internal/config"
	"github.com/raphael-guer1n/AREA/AuthService/internal/domain"
)

type OAuth2StorageService struct {
	profileRepo       domain.UserProfileRepository
	fieldRepo         domain.UserServiceFieldRepository
	configCache       map[string]*config.ProviderConfig
	configMutex       sync.RWMutex
	serviceServiceURL string
	internalSecret    string
	httpClient        *http.Client
}

func NewOAuth2StorageService(
	profileRepo domain.UserProfileRepository,
	fieldRepo domain.UserServiceFieldRepository,
	serviceServiceURL string,
	internalSecret string,
) *OAuth2StorageService {
	return &OAuth2StorageService{
		profileRepo:       profileRepo,
		fieldRepo:         fieldRepo,
		configCache:       make(map[string]*config.ProviderConfig),
		serviceServiceURL: serviceServiceURL,
		internalSecret:    internalSecret,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// getProviderConfig retrieves provider config from cache or fetches from ServiceService API
func (s *OAuth2StorageService) getProviderConfig(serviceName string) (*config.ProviderConfig, error) {
	// Check cache first (read lock)
	s.configMutex.RLock()
	if cached, exists := s.configCache[serviceName]; exists {
		s.configMutex.RUnlock()
		return cached, nil
	}
	s.configMutex.RUnlock()

	// Not in cache, fetch from API (write lock)
	s.configMutex.Lock()
	defer s.configMutex.Unlock()

	// Double-check in case another goroutine loaded it
	if cached, exists := s.configCache[serviceName]; exists {
		return cached, nil
	}

	// Fetch from ServiceService API
	url := fmt.Sprintf("%s/providers/config?service=%s", s.serviceServiceURL, serviceName)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider config request: %w", err)
	}
	if strings.TrimSpace(s.internalSecret) != "" {
		req.Header.Set("X-Internal-Secret", s.internalSecret)
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch provider config from ServiceService: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ServiceService returned status %d for service %s", resp.StatusCode, serviceName)
	}

	// Parse response
	var apiResp struct {
		Success bool                  `json:"success"`
		Data    config.ProviderConfig `json:"data"`
		Error   string                `json:"error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode provider config response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("ServiceService error: %s", apiResp.Error)
	}

	// Cache the config
	s.configCache[serviceName] = &apiResp.Data

	return &apiResp.Data, nil
}

// StoreOAuth2Response stores the OAuth2 user info response in the database
// It creates a user_service_profile entry and extracts fields based on the provider's mapping configuration
func (s *OAuth2StorageService) StoreOAuth2Response(
	userId int,
	serviceName string,
	accessToken string,
	refreshToken string,
	expiresAt time.Time,
	userInfoJSON []byte,
) error {
	// Get provider config (lazy-loaded from ServiceService API)
	providerConfig, err := s.getProviderConfig(serviceName)
	if err != nil {
		return fmt.Errorf("failed to get provider config: %w", err)
	}

	// Parse user info JSON
	var userInfo map[string]interface{}
	if err := json.Unmarshal(userInfoJSON, &userInfo); err != nil {
		return fmt.Errorf("failed to parse user info JSON: %w", err)
	}

	// Extract provider_user_id from the mappings
	providerUserId, err := s.extractProviderUserId(userInfo, providerConfig.Mappings)
	if err != nil {
		return fmt.Errorf("failed to extract provider_user_id: %w", err)
	}

	// Create user profile
	profile, err := s.profileRepo.Create(
		userId,
		serviceName,
		providerUserId,
		accessToken,
		refreshToken,
		expiresAt,
		userInfoJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to create user profile: %w", err)
	}

	// Extract and store fields based on mappings
	fields, err := s.extractFields(profile.ID, userInfo, providerConfig.Mappings)
	if err != nil {
		return fmt.Errorf("failed to extract fields: %w", err)
	}

	if err := s.fieldRepo.CreateBatch(fields); err != nil {
		return fmt.Errorf("failed to create fields: %w", err)
	}

	return nil
}

// extractProviderUserId extracts the provider_user_id from the user info
func (s *OAuth2StorageService) extractProviderUserId(userInfo map[string]interface{}, mappings []config.FieldConfig) (string, error) {
	for _, mapping := range mappings {
		if mapping.FieldKey == "provider_user_id" {
			value, exists := userInfo[mapping.JSONPath]
			if !exists {
				return "", fmt.Errorf("provider_user_id field '%s' not found in user info", mapping.JSONPath)
			}
			return fmt.Sprintf("%v", value), nil
		}
	}
	return "", fmt.Errorf("provider_user_id mapping not found in provider config")
}

// extractFields extracts all fields from user info based on the provider's field mappings
func (s *OAuth2StorageService) extractFields(profileId int, userInfo map[string]interface{}, mappings []config.FieldConfig) ([]domain.UserServiceField, error) {
	fields := make([]domain.UserServiceField, 0, len(mappings))

	for _, mapping := range mappings {
		value, exists := userInfo[mapping.JSONPath]
		if !exists {
			// Skip fields that don't exist in the response
			continue
		}

		field := domain.UserServiceField{
			ProfileId: profileId,
			FieldKey:  mapping.FieldKey,
		}

		// Set value based on type
		switch mapping.Type {
		case "string":
			if strVal, ok := value.(string); ok {
				field.StringValue = strVal
			} else {
				field.StringValue = fmt.Sprintf("%v", value)
			}

		case "number":
			switch v := value.(type) {
			case float64:
				field.NumberValue = v
			case int:
				field.NumberValue = float64(v)
			case string:
				if numVal, err := strconv.ParseFloat(v, 64); err == nil {
					field.NumberValue = numVal
				}
			}

		case "boolean":
			if boolVal, ok := value.(bool); ok {
				field.BoolValue = boolVal
			} else if strVal, ok := value.(string); ok {
				field.BoolValue = strVal == "true" || strVal == "1"
			}

		case "json":
			jsonBytes, err := json.Marshal(value)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal JSON field %s: %w", mapping.FieldKey, err)
			}
			raw := json.RawMessage(jsonBytes)
			field.JsonValue = &raw

		default:
			return nil, fmt.Errorf("unsupported field type: %s for field %s", mapping.Type, mapping.FieldKey)
		}

		fields = append(fields, field)
	}

	return fields, nil
}

// GetUserServicesStatus retrieves all available providers and their login status for a user
func (s *OAuth2StorageService) GetUserServicesStatus(userId int) ([]map[string]interface{}, error) {
	// Fetch all available providers from ServiceService API
	url := fmt.Sprintf("%s/providers/services", s.serviceServiceURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create providers request: %w", err)
	}
	if strings.TrimSpace(s.internalSecret) != "" {
		req.Header.Set("X-Internal-Secret", s.internalSecret)
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch providers from ServiceService: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ServiceService returned status %d", resp.StatusCode)
	}

	var apiResp struct {
		Success bool `json:"success"`
		Data    struct {
			Services []string `json:"services"`
		} `json:"data"`
		Error string `json:"error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode providers response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("ServiceService error: %s", apiResp.Error)
	}

	// Get user's logged services with reconnect status
	loggedServices, err := s.profileRepo.GetServicesStatusByUserId(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user services: %w", err)
	}

	// Create a map for a quick lookup
	loggedServicesMap := make(map[string]bool)
	reconnectMap := make(map[string]bool)
	for _, service := range loggedServices {
		loggedServicesMap[service.Service] = true
		reconnectMap[service.Service] = service.NeedsReconnect
	}

	// Build response
	result := make([]map[string]interface{}, 0, len(apiResp.Data.Services))
	for _, serviceName := range apiResp.Data.Services {
		needsReconnect := reconnectMap[serviceName]
		isLogged := loggedServicesMap[serviceName] && !needsReconnect
		logoURL := ""

		if providerCfg, err := s.getProviderConfig(serviceName); err == nil && providerCfg != nil {
			logoURL = providerCfg.LogoURL
		}

		result = append(result, map[string]interface{}{
			"provider":          serviceName,
			"is_logged":         isLogged,
			"need_reconnecting": needsReconnect,
			"logo_url":          logoURL,
		})
	}

	return result, nil
}

func (s *OAuth2StorageService) GetProviderTokenByServiceByUser(userId int, serviceName string) (string, error) {
	return s.profileRepo.GetProviderUserTokenByServiceByUserId(userId, serviceName)
}

func (s *OAuth2StorageService) GetProviderProfileByServiceByUser(userId int, serviceName string) (domain.UserProfile, error) {
	return s.profileRepo.GetProviderProfileProfileByServiceByUser(userId, serviceName)
}

func (s *OAuth2StorageService) GetProviderFieldsByProfileId(profileId int) ([]domain.UserServiceField, error) {
	return s.fieldRepo.GetFieldsByProfileId(profileId)
}

func (s *OAuth2StorageService) DeleteProviderConnection(userId int, serviceName string) error {
	return s.profileRepo.DeleteByUserIdAndService(userId, serviceName)
}
