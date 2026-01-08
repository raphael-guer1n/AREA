package http

import (
	"encoding/json"
	"net/http"
)

type Router struct {
	mux                    *http.ServeMux
	providerHandler        *ProviderHandler
	webhookProviderHandler *WebhookProviderHandler
}

func NewRouter(providerHandler *ProviderHandler, webhookProviderHandler *WebhookProviderHandler) *Router {
	r := &Router{
		mux:                    http.NewServeMux(),
		providerHandler:        providerHandler,
		webhookProviderHandler: webhookProviderHandler,
	}

	r.routes()
	return r
}

func (r *Router) routes() {
	r.mux.HandleFunc("/health", r.handleHealth)

	// Provider configuration endpoints
	r.mux.HandleFunc("/providers/services", r.providerHandler.HandleGetServices)
	r.mux.HandleFunc("/providers/oauth2-config", r.providerHandler.HandleGetOAuth2Config)
	r.mux.HandleFunc("/providers/config", r.providerHandler.HandleGetProviderConfig)

	// Webhook provider configuration endpoints
	r.mux.HandleFunc("/webhooks/providers", r.webhookProviderHandler.HandleGetProviders)
	r.mux.HandleFunc("/webhooks/providers/config", r.webhookProviderHandler.HandleGetProviderConfig)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) handleHealth(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]string{
			"status": "healthy",
		},
	})
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
