package http

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/raphael-guer1n/AREA/WebhookService/internal/config"
	"github.com/raphael-guer1n/AREA/WebhookService/internal/service"
	"github.com/raphael-guer1n/AREA/WebhookService/internal/utils"
)

const maxBodyBytes int64 = 1 << 20

type WebhookHandler struct {
	subscriptionSvc *service.SubscriptionService
	providerSvc     *service.ProviderConfigService
}

func NewWebhookHandler(subscriptionSvc *service.SubscriptionService, providerSvc *service.ProviderConfigService) *WebhookHandler {
	return &WebhookHandler{
		subscriptionSvc: subscriptionSvc,
		providerSvc:     providerSvc,
	}
}

func (h *WebhookHandler) HandleReceiveWebhook(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}

	provider, hookID := parseWebhookPath(req.URL.Path)
	if provider == "" || hookID == "" {
		respondJSON(w, http.StatusNotFound, map[string]any{
			"success": false,
			"error":   "webhook not found",
		})
		return
	}

	subscription, err := h.subscriptionSvc.GetSubscriptionByHookID(hookID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to load subscription",
		})
		return
	}
	if subscription == nil || subscription.Provider != provider {
		respondJSON(w, http.StatusNotFound, map[string]any{
			"success": false,
			"error":   "webhook not found",
		})
		return
	}

	providerConfig, err := h.providerSvc.GetProviderConfig(provider)
	if err != nil {
		status := http.StatusInternalServerError
		message := "failed to load provider config"
		if errors.Is(err, service.ErrProviderConfigNotFound) {
			status = http.StatusNotFound
			message = "provider not supported"
		}
		respondJSON(w, status, map[string]any{
			"success": false,
			"error":   message,
		})
		return
	}

	body, err := readBody(req)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var payload any = map[string]any{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &payload); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]any{
				"success": false,
				"error":   "invalid JSON body",
			})
			return
		}
	}

	var subscriptionConfig any = map[string]any{}
	if len(subscription.Config) > 0 {
		if err := json.Unmarshal(subscription.Config, &subscriptionConfig); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]any{
				"success": false,
				"error":   "invalid subscription config",
			})
			return
		}
	}

	if providerConfig.Signature != nil {
		secretValue, ok := utils.ExtractJSONPath(subscriptionConfig, providerConfig.Signature.SecretJSONPath)
		if !ok || fmt.Sprint(secretValue) == "" {
			respondJSON(w, http.StatusUnauthorized, map[string]any{
				"success": false,
				"error":   "missing signature secret",
			})
			return
		}

		if err := validateSignature(providerConfig.Signature, req.Header, fmt.Sprint(secretValue), body); err != nil {
			respondJSON(w, http.StatusUnauthorized, map[string]any{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	}

	eventType := ""
	if providerConfig.EventHeader != "" {
		eventType = req.Header.Get(providerConfig.EventHeader)
	}
	if eventType == "" && providerConfig.EventJSONPath != "" {
		if value, ok := utils.ExtractJSONPath(payload, providerConfig.EventJSONPath); ok {
			eventType = fmt.Sprint(value)
		}
	}

	mapped := map[string]any{}
	if len(providerConfig.Mappings) > 0 {
		mapped, err = buildMappings(payload, providerConfig.Mappings)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]any{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	}

	event := map[string]any{
		"hook_id":     subscription.HookID,
		"user_id":     subscription.UserID,
		"area_id":     subscription.AreaID,
		"provider":    subscription.Provider,
		"event":       eventType,
		"mapped":      mapped,
		"payload":     payload,
		"config":      subscriptionConfig,
		"received_at": time.Now().UTC(),
	}

	// TODO: dispatch event to AreaService using event payload
	_ = event

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"hook_id": subscription.HookID,
			"event":   eventType,
			"mapped":  mapped,
		},
	})
}

func parseWebhookPath(path string) (string, string) {
	trimmed := strings.TrimPrefix(path, "/webhooks/")
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func readBody(req *http.Request) ([]byte, error) {
	defer req.Body.Close()
	body, err := io.ReadAll(io.LimitReader(req.Body, maxBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to read body")
	}
	return body, nil
}

func validateSignature(sig *config.WebhookSignatureConfig, headers http.Header, secret string, body []byte) error {
	received := headers.Get(sig.Header)
	if received == "" {
		return fmt.Errorf("missing signature header")
	}

	expected, err := computeSignature(sig.Type, secret, body)
	if err != nil {
		return err
	}

	if sig.Prefix != "" {
		if !strings.HasPrefix(received, sig.Prefix) {
			return fmt.Errorf("signature prefix mismatch")
		}
		received = strings.TrimPrefix(received, sig.Prefix)
	}

	if !hmac.Equal([]byte(received), []byte(expected)) {
		return fmt.Errorf("signature mismatch")
	}

	return nil
}

func computeSignature(sigType, secret string, body []byte) (string, error) {
	var hash []byte
	switch sigType {
	case "hmac-sha256":
		mac := hmac.New(sha256.New, []byte(secret))
		_, _ = mac.Write(body)
		hash = mac.Sum(nil)
	case "hmac-sha1":
		mac := hmac.New(sha1.New, []byte(secret))
		_, _ = mac.Write(body)
		hash = mac.Sum(nil)
	default:
		return "", fmt.Errorf("unsupported signature type")
	}

	return hex.EncodeToString(hash), nil
}

func buildMappings(payload any, mappings []config.FieldConfig) (map[string]any, error) {
	mapped := make(map[string]any, len(mappings))
	for _, mapping := range mappings {
		value, ok := utils.ExtractJSONPath(payload, mapping.JSONPath)
		if !ok {
			if mapping.Optional {
				continue
			}
			return nil, fmt.Errorf("missing json path %s", mapping.JSONPath)
		}
		coerced, err := coerceValue(value, mapping.Type)
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", mapping.FieldKey, err)
		}
		mapped[mapping.FieldKey] = coerced
	}
	return mapped, nil
}

func coerceValue(value any, valueType string) (any, error) {
	switch valueType {
	case "string":
		if v, ok := value.(string); ok {
			return v, nil
		}
		return nil, fmt.Errorf("expected string")
	case "number":
		switch v := value.(type) {
		case float64:
			return v, nil
		case float32:
			return float64(v), nil
		case int:
			return float64(v), nil
		case int64:
			return float64(v), nil
		default:
			return nil, fmt.Errorf("expected number")
		}
	case "boolean":
		if v, ok := value.(bool); ok {
			return v, nil
		}
		return nil, fmt.Errorf("expected boolean")
	case "json":
		return value, nil
	default:
		return nil, fmt.Errorf("unsupported type")
	}
}
