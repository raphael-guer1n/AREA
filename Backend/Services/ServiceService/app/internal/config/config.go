package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	HTTPPort string
	DBHost   string
	DBPort   string
	DBUser   string
	DBPass   string
	DBName   string
}

func Load() Config {
	return Config{
		HTTPPort: getEnv("SERVER_PORT", "8080"),
		DBHost:   getEnv("DB_HOST", "localhost"),
		DBPort:   getEnv("DB_PORT", "5432"),
		DBUser:   getEnv("DB_USER", "postgres"),
		DBPass:   getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "myservice_db"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

type OAuth2Config struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	AuthURL      string   `json:"auth_url"`
	TokenURL     string   `json:"token_url"`
	RedirectURI  string   `json:"redirect_uri"`
	Scopes       []string `json:"scopes"`
	UserInfoURL  string   `json:"user_info_url"`
}

type MappingConfig struct {
	FieldKey string `json:"field_key"`
	JSONPath string `json:"json_path"`
	Type     string `json:"type"`
	Optional bool   `json:"optional,omitempty"`
}

type ProviderConfig struct {
	Name     string          `json:"name"`
	OAuth2   OAuth2Config    `json:"oauth2"`
	Mappings []MappingConfig `json:"mappings"`
}

type WebhookSignatureConfig struct {
	Type                      string `json:"type"`
	Header                    string `json:"header"`
	Prefix                    string `json:"prefix"`
	SecretJSONPath            string `json:"secret_json_path"`
	Algorithm                 string `json:"algorithm,omitempty"`
	Encoding                  string `json:"encoding,omitempty"`
	SigningStringTemplate     string `json:"signing_string_template,omitempty"`
	TimestampHeader           string `json:"timestamp_header,omitempty"`
	TimestampToleranceSeconds int    `json:"timestamp_tolerance_seconds,omitempty"`
}

type WebhookProviderAuthConfig struct {
	Type     string `json:"type"`
	Header   string `json:"header"`
	Prefix   string `json:"prefix"`
	Provider string `json:"provider,omitempty"`
}

type WebhookProviderSetupConfig struct {
	Method             string                     `json:"method"`
	URLTemplate        string                     `json:"url_template"`
	Headers            map[string]string          `json:"headers,omitempty"`
	Auth               *WebhookProviderAuthConfig `json:"auth,omitempty"`
	BodyTemplate       json.RawMessage            `json:"body_template,omitempty"`
	BodyEncoding       string                     `json:"body_encoding,omitempty"`
	RepeatFor          string                     `json:"repeat_for,omitempty"`
	ResponseIDJSONPath string                     `json:"response_id_json_path"`
}

type WebhookProviderConfig struct {
	Name          string                        `json:"name"`
	PayloadFormat string                        `json:"payload_format,omitempty"`
	OAuthProvider string                        `json:"oauth_provider,omitempty"`
	TopicTemplate string                        `json:"topic_template,omitempty"`
	Signature     *WebhookSignatureConfig       `json:"signature,omitempty"`
	EventHeader   string                        `json:"event_header"`
	EventJSONPath string                        `json:"event_json_path"`
	Mappings      []MappingConfig               `json:"mappings,omitempty"`
	Prepare       []WebhookProviderPrepareStep  `json:"prepare,omitempty"`
	Renewal       *WebhookProviderRenewalConfig `json:"renewal,omitempty"`
	Setup         *WebhookProviderSetupConfig   `json:"setup,omitempty"`
	Teardown      *WebhookProviderSetupConfig   `json:"teardown,omitempty"`
}

type WebhookPrepareCondition struct {
	JSONPath string   `json:"json_path"`
	Equals   string   `json:"equals,omitempty"`
	In       []string `json:"in,omitempty"`
	Exists   *bool    `json:"exists,omitempty"`
}

type WebhookProviderPrepareStep struct {
	When         *WebhookPrepareCondition           `json:"when,omitempty"`
	Fetch        *WebhookProviderFetchConfig        `json:"fetch,omitempty"`
	TemplateList *WebhookProviderTemplateListConfig `json:"template_list,omitempty"`
	Extract      *WebhookProviderExtractConfig      `json:"extract,omitempty"`
}

type WebhookProviderFetchConfig struct {
	Method           string                           `json:"method"`
	URLTemplate      string                           `json:"url_template"`
	Headers          map[string]string                `json:"headers,omitempty"`
	Auth             *WebhookProviderAuthConfig       `json:"auth,omitempty"`
	BodyTemplate     json.RawMessage                  `json:"body_template,omitempty"`
	BodyEncoding     string                           `json:"body_encoding,omitempty"`
	ResponseJSONPath string                           `json:"response_json_path,omitempty"`
	ItemJSONPath     string                           `json:"item_json_path,omitempty"`
	StorePath        string                           `json:"store_path"`
	Pagination       *WebhookProviderPaginationConfig `json:"pagination,omitempty"`
}

type WebhookProviderPaginationConfig struct {
	RequestParam     string `json:"request_param"`
	ResponseJSONPath string `json:"response_json_path"`
}

type WebhookProviderTemplateListConfig struct {
	RepeatFor string `json:"repeat_for"`
	Template  string `json:"template"`
	StorePath string `json:"store_path"`
	Unique    bool   `json:"unique,omitempty"`
}

type WebhookProviderExtractConfig struct {
	SourceJSONPath string `json:"source_json_path"`
	Regex          string `json:"regex"`
	Group          int    `json:"group,omitempty"`
	StorePath      string `json:"store_path"`
	Optional       bool   `json:"optional,omitempty"`
}

type WebhookProviderRenewalConfig struct {
	AfterSeconds int `json:"after_seconds"`
}

// LoadProviderConfigs loads all *.json provider configs from the given directory.
// Returns a map keyed by the provider name (e.g. "google", "discord").
func LoadProviderConfigs(dir string) (map[string]ProviderConfig, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read providers dir: %w", err)
	}

	providers := make(map[string]ProviderConfig)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		path := filepath.Join(dir, entry.Name())

		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read provider file %s: %w", path, err)
		}

		var cfg ProviderConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("unmarshal provider file %s: %w", path, err)
		}

		if cfg.Name == "" {
			return nil, fmt.Errorf("provider file %s: missing name", path)
		}

		for _, f := range cfg.Mappings {
			if f.FieldKey == "" {
				return nil, fmt.Errorf("webhook provider %s: mapping missing field_key", cfg.Name)
			}
			if f.JSONPath == "" {
				return nil, fmt.Errorf("webhook provider %s: mapping %s missing json_path", cfg.Name, f.FieldKey)
			}
			switch f.Type {
			case "string", "number", "boolean", "json":
			default:
				return nil, fmt.Errorf("provider %s: field %s has invalid type %q", cfg.Name, f.FieldKey, f.Type)
			}
		}

		providers[cfg.Name] = cfg
	}

	return providers, nil
}

