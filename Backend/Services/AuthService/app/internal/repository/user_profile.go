package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/raphael-guer1n/AREA/AuthService/internal/domain"
)

type userProfileRepository struct {
	db *sql.DB
}

func (r *userProfileRepository) GetProviderProfileProfileByServiceByUser(userId int, service string) (domain.UserProfile, error) {
	var userProfile domain.UserProfile
	var lastRefreshError sql.NullString
	var lastRefreshAt sql.NullTime
	err := r.db.QueryRow(
		`SELECT id, user_id, service, provider_user_id, access_token, refresh_token, expires_at, raw_profile, needs_reconnect, last_refresh_error, last_refresh_at, created_at, updated_at FROM user_service_profiles WHERE user_id = $1 AND service = $2`,
		userId, service,
	).Scan(
		&userProfile.ID,
		&userProfile.UserId,
		&userProfile.Service,
		&userProfile.ProviderUserId,
		&userProfile.AccessToken,
		&userProfile.RefreshToken,
		&userProfile.ExpiresAt,
		&userProfile.RawProfile,
		&userProfile.NeedsReconnect,
		&lastRefreshError,
		&lastRefreshAt,
		&userProfile.CreatedAt,
		&userProfile.UpdatedAt,
	)
	if lastRefreshError.Valid {
		userProfile.LastRefreshError = &lastRefreshError.String
	}
	if lastRefreshAt.Valid {
		userProfile.LastRefreshAt = &lastRefreshAt.Time
	}
	return userProfile, err
}

func (r *userProfileRepository) GetProviderUserTokenByServiceByUserId(userId int, service string) (string, error) {
	var providerToken string

	err := r.db.QueryRow(
		`SELECT access_token FROM user_service_profiles WHERE user_id = $1 AND service = $2`,
		userId, service,
	).Scan(&providerToken)
	return providerToken, err
}

func NewUserProfileRepository(db *sql.DB) domain.UserProfileRepository {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) Create(
	userId int,
	service, providerUserId, accessToken, refreshToken string,
	expiresAt time.Time,
	rawProfile json.RawMessage,
) (domain.UserProfile, error) {
	var u domain.UserProfile
	var lastRefreshError sql.NullString
	var lastRefreshAt sql.NullTime
	err := r.db.QueryRow(
		`INSERT INTO user_service_profiles (
			 user_id,
			 service,
			 provider_user_id,
			 access_token,
			 refresh_token,
			 expires_at,
			 raw_profile,
			 needs_reconnect,
			 last_refresh_error,
			 last_refresh_at
		 )
		 VALUES ($1, $2, $3, $4, $5, $6, $7, false, NULL, NULL)
		 ON CONFLICT (user_id, service)
		 DO UPDATE SET
			 provider_user_id = EXCLUDED.provider_user_id,
			 access_token     = EXCLUDED.access_token,
			 refresh_token    = COALESCE(NULLIF(EXCLUDED.refresh_token, ''), user_service_profiles.refresh_token),
			 expires_at       = EXCLUDED.expires_at,
			 raw_profile      = EXCLUDED.raw_profile,
			 needs_reconnect  = false,
			 last_refresh_error = NULL,
			 last_refresh_at  = NULL,
			 updated_at       = NOW()
		 RETURNING
			 id,
			 user_id,
			 service,
			 provider_user_id,
			 access_token,
			 refresh_token,
			 expires_at,
			 raw_profile,
			 needs_reconnect,
			 last_refresh_error,
			 last_refresh_at,
			 created_at,
			 updated_at`,
		userId, service, providerUserId, accessToken, refreshToken, expiresAt, rawProfile,
	).Scan(
		&u.ID,
		&u.UserId,
		&u.Service,
		&u.ProviderUserId,
		&u.AccessToken,
		&u.RefreshToken,
		&u.ExpiresAt,
		&u.RawProfile,
		&u.NeedsReconnect,
		&lastRefreshError,
		&lastRefreshAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if lastRefreshError.Valid {
		u.LastRefreshError = &lastRefreshError.String
	}
	if lastRefreshAt.Valid {
		u.LastRefreshAt = &lastRefreshAt.Time
	}
	return u, err
}

func (r *userProfileRepository) GetServicesStatusByUserId(userId int) ([]domain.ServiceStatus, error) {
	rows, err := r.db.Query(
		`SELECT service, needs_reconnect FROM user_service_profiles WHERE user_id = $1`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []domain.ServiceStatus
	for rows.Next() {
		var status domain.ServiceStatus
		if err := rows.Scan(&status.Service, &status.NeedsReconnect); err != nil {
			return nil, err
		}
		services = append(services, status)
	}

	return services, rows.Err()
}

func (r *userProfileRepository) ListRefreshCandidates(expireBefore time.Time) ([]domain.RefreshCandidate, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, service, refresh_token, expires_at
		 FROM user_service_profiles
		 WHERE needs_reconnect = false
		   AND expires_at <= $1`,
		expireBefore,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candidates []domain.RefreshCandidate
	for rows.Next() {
		var candidate domain.RefreshCandidate
		if err := rows.Scan(
			&candidate.ID,
			&candidate.UserId,
			&candidate.Service,
			&candidate.RefreshToken,
			&candidate.ExpiresAt,
		); err != nil {
			return nil, err
		}
		candidates = append(candidates, candidate)
	}

	return candidates, rows.Err()
}

func (r *userProfileRepository) UpdateTokens(profileID int, accessToken, refreshToken string, expiresAt time.Time) error {
	_, err := r.db.Exec(
		`UPDATE user_service_profiles
		 SET access_token = $1,
		     refresh_token = COALESCE(NULLIF($2, ''), refresh_token),
		     expires_at = $3,
		     needs_reconnect = false,
		     last_refresh_error = NULL,
		     last_refresh_at = NOW(),
		     updated_at = NOW()
		 WHERE id = $4`,
		accessToken,
		refreshToken,
		expiresAt,
		profileID,
	)
	return err
}

func (r *userProfileRepository) MarkNeedsReconnect(profileID int, reason string) error {
	_, err := r.db.Exec(
		`UPDATE user_service_profiles
		 SET needs_reconnect = true,
		     last_refresh_error = $1,
		     last_refresh_at = NOW(),
		     updated_at = NOW()
		 WHERE id = $2`,
		reason,
		profileID,
	)
	return err
}
