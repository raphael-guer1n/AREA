package service

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/raphael-guer1n/AREA/WebhookService/internal/config"
	"github.com/raphael-guer1n/AREA/WebhookService/internal/utils"
)

func (s *SubscriptionService) applyPrepareSteps(userID int, providerConfig *config.WebhookProviderConfig, cfg map[string]any) (map[string]any, error) {
	if providerConfig == nil || len(providerConfig.Prepare) == 0 {
		return cfg, nil
	}

	for _, step := range providerConfig.Prepare {
		if step.When != nil {
			ok, err := matchesPrepareCondition(step.When, cfg)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}
		}

		switch {
		case step.Fetch != nil:
			if err := s.applyFetchStep(userID, providerConfig.Name, step.Fetch, cfg); err != nil {
				return nil, err
			}
		case step.TemplateList != nil:
			if err := applyTemplateListStep(step.TemplateList, cfg); err != nil {
				return nil, err
			}
		case step.Extract != nil:
			if err := applyExtractStep(step.Extract, cfg); err != nil {
				return nil, err
			}
		}
	}

	return cfg, nil
}

func matchesPrepareCondition(cond *config.WebhookPrepareCondition, cfg map[string]any) (bool, error) {
	if cond == nil {
		return true, nil
	}
	path := strings.TrimSpace(cond.JSONPath)
	path = strings.TrimPrefix(path, "config.")
	if path == "" {
		return false, fmt.Errorf("%w: prepare condition missing json_path", ErrInvalidConfig)
	}

	value, ok := utils.ExtractJSONPath(cfg, path)

	if cond.Exists != nil {
		if ok != *cond.Exists {
			return false, nil
		}
	}

	if cond.Equals != "" {
		if !ok {
			return false, nil
		}
		if !strings.EqualFold(fmt.Sprint(value), cond.Equals) {
			return false, nil
		}
	}

	if len(cond.In) > 0 {
		if !ok {
			return false, nil
		}
		matched := false
		current := strings.ToLower(fmt.Sprint(value))
		for _, candidate := range cond.In {
			if strings.EqualFold(current, strings.ToLower(candidate)) {
				matched = true
				break
			}
		}
		if !matched {
			return false, nil
		}
	}

	return true, nil
}

func (s *SubscriptionService) applyFetchStep(userID int, provider string, fetch *config.WebhookProviderFetchConfig, cfg map[string]any) error {
	if fetch == nil {
		return nil
	}

	action := &config.WebhookProviderSetupConfig{
		Method:       fetch.Method,
		URLTemplate:  fetch.URLTemplate,
		Headers:      fetch.Headers,
		Auth:         fetch.Auth,
		BodyTemplate: fetch.BodyTemplate,
		BodyEncoding: fetch.BodyEncoding,
	}

	ctx := utils.TemplateContext{
		Config:   cfg,
		UserID:   userID,
		Provider: provider,
	}

	var items []any
	pageToken := ""

	for {
		queryOverrides := map[string]string{}
		if fetch.Pagination != nil && pageToken != "" {
			queryOverrides[fetch.Pagination.RequestParam] = pageToken
		}

		responseBody, err := s.webhookSetupSvc.executeActionOnce(action, provider, userID, ctx, "prepare", queryOverrides)
		if err != nil {
			return err
		}

		var payload any
		if err := json.Unmarshal(responseBody, &payload); err != nil {
			return fmt.Errorf("%w: invalid prepare response", ErrInvalidConfig)
		}

		value := payload
		if fetch.ResponseJSONPath != "" {
			val, ok := utils.ExtractJSONPath(payload, fetch.ResponseJSONPath)
			if !ok {
				return fmt.Errorf("%w: response json_path not found: %s", ErrInvalidConfig, fetch.ResponseJSONPath)
			}
			value = val
		}

		extracted, err := extractItems(value, fetch.ItemJSONPath)
		if err != nil {
			return err
		}
		items = append(items, extracted...)

		if fetch.Pagination == nil {
			break
		}
		nextToken, ok := utils.ExtractJSONPath(payload, fetch.Pagination.ResponseJSONPath)
		if !ok {
			break
		}
		pageToken = strings.TrimSpace(fmt.Sprint(nextToken))
		if pageToken == "" {
			break
		}
	}

	if fetch.StorePath == "" {
		return fmt.Errorf("%w: prepare fetch store_path is required", ErrInvalidConfig)
	}
	return setConfigValue(cfg, fetch.StorePath, items)
}

