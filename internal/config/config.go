package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/gallyamow/go-fias-exporter/internal/errors"
)

const (
	ModeBulk   = "bulk"
	ModeUpsert = "upsert"
	ModeSchema = "schema"
	ModeKeys   = "keys"

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
	mode := flag.String("mode", ModeBulk, "mode create|bulk|upsert|keys")
	dbType := flag.String("db-type", DBPostgres, "database type: postgres|mysql")
	dbSchema := flag.String("db-schema", "", "database dbSchema")
	batchSize := flag.Int("batch-size", 1000000, "batch size")
	ignoreRequired := flag.Bool("ignore-required", true, "ignore NOT NULL in CREATE TABLE")
	ignorePrimaryKey := flag.Bool("ignore-primary-key", true, "ignore PRIMARY KEY in CREATE TABLE")

	flag.Parse()

	if flag.NArg() < 1 {
		return nil, errors.ErrPathRequired
	}

	path := flag.Arg(0)
	if path == "" {
		return nil, errors.ErrPathRequired
	}

	// совместимость со старой версией
	if *mode == "copy" {
		*mode = ModeBulk
	}

	if *mode != ModeBulk && *mode != ModeUpsert && *mode != ModeSchema && *mode != ModeKeys {
		return nil, errors.ErrInvalidMode
	}

	if *dbType != DBPostgres && *dbType != DBMySQL {
		return nil, errors.ErrInvalidDBType
	}

	if *batchSize <= 0 {
		return nil, errors.ErrInvalidBatchSize
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
