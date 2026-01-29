package config

import (
	"errors"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestParseFlags(t *testing.T) {
	t.Run("no path", func(t *testing.T) {
		resetFlags()
		os.Args = []string{"cmd"}

		_, err := ParseFlags()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, ErrorPathRequired) {
			t.Fatalf("got %v, want %v", err, ErrorPathRequired)
		}
	})

	t.Run("defaults", func(t *testing.T) {
		resetFlags()
		os.Args = []string{"cmd", "/test"}

		cfg, err := ParseFlags()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Path != "/test" {
			t.Fatalf("got Path=%v, want %v", cfg.Path, "/test")
		}
		if cfg.Database != "duplicates.db" {
			t.Fatalf("got Database=%v, want %v", cfg.Database, "duplicates.db")
		}
		wantBatch := runtime.NumCPU() * 10
		if cfg.BatchSize != wantBatch {
			t.Fatalf("got BatchSize=%v, want %v", cfg.BatchSize, wantBatch)
		}
		if cfg.Delta != 0 {
			t.Fatalf("got Delta=%v, want %v", cfg.Delta, 0)
		}
		if cfg.Replace != false {
			t.Fatalf("got Replace=%v, want %v", cfg.Replace, false)
		}
	})

	t.Run("parse all flags", func(t *testing.T) {
		resetFlags()
		os.Args = []string{
			"cmd",
			"-db=custom.db",
			"-batch-size=5",
			"-delta=3",
			"-replace",
			"/my/path",
		}

		cfg, err := ParseFlags()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Path != "/my/path" {
			t.Fatalf("got Path=%v, want %v", cfg.Path, "/my/path")
		}
		if cfg.Database != "custom.db" {
			t.Fatalf("got Database=%v, want %v", cfg.Database, "custom.db")
		}
		if cfg.BatchSize != 5 {
			t.Fatalf("got BatchSize=%v, want %v", cfg.BatchSize, 5)
		}
		if cfg.Delta != 3 {
			t.Fatalf("got Delta=%v, want %v", cfg.Delta, 3)
		}
		if cfg.Replace != true {
			t.Fatalf("got Replace=%v, want %v", cfg.Replace, true)
		}
	})

	t.Run("invalid batch-size", func(t *testing.T) {
		resetFlags()
		os.Args = []string{"cmd", "-batch-size=0", "/path"}

		_, err := ParseFlags()
		if err == nil {
			t.Fatal("expected error for invalid batch-size, got nil")
		}
		if !strings.Contains(err.Error(), "batch-size must be > 0") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
