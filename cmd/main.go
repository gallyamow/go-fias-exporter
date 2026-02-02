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

	"github.com/gallyamow/go-fias-exporter/internal/itemiterator"

	"github.com/gallyamow/go-fias-exporter/internal/config"
	"github.com/gallyamow/go-fias-exporter/internal/sqlbuilder"
	"github.com/gallyamow/go-fias-exporter/pkg/filescanner"
)

var version = "unknown"

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

		fmt.Printf("-- Started: file %q (%d bytes), table %q\n", fileName, fileInfo.Size, tableName)

		switch cfg.Mode {
		case config.ModeCopy, config.ModeUpsert:
			totalRows, err := handleDataFile(ctx, cfg, tableName, fileInfo.Path)
			if err != nil && err != io.EOF {
				if errors.Is(err, context.Canceled) {
					fmt.Println("-- Canceled")
				}
				log.Panic(err)
			}
			fmt.Printf("-- %d rows\n", totalRows)
		case config.ModeSchema:
			handleSchemaFile(ctx, cfg, tableName, fileInfo.Path)
		}

		fmt.Printf("-- Ended: %q\n", fileInfo.Path)
		fmt.Println()
	}
}

func handleDataFile(ctx context.Context, cfg *config.Config, tableName string, filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	iterator := itemiterator.New(file)
	var sqlBuilder sqlbuilder.Builder
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
			sqlBuilder = resolveDataBuilder(cfg, tableName, items[0])
		}

		if len(items) > 0 {
			sql, err := sqlBuilder.Build(items)
			if err != nil {
				return totalRows, err
			}

			fmt.Println(sql)
			totalRows++
		}

		if err == io.EOF {
			return totalRows, err
		}
	}
}

func handleSchemaFile(ctx context.Context, cfg *config.Config, tableName string, filePath string) {
	sqlBuilder := sqlbuilder.NewSchemaBuilder(cfg.DbSchema, tableName, cfg.IgnoreNotNull)

	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	sql, err := sqlBuilder.Build(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(sql)
}

func resolveDataBuilder(cfg *config.Config, tableName string, item map[string]string) sqlbuilder.Builder {
	attrs := sqlbuilder.ResolveAttrs(item)

	switch cfg.Mode {
	case config.ModeCopy:
		return sqlbuilder.NewCopyBuilder(cfg.DbSchema, tableName, attrs)
	case config.ModeUpsert:
		return sqlbuilder.NewUpsertBuilder(cfg.DbSchema, tableName, attrs)
	default:
		panic(fmt.Sprintf("failed to resolve builder for %q", cfg.Mode))
	}
}