func extractItems(value any, itemPath string) ([]any, error) {
	if itemPath == "" {
		switch v := value.(type) {
		case []any:
			return v, nil
		case []string:
			out := make([]any, 0, len(v))
			for _, item := range v {
				out = append(out, item)
			}
			return out, nil
		default:
			return []any{value}, nil
		}
	}

	list, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("%w: expected array response", ErrInvalidConfig)
	}

	out := make([]any, 0, len(list))
	for _, item := range list {
		val, ok := utils.ExtractJSONPath(item, itemPath)
		if !ok {
			continue
		}
		out = append(out, val)
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("%w: no items extracted", ErrInvalidConfig)
	}

	return out, nil
}

func applyTemplateListStep(step *config.WebhookProviderTemplateListConfig, cfg map[string]any) error {
	if step == nil {
		return nil
	}
	if step.RepeatFor == "" || step.Template == "" || step.StorePath == "" {
		return fmt.Errorf("%w: prepare template_list requires repeat_for, template and store_path", ErrInvalidConfig)
	}

	items, err := resolveRepeatItems(cfg, step.RepeatFor)
	if err != nil {
		return err
	}

	results := make([]string, 0, len(items))
	seen := map[string]struct{}{}
	for idx, item := range items {
		ctx := utils.TemplateContext{
			Config:      cfg,
			Item:        item,
			RepeatIndex: idx,
		}
		rendered, err := utils.RenderTemplateString(step.Template, ctx)
		if err != nil {
			return err
		}
		value := strings.TrimSpace(fmt.Sprint(rendered))
		if value == "" {
			continue
		}
		if step.Unique {
			if _, ok := seen[value]; ok {
				continue
			}
			seen[value] = struct{}{}
		}
		results = append(results, value)
	}

	if len(results) == 0 {
		return fmt.Errorf("%w: template_list produced no values", ErrInvalidConfig)
	}

	return setConfigValue(cfg, step.StorePath, results)
}

func applyExtractStep(step *config.WebhookProviderExtractConfig, cfg map[string]any) error {
	if step == nil {
		return nil
	}
	if strings.TrimSpace(step.SourceJSONPath) == "" || strings.TrimSpace(step.Regex) == "" || strings.TrimSpace(step.StorePath) == "" {
		return fmt.Errorf("%w: prepare extract requires source_json_path, regex and store_path", ErrInvalidConfig)
	}

	path := strings.TrimPrefix(strings.TrimSpace(step.SourceJSONPath), "config.")
	value, ok := utils.ExtractJSONPath(cfg, path)
	if !ok {
		return nil
	}

	input := fmt.Sprint(value)
	if strings.TrimSpace(input) == "" {
		return nil
	}

	re, err := regexp.Compile(step.Regex)
	if err != nil {
		return fmt.Errorf("%w: invalid extract regex", ErrInvalidConfig)
	}

	matches := re.FindStringSubmatch(input)
	if matches == nil {
		if step.Optional {
			return nil
		}
		return fmt.Errorf("%w: extract regex did not match", ErrInvalidConfig)
	}

	group := step.Group
	if group == 0 {
		group = 1
	}
	if group < 0 || group >= len(matches) {
		return fmt.Errorf("%w: extract group out of range", ErrInvalidConfig)
	}

	return setConfigValue(cfg, step.StorePath, matches[group])
}

func resolveRepeatItems(cfg map[string]any, path string) ([]any, error) {
	trimmed := strings.TrimSpace(path)
	trimmed = strings.TrimPrefix(trimmed, "config.")
	if trimmed == "" {
		return nil, fmt.Errorf("%w: repeat_for path is empty", ErrInvalidConfig)
	}

	value, ok := utils.ExtractJSONPath(cfg, trimmed)
	if !ok {
		return nil, fmt.Errorf("%w: repeat_for path not found: %s", ErrInvalidConfig, path)
	}

	switch v := value.(type) {
	case []any:
		return v, nil
	case []string:
		out := make([]any, 0, len(v))
		for _, item := range v {
			out = append(out, item)
		}
		return out, nil
	default:
		return []any{v}, nil
	}
}

func setConfigValue(cfg map[string]any, path string, value any) error {
	trimmed := strings.TrimSpace(path)
	trimmed = strings.TrimPrefix(trimmed, "config.")
	if trimmed == "" {
		return fmt.Errorf("%w: invalid store_path", ErrInvalidConfig)
	}

	segments := strings.Split(trimmed, ".")
	current := cfg
	for i, segment := range segments {
		if segment == "" {
			return fmt.Errorf("%w: invalid store_path", ErrInvalidConfig)
		}
		if i == len(segments)-1 {
			current[segment] = value
			return nil
		}
		next, ok := current[segment].(map[string]any)
		if !ok {
			next = map[string]any{}
			current[segment] = next
		}
		current = next
	}

	return nil
}
