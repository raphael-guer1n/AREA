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
	authSvc       *service.AuthService
	oauth2Manager *oauth2.Manager
}

func NewRouter(authSvc *service.AuthService, oauth2Manager *oauth2.Manager) *Router {
	r := &Router{
		mux:           http.NewServeMux(),
		authSvc:       authSvc,
		oauth2Manager: oauth2Manager,
	}

	r.routes()
	return r
}

func (r *Router) routes() {
	r.mux.HandleFunc("/health", r.handleHealth)
	r.mux.HandleFunc("/auth/register", r.handleRegister)
	r.mux.HandleFunc("/auth/login", r.handleLogin)
	r.mux.HandleFunc("/auth/me", r.handleMe)

	// OAuth2 routes
	r.mux.HandleFunc("/auth/oauth2/providers", r.handleListProviders)
	r.mux.HandleFunc("/auth/oauth2/authorize", r.handleOAuth2Authorize)
	r.mux.HandleFunc("/auth/oauth2/callback", r.handleOAuth2Callback)
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

// POST /auth/register
func (r *Router) handleRegister(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}

	var body struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid request body",
		})
		return
	}

	user, token, err := r.authSvc.Register(body.Email, body.Username, body.Password)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, service.ErrInvalidEmail),
			errors.Is(err, service.ErrInvalidUsername),
			errors.Is(err, service.ErrInvalidPassword):
			status = http.StatusBadRequest
		case errors.Is(err, service.ErrEmailAlreadyExists),
			errors.Is(err, service.ErrUsernameExists):
			status = http.StatusConflict
		}

		respondJSON(w, status, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusCreated, map[string]any{
		"success": true,
		"data": map[string]any{
			"user":  user,
			"token": token,
		},
	})
}

// POST /auth/login
func (r *Router) handleLogin(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}

	var body struct {
		EmailOrUsername string `json:"emailOrUsername"`
		Password        string `json:"password"`
	}

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid request body",
		})
		return
	}

	user, token, err := r.authSvc.Login(body.EmailOrUsername, body.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrInvalidCredentials) {
			status = http.StatusUnauthorized
		}

		respondJSON(w, status, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"user":  user,
			"token": token,
		},
	})
}

// GET /auth/me - requires JWT authentication
func (r *Router) handleMe(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}

	// Extract token from Authorization header
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		respondJSON(w, http.StatusUnauthorized, map[string]any{
			"success": false,
			"error":   "missing authorization header",
		})
		return
	}

	// Extract Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		respondJSON(w, http.StatusUnauthorized, map[string]any{
			"success": false,
			"error":   "invalid authorization header format",
		})
		return
	}

	token := parts[1]

	// Validate token and extract user ID
	userID, err := auth.ValidateToken(token)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]any{
			"success": false,
			"error":   "invalid or expired token",
		})
		return
	}

	// Get user profile
	user, err := r.authSvc.GetUserByID(userID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrUserNotFound) {
			status = http.StatusNotFound
		}

		respondJSON(w, status, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"user": user,
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
