package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/raphael-guer1n/AREA/AuthService/internal/config"
	"github.com/raphael-guer1n/AREA/AuthService/internal/domain"
)

type OAuth2StorageService struct {
	profileRepo domain.UserProfileRepository
	fieldRepo   domain.UserServiceFieldRepository
	configSvc   *ProviderConfigService
}

func NewOAuth2StorageService(
	profileRepo domain.UserProfileRepository,
	fieldRepo domain.UserServiceFieldRepository,
	configSvc *ProviderConfigService,
) *OAuth2StorageService {
	return &OAuth2StorageService{
		profileRepo: profileRepo,
		fieldRepo:   fieldRepo,
		configSvc:   configSvc,
	}
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
	// Get provider config
	providerConfig, exists := s.configSvc.GetProviderConfig(serviceName)
	if !exists {
		return fmt.Errorf("unknown service: %s", serviceName)
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
		json.RawMessage(userInfoJSON),
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
			field.JsonValue = json.RawMessage(jsonBytes)

		default:
			return nil, fmt.Errorf("unsupported field type: %s for field %s", mapping.Type, mapping.FieldKey)
		}

		fields = append(fields, field)
	}

	return fields, nil
}
