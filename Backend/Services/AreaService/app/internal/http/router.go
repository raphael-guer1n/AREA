package http

import (
	"encoding/json"
	"net/http"
)

type Router struct {
	mux         *http.ServeMux
	areaHandler *AreaHandler
}

func NewRouter(areaHandler *AreaHandler) *Router {
	r := &Router{
		mux:         http.NewServeMux(),
		areaHandler: areaHandler,
	}

	r.routes()
	return r
}

func (r *Router) routes() {
	r.mux.HandleFunc("/health", r.handleHealth)
	r.mux.HandleFunc("/createEvent", r.areaHandler.HandleCreateEventArea)
	r.mux.HandleFunc("/saveArea", r.areaHandler.SaveArea)
	r.mux.HandleFunc("/getAreas", r.areaHandler.GetAreas)
	r.mux.HandleFunc("/triggerArea", r.areaHandler.HandleActionTrigger)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
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
