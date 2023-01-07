package main

import (
	"bot/config"
	"bot/internal/app"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Info("Initializing bot...")

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	an, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Cannot initialize sBot: %s", err)
	}

	configureCloseSignal(an)

	an.Start()
}

func configureCloseSignal(an *app.An) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		an.Shutdown()
		os.Exit(0)
	}()
}
