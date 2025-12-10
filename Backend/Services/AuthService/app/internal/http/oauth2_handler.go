package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/raphael-guer1n/AREA/AuthService/internal/auth"
	"github.com/raphael-guer1n/AREA/AuthService/internal/oauth2"
	"github.com/raphael-guer1n/AREA/AuthService/internal/service"
)

type OAuth2Handler struct {
	oauth2StorageSvc *service.OAuth2StorageService
	oauth2Manager    *oauth2.Manager
	authSvc          *service.AuthService
}

func NewOAuth2Handler(
	oauth2StorageSvc *service.OAuth2StorageService,
	oauth2Manager *oauth2.Manager,
	authSvc *service.AuthService,
) *OAuth2Handler {
	return &OAuth2Handler{
		oauth2StorageSvc: oauth2StorageSvc,
		oauth2Manager:    oauth2Manager,
		authSvc:          authSvc,
	}
}

type StoreOAuth2Request struct {
	UserId       int    `json:"user_id"`
	Service      string `json:"service"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
	UserInfo     any    `json:"user_info"`
}

func (h *OAuth2Handler) HandleStoreOAuth2(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "Method not allowed",
		})
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "Failed to read request body",
		})
		return
	}
	defer req.Body.Close()

	var storeReq StoreOAuth2Request
	if err := json.Unmarshal(body, &storeReq); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "Invalid JSON format",
		})
		return
	}

	// Validate required fields
	if storeReq.UserId == 0 {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "user_id is required",
		})
		return
	}
	if storeReq.Service == "" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "service is required",
		})
		return
	}
	if storeReq.AccessToken == "" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "access_token is required",
		})
		return
	}

	// Convert user_info back to JSON bytes
	userInfoJSON, err := json.Marshal(storeReq.UserInfo)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "Invalid user_info format",
		})
		return
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(time.Duration(storeReq.ExpiresIn) * time.Second)

	// Store OAuth2 data
	err = h.oauth2StorageSvc.StoreOAuth2Response(
		storeReq.UserId,
		storeReq.Service,
		storeReq.AccessToken,
		storeReq.RefreshToken,
		expiresAt,
		userInfoJSON,
	)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusCreated, map[string]any{
		"success": true,
		"message": "OAuth2 data stored successfully",
	})
}

// GET /oauth2/providers - List available OAuth2 providers
func (h *OAuth2Handler) handleListProviders(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}

	if h.oauth2Manager == nil {
		respondJSON(w, http.StatusServiceUnavailable, map[string]any{
			"success": false,
			"error":   "OAuth2 not configured",
		})
		return
	}

	providers, err := h.oauth2Manager.ListProviders()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"providers": providers,
		},
	})
}

// GET /auth/oauth2/login?provider=<provider_name>&callback_url=<url>&platform=<platform>
func (h *OAuth2Handler) handleOAuth2Login(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}

	if h.oauth2Manager == nil {
		respondJSON(w, http.StatusServiceUnavailable, map[string]any{
			"success": false,
			"error":   "OAuth2 not configured",
		})
		return
	}

	provider := req.URL.Query().Get("provider")
	callbackURL := req.URL.Query().Get("callback_url")
	platform := req.URL.Query().Get("platform")

	if provider == "" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "provider parameter is required",
		})
		return
	}

	if platform == "" {
		platform = "web"
	}

	authURL, err := h.oauth2Manager.GetAuthURL(provider, 0, false, callbackURL, platform, "login")
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

// GET /auth/oauth2/authorize?provider=<provider_name>&callback_url=<url>&platform=<platform>
func (h *OAuth2Handler) handleOAuth2Authorize(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}

	if h.oauth2Manager == nil {
		respondJSON(w, http.StatusServiceUnavailable, map[string]any{
			"success": false,
			"error":   "OAuth2 not configured",
		})
		return
	}

	provider := req.URL.Query().Get("provider")
	callbackURL := req.URL.Query().Get("callback_url")
	platform := req.URL.Query().Get("platform")

	if provider == "" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "provider parameter is required",
		})
		return
	}

	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		respondJSON(w, http.StatusUnauthorized, map[string]any{
			"success": false,
			"error":   "missing authorization header",
		})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		respondJSON(w, http.StatusUnauthorized, map[string]any{
			"success": false,
			"error":   "invalid authorization header format",
		})
		return
	}

	token := parts[1]
	userID, err := auth.ValidateToken(token)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]any{
			"success": false,
			"error":   "invalid or expired token",
		})
		return
	}

	// Default platform to "web" if not specified
	if platform == "" {
		platform = "web"
	}

	authURL, err := h.oauth2Manager.GetAuthURL(provider, userID, true, callbackURL, platform, "link")
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

// GET /oauth2/callback?code=<code>&state=<state>
func (h *OAuth2Handler) handleOAuth2Callback(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}

	if h.oauth2Manager == nil {
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

	// Handle OAuth2 callback - returns StateData with user_id
	userInfo, tokenResp, stateData, err := h.oauth2Manager.HandleCallback(state, code)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Marshal user info to JSON
	userInfoJSON, err := json.Marshal(userInfo.RawData)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to marshal user info",
		})
		return
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	userID := stateData.UserID
	purpose := stateData.Purpose
	if purpose == "" {
		purpose = "link"
	}

	var loginUser interface{}
	var loginToken string

	if purpose == "login" {
		if h.authSvc == nil {
			respondJSON(w, http.StatusInternalServerError, map[string]any{
				"success": false,
				"error":   "auth service not available",
			})
			return
		}

		user, token, err := h.authSvc.LoginWithOAuth(stateData.Provider, userInfo)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]any{
				"success": false,
				"error":   fmt.Sprintf("failed to authenticate user: %v", err),
			})
			return
		}

		loginUser = user
		loginToken = token
		userID = user.ID
	}

	if userID == 0 {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "user identifier missing from OAuth2 state",
		})
		return
	}

	// Store OAuth2 data using user_id
	err = h.oauth2StorageSvc.StoreOAuth2Response(
		userID,
		stateData.Provider,
		tokenResp.AccessToken,
		tokenResp.RefreshToken,
		expiresAt,
		userInfoJSON,
	)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   fmt.Sprintf("failed to store OAuth2 data: %v", err),
		})
		return
	}

	responseData := map[string]any{
		"provider":     stateData.Provider,
		"user_info":    userInfo,
		"access_token": tokenResp.AccessToken,
		"token_type":   tokenResp.TokenType,
		"expires_in":   tokenResp.ExpiresIn,
		"message":      "OAuth2 data stored successfully",
		"callback_url": stateData.CallbackURL,
		"platform":     stateData.Platform,
	}

	if loginUser != nil {
		responseData["user"] = loginUser
		responseData["token"] = loginToken
	}

	// Return success with user info
	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    responseData,
	})
}

// GET /oauth2/providers/{userId}
func (h *OAuth2Handler) handleGetUserServices(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}

	// Extract userId from path
	path := req.URL.Path
	userIDStr := path[len("/oauth2/providers/"):]
	if userIDStr == "" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "user_id is required",
		})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid user_id",
		})
		return
	}

	providers, err := h.oauth2StorageSvc.GetUserServicesStatus(userID)
	if err != nil {
		fmt.Printf("Error getting user services status: %v\n", err)
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"providers": providers,
		},
	})
}

func (h *OAuth2Handler) handleGetProviderTokenByServiceByUserId(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}
	userIdStr := req.URL.Query().Get("user_id")
	serviceName := req.URL.Query().Get("service")
	if userIdStr == "" || serviceName == "" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "user_id and service query parameters are required",
		})
		return
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid user_id",
		})
		return
	}
	providerToken, err := h.oauth2StorageSvc.GetProviderTokenByServiceByUser(userId, serviceName)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
	}
	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"providerToken": providerToken,
		},
	})
}
