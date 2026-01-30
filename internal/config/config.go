package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	ErrorPathRequired = errors.New("path is required")
)

const (
	ModeCopy   = "copy"
	ModeUpsert = "upsert"
)

type Config struct {
	Path      string
	Mode      string
	Schema    string
	BatchSize int
}

func (c *Config) String() string {
	return fmt.Sprintf("Path: %s, Mode: %s, Schema: %s, BatchSize: %d",
		c.Path,
		c.Mode,
		c.Schema,
		c.BatchSize,
	)
}

func ParseFlags() (*Config, error) {
	mode := flag.String("mode", ModeCopy, "mode copy|upsert")
	schema := flag.String("schema", "", "database schema")
	batchSize := flag.Int("batch-size", 1000, "batch size")

	flag.Parse()

	if flag.NArg() < 1 {
		return nil, ErrorPathRequired
	}

	path := flag.Arg(0)
	if path == "" {
		return nil, ErrorPathRequired
	}

	if *mode != ModeCopy && *mode != ModeUpsert {
		return nil, fmt.Errorf("invalid mode '%s'", *mode)
	}

	if *batchSize <= 0 {
		return nil, fmt.Errorf("batch-size must be > 0")
	}

	return &Config{
		Path:      path,
		Mode:      *mode,
		Schema:    *schema,
		BatchSize: *batchSize,
	}, nil
}

// testing purposes
func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}
