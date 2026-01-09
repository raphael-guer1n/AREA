package main

import (
	"log"
	"net/http"

	"github.com/raphael-guer1n/AREA/AreaService/internal/config"
	"github.com/raphael-guer1n/AREA/AreaService/internal/db"
	httphandler "github.com/raphael-guer1n/AREA/AreaService/internal/http"
	"github.com/raphael-guer1n/AREA/AreaService/internal/repository"
	"github.com/raphael-guer1n/AREA/AreaService/internal/service"
)

func main() {
	cfg := config.Load()
	dbConn := db.Connect(cfg)

	areaRepository := repository.NewAreaRepository(dbConn)

	areaSvc := service.NewAreaService(areaRepository)
	areaHandler := httphandler.NewAreaHandler(areaSvc)
	router := httphandler.NewRouter(areaHandler)

	addr := ":" + cfg.HTTPPort
	log.Printf("Starting server on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
