package http

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/raphael-guer1n/AREA/AuthService/internal/service"
)

type OAuth2Handler struct {
	oauth2StorageSvc *service.OAuth2StorageService
}

func NewOAuth2Handler(oauth2StorageSvc *service.OAuth2StorageService) *OAuth2Handler {
	return &OAuth2Handler{
		oauth2StorageSvc: oauth2StorageSvc,
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
