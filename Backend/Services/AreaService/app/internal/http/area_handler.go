package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/raphael-guer1n/AREA/AreaService/internal/config"
	"github.com/raphael-guer1n/AREA/AreaService/internal/domain"
	"github.com/raphael-guer1n/AREA/AreaService/internal/service"
)

type AreaHandler struct {
	areaService *service.AreaService
	cfg         config.Config
}

func NewAreaHandler(authSvc *service.AreaService, cfg config.Config) *AreaHandler {
	return &AreaHandler{
		areaService: authSvc,
		cfg:         cfg,
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
	userId, err := h.getUserId(req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "Error getting user ID," + err.Error(),
		})
		return
	}
	token, err := h.getUserServiceToken(userId, "google")
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

	if body.Delay > 0 {
		time.AfterFunc(time.Duration(body.Delay)*time.Second, func() {
			event, err := h.areaService.CreateCalendarEvent(token, body.Event)
			if err != nil {
				log.Printf("Error creating delayed calendar event: %v", err)
				return
			}
			log.Printf("Delayed calendar event created successfully: %s", event.Summary)
		})

		respondJSON(w, http.StatusAccepted, map[string]any{
			"success": true,
			"message": fmt.Sprintf("Event will be created in %d seconds", body.Delay),
		})
		return
	}

	event, err := h.areaService.CreateCalendarEvent(token, body.Event)
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

func (h *AreaHandler) getUserServiceProfile(userId int, service string) (domain.UserService, error) {
	baseURL := strings.TrimRight(h.cfg.AuthServiceURL, "/") + "/oauth2/provider/profile/"
	params := url.Values{}
	params.Add("user_id", fmt.Sprintf("%d", userId))
	params.Add("service", service)

	fullURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return domain.UserService{}, err
	}
	if h.cfg.InternalSecret != "" {
		req.Header.Set("X-Internal-Secret", h.cfg.InternalSecret)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return domain.UserService{}, err
	}
	defer resp.Body.Close()
	var body struct {
		Data domain.UserService `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return domain.UserService{}, err
	}
	return body.Data, nil
}

func (h *AreaHandler) getUserServiceToken(userId int, service string) (string, error) {
	baseURL := strings.TrimRight(h.cfg.AuthServiceURL, "/") + "/oauth2/provider/token/"
	params := url.Values{}
	params.Add("user_id", fmt.Sprintf("%d", userId))
	params.Add("service", service)

	fullURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return "", err
	}
	if h.cfg.InternalSecret != "" {
		req.Header.Set("X-Internal-Secret", h.cfg.InternalSecret)
	}
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
	return body.Data.Token, nil
}

func (h *AreaHandler) getAreaConfiguration(area domain.Area) (domain.AreaConfig, error) {
	var areaConfig domain.AreaConfig

	for _, action := range area.Actions {
		actionConfig, err := h.getActionDetails(action)
		if err != nil {
			return areaConfig, err
		}
		areaConfig.Actions = append(areaConfig.Actions, actionConfig)
	}
	for _, reaction := range area.Reactions {
		reactionConfig, err := h.getReactionDetails(reaction)
		if err != nil {
			return areaConfig, err
		}
		areaConfig.Reactions = append(areaConfig.Reactions, reactionConfig)
	}
	return areaConfig, nil
}

func (h *AreaHandler) getActionDetails(action domain.AreaAction) (domain.ActionConfig, error) {
	var actionConfig domain.ActionConfig

	baseUrl := strings.TrimRight(h.cfg.ServiceServiceURL, "/") + "/services/service-config"
	params := url.Values{}
	params.Add("service", action.Service)
	fullUrl := baseUrl + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
	if err != nil {
		return actionConfig, err
	}
	if h.cfg.InternalSecret != "" {
		req.Header.Set("X-Internal-Secret", h.cfg.InternalSecret)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return actionConfig, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return actionConfig, fmt.Errorf("failed to get service config: %w", err)
	}
	var body struct {
		Data struct {
			Actions []domain.ActionConfig `json:"actions"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return actionConfig, err
	}
	var exist bool
	for _, actConfig := range body.Data.Actions {
		if actConfig.Title == action.Title {
			actionConfig = actConfig
			exist = true
			break
		}
	}
	if !exist {
		return actionConfig, fmt.Errorf("action not found")
	}
	return actionConfig, nil
}

