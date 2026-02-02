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
	ModeSchema = "schema"
)

type Config struct {
	Path      string
	Mode      string
	DbSchema  string
	BatchSize int
}

func (c *Config) String() string {
	return fmt.Sprintf("Path: %s, Mode: %s, DbSchema: %s, BatchSize: %d",
		c.Path,
		c.Mode,
		c.DbSchema,
		c.BatchSize,
	)
}

func ParseFlags() (*Config, error) {
	mode := flag.String("mode", ModeCopy, "mode create|copy|upsert")
	dbSchema := flag.String("db-schema", "", "database dbSchema")
	batchSize := flag.Int("batch-size", 1000000, "batch size")

	flag.Parse()

	if flag.NArg() < 1 {
		return nil, ErrorPathRequired
	}

	path := flag.Arg(0)
	if path == "" {
		return nil, ErrorPathRequired
	}

	if *mode != ModeCopy && *mode != ModeUpsert && *mode != ModeSchema {
		return nil, fmt.Errorf("invalid mode '%s'", *mode)
	}

	if *batchSize <= 0 {
		return nil, fmt.Errorf("batch-size must be > 0")
	}

	return &Config{
		Path:      path,
		Mode:      *mode,
		DbSchema:  *dbSchema,
		BatchSize: *batchSize,
	}, nil
}

// testing purposes
func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}
