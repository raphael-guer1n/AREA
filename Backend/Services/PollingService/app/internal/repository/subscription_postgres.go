package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/raphael-guer1n/AREA/PollingService/internal/domain"
)

type subscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) domain.SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(sub *domain.Subscription) (*domain.Subscription, error) {
	cfg := []byte("{}")
	if len(sub.Config) > 0 {
		cfg = []byte(sub.Config)
	}

	var created domain.Subscription
	var configBytes []byte
	var lastItemID sql.NullString
	var lastPolledAt sql.NullTime
	var nextRunAt sql.NullTime
	var lastError sql.NullString

	if sub.LastItemID != "" {
		lastItemID = sql.NullString{String: sub.LastItemID, Valid: true}
	}
	if sub.LastPolledAt != nil {
		lastPolledAt = sql.NullTime{Time: *sub.LastPolledAt, Valid: true}
	}
	if sub.NextRunAt != nil {
		nextRunAt = sql.NullTime{Time: *sub.NextRunAt, Valid: true}
	}
	if sub.LastError != "" {
		lastError = sql.NullString{String: sub.LastError, Valid: true}
	}

	err := r.db.QueryRow(
		`INSERT INTO polling_subscriptions (user_id, action_id, provider, service, active, config, interval_seconds, last_item_id, last_polled_at, next_run_at, last_error)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		 RETURNING id, user_id, action_id, provider, service, active, config, interval_seconds, last_item_id, last_polled_at, next_run_at, last_error, created_at, updated_at`,
		sub.UserID,
		sub.ActionID,
		sub.Provider,
		sub.Service,
		sub.Active,
		cfg,
		sub.IntervalSeconds,
		lastItemID,
		lastPolledAt,
		nextRunAt,
		lastError,
	).Scan(
		&created.ID,
		&created.UserID,
		&created.ActionID,
		&created.Provider,
		&created.Service,
		&created.Active,
		&configBytes,
		&created.IntervalSeconds,
		&lastItemID,
		&lastPolledAt,
		&nextRunAt,
		&lastError,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if configBytes != nil {
		created.Config = configBytes
	}
	if lastItemID.Valid {
		created.LastItemID = lastItemID.String
	}
	if lastPolledAt.Valid {
		created.LastPolledAt = &lastPolledAt.Time
	}
	if nextRunAt.Valid {
		created.NextRunAt = &nextRunAt.Time
	}
	if lastError.Valid {
		created.LastError = lastError.String
	}

	return &created, nil
}

func (r *subscriptionRepository) FindByActionID(actionID int) (*domain.Subscription, error) {
	var sub domain.Subscription
	var configBytes []byte
	var lastItemID sql.NullString
	var lastPolledAt sql.NullTime
	var nextRunAt sql.NullTime
	var lastError sql.NullString

	err := r.db.QueryRow(
		`SELECT id, user_id, action_id, provider, service, active, config, interval_seconds, last_item_id, last_polled_at, next_run_at, last_error, created_at, updated_at
		 FROM polling_subscriptions WHERE action_id = $1`,
		actionID,
	).Scan(
		&sub.ID,
		&sub.UserID,
		&sub.ActionID,
		&sub.Provider,
		&sub.Service,
		&sub.Active,
		&configBytes,
		&sub.IntervalSeconds,
		&lastItemID,
		&lastPolledAt,
		&nextRunAt,
		&lastError,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if configBytes != nil {
		sub.Config = configBytes
	}
	if lastItemID.Valid {
		sub.LastItemID = lastItemID.String
	}
	if lastPolledAt.Valid {
		sub.LastPolledAt = &lastPolledAt.Time
	}
	if nextRunAt.Valid {
		sub.NextRunAt = &nextRunAt.Time
	}
	if lastError.Valid {
		sub.LastError = lastError.String
	}
	return &sub, nil
}

func (r *subscriptionRepository) ListDue(now time.Time) ([]domain.Subscription, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, action_id, provider, service, active, config, interval_seconds, last_item_id, last_polled_at, next_run_at, last_error, created_at, updated_at
		 FROM polling_subscriptions
		 WHERE active = true AND (next_run_at IS NULL OR next_run_at <= $1)
		 ORDER BY next_run_at NULLS FIRST, created_at ASC`,
		now,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []domain.Subscription
	for rows.Next() {
		var sub domain.Subscription
		var configBytes []byte
		var lastItemID sql.NullString
		var lastPolledAt sql.NullTime
		var nextRunAt sql.NullTime
		var lastError sql.NullString

		if err := rows.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.ActionID,
			&sub.Provider,
			&sub.Service,
			&sub.Active,
			&configBytes,
			&sub.IntervalSeconds,
			&lastItemID,
			&lastPolledAt,
			&nextRunAt,
			&lastError,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if configBytes != nil {
			sub.Config = configBytes
		}
		if lastItemID.Valid {
			sub.LastItemID = lastItemID.String
		}
		if lastPolledAt.Valid {
			sub.LastPolledAt = &lastPolledAt.Time
		}
		if nextRunAt.Valid {
			sub.NextRunAt = &nextRunAt.Time
		}
		if lastError.Valid {
			sub.LastError = lastError.String
		}
		subs = append(subs, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if subs == nil {
		subs = []domain.Subscription{}
	}
	return subs, nil
}

func (r *subscriptionRepository) UpdateByActionID(sub *domain.Subscription) (*domain.Subscription, error) {
	cfg := []byte("{}")
	if len(sub.Config) > 0 {
		cfg = []byte(sub.Config)
	}

	var updated domain.Subscription
	var configBytes []byte
	var lastItemID sql.NullString
	var lastPolledAt sql.NullTime
	var nextRunAt sql.NullTime
	var lastError sql.NullString

	if sub.LastItemID != "" {
		lastItemID = sql.NullString{String: sub.LastItemID, Valid: true}
	}
	if sub.LastPolledAt != nil {
		lastPolledAt = sql.NullTime{Time: *sub.LastPolledAt, Valid: true}
	}
	if sub.NextRunAt != nil {
		nextRunAt = sql.NullTime{Time: *sub.NextRunAt, Valid: true}
	}
	if sub.LastError != "" {
		lastError = sql.NullString{String: sub.LastError, Valid: true}
	}

	err := r.db.QueryRow(
		`UPDATE polling_subscriptions
		 SET provider = $1, service = $2, active = $3, config = $4, interval_seconds = $5, last_item_id = $6,
		     last_polled_at = $7, next_run_at = $8, last_error = $9, updated_at = NOW()
		 WHERE action_id = $10
		 RETURNING id, user_id, action_id, provider, service, active, config, interval_seconds, last_item_id, last_polled_at, next_run_at, last_error, created_at, updated_at`,
		sub.Provider,
		sub.Service,
		sub.Active,
		cfg,
		sub.IntervalSeconds,
		lastItemID,
		lastPolledAt,
		nextRunAt,
		lastError,
		sub.ActionID,
	).Scan(
		&updated.ID,
		&updated.UserID,
		&updated.ActionID,
		&updated.Provider,
		&updated.Service,
		&updated.Active,
		&configBytes,
		&updated.IntervalSeconds,
		&lastItemID,
		&lastPolledAt,
		&nextRunAt,
		&lastError,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if configBytes != nil {
		updated.Config = configBytes
	}
	if lastItemID.Valid {
		updated.LastItemID = lastItemID.String
	}
	if lastPolledAt.Valid {
		updated.LastPolledAt = &lastPolledAt.Time
	}
	if nextRunAt.Valid {
		updated.NextRunAt = &nextRunAt.Time
	}
	if lastError.Valid {
		updated.LastError = lastError.String
	}
	return &updated, nil
}

func (r *subscriptionRepository) UpdatePollingState(actionID int, lastItemID string, nextRunAt time.Time, lastError string, lastPolledAt time.Time) error {
	var lastItemIDValue sql.NullString
	var lastErrorValue sql.NullString
	if lastItemID != "" {
		lastItemIDValue = sql.NullString{String: lastItemID, Valid: true}
	}
	if lastError != "" {
		lastErrorValue = sql.NullString{String: lastError, Valid: true}
	}

	_, err := r.db.Exec(
		`UPDATE polling_subscriptions
		 SET last_item_id = $1, last_polled_at = $2, next_run_at = $3, last_error = $4, updated_at = NOW()
		 WHERE action_id = $5`,
		lastItemIDValue,
		lastPolledAt,
		nextRunAt,
		lastErrorValue,
		actionID,
	)
	return err
}

func (r *subscriptionRepository) DeleteByActionID(actionID int) error {
	_, err := r.db.Exec(
		`DELETE FROM polling_subscriptions WHERE action_id = $1`,
		actionID,
	)
	return err
}
