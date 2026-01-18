package domain

import "time"

type InputField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type OutputField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Action struct {
	ActionID  int          `json:"action_id"`
	Active    bool         `json:"active"`
	Type      string       `json:"type"`
	Provider  string       `json:"provider"`
	Service   string       `json:"service"`
	Title     string       `json:"title"`
	Input     []InputField `json:"input"`
	CronJobID *int         `json:"-"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type CreateActionsRequest struct {
	Actions []Action `json:"actions"`
}

type TriggerAreaRequest struct {
	ActionID     int           `json:"action_id"`
	OutputFields []OutputField `json:"output_fields"`
}
