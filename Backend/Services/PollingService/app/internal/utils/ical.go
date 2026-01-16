package utils

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type icalProperty struct {
	Name   string
	Params map[string][]string
	Value  string
}

type icalEvent struct {
	Props map[string][]icalProperty
}

type icalItem struct {
	Fields    map[string]any
	UpdatedAt time.Time
	UID       string
}

func (e *icalEvent) addProp(prop icalProperty) {
	if e.Props == nil {
		e.Props = make(map[string][]icalProperty)
	}
	e.Props[prop.Name] = append(e.Props[prop.Name], prop)
}

func (e icalEvent) firstProp(name string) (icalProperty, bool) {
	props := e.Props[name]
	if len(props) == 0 {
		return icalProperty{}, false
	}
	return props[0], true
}

func (e icalEvent) firstValue(name string) string {
	prop, ok := e.firstProp(name)
	if !ok {
		return ""
	}
	return prop.Value
}

// ParseICalToItems parses iCalendar data into a list of event items.
// Each item is a flat map suitable for JSON path extraction.
func ParseICalToItems(data []byte) ([]map[string]any, error) {
	lines, err := unfoldICalLines(data)
	if err != nil {
		return nil, err
	}

	var (
		stack      []string
		events     []icalEvent
		current    *icalEvent
		method     string
		defaultTZ  string
		insideCal  bool
	)

	for _, rawLine := range lines {
		line := rawLine
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.HasPrefix(line, "BEGIN:") {
			component := strings.ToUpper(strings.TrimSpace(strings.TrimPrefix(line, "BEGIN:")))
			stack = append(stack, component)
			if component == "VCALENDAR" {
				insideCal = true
			}
			if component == "VEVENT" {
				current = &icalEvent{Props: make(map[string][]icalProperty)}
			}
			continue
		}
		if strings.HasPrefix(line, "END:") {
			component := strings.ToUpper(strings.TrimSpace(strings.TrimPrefix(line, "END:")))
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
			if component == "VEVENT" && current != nil {
				events = append(events, *current)
				current = nil
			}
			if component == "VCALENDAR" {
				insideCal = false
			}
			continue
		}
		if len(stack) == 0 || !insideCal {
			continue
		}
		prop, ok := parseICalProperty(line)
		if !ok {
			continue
		}

		switch stack[len(stack)-1] {
		case "VCALENDAR":
			switch prop.Name {
			case "METHOD":
				method = strings.ToUpper(strings.TrimSpace(prop.Value))
			case "X-WR-TIMEZONE":
				defaultTZ = strings.TrimSpace(prop.Value)
			}
		case "VEVENT":
			if current != nil {
				current.addProp(prop)
			}
		}
	}

	items := make([]icalItem, 0, len(events))
	for _, ev := range events {
		item, ok := buildICalItem(ev, method, defaultTZ)
		if !ok {
			continue
		}
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		left := items[i].UpdatedAt
		right := items[j].UpdatedAt
		if !left.Equal(right) {
			return left.After(right)
		}
		return items[i].UID > items[j].UID
	})

	out := make([]map[string]any, 0, len(items))
	for _, item := range items {
		out = append(out, item.Fields)
	}
	return out, nil
}

