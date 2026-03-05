package config

import (
	"errors"
	"os"
	"testing"

	apperrors "github.com/gallyamow/go-fias-exporter/internal/errors"
)

func TestParseFlags_Mode(t *testing.T) {
	resetFlags()
	os.Args = []string{
		"cmd",
		"-mode=keys",
		"/data/input",
	}

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Mode != ModeKeys {
		t.Fatalf("unexpected mode: %s", cfg.Mode)
	}
}

func TestParseFlags_BulkMode(t *testing.T) {
	resetFlags()
	os.Args = []string{
		"cmd",
		"-mode=bulk",
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

	if cfg.Mode != ModeBulk {
		t.Fatalf("unexpected mode: %s", cfg.Mode)
	}
}

func TestParseFlags_DeprecatedBulkMode(t *testing.T) {
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

	if cfg.Mode != ModeBulk {
		t.Fatalf("unexpected mode: %s", cfg.Mode)
	}
}

func TestParseFlags_PathRequired(t *testing.T) {
	resetFlags()
	os.Args = []string{
		"cmd",
	}

	_, err := ParseFlags()
	if !errors.Is(err, apperrors.ErrPathRequired) {
		t.Fatalf("expected ErrPathRequired, got %v", err)
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

func TestParseFlags_BatchSize(t *testing.T) {
	resetFlags()
	os.Args = []string{
		"cmd",
		"-mode=bulk",
		"-batch-size=1100",
		"/data/input",
	}

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.BatchSize != 1100 {
		t.Fatalf("unexpected batch size: %d", cfg.BatchSize)
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
