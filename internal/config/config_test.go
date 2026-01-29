package config

import (
	"errors"
	"os"
	"testing"
)

func TestParseFlags_OutputMode(t *testing.T) {
	resetFlags()
	os.Args = []string{
		"cmd",
		"-mode=output",
		"-batch-size=100",
		"-delta=5",
		"/data/input",
	}

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Path != "/data/input" {
		t.Fatalf("unexpected path: %s", cfg.Path)
	}

	if cfg.Mode != ModeOutput {
		t.Fatalf("unexpected mode: %s", cfg.Mode)
	}

	if cfg.BatchSize != 100 {
		t.Fatalf("unexpected batch size: %d", cfg.BatchSize)
	}

	if cfg.Delta != 5 {
		t.Fatalf("unexpected delta: %d", cfg.Delta)
	}
}

func TestParseFlags_ExecuteModeRequiresDB(t *testing.T) {
	resetFlags()
	os.Args = []string{
		"cmd",
		"-mode=execute",
		"/data/input",
	}

	_, err := ParseFlags()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestParseFlags_ExecuteModeWithDB(t *testing.T) {
	resetFlags()
	os.Args = []string{
		"cmd",
		"-mode=execute",
		"-db=postgres://localhost/db",
		"/data/input",
	}

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Mode != ModeExecute {
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

func TestParseFlags_InvalidDelta(t *testing.T) {
	resetFlags()
	os.Args = []string{
		"cmd",
		"-delta=-1",
		"/data/input",
	}

	_, err := ParseFlags()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