type FieldConfig struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	Label         string `json:"label"`
	Required      bool   `json:"required"`
	DefaultValuer string `json:"default"`
}

type OutputField struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Label string `json:"label"`
}

type ActionConfig struct {
	Title        string        `json:"title"`
	Label        string        `json:"label"`
	Type         string        `json:"type"`
	Fields       []FieldConfig `json:"fields"`
	OutputFields []OutputField `json:"output_fields"`
}

type ReactionConfig struct {
	Title      string          `json:"title"`
	Label      string          `json:"label"`
	Url        string          `json:"url"`
	Fields     []FieldConfig   `json:"fields"`
	Method     string          `json:"method"`
	BodyType   string          `json:"bodyType"`
	BodyStruct json.RawMessage `json:"body_struct"`
}

type ServiceConfig struct {
	Provider  string           `json:"provider"`
	Name      string           `json:"name"`
	IconURL   string           `json:"icon_url"`
	Label     string           `json:"label"`
	Actions   []ActionConfig   `json:"actions"`
	Reactions []ReactionConfig `json:"reactions"`
}

func LoadServiceConfig(dir string) (map[string]ServiceConfig, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read services dir: %w", err)
	}

	services := make(map[string]ServiceConfig)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read service file %s: %w", path, err)
		}
		var cfg ServiceConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("unmarshal service file %s: %w", path, err)
		}
		services[cfg.Name] = cfg
	}
	return services, nil
}

