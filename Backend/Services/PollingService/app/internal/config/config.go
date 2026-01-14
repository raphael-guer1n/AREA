package config

import (
	"encoding/json"
	"os"
	"strconv"
)

type Config struct {
	HTTPPort           string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPass             string
	DBName             string
	InternalSecret     string
	ServiceServiceURL  string
	AuthServiceURL     string
	AreaServiceURL     string
	PollingTickSeconds int
}

func Load() Config {
	return Config{
		HTTPPort:           getEnv("SERVER_PORT", "8087"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPass:             getEnv("DB_PASSWORD", "postgres"),
		DBName:             getEnv("DB_NAME", "polling_service_db"),
		InternalSecret:     getEnv("INTERNAL_SECRET", ""),
		ServiceServiceURL:  getEnv("SERVICE_SERVICE_URL", "http://gateway:8080/area_service_api"),
		AuthServiceURL:     getEnv("AUTH_SERVICE_URL", "http://gateway:8080/area_auth_api"),
		AreaServiceURL:     getEnv("AREA_SERVICE_URL", "http://gateway:8080/area_area_api"),
		PollingTickSeconds: getEnvInt("POLLING_TICK_SECONDS", 60),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
	}
	return def
}

type MappingConfig struct {
	FieldKey string `json:"field_key"`
	JSONPath string `json:"json_path"`
	Type     string `json:"type"`
	Optional bool   `json:"optional,omitempty"`
}

type PollingProviderAuthConfig struct {
	Type     string `json:"type"`
	Header   string `json:"header"`
	Prefix   string `json:"prefix"`
	Provider string `json:"provider,omitempty"`
}

type PollingProviderRequestConfig struct {
	Method       string                     `json:"method"`
	URLTemplate  string                     `json:"url_template"`
	QueryParams  map[string]string          `json:"query_params,omitempty"`
	Headers      map[string]string          `json:"headers,omitempty"`
	Auth         *PollingProviderAuthConfig `json:"auth,omitempty"`
	BodyTemplate json.RawMessage            `json:"body_template,omitempty"`
	BodyEncoding string                     `json:"body_encoding,omitempty"`
}

type PollingFilterRule struct {
	JSONPath        string `json:"json_path"`
	Operator        string `json:"operator,omitempty"`
	Value           any    `json:"value,omitempty"`
	Values          []any  `json:"values,omitempty"`
	CaseInsensitive bool   `json:"case_insensitive,omitempty"`
}

type PollingFilterConfig struct {
	Mode  string              `json:"mode,omitempty"`
	Rules []PollingFilterRule `json:"rules,omitempty"`
}

type PollingProviderConfig struct {
	Name            string                        `json:"name"`
	PayloadFormat   string                        `json:"payload_format,omitempty"`
	IntervalSeconds int                           `json:"interval_seconds"`
	Request         PollingProviderRequestConfig  `json:"request"`
	ItemsPath       string                        `json:"items_path,omitempty"`
	ItemIDPath      string                        `json:"item_id_path,omitempty"`
	ChangeDetection *PollingChangeDetectionConfig `json:"change_detection,omitempty"`
	Filters         *PollingFilterConfig          `json:"filters,omitempty"`
	Mappings        []MappingConfig               `json:"mappings,omitempty"`
	Prepare         []PollingProviderPrepareStep  `json:"prepare,omitempty"`
}

type PollingChangeDetectionConfig struct {
	ValueJSONPath string `json:"value_json_path,omitempty"`
	MinDelta      any    `json:"min_delta,omitempty"`
	MinPercent    any    `json:"min_percent,omitempty"`
}

type PollingPrepareCondition struct {
	JSONPath string   `json:"json_path"`
	Equals   string   `json:"equals,omitempty"`
	In       []string `json:"in,omitempty"`
	Exists   *bool    `json:"exists,omitempty"`
}

type PollingProviderPrepareStep struct {
	When         *PollingPrepareCondition           `json:"when,omitempty"`
	Fetch        *PollingProviderFetchConfig        `json:"fetch,omitempty"`
	TemplateList *PollingProviderTemplateListConfig `json:"template_list,omitempty"`
	Extract      *PollingProviderExtractConfig      `json:"extract,omitempty"`
	Generate     *PollingProviderGenerateConfig     `json:"generate,omitempty"`
}

type PollingProviderFetchConfig struct {
	Method           string                           `json:"method"`
	URLTemplate      string                           `json:"url_template"`
	Headers          map[string]string                `json:"headers,omitempty"`
	Auth             *PollingProviderAuthConfig       `json:"auth,omitempty"`
	BodyTemplate     json.RawMessage                  `json:"body_template,omitempty"`
	BodyEncoding     string                           `json:"body_encoding,omitempty"`
	ResponseJSONPath string                           `json:"response_json_path,omitempty"`
	ItemJSONPath     string                           `json:"item_json_path,omitempty"`
	StorePath        string                           `json:"store_path"`
	Pagination       *PollingProviderPaginationConfig `json:"pagination,omitempty"`
}

type PollingProviderPaginationConfig struct {
	RequestParam     string `json:"request_param"`
	ResponseJSONPath string `json:"response_json_path"`
}

type PollingProviderTemplateListConfig struct {
	RepeatFor string `json:"repeat_for"`
	Template  string `json:"template"`
	StorePath string `json:"store_path"`
	Unique    bool   `json:"unique,omitempty"`
}

type PollingProviderExtractConfig struct {
	SourceJSONPath string `json:"source_json_path"`
	Regex          string `json:"regex"`
	Group          int    `json:"group,omitempty"`
	StorePath      string `json:"store_path"`
	Optional       bool   `json:"optional,omitempty"`
}

type PollingProviderGenerateConfig struct {
	StorePath     string `json:"store_path"`
	Length        int    `json:"length,omitempty"`
	Encoding      string `json:"encoding,omitempty"`
	OnlyIfMissing bool   `json:"only_if_missing,omitempty"`
}
