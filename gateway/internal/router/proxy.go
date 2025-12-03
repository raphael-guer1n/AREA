package router

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewReverseProxy(baseURL string) http.Handler {
	target, err := url.Parse(baseURL)
	if err != nil {
		panic("invalid service base URL: " + baseURL)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

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
		req.URL.Path = basePath + incomingPath
	}

	proxy.ModifyResponse = func(res *http.Response) error {
		removeHopByHopHeaders(res.Header)
		res.Header.Del("Server")
		res.Header.Del("X-Powered-By")
		return nil
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
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
