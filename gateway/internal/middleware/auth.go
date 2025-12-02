package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/raphael-guer1n/AREA/area-gateway/internal/config"
	"github.com/raphael-guer1n/AREA/area-gateway/internal/core"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userContextKey = contextKey("user")

type UserContext struct {
	UserID      string
	Email       string
	Permissions []string
}

func GetUserFromContext(ctx context.Context) *UserContext {
	v := ctx.Value(userContextKey)
	if v == nil {
		return nil
	}
	uc, ok := v.(*UserContext)
	if !ok {
		return nil
	}
	return uc
}

type AuthMiddleware struct {
	Algorithm string
	PublicKey []byte // RS256
	Secret    []byte // HS256
}

func NewAuthMiddleware(cfg *config.GatewayConfig) *AuthMiddleware {
	return &AuthMiddleware{
		Algorithm: cfg.JwtAlgorithm,
		PublicKey: []byte(cfg.JwtPublicKey),
		Secret:    []byte(cfg.JwtSecret),
	}
}

func (a *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth := r.Header.Get("Authorization")
		if auth == "" {
			core.WriteError(w, http.StatusUnauthorized, core.ErrMissingToken, "Missing Authorization header")
			return
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			core.WriteError(w, http.StatusUnauthorized, core.ErrInvalidAuthHeader, "Authorization must be Bearer <token>")
			return
		}

		rawToken := parts[1]

		if a.Algorithm == "RS256" && len(a.PublicKey) == 0 {
			core.WriteError(w, 500, core.ErrInternalError, "RS256 requires JWT_PUBLIC_KEY")
			return
		}

		if a.Algorithm == "HS256" && len(a.Secret) == 0 {
			core.WriteError(w, 500, core.ErrInternalError, "HS256 requires JWT_SECRET")
			return
		}

		token, err := jwt.Parse(rawToken, func(t *jwt.Token) (interface{}, error) {

			if t.Method.Alg() != a.Algorithm {
				return nil, fmt.Errorf("unexpected signing method: %s", t.Method.Alg())
			}

			if a.Algorithm == "RS256" {
				key, err := jwt.ParseRSAPublicKeyFromPEM(a.PublicKey)
				if err != nil {
					return nil, errors.New("invalid RSA public key")
				}
				return key, nil
			}

			if a.Algorithm == "HS256" {
				return a.Secret, nil
			}

			return nil, fmt.Errorf("unsupported JWT algorithm: %s", a.Algorithm)
		})

		if err != nil || !token.Valid {
			core.WriteError(w, http.StatusUnauthorized, core.ErrInvalidToken, "Invalid or expired token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			core.WriteError(w, http.StatusUnauthorized, core.ErrInvalidToken, "Token claims malformed")
			return
		}

		exp, ok := claims["exp"].(float64)
		if !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
			core.WriteError(w, http.StatusUnauthorized, core.ErrInvalidToken, "Token expired")
			return
		}

		uid, _ := claims["sub"].(string)
		email, _ := claims["email"].(string)

		var perms []string
		if arr, ok := claims["permissions"].([]interface{}); ok {
			for _, p := range arr {
				if s, ok := p.(string); ok {
					perms = append(perms, s)
				}
			}
		}

		userCtx := &UserContext{
			UserID:      uid,
			Email:       email,
			Permissions: perms,
		}

		ctx := context.WithValue(r.Context(), userContextKey, userCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}