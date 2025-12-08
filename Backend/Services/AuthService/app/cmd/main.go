package main

import (
	"log"
	"net/http"

	"github.com/raphael-guer1n/AREA/AuthService/internal/config"
	"github.com/raphael-guer1n/AREA/AuthService/internal/db"
	httphandler "github.com/raphael-guer1n/AREA/AuthService/internal/http"
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
	oauth2StorageSvc := service.NewOAuth2StorageService(userProfileRepo, userFieldRepo, cfg.ServiceServiceURL)
	authSvc := service.NewAuthService(userRepo)

	// Load OAuth2 configuration (optional)
	var oauth2Manager *oauth2.Manager
	oauth2Config, err := oauth2.LoadConfig(cfg.OAuth2ConfigPath)
	if err != nil {
		log.Printf("Warning: OAuth2 config not loaded (%v). OAuth2 endpoints will be disabled.", err)
	} else {
		oauth2Manager = oauth2.NewManager(oauth2Config)
		log.Printf("OAuth2 enabled with providers: %v", oauth2Manager.ListProviders())
	}

	// Build handlers
	oauth2Handler := httphandler.NewOAuth2Handler(oauth2StorageSvc)
	authHandler := httphandler.NewAuthHandler(authSvc)

	// Build router
	router := httphandler.NewRouter(authHandler, oauth2Handler)

	addr := ":" + cfg.HTTPPort
	log.Printf("Starting server on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
