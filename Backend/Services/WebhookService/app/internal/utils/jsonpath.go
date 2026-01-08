package utils

import (
	"strconv"
	"strings"
)

func ExtractJSONPath(data any, path string) (any, bool) {
	path = strings.TrimSpace(path)
	if strings.HasPrefix(path, "$.") {
		path = strings.TrimPrefix(path, "$.")
	}
	path = strings.TrimPrefix(path, "$")
	path = strings.TrimPrefix(path, ".")
	if path == "" {
		return data, true
	}

	segments := strings.Split(path, ".")
	current := data

	for _, segment := range segments {
		if segment == "" {
			return nil, false
		}
		key := segment
		for {
			bracket := strings.IndexRune(key, '[')
			if bracket == -1 {
				if key != "" {
					obj, ok := current.(map[string]any)
					if !ok {
						return nil, false
					}
					var exists bool
					current, exists = obj[key]
					if !exists {
						return nil, false
					}
				}
				break
			}

			if bracket > 0 {
				base := key[:bracket]
				obj, ok := current.(map[string]any)
				if !ok {
					return nil, false
				}
				var exists bool
				current, exists = obj[base]
				if !exists {
					return nil, false
				}
			}

			rest := key[bracket:]
			end := strings.IndexRune(rest, ']')
			if end == -1 {
				return nil, false
			}

			indexStr := rest[1:end]
			idx, err := strconv.Atoi(indexStr)
			if err != nil {
				return nil, false
			}

			arr, ok := current.([]any)
			if !ok {
				return nil, false
			}
			if idx < 0 || idx >= len(arr) {
				return nil, false
			}
			current = arr[idx]

			key = rest[end+1:]
			if key == "" {
				break
			}
		}
	}

	return current, true
}
