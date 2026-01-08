package middleware

import (
	"net/http"

	"github.com/raphael-guer1n/AREA/area-gateway/internal/config"
	"github.com/raphael-guer1n/AREA/area-gateway/internal/core"
	"github.com/raphael-guer1n/AREA/area-gateway/internal/registry"
)

type InternalMiddleware struct {
	reg            *registry.Registry
	internalSecret string
}

func NewInternalMiddleware(cfg *config.GatewayConfig, reg *registry.Registry) *InternalMiddleware {
	return &InternalMiddleware{
		reg:            reg,
		internalSecret: cfg.InternalSecret,
	}
}

func (im *InternalMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		rt, err := im.reg.FindRoute(r.URL.Path, r.Method)
		if err != nil {
			core.WriteError(w, http.StatusInternalServerError, core.ErrNotFound, "Route not found in registry")
			return
		}

		if !rt.InternalOnly {
			next.ServeHTTP(w, r)
			return
		}

		if im.internalSecret == "" {
			core.WriteError(w, http.StatusInternalServerError, core.ErrInternalError, "Internal secret not configured")
			return
		}

		headerSecret := r.Header.Get("X-Internal-Secret")
		if headerSecret == "" || headerSecret != im.internalSecret {
			core.WriteError(
				w,
				http.StatusForbidden,
				core.ErrForbidden,
				"Access to internal endpoint denied",
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}
