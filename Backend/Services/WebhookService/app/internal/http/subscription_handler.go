package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/raphael-guer1n/AREA/WebhookService/internal/config"
	"github.com/raphael-guer1n/AREA/WebhookService/internal/service"
)

type ActionHandler struct {
	subscriptionSvc *service.SubscriptionService
	authSvc         *service.AuthService
	cfg             config.Config
}

func NewActionHandler(subscriptionSvc *service.SubscriptionService, authSvc *service.AuthService, cfg config.Config) *ActionHandler {
	return &ActionHandler{
		subscriptionSvc: subscriptionSvc,
		authSvc:         authSvc,
		cfg:             cfg,
	}
}

type actionInput struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type actionRequest struct {
	Active   bool          `json:"active"`
	ActionID int           `json:"action_id"`
	Type     string        `json:"type"`
	Provider string        `json:"provider"`
	Service  string        `json:"service"`
	Title    string        `json:"title"`
	Input    []actionInput `json:"input"`
}

func (h *ActionHandler) HandleActions(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		h.handleCreateActions(w, req)
	case http.MethodPut:
		h.handleUpdateActions(w, req)
	default:
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
	}
}

func (h *ActionHandler) handleCreateActions(w http.ResponseWriter, req *http.Request) {
	userID, err := h.resolveUser(req)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var body struct {
		Actions []actionRequest `json:"actions"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid request body",
		})
		return
	}
	if len(body.Actions) == 0 {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "actions are required",
		})
		return
	}

	webhookBaseURL := buildWebhookBaseURL(req, h.cfg.PublicBaseURL)
	created := make([]map[string]any, 0, len(body.Actions))
	createdActionIDs := make([]int, 0, len(body.Actions))

	for _, action := range body.Actions {
		if action.Type != "webhook" {
			respondJSON(w, http.StatusBadRequest, map[string]any{
				"success": false,
				"error":   "only webhook actions are supported",
			})
			return
		}
		if action.ActionID <= 0 || (action.Provider == "" && action.Service == "") {
			respondJSON(w, http.StatusBadRequest, map[string]any{
				"success": false,
				"error":   "action_id and provider/service are required",
			})
			return
		}
		cfgPayload, err := buildConfigPayload(action.Input)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]any{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		subscription, err := h.subscriptionSvc.CreateSubscription(userID, action.ActionID, action.Provider, action.Service, cfgPayload, action.Active, webhookBaseURL)
		if err != nil {
			for _, actionID := range createdActionIDs {
				_ = h.subscriptionSvc.DeleteSubscription(actionID, webhookBaseURL)
			}
			status := http.StatusInternalServerError
			switch {
			case errors.Is(err, service.ErrProviderNotSupported):
				status = http.StatusNotFound
			case errors.Is(err, service.ErrInvalidConfig), errors.Is(err, service.ErrMissingSecret):
				status = http.StatusBadRequest
			case errors.Is(err, service.ErrActionExists):
				status = http.StatusConflict
			}
			respondJSON(w, status, map[string]any{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		webhookURL := buildWebhookURL(webhookBaseURL, subscription.Service, subscription.HookID)
		created = append(created, map[string]any{
			"action_id":        subscription.ActionID,
			"active":           subscription.Active,
			"hook_id":          subscription.HookID,
			"provider_hook_id": subscription.ProviderHookID,
			"provider":         subscription.Provider,
			"service":          subscription.Service,
			"webhook_url":      webhookURL,
		})
		createdActionIDs = append(createdActionIDs, subscription.ActionID)
	}

	respondJSON(w, http.StatusCreated, map[string]any{
		"success": true,
		"data": map[string]any{
			"actions": created,
		},
	})
}

func (h *ActionHandler) handleUpdateActions(w http.ResponseWriter, req *http.Request) {
	userID, err := h.resolveUser(req)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var body struct {
		Actions []actionRequest `json:"actions"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid request body",
		})
		return
	}
	if len(body.Actions) == 0 {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "actions are required",
		})
		return
	}

	webhookBaseURL := buildWebhookBaseURL(req, h.cfg.PublicBaseURL)
	updated := make([]map[string]any, 0, len(body.Actions))

	for _, action := range body.Actions {
		if action.Type != "webhook" {
			respondJSON(w, http.StatusBadRequest, map[string]any{
				"success": false,
				"error":   "only webhook actions are supported",
			})
			return
		}
		if action.ActionID <= 0 || (action.Provider == "" && action.Service == "") {
			respondJSON(w, http.StatusBadRequest, map[string]any{
				"success": false,
				"error":   "action_id and provider/service are required",
			})
			return
		}
		cfgPayload, err := buildConfigPayload(action.Input)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]any{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		subscription, err := h.subscriptionSvc.UpdateSubscription(userID, action.ActionID, action.Provider, action.Service, cfgPayload, action.Active, webhookBaseURL)
		if err != nil {
			status := http.StatusInternalServerError
			switch {
			case errors.Is(err, service.ErrProviderNotSupported):
				status = http.StatusNotFound
			case errors.Is(err, service.ErrInvalidConfig), errors.Is(err, service.ErrMissingSecret):
				status = http.StatusBadRequest
			case errors.Is(err, service.ErrSubscriptionNotFound):
				status = http.StatusNotFound
			case errors.Is(err, service.ErrUnauthorizedAction):
				status = http.StatusForbidden
			}
			respondJSON(w, status, map[string]any{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		webhookURL := buildWebhookURL(webhookBaseURL, subscription.Service, subscription.HookID)
		updated = append(updated, map[string]any{
			"action_id":        subscription.ActionID,
			"active":           subscription.Active,
			"hook_id":          subscription.HookID,
			"provider_hook_id": subscription.ProviderHookID,
			"provider":         subscription.Provider,
			"service":          subscription.Service,
			"webhook_url":      webhookURL,
		})
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"actions": updated,
		},
	})
}