func (h *AreaHandler) getReactionDetails(reaction domain.AreaReaction) (domain.ReactionConfig, error) {
	var reactionConfig domain.ReactionConfig

	baseUrl := strings.TrimRight(h.cfg.ServiceServiceURL, "/") + "/services/service-config"
	params := url.Values{}
	params.Add("service", reaction.Service)
	fullUrl := baseUrl + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
	if err != nil {
		return reactionConfig, err
	}
	if h.cfg.InternalSecret != "" {
		req.Header.Set("X-Internal-Secret", h.cfg.InternalSecret)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return reactionConfig, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return reactionConfig, fmt.Errorf("failed to get service config")
	}
	var body struct {
		Data struct {
			Reactions []domain.ReactionConfig `json:"reactions"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return reactionConfig, err
	}
	var exist bool
	for _, reactConfig := range body.Data.Reactions {
		if reactConfig.Title == reaction.Title {
			reactionConfig = reactConfig
			exist = true
			break
		}
	}
	if !exist {
		return reactionConfig, fmt.Errorf("reaction not found")
	}
	return reactionConfig, nil
}

func (h *AreaHandler) getUserId(r *http.Request) (int, error) {
	endpoint := strings.TrimRight(h.cfg.AuthServiceURL, "/") + "/auth/me"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
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
	if authResp.Data.User.ID == 0 {
		return 0, errors.New("user not found")
	}
	return authResp.Data.User.ID, nil
}

func (h *AreaHandler) checkUserProviderConnections(userId int, area domain.Area) ([]string, error) {
	providersNeeded := make(map[string]bool)

	for _, action := range area.Actions {
		if action.Provider != "" {
			providersNeeded[action.Provider] = true
		}
	}

	for _, reaction := range area.Reactions {
		if reaction.Provider != "" {
			providersNeeded[reaction.Provider] = true
		}
	}

	if len(providersNeeded) == 0 {
		return nil, nil
	}

	baseURL := strings.TrimRight(h.cfg.AuthServiceURL, "/") + fmt.Sprintf("/oauth2/providers/%d", userId)
	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, err
	}
	if h.cfg.InternalSecret != "" {
		req.Header.Set("X-Internal-Secret", h.cfg.InternalSecret)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body struct {
		Success bool `json:"success"`
		Data    struct {
			Providers []struct {
				Provider         string `json:"provider"`
				IsLogged         bool   `json:"is_logged"`
				NeedReconnecting bool   `json:"need_reconnecting"`
			} `json:"providers"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	var missingProviders []string
	for provider := range providersNeeded {
		found := false
		for _, p := range body.Data.Providers {
			if p.Provider == provider && p.IsLogged && !p.NeedReconnecting {
				found = true
				break
			}
		}
		if !found {
			missingProviders = append(missingProviders, provider)
		}
	}

	return missingProviders, nil
}

func (h *AreaHandler) SaveArea(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}
	var body domain.Area
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid request body " + err.Error(),
		})
		return
	}
	userId, err := h.getUserId(req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "Error getting user ID," + err.Error(),
		})
		return
	}
	if userId == 0 {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "Error getting user ID",
		})
		return
	}
	body.UserID = userId
	areaConfig, err := h.getAreaConfiguration(body)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	err = CheckAreaValidity(body, areaConfig)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	missingProviders, err := h.checkUserProviderConnections(userId, body)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "Error checking provider connections: " + err.Error(),
		})
		return
	}

	area, err := h.areaService.SaveArea(body)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	err = h.TriggerAction(area.Actions, area.Active, req.Header.Get("Authorization"))
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if len(missingProviders) > 0 {
		respondJSON(w, http.StatusOK, map[string]any{
			"success":           true,
			"message":           "Area saved but some provider connections are missing; actions may not run until you connect them.",
			"missing_providers": missingProviders,
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{})
}

func (h *AreaHandler) HandleActivateArea(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}
	var body struct {
		AreaId int `json:"area_id"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid request body " + err.Error(),
		})
		return
	}
	userId, err := h.getUserId(req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "Error getting user ID," + err.Error(),
		})
		return
	}
	area, err := h.areaService.GetArea(body.AreaId)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if area.UserID != userId {
		respondJSON(w, http.StatusForbidden, map[string]any{
			"success": false,
			"error":   "You are not allowed to activate this area",
		})
		return
	}
	if area.Active == true {
		respondJSON(w, http.StatusOK, map[string]any{
			"success": false,
			"error":   "Area already active",
		})
		return
	}

	missingProviders, err := h.checkUserProviderConnections(userId, area)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "Error checking provider connections: " + err.Error(),
		})
		return
	}

	if len(missingProviders) > 0 {
		respondJSON(w, http.StatusOK, map[string]any{
			"success":           false,
			"message":           "Cannot activate area due to missing provider connections",
			"missing_providers": missingProviders,
		})
		return
	}

	err = h.areaService.ToggleArea(body.AreaId, true)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	for _, action := range area.Actions {
		err := h.ActivateAction(req, action)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]any{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	}
	respondJSON(w, http.StatusOK, map[string]any{})
}

func (h *AreaHandler) HandleDeactivateArea(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}
	var body struct {
		AreaId int `json:"area_id"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid request body " + err.Error(),
		})
		return
	}
	userId, err := h.getUserId(req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{})
	}
	area, err := h.areaService.GetArea(body.AreaId)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if area.UserID != userId {
		respondJSON(w, http.StatusForbidden, map[string]any{
			"success": false,
			"error":   "You are not allowed to deactivate this area",
		})
		return
	}
	if area.Active == false {
		respondJSON(w, http.StatusOK, map[string]any{
			"success": false,
			"error":   "Area already inactive",
		})
		return
	}
	err = h.areaService.ToggleArea(body.AreaId, false)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	for _, action := range area.Actions {
		err := h.DeactivateAction(req, action)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]any{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	}
	respondJSON(w, http.StatusOK, map[string]any{})
}

