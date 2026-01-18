package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/raphael-guer1n/AREA/CronService/internal/domain"
	"github.com/raphael-guer1n/AREA/CronService/internal/service"
)

type ActionHandler struct {
	cronService *service.CronService
}

func NewActionHandler(cronService *service.CronService) *ActionHandler {
	return &ActionHandler{
		cronService: cronService,
	}
}

func (h *ActionHandler) CreateActions(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateActionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Actions) == 0 {
		http.Error(w, "No actions provided", http.StatusBadRequest)
		return
	}

	for i := range req.Actions {
		action := &req.Actions[i]
		if err := h.cronService.CreateAction(action); err != nil {
			log.Printf("Failed to create action %d: %v", action.ActionID, err)
			http.Error(w, "Failed to create action: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Actions created successfully"})
}

func (h *ActionHandler) ActivateAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	actionIDStr := vars["actionId"]

	actionID, err := strconv.Atoi(actionIDStr)
	if err != nil {
		http.Error(w, "Invalid action ID", http.StatusBadRequest)
		return
	}

	if err := h.cronService.ActivateAction(actionID); err != nil {
		log.Printf("Failed to activate action %d: %v", actionID, err)
		http.Error(w, "Failed to activate action: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Action activated successfully"})
}

func (h *ActionHandler) DeactivateAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	actionIDStr := vars["actionId"]

	actionID, err := strconv.Atoi(actionIDStr)
	if err != nil {
		http.Error(w, "Invalid action ID", http.StatusBadRequest)
		return
	}

	if err := h.cronService.DeactivateAction(actionID); err != nil {
		log.Printf("Failed to deactivate action %d: %v", actionID, err)
		http.Error(w, "Failed to deactivate action: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Action deactivated successfully"})
}

func (h *ActionHandler) DeleteAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	actionIDStr := vars["actionId"]

	actionID, err := strconv.Atoi(actionIDStr)
	if err != nil {
		http.Error(w, "Invalid action ID", http.StatusBadRequest)
		return
	}

	if err := h.cronService.DeleteAction(actionID); err != nil {
		log.Printf("Failed to delete action %d: %v", actionID, err)
		http.Error(w, "Failed to delete action: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Action deleted successfully"})
}
