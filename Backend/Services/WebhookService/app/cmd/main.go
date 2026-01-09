package main

import (
	"log"
	"net/http"

	"github.com/raphael-guer1n/AREA/WebhookService/internal/config"
	"github.com/raphael-guer1n/AREA/WebhookService/internal/db"
	httphandler "github.com/raphael-guer1n/AREA/WebhookService/internal/http"
	"github.com/raphael-guer1n/AREA/WebhookService/internal/repository"
	"github.com/raphael-guer1n/AREA/WebhookService/internal/service"
)

func main() {
	cfg := config.Load()
	conn := db.Connect(cfg)

	repo := repository.NewSubscriptionRepository(conn)
	providerConfigSvc := service.NewProviderConfigService(cfg.ServiceServiceURL)
	oauth2TokenSvc := service.NewOAuth2TokenService(cfg.AuthServiceURL)
	webhookSetupSvc := service.NewWebhookSetupService(oauth2TokenSvc)
	subscriptionSvc := service.NewSubscriptionService(repo, providerConfigSvc, webhookSetupSvc)
	renewalSvc := service.NewSubscriptionRenewalService(repo, providerConfigSvc, webhookSetupSvc, cfg.PublicBaseURL)
	go renewalSvc.Start()

	subscriptionHandler := httphandler.NewSubscriptionHandler(subscriptionSvc, cfg)
	webhookHandler := httphandler.NewWebhookHandler(subscriptionSvc, providerConfigSvc)

	router := httphandler.NewRouter(subscriptionHandler, webhookHandler)

	addr := ":" + cfg.HTTPPort
	log.Printf("Starting server on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
