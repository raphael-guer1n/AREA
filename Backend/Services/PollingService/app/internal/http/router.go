package http

import (
	"encoding/json"
	"log"
	"net/http"
)

type Router struct {
	mux           *http.ServeMux
	actionHandler *ActionHandler
}

func NewRouter(actionHandler *ActionHandler) *Router {
	r := &Router{
		mux:           http.NewServeMux(),
		actionHandler: actionHandler,
	}

	r.routes()
	return r
}

func (r *Router) routes() {
	r.mux.HandleFunc("/health", r.handleHealth)

	// Polling action management (intended for internal use via AreaService)
	r.mux.HandleFunc("/actions", r.actionHandler.HandleActions)
	r.mux.HandleFunc("/actions/", r.actionHandler.HandleAction)
	r.mux.HandleFunc("/activate/", r.actionHandler.HandleActivateAction)
	r.mux.HandleFunc("/deactivate/", r.actionHandler.HandleDeactivateAction)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	recorder := &statusRecorder{ResponseWriter: w}
	r.mux.ServeHTTP(recorder, req)
	if recorder.status == 0 {
		recorder.status = http.StatusOK
	}
	log.Printf("request method=%s url=%s status=%d", req.Method, req.URL.RequestURI(), recorder.status)
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

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(data []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(data)
}
