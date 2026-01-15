package service

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/raphael-guer1n/AREA/PollingService/internal/config"
	"github.com/raphael-guer1n/AREA/PollingService/internal/domain"
	"github.com/raphael-guer1n/AREA/PollingService/internal/utils"
)

type PollingWorker struct {
	repo           domain.SubscriptionRepository
	providerSvc    *ProviderConfigService
	requestSvc     *RequestService
	areaTriggerSvc *AreaTriggerService
	tick           time.Duration
}

func NewPollingWorker(repo domain.SubscriptionRepository, providerSvc *ProviderConfigService, requestSvc *RequestService, areaTriggerSvc *AreaTriggerService, tickSeconds int) *PollingWorker {
	if tickSeconds <= 0 {
		tickSeconds = 60
	}
	return &PollingWorker{
		repo:           repo,
		providerSvc:    providerSvc,
		requestSvc:     requestSvc,
		areaTriggerSvc: areaTriggerSvc,
		tick:           time.Duration(tickSeconds) * time.Second,
	}
}

func (w *PollingWorker) Start() {
	w.pollDue()
	ticker := time.NewTicker(w.tick)
	defer ticker.Stop()

	for range ticker.C {
		w.pollDue()
	}
}

func (w *PollingWorker) pollDue() {
	now := time.Now().UTC()
	subs, err := w.repo.ListDue(now)
	if err != nil {
		log.Printf("polling: failed to list subscriptions: %v", err)
		return
	}
	for _, sub := range subs {
		if err := w.processSubscription(&sub); err != nil {
			log.Printf("polling: subscription action_id=%d error=%v", sub.ActionID, err)
		}
	}
}

func (w *PollingWorker) processSubscription(sub *domain.Subscription) error {
	if sub == nil || !sub.Active {
		return nil
	}
	if w.providerSvc == nil || w.requestSvc == nil {
		return errors.New("polling services not configured")
	}

	providerConfig, err := w.providerSvc.GetProviderConfig(sub.Service)
	if err != nil {
		return err
	}
	if providerConfig.IntervalSeconds <= 0 || strings.TrimSpace(providerConfig.Request.Method) == "" || strings.TrimSpace(providerConfig.Request.URLTemplate) == "" {
		return w.finishWithError(sub, providerConfig, fmt.Errorf("invalid provider config"))
	}

	var cfgPayload any = map[string]any{}
	if len(sub.Config) > 0 {
		if err := json.Unmarshal(sub.Config, &cfgPayload); err != nil {
			return fmt.Errorf("invalid subscription config")
		}
	}
	cfgMap, ok := cfgPayload.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid subscription config")
	}

	ctx := utils.TemplateContext{
		Config:   cfgMap,
		UserID:   sub.UserID,
		Provider: sub.Provider,
		Env:      utils.EnvMap(),
	}

	payloadBody, err := w.requestSvc.ExecuteRequest(providerConfig.Request, sub.Provider, sub.UserID, ctx, nil)
	if err != nil {
		return w.finishWithError(sub, providerConfig, err)
	}

	payload, err := parsePayload(payloadBody, providerConfig.PayloadFormat)
	if err != nil {
		return w.finishWithError(sub, providerConfig, err)
	}

	itemsPath, err := renderTemplatePath(providerConfig.ItemsPath, ctx)
	if err != nil {
		return w.finishWithError(sub, providerConfig, err)
	}
	itemIDPath, err := renderTemplatePath(providerConfig.ItemIDPath, ctx)
	if err != nil {
		return w.finishWithError(sub, providerConfig, err)
	}

	items, err := extractItemsFromPayload(payload, itemsPath)
	if err != nil {
		return w.finishWithError(sub, providerConfig, err)
	}

	filtered := filterItems(items, providerConfig.Filters)
	newItems, newLastID, err := selectNewItemsWithConfig(filtered, sub.LastItemID, itemIDPath, providerConfig.ChangeDetection, ctx)
	if err != nil {
		return w.finishWithError(sub, providerConfig, err)
	}

	if len(newItems) > 0 && w.areaTriggerSvc != nil {
		for i := len(newItems) - 1; i >= 0; i-- {
			item := newItems[i]
			mapped, err := buildMappings(item, providerConfig.Mappings, ctx)
			if err != nil {
				log.Printf("polling: mapping error action_id=%d provider=%s err=%v", sub.ActionID, sub.Service, err)
				continue
			}
			outputFields := buildOutputFields(providerConfig.Mappings, mapped)
			if err := w.areaTriggerSvc.Trigger(sub.ActionID, outputFields); err != nil {
				log.Printf("polling: trigger failed action_id=%d provider=%s err=%v", sub.ActionID, sub.Service, err)
			}
		}
	}

	if newLastID == "" && len(filtered) > 0 {
		if id, err := resolveItemID(filtered[0], itemIDPath); err == nil {
			newLastID = id
		}
	}

	return w.finishWithSuccess(sub, providerConfig, newLastID)
}

