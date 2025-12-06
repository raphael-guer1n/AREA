package http

import (
	"encoding/json"
	"net/http"

	"github.com/raphael-guer1n/AREA/ServiceService/internal/service"
)

type Router struct {
	mux     *http.ServeMux
	userSvc *service.UserService
}

func NewRouter(userSvc *service.UserService) *Router {
	r := &Router{
		mux:     http.NewServeMux(),
		userSvc: userSvc,
	}

	r.routes()
	return r
}

func (r *Router) routes() {
	r.mux.HandleFunc("/health", r.handleHealth)
	r.mux.HandleFunc("/users", r.handleUsers)
	r.mux.HandleFunc("/users/create", r.handleCreateUser)
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

func (r *Router) handleUsers(w http.ResponseWriter, _ *http.Request) {
	_, err := r.userSvc.ListUsers()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
	} else {
		respondJSON(w, http.StatusOK, map[string]any{
			"success": true,
		})
	}
}

func (r *Router) handleCreateUser(w http.ResponseWriter, _ *http.Request) {
	_, err := r.userSvc.CreateUser("", "John", "Doe")
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
	} else {
		respondJSON(w, http.StatusOK, map[string]any{
			"success": true,
		})
	}
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
