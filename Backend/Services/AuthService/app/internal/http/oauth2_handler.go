package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/raphael-guer1n/AREA/AuthService/internal/oauth2"
	"github.com/raphael-guer1n/AREA/AuthService/internal/service"
)

type OAuth2Handler struct {
	oauth2StorageSvc *service.OAuth2StorageService
	oauth2Manager    *oauth2.Manager
}

func NewOAuth2Handler(oauth2StorageSvc *service.OAuth2StorageService, oauth2Manager *oauth2.Manager) *OAuth2Handler {
	return &OAuth2Handler{
		oauth2StorageSvc: oauth2StorageSvc,
		oauth2Manager:    oauth2Manager,
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

// GET /auth/oauth2/authorize?provider=<provider_name>&user_id=<user_id>&callback_url=<url>&platform=<platform>
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
	userIDStr := req.URL.Query().Get("user_id")
	callbackURL := req.URL.Query().Get("callback_url")
	platform := req.URL.Query().Get("platform")

	if provider == "" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "provider parameter is required",
		})
		return
	}

	if userIDStr == "" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "user_id parameter is required",
		})
		return
	}

	// Parse user ID
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid user_id parameter",
		})
		return
	}

	// Default platform to "web" if not specified
	if platform == "" {
		platform = "web"
	}

	authURL, err := h.oauth2Manager.GetAuthURL(provider, userID, callbackURL, platform)
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

	// Store OAuth2 data using user_id from StateData
	err = h.oauth2StorageSvc.StoreOAuth2Response(
		stateData.UserID,
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

	// Return success with user info
	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"provider":     stateData.Provider,
			"user_info":    userInfo,
			"access_token": tokenResp.AccessToken,
			"token_type":   tokenResp.TokenType,
			"expires_in":   tokenResp.ExpiresIn,
			"message":      "OAuth2 data stored successfully",
			"callback_url": stateData.CallbackURL,
			"platform":     stateData.Platform,
		},
	})
}

// GET /oauth2/services/{userId}
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
	userIDStr := path[len("/oauth2/services/"):]
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