func unfoldICalLines(data []byte) ([]string, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var lines []string
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")
		if len(line) > 0 && (line[0] == ' ' || line[0] == '\t') {
			if len(lines) == 0 {
				lines = append(lines, strings.TrimLeft(line, " \t"))
				continue
			}
			lines[len(lines)-1] += strings.TrimLeft(line, " \t")
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func parseICalProperty(line string) (icalProperty, bool) {
	colon := strings.Index(line, ":")
	if colon == -1 {
		return icalProperty{}, false
	}
	left := line[:colon]
	value := line[colon+1:]
	parts := splitICalTokens(left)
	if len(parts) == 0 {
		return icalProperty{}, false
	}

	name := strings.ToUpper(strings.TrimSpace(parts[0]))
	if name == "" {
		return icalProperty{}, false
	}
	if idx := strings.LastIndex(name, "."); idx != -1 && idx < len(name)-1 {
		name = name[idx+1:]
	}

	params := make(map[string][]string)
	for _, param := range parts[1:] {
		key, rawValue, ok := strings.Cut(param, "=")
		if !ok {
			continue
		}
		key = strings.ToUpper(strings.TrimSpace(key))
		if key == "" {
			continue
		}
		for _, val := range splitICalParamValues(strings.TrimSpace(rawValue)) {
			val = strings.Trim(val, `"`)
			if val == "" {
				continue
			}
			params[key] = append(params[key], val)
		}
	}

	return icalProperty{
		Name:   name,
		Params: params,
		Value:  value,
	}, true
}

func splitICalTokens(input string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false
	for _, r := range input {
		switch r {
		case '"':
			inQuotes = !inQuotes
			current.WriteRune(r)
		case ';':
			if inQuotes {
				current.WriteRune(r)
				continue
			}
			parts = append(parts, current.String())
			current.Reset()
		default:
			current.WriteRune(r)
		}
	}
	parts = append(parts, current.String())
	return parts
}

func splitICalParamValues(input string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false
	for _, r := range input {
		switch r {
		case '"':
			inQuotes = !inQuotes
			current.WriteRune(r)
		case ',':
			if inQuotes {
				current.WriteRune(r)
				continue
			}
			parts = append(parts, current.String())
			current.Reset()
		default:
			current.WriteRune(r)
		}
	}
	parts = append(parts, current.String())
	return parts
}

func buildICalItem(ev icalEvent, method string, defaultTZ string) (icalItem, bool) {
	uid := strings.TrimSpace(ev.firstValue("UID"))
	if uid == "" {
		return icalItem{}, false
	}

	summary := unescapeICalText(ev.firstValue("SUMMARY"))
	location := unescapeICalText(ev.firstValue("LOCATION"))
	description := unescapeICalText(ev.firstValue("DESCRIPTION"))
	url := strings.TrimSpace(ev.firstValue("URL"))

	organizer := extractOrganizer(ev)
	status := strings.ToUpper(strings.TrimSpace(ev.firstValue("STATUS")))
	if status == "" && strings.EqualFold(method, "CANCEL") {
		status = "CANCELLED"
	}

	sequence := parseICalInt(ev.firstValue("SEQUENCE"))

	startProp, hasStart := ev.firstProp("DTSTART")
	startTime, allDay, tzName, startRaw := time.Time{}, false, "", ""
	if hasStart {
		startTime, allDay, tzName, startRaw = parseICalPropertyTime(startProp, defaultTZ)
	}

	endProp, hasEnd := ev.firstProp("DTEND")
	endTime, _, _, endRaw := time.Time{}, false, "", ""
	if hasEnd {
		endTime, _, _, endRaw = parseICalPropertyTime(endProp, defaultTZ)
	}
	if endTime.IsZero() && !startTime.IsZero() {
		if duration := strings.TrimSpace(ev.firstValue("DURATION")); duration != "" {
			if parsed, err := parseICalDuration(duration); err == nil {
				endTime = startTime.Add(parsed)
			}
		}
	}
	if endTime.IsZero() && allDay && !startTime.IsZero() {
		endTime = startTime.AddDate(0, 0, 1)
	}

	recurrenceProp, hasRecurrence := ev.firstProp("RECURRENCE-ID")
	recurrenceTime, _, _, recurrenceRaw := time.Time{}, false, "", ""
	if hasRecurrence {
		recurrenceTime, _, _, recurrenceRaw = parseICalPropertyTime(recurrenceProp, defaultTZ)
	}

	lastMod := parseICalTimestamp(ev, "LAST-MODIFIED", defaultTZ)
	dtStamp := parseICalTimestamp(ev, "DTSTAMP", defaultTZ)
	created := parseICalTimestamp(ev, "CREATED", defaultTZ)
	updatedAt := maxICalTime(lastMod, dtStamp, created)

	updatedAtStr := ""
	if !updatedAt.IsZero() {
		updatedAtStr = updatedAt.Format(time.RFC3339)
	}

	startStr := ""
	if !startTime.IsZero() {
		startStr = startTime.Format(time.RFC3339)
	}

	endStr := ""
	if !endTime.IsZero() {
		endStr = endTime.Format(time.RFC3339)
	}

	recurrenceStr := ""
	if !recurrenceTime.IsZero() {
		recurrenceStr = recurrenceTime.Format(time.RFC3339)
	} else if recurrenceRaw != "" {
		recurrenceStr = recurrenceRaw
	}

	if tzName == "" {
		tzName = strings.TrimSpace(defaultTZ)
	}
	if tzName == "" {
		tzName = "UTC"
	}

	id := buildICalItemID(uid, recurrenceStr, sequence, updatedAtStr, summary, startRaw, endRaw, status, location, description)

	fields := map[string]any{
		"id":            id,
		"uid":           uid,
		"recurrence_id": recurrenceStr,
		"summary":       summary,
		"location":      location,
		"description":   description,
		"start":         startStr,
		"end":           endStr,
		"start_raw":     startRaw,
		"end_raw":       endRaw,
		"all_day":       allDay,
		"timezone":      tzName,
		"status":        status,
		"url":           url,
		"organizer":     organizer,
		"updated_at":    updatedAtStr,
		"sequence":      sequence,
	}

	return icalItem{
		Fields:    fields,
		UpdatedAt: updatedAt,
		UID:       uid,
	}, true
}

func extractOrganizer(ev icalEvent) string {
	prop, ok := ev.firstProp("ORGANIZER")
	if !ok {
		return ""
	}
	cn := getICalParam(prop.Params, "CN")
	if cn != "" {
		return unescapeICalText(cn)
	}
	value := strings.TrimSpace(prop.Value)
	valueLower := strings.ToLower(value)
	if strings.HasPrefix(valueLower, "mailto:") {
		value = value[len("mailto:"):]
	}
	return value
}

func parseICalPropertyTime(prop icalProperty, defaultTZ string) (time.Time, bool, string, string) {
	raw := strings.TrimSpace(prop.Value)
	if raw == "" {
		return time.Time{}, false, "", ""
	}
	parsed, allDay, tzName, err := parseICalDateTime(raw, prop.Params, defaultTZ)
	if err != nil {
		return time.Time{}, false, tzName, raw
	}
	return parsed, allDay, tzName, raw
}

func parseICalTimestamp(ev icalEvent, name string, defaultTZ string) time.Time {
	prop, ok := ev.firstProp(name)
	if !ok {
		return time.Time{}
	}
	parsed, _, _, err := parseICalDateTime(strings.TrimSpace(prop.Value), prop.Params, defaultTZ)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

func parseICalDateTime(raw string, params map[string][]string, defaultTZ string) (time.Time, bool, string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, false, "", fmt.Errorf("empty datetime")
	}

	valueType := strings.ToUpper(getICalParam(params, "VALUE"))
	tzid := strings.TrimSpace(getICalParam(params, "TZID"))
	loc, tzName := resolveICalLocation(tzid, defaultTZ)

	if valueType == "DATE" || (len(raw) == 8 && !strings.Contains(raw, "T")) {
		parsed, err := time.ParseInLocation("20060102", raw, loc)
		if err != nil {
			return time.Time{}, true, tzName, err
		}
		return parsed, true, tzName, nil
	}

	parsed, err := parseICalDateTimeValue(raw, loc)
	if err != nil {
		return time.Time{}, false, tzName, err
	}
	return parsed, false, tzName, nil
}

func parseICalDateTimeValue(raw string, loc *time.Location) (time.Time, error) {
	if strings.HasSuffix(raw, "Z") {
		for _, layout := range []string{"20060102T150405Z", "20060102T1504Z"} {
			if parsed, err := time.Parse(layout, raw); err == nil {
				return parsed, nil
			}
		}
	}

	if hasICalOffset(raw) {
		normalized := raw
		if len(raw) >= 6 {
			offset := raw[len(raw)-6:]
			if strings.Contains(offset, ":") {
				normalized = raw[:len(raw)-6] + strings.Replace(offset, ":", "", 1)
			}
		}
		for _, layout := range []string{"20060102T150405-0700", "20060102T1504-0700"} {
			if parsed, err := time.Parse(layout, normalized); err == nil {
				return parsed, nil
			}
		}
	}

	for _, layout := range []string{"20060102T150405", "20060102T1504"} {
		if parsed, err := time.ParseInLocation(layout, raw, loc); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid datetime %q", raw)
}

func hasICalOffset(raw string) bool {
	if len(raw) < len("20060102T1504")+5 {
		return false
	}
	idx := strings.LastIndexAny(raw, "+-")
	if idx == -1 {
		return false
	}
	return idx > len("20060102T1504")
}

func resolveICalLocation(tzid string, defaultTZ string) (*time.Location, string) {
	tz := strings.TrimSpace(tzid)
	if tz == "" {
		tz = strings.TrimSpace(defaultTZ)
	}
	if tz == "" || strings.EqualFold(tz, "UTC") || strings.EqualFold(tz, "GMT") {
		return time.UTC, "UTC"
	}
	if loc, err := time.LoadLocation(tz); err == nil {
		return loc, tz
	}
	if offset, ok := parseICalOffset(tz); ok {
		return time.FixedZone(tz, offset), tz
	}
	return time.UTC, tz
}

func parseICalOffset(value string) (int, bool) {
	clean := strings.TrimSpace(value)
	clean = strings.TrimPrefix(strings.TrimPrefix(clean, "UTC"), "GMT")
	clean = strings.TrimSpace(clean)
	if len(clean) == 6 && clean[3] == ':' {
		clean = clean[:3] + clean[4:]
	}
	if len(clean) != 5 {
		return 0, false
	}
	sign := clean[0]
	if sign != '+' && sign != '-' {
		return 0, false
	}
	hours, err := strconv.Atoi(clean[1:3])
	if err != nil {
		return 0, false
	}
	mins, err := strconv.Atoi(clean[3:5])
	if err != nil {
		return 0, false
	}
	offset := hours*3600 + mins*60
	if sign == '-' {
		offset = -offset
	}
	return offset, true
}

func parseICalDuration(value string) (time.Duration, error) {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return 0, fmt.Errorf("empty duration")
	}
	sign := 1
	if raw[0] == '-' {
		sign = -1
		raw = raw[1:]
	}
	if raw == "" || raw[0] != 'P' {
		return 0, fmt.Errorf("invalid duration %q", value)
	}
	raw = raw[1:]

	var (
		number strings.Builder
		inTime bool
		weeks  int
		days   int
		hours  int
		mins   int
		secs   int
	)

	flush := func(unit byte) error {
		if number.Len() == 0 {
			return fmt.Errorf("invalid duration %q", value)
		}
		val, err := strconv.Atoi(number.String())
		if err != nil {
			return err
		}
		number.Reset()
		switch unit {
		case 'W':
			weeks = val
		case 'D':
			days = val
		case 'H':
			hours = val
		case 'M':
			if !inTime {
				return fmt.Errorf("months not supported in duration %q", value)
			}
			mins = val
		case 'S':
			secs = val
		default:
			return fmt.Errorf("invalid duration %q", value)
		}
		return nil
	}

	for i := 0; i < len(raw); i++ {
		ch := raw[i]
		switch {
		case ch >= '0' && ch <= '9':
			number.WriteByte(ch)
		case ch == 'T':
			inTime = true
		default:
			if err := flush(ch); err != nil {
				return 0, err
			}
		}
	}
	if number.Len() > 0 {
		return 0, fmt.Errorf("invalid duration %q", value)
	}

	total := time.Duration(weeks*7*24)*time.Hour +
		time.Duration(days*24)*time.Hour +
		time.Duration(hours)*time.Hour +
		time.Duration(mins)*time.Minute +
		time.Duration(secs)*time.Second

	if sign < 0 {
		total = -total
	}
	return total, nil
}

func getICalParam(params map[string][]string, key string) string {
	values := params[strings.ToUpper(key)]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func parseICalInt(value string) int {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	if parsed, err := strconv.Atoi(value); err == nil {
		return parsed
	}
	return 0
}

func buildICalItemID(uid string, recurrenceID string, sequence int, updatedAt string, summary string, startRaw string, endRaw string, status string, location string, description string) string {
	parts := []string{
		uid,
		recurrenceID,
		strconv.Itoa(sequence),
		updatedAt,
		summary,
		startRaw,
		endRaw,
		status,
		location,
		description,
	}
	base := strings.Join(parts, "|")
	hash := sha256.Sum256([]byte(base))
	return hex.EncodeToString(hash[:])
}

func maxICalTime(times ...time.Time) time.Time {
	var out time.Time
	for _, t := range times {
		if t.After(out) {
			out = t
		}
	}
	return out
}

func unescapeICalText(value string) string {
	if value == "" {
		return ""
	}
	replacer := strings.NewReplacer(
		`\\`, `\`,
		`\n`, "\n",
		`\N`, "\n",
		`\;`, ";",
		`\,`, ",",
	)
	return replacer.Replace(value)
}