func (w *PollingWorker) finishWithSuccess(sub *domain.Subscription, providerConfig *config.PollingProviderConfig, lastItemID string) error {
	now := time.Now().UTC()
	next := computeNextRunAt(sub, providerConfig.IntervalSeconds, now)
	return w.repo.UpdatePollingState(sub.ActionID, lastItemID, next, "", now)
}

func (w *PollingWorker) finishWithError(sub *domain.Subscription, providerConfig *config.PollingProviderConfig, err error) error {
	now := time.Now().UTC()
	next := computeNextRunAt(sub, providerConfig.IntervalSeconds, now)
	updateErr := w.repo.UpdatePollingState(sub.ActionID, sub.LastItemID, next, err.Error(), now)
	if updateErr != nil {
		log.Printf("polling: failed to update error state action_id=%d err=%v", sub.ActionID, updateErr)
	}
	return err
}

func computeNextRunAt(sub *domain.Subscription, intervalSeconds int, now time.Time) time.Time {
	interval := time.Duration(intervalSeconds) * time.Second
	if interval <= 0 {
		interval = time.Minute
	}

	base := now
	if sub != nil && sub.NextRunAt != nil && !sub.NextRunAt.IsZero() {
		base = sub.NextRunAt.UTC()
	}

	next := base.Add(interval)
	for next.Before(now) {
		next = next.Add(interval)
	}
	return next
}

func parsePayload(body []byte, format string) (any, error) {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "", "json":
		var payload any
		if len(body) == 0 {
			return map[string]any{}, nil
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			return nil, fmt.Errorf("invalid json payload")
		}
		return payload, nil
	case "xml":
		payload, err := utils.ParseXMLToMap(body)
		if err != nil {
			return nil, fmt.Errorf("invalid xml payload")
		}
		return payload, nil
	default:
		return nil, fmt.Errorf("unsupported payload format")
	}
}

func extractItemsFromPayload(payload any, path string) ([]any, error) {
	value := payload
	if strings.TrimSpace(path) != "" {
		val, ok := utils.ExtractJSONPath(payload, path)
		if !ok {
			return nil, fmt.Errorf("items path not found")
		}
		value = val
	}

	switch v := value.(type) {
	case []any:
		return v, nil
	case []map[string]any:
		items := make([]any, 0, len(v))
		for _, item := range v {
			items = append(items, item)
		}
		return items, nil
	case nil:
		return []any{}, nil
	default:
		return []any{v}, nil
	}
}

func filterItems(items []any, filters *config.PollingFilterConfig) []any {
	if filters == nil || len(filters.Rules) == 0 {
		return items
	}
	mode := strings.ToLower(strings.TrimSpace(filters.Mode))
	if mode == "" {
		mode = "all"
	}

	out := make([]any, 0, len(items))
	for _, item := range items {
		if matchesFilters(item, filters, mode) {
			out = append(out, item)
		}
	}
	return out
}

func matchesFilters(item any, filters *config.PollingFilterConfig, mode string) bool {
	if filters == nil || len(filters.Rules) == 0 {
		return true
	}
	switch mode {
	case "any":
		for _, rule := range filters.Rules {
			if matchesRule(item, rule) {
				return true
			}
		}
		return false
	default:
		for _, rule := range filters.Rules {
			if !matchesRule(item, rule) {
				return false
			}
		}
		return true
	}
}

