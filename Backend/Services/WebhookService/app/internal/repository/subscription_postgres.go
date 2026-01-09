package repository

import (
	"database/sql"
	"errors"

	"github.com/raphael-guer1n/AREA/WebhookService/internal/domain"
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
	err := r.db.QueryRow(
		`INSERT INTO webhook_subscriptions (hook_id, user_id, area_id, provider, config, provider_hook_id)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, hook_id, user_id, area_id, provider, config, provider_hook_id, created_at, updated_at`,
		sub.HookID, sub.UserID, sub.AreaID, sub.Provider, cfg, sub.ProviderHookID,
	).Scan(
		&created.ID,
		&created.HookID,
		&created.UserID,
		&created.AreaID,
		&created.Provider,
		&configBytes,
		&created.ProviderHookID,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if configBytes != nil {
		created.Config = configBytes
	}
	return &created, nil
}

func (r *subscriptionRepository) FindByHookID(hookID string) (*domain.Subscription, error) {
	var sub domain.Subscription
	var configBytes []byte
	err := r.db.QueryRow(
		`SELECT id, hook_id, user_id, area_id, provider, config, provider_hook_id, created_at, updated_at
		 FROM webhook_subscriptions WHERE hook_id = $1`,
		hookID,
	).Scan(
		&sub.ID,
		&sub.HookID,
		&sub.UserID,
		&sub.AreaID,
		&sub.Provider,
		&configBytes,
		&sub.ProviderHookID,
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
	return &sub, nil
}

func (r *subscriptionRepository) ListByUserID(userID int) ([]domain.Subscription, error) {
	rows, err := r.db.Query(
		`SELECT id, hook_id, user_id, area_id, provider, config, provider_hook_id, created_at, updated_at
		 FROM webhook_subscriptions WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []domain.Subscription
	for rows.Next() {
		var sub domain.Subscription
		var configBytes []byte
		if err := rows.Scan(
			&sub.ID,
			&sub.HookID,
			&sub.UserID,
			&sub.AreaID,
			&sub.Provider,
			&configBytes,
			&sub.ProviderHookID,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if configBytes != nil {
			sub.Config = configBytes
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

func (r *subscriptionRepository) UpdateProviderHookID(hookID, providerHookID string) error {
	_, err := r.db.Exec(
		`UPDATE webhook_subscriptions SET provider_hook_id = $1, updated_at = NOW() WHERE hook_id = $2`,
		providerHookID, hookID,
	)
	return err
}

func (r *subscriptionRepository) DeleteByHookID(hookID string) error {
	_, err := r.db.Exec(`DELETE FROM webhook_subscriptions WHERE hook_id = $1`, hookID)
	return err
}
