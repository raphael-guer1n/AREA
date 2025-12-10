package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/raphael-guer1n/AREA/AreaService/internal/domain"
	"github.com/raphael-guer1n/AREA/AreaService/internal/service"
)

type AreaHandler struct {
	authSvc *service.AreaService
}

func NewAreaHandler(authSvc *service.AreaService) *AreaHandler {
	return &AreaHandler{
		authSvc: authSvc,
	}
}

func (h *AreaHandler) HandleCreateEventArea(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}
	userId, err := getUserId(req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "Error getting user ID," + err.Error(),
		})
		return
	}
	token, err := getUserServiceToken(req, userId, "google")
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "Error getting user service token, the user is not linked to Google," + err.Error(),
		})
		return
	}
	var body struct {
		Delay int          `json:"delay"`
		Event domain.Event `json:"event"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid request body",
		})
		return
	}
	event, err := h.authSvc.CreateCalendarEvent(token, body.Event)
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
			"event": event,
		},
	})
}

func getUserServiceToken(r *http.Request, userId int, service string) (string, error) {
	baseURL := "http://area_auth_api:8083/oauth2/provider/token/"
	params := url.Values{}
	params.Add("user_id", fmt.Sprintf("%d", userId))
	params.Add("service", service)

	fullURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var body struct {
		Data struct {
			Token string `json:"providerToken"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	log.Default().Printf("Got token for user %d from service %s: %s", userId, service, body.Data.Token)
	return body.Data.Token, nil
}

func getUserId(r *http.Request) (int, error) {
	req, err := http.NewRequest(http.MethodGet, "http://area_auth_api:8083/auth/me", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var authResp struct {
		Data struct {
			User struct {
				ID int `json:"id"`
			} `json:"user"`
		} `json:"data"`
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return 0, err
	}
	return authResp.Data.User.ID, nil
}
