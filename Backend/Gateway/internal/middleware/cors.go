package middleware

import (
	"net/http"
	"strings"

	"github.com/raphael-guer1n/AREA/area-gateway/internal/config"
)

type CORSMiddleware struct {
	allowedOrigins map[string]struct{}
	allowAll       bool
}

func NewCORSMiddleware(cfg *config.GatewayConfig) *CORSMiddleware {
	origins := make(map[string]struct{})
	allowAll := false

	for _, o := range cfg.AllowedOrigins {
		trimmed := strings.TrimSpace(o)
		if trimmed == "" {
			continue
		}
		if trimmed == "*" {
			allowAll = true
			continue
		}
		origins[trimmed] = struct{}{}
	}

	return &CORSMiddleware{
		allowedOrigins: origins,
		allowAll:       allowAll,
	}
}

func (c *CORSMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		origin := r.Header.Get("Origin")
		if origin != "" {
			if c.allowAll {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Vary", "Origin")
			} else if _, ok := c.allowedOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Vary", "Origin")
			}
		}

		w.Header().Set(
			"Access-Control-Allow-Methods",
			"GET, POST, PUT, PATCH, DELETE, OPTIONS",
		)

		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Authorization",
		)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