func matchesRule(item any, rule config.PollingFilterRule) bool {
	path := strings.TrimSpace(rule.JSONPath)
	if path == "" {
		return false
	}
	value, ok := utils.ExtractJSONPath(item, path)
	operator := strings.ToLower(strings.TrimSpace(rule.Operator))
	if operator == "" {
		operator = "equals"
	}

	switch operator {
	case "exists":
		return ok
	}
	if !ok {
		return false
	}

	switch operator {
	case "equals":
		return compareString(value, rule.Value, rule.CaseInsensitive)
	case "contains":
		left := normalizeString(value, rule.CaseInsensitive)
		right := normalizeString(rule.Value, rule.CaseInsensitive)
		return left != "" && strings.Contains(left, right)
	case "in":
		candidates := rule.Values
		if len(candidates) == 0 && rule.Value != nil {
			candidates = []any{rule.Value}
		}
		for _, candidate := range candidates {
			if compareString(value, candidate, rule.CaseInsensitive) {
				return true
			}
		}
		return false
	case "regex":
		pattern := fmt.Sprint(rule.Value)
		if pattern == "" {
			return false
		}
		re, err := compileRegex(pattern, rule.CaseInsensitive)
		if err != nil {
			return false
		}
		return re.MatchString(fmt.Sprint(value))
	case "gt", "gte", "lt", "lte":
		left, okLeft := toFloat(value)
		right, okRight := toFloat(rule.Value)
		if !okLeft || !okRight {
			return false
		}
		switch operator {
		case "gt":
			return left > right
		case "gte":
			return left >= right
		case "lt":
			return left < right
		case "lte":
			return left <= right
		}
	}

	return false
}

func normalizeString(value any, caseInsensitive bool) string {
	str := fmt.Sprint(value)
	if caseInsensitive {
		return strings.ToLower(str)
	}
	return str
}

func compareString(left any, right any, caseInsensitive bool) bool {
	l := normalizeString(left, caseInsensitive)
	r := normalizeString(right, caseInsensitive)
	if caseInsensitive {
		return l == r
	}
	return l == r
}

func compileRegex(pattern string, caseInsensitive bool) (*regexp.Regexp, error) {
	if caseInsensitive {
		pattern = "(?i)" + pattern
	}
	return regexp.Compile(pattern)
}

func toFloat(value any) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case json.Number:
		f, err := v.Float64()
		if err != nil {
			return 0, false
		}
		return f, true
	case string:
		trimmed := strings.TrimSpace(v)
		if trimmed == "" {
			return 0, false
		}
		f, err := strconv.ParseFloat(trimmed, 64)
		if err != nil {
			return 0, false
		}
		return f, true
	default:
		return 0, false
	}
}

func selectNewItems(items []any, lastItemID string, itemIDPath string) ([]any, string) {
	if len(items) == 0 {
		return nil, lastItemID
	}

	newItems := make([]any, 0, len(items))
	for _, item := range items {
		id, err := resolveItemID(item, itemIDPath)
		if err != nil {
			continue
		}
		if lastItemID != "" && id == lastItemID {
			break
		}
		newItems = append(newItems, item)
	}

	if len(newItems) == 0 {
		return nil, lastItemID
	}

	newLastID, err := resolveItemID(newItems[0], itemIDPath)
	if err != nil {
		return newItems, lastItemID
	}
	return newItems, newLastID
}

func selectNewItemsWithConfig(items []any, lastItemID string, itemIDPath string, changeCfg *config.PollingChangeDetectionConfig, ctx utils.TemplateContext) ([]any, string, error) {
	if changeCfg == nil {
		newItems, newLastID := selectNewItems(items, lastItemID, itemIDPath)
		return newItems, newLastID, nil
	}
	if len(items) == 0 {
		return nil, lastItemID, nil
	}

	item := items[0]
	valuePath := changeCfg.ValueJSONPath
	if strings.TrimSpace(valuePath) == "" {
		valuePath = itemIDPath
	}
	valuePath, err := renderTemplatePath(valuePath, ctx)
	if err != nil {
		return nil, lastItemID, err
	}
	if strings.TrimSpace(valuePath) == "" {
		return nil, lastItemID, fmt.Errorf("change detection value path is empty")
	}

	value, ok := utils.ExtractJSONPath(item, valuePath)
	if !ok {
		return nil, lastItemID, fmt.Errorf("change detection path not found")
	}
	current, ok := toFloat(value)
	if !ok {
		return nil, lastItemID, fmt.Errorf("change detection value is not numeric")
	}

	if lastItemID == "" {
		return []any{item}, formatFloat(current), nil
	}

	prev, ok := toFloat(lastItemID)
	if !ok {
		return nil, lastItemID, fmt.Errorf("change detection previous value is not numeric")
	}

	minPercent, hasPercent, err := resolveThreshold(changeCfg.MinPercent, ctx)
	if err != nil {
		return nil, lastItemID, err
	}
	minDelta, hasDelta, err := resolveThreshold(changeCfg.MinDelta, ctx)
	if err != nil {
		return nil, lastItemID, err
	}

	diff := math.Abs(current - prev)
	shouldTrigger := false
	switch {
	case hasPercent && minPercent > 0:
		if prev == 0 {
			shouldTrigger = current != 0
		} else {
			shouldTrigger = (diff/math.Abs(prev))*100 >= minPercent
		}
	case hasDelta && minDelta > 0:
		shouldTrigger = diff >= minDelta
	default:
		shouldTrigger = diff > 0
	}

	if !shouldTrigger {
		return nil, lastItemID, nil
	}

	return []any{item}, formatFloat(current), nil
}

