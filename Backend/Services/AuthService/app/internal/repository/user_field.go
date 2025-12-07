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
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", argPos, argPos+1, argPos+2, argPos+3, argPos+4, argPos+5))
		valueArgs = append(valueArgs, field.ProfileId, field.FieldKey, field.StringValue, field.NumberValue, field.BoolValue, field.JsonValue)
		argPos += 6
	}

	query := fmt.Sprintf(`
		INSERT INTO user_service_fields (profile_id, field_key, value_string, value_number, value_boolean, value_json)
		VALUES %s
	`, strings.Join(valueStrings, ","))

	_, err := r.db.Exec(query, valueArgs...)
	return err
}
