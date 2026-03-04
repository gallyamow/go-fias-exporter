package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/gallyamow/go-fias-exporter/internal/itemiterator"

	"github.com/gallyamow/go-fias-exporter/internal/config"
	"github.com/gallyamow/go-fias-exporter/internal/sqlbuilder"
	"github.com/gallyamow/go-fias-exporter/pkg/filescanner"
)

var version = "unknown"

//nolint:gocyclo
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.ParseFlags()
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	fmt.Println("-- >>>")
	fmt.Printf("-- Version: %s\n", version)
	fmt.Printf("-- %s\n", cfg)
	fmt.Printf("-- Started at: %s\n", time.Now())
	fmt.Println("-- <<<")
	fmt.Println()

	files, err := filescanner.ScanDir(ctx, cfg.Path, filescanner.Filter{IncludeExts: []string{"xml"}})
	if err != nil {
		log.Fatalf("Failed to scan dir: %v", err)
	}

	if len(files) == 0 {
		log.Fatalf("No files found")
	}

	for _, fileInfo := range files {
		if ctx.Err() != nil {
			return
		}

		fileName := filepath.Base(fileInfo.Path)

		tableName, err := sqlbuilder.ResolveTableName(fileName)
		if err != nil {
			log.Fatalf("Failed to resolve table name: %v", err)
			return
		}

		startedAt := time.Now()
		_, _ = fmt.Fprintf(os.Stderr, "Started file %q (%s) to table %q at %q.\n", fileInfo.Path, humanize.Bytes(uint64(fileInfo.Size)), tableName, startedAt.Format(time.RFC3339))

		switch cfg.Mode {
		case config.ModeBulk, config.ModeUpsert:
			totalRows, err := handleDataFile(ctx, cfg, tableName, fileInfo.Path)
			if err != nil && err != io.EOF {
				if errors.Is(err, context.Canceled) {
					// to mark an unfinished process
					fmt.Println("-- Canceled")
				}
				log.Fatalf("Failed to handle data file: %v", err)
			}
			_, _ = fmt.Fprintf(os.Stderr, "Handled %d rows\n", totalRows)
		case config.ModeSchema, config.ModeKeys:
			if err = handleSchemaFile(ctx, cfg, tableName, fileInfo.Path); err != nil {
				log.Fatalf("Failed to handle schema file: %v", err)
			}
		}

		_, _ = fmt.Fprintf(os.Stderr, "Ended after %s.\n", time.Since(startedAt).Round(time.Second))
	}
}

//nolint:gocyclo
func handleDataFile(ctx context.Context, cfg *config.Config, tableName string, filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	iterator := itemiterator.New(file)
	var sqlBuilder sqlbuilder.ImportBuilder
	totalRows := 0

	for {
		items, err := iterator.Next(ctx, cfg.BatchSize)
		if err != nil && err != io.EOF {
			if errors.Is(err, context.Canceled) {
				return totalRows, err
			}
		}

		// <ITEMS> can be empty, wait EOF
		if len(items) == 0 {
			return totalRows, nil
		}

		if sqlBuilder == nil {
			sqlBuilder = resolveImportBuilder(cfg, tableName, items[0])
		}

		if len(items) > 0 {
			sql, err := sqlBuilder.Build(items)
			if err != nil {
				return totalRows, err
			}

			fmt.Println(sql)
			totalRows += len(items)
		}

		if err == io.EOF {
			return totalRows, err
		}
	}
}

func handleSchemaFile(ctx context.Context, cfg *config.Config, tableName string, filePath string) error {
	var sqlBuilder sqlbuilder.SchemaBuilder

	switch cfg.DbType {
	case config.DBPostgres:
		sqlBuilder = sqlbuilder.NewPostgreSQLSchemaBuilder(cfg.DbSchema, tableName, cfg.IgnoreRequired, cfg.IgnorePrimaryKey)
	case config.DBMySQL:
		sqlBuilder = sqlbuilder.NewMySQLSchemaBuilder(cfg.DbSchema, tableName, cfg.IgnoreRequired, cfg.IgnorePrimaryKey)
	default:
		return fmt.Errorf("unsupported database type: %s", cfg.DbType)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var sql string

	switch cfg.Mode {
	case config.ModeKeys:
		sql, err = sqlBuilder.BuildPrimaryKey()
	default:
		sql, err = sqlBuilder.Build(data)
	}

	if err != nil {
		return err
	}

	fmt.Println(sql)
	return nil
}

func resolveImportBuilder(cfg *config.Config, tableName string, item map[string]string) sqlbuilder.ImportBuilder {
	attrs := sqlbuilder.ResolveAttrs(item)

	switch cfg.Mode {
	case config.ModeBulk:
		switch cfg.DbType {
		case config.DBPostgres:
			return sqlbuilder.NewPostgreSQLCopyBuilder(cfg.DbSchema, tableName, attrs)
		case config.DBMySQL:
			return sqlbuilder.NewMySQLLoadDataBuilder(cfg.DbSchema, tableName, attrs)
		}
	case config.ModeUpsert:
		switch cfg.DbType {
		case config.DBPostgres:
			return sqlbuilder.NewPostgreSQLUpsertBuilder(cfg.DbSchema, tableName, attrs)
		case config.DBMySQL:
			return sqlbuilder.NewMySQLUpsertBuilder(cfg.DbSchema, tableName, attrs)
		}
	}
	panic(fmt.Sprintf("failed to resolve import builder for %q", cfg.Mode))
}
