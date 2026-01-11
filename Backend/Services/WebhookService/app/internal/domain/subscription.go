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
	ActionID       int             `json:"action_id"`
	Provider       string          `json:"provider"`
	Service        string          `json:"service"`
	Active         bool            `json:"active"`
	Config         json.RawMessage `json:"config"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type SubscriptionRepository interface {
	Create(sub *Subscription) (*Subscription, error)
	FindByHookID(hookID string) (*Subscription, error)
	FindByActionID(actionID int) (*Subscription, error)
	ListByUserID(userID int) ([]Subscription, error)
	ListByProvider(provider string) ([]Subscription, error)
	UpdateByActionID(sub *Subscription) (*Subscription, error)
	UpdateProviderHookID(hookID, providerHookID string) error
	TouchByHookID(hookID string) error
	DeleteByActionID(actionID int) error
}
