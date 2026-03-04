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

	DBPostgres = "postgres"
	DBMySQL    = "mysql"
)

type Config struct {
	Path             string
	Mode             string
	DbType           string
	DbSchema         string
	BatchSize        int
	IgnoreRequired   bool
	IgnorePrimaryKey bool
}

func (c *Config) String() string {
	return fmt.Sprintf("Path: %s, Mode: %s, DbType: %s, DbSchema: %s, BatchSize: %d, IgnoreRequired: %v, IgnorePrimaryKey: %v",
		c.Path,
		c.Mode,
		c.DbType,
		c.DbSchema,
		c.BatchSize,
		c.IgnoreRequired,
		c.IgnorePrimaryKey,
	)
}

func ParseFlags() (*Config, error) {
	mode := flag.String("mode", ModeCopy, "mode create|copy|upsert")
	dbType := flag.String("db-type", DBPostgres, "database type: postgres|mysql")
	dbSchema := flag.String("db-schema", "", "database dbSchema")
	batchSize := flag.Int("batch-size", 1000000, "batch size")
	ignoreRequired := flag.Bool("ignore-required", true, "ignore NOT NULL in CREATE TABLE")
	ignorePrimaryKey := flag.Bool("ignore-primary-key", true, "ignore PRIMARY KEY in CREATE TABLE")

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

	if *dbType != DBPostgres && *dbType != DBMySQL {
		return nil, fmt.Errorf("invalid db-type '%s'", *dbType)
	}

	if *batchSize <= 0 {
		return nil, fmt.Errorf("batch-size must be > 0")
	}

	return &Config{
		Path:             path,
		Mode:             *mode,
		DbType:           *dbType,
		DbSchema:         *dbSchema,
		BatchSize:        *batchSize,
		IgnoreRequired:   *ignoreRequired,
		IgnorePrimaryKey: *ignorePrimaryKey,
	}, nil
}

// testing purposes
func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}