// LoadWebhookProviderConfigs loads all *.json webhook provider configs from the given directory.
// Returns a map keyed by the provider name (e.g. "github", "generic").
func LoadWebhookProviderConfigs(dir string) (map[string]WebhookProviderConfig, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read webhook providers dir: %w", err)
	}

	providers := make(map[string]WebhookProviderConfig)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		path := filepath.Join(dir, entry.Name())

		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read webhook provider file %s: %w", path, err)
		}

		var cfg WebhookProviderConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("unmarshal webhook provider file %s: %w", path, err)
		}

		if cfg.Name == "" {
			return nil, fmt.Errorf("webhook provider file %s: missing name", path)
		}

		if err := validateWebhookSignatureConfig(cfg.Name, cfg.Signature); err != nil {
			return nil, err
		}
		if cfg.PayloadFormat != "" {
			switch strings.ToLower(cfg.PayloadFormat) {
			case "json", "xml":
			default:
				return nil, fmt.Errorf("webhook provider %s: unsupported payload_format %q", cfg.Name, cfg.PayloadFormat)
			}
		}
		if err := validateWebhookProviderPrepare(cfg.Name, cfg.Prepare); err != nil {
			return nil, err
		}
		if err := validateWebhookProviderRenewal(cfg.Name, cfg.Renewal); err != nil {
			return nil, err
		}

		if err := validateWebhookProviderAction(cfg.Name, "setup", cfg.Setup); err != nil {
			return nil, err
		}
		if err := validateWebhookProviderAction(cfg.Name, "teardown", cfg.Teardown); err != nil {
			return nil, err
		}

		for _, f := range cfg.Mappings {
			if f.FieldKey == "" {
				return nil, fmt.Errorf("webhook provider %s: mapping missing field_key", cfg.Name)
			}
			if f.JSONPath == "" {
				return nil, fmt.Errorf("webhook provider %s: mapping %s missing json_path", cfg.Name, f.FieldKey)
			}
			switch f.Type {
			case "string", "number", "boolean", "json":
			default:
				return nil, fmt.Errorf("webhook provider %s: field %s has invalid type %q", cfg.Name, f.FieldKey, f.Type)
			}
		}

		providers[cfg.Name] = cfg
	}

	return providers, nil
}

func validateWebhookProviderAction(providerName, label string, action *WebhookProviderSetupConfig) error {
	if action == nil {
		return nil
	}
	if action.Method == "" {
		return fmt.Errorf("webhook provider %s: %s method is required", providerName, label)
	}
	if action.URLTemplate == "" {
		return fmt.Errorf("webhook provider %s: %s url_template is required", providerName, label)
	}
	if action.Auth != nil {
		if action.Auth.Type == "" {
			return fmt.Errorf("webhook provider %s: %s auth type is required", providerName, label)
		}
		if action.Auth.Header == "" {
			return fmt.Errorf("webhook provider %s: %s auth header is required", providerName, label)
		}
		switch action.Auth.Type {
		case "oauth2":
		default:
			return fmt.Errorf("webhook provider %s: unsupported auth type %q", providerName, action.Auth.Type)
		}
		if action.Auth.Provider != "" && action.Auth.Type != "oauth2" {
			return fmt.Errorf("webhook provider %s: auth provider is only supported for oauth2", providerName)
		}
	}
	if action.BodyEncoding != "" {
		switch strings.ToLower(action.BodyEncoding) {
		case "json", "form", "x-www-form-urlencoded":
		default:
			return fmt.Errorf("webhook provider %s: unsupported body_encoding %q", providerName, action.BodyEncoding)
		}
	}
	if action.RepeatFor != "" && len(action.BodyTemplate) == 0 {
		return fmt.Errorf("webhook provider %s: %s repeat_for requires a body_template", providerName, label)
	}
	if len(action.BodyTemplate) > 0 {
		var payload any
		if err := json.Unmarshal(action.BodyTemplate, &payload); err != nil {
			return fmt.Errorf("webhook provider %s: invalid %s body_template: %w", providerName, label, err)
		}
	}
	return nil
}

