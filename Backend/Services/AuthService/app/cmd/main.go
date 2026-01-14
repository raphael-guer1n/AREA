package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/raphael-guer1n/AREA/AuthService/internal/config"
	"github.com/raphael-guer1n/AREA/AuthService/internal/db"
	httphandler "github.com/raphael-guer1n/AREA/AuthService/internal/http"
	"github.com/raphael-guer1n/AREA/AuthService/internal/oauth2"
	"github.com/raphael-guer1n/AREA/AuthService/internal/repository"
	"github.com/raphael-guer1n/AREA/AuthService/internal/service"
)

func main() {
	cfg := config.Load()
	dbConn := db.Connect(cfg)

	// Build repositories
	userProfileRepo := repository.NewUserProfileRepository(dbConn)
	userFieldRepo := repository.NewUserServiceFieldRepository(dbConn)
	userRepo := repository.NewUserRepository(dbConn)

	// Build services
	oauth2StorageSvc := service.NewOAuth2StorageService(userProfileRepo, userFieldRepo, cfg.ServiceServiceURL, cfg.InternalSecret)
	authSvc := service.NewAuthService(userRepo)

	// Initialize OAuth2 manager with service-service URL (lazy loading)
	oauth2Manager := oauth2.NewManager(cfg.ServiceServiceURL, cfg.InternalSecret)
	log.Printf("OAuth2 manager initialized (providers will be loaded on-demand from service-service)")

	refreshWorker := service.NewOAuth2RefreshWorker(
		userProfileRepo,
		oauth2Manager,
		time.Duration(cfg.OAuth2RefreshIntervalSeconds)*time.Second,
		time.Duration(cfg.OAuth2RefreshLeewayMinutes)*time.Minute,
	)
	go refreshWorker.Start(context.Background())

	// Build handlers
	oauth2Handler := httphandler.NewOAuth2Handler(oauth2StorageSvc, oauth2Manager, authSvc, cfg)
	authHandler := httphandler.NewAuthHandler(authSvc)

	// Build router
	router := httphandler.NewRouter(authHandler, oauth2Handler)

	addr := ":" + cfg.HTTPPort
	log.Printf("Starting server on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
