package router

import (
	"net/http"

	"github.com/raphael-guer1n/AREA/area-gateway/internal/config"
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
		_ = route
	}

	return rt.mux, nil
}
