package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/raphael-guer1n/AREA/AuthService/internal/auth"
	"github.com/raphael-guer1n/AREA/AuthService/internal/service"
)

var (
	errMissingAuthorizationHeader = errors.New("missing authorization header")
	errInvalidAuthorizationHeader = errors.New("invalid authorization header format")
	errInvalidOrExpiredToken      = errors.New("invalid or expired token")
)

func getUserIDFromRequest(req *http.Request) (int, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return 0, errMissingAuthorizationHeader
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, errInvalidAuthorizationHeader
	}

	userID, err := auth.ValidateToken(parts[1])
	if err != nil {
		return 0, errInvalidOrExpiredToken
	}

	return userID, nil
}

func getUserIDFromAuth(req *http.Request, authSvc *service.AuthService) (int, error) {
	return getUserIDFromRequest(req)
}
