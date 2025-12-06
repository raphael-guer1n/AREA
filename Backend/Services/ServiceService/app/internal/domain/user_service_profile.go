package domain

import (
	"encoding/json"
	"time"
)

type UserServiceProfile struct {
	ID             int
	UserId         int
	Service        string
	ProviderUserId string
	AccessToken    string
	RefreshToken   string
	ExpiresAt      time.Time
	RawProfile     json.RawMessage
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type UserServiceProfileRepository interface {
	Create(userId int, service, providerUserId, accessToken, refreshToken string, expiresAt time.Time, rawProfile json.RawMessage) (*UserServiceProfile, error)
	FindByUserId(userId int) ([]*UserServiceProfile, error)
}
