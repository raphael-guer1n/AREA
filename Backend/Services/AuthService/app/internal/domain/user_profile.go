package domain

import (
	"encoding/json"
	"time"
)

type UserProfile struct {
	ID               int             `json:"id"`
	UserId           int             `json:"user_id"`
	Service          string          `json:"service"`
	ProviderUserId   string          `json:"provider_user_id"`
	AccessToken      string          `json:"access_token"`
	RefreshToken     string          `json:"refresh_token"`
	ExpiresAt        time.Time       `json:"expires_at"`
	RawProfile       json.RawMessage `json:"raw_profile"`
	NeedsReconnect   bool            `json:"needs_reconnect"`
	LastRefreshError *string         `json:"last_refresh_error,omitempty"`
	LastRefreshAt    *time.Time      `json:"last_refresh_at,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

type ServiceStatus struct {
	Service        string
	NeedsReconnect bool
}

type RefreshCandidate struct {
	ID           int
	UserId       int
	Service      string
	RefreshToken string
	ExpiresAt    time.Time
}

type UserProfileRepository interface {
	Create(userId int, service, providerUserId, accessToken, refreshToken string, expiresAt time.Time, rawProfile json.RawMessage) (UserProfile, error)
	GetServicesStatusByUserId(userId int) ([]ServiceStatus, error)
	GetProviderUserTokenByServiceByUserId(userId int, service string) (string, error)
	GetProviderProfileProfileByServiceByUser(userId int, service string) (UserProfile, error)
	ListRefreshCandidates(expireBefore time.Time) ([]RefreshCandidate, error)
	UpdateTokens(profileID int, accessToken, refreshToken string, expiresAt time.Time) error
	MarkNeedsReconnect(profileID int, reason string) error
}
