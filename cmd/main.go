package main

import (
	"context"
	"fmt"
	"github.com/gallyamow/go-fias-exporter/internal/config"
	"github.com/gallyamow/go-fias-exporter/internal/itemiterator"
	"github.com/gallyamow/go-fias-exporter/internal/sqlbuilder"
	"github.com/gallyamow/go-fias-exporter/pkg/filescanner"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var version = "unknown"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.ParseFlags()
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	//_, _ = fmt.Fprintf(os.Stderr, "Version: %s\n", version)
	//_, _ = fmt.Fprintln(os.Stderr, cfg)

	files, err := filescanner.ScanDir(ctx, cfg.Path, filescanner.Filter{IncludeExts: []string{"xml"}})
	if err != nil {
		log.Fatalf("Failed to scan dir: %v", err)
	}

	_, _ = fmt.Fprintf(os.Stderr, "Found files: %d\n", len(files))
	if len(files) == 0 {
		log.Fatalf("No files found")
	}

	for _, f := range files {
		if ctx.Err() != nil {
			return
		}

		fileName := filepath.Base(f.Path)

		tableName, err := sqlbuilder.ResolveTableName(fileName)
		if err != nil {
			log.Fatalf("Failed to resolve tablename: %v", err)
			return
		}

		fmt.Printf("-- Started: file %q (%d bytes), table %q\n", fileName, f.Size, tableName)

		func() {
			file, err := os.Open(f.Path)
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
					log.Panicf("Failed to read file: %v", err)
					return
				}

				if sqlBuilder == nil {
					primaryKey := sqlbuilder.ResolvePrimaryKey(tableName, items[0])
					attrs := sqlbuilder.ResolveAttrs(items[0])

					switch cfg.Mode {
					case config.ModeCopyFrom:
						sqlBuilder = sqlbuilder.NewCopyBuilder(tableName, primaryKey, attrs)
					case config.ModeUpsert:
						sqlBuilder = sqlbuilder.NewUpsertBuilder(tableName, primaryKey, attrs)
					default:
						log.Panicf("Failed to resolve builder")
						return
					}
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
					fmt.Printf("-- Ended: %q (%d rows)\n", f.Path, totalRows)
					fmt.Println()
					break
				}
			}
		}()
	}
}
