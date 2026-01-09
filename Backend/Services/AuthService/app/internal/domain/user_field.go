package domain

import (
	"encoding/json"
	"time"
)

type UserServiceField struct {
	ID          int              `json:"id"`
	ProfileId   int              `json:"profile_id"`
	FieldKey    string           `json:"field_key"`
	StringValue string           `json:"string_value"`
	NumberValue float64          `json:"number_value"`
	BoolValue   bool             `json:"bool_value"`
	JsonValue   *json.RawMessage `json:"json_value"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type UserServiceFieldRepository interface {
	CreateBatch(fields []UserServiceField) error
	GetFieldsByProfileId(profileId int) ([]UserServiceField, error)
}
