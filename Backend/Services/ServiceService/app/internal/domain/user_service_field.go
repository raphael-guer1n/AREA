package domain

import (
	"encoding/json"
	"time"
)

type UserServiceField struct {
	ID          int
	ProfileId   int
	FieldKey    string
	StringValue string
	NumberValue float64
	BoolValue   bool
	JsonValue   json.RawMessage
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
