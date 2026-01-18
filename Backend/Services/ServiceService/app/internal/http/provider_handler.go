package http

import (
	"net"
	"net/http"
	"sort"
	"strings"
	"time"

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
	providers := h.providerConfigSvc.GetAllProviderSummaries()

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"services":  services,
			"providers": providers,
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

func (h *ProviderHandler) HandleGetAboutJSON(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "Method not allowed",
		})
		return
	}

	serviceNames := h.providerConfigSvc.GetAllServicesNames()
	sort.Strings(serviceNames)

	services := make([]aboutService, 0, len(serviceNames))
	for _, serviceName := range serviceNames {
		serviceConfig, exists := h.providerConfigSvc.GetServiceConfig(serviceName)
		if !exists {
			continue
		}

		actions := make([]aboutServiceItem, 0, len(serviceConfig.Actions))
		for _, action := range serviceConfig.Actions {
			description := action.Label
			if description == "" {
				description = action.Title
			}
			actions = append(actions, aboutServiceItem{
				Name:        action.Title,
				Description: description,
			})
		}

		reactions := make([]aboutServiceItem, 0, len(serviceConfig.Reactions))
		for _, reaction := range serviceConfig.Reactions {
			description := reaction.Label
			if description == "" {
				description = reaction.Title
			}
			reactions = append(reactions, aboutServiceItem{
				Name:        reaction.Title,
				Description: description,
			})
		}

		services = append(services, aboutService{
			Name:      serviceConfig.Name,
			Actions:   actions,
			Reactions: reactions,
		})
	}

	resp := aboutResponse{
		Client: aboutClient{
			Host: clientIPFromRequest(req),
		},
		Server: aboutServer{
			CurrentTime: time.Now().Unix(),
			Services:    services,
		},
	}

	respondJSON(w, http.StatusOK, resp)
}

type aboutResponse struct {
	Client aboutClient `json:"client"`
	Server aboutServer `json:"server"`
}

type aboutClient struct {
	Host string `json:"host"`
}

type aboutServer struct {
	CurrentTime int64          `json:"current_time"`
	Services    []aboutService `json:"services"`
}

type aboutService struct {
	Name      string             `json:"name"`
	Actions   []aboutServiceItem `json:"actions"`
	Reactions []aboutServiceItem `json:"reactions"`
}

type aboutServiceItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func clientIPFromRequest(req *http.Request) string {
	if forwarded := req.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		if len(parts) > 0 {
			if ip := strings.TrimSpace(parts[0]); ip != "" {
				return ip
			}
		}
	}

	if realIP := strings.TrimSpace(req.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}

	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil && host != "" {
		return host
	}

	return req.RemoteAddr
}
