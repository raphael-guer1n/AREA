package middleware

import (
	"net/http"

	"github.com/raphael-guer1n/AREA/area-gateway/internal/core"
	"github.com/raphael-guer1n/AREA/area-gateway/internal/registry"
)

type PermissionsMiddleware struct {
	reg *registry.Registry
}

func NewPermissionsMiddleware(reg *registry.Registry) *PermissionsMiddleware {
	return &PermissionsMiddleware{
		reg: reg,
	}
}

func (pm *PermissionsMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		rt, err := pm.reg.FindRoute(r.URL.Path, r.Method)
		if err != nil {
			core.WriteError(w, http.StatusInternalServerError, core.ErrNotFound, "Route not found in registry")
			return
		}

		if !rt.AuthRequired {
			next.ServeHTTP(w, r)
			return
		}

		user := GetUserFromContext(r.Context())
		if user == nil {
			core.WriteError(w, http.StatusUnauthorized, core.ErrUnauthorized, "Authentication required")
			return
		}

		if len(rt.Permissions) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		userPerms := make(map[string]bool, len(user.Permissions))
		for _, p := range user.Permissions {
			userPerms[p] = true
		}

		for _, required := range rt.Permissions {
			if !userPerms[required] {
				core.WriteError(
					w,
					http.StatusForbidden,
					core.ErrForbidden,
					"Missing required permission: "+required,
				)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
