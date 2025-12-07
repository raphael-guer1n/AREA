package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/raphael-guer1n/AREA/ServiceService/internal/domain"
)

type userProfileRepository struct {
	db *sql.DB
}

func newUserProfileRepository(db *sql.DB) domain.UserProfileRepository {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) Create(userId int, service, providerUserId, accessToken, refreshToken string, expiresAt time.Time, rawProfile json.RawMessage) (domain.UserProfile, error) {
	var u domain.UserProfile
	err := r.db.QueryRow(
		`INSERT INTO user_service_profiles (user_id, service, provider_user_id,access_token ,refresh_token, expires_at, rawProfile)
         VALUES ($1, $2, $3, $4, $5, $6, $7)
         RETURNING id, email, first_name, last_name, created_at`,
		userId, service, providerUserId, accessToken, refreshToken, expiresAt, rawProfile,
	).Scan(&u.ID, &u.UserId, &u.Service, &u.ProviderUserId, &u.AccessToken, &u.RefreshToken, &u.ExpiresAt, &u.RawProfile)
	return u, err
}