func (h *AreaHandler) ActivateAction(req *http.Request, areaAction domain.AreaAction) error {
	if areaAction.Type == "cron" {
		err := h.TriggerCronAction(areaAction)
		if err != nil {
			return err
		}
		return nil
	}

	activateURL, exists := h.cfg.ActivateActionsUrls[areaAction.Type]
	if !exists {
		return fmt.Errorf("action type %s not supported or not configured", areaAction.Type)
	}

	if activateURL == "nil" || activateURL == "" {
		return nil
	}

	activateURL = strings.TrimRight(activateURL, "/") + fmt.Sprintf("/%d", areaAction.ID)

	activateReq, err := http.NewRequest(http.MethodPost, activateURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create activate request: %w", err)
	}

	if authHeader := req.Header.Get("Authorization"); authHeader != "" {
		activateReq.Header.Set("Authorization", authHeader)
	}
	activateReq.Header.Set("X-Internal-Secret", h.cfg.InternalSecret)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(activateReq)
	if err != nil {
		return fmt.Errorf("failed to activate %s action: %w", areaAction.Type, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]any
		json.NewDecoder(resp.Body).Decode(&errResp)
		if errMsg, ok := errResp["error"].(string); ok {
			return fmt.Errorf("failed to activate %s action: %s", areaAction.Type, errMsg)
		}
		return fmt.Errorf("failed to activate %s action: status %d", areaAction.Type, resp.StatusCode)
	}

	return nil
}

func (h *AreaHandler) DeactivateAction(req *http.Request, areaAction domain.AreaAction) error {
	deactivateURL, exists := h.cfg.DeactivateActionsUrls[areaAction.Type]
	if !exists {
		return fmt.Errorf("action type %s not supported or not configured", areaAction.Type)
	}

	if deactivateURL == "nil" || deactivateURL == "" {
		return nil
	}

	deactivateURL = strings.TrimRight(deactivateURL, "/") + fmt.Sprintf("/%d", areaAction.ID)

	deactivateReq, err := http.NewRequest(http.MethodPost, deactivateURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create deactivate request: %w", err)
	}

	if authHeader := req.Header.Get("Authorization"); authHeader != "" {
		deactivateReq.Header.Set("Authorization", authHeader)
	}

	deactivateReq.Header.Set("X-Internal-Secret", h.cfg.InternalSecret)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(deactivateReq)
	if err != nil {
		return fmt.Errorf("failed to deactivate %s action: %w", areaAction.Type, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]any
		json.NewDecoder(resp.Body).Decode(&errResp)
		if errMsg, ok := errResp["error"].(string); ok {
			return fmt.Errorf("failed to deactivate %s action: %s", areaAction.Type, errMsg)
		}
		return fmt.Errorf("failed to deactivate %s action: status %d", areaAction.Type, resp.StatusCode)
	}

	return nil
}

