package http

import (
	"encoding/json"
	"net/http"
)

type Router struct {
	mux             *http.ServeMux
	providerHandler *ProviderHandler
}

func NewRouter(providerHandler *ProviderHandler) *Router {
	r := &Router{
		mux:             http.NewServeMux(),
		providerHandler: providerHandler,
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
