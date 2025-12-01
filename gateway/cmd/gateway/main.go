package main

import (
	"fmt"
	"github.com/raphael-guer1n/AREA/area-gateway/internal/config"
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
}
