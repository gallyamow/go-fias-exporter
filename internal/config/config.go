package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
)

var (
	ErrorPathRequired = errors.New("path is required")
)

const (
	ModeOutput  = "output"
	ModeExecute = "execute"
)

type Config struct {
	Path      string
	Mode      string
	BatchSize int
	Delta     int
}

func (c *Config) String() string {
	return fmt.Sprintf("Path: %s, Mode: %s, BatchSize: %d, Delta: %v",
		c.Path,
		c.Mode,
		c.BatchSize,
		c.Delta,
	)
}

func ParseFlags() (*Config, error) {
	defaultBatch := runtime.NumCPU() * 10

	mode := flag.String("mode", ModeOutput, "mode output|execute")
	batchSize := flag.Int("batch-size", defaultBatch, "batch size")
	deltaKey := flag.Int("delta", 0, "delta key")

	flag.Parse()

	if flag.NArg() < 1 {
		return nil, ErrorPathRequired
	}

	path := flag.Arg(0)
	if path == "" {
		return nil, ErrorPathRequired
	}

	if *mode != ModeOutput && *mode != ModeExecute {
		return nil, fmt.Errorf("invalid mode")
	}

	if *batchSize <= 0 {
		return nil, fmt.Errorf("batch-size must be > 0")
	}

	if *deltaKey <= 0 {
		return nil, fmt.Errorf("deltaKey must be > 0")
	}

	return &Config{
		Path:      path,
		Mode:      *mode,
		BatchSize: *batchSize,
		Delta:     *deltaKey,
	}, nil
}

// testing purposes
func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}
