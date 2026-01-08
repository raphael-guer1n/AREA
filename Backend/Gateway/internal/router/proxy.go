package router

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/raphael-guer1n/AREA/area-gateway/internal/core"
)

func NewReverseProxy(baseURL string, stripPrefix string) http.Handler {
	target, err := url.Parse(baseURL)
	if err != nil {
		panic("invalid service base URL: " + baseURL)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   3 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		MaxIdleConns:          100,
		MaxConnsPerHost:       100,
	}

	proxy.Director = func(req *http.Request) {
		clientIP := clientIPFromRemoteAddr(req.RemoteAddr)

		appendForwardedHeader(req, "X-Forwarded-For", clientIP)
		appendForwardedHeader(req, "X-Forwarded-Host", req.Host)
		appendForwardedHeader(req, "X-Forwarded-Proto", schemeFromRequest(req))

		req.Header.Set("X-Real-IP", clientIP)
		req.Header.Del("X-Internal-Secret")
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host

		basePath := strings.TrimSuffix(target.Path, "/")
		incomingPath := req.URL.Path
		if !strings.HasPrefix(incomingPath, "/") {
			incomingPath = "/" + incomingPath
		}
		if stripPrefix != "" {
			sp := stripPrefix
			if !strings.HasPrefix(sp, "/") {
				sp = "/" + sp
			}
			sp = strings.TrimSuffix(sp, "/")
			if strings.HasPrefix(incomingPath, sp) {
				incomingPath = strings.TrimPrefix(incomingPath, sp)
				if incomingPath == "" {
					incomingPath = "/"
				}
			}
		}
		req.URL.Path = basePath + incomingPath
	}

	proxy.ModifyResponse = func(res *http.Response) error {
		removeHopByHopHeaders(res.Header)
		res.Header.Del("Server")
		res.Header.Del("X-Powered-By")
		return nil
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("proxy error for %s %s: %v", r.Method, r.URL.Path, err)

		if errors.Is(err, context.DeadlineExceeded) {
			core.WriteError(
				w,
				http.StatusGatewayTimeout,
				core.ErrGatewayTimeout,
				"Upstream service timeout",
			)
			return
		}

		var netErr net.Error
		if errors.As(err, &netErr) {
			if netErr.Timeout() {
				core.WriteError(
					w,
					http.StatusGatewayTimeout,
					core.ErrGatewayTimeout,
					"Upstream service timeout",
				)
				return
			}
			core.WriteError(
				w,
				http.StatusBadGateway,
				core.ErrBadGateway,
				"Upstream service unreachable",
			)
			return
		}

		core.WriteError(
			w,
			http.StatusBadGateway,
			core.ErrBadGateway,
			"Upstream service unreachable",
		)
	}

	return proxy
}

func clientIPFromRemoteAddr(addr string) string {
	if strings.HasPrefix(addr, "[") {
		end := strings.Index(addr, "]")
		if end != -1 {
			return addr[1:end]
		}
	}

	ip, _, err := net.SplitHostPort(addr)
	if err == nil {
		return ip
	}
	return addr
}

func appendForwardedHeader(req *http.Request, header string, value string) {
	if prev := req.Header.Get(header); prev != "" {
		req.Header.Set(header, prev+", "+value)
	} else {
		req.Header.Set(header, value)
	}
}

func schemeFromRequest(req *http.Request) string {
	if req.TLS != nil {
		return "https"
	}
	return "http"
}

func removeHopByHopHeaders(h http.Header) {
	hopByHop := []string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"TE",
		"Trailer",
		"Transfer-Encoding",
		"Upgrade",
	}
	for _, k := range hopByHop {
		h.Del(k)
	}
}
