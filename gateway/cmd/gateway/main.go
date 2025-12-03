package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/raphael-guer1n/AREA/area-gateway/internal/config"
	"github.com/raphael-guer1n/AREA/area-gateway/internal/core"
	"github.com/raphael-guer1n/AREA/area-gateway/internal/middleware"
	"github.com/raphael-guer1n/AREA/area-gateway/internal/registry"
	"github.com/raphael-guer1n/AREA/area-gateway/internal/router"
)

func main() {
	// Load all service configs

	services, err := config.LoadAllServiceConfigs("./services-config")
	if err != nil {
		panic(err)
	}

	if err := config.ValidateAll(services); err != nil {
		panic(err)
	}

	fmt.Println("All configs are valid.")

	// Load registry

	reg := registry.NewRegistry()
	if err := reg.Load(services); err != nil {
		panic(err)
	}

	fmt.Printf("Loaded %d services\n", len(reg.ListServices()))
	fmt.Printf("Loaded %d routes\n", len(reg.ListAllRoutes()))

	// Load Gateway config

	config.LoadDotEnv("configs/gateway.env")

	cfg, err := config.LoadGatewayConfig()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Gateway starting on port %d (debug=%v)\n", cfg.Port, cfg.DebugMode)

	// Build all middlewares

	authMW := middleware.NewAuthMiddleware(cfg, reg)
	permMW := middleware.NewPermissionsMiddleware(reg)
	internalMW := middleware.NewInternalMiddleware(cfg, reg)
	loggingMW := middleware.NewLoggingMiddleware(reg)

	// Build dynamic router

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
		panic(err)
	}

	// Standalone JWT testing endpoint

	mux.Handle("/secure", authMW.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := middleware.GetUserFromContext(r.Context())
		if user == nil {
			core.WriteError(w, http.StatusUnauthorized, core.ErrUnauthorized, "Missing user in context")
			return
		}
		w.Write([]byte("Hello " + user.Email))
	})))

	// Start HTTP server

	addr := fmt.Sprintf(":%d", cfg.Port)
	fmt.Println("Listening on", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("Server crashed:", err)
	}
}
