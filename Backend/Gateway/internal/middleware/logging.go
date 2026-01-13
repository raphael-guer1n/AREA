package middleware

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/raphael-guer1n/AREA/area-gateway/internal/registry"
)

func firstHeaderValue(headerValue string) string {
	if headerValue == "" {
		return ""
	}
	parts := strings.Split(headerValue, ",")
	return strings.TrimSpace(parts[0])
}

func requestSource(r *http.Request) string {
	if origin := strings.TrimSpace(r.Header.Get("Origin")); origin != "" {
		return origin
	}

	if ref := strings.TrimSpace(r.Header.Get("Referer")); ref != "" {
		if u, err := url.Parse(ref); err == nil {
			if u.Scheme != "" && u.Host != "" {
				return u.Scheme + "://" + u.Host
			}
		}
		return ref
	}

	if xff := firstHeaderValue(r.Header.Get("X-Forwarded-For")); xff != "" {
		return xff
	}

	return r.RemoteAddr
}

type wrappedResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *wrappedResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

type LoggingMiddleware struct {
	reg *registry.Registry
}

func NewLoggingMiddleware(reg *registry.Registry) *LoggingMiddleware {
	return &LoggingMiddleware{
		reg: reg,
	}
}

func (lm *LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		ww := &wrappedResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		internalFlag := ""
		if rt, err := lm.reg.FindRoute(r.URL.Path, r.Method); err == nil && rt.InternalOnly {
			internalFlag = "[INTERNAL] "
		}

		user := GetUserFromContext(r.Context())
		userInfo := ""
		if user != nil && user.UserID != "" {
			userInfo = "(auth:" + user.UserID + ")"
		}

		level := "[INFO]"
		if ww.status >= 400 {
			level = "[ERROR]"
		}

		log.Printf(
			"%s %s%s %d %s src=%s %dms",
			level,
			internalFlag,
			r.Method+" "+r.URL.Path,
			ww.status,
			userInfo,
			requestSource(r),
			duration.Milliseconds(),
		)
	})
}
