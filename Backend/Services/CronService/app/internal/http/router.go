package http

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(actionHandler *ActionHandler, logAllRequests bool) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/actions", actionHandler.CreateActions).Methods("POST")
	router.HandleFunc("/activate/{actionId}", actionHandler.ActivateAction).Methods("POST")
	router.HandleFunc("/deactivate/{actionId}", actionHandler.DeactivateAction).Methods("POST")
	router.HandleFunc("/actions/{actionId}", actionHandler.DeleteAction).Methods("DELETE")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	if logAllRequests {
		return loggingMiddleware(router)
	}

	return router
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
