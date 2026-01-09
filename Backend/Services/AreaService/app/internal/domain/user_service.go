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

type UserServiceField struct {
	ID          int              `json:"id"`
	ProfileId   int              `json:"profile_id"`
	FieldKey    string           `json:"field_key"`
	StringValue string           `json:"string_value"`
	NumberValue float64          `json:"number_value"`
	BoolValue   bool             `json:"bool_value"`
	JsonValue   *json.RawMessage `json:"json_value"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type UserService struct {
	Profile UserProfile        `json:"userProfile"`
	Fields  []UserServiceField `json:"userFields"`
}
