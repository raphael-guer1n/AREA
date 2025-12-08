package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/raphael-guer1n/AREA/AuthService/internal/auth"
	"github.com/raphael-guer1n/AREA/AuthService/internal/oauth2"
	"github.com/raphael-guer1n/AREA/AuthService/internal/service"
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
	r.mux.HandleFunc("/auth/oauth2/providers", r.handleListProviders)
	r.mux.HandleFunc("/auth/oauth2/authorize", r.handleOAuth2Authorize)
	r.mux.HandleFunc("/auth/oauth2/callback", r.handleOAuth2Callback)
	r.mux.HandleFunc("/oauth2/store", r.oauth2Handler.HandleStoreOAuth2)

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

// GET /auth/oauth2/providers - List available OAuth2 providers
func (r *Router) handleListProviders(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}

	if r.oauth2Manager == nil {
		respondJSON(w, http.StatusServiceUnavailable, map[string]any{
			"success": false,
			"error":   "OAuth2 not configured",
		})
		return
	}

	providers := r.oauth2Manager.ListProviders()
	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"providers": providers,
		},
	})
}

// GET /auth/oauth2/authorize?provider=<provider_name>
func (r *Router) handleOAuth2Authorize(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}

	if r.oauth2Manager == nil {
		respondJSON(w, http.StatusServiceUnavailable, map[string]any{
			"success": false,
			"error":   "OAuth2 not configured",
		})
		return
	}

	provider := req.URL.Query().Get("provider")
	if provider == "" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "provider parameter is required",
		})
		return
	}

	authURL, err := r.oauth2Manager.GetAuthURL(provider)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"auth_url": authURL,
			"provider": provider,
		},
	})
}

// GET /auth/oauth2/callback?code=<code>&state=<state>
func (r *Router) handleOAuth2Callback(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}

	if r.oauth2Manager == nil {
		respondJSON(w, http.StatusServiceUnavailable, map[string]any{
			"success": false,
			"error":   "OAuth2 not configured",
		})
		return
	}

	code := req.URL.Query().Get("code")
	state := req.URL.Query().Get("state")

	if code == "" || state == "" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "code and state parameters are required",
		})
		return
	}

	// Handle OAuth2 callback
	userInfo, tokenResp, provider, err := r.oauth2Manager.HandleCallback(state, code)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// TODO: Link OAuth2 account to user or create new user
	// For now, return the OAuth2 user info and tokens
	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"provider":     provider,
			"user_info":    userInfo,
			"access_token": tokenResp.AccessToken,
			"token_type":   tokenResp.TokenType,
			"expires_in":   tokenResp.ExpiresIn,
		},
	})
}
