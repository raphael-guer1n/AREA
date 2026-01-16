package domain

import (
	"encoding/json"
	"time"
)

type Subscription struct {
	ID              int             `json:"id"`
	UserID          int             `json:"user_id"`
	ActionID        int             `json:"action_id"`
	Provider        string          `json:"provider"`
	Service         string          `json:"service"`
	Active          bool            `json:"active"`
	Config          json.RawMessage `json:"config"`
	IntervalSeconds int             `json:"interval_seconds"`
	LastItemID      string          `json:"last_item_id,omitempty"`
	LastPolledAt    *time.Time      `json:"last_polled_at,omitempty"`
	NextRunAt       *time.Time      `json:"next_run_at,omitempty"`
	LastError       string          `json:"last_error,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type SubscriptionRepository interface {
	Create(sub *Subscription) (*Subscription, error)
	FindByActionID(actionID int) (*Subscription, error)
	ListDue(now time.Time) ([]Subscription, error)
	UpdateByActionID(sub *Subscription) (*Subscription, error)
	UpdatePollingState(actionID int, lastItemID string, nextRunAt time.Time, lastError string, lastPolledAt time.Time) error
	DeleteByActionID(actionID int) error
}
