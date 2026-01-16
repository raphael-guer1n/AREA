package main

import (
	"log"
	"net/http"

	"github.com/raphael-guer1n/AREA/CronService/internal/config"
	"github.com/raphael-guer1n/AREA/CronService/internal/db"
	httphandler "github.com/raphael-guer1n/AREA/CronService/internal/http"
	"github.com/raphael-guer1n/AREA/CronService/internal/repository"
	"github.com/raphael-guer1n/AREA/CronService/internal/service"
)

func main() {
	cfg := config.Load()
	conn := db.Connect(cfg)
	defer conn.Close()

	repo := repository.NewActionRepository(conn)
	cronService := service.NewCronService(repo, cfg.AreaServiceURL, cfg.InternalSecret)

	cronService.Start()
	defer cronService.Stop()

	actionHandler := httphandler.NewActionHandler(cronService)
	router := httphandler.NewRouter(actionHandler, cfg.LogAllRequests)

	addr := ":" + cfg.HTTPPort
	log.Printf("Starting Cron Service on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
