package domain

import (
	"encoding/json"
	"time"
)

type Subscription struct {
	ID             int             `json:"id"`
	HookID         string          `json:"hook_id"`
	ProviderHookID string          `json:"provider_hook_id,omitempty"`
	UserID         int             `json:"user_id"`
	AreaID         int             `json:"area_id"`
	Provider       string          `json:"provider"`
	Config         json.RawMessage `json:"config"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type SubscriptionRepository interface {
	Create(sub *Subscription) (*Subscription, error)
	FindByHookID(hookID string) (*Subscription, error)
	ListByUserID(userID int) ([]Subscription, error)
	UpdateProviderHookID(hookID, providerHookID string) error
	DeleteByHookID(hookID string) error
}
