package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings" // Ajout de strings

	"github.com/raphael-guer1n/AREA/AreaService/internal/domain"
)

type AreaService struct {
	areaRepo domain.AreaRepository
}

func NewAreaService(repository domain.AreaRepository) *AreaService {
	return &AreaService{areaRepo: repository}
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

	var buildPayload func(fields []domain.BodyField) (map[string]any, error)
	buildPayload = func(fields []domain.BodyField) (map[string]any, error) {
		result := make(map[string]any)
		for _, bf := range fields {
			if bf.Type == "object" {
				var subFields []domain.BodyField
				if err := json.Unmarshal(bf.Value, &subFields); err == nil {
					result[bf.Path], err = buildPayload(subFields)
				}
			} else {
				valStr := string(bf.Value)
				valStr = string(bytes.Trim(bf.Value, `"`))

				finalVal := valStr
				for key, value := range fieldValues {
					placeholder := "{{" + key + "}}"
					finalVal = strings.ReplaceAll(finalVal, placeholder, value)
				}

				result[bf.Path] = finalVal
			}
		}
		return result, nil
	}

	payload, err := buildPayload(reaction.BodyStruct)

	if err != nil {
		return err
	}
	url := reaction.Url

	for key, value := range fieldValues {
		url = strings.ReplaceAll(url, "{{"+key+"}}", value)
	}
	log.Println(url)
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+userToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	fmt.Println(payload)
	if err != nil {
		return fmt.Errorf("failed to call Google Calendar API: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println(resp)
		return fmt.Errorf("google calendar API returned status %d", resp.StatusCode)
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
