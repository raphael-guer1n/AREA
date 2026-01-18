package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/raphael-guer1n/AREA/MailService/internal/config"
	"github.com/raphael-guer1n/AREA/MailService/internal/service"
)

type Router struct {
	mux     *http.ServeMux
	mailer  *service.Mailer
}

func NewRouter(mailer *service.Mailer, _ config.Config) *Router {
	r := &Router{
		mux:     http.NewServeMux(),
		mailer:  mailer,
	}

	r.routes()
	return r
}

func (r *Router) routes() {
	r.mux.HandleFunc("/health", r.handleHealth)
	r.mux.HandleFunc("/send", r.handleSend)
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

type sendRequest struct {
	To      any    `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (r *Router) handleSend(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"success": false,
			"error":   "method not allowed",
		})
		return
	}

	var body sendRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "invalid request body",
		})
		return
	}

	recipients := parseRecipients(body.To)
	subject := strings.TrimSpace(body.Subject)
	message := strings.TrimSpace(body.Body)
	if len(recipients) == 0 || subject == "" || message == "" {
		respondJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "to, subject and body are required",
		})
		return
	}

	if err := r.mailer.Send(recipients, subject, message); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
	})
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func parseRecipients(raw any) []string {
	recipients := []string{}
	switch v := raw.(type) {
	case string:
		recipients = splitRecipients(v)
	case []any:
		for _, item := range v {
			if str, ok := item.(string); ok {
				recipients = append(recipients, strings.TrimSpace(str))
			}
		}
	}

	filtered := recipients[:0]
	for _, recipient := range recipients {
		trimmed := strings.TrimSpace(recipient)
		if trimmed != "" {
			filtered = append(filtered, trimmed)
		}
	}
	return filtered
}

func splitRecipients(raw string) []string {
	raw = strings.ReplaceAll(raw, ";", ",")
	parts := strings.Split(raw, ",")
	recipients := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			recipients = append(recipients, trimmed)
		}
	}
	return recipients
}
