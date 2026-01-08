package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/raphael-guer1n/AREA/area-gateway/internal/config"
	"github.com/raphael-guer1n/AREA/area-gateway/internal/middleware"
	"github.com/raphael-guer1n/AREA/area-gateway/internal/registry"
	"github.com/raphael-guer1n/AREA/area-gateway/internal/router"
)

func main() {
	services, err := config.LoadAllServiceConfigs("./services-config")
	if err != nil {
		log.Fatalf("[FATAL] failed to load service configs: %v", err)
	}

	if err := config.ValidateAll(services); err != nil {
		log.Fatalf("[FATAL] invalid service configuration: %v", err)
	}

	log.Println("[INFO] All service configs are valid.")

	reg := registry.NewRegistry()
	if err := reg.Load(services); err != nil {
		log.Fatalf("[FATAL] failed to load registry: %v", err)
	}

	log.Printf("[INFO] Loaded %d services", len(reg.ListServices()))
	log.Printf("[INFO] Loaded %d routes", len(reg.ListAllRoutes()))

	config.LoadDotEnv("configs/gateway.env")

	cfg, err := config.LoadGatewayConfig()
	if err != nil {
		log.Fatalf("[FATAL] failed to load gateway config: %v", err)
	}

	log.Printf("[INFO] Gateway starting on port %d (debug=%v)", cfg.Port, cfg.DebugMode)

	authMW := middleware.NewAuthMiddleware(cfg, reg)
	permMW := middleware.NewPermissionsMiddleware(reg)
	internalMW := middleware.NewInternalMiddleware(cfg, reg)
	loggingMW := middleware.NewLoggingMiddleware(reg)

	rt := router.NewRouter(
		reg,
		cfg,
		authMW,
		permMW,
		internalMW,
		loggingMW,
	)

	mux, err := rt.Build()
	if err != nil {
		log.Fatalf("[FATAL] failed to build router: %v", err)
	}

	corsMW := middleware.NewCORSMiddleware(cfg)
	handler := corsMW.Handler(mux)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Println("[INFO] Listening on", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal("[FATAL] server crashed:", err)
	}
}
