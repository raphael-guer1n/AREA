package http

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Router struct {
	mux            *http.ServeMux
	actionHandler  *ActionHandler
	webhookHandler *WebhookHandler
	logAll         bool
}

func NewRouter(actionHandler *ActionHandler, webhookHandler *WebhookHandler, logAll bool) *Router {
	r := &Router{
		mux:            http.NewServeMux(),
		actionHandler:  actionHandler,
		webhookHandler: webhookHandler,
		logAll:         logAll,
	}

	r.routes()
	return r
}

func (r *Router) routes() {
	r.mux.HandleFunc("/health", r.handleHealth)

	// Action webhook management (intended for internal use via AreaService)
	r.mux.HandleFunc("/actions", r.actionHandler.HandleActions)
	r.mux.HandleFunc("/actions/", r.actionHandler.HandleAction)
	r.mux.HandleFunc("/activate/", r.actionHandler.HandleActivateAction)
	r.mux.HandleFunc("/deactivate/", r.actionHandler.HandleDeactivateAction)

	// Webhook receivers
	r.mux.HandleFunc("/webhooks/", r.webhookHandler.HandleReceiveWebhook)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var bodyBytes []byte
	if r.logAll && req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
		_ = req.Body.Close()
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	recorder := &responseRecorder{ResponseWriter: w}
	r.mux.ServeHTTP(recorder, req)
	if recorder.status == 0 {
		recorder.status = http.StatusOK
	}

	log.Printf("request method=%s url=%s status=%d", req.Method, req.URL.RequestURI(), recorder.status)
	if r.logAll {
		log.Printf("request headers=%v body=%s", req.Header, string(bodyBytes))
		log.Printf("response headers=%v body=%s", recorder.Header(), recorder.body.String())
	}
}

func (r *Router) handleHealth(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]string{
			"status": "healthy",
		},
	})
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Write(data []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	_, _ = r.body.Write(data)
	return r.ResponseWriter.Write(data)
}
