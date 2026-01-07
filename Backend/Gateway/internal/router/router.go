package router

import (
	"net/http"
	"strings"

	"github.com/raphael-guer1n/AREA/area-gateway/internal/config"
	"github.com/raphael-guer1n/AREA/area-gateway/internal/core"
	"github.com/raphael-guer1n/AREA/area-gateway/internal/middleware"
	"github.com/raphael-guer1n/AREA/area-gateway/internal/registry"
)

type Router struct {
	registry   *registry.Registry
	config     *config.GatewayConfig
	authMW     *middleware.AuthMiddleware
	permMW     *middleware.PermissionsMiddleware
	internalMW *middleware.InternalMiddleware
	loggingMW  *middleware.LoggingMiddleware
	mux        *http.ServeMux
}

func NewRouter(
	reg *registry.Registry,
	cfg *config.GatewayConfig,
	auth *middleware.AuthMiddleware,
	perm *middleware.PermissionsMiddleware,
	internal *middleware.InternalMiddleware,
	logging *middleware.LoggingMiddleware,
) *Router {
	return &Router{
		registry:   reg,
		config:     cfg,
		authMW:     auth,
		permMW:     perm,
		internalMW: internal,
		loggingMW:  logging,
		mux:        http.NewServeMux(),
	}
}

func (rt *Router) Build() (*http.ServeMux, error) {

	routes := rt.registry.ListAllRoutes()

	proxies := make(map[string]http.Handler)
	for _, route := range routes {
		if _, exists := proxies[route.ServiceName]; exists {
			continue
		}
		proxies[route.ServiceName] = NewReverseProxy(route.BaseURL, "/"+route.ServiceName)
	}

	rt.mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, err := rt.registry.FindRouteByPath(r.URL.Path)
		if err != nil {
			core.WriteError(
				w,
				http.StatusNotFound,
				core.ErrNotFound,
				"Route not found",
			)
			return
		}

		if !methodAllowed(r.Method, route.Methods) {
			core.WriteError(
				w,
				http.StatusMethodNotAllowed,
				core.ErrForbidden,
				"Method not allowed",
			)
			return
		}

		proxy, ok := proxies[route.ServiceName]
		if !ok {
			proxy = NewReverseProxy(route.BaseURL, "/"+route.ServiceName)
			proxies[route.ServiceName] = proxy
		}

		handler := rt.permMW.Handler(proxy)
		handler = rt.authMW.Handler(handler)
		handler = rt.internalMW.Handler(handler)
		handler = rt.loggingMW.Handler(handler)

		handler.ServeHTTP(w, r)
	}))

	return rt.mux, nil
}

func methodAllowed(method string, allowed []string) bool {
	for _, m := range allowed {
		if strings.EqualFold(m, method) {
			return true
		}
	}
	return false
}
