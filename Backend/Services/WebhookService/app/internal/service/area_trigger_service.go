package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type TriggerOutputField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type AreaTriggerService struct {
	baseURL string
	client  *http.Client
}

func NewAreaTriggerService(baseURL string) *AreaTriggerService {
	return &AreaTriggerService{
		baseURL: strings.TrimRight(baseURL, "/"),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *AreaTriggerService) Trigger(authToken string, actionID int, outputFields []TriggerOutputField) error {
	if actionID <= 0 {
		return errors.New("action_id is required")
	}
	endpoint := s.baseURL + "/triggerArea"
	payload := map[string]any{
		"action_id":     actionID,
		"output_fields": outputFields,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if strings.TrimSpace(authToken) != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("trigger area: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		var rawErr json.RawMessage
		if len(bodyBytes) > 0 {
			rawErr = bodyBytes
		}
		message, _ := parseRemoteError(rawErr)
		if message == "" {
			message = "failed to trigger area"
		}
		return errors.New(message)
	}

	return nil
}
