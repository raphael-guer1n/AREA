package main

import (
	"fmt"
	"github.com/raphael-guer1n/AREA/area-gateway/internal/config"
	"github.com/raphael-guer1n/AREA/area-gateway/internal/registry"
	"github.com/raphael-guer1n/AREA/area-gateway/internal/core"
	"net/http/httptest"
)

func main() {
	//service loader
	services, err := config.LoadAllServiceConfigs("./services-config")
	if err != nil {
		panic(err)
	}

	// service validator
	if err := config.ValidateAll(services); err != nil {
		panic(err)
	}

	fmt.Println("All configs are valid.")
	reg := registry.NewRegistry()
	if err := reg.Load(services); err != nil {
		panic(err)
	}

	// service registry
	fmt.Printf("Loaded %d services\n", len(reg.ListServices()))
	fmt.Printf("Loaded %d routes\n", len(reg.ListAllRoutes()))

	// gateway env loader
	config.LoadDotEnv("configs/gateway.env")
	cfg, err := config.LoadGatewayConfig()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Gateway starting on port %d (debug=%v)\n", cfg.Port, cfg.DebugMode)

	// graceful https erros
	w := httptest.NewRecorder()
	core.WriteError(w, 401, core.ErrUnauthorized, "Missing token")

	fmt.Println("Status code :", w.Code)
	fmt.Println("Headers :", w.Header())
	fmt.Println("Body :", w.Body.String())

}
