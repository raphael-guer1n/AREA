package config

type ServiceConfig struct {
	Name    string        `json:"name"`
	BaseURL string        `json:"base_url"`
	Routes  []RouteConfig `json:"routes"`
}

type RouteConfig struct {
	Path    string   `json:"path"`
	Methods []string `json:"methods"`

	AuthRequired bool     `json:"auth_required"`
	Permissions  []string `json:"permissions,omitempty"`

	InternalOnly bool `json:"internal_only"`

	RateLimit *RateLimitConfig `json:"rate_limit,omitempty"`
}

type RateLimitConfig struct {
	Requests int `json:"requests"`
	PerSec   int `json:"per_sec"`
}
