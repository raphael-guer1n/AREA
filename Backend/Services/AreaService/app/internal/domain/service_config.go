package domain

type FieldConfig struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	Label         string `json:"label"`
	Required      bool   `json:"required"`
	DefaultValuer string `json:"default"`
}

type OutputFieldConfig struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Label string `json:"label"`
}

type ActionConfig struct {
	Title        string              `json:"title"`
	Label        string              `json:"label"`
	Type         string              `json:"type"`
	Fields       []FieldConfig       `json:"fields"`
	OutputFields []OutputFieldConfig `json:"output_fields"`
}

type ReactionConfig struct {
	Title      string        `json:"title"`
	Label      string        `json:"label"`
	Url        string        `json:"url"`
	Fields     []FieldConfig `json:"fields"`
	Method     string        `json:"method"`
	BodyType   string        `json:"bodyType"`
	BodyStruct []BodyField   `json:"body_struct"`
	Headers    map[string]string `json:"headers,omitempty"`
}

type ServiceConfig struct {
	Provider  string           `json:"provider"`
	Name      string           `json:"name"`
	IconURL   string           `json:"icon_url"`
	LogoURL   string           `json:"logo_url"`
	Label     string           `json:"label"`
	Actions   []ActionConfig   `json:"actions"`
	Reactions []ReactionConfig `json:"reactions"`
}

type AreaConfig struct {
	Actions   []ActionConfig
	Reactions []ReactionConfig
}
