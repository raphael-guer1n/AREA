package http

import (
	"net/http"

	"github.com/raphael-guer1n/AREA/ServiceService/internal/service"
)

type PollingProviderHandler struct {
	providerConfigSvc *service.PollingProviderConfigService
}

func NewPollingProviderHandler(providerConfigSvc *service.PollingProviderConfigService) *PollingProviderHandler {
	return &PollingProviderHandler{
		providerConfigSvc: providerConfigSvc,
	}
}

func (h *PollingProviderHandler) HandleGetProviders(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}

	providers := h.providerConfigSvc.GetAllProviderNames()

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"providers": providers,
		},
	})
}

func (h *PollingProviderHandler) HandleGetProviderConfig(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}

	providerName := req.URL.Query().Get("provider")
	if providerName == "" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "provider query parameter is required",
		})
		return
	}

	providerConfig, exists := h.providerConfigSvc.GetProviderConfig(providerName)
	if !exists {
		respondJSON(w, http.StatusNotFound, map[string]any{
			"success": false,
			"error":   "provider not found",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    providerConfig,
	})
}
