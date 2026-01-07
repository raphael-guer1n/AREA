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

	for _, rt := range r.routes {
		if method != "" && !slices.Contains(rt.Methods, method) {
			continue
		}

		if pathMatchesPattern(path, rt.Path) {
			return &rt, nil
		}

		if pathMatchesPattern(path, rt.NamespacedPath) {
			return &rt, nil
		}
	}

	return nil, fmt.Errorf("no route found for %s %s", method, path)
}

func (r *Registry) FindRouteByPath(path string) (*RegisteredRoute, error) {
	return r.FindRoute(path, "")
}

func pathMatchesPattern(path, pattern string) bool {
	if pattern == "" {
		return false
	}

	if path == pattern {
		return true
	}

	trimmedPath := strings.Trim(path, "/")
	trimmedPattern := strings.Trim(pattern, "/")

	if trimmedPath == trimmedPattern {
		return true
	}

	pathSegments := splitPathSegments(trimmedPath)
	patternSegments := splitPathSegments(trimmedPattern)

	if len(pathSegments) != len(patternSegments) {
		return false
	}

	for i, segment := range patternSegments {
		if isParamSegment(segment) {
			if pathSegments[i] == "" {
				return false
			}
			continue
		}
		if segment != pathSegments[i] {
			return false
		}
	}

	return true
}

func splitPathSegments(path string) []string {
	if path == "" {
		return []string{}
	}
	return strings.Split(path, "/")
}

func isParamSegment(segment string) bool {
	return strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") && len(segment) > 2
}

func (r *Registry) ListAllRoutes() []RegisteredRoute {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return append([]RegisteredRoute{}, r.routes...)
}
