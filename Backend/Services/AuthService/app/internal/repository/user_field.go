package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/raphael-guer1n/AREA/AuthService/internal/domain"
)

type userServiceFieldRepository struct {
	db *sql.DB
}

func (r *userServiceFieldRepository) GetFieldsByProfileId(profileId int) ([]domain.UserServiceField, error) {
	rows, err := r.db.Query(`SELECT id, profile_id, field_key, value_string, value_number, value_boolean, value_json, created_at, updated_at FROM user_service_fields WHERE profile_id = $1`, profileId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var fields []domain.UserServiceField
	for rows.Next() {
		var field domain.UserServiceField
		if err := rows.Scan(&field.ID, &field.ProfileId, &field.FieldKey, &field.StringValue, &field.NumberValue, &field.BoolValue, &field.JsonValue, &field.CreatedAt, &field.UpdatedAt); err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}
	return fields, rows.Err()
}

func NewUserServiceFieldRepository(db *sql.DB) domain.UserServiceFieldRepository {
	return &userServiceFieldRepository{db: db}
}

func (r *userServiceFieldRepository) CreateBatch(fields []domain.UserServiceField) error {
	if len(fields) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(fields))
	valueArgs := make([]interface{}, 0, len(fields)*6)
	argPos := 1

	for _, field := range fields {
		valueStrings = append(
			valueStrings,
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", argPos, argPos+1, argPos+2, argPos+3, argPos+4, argPos+5),
		)

		var jsonArg interface{}
		if field.JsonValue != nil && len(*field.JsonValue) > 0 {
			jsonArg = string(*field.JsonValue) // valid JSON text
		} else {
			jsonArg = nil // will be NULL in DB
		}

		valueArgs = append(
			valueArgs,
			field.ProfileId,
			field.FieldKey,
			field.StringValue,
			field.NumberValue,
			field.BoolValue,
			jsonArg,
		)
		argPos += 6
	}

	query := fmt.Sprintf(`
		INSERT INTO user_service_fields (
			profile_id,
			field_key,
			value_string,
			value_number,
			value_boolean,
			value_json
		)
		VALUES %s
		ON CONFLICT (profile_id, field_key)
		DO UPDATE SET
			value_string = EXCLUDED.value_string,
			value_number = EXCLUDED.value_number,
			value_boolean = EXCLUDED.value_boolean,
			value_json = EXCLUDED.value_json,
			updated_at = NOW()
	`, strings.Join(valueStrings, ","))

	_, err := r.db.Exec(query, valueArgs...)
	return err
}