func (h *ActionHandler) HandleAction(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		h.handleGetAction(w, req)
	case http.MethodDelete:
		h.handleDeleteAction(w, req)
	default:
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
	}
}

func (h *ActionHandler) HandleActivateAction(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}

	userID, err := h.resolveUser(req)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	actionID, err := parseActionID(req.URL.Path, "/activate/")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid action_id",
		})
		return
	}

	webhookBaseURL := buildWebhookBaseURL(req, h.cfg.PublicBaseURL)
	subscription, err := h.subscriptionSvc.ActivateSubscription(userID, actionID, webhookBaseURL)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, service.ErrProviderNotSupported):
			status = http.StatusNotFound
		case errors.Is(err, service.ErrInvalidConfig), errors.Is(err, service.ErrMissingSecret):
			status = http.StatusBadRequest
		case errors.Is(err, service.ErrSubscriptionNotFound):
			status = http.StatusNotFound
		case errors.Is(err, service.ErrUnauthorizedAction):
			status = http.StatusForbidden
		}
		respondJSON(w, status, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	webhookURL := buildWebhookURL(webhookBaseURL, subscription.Service, subscription.HookID)
	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"action_id":        subscription.ActionID,
			"active":           subscription.Active,
			"hook_id":          subscription.HookID,
			"provider_hook_id": subscription.ProviderHookID,
			"provider":         subscription.Provider,
			"service":          subscription.Service,
			"webhook_url":      webhookURL,
		},
	})
}

func (h *ActionHandler) HandleDeactivateAction(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}

	userID, err := h.resolveUser(req)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	actionID, err := parseActionID(req.URL.Path, "/deactivate/")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid action_id",
		})
		return
	}

	webhookBaseURL := buildWebhookBaseURL(req, h.cfg.PublicBaseURL)
	subscription, err := h.subscriptionSvc.DeactivateSubscription(userID, actionID, webhookBaseURL)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, service.ErrProviderNotSupported):
			status = http.StatusNotFound
		case errors.Is(err, service.ErrInvalidConfig), errors.Is(err, service.ErrProviderHookMissing):
			status = http.StatusBadRequest
		case errors.Is(err, service.ErrSubscriptionNotFound):
			status = http.StatusNotFound
		case errors.Is(err, service.ErrUnauthorizedAction):
			status = http.StatusForbidden
		}
		respondJSON(w, status, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	webhookURL := buildWebhookURL(webhookBaseURL, subscription.Service, subscription.HookID)
	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"action_id":        subscription.ActionID,
			"active":           subscription.Active,
			"hook_id":          subscription.HookID,
			"provider_hook_id": subscription.ProviderHookID,
			"provider":         subscription.Provider,
			"service":          subscription.Service,
			"webhook_url":      webhookURL,
		},
	})
}

