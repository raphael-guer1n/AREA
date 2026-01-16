package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/raphael-guer1n/AREA/CronService/internal/domain"
)

type ActionRepository struct {
	db *sql.DB
}

func NewActionRepository(db *sql.DB) *ActionRepository {
	return &ActionRepository{db: db}
}

func (r *ActionRepository) Create(action *domain.Action) error {
	inputJSON, err := json.Marshal(action.Input)
	if err != nil {
		return fmt.Errorf("failed to marshal input: %w", err)
	}

	query := `
		INSERT INTO cron_actions (action_id, active, type, provider, service, title, input, cron_job_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`
	err = r.db.QueryRow(
		query,
		action.ActionID,
		action.Active,
		action.Type,
		action.Provider,
		action.Service,
		action.Title,
		inputJSON,
		action.CronJobID,
	).Scan(&action.CreatedAt, &action.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create action: %w", err)
	}

	return nil
}

func (r *ActionRepository) GetByActionID(actionID int) (*domain.Action, error) {
	query := `
		SELECT action_id, active, type, provider, service, title, input, cron_job_id, created_at, updated_at
		FROM cron_actions
		WHERE action_id = $1
	`

	action := &domain.Action{}
	var inputJSON []byte

	err := r.db.QueryRow(query, actionID).Scan(
		&action.ActionID,
		&action.Active,
		&action.Type,
		&action.Provider,
		&action.Service,
		&action.Title,
		&inputJSON,
		&action.CronJobID,
		&action.CreatedAt,
		&action.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("action not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get action: %w", err)
	}

	if err := json.Unmarshal(inputJSON, &action.Input); err != nil {
		return nil, fmt.Errorf("failed to unmarshal input: %w", err)
	}

	return action, nil
}

func (r *ActionRepository) GetAll() ([]*domain.Action, error) {
	query := `
		SELECT action_id, active, type, provider, service, title, input, cron_job_id, created_at, updated_at
		FROM cron_actions
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get actions: %w", err)
	}
	defer rows.Close()

	var actions []*domain.Action
	for rows.Next() {
		action := &domain.Action{}
		var inputJSON []byte

		err := rows.Scan(
			&action.ActionID,
			&action.Active,
			&action.Type,
			&action.Provider,
			&action.Service,
			&action.Title,
			&inputJSON,
			&action.CronJobID,
			&action.CreatedAt,
			&action.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan action: %w", err)
		}

		if err := json.Unmarshal(inputJSON, &action.Input); err != nil {
			return nil, fmt.Errorf("failed to unmarshal input: %w", err)
		}

		actions = append(actions, action)
	}

	return actions, nil
}

func (r *ActionRepository) Update(action *domain.Action) error {
	inputJSON, err := json.Marshal(action.Input)
	if err != nil {
		return fmt.Errorf("failed to marshal input: %w", err)
	}

	query := `
		UPDATE cron_actions
		SET active = $1, type = $2, provider = $3, service = $4, title = $5, input = $6, cron_job_id = $7, updated_at = NOW()
		WHERE action_id = $8
		RETURNING updated_at
	`

	err = r.db.QueryRow(
		query,
		action.Active,
		action.Type,
		action.Provider,
		action.Service,
		action.Title,
		inputJSON,
		action.CronJobID,
		action.ActionID,
	).Scan(&action.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update action: %w", err)
	}

	return nil
}

func (r *ActionRepository) Delete(actionID int) error {
	query := `DELETE FROM cron_actions WHERE action_id = $1`

	result, err := r.db.Exec(query, actionID)
	if err != nil {
		return fmt.Errorf("failed to delete action: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("action not found")
	}

	return nil
}
