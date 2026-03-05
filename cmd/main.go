package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gallyamow/go-fias-exporter/internal/app"
	"github.com/gallyamow/go-fias-exporter/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.ParseFlags()
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	application := app.New(cfg)
	if err := application.Run(ctx); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