func (h *AreaHandler) GetAreas(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}
	userId, err := h.getUserId(req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "Error getting user ID," + err.Error(),
		})
		return
	}
	if userId == 0 {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "Error getting user ID",
		})
		return
	}
	areas, err := h.areaService.GetUserAreas(userId)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    areas,
	})
}

func (h *AreaHandler) HandleActionTrigger(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}
	var body struct {
		ActionId     int                 `json:"action_id"`
		OutputFields []domain.InputField `json:"output_fields"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid request body " + err.Error(),
		})
		return
	}
	area, err := h.areaService.GetAreaFromAction(body.ActionId)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	userId := area.UserID
	if userId == 0 {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "missing user for action",
		})
		return
	}
	if area.Active != true {
		respondJSON(w, http.StatusOK, map[string]any{
			"success": false,
			"error":   "Inactive area, you can't trigger actions on it",
		})
		return
	}
	for _, reaction := range area.Reactions {
		err := h.TriggerReaction(reaction, body.OutputFields, userId)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]any{
				"success": false,
				"error":   err.Error(),
			})
		}
	}
	respondJSON(w, http.StatusOK, map[string]any{})
	return
}

func checkFieldsValidity(fields []domain.InputField, config []domain.FieldConfig) error {
	for _, field := range config {
		if field.Required {
			var found bool
			for _, inputField := range fields {
				if inputField.Name == field.Name && inputField.Value != "" {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("field %s is required", field.Name)
			}
		}
	}
	return nil
}

func CheckAreaValidity(area domain.Area, config domain.AreaConfig) error {
	for i, action := range area.Actions {
		configAction := config.Actions[i]
		err := checkFieldsValidity(action.Input, configAction.Fields)
		if err != nil {
			return err
		}
	}
	for i, reaction := range area.Reactions {
		configReaction := config.Reactions[i]
		err := checkFieldsValidity(reaction.Input, configReaction.Fields)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *AreaHandler) TriggerReaction(areaReaction domain.AreaReaction, outputFields []domain.InputField, userId int) error {
	serviceProfile := domain.UserService{}
	var err error
	if strings.TrimSpace(areaReaction.Provider) != "" {
		serviceProfile, err = h.getUserServiceProfile(userId, areaReaction.Provider)
		if err != nil {
			return err
		}
	}
	reactionConfig, err := h.getReactionDetails(areaReaction)
	if err != nil {
		return err
	}

	fieldValues := make(map[string]string)
	for _, field := range areaReaction.Input {
		for _, outputField := range outputFields {
			field.Value = strings.ReplaceAll(field.Value, "{{"+outputField.Name+"}}", outputField.Value)
		}
		if strings.TrimSpace(areaReaction.Provider) != "" {
			for _, serviceField := range serviceProfile.Fields {
				field.Value = strings.ReplaceAll(field.Value, "{{"+serviceField.FieldKey+"}}", serviceField.StringValue)
			}
		}
		fieldValues[field.Name] = field.Value
	}
	if strings.TrimSpace(areaReaction.Provider) != "" {
		for _, serviceField := range serviceProfile.Fields {
			fieldValues[serviceField.FieldKey] = serviceField.StringValue
		}
	}
	userToken := ""
	if strings.TrimSpace(areaReaction.Provider) != "" {
		userToken = serviceProfile.Profile.AccessToken
	}
	return h.areaService.LaunchReactions(userToken, fieldValues, reactionConfig)
}

type actionRequest struct {
	Active   bool                `json:"active"`
	ActionID int                 `json:"action_id"`
	Type     string              `json:"type"`
	Provider string              `json:"provider"`
	Service  string              `json:"service"`
	Title    string              `json:"title"`
	Input    []domain.InputField `json:"input"`
}

func (h *AreaHandler) TriggerAction(areaAction []domain.AreaAction, isActive bool, authHeader string) error {
	cronActions := make([]domain.AreaAction, 0)
	otherActions := make([]domain.AreaAction, 0)
	for _, action := range areaAction {
		action.Active = isActive
		if action.Type == "cron" {
			cronActions = append(cronActions, action)
		} else {
			otherActions = append(otherActions, action)
		}
	}
	if isActive {
		for _, action := range cronActions {
			err := h.TriggerCronAction(action)
			if err != nil {
				return err
			}
		}
	}

	actionsByType := make(map[string][]actionRequest)
	for _, action := range otherActions {
		action.Active = isActive
		actionsByType[action.Type] = append(actionsByType[action.Type], actionRequest{
			Active:   action.Active,
			ActionID: action.ID,
			Type:     action.Type,
			Provider: action.Provider,
			Service:  action.Service,
			Title:    action.Title,
			Input:    action.Input,
		})
	}

	for actionType, actions := range actionsByType {
		createUrl, exist := h.cfg.CreateActionsUrls[actionType]
		if !exist {
			return fmt.Errorf("action type %s not supported", actionType)
		}
		var body struct {
			Actions []actionRequest `json:"actions"`
		}
		body.Actions = actions
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}
		resp, err := http.NewRequest(http.MethodPost, createUrl, bytes.NewBuffer(payload))
		if err != nil {
			return err
		}
		if strings.TrimSpace(authHeader) != "" {
			resp.Header.Set("Authorization", authHeader)
		} else if h.cfg.InternalSecret != "" {
			resp.Header.Set("Authorization", "Bearer "+h.cfg.InternalSecret)
		}
		if h.cfg.InternalSecret != "" {
			resp.Header.Set("X-Internal-Secret", h.cfg.InternalSecret)
		}
		client := &http.Client{}
		_, err = client.Do(resp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *AreaHandler) TriggerCronAction(areaAction domain.AreaAction) error {
	if areaAction.Type != "cron" {
		return fmt.Errorf("only cron actions are supported")
	}
	var delay int
	for _, input := range areaAction.Input {
		if input.Name == "delay" {
			delay, _ = strconv.Atoi(input.Value)
		}
	}
	var outputFields []domain.InputField
	outputFields = append(outputFields, domain.InputField{Name: "delay", Value: strconv.Itoa(delay)})
	log.Printf("Triggering action %d in %d seconds", areaAction.ID, delay)
	time.AfterFunc(time.Duration(delay)*time.Second, func() {
		var body struct {
			OutputFields []domain.InputField `json:"output_fields"`
			ActionId     int                 `json:"action_id"`
		}
		body.OutputFields = outputFields
		body.ActionId = areaAction.ID
		payload, err := json.Marshal(body)
		if err != nil {
			log.Fatal(err)
			return
		}
		endpoint := strings.TrimRight(h.cfg.AreaServiceURL, "/") + "/triggerArea"
		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(payload))
		if err != nil {
			log.Fatal(err)
		}
		if h.cfg.InternalSecret != "" {
			req.Header.Set("X-Internal-Secret", h.cfg.InternalSecret)
		}
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
			return
		}
	})
	return nil
}

func (h *AreaHandler) HandleDeleteArea(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}
	var body struct {
		AreaId int `json:"area_id"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid request body " + err.Error(),
		})
		return
	}
	userId, err := h.getUserId(req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "Error getting user ID," + err.Error(),
		})
		return
	}
	area, err := h.areaService.GetArea(body.AreaId)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	log.Printf("Deleting area of %d by user %d", area.UserID, userId)
	if area.UserID != userId {
		respondJSON(w, http.StatusForbidden, map[string]any{
			"success": false,
			"error":   "You are not allowed to delete this area",
		})
		return
	}
	for _, action := range area.Actions {
		delUrl, exist := h.cfg.DelActionsUrls[action.Type]
		if !exist || delUrl == "nil" {
			continue
		}
		deleteUrl := delUrl + "/" + strconv.Itoa(action.ID)
		delReq, err := http.NewRequest(http.MethodDelete, deleteUrl, nil)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]any{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		if authHeader := req.Header.Get("Authorization"); authHeader != "" {
			delReq.Header.Set("Authorization", authHeader)
		}
		if h.cfg.InternalSecret != "" {
			delReq.Header.Set("X-Internal-Secret", h.cfg.InternalSecret)
		}
		client := &http.Client{}
		_, err = client.Do(delReq)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]any{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
	}
	err = h.areaService.DeleteArea(body.AreaId)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{})
}

func (h *AreaHandler) HandleDeactivateAreasByProvider(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}
	var body struct {
		UserId   int    `json:"user_id"`
		Provider string `json:"provider"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid request body " + err.Error(),
		})
		return
	}
	if body.UserId == 0 {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "user_id is required",
		})
		return
	}
	if body.Provider == "" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "provider is required",
		})
		return
	}
	deactivatedCount, err := h.areaService.DeactivateAreasByProvider(body.UserId, body.Provider)
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
			"deactivated_count": deactivatedCount,
		},
	})
}
