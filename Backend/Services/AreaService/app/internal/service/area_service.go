package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings" // Ajout de strings

	"github.com/raphael-guer1n/AREA/AreaService/internal/domain"
)

type AreaService struct {
	areaRepo       domain.AreaRepository
	internalSecret string
}

func NewAreaService(repository domain.AreaRepository, internalSecret string) *AreaService {
	return &AreaService{
		areaRepo:       repository,
		internalSecret: internalSecret,
	}
}

func (s *AreaService) CreateCalendarEvent(authToken string, event domain.Event) (domain.Event, error) {
	payload := map[string]any{
		"summary":     event.Summary,
		"description": event.Description,
		"start": map[string]any{
			"dateTime": event.StartTime,
		},
		"end": map[string]any{
			"dateTime": event.EndTime,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return domain.Event{}, fmt.Errorf("failed to marshal event payload: %w", err)
	}

	url := "https://www.googleapis.com/calendar/v3/calendars/primary/events"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return domain.Event{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return domain.Event{}, fmt.Errorf("failed to call Google Calendar API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return domain.Event{}, fmt.Errorf("google calendar API returned status %d", resp.StatusCode)
	}
	return event, nil
}

func (s *AreaService) LaunchReactions(userToken string, fieldValues map[string]string, reaction domain.ReactionConfig) error {

	apiKey := fieldValues["api_key"]
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" && strings.Contains(reaction.Url, "slack.com/api/") {
		apiKey = os.Getenv("SLACK_BOT_TOKEN")
	}

	envPlaceholderRegexp := regexp.MustCompile(`\{\{\s*env\.([A-Za-z0-9_]+)\s*\}\}`)
	replacePlaceholders := func(input string) string {
		result := input
		for key, value := range fieldValues {
			placeholder := "{{" + key + "}}"
			result = strings.ReplaceAll(result, placeholder, value)
		}
		result = envPlaceholderRegexp.ReplaceAllStringFunc(result, func(match string) string {
			parts := envPlaceholderRegexp.FindStringSubmatch(match)
			if len(parts) != 2 {
				return match
			}
			envValue := os.Getenv(parts[1])
			if envValue == "" {
				return match
			}
			return envValue
		})
		return result
	}

	setNestedValue := func(target map[string]any, path string, value any) {
		if path == "" {
			if nested, ok := value.(map[string]any); ok {
				for k, v := range nested {
					target[k] = v
				}
			}
			return
		}
		if strings.HasPrefix(path, "@") {
			target[path] = value
			return
		}
		parts := strings.Split(path, ".")
		current := target
		for i, part := range parts {
			if i == len(parts)-1 {
				current[part] = value
				return
			}
			next, ok := current[part].(map[string]any)
			if !ok {
				next = make(map[string]any)
				current[part] = next
			}
			current = next
		}
	}

	var buildValue func(field domain.BodyField) (any, error)
	var buildPayload func(fields []domain.BodyField) (map[string]any, error)

	buildPayload = func(fields []domain.BodyField) (map[string]any, error) {
		result := make(map[string]any)
		for _, bf := range fields {
			val, err := buildValue(bf)
			if err != nil {
				return nil, err
			}
			setNestedValue(result, bf.Path, val)
		}
		return result, nil
	}

	buildValue = func(field domain.BodyField) (any, error) {
		switch strings.ToLower(field.Type) {
		case "object":
			var subFields []domain.BodyField
			if err := json.Unmarshal(field.Value, &subFields); err != nil {
				return nil, fmt.Errorf("failed to parse object for path %s: %w", field.Path, err)
			}
			return buildPayload(subFields)

		case "array":
			var rawItems []json.RawMessage
			if err := json.Unmarshal(field.Value, &rawItems); err != nil {
				return nil, fmt.Errorf("failed to parse array for path %s: %w", field.Path, err)
			}

			result := make([]any, 0, len(rawItems))
			for _, rawItem := range rawItems {
				var subField domain.BodyField
				if err := json.Unmarshal(rawItem, &subField); err == nil && subField.Type != "" {
					val, err := buildValue(subField)
					if err != nil {
						return nil, err
					}
					if subField.Path == "" {
						result = append(result, val)
					} else {
						obj := make(map[string]any)
						setNestedValue(obj, subField.Path, val)
						result = append(result, obj)
					}
					continue
				}

				var strVal string
				if err := json.Unmarshal(rawItem, &strVal); err == nil {
					result = append(result, replacePlaceholders(strVal))
					continue
				}

				var generic any
				if err := json.Unmarshal(rawItem, &generic); err == nil {
					result = append(result, generic)
					continue
				}

				return nil, fmt.Errorf("unsupported array item for path %s", field.Path)
			}
			return result, nil

		default:
			valStr := strings.Trim(string(field.Value), `"`)
			finalVal := replacePlaceholders(valStr)

			switch strings.ToLower(field.Type) {
			case "boolean":
				return strings.EqualFold(finalVal, "true") || finalVal == "1", nil
			case "number":
				if numVal, err := strconv.ParseFloat(finalVal, 64); err == nil {
					return numVal, nil
				}
				return finalVal, nil
			default:
				return finalVal, nil
			}
		}
	}

	var bodyReader io.Reader
	contentType := ""
	if strings.EqualFold(reaction.BodyType, "binary") {
		if len(reaction.BodyStruct) == 1 {
			val, err := buildValue(reaction.BodyStruct[0])
			if err != nil {
				return err
			}
			bodyReader = strings.NewReader(fmt.Sprint(val))
		} else if len(reaction.BodyStruct) > 1 {
			payload, err := buildPayload(reaction.BodyStruct)
			if err != nil {
				return err
			}
			bodyReader = strings.NewReader(fmt.Sprint(payload))
		}
	} else {
		payload, err := buildPayload(reaction.BodyStruct)
		if err != nil {
			return err
		}
		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal event payload: %w", err)
		}
		bodyReader = bytes.NewReader(body)
		contentType = "application/json"
	}
	url := reaction.Url

	for key, value := range fieldValues {
		url = strings.ReplaceAll(url, "{{"+key+"}}", value)
	}
	url = replacePlaceholders(url)
	log.Println(url)

	method := reaction.Method
	if method == "" {
		method = http.MethodPost
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if strings.TrimSpace(userToken) != "" {
		req.Header.Set("Authorization", "Bearer "+userToken)
	} else if strings.TrimSpace(apiKey) != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	if clientID := strings.TrimSpace(fieldValues["client_id"]); clientID != "" {
		req.Header.Set("Client-Id", clientID)
	}
	if s.internalSecret != "" {
		req.Header.Set("X-Internal-Secret", s.internalSecret)
	}
	for key, value := range reaction.Headers {
		renderedValue := replacePlaceholders(value)
		req.Header.Set(key, renderedValue)
	}
	if contentType != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call reaction endpoint: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("reaction request returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

func (s *AreaService) GetUserAreas(userId int) ([]domain.Area, error) {
	return s.areaRepo.GetUserAreas(userId)
}

func (s *AreaService) SaveArea(area domain.Area) (domain.Area, error) {
	return s.areaRepo.SaveArea(area)
}

func (s *AreaService) GetAreaFromAction(actionId int) (domain.Area, error) {
	return s.areaRepo.GetAreaFromAction(actionId)
}

func (s *AreaService) GetAreaReactions(areaID int) ([]domain.AreaReaction, error) {
	return s.areaRepo.GetAreaReactions(areaID)
}

func (s *AreaService) GetArea(areaID int) (domain.Area, error) {
	return s.areaRepo.GetArea(areaID)
}

func (s *AreaService) ToggleArea(areaID int, isActive bool) error {
	return s.areaRepo.ToggleArea(areaID, isActive)
}

func (s *AreaService) DeleteArea(areaID int) error {
	return s.areaRepo.DeleteArea(areaID)
}

func (s *AreaService) DeactivateAreasByProvider(userID int, provider string) (int, error) {
	return s.areaRepo.DeactivateAreasByProvider(userID, provider)
}
