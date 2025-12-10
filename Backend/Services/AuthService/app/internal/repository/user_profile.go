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
	err := r.db.QueryRow(
		`INSERT INTO user_service_profiles (
             user_id,
             service,
             provider_user_id,
             access_token,
             refresh_token,
             expires_at,
             raw_profile
         )
         VALUES ($1, $2, $3, $4, $5, $6, $7)
         RETURNING
             id,
             user_id,
             service,
             provider_user_id,
             access_token,
             refresh_token,
             expires_at,
             raw_profile,
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
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	return u, err
}

func (r *userProfileRepository) GetServicesByUserId(userId int) ([]string, error) {
	rows, err := r.db.Query(
		`SELECT DISTINCT service FROM user_service_profiles WHERE user_id = $1`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []string
	for rows.Next() {
		var service string
		if err := rows.Scan(&service); err != nil {
			return nil, err
		}
		services = append(services, service)
	}

	return services, rows.Err()
}
