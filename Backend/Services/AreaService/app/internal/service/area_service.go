package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/raphael-guer1n/AREA/AreaService/internal/domain"
)

type AreaService struct{}

func NewAreaService() *AreaService {
	return &AreaService{}
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
	log.Println(req)
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