func (h *ActionHandler) handleGetAction(w http.ResponseWriter, req *http.Request) {
	userID, err := h.resolveUser(req)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	actionID, err := parseActionID(req.URL.Path, "/actions/")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid action_id",
		})
		return
	}

	subscription, err := h.subscriptionSvc.GetSubscriptionByActionID(actionID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to load subscription",
		})
		return
	}
	if subscription == nil {
		respondJSON(w, http.StatusNotFound, map[string]any{
			"success": false,
			"error":   "subscription not found",
		})
		return
	}
	if subscription.UserID != userID {
		respondJSON(w, http.StatusForbidden, map[string]any{
			"success": false,
			"error":   "action does not belong to user",
		})
		return
	}

	webhookBaseURL := buildWebhookBaseURL(req, h.cfg.PublicBaseURL)
	webhookURL := buildWebhookURL(webhookBaseURL, subscription.Service, subscription.HookID)

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"action_id":        subscription.ActionID,
			"active":           subscription.Active,
			"hook_id":          subscription.HookID,
			"provider_hook_id": subscription.ProviderHookID,
			"provider":         subscription.Provider,
			"service":          subscription.Service,
			"webhook_url":      webhookURL,
		},
	})
}

func (h *ActionHandler) handleDeleteAction(w http.ResponseWriter, req *http.Request) {
	userID, err := h.resolveUser(req)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	actionID, err := parseActionID(req.URL.Path, "/actions/")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid action_id",
		})
		return
	}

	subscription, err := h.subscriptionSvc.GetSubscriptionByActionID(actionID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to load subscription",
		})
		return
	}
	if subscription == nil {
		respondJSON(w, http.StatusNotFound, map[string]any{
			"success": false,
			"error":   "subscription not found",
		})
		return
	}
	if subscription.UserID != userID {
		respondJSON(w, http.StatusForbidden, map[string]any{
			"success": false,
			"error":   "action does not belong to user",
		})
		return
	}

	webhookBaseURL := buildWebhookBaseURL(req, h.cfg.PublicBaseURL)
	if err := h.subscriptionSvc.DeleteSubscription(actionID, webhookBaseURL); err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, service.ErrSubscriptionNotFound):
			status = http.StatusNotFound
		case errors.Is(err, service.ErrProviderNotSupported):
			status = http.StatusNotFound
		case errors.Is(err, service.ErrProviderHookMissing):
			status = http.StatusBadRequest
		case errors.Is(err, service.ErrInvalidConfig):
			status = http.StatusBadRequest
		}
		respondJSON(w, status, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
	})
}

func (h *ActionHandler) resolveUser(req *http.Request) (int, error) {
	authHeader := strings.TrimSpace(req.Header.Get("Authorization"))
	if authHeader == "" {
		return 0, errors.New("missing Authorization header")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, errors.New("authorization must be Bearer <token>")
	}
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if token == "" {
		return 0, errors.New("empty token")
	}
	if h.authSvc == nil {
		return 0, errors.New("auth service not configured")
	}
	userID, err := h.authSvc.GetUserID(authHeader)
	if err != nil {
		return 0, err
	}
	if userID <= 0 {
		return 0, errors.New("invalid user")
	}
	return userID, nil
}

func buildConfigPayload(inputs []actionInput) (json.RawMessage, error) {
	cfg := make(map[string]any, len(inputs))
	for _, input := range inputs {
		name := strings.TrimSpace(input.Name)
		if name == "" {
			continue
		}
		cfg[name] = normalizeInputValue(input.Value)
	}
	payload, err := json.Marshal(cfg)
	if err != nil {
		return nil, errors.New("invalid input values")
	}
	return payload, nil
}

func normalizeInputValue(value any) any {
	str, ok := value.(string)
	if !ok {
		return value
	}
	trimmed := strings.TrimSpace(str)
	if trimmed == "" {
		return str
	}
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		var parsed any
		if err := json.Unmarshal([]byte(trimmed), &parsed); err == nil {
			return parsed
		}
	}
	return str
}

func parseActionID(path, prefix string) (int, error) {
	raw := strings.TrimPrefix(path, prefix)
	raw = strings.Trim(raw, "/")
	if raw == "" {
		return 0, errors.New("missing action id")
	}
	return strconv.Atoi(raw)
}

func buildWebhookBaseURL(req *http.Request, publicBaseURL string) string {
	base := publicBaseURL
	if base == "" {
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

		base = scheme + "://" + host
	}

	return strings.TrimRight(base, "/")
}

func buildWebhookURL(baseURL, provider, hookID string) string {
	return strings.TrimRight(baseURL, "/") + "/webhooks/" + provider + "/" + hookID
}
