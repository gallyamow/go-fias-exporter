package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gallyamow/go-fias-exporter/internal/app"
	"github.com/gallyamow/go-fias-exporter/internal/config"
)

var Version = "unknown"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.ParseFlags()

	printHeader(cfg)

	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	application := app.New(cfg)
	if err := application.Run(ctx); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func printHeader(cfg *config.Config) {
	fmt.Println("-- >>>")
	fmt.Printf("-- Version: %s\n", Version)
	if cfg != nil {
		fmt.Printf("-- %s\n", cfg)
	}
	fmt.Printf("-- Started at: %s\n", time.Now())
	fmt.Println("-- <<<")
	fmt.Println()
}
