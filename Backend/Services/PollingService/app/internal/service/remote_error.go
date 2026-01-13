package service

import (
	"encoding/json"
	"strings"
)

type remoteAPIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func parseRemoteError(raw json.RawMessage) (string, string) {
	if len(raw) == 0 {
		return "", ""
	}

	var msg string
	if err := json.Unmarshal(raw, &msg); err == nil {
		return strings.TrimSpace(msg), ""
	}

	var apiErr remoteAPIError
	if err := json.Unmarshal(raw, &apiErr); err == nil {
		message := strings.TrimSpace(apiErr.Message)
		if message == "" {
			message = strings.TrimSpace(apiErr.Code)
		}
		return message, strings.TrimSpace(apiErr.Code)
	}

	return strings.TrimSpace(string(raw)), ""
}
