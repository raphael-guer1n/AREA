package domain

import "encoding/json"

type BodyField struct {
	Path  string          `json:"path"`
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}
