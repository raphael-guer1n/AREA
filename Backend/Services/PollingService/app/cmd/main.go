package main

import (
	"log"
	"net/http"

	"github.com/raphael-guer1n/AREA/PollingService/internal/config"
	"github.com/raphael-guer1n/AREA/PollingService/internal/db"
	httphandler "github.com/raphael-guer1n/AREA/PollingService/internal/http"
	"github.com/raphael-guer1n/AREA/PollingService/internal/repository"
	"github.com/raphael-guer1n/AREA/PollingService/internal/service"
)

func main() {
	cfg := config.Load()
	conn := db.Connect(cfg)

	repo := repository.NewSubscriptionRepository(conn)
	providerConfigSvc := service.NewProviderConfigService(cfg.ServiceServiceURL, cfg.InternalSecret)
	oauth2TokenSvc := service.NewOAuth2TokenService(cfg.AuthServiceURL, cfg.InternalSecret)
	requestSvc := service.NewRequestService(oauth2TokenSvc, cfg.LogProviderRequests)
	subscriptionSvc := service.NewSubscriptionService(repo, providerConfigSvc, requestSvc)
	authSvc := service.NewAuthService(cfg.AuthServiceURL)
	areaTriggerSvc := service.NewAreaTriggerService(cfg.AreaServiceURL, cfg.InternalSecret)
	pollingWorker := service.NewPollingWorker(repo, providerConfigSvc, requestSvc, areaTriggerSvc, cfg.PollingTickSeconds)
	go pollingWorker.Start()

	actionHandler := httphandler.NewActionHandler(subscriptionSvc, authSvc)
	router := httphandler.NewRouter(actionHandler)

	addr := ":" + cfg.HTTPPort
	log.Printf("Starting server on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