func resolveItemID(item any, itemIDPath string) (string, error) {
	if strings.TrimSpace(itemIDPath) != "" {
		val, ok := utils.ExtractJSONPath(item, itemIDPath)
		if !ok {
			return "", fmt.Errorf("item id path not found")
		}
		id := strings.TrimSpace(fmt.Sprint(val))
		if id == "" {
			return "", fmt.Errorf("item id is empty")
		}
		return id, nil
	}

	payload, err := json.Marshal(item)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return fmt.Sprintf("%x", sum[:]), nil
}

func buildMappings(payload any, mappings []config.MappingConfig, ctx utils.TemplateContext) (map[string]any, error) {
	mapped := make(map[string]any, len(mappings))
	for _, mapping := range mappings {
		jsonPath, err := renderTemplatePath(mapping.JSONPath, ctx)
		if err != nil {
			return nil, err
		}
		value, ok := utils.ExtractJSONPath(payload, jsonPath)
		if !ok {
			if mapping.Optional {
				continue
			}
			return nil, fmt.Errorf("missing json path %s", jsonPath)
		}
		coerced, err := coerceValue(value, mapping.Type)
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", mapping.FieldKey, err)
		}
		mapped[mapping.FieldKey] = coerced
	}
	return mapped, nil
}

func buildOutputFields(mappings []config.MappingConfig, mapped map[string]any) []TriggerOutputField {
	if len(mapped) == 0 {
		return []TriggerOutputField{}
	}
	fields := make([]TriggerOutputField, 0, len(mapped))
	for _, mapping := range mappings {
		value, ok := mapped[mapping.FieldKey]
		if !ok {
			continue
		}
		fields = append(fields, TriggerOutputField{
			Name:  mapping.FieldKey,
			Value: stringifyOutputValue(value),
		})
	}
	return fields
}

func stringifyOutputValue(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case nil:
		return ""
	default:
		switch v.(type) {
		case map[string]any, []any:
			encoded, err := json.Marshal(v)
			if err == nil {
				return string(encoded)
			}
		}
		return fmt.Sprint(v)
	}
}

func coerceValue(value any, valueType string) (any, error) {
	switch valueType {
	case "string":
		if v, ok := value.(string); ok {
			return v, nil
		}
		return nil, fmt.Errorf("expected string")
	case "number":
		switch v := value.(type) {
		case float64:
			return v, nil
		case float32:
			return float64(v), nil
		case int:
			return float64(v), nil
		case int64:
			return float64(v), nil
		default:
			return nil, fmt.Errorf("expected number")
		}
	case "boolean":
		if v, ok := value.(bool); ok {
			return v, nil
		}
		return nil, fmt.Errorf("expected boolean")
	case "json":
		return value, nil
	default:
		return nil, fmt.Errorf("unsupported type")
	}
}

func renderTemplatePath(value string, ctx utils.TemplateContext) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", nil
	}
	rendered, err := utils.RenderTemplateString(value, ctx)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(fmt.Sprint(rendered)), nil
}

func resolveThreshold(raw any, ctx utils.TemplateContext) (float64, bool, error) {
	if raw == nil {
		return 0, false, nil
	}
	rendered, err := utils.RenderTemplateValue(raw, ctx)
	if err != nil {
		return 0, false, err
	}
	if rendered == nil {
		return 0, false, nil
	}
	value, ok := toFloat(rendered)
	if !ok {
		return 0, false, fmt.Errorf("threshold is not numeric")
	}
	return value, true, nil
}

func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
