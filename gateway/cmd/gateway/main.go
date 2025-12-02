package main

import (
	"fmt"
	"github.com/raphael-guer1n/AREA/area-gateway/internal/config"
	"github.com/raphael-guer1n/AREA/area-gateway/internal/registry"
)

func main() {
	services, err := config.LoadAllServiceConfigs("./services-config")
	if err != nil {
		panic(err)
	}

	if err := config.ValidateAll(services); err != nil {
		panic(err)
	}

	fmt.Println("All configs are valid.")
	reg := registry.NewRegistry()
	if err := reg.Load(services); err != nil {
		panic(err)
	}

	fmt.Printf("Loaded %d services\n", len(reg.ListServices()))
	fmt.Printf("Loaded %d routes\n", len(reg.ListAllRoutes()))

}
