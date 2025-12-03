package main

import (
	"log"
	"net/http"

	"github.com/raphael-guer1n/AREA/AuthService/internal/config"
	"github.com/raphael-guer1n/AREA/AuthService/internal/db"
	httphandler "github.com/raphael-guer1n/AREA/AuthService/internal/http"
	"github.com/raphael-guer1n/AREA/AuthService/internal/repository"
	"github.com/raphael-guer1n/AREA/AuthService/internal/service"
)

func main() {
	cfg := config.Load()
	dbConn := db.Connect(cfg)

	userRepo := repository.NewUserRepository(dbConn)
	authSvc := service.NewAuthService(userRepo)
	router := httphandler.NewRouter(authSvc)

	addr := ":" + cfg.HTTPPort
	log.Printf("Starting server on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
