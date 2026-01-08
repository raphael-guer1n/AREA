package domain

import (
	"encoding/json"
	"time"
)

type UserProfile struct {
	ID             int             `json:"id"`
	UserId         int             `json:"user_id"`
	Service        string          `json:"service"`
	ProviderUserId string          `json:"provider_user_id"`
	AccessToken    string          `json:"access_token"`
	RefreshToken   string          `json:"refresh_token"`
	ExpiresAt      time.Time       `json:"expires_at"`
	RawProfile     json.RawMessage `json:"raw_profile"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type UserProfileRepository interface {
	Create(userId int, service, providerUserId, accessToken, refreshToken string, expiresAt time.Time, rawProfile json.RawMessage) (UserProfile, error)
	GetServicesByUserId(userId int) ([]string, error)
	GetProviderUserTokenByServiceByUserId(userId int, service string) (string, error)
	GetProviderProfileProfileByServiceByUser(userId int, service string) (UserProfile, error)
}
