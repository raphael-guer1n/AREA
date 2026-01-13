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
	log.Println(req.URL.String())
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
	area, err := h.areaService.SaveArea(body)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	err = h.TriggerAction(area.Actions, area.Active)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
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
	switch areaAction.Type {
	case "cron":
		err := h.TriggerCronAction(areaAction)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("action type %s not supported", areaAction.Type)
	}
}

func (h *AreaHandler) DeactivateAction(req *http.Request, areaAction domain.AreaAction) error {
	switch areaAction.Type {
	case "cron":
		return nil
	default:
		return fmt.Errorf("action type %s not supported", areaAction.Type)
	}
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
	serviceProfile, err := h.getUserServiceProfile(userId, areaReaction.Provider)
	if err != nil {
		return err
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
		for _, serviceField := range serviceProfile.Fields {
			field.Value = strings.ReplaceAll(field.Value, "{{"+serviceField.FieldKey+"}}", serviceField.StringValue)
		}
		fieldValues[field.Name] = field.Value
	}
	for _, serviceField := range serviceProfile.Fields {
		fieldValues[serviceField.FieldKey] = serviceField.StringValue
	}
	return h.areaService.LaunchReactions(serviceProfile.Profile.AccessToken, fieldValues, reactionConfig)
}

func (h *AreaHandler) TriggerAction(areaAction []domain.AreaAction, isActive bool) error {
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
	for _, action := range otherActions {
		createUrl, exist := h.cfg.CreateActionsUrls[action.Type]
		log.Printf("url configs : %v", h.cfg.CreateActionsUrls)
		log.Printf("url configs : %v", h.cfg.CreateActionsUrls["webhook"])
		if !exist {
			return fmt.Errorf("action type %s not supported", action.Type)
		}
		var body struct {
			Actions []domain.AreaAction `json:"actions"`
		}
		body.Actions = areaAction
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}
		resp, err := http.NewRequest(http.MethodPost, createUrl, bytes.NewBuffer(payload))
		if err != nil {
			return err
		}
		if h.cfg.InternalSecret != "" {
			resp.Header.Set("Authorization", "Bearer "+h.cfg.InternalSecret)
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
