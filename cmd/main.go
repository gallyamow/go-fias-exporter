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
	fmt.Printf("-- %s", cfg)
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
			log.Fatalf("Failed to resolve tablename: %v", err)
			return
		}

		fmt.Printf("-- Started: file %q (%d bytes), table %q\n", fileName, fileInfo.Size, tableName)
		handleFile(ctx, cfg, tableName, fileInfo.Path)
	}
}

func handleFile(ctx context.Context, cfg *config.Config, tableName string, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Panicf("Failed to open file: %v", err)
		return
	}
	defer file.Close()

	iterator := itemiterator.New(file)
	var sqlBuilder sqlbuilder.Builder
	totalRows := 0

	for {
		items, err := iterator.Next(ctx, cfg.BatchSize)
		if err != nil && err != io.EOF {
			if errors.Is(err, context.Canceled) {
				fmt.Println("-- Canceled")
				return
			}

			log.Panicf("Failed to read file: %v", err)
			return
		}

		if sqlBuilder == nil {
			sqlBuilder = resolveSQLBuilder(cfg, tableName, items[0])
		}

		if len(items) > 0 {
			sql, err := sqlBuilder.Build(items)
			if err != nil {
				log.Panicf("Failed to build sql: %v", err)
				return
			}

			fmt.Println(sql)
			totalRows++
		}

		if err == io.EOF {
			fmt.Printf("-- Ended: %q (%d rows)\n", filePath, totalRows)
			fmt.Println()
			break
		}
	}
}

func resolveSQLBuilder(cfg *config.Config, tableName string, item map[string]string) sqlbuilder.Builder {
	primaryKey := sqlbuilder.ResolvePrimaryKey(tableName, item)
	attrs := sqlbuilder.ResolveAttrs(item)

	switch cfg.Mode {
	case config.ModeCopy:
		return sqlbuilder.NewCopyBuilder(cfg.Schema, tableName, primaryKey, attrs)
	case config.ModeUpsert:
		return sqlbuilder.NewUpsertBuilder(cfg.Schema, tableName, primaryKey, attrs)
	default:
		panic("Failed to resolve builder")
	}
}
