package http

import (
	"net/http"

	"github.com/raphael-guer1n/AREA/ServiceService/internal/service"
)

type ProviderHandler struct {
	providerConfigSvc *service.ProviderConfigService
}

func NewProviderHandler(providerConfigSvc *service.ProviderConfigService) *ProviderHandler {
	return &ProviderHandler{
		providerConfigSvc: providerConfigSvc,
	}
}

func (h *ProviderHandler) HandleGetProviders(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "Method not allowed",
		})
		return
	}

	services := h.providerConfigSvc.GetAllProvidersNames()

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"services": services,
		},
	})
}

func (h *ProviderHandler) HandleGetOAuth2Config(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "Method not allowed",
		})
		return
	}

	serviceName := req.URL.Query().Get("service")
	if serviceName == "" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "service query parameter is required",
		})
		return
	}

	oauth2Config, exists := h.providerConfigSvc.GetOAuth2Config(serviceName)
	if !exists {
		respondJSON(w, http.StatusNotFound, map[string]any{
			"success": false,
			"error":   "Service not found",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    oauth2Config,
	})
}

func (h *ProviderHandler) HandleGetProviderConfig(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "Method not allowed",
		})
		return
	}

	serviceName := req.URL.Query().Get("service")
	if serviceName == "" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "service query parameter is required",
		})
		return
	}

	providerConfig, exists := h.providerConfigSvc.GetProviderConfig(serviceName)
	if !exists {
		respondJSON(w, http.StatusNotFound, map[string]any{
			"success": false,
			"error":   "Service not found",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    providerConfig,
	})
}

func (h *ProviderHandler) HandleGetServiceConfig(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "Method not allowed",
		})
		return
	}
	serviceName := req.URL.Query().Get("service")
	if serviceName == "" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "service query parameter is required",
		})
	}
	serviceConfig, exists := h.providerConfigSvc.GetServiceConfig(serviceName)
	if !exists {
		respondJSON(w, http.StatusNotFound, map[string]any{
			"success": false,
			"error":   "Service not found",
		})
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    serviceConfig,
	})
}

func (h *ProviderHandler) HandleGetServices(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "Method not allowed",
		})
		return
	}
	servicesNames := h.providerConfigSvc.GetAllServicesNames()
	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"services": servicesNames,
		},
	})
}
