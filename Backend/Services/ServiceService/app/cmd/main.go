package main

import (
	"log"
	"net/http"

	"github.com/raphael-guer1n/AREA/ServiceService/internal/config"
	httphandler "github.com/raphael-guer1n/AREA/ServiceService/internal/http"
	"github.com/raphael-guer1n/AREA/ServiceService/internal/service"
)

func main() {
	cfg := config.Load()

	// Provider config service
	providerConfigSvc, err := service.NewProviderConfigService("internal/config/providers")
	if err != nil {
		log.Fatalf("Failed to load provider configs: %v", err)
	}

	// HTTP handlers
	providerHandler := httphandler.NewProviderHandler(providerConfigSvc)

	router := httphandler.NewRouter(providerHandler)

	addr := ":" + cfg.HTTPPort
	log.Printf("Starting server on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
