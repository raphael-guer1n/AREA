package http

import (
	"encoding/json"
	"net/http"
)

type Router struct {
	mux           *http.ServeMux
	oauth2Handler *OAuth2Handler
	authHandler   *AuthHandler
}

func NewRouter(handler *AuthHandler, auth2Handler *OAuth2Handler) *Router {
	r := &Router{
		mux:           http.NewServeMux(),
		oauth2Handler: auth2Handler,
		authHandler:   handler,
	}

	r.routes()
	return r
}

func (r *Router) routes() {
	r.mux.HandleFunc("/health", r.handleHealth)
	r.mux.HandleFunc("/auth/register", r.authHandler.handleRegister)
	r.mux.HandleFunc("/auth/login", r.authHandler.handleLogin)
	r.mux.HandleFunc("/auth/me", r.authHandler.handleMe)

	// OAuth2 routes
	r.mux.HandleFunc("/oauth2/providers", r.oauth2Handler.handleListProviders)
	r.mux.HandleFunc("/oauth2/authorize", r.oauth2Handler.handleOAuth2Authorize)
	r.mux.HandleFunc("/oauth2/callback", r.oauth2Handler.handleOAuth2Callback)
	r.mux.HandleFunc("/oauth2/store", r.oauth2Handler.HandleStoreOAuth2)
	r.mux.HandleFunc("/oauth2/providers/", r.oauth2Handler.handleGetUserServices)
	r.mux.HandleFunc("/oauth2/provider/token/", r.oauth2Handler.handleGetProviderTokenByServiceByUserId)
	r.mux.HandleFunc("/oauth2/provider/profile/", r.oauth2Handler.handleGetProviderProfileByServiceByUserId)
	r.mux.HandleFunc("/oauth2/disconnect", r.oauth2Handler.handleDisconnectProvider)
	r.mux.HandleFunc("/loginwith", r.oauth2Handler.handleLoginWithAuthorize)

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
