package main

import (
	"context"
	"fmt"
	"github.com/gallyamow/go-fias-exporter/internal/config"
	"github.com/gallyamow/go-fias-exporter/pkg/filescanner"
	"os"
	"os/signal"
	"syscall"
)

var version = "unknown"

func main() {
	_, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed parse config:", err)
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(os.Stderr, "Version: %s\n", version)
	_, _ = fmt.Fprintln(os.Stderr, cfg)

	files, err := filescanner.ScanDir(cfg.Path, filescanner.Filter{IncludeExts: []string{"xml"}})
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(os.Stderr, "Found total files: %d\n", len(files))
}
