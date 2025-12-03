package main

import (
	"log"
	"net/http"

	"github.com/raphael-guer1n/AREA/Template/internal/config"
	"github.com/raphael-guer1n/AREA/Template/internal/db"
	httphandler "github.com/raphael-guer1n/AREA/Template/internal/http"
	"github.com/raphael-guer1n/AREA/Template/internal/repository"
	"github.com/raphael-guer1n/AREA/Template/internal/service"
)

func main() {
	cfg := config.Load()
	dbConn := db.Connect(cfg)

	userRepo := repository.NewUserRepository(dbConn)
	userSvc := service.NewUserService(userRepo)
	router := httphandler.NewRouter(userSvc)

	addr := ":" + cfg.HTTPPort
	log.Printf("Starting server on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
