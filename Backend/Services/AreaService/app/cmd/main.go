package main

import (
	"log"
	"net/http"

	"github.com/raphael-guer1n/AREA/AreaService/internal/config"
	httphandler "github.com/raphael-guer1n/AREA/AreaService/internal/http"
	"github.com/raphael-guer1n/AREA/AreaService/internal/service"
)

func main() {
	cfg := config.Load()

	areaSvc := service.NewAreaService()
	areaHandler := httphandler.NewAreaHandler(areaSvc)
	router := httphandler.NewRouter(areaHandler)

	addr := ":" + cfg.HTTPPort
	log.Printf("Starting server on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
