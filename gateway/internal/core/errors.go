package core

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

func WriteError(w http.ResponseWriter, status int, code string, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := ErrorResponse{
		Error: APIError{
			Code:    code,
			Message: msg,
		},
	}

	json.NewEncoder(w).Encode(resp)
}

const (
	ErrUnauthorized      = "unauthorized"
	ErrForbidden         = "forbidden"
	ErrInvalidToken      = "invalid_token"
	ErrMissingToken      = "missing_token"
	ErrInvalidAuthHeader = "invalid_auth_header"
	ErrInternalError     = "internal_error"
)
