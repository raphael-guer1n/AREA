package main

import (
	"log"
	"net/http"

	"github.com/raphael-guer1n/AREA/MailService/internal/config"
	httphandler "github.com/raphael-guer1n/AREA/MailService/internal/http"
	"github.com/raphael-guer1n/AREA/MailService/internal/service"
)

func main() {
	cfg := config.Load()
	mailer := service.NewMailer(cfg)
	router := httphandler.NewRouter(mailer, cfg)

	addr := ":" + cfg.HTTPPort
	log.Printf("Starting server on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
