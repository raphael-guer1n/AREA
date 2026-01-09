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

type SubscriptionHandler struct {
	subscriptionSvc *service.SubscriptionService
	cfg             config.Config
}

func NewSubscriptionHandler(subscriptionSvc *service.SubscriptionService, cfg config.Config) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionSvc: subscriptionSvc,
		cfg:             cfg,
	}
}

func (h *SubscriptionHandler) HandleCreateSubscription(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		h.handleCreateSubscription(w, req)
	case http.MethodGet:
		h.handleListSubscriptions(w, req)
	default:
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
	}
}

func (h *SubscriptionHandler) handleCreateSubscription(w http.ResponseWriter, req *http.Request) {

	var body struct {
		UserID   int             `json:"user_id"`
		AreaID   int             `json:"area_id"`
		Provider string          `json:"provider"`
		Config   json.RawMessage `json:"config"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid request body",
		})
		return
	}

	if body.UserID <= 0 || body.AreaID <= 0 || body.Provider == "" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "user_id, area_id and provider are required",
		})
		return
	}

	webhookBaseURL := buildWebhookBaseURL(req, h.cfg.PublicBaseURL)
	subscription, err := h.subscriptionSvc.CreateSubscription(body.UserID, body.AreaID, body.Provider, body.Config, webhookBaseURL)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, service.ErrProviderNotSupported):
			status = http.StatusNotFound
		case errors.Is(err, service.ErrInvalidConfig), errors.Is(err, service.ErrMissingSecret):
			status = http.StatusBadRequest
		}
		respondJSON(w, status, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	webhookURL := buildWebhookURL(webhookBaseURL, subscription.Provider, subscription.HookID)

	respondJSON(w, http.StatusCreated, map[string]any{
		"success": true,
		"data": map[string]any{
			"hook_id":          subscription.HookID,
			"provider_hook_id": subscription.ProviderHookID,
			"provider":         subscription.Provider,
			"webhook_url":      webhookURL,
		},
	})
}

func (h *SubscriptionHandler) HandleSubscription(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		h.handleGetSubscription(w, req)
	case http.MethodDelete:
		h.handleDeleteSubscription(w, req)
	default:
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
	}
}

func (h *SubscriptionHandler) handleListSubscriptions(w http.ResponseWriter, req *http.Request) {
	userIDParam := strings.TrimSpace(req.URL.Query().Get("user_id"))
	if userIDParam == "" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "user_id is required",
		})
		return
	}
	userID, err := strconv.Atoi(userIDParam)
	if err != nil || userID <= 0 {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid user_id",
		})
		return
	}

	subscriptions, err := h.subscriptionSvc.ListSubscriptionsByUserID(userID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to load subscriptions",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"subscriptions": subscriptions,
		},
	})
}

func (h *SubscriptionHandler) handleGetSubscription(w http.ResponseWriter, req *http.Request) {
	hookID := strings.TrimPrefix(req.URL.Path, "/subscriptions/")
	if hookID == "" || hookID == "/" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "hook id is required",
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
	if subscription == nil {
		respondJSON(w, http.StatusNotFound, map[string]any{
			"success": false,
			"error":   "subscription not found",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    subscription,
	})
}

func (h *SubscriptionHandler) handleDeleteSubscription(w http.ResponseWriter, req *http.Request) {
	hookID := strings.TrimPrefix(req.URL.Path, "/subscriptions/")
	if hookID == "" || hookID == "/" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "hook id is required",
		})
		return
	}

	webhookBaseURL := buildWebhookBaseURL(req, h.cfg.PublicBaseURL)
	if err := h.subscriptionSvc.DeleteSubscription(hookID, webhookBaseURL); err != nil {
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
