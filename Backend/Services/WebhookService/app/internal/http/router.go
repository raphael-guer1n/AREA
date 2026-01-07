package http

import (
	"encoding/json"
	"net/http"
)

type Router struct {
	mux                 *http.ServeMux
	subscriptionHandler *SubscriptionHandler
	webhookHandler      *WebhookHandler
}

func NewRouter(subscriptionHandler *SubscriptionHandler, webhookHandler *WebhookHandler) *Router {
	r := &Router{
		mux:                 http.NewServeMux(),
		subscriptionHandler: subscriptionHandler,
		webhookHandler:      webhookHandler,
	}

	r.routes()
	return r
}

func (r *Router) routes() {
	r.mux.HandleFunc("/health", r.handleHealth)

	// Subscription management (intended for internal use via AreaService)
	r.mux.HandleFunc("/subscriptions", r.subscriptionHandler.HandleCreateSubscription)
	r.mux.HandleFunc("/subscriptions/", r.subscriptionHandler.HandleSubscription)

	// Webhook receivers
	r.mux.HandleFunc("/webhooks/", r.webhookHandler.HandleReceiveWebhook)
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
