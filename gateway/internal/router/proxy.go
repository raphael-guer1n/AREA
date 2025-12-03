package router

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(baseURL string) http.Handler {
	target, err := url.Parse(baseURL)
	if err != nil {
		panic("invalid service base URL: " + baseURL)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Forwarded-Proto", req.URL.Scheme)
		req.Header.Set("X-Forwarded-For", req.RemoteAddr)

		req.Host = target.Host
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

	return proxy
}
