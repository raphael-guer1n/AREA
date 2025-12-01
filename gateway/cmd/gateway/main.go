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

	fmt.Println("Loaded services:")
	for _, s := range services {
		fmt.Printf(" - %s (%s)\n", s.Name, s.BaseURL)
	}
}
