package config

import (
	"errors"
	"os"
	"testing"
)

func TestParseFlags_CopyFromMode(t *testing.T) {
	resetFlags()
	os.Args = []string{
		"cmd",
		"-mode=copy",
		"-batch-size=100",
		"/data/input",
	}

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Path != "/data/input" {
		t.Fatalf("unexpected path: %s", cfg.Path)
	}

	if cfg.Mode != ModeCopy {
		t.Fatalf("unexpected mode: %s", cfg.Mode)
	}

	if cfg.BatchSize != 100 {
		t.Fatalf("unexpected batch size: %d", cfg.BatchSize)
	}
}

func TestParseFlags_UpsertModeWithDB(t *testing.T) {
	resetFlags()
	os.Args = []string{
		"cmd",
		"-mode=upsert",
		"/data/input",
	}

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Mode != ModeUpsert {
		t.Fatalf("unexpected mode: %s", cfg.Mode)
	}
}

func TestParseFlags_PathRequired(t *testing.T) {
	resetFlags()
	os.Args = []string{
		"cmd",
	}

	_, err := ParseFlags()
	if !errors.Is(err, ErrorPathRequired) {
		t.Fatalf("expected ErrorPathRequired, got %v", err)
	}
}

func TestParseFlags_InvalidMode(t *testing.T) {
	resetFlags()
	os.Args = []string{
		"cmd",
		"-mode=unknown",
		"/data/input",
	}

	_, err := ParseFlags()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestParseFlags_InvalidBatchSize(t *testing.T) {
	resetFlags()
	os.Args = []string{
		"cmd",
		"-batch-size=0",
		"/data/input",
	}

	_, err := ParseFlags()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
