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

	for _, route := range routes {
		proxy := NewReverseProxy(route.BaseURL)

		var handler http.Handler = proxy

		handler = rt.permMW.Handler(handler)
		handler = rt.authMW.Handler(handler)
		handler = rt.internalMW.Handler(handler)
		handler = rt.loggingMW.Handler(handler)

		methods := make(map[string]struct{})
		for _, m := range route.Methods {
			methods[strings.ToUpper(m)] = struct{}{}
		}

		finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := methods[r.Method]; !ok {
				core.WriteError(
					w,
					http.StatusMethodNotAllowed,
					core.ErrForbidden,
					"Method not allowed",
				)
				return
			}
			handler.ServeHTTP(w, r)
		})

		rt.mux.Handle(route.Path, finalHandler)
	}

	rt.mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		core.WriteError(
			w,
			http.StatusNotFound,
			core.ErrNotFound,
			"Route not found",
		)
	}))

	return rt.mux, nil
}
