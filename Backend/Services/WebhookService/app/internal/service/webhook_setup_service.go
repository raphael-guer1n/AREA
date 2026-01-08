package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/raphael-guer1n/AREA/WebhookService/internal/config"
	"github.com/raphael-guer1n/AREA/WebhookService/internal/domain"
	"github.com/raphael-guer1n/AREA/WebhookService/internal/utils"
)

type WebhookSetupService struct {
	oauth2TokenSvc *OAuth2TokenService
	client         *http.Client
}

func NewWebhookSetupService(oauth2TokenSvc *OAuth2TokenService) *WebhookSetupService {
	return &WebhookSetupService{
		oauth2TokenSvc: oauth2TokenSvc,
		client: &http.Client{
			Timeout: 8 * time.Second,
		},
	}
}

func (s *WebhookSetupService) RegisterWebhook(providerConfig *config.WebhookProviderConfig, sub *domain.Subscription, webhookURL string, subscriptionConfig any) (string, error) {
	if providerConfig == nil || providerConfig.Setup == nil {
		return "", nil
	}
	if webhookURL == "" {
		return "", fmt.Errorf("webhook URL is required for provider setup")
	}

	ctx := buildTemplateContext(sub, webhookURL, subscriptionConfig)
	responseBody, err := s.executeAction(providerConfig.Setup, sub, ctx, "setup")
	if err != nil {
		return "", err
	}

	if providerConfig.Setup.ResponseIDJSONPath == "" {
		return "", nil
	}

	var responsePayload any
	if err := json.Unmarshal(responseBody, &responsePayload); err != nil {
		return "", fmt.Errorf("decode provider response: %w", err)
	}

	idValue, ok := utils.ExtractJSONPath(responsePayload, providerConfig.Setup.ResponseIDJSONPath)
	if !ok {
		return "", fmt.Errorf("response id not found at %s", providerConfig.Setup.ResponseIDJSONPath)
	}

	return formatProviderHookID(idValue)
}

func (s *WebhookSetupService) DeleteWebhook(providerConfig *config.WebhookProviderConfig, sub *domain.Subscription, webhookURL string, subscriptionConfig any) error {
	if providerConfig == nil || providerConfig.Teardown == nil {
		return nil
	}

	ctx := buildTemplateContext(sub, webhookURL, subscriptionConfig)
	_, err := s.executeAction(providerConfig.Teardown, sub, ctx, "teardown")
	if err != nil {
		var missingValueErr utils.MissingTemplateValueError
		if errors.As(err, &missingValueErr) && missingValueErr.Key == "provider_hook_id" {
			return ErrProviderHookMissing
		}
		return err
	}

	return nil
}

func (s *WebhookSetupService) executeAction(action *config.WebhookProviderSetupConfig, sub *domain.Subscription, ctx utils.TemplateContext, label string) ([]byte, error) {
	if action == nil {
		return nil, nil
	}

	urlValue, err := utils.RenderTemplateString(action.URLTemplate, ctx)
	if err != nil {
		return nil, err
	}
	urlStr, ok := urlValue.(string)
	if !ok {
		urlStr = fmt.Sprint(urlValue)
	}

	var body io.Reader
	if len(action.BodyTemplate) > 0 {
		var template any
		if err := json.Unmarshal(action.BodyTemplate, &template); err != nil {
			return nil, fmt.Errorf("invalid body template: %w", err)
		}
		rendered, err := utils.RenderTemplateValue(template, ctx)
		if err != nil {
			return nil, err
		}
		payload, err := json.Marshal(rendered)
		if err != nil {
			return nil, fmt.Errorf("marshal webhook payload: %w", err)
		}
		body = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(action.Method, urlStr, body)
	if err != nil {
		return nil, err
	}

	if len(action.BodyTemplate) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	for key, value := range action.Headers {
		renderedValue, err := utils.RenderTemplateString(value, ctx)
		if err != nil {
			return nil, err
		}
		req.Header.Set(key, fmt.Sprint(renderedValue))
	}

	if action.Auth != nil {
		switch action.Auth.Type {
		case "oauth2":
			token, err := s.oauth2TokenSvc.GetProviderToken(sub.UserID, sub.Provider)
			if err != nil {
				return nil, err
			}
			prefix := action.Auth.Prefix
			req.Header.Set(action.Auth.Header, prefix+token)
		default:
			return nil, fmt.Errorf("unsupported auth type %s", action.Auth.Type)
		}
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("provider webhook %s failed: %w", label, err)
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("provider webhook %s failed: status %d: %s", label, resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

func buildTemplateContext(sub *domain.Subscription, webhookURL string, subscriptionConfig any) utils.TemplateContext {
	return utils.TemplateContext{
		HookURL:        webhookURL,
		HookID:         sub.HookID,
		Provider:       sub.Provider,
		ProviderHookID: normalizeProviderHookID(sub.ProviderHookID),
		UserID:         sub.UserID,
		AreaID:         sub.AreaID,
		Config:         subscriptionConfig,
	}
}

func formatProviderHookID(value any) (string, error) {
	switch v := value.(type) {
	case string:
		if v == "" {
			return "", fmt.Errorf("provider hook id is empty")
		}
		return v, nil
	case json.Number:
		if id, err := v.Int64(); err == nil {
			return strconv.FormatInt(id, 10), nil
		}
		if f, err := v.Float64(); err == nil && f == math.Trunc(f) {
			return strconv.FormatInt(int64(f), 10), nil
		}
		return "", fmt.Errorf("provider hook id is invalid")
	case float64:
		if v != math.Trunc(v) {
			return "", fmt.Errorf("provider hook id is not an integer")
		}
		return strconv.FormatInt(int64(v), 10), nil
	case float32:
		f := float64(v)
		if f != math.Trunc(f) {
			return "", fmt.Errorf("provider hook id is not an integer")
		}
		return strconv.FormatInt(int64(f), 10), nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case uint64:
		return strconv.FormatUint(v, 10), nil
	default:
		return "", fmt.Errorf("unsupported provider hook id type %T", value)
	}
}

func normalizeProviderHookID(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return value
	}
	if !strings.ContainsAny(trimmed, "eE") {
		return trimmed
	}
	parsed, err := strconv.ParseFloat(trimmed, 64)
	if err != nil || parsed != math.Trunc(parsed) {
		return trimmed
	}
	return strconv.FormatInt(int64(parsed), 10)
}
