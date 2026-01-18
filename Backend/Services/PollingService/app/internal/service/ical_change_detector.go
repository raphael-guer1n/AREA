package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	icalStateVersion           = 1
	icalDeleteMissingThreshold = 2
)

type icalState struct {
	Version int                      `json:"version"`
	Items   map[string]icalStateItem `json:"items"`
}

type icalStateItem struct {
	ID           string `json:"id,omitempty"`
	UID          string `json:"uid"`
	RecurrenceID string `json:"recurrence_id,omitempty"`
	Summary      string `json:"summary,omitempty"`
	Location     string `json:"location,omitempty"`
	Description  string `json:"description,omitempty"`
	Start        string `json:"start,omitempty"`
	End          string `json:"end,omitempty"`
	StartRaw     string `json:"start_raw,omitempty"`
	EndRaw       string `json:"end_raw,omitempty"`
	AllDay       bool   `json:"all_day,omitempty"`
	Timezone     string `json:"timezone,omitempty"`
	Status       string `json:"status,omitempty"`
	URL          string `json:"url,omitempty"`
	Organizer    string `json:"organizer,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
	Sequence     int    `json:"sequence,omitempty"`
	Fingerprint  string `json:"fingerprint"`
	Missing      int    `json:"missing,omitempty"`
	LastSeen     string `json:"last_seen,omitempty"`
}

func selectICalChanges(items []any, lastItemID string) ([]any, string) {
	state, legacy := parseICalState(lastItemID)
	if state == nil {
		state = &icalState{
			Version: icalStateVersion,
			Items:   make(map[string]icalStateItem),
		}
	}
	if state.Items == nil {
		state.Items = make(map[string]icalStateItem)
	}

	now := time.Now().UTC()
	newItems := make([]any, 0, len(items))
	currentKeys := make(map[string]struct{}, len(items))

	for _, raw := range items {
		item, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		uid := getICalString(item, "uid")
		if uid == "" {
			continue
		}
		recurrenceID := getICalString(item, "recurrence_id")
		key := buildICalKey(uid, recurrenceID)
		currentKeys[key] = struct{}{}

		curr := buildICalStateItem(item, now)
		prev, hasPrev := state.Items[key]
		updateType := ""

		if !hasPrev {
			if isCancelledStatus(curr.Status) {
				updateType = "Cancelled"
			} else {
				updateType = "Created"
			}
		} else if !legacy {
			if statusChangedToCancelled(prev.Status, curr.Status) {
				updateType = "Cancelled"
			} else if isICalMoved(prev, curr) {
				updateType = "Moved"
			} else if prev.Fingerprint != curr.Fingerprint {
				updateType = "Updated"
			}
		}

		if updateType != "" && !legacy {
			item["update_type"] = updateType
			newItems = append(newItems, item)
		}

		curr.Missing = 0
		state.Items[key] = curr
	}

	if !legacy {
		for key, prev := range state.Items {
			if _, ok := currentKeys[key]; ok {
				continue
			}
			prev.Missing++
			if prev.Missing >= icalDeleteMissingThreshold {
				if payload := buildICalItemFromState(prev); payload != nil {
					payload["update_type"] = "deleted"
					newItems = append(newItems, payload)
				}
				delete(state.Items, key)
				continue
			}
			state.Items[key] = prev
		}
	}

	state.Version = icalStateVersion
	if len(state.Items) == 0 {
		return newItems, ""
	}
	payload, err := json.Marshal(state)
	if err != nil {
		return newItems, lastItemID
	}
	return newItems, string(payload)
}

func parseICalState(raw string) (*icalState, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	if strings.HasPrefix(raw, "{") {
		var state icalState
		if err := json.Unmarshal([]byte(raw), &state); err == nil && state.Version > 0 {
			if state.Items == nil {
				state.Items = make(map[string]icalStateItem)
			}
			return &state, false
		}
	}
	return nil, true
}

func buildICalStateItem(item map[string]any, now time.Time) icalStateItem {
	entry := icalStateItem{
		ID:           getICalString(item, "id"),
		UID:          getICalString(item, "uid"),
		RecurrenceID: getICalString(item, "recurrence_id"),
		Summary:      getICalString(item, "summary"),
		Location:     getICalString(item, "location"),
		Description:  getICalString(item, "description"),
		Start:        getICalString(item, "start"),
		End:          getICalString(item, "end"),
		StartRaw:     getICalString(item, "start_raw"),
		EndRaw:       getICalString(item, "end_raw"),
		AllDay:       getICalBool(item, "all_day"),
		Timezone:     getICalString(item, "timezone"),
		Status:       strings.ToUpper(getICalString(item, "status")),
		URL:          getICalString(item, "url"),
		Organizer:    getICalString(item, "organizer"),
		UpdatedAt:    getICalString(item, "updated_at"),
		Sequence:     getICalInt(item, "sequence"),
		LastSeen:     now.Format(time.RFC3339),
	}
	entry.Fingerprint = buildICalFingerprint(entry)
	return entry
}

func buildICalFingerprint(item icalStateItem) string {
	parts := []string{
		item.UID,
		item.RecurrenceID,
		item.Summary,
		item.Location,
		item.Description,
		item.Status,
		item.URL,
		item.Organizer,
		item.StartRaw,
		item.EndRaw,
		item.Start,
		item.End,
		item.Timezone,
		strconv.FormatBool(item.AllDay),
		strconv.Itoa(item.Sequence),
	}
	raw := strings.Join(parts, "\x1f")
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func buildICalItemFromState(item icalStateItem) map[string]any {
	if item.UID == "" {
		return nil
	}
	return map[string]any{
		"id":            item.ID,
		"uid":           item.UID,
		"recurrence_id": item.RecurrenceID,
		"summary":       item.Summary,
		"location":      item.Location,
		"description":   item.Description,
		"start":         item.Start,
		"end":           item.End,
		"start_raw":     item.StartRaw,
		"end_raw":       item.EndRaw,
		"all_day":       item.AllDay,
		"timezone":      item.Timezone,
		"status":        item.Status,
		"url":           item.URL,
		"organizer":     item.Organizer,
		"updated_at":    item.UpdatedAt,
		"sequence":      item.Sequence,
	}
}

func buildICalKey(uid, recurrenceID string) string {
	if recurrenceID == "" {
		return uid
	}
	return uid + "|" + recurrenceID
}

func isICalMoved(prev, curr icalStateItem) bool {
	return prev.StartRaw != curr.StartRaw ||
		prev.EndRaw != curr.EndRaw ||
		prev.Start != curr.Start ||
		prev.End != curr.End ||
		prev.Timezone != curr.Timezone ||
		prev.AllDay != curr.AllDay
}

func statusChangedToCancelled(prev, curr string) bool {
	return !isCancelledStatus(prev) && isCancelledStatus(curr)
}

func isCancelledStatus(status string) bool {
	status = strings.ToUpper(strings.TrimSpace(status))
	return status == "CANCELLED" || status == "CANCELED"
}

func getICalString(item map[string]any, key string) string {
	val, ok := item[key]
	if !ok || val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		return strings.TrimSpace(v)
	case fmt.Stringer:
		return strings.TrimSpace(v.String())
	default:
		return strings.TrimSpace(fmt.Sprint(v))
	}
}

func getICalBool(item map[string]any, key string) bool {
	val, ok := item[key]
	if !ok || val == nil {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	case string:
		parsed, err := strconv.ParseBool(strings.TrimSpace(v))
		return err == nil && parsed
	case float64:
		return v != 0
	case int:
		return v != 0
	case int64:
		return v != 0
	default:
		return false
	}
}

func getICalInt(item map[string]any, key string) int {
	val, ok := item[key]
	if !ok || val == nil {
		return 0
	}
	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case json.Number:
		i, err := v.Int64()
		if err != nil {
			return 0
		}
		return int(i)
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return 0
		}
		return parsed
	default:
		return 0
	}
}
