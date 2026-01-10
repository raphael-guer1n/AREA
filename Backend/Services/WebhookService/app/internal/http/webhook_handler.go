package http

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"log"
	"net/http"
	"strconv"
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
	areaTriggerSvc  *service.AreaTriggerService
}

func NewWebhookHandler(subscriptionSvc *service.SubscriptionService, providerSvc *service.ProviderConfigService, areaTriggerSvc *service.AreaTriggerService) *WebhookHandler {
	return &WebhookHandler{
		subscriptionSvc: subscriptionSvc,
		providerSvc:     providerSvc,
		areaTriggerSvc:  areaTriggerSvc,
	}
}

func (h *WebhookHandler) HandleReceiveWebhook(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost && req.Method != http.MethodGet {
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

	if req.Method == http.MethodGet {
		challenge := req.URL.Query().Get("hub.challenge")
		if challenge == "" {
			respondJSON(w, http.StatusBadRequest, map[string]any{
				"success": false,
				"error":   "missing hub.challenge",
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
		if subscription == nil || subscription.Service != provider {
			mode := strings.ToLower(req.URL.Query().Get("hub.mode"))
			if mode != "unsubscribe" {
				respondJSON(w, http.StatusNotFound, map[string]any{
					"success": false,
					"error":   "webhook not found",
				})
				return
			}
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(challenge))
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
	if subscription == nil || subscription.Service != provider {
		respondJSON(w, http.StatusNotFound, map[string]any{
			"success": false,
			"error":   "webhook not found",
		})
		return
	}
	if !subscription.Active {
		_, _ = io.Copy(io.Discard, req.Body)
		log.Printf(
			"webhook ignored: hook_id=%s action_id=%d provider=%s inactive=true",
			subscription.HookID,
			subscription.ActionID,
			subscription.Service,
		)
		respondJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data": map[string]any{
				"hook_id":  subscription.HookID,
				"inactive": true,
			},
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

	payload, err := parsePayload(body, providerConfig.PayloadFormat)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
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
		secret := strings.TrimSpace(fmt.Sprint(secretValue))
		if !ok || secret == "" {
			respondJSON(w, http.StatusUnauthorized, map[string]any{
				"success": false,
				"error":   "missing signature secret",
			})
			return
		}

		if err := validateSignature(providerConfig.Signature, req, secret, body); err != nil {
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

	if ok, reason := shouldProcessEvent(eventType, providerConfig, subscriptionConfig); !ok {
		log.Printf(
			"webhook ignored: hook_id=%s action_id=%d provider=%s event=%s reason=%s",
			subscription.HookID,
			subscription.ActionID,
			subscription.Service,
			eventType,
			reason,
		)
		respondJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data": map[string]any{
				"hook_id": subscription.HookID,
				"event":   eventType,
				"ignored": true,
				"reason":  reason,
			},
		})
		return
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

	outputFields := buildOutputFields(providerConfig.Mappings, mapped)
	if h.areaTriggerSvc != nil {
		if err := h.areaTriggerSvc.Trigger(subscription.AuthToken, subscription.ActionID, outputFields); err != nil {
			log.Printf(
				"webhook dispatch failed: hook_id=%s action_id=%d provider=%s error=%v",
				subscription.HookID,
				subscription.ActionID,
				subscription.Service,
				err,
			)
		}
	}

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

func shouldProcessEvent(eventType string, providerConfig *config.WebhookProviderConfig, subscriptionConfig any) (bool, string) {
	if providerConfig == nil {
		return true, ""
	}

	event := strings.ToLower(strings.TrimSpace(eventType))
	if event == "" {
		return true, ""
	}

	if len(providerConfig.EventIgnore) > 0 && eventInList(event, providerConfig.EventIgnore) {
		return false, "ignored"
	}

	if providerConfig.EventAllowServiceConfigFilePath == "" || subscriptionConfig == nil {
		return true, ""
	}

	path := strings.TrimSpace(providerConfig.EventAllowServiceConfigFilePath)
	path = strings.TrimPrefix(path, "config.")
	if path == "" {
		return true, ""
	}

	value, ok := utils.ExtractJSONPath(subscriptionConfig, path)
	if !ok {
		return true, ""
	}

	allowed, ok := normalizeEventList(value)
	if !ok {
		return true, ""
	}
	if len(allowed) == 0 {
		return false, "filtered"
	}
	if eventInList("*", allowed) {
		return true, ""
	}
	if !eventInList(event, allowed) {
		return false, "filtered"
	}

	return true, ""
}

func normalizeEventList(value any) ([]string, bool) {
	switch v := value.(type) {
	case []string:
		out := make([]string, 0, len(v))
		for _, item := range v {
			out = append(out, strings.TrimSpace(item))
		}
		return out, true
	case []any:
		out := make([]string, 0, len(v))
		for _, item := range v {
			out = append(out, strings.TrimSpace(fmt.Sprint(item)))
		}
		return out, true
	case string:
		return []string{strings.TrimSpace(v)}, true
	default:
		return nil, false
	}
}

func eventInList(event string, list []string) bool {
	needle := strings.ToLower(strings.TrimSpace(event))
	for _, item := range list {
		if strings.ToLower(strings.TrimSpace(item)) == needle {
			return true
		}
	}
	return false
}

func readBody(req *http.Request) ([]byte, error) {
	defer req.Body.Close()
	body, err := io.ReadAll(io.LimitReader(req.Body, maxBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to read body")
	}
	return body, nil
}

func validateSignature(sig *config.WebhookSignatureConfig, req *http.Request, secret string, body []byte) error {
	signatureType, algorithm := normalizeSignatureType(sig)
	if signatureType == "token" {
		signatureType = "header"
	}

	switch signatureType {
	case "hmac":
		return validateHMACSignature(sig, algorithm, req, secret, body)
	case "header":
		return validateHeaderSignature(sig, req.Header, secret)
	default:
		return fmt.Errorf("unsupported signature type")
	}
}

func validateHeaderSignature(sig *config.WebhookSignatureConfig, headers http.Header, secret string) error {
	received := headers.Get(sig.Header)
	if received == "" {
		return fmt.Errorf("missing signature header")
	}

	if sig.Prefix != "" {
		if !strings.HasPrefix(received, sig.Prefix) {
			return fmt.Errorf("signature prefix mismatch")
		}
		received = strings.TrimPrefix(received, sig.Prefix)
	}

	if !hmac.Equal([]byte(received), []byte(secret)) {
		return fmt.Errorf("signature mismatch")
	}

	return nil
}

func validateHMACSignature(sig *config.WebhookSignatureConfig, algorithm string, req *http.Request, secret string, body []byte) error {
	received := req.Header.Get(sig.Header)
	if received == "" {
		return fmt.Errorf("missing signature header")
	}

	if sig.Prefix != "" {
		if !strings.HasPrefix(received, sig.Prefix) {
			return fmt.Errorf("signature prefix mismatch")
		}
		received = strings.TrimPrefix(received, sig.Prefix)
	}

	if sig.TimestampHeader != "" {
		if err := validateTimestamp(sig, req.Header); err != nil {
			return err
		}
	}

	signingInput := body
	if sig.SigningStringTemplate != "" {
		ctx := utils.TemplateContext{
			Body:    string(body),
			Headers: req.Header,
			Method:  req.Method,
			Path:    req.URL.Path,
			URL:     buildRequestURL(req),
			Query:   req.URL.RawQuery,
		}
		rendered, err := utils.RenderTemplateString(sig.SigningStringTemplate, ctx)
		if err != nil {
			return err
		}
		signingInput = []byte(fmt.Sprint(rendered))
	}

	expected, err := computeHMACSignature(algorithm, sig.Encoding, secret, signingInput)
	if err != nil {
		return err
	}

	if !hmac.Equal([]byte(received), []byte(expected)) {
		return fmt.Errorf("signature mismatch")
	}

	return nil
}

func validateTimestamp(sig *config.WebhookSignatureConfig, headers http.Header) error {
	value := headers.Get(sig.TimestampHeader)
	if value == "" {
		return fmt.Errorf("missing timestamp header")
	}

	timestamp, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid timestamp header")
	}

	tolerance := sig.TimestampToleranceSeconds
	if tolerance <= 0 {
		tolerance = 300
	}

	now := time.Now().Unix()
	diff := now - timestamp
	if diff < 0 {
		diff = -diff
	}

	if diff > int64(tolerance) {
		return fmt.Errorf("timestamp outside allowed window")
	}

	return nil
}

func computeHMACSignature(algorithm, encoding, secret string, payload []byte) (string, error) {
	if algorithm == "" {
		algorithm = "sha256"
	}

	var mac hashFunc
	switch strings.ToLower(algorithm) {
	case "sha1":
		mac = sha1.New
	case "sha256":
		mac = sha256.New
	case "sha512":
		mac = sha512.New
	default:
		return "", fmt.Errorf("unsupported signature algorithm")
	}

	h := hmac.New(mac, []byte(secret))
	_, _ = h.Write(payload)
	sum := h.Sum(nil)

	if encoding == "" {
		encoding = "hex"
	}

	switch strings.ToLower(encoding) {
	case "hex":
		return hex.EncodeToString(sum), nil
	case "base64":
		return base64.StdEncoding.EncodeToString(sum), nil
	default:
		return "", fmt.Errorf("unsupported signature encoding")
	}
}

type hashFunc func() hash.Hash

func normalizeSignatureType(sig *config.WebhookSignatureConfig) (string, string) {
	signatureType := strings.ToLower(sig.Type)
	algorithm := strings.ToLower(sig.Algorithm)

	if strings.HasPrefix(signatureType, "hmac-") {
		algorithm = strings.TrimPrefix(signatureType, "hmac-")
		signatureType = "hmac"
	}

	return signatureType, algorithm
}

func buildRequestURL(req *http.Request) string {
	scheme := "http"
	if req.TLS != nil {
		scheme = "https"
	}
	if forwardedProto := req.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
		scheme = forwardedProto
	}

	host := req.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = req.Host
	}

	if host == "" {
		return req.URL.RequestURI()
	}

	return scheme + "://" + host + req.URL.RequestURI()
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

func buildOutputFields(mappings []config.FieldConfig, mapped map[string]any) []service.TriggerOutputField {
	if len(mapped) == 0 {
		return []service.TriggerOutputField{}
	}
	fields := make([]service.TriggerOutputField, 0, len(mapped))
	for _, mapping := range mappings {
		value, ok := mapped[mapping.FieldKey]
		if !ok {
			continue
		}
		fields = append(fields, service.TriggerOutputField{
			Name:  mapping.FieldKey,
			Value: stringifyOutputValue(value),
		})
	}
	return fields
}

func stringifyOutputValue(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case nil:
		return ""
	default:
		switch v.(type) {
		case map[string]any, []any:
			encoded, err := json.Marshal(v)
			if err == nil {
				return string(encoded)
			}
		}
		return fmt.Sprint(v)
	}
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

func parsePayload(body []byte, format string) (any, error) {
	if len(body) == 0 {
		return map[string]any{}, nil
	}

	switch strings.ToLower(strings.TrimSpace(format)) {
	case "", "json":
		var payload any
		if err := json.Unmarshal(body, &payload); err != nil {
			return nil, fmt.Errorf("invalid JSON body")
		}
		return payload, nil
	case "xml":
		payload, err := utils.ParseXMLToMap(body)
		if err != nil {
			return nil, fmt.Errorf("invalid XML body")
		}
		return payload, nil
	default:
		return nil, fmt.Errorf("unsupported payload format")
	}
}
