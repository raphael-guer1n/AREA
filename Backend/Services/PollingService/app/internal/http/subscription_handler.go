package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/raphael-guer1n/AREA/PollingService/internal/service"
)

type ActionHandler struct {
	subscriptionSvc *service.SubscriptionService
	authSvc         *service.AuthService
}

func NewActionHandler(subscriptionSvc *service.SubscriptionService, authSvc *service.AuthService) *ActionHandler {
	return &ActionHandler{
		subscriptionSvc: subscriptionSvc,
		authSvc:         authSvc,
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

	created := make([]map[string]any, 0, len(body.Actions))
	createdActionIDs := make([]int, 0, len(body.Actions))

	for _, action := range body.Actions {
		if action.Type != "polling" {
			respondJSON(w, http.StatusBadRequest, map[string]any{
				"success": false,
				"error":   "only polling actions are supported",
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

		subscription, err := h.subscriptionSvc.CreateSubscription(userID, action.ActionID, action.Provider, action.Service, cfgPayload, action.Active)
		if err != nil {
			for _, actionID := range createdActionIDs {
				_ = h.subscriptionSvc.DeleteSubscription(actionID)
			}
			status := http.StatusInternalServerError
			switch {
			case errors.Is(err, service.ErrProviderNotSupported):
				status = http.StatusNotFound
			case errors.Is(err, service.ErrInvalidConfig):
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

		created = append(created, map[string]any{
			"action_id":        subscription.ActionID,
			"active":           subscription.Active,
			"provider":         subscription.Provider,
			"service":          subscription.Service,
			"interval_seconds": subscription.IntervalSeconds,
			"next_run_at":      subscription.NextRunAt,
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

	updated := make([]map[string]any, 0, len(body.Actions))

	for _, action := range body.Actions {
		if action.Type != "polling" {
			respondJSON(w, http.StatusBadRequest, map[string]any{
				"success": false,
				"error":   "only polling actions are supported",
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

		subscription, err := h.subscriptionSvc.UpdateSubscription(userID, action.ActionID, action.Provider, action.Service, cfgPayload, action.Active)
		if err != nil {
			status := http.StatusInternalServerError
			switch {
			case errors.Is(err, service.ErrProviderNotSupported):
				status = http.StatusNotFound
			case errors.Is(err, service.ErrInvalidConfig):
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

		updated = append(updated, map[string]any{
			"action_id":        subscription.ActionID,
			"active":           subscription.Active,
			"provider":         subscription.Provider,
			"service":          subscription.Service,
			"interval_seconds": subscription.IntervalSeconds,
			"next_run_at":      subscription.NextRunAt,
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

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"action_id":        subscription.ActionID,
			"active":           subscription.Active,
			"provider":         subscription.Provider,
			"service":          subscription.Service,
			"interval_seconds": subscription.IntervalSeconds,
			"next_run_at":      subscription.NextRunAt,
			"last_item_id":     subscription.LastItemID,
			"last_polled_at":   subscription.LastPolledAt,
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

	if err := h.subscriptionSvc.DeleteSubscription(actionID); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
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
