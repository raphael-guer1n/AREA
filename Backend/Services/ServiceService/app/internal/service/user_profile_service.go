package service

import (
	"encoding/json"
	"time"

	"github.com/raphael-guer1n/AREA/ServiceService/internal/domain"
)

type UserProfileService struct {
	repo domain.UserProfileRepository
}

func NewUserProfileService(repo domain.UserProfileRepository) *UserProfileService {
	return &UserProfileService{repo: repo}
}

func (s *UserProfileService) Create(userId int, service, providerUserId, accessToken, refreshToken string, expiresAt time.Time, rawProfile json.RawMessage) (domain.UserProfile, error) {
	return s.repo.Create(userId, service, providerUserId, accessToken, refreshToken, expiresAt, rawProfile)
}
