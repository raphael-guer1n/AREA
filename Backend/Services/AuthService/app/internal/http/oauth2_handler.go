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

func NewOAuth2Handler(oauth2StorageSvc *service.OAuth2StorageService, oauth2Manager *oauth2.Manager, authSvc *service.AuthService) *OAuth2Handler {
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
	expiresAt := oauth2.ResolveExpiresAt(storeReq.ExpiresIn)

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
	callbackURL := req.URL.Query().Get("callback_url")
	platform := req.URL.Query().Get("platform")

	if provider == "" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "provider parameter is required",
		})
		return
	}

	userID, err := getUserIDFromRequest(req)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

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

// GET /loginwith?provider=<provider_name>&callback_url=<url>&platform=<platform>
// Starts an OAuth2 login flow that will create or connect a user on callback
func (h *OAuth2Handler) handleLoginWithAuthorize(w http.ResponseWriter, req *http.Request) {
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

	// Use userID = 0 to indicate sign-in/sign-up flow
	authURL, err := h.oauth2Manager.GetAuthURL(provider, 0, callbackURL, platform)
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
	redirectOverride := req.URL.Query().Get("redirect_uri")

	userInfo, tokenResp, stateData, err := h.oauth2Manager.HandleCallback(state, code, redirectOverride)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	userInfoJSON, err := json.Marshal(userInfo.RawData)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "failed to marshal user info",
		})
		return
	}

	expiresAt := oauth2.ResolveExpiresAt(tokenResp.ExpiresIn)

	var userIDForStorage int = stateData.UserID
	var jwtToken string
	if stateData.UserID == 0 {
		email := strings.TrimSpace(userInfo.Email)
		username := strings.TrimSpace(userInfo.Username)
		if username == "" {
			username = strings.TrimSpace(userInfo.Name)
		}
		if username == "" {
			username = fmt.Sprintf("%s_user", stateData.Provider)
		}
		if email == "" {
			if userInfo.ID != "" {
				email = fmt.Sprintf("%s_%s@oauth.local", stateData.Provider, userInfo.ID)
			} else {
				email = fmt.Sprintf("%s_%d@oauth.local", stateData.Provider, time.Now().Unix())
			}
		}

		existingUser, findErr := h.authSvc.GetUserByEmail(email)
		if findErr == nil && existingUser != nil {
			token, genErr := auth.GenerateToken(existingUser.ID)
			if genErr != nil {
				respondJSON(w, http.StatusInternalServerError, map[string]any{
					"success": false,
					"error":   genErr.Error(),
				})
				return
			}
			userIDForStorage = existingUser.ID
			jwtToken = token
		} else {
			randPass := fmt.Sprintf("oauth_%d_%s", time.Now().UnixNano(), userInfo.ID)
			newUser, token, regErr := h.authSvc.Register(email, username, randPass)
			if regErr != nil {
				respondJSON(w, http.StatusInternalServerError, map[string]any{
					"success": false,
					"error":   regErr.Error(),
				})
				return
			}
			userIDForStorage = newUser.ID
			jwtToken = token
		}
	}

	err = h.oauth2StorageSvc.StoreOAuth2Response(
		userIDForStorage,
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

	// NEW: handle mobile platforms by redirecting to deep link
	if stateData.Platform == "android" || stateData.Platform == "ios" {
		redirect := fmt.Sprintf("area://auth?provider=%s&code=%s&state=%s&token=%s",
			stateData.Provider, code, state, jwtToken)

		html := fmt.Sprintf(`
			<!DOCTYPE html>
			<html>
			<head><meta charset="utf-8"><title>AREA Redirect</title></head>
			<body style="font-family:sans-serif;text-align:center;margin-top:2em;">
				<h2>Returning to AREA app...</h2>
				<p>If you are not redirected, <a href="%s">tap here</a>.</p>
				<script>window.onload=function(){window.location="%s";}</script>
			</body>
			</html>`, redirect, redirect)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(html))
		return
	}

	respData := map[string]any{
		"provider":     stateData.Provider,
		"user_info":    userInfo,
		"access_token": tokenResp.AccessToken,
		"token_type":   tokenResp.TokenType,
		"expires_in":   tokenResp.ExpiresIn,
		"message":      "OAuth2 data stored successfully",
		"callback_url": stateData.CallbackURL,
		"platform":     stateData.Platform,
	}
	if jwtToken != "" {
		respData["token"] = jwtToken
	}
	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    respData,
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
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"providerToken": providerToken,
		},
	})
}

func (h *OAuth2Handler) handleGetProviderProfileByServiceByUserId(w http.ResponseWriter, req *http.Request) {
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

	userProfile, err := h.oauth2StorageSvc.GetProviderProfileByServiceByUser(userId, serviceName)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	userFields, err := h.oauth2StorageSvc.GetProviderFieldsByProfileId(userProfile.ID)
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
			"userProfile": userProfile,
			"userFields":  userFields,
		},
	})
	return
}
