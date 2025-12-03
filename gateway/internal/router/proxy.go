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
		req.Host = target.Host
	}

	return proxy
}
