package registry

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/raphael-guer1n/AREA/area-gateway/internal/config"
)

type RegisteredRoute struct {
	ServiceName    string
	BaseURL        string
	Path           string
	NamespacedPath string
	Methods        []string
	AuthRequired   bool
	Permissions    []string
	InternalOnly   bool
}

type Registry struct {
	services map[string]config.ServiceConfig
	routes   []RegisteredRoute
	mu       sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		services: make(map[string]config.ServiceConfig),
		routes:   []RegisteredRoute{},
	}
}

func (r *Registry) Load(services []config.ServiceConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, svc := range services {
		r.services[svc.Name] = svc

		for _, route := range svc.Routes {
			r.routes = append(r.routes, RegisteredRoute{
				ServiceName:    svc.Name,
				BaseURL:        svc.BaseURL,
				Path:           route.Path,
				NamespacedPath: "/" + svc.Name + route.Path,
				Methods:        route.Methods,
				AuthRequired:   route.AuthRequired,
				Permissions:    route.Permissions,
				InternalOnly:   route.InternalOnly,
			})
		}
	}

	return nil
}

func (r *Registry) GetService(name string) (config.ServiceConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	svc, ok := r.services[name]
	return svc, ok
}

func (r *Registry) ListServices() []config.ServiceConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []config.ServiceConfig
	for _, svc := range r.services {
		list = append(list, svc)
	}
	return list
}

func (r *Registry) FindRoute(path, method string) (*RegisteredRoute, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Exact match first.
	for _, rt := range r.routes {
		if rt.Path == path || rt.NamespacedPath == path {
			if slices.Contains(rt.Methods, method) {
				return &rt, nil
			}
		}
	}

	// Fallback for routes that allow dynamic suffixes (e.g. /oauth2/providers/{id}).
	for _, rt := range r.routes {
		if !slices.Contains(rt.Methods, method) {
			continue
		}

		if strings.HasPrefix(path, rt.Path) && len(path) > len(rt.Path) && strings.HasPrefix(path[len(rt.Path):], "/") {
			return &rt, nil
		}

		if strings.HasPrefix(path, rt.NamespacedPath) && len(path) > len(rt.NamespacedPath) && strings.HasPrefix(path[len(rt.NamespacedPath):], "/") {
			return &rt, nil
		}
	}

	return nil, fmt.Errorf("no route found for %s %s", method, path)
}

func (r *Registry) ListAllRoutes() []RegisteredRoute {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return append([]RegisteredRoute{}, r.routes...)
}