func validateWebhookProviderPrepare(providerName string, steps []WebhookProviderPrepareStep) error {
	for idx, step := range steps {
		hasFetch := step.Fetch != nil
		hasTemplate := step.TemplateList != nil
		hasExtract := step.Extract != nil
		stepCount := 0
		if hasFetch {
			stepCount++
		}
		if hasTemplate {
			stepCount++
		}
		if hasExtract {
			stepCount++
		}
		if stepCount != 1 {
			return fmt.Errorf("webhook provider %s: prepare[%d] must define exactly one of fetch, template_list, or extract", providerName, idx)
		}
		if step.When != nil && strings.TrimSpace(step.When.JSONPath) == "" {
			return fmt.Errorf("webhook provider %s: prepare[%d] when.json_path is required", providerName, idx)
		}
		if step.Fetch != nil {
			if err := validateWebhookProviderFetch(providerName, idx, step.Fetch); err != nil {
				return err
			}
		}
		if step.TemplateList != nil {
			if err := validateWebhookProviderTemplateList(providerName, idx, step.TemplateList); err != nil {
				return err
			}
		}
		if step.Extract != nil {
			if err := validateWebhookProviderExtract(providerName, idx, step.Extract); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateWebhookProviderFetch(providerName string, idx int, fetch *WebhookProviderFetchConfig) error {
	if fetch.Method == "" {
		return fmt.Errorf("webhook provider %s: prepare[%d] fetch method is required", providerName, idx)
	}
	if fetch.URLTemplate == "" {
		return fmt.Errorf("webhook provider %s: prepare[%d] fetch url_template is required", providerName, idx)
	}
	if strings.TrimSpace(fetch.StorePath) == "" {
		return fmt.Errorf("webhook provider %s: prepare[%d] fetch store_path is required", providerName, idx)
	}
	if fetch.Auth != nil {
		if fetch.Auth.Type == "" {
			return fmt.Errorf("webhook provider %s: prepare[%d] fetch auth type is required", providerName, idx)
		}
		if fetch.Auth.Header == "" {
			return fmt.Errorf("webhook provider %s: prepare[%d] fetch auth header is required", providerName, idx)
		}
		switch fetch.Auth.Type {
		case "oauth2":
		default:
			return fmt.Errorf("webhook provider %s: prepare[%d] fetch unsupported auth type %q", providerName, idx, fetch.Auth.Type)
		}
		if fetch.Auth.Provider != "" && fetch.Auth.Type != "oauth2" {
			return fmt.Errorf("webhook provider %s: prepare[%d] fetch auth provider is only supported for oauth2", providerName, idx)
		}
	}
	if fetch.BodyEncoding != "" {
		switch strings.ToLower(fetch.BodyEncoding) {
		case "json", "form", "x-www-form-urlencoded":
		default:
			return fmt.Errorf("webhook provider %s: prepare[%d] fetch unsupported body_encoding %q", providerName, idx, fetch.BodyEncoding)
		}
	}
	if len(fetch.BodyTemplate) > 0 {
		var payload any
		if err := json.Unmarshal(fetch.BodyTemplate, &payload); err != nil {
			return fmt.Errorf("webhook provider %s: prepare[%d] invalid body_template: %w", providerName, idx, err)
		}
	}
	if fetch.Pagination != nil {
		if strings.TrimSpace(fetch.Pagination.RequestParam) == "" {
			return fmt.Errorf("webhook provider %s: prepare[%d] pagination request_param is required", providerName, idx)
		}
		if strings.TrimSpace(fetch.Pagination.ResponseJSONPath) == "" {
			return fmt.Errorf("webhook provider %s: prepare[%d] pagination response_json_path is required", providerName, idx)
		}
	}
	return nil
}

func validateWebhookProviderTemplateList(providerName string, idx int, step *WebhookProviderTemplateListConfig) error {
	if strings.TrimSpace(step.RepeatFor) == "" {
		return fmt.Errorf("webhook provider %s: prepare[%d] template_list repeat_for is required", providerName, idx)
	}
	if strings.TrimSpace(step.Template) == "" {
		return fmt.Errorf("webhook provider %s: prepare[%d] template_list template is required", providerName, idx)
	}
	if strings.TrimSpace(step.StorePath) == "" {
		return fmt.Errorf("webhook provider %s: prepare[%d] template_list store_path is required", providerName, idx)
	}
	return nil
}

func validateWebhookProviderExtract(providerName string, idx int, step *WebhookProviderExtractConfig) error {
	if strings.TrimSpace(step.SourceJSONPath) == "" {
		return fmt.Errorf("webhook provider %s: prepare[%d] extract source_json_path is required", providerName, idx)
	}
	if strings.TrimSpace(step.Regex) == "" {
		return fmt.Errorf("webhook provider %s: prepare[%d] extract regex is required", providerName, idx)
	}
	if strings.TrimSpace(step.StorePath) == "" {
		return fmt.Errorf("webhook provider %s: prepare[%d] extract store_path is required", providerName, idx)
	}
	if step.Group < 0 {
		return fmt.Errorf("webhook provider %s: prepare[%d] extract group must be >= 0", providerName, idx)
	}
	return nil
}

func validateWebhookProviderRenewal(providerName string, renewal *WebhookProviderRenewalConfig) error {
	if renewal == nil {
		return nil
	}
	if renewal.AfterSeconds <= 0 {
		return fmt.Errorf("webhook provider %s: renewal.after_seconds must be > 0", providerName)
	}
	return nil
}

func validateWebhookSignatureConfig(providerName string, sig *WebhookSignatureConfig) error {
	if sig == nil {
		return nil
	}

	if sig.Type == "" {
		return fmt.Errorf("webhook provider %s: signature type is required", providerName)
	}

	signatureType, algorithm := normalizeSignatureType(sig)
	if signatureType == "token" {
		signatureType = "header"
	}

	switch signatureType {
	case "hmac":
		if sig.Header == "" {
			return fmt.Errorf("webhook provider %s: signature header is required", providerName)
		}
		if sig.SecretJSONPath == "" {
			return fmt.Errorf("webhook provider %s: signature secret_json_path is required", providerName)
		}
		if algorithm == "" {
			algorithm = "sha256"
		}
		switch algorithm {
		case "sha1", "sha256", "sha512":
		default:
			return fmt.Errorf("webhook provider %s: unsupported signature algorithm %q", providerName, algorithm)
		}
		if sig.Encoding != "" {
			switch strings.ToLower(sig.Encoding) {
			case "hex", "base64":
			default:
				return fmt.Errorf("webhook provider %s: unsupported signature encoding %q", providerName, sig.Encoding)
			}
		}
		if sig.TimestampToleranceSeconds < 0 {
			return fmt.Errorf("webhook provider %s: timestamp_tolerance_seconds must be >= 0", providerName)
		}
	case "header":
		if sig.Header == "" {
			return fmt.Errorf("webhook provider %s: signature header is required", providerName)
		}
		if sig.SecretJSONPath == "" {
			return fmt.Errorf("webhook provider %s: signature secret_json_path is required", providerName)
		}
	default:
		return fmt.Errorf("webhook provider %s: unsupported signature type %q", providerName, sig.Type)
	}

	return nil
}

func normalizeSignatureType(sig *WebhookSignatureConfig) (string, string) {
	signatureType := strings.ToLower(sig.Type)
	algorithm := strings.ToLower(sig.Algorithm)

	if strings.HasPrefix(signatureType, "hmac-") {
		algorithm = strings.TrimPrefix(signatureType, "hmac-")
		signatureType = "hmac"
	}

	return signatureType, algorithm
}
