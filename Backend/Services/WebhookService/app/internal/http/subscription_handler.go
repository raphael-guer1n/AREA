package http

import (
	"encoding/json"
	"net/http"
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
	if req.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}

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
		switch err {
		case service.ErrProviderNotSupported:
			status = http.StatusNotFound
		case service.ErrInvalidConfig, service.ErrMissingSecret:
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
		switch err {
		case service.ErrSubscriptionNotFound:
			status = http.StatusNotFound
		case service.ErrProviderNotSupported:
			status = http.StatusNotFound
		case service.ErrProviderHookMissing:
			status = http.StatusBadRequest
		case service.ErrInvalidConfig:
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
