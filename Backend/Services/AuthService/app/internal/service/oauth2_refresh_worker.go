package service

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/raphael-guer1n/AREA/AuthService/internal/domain"
	"github.com/raphael-guer1n/AREA/AuthService/internal/oauth2"
)

type OAuth2RefreshWorker struct {
	profileRepo domain.UserProfileRepository
	manager     *oauth2.Manager
	interval    time.Duration
	leeway      time.Duration
}

func NewOAuth2RefreshWorker(
	profileRepo domain.UserProfileRepository,
	manager *oauth2.Manager,
	interval time.Duration,
	leeway time.Duration,
) *OAuth2RefreshWorker {
	return &OAuth2RefreshWorker{
		profileRepo: profileRepo,
		manager:     manager,
		interval:    interval,
		leeway:      leeway,
	}
}

func (w *OAuth2RefreshWorker) Start(ctx context.Context) {
	if w.interval <= 0 {
		w.interval = time.Minute
	}
	if w.leeway <= 0 {
		w.leeway = 5 * time.Minute
	}

	w.runOnce()

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.runOnce()
		}
	}
}

func (w *OAuth2RefreshWorker) runOnce() {
	cutoff := time.Now().Add(w.leeway)
	candidates, err := w.profileRepo.ListRefreshCandidates(cutoff)
	if err != nil {
		log.Printf("oauth2 refresh: failed to list candidates: %v", err)
		return
	}
	if len(candidates) == 0 {
		return
	}

	for _, candidate := range candidates {
		if strings.TrimSpace(candidate.RefreshToken) == "" {
			w.markReconnect(candidate.ID, "missing refresh token")
			continue
		}

		provider, err := w.manager.GetProvider(candidate.Service)
		if err != nil {
			w.markReconnect(candidate.ID, "failed to load provider config: "+err.Error())
			continue
		}

		tokenResp, err := provider.RefreshToken(candidate.RefreshToken)
		if err != nil {
			w.markReconnect(candidate.ID, err.Error())
			continue
		}

		expiresAt := oauth2.ResolveExpiresAt(tokenResp.ExpiresIn)
		if err := w.profileRepo.UpdateTokens(candidate.ID, tokenResp.AccessToken, tokenResp.RefreshToken, expiresAt); err != nil {
			log.Printf("oauth2 refresh: failed to update tokens for profile %d: %v", candidate.ID, err)
			continue
		}

		log.Printf("oauth2 refresh: refreshed token for service=%s profile=%d", candidate.Service, candidate.ID)
	}
}

func (w *OAuth2RefreshWorker) markReconnect(profileID int, reason string) {
	if err := w.profileRepo.MarkNeedsReconnect(profileID, reason); err != nil {
		log.Printf("oauth2 refresh: failed to mark reconnect for profile %d: %v", profileID, err)
		return
	}
	log.Printf("oauth2 refresh: marked reconnect for profile %d (%s)", profileID, reason)
}
