package utils

import (
	"fmt"
	"regexp"
	"strings"
)

type TemplateContext struct {
	HookURL        string
	HookID         string
	Provider       string
	ProviderHookID string
	UserID         int
	AreaID         int
	Config         any
}

var placeholderRegexp = regexp.MustCompile(`\{\{\s*([^}]+)\s*\}\}`)

type MissingTemplateValueError struct {
	Key string
}

func (e MissingTemplateValueError) Error() string {
	return fmt.Sprintf("missing template value for %s", e.Key)
}

func RenderTemplateValue(value any, ctx TemplateContext) (any, error) {
	switch v := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(v))
		for key, item := range v {
			rendered, err := RenderTemplateValue(item, ctx)
			if err != nil {
				return nil, err
			}
			out[key] = rendered
		}
		return out, nil
	case []any:
		out := make([]any, len(v))
		for i, item := range v {
			rendered, err := RenderTemplateValue(item, ctx)
			if err != nil {
				return nil, err
			}
			out[i] = rendered
		}
		return out, nil
	case string:
		return RenderTemplateString(v, ctx)
	default:
		return value, nil
	}
}

func RenderTemplateString(input string, ctx TemplateContext) (any, error) {
	trimmed := strings.TrimSpace(input)
	if matches := placeholderRegexp.FindStringSubmatch(trimmed); len(matches) == 2 && matches[0] == trimmed {
		val, ok := resolvePlaceholder(matches[1], ctx)
		if !ok {
			return nil, MissingTemplateValueError{Key: matches[1]}
		}
		return val, nil
	}

	result := input
	for _, match := range placeholderRegexp.FindAllStringSubmatch(input, -1) {
		if len(match) != 2 {
			continue
		}
		val, ok := resolvePlaceholder(match[1], ctx)
		if !ok {
			return nil, MissingTemplateValueError{Key: match[1]}
		}
		switch val.(type) {
		case map[string]any, []any:
			return nil, fmt.Errorf("cannot embed non-scalar value for %s", match[1])
		}
		replacement := fmt.Sprint(val)
		result = strings.ReplaceAll(result, match[0], replacement)
	}

	return result, nil
}

func resolvePlaceholder(key string, ctx TemplateContext) (any, bool) {
	key = strings.TrimSpace(key)
	switch key {
	case "hook_url":
		return ctx.HookURL, ctx.HookURL != ""
	case "hook_id":
		return ctx.HookID, ctx.HookID != ""
	case "provider":
		return ctx.Provider, ctx.Provider != ""
	case "provider_hook_id":
		return ctx.ProviderHookID, ctx.ProviderHookID != ""
	case "user_id":
		return ctx.UserID, true
	case "area_id":
		return ctx.AreaID, true
	case "config":
		return ctx.Config, ctx.Config != nil
	}

	if strings.HasPrefix(key, "config.") {
		path := strings.TrimPrefix(key, "config.")
		val, ok := ExtractJSONPath(ctx.Config, path)
		return val, ok
	}

	return nil, false
}
