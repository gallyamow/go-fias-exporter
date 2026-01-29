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

type Config struct {
	Path      string
	Database  string
	BatchSize int
	Delta     int
	Replace   bool
}

func (c *Config) String() string {
	return fmt.Sprintf("Path: %s, Database: %s, BatchSize: %d, Delta: %v, Replace: %v",
		c.Path,
		c.Database,
		c.BatchSize,
		c.Delta,
		c.Replace,
	)
}

func ParseFlags() (*Config, error) {
	defaultBatch := runtime.NumCPU() * 10

	dbFlag := flag.String("db", "duplicates.db", "database path")
	batchFlag := flag.Int("batch-size", defaultBatch, "batch size")
	deltaFlag := flag.Int("delta", 0, "delta value")
	replaceFlag := flag.Bool("replace", false, "replace duplicate files")

	flag.Parse()

	if flag.NArg() < 1 {
		return nil, ErrorPathRequired
	}

	path := flag.Arg(0)
	if path == "" {
		return nil, ErrorPathRequired
	}

	if *batchFlag <= 0 {
		return nil, fmt.Errorf("batch-size must be > 0")
	}

	return &Config{
		Path:      path,
		Database:  *dbFlag,
		BatchSize: *batchFlag,
		Delta:     *deltaFlag,
		Replace:   *replaceFlag,
	}, nil
}

// testing purposes
func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}
