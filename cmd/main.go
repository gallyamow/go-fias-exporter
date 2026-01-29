package main

import (
	"context"
	"fmt"
	"github.com/gallyamow/go-fias-exporter/internal/config"
	"github.com/gallyamow/go-fias-exporter/internal/itemiterator"
	"github.com/gallyamow/go-fias-exporter/internal/sqlbuilder"
	"github.com/gallyamow/go-fias-exporter/pkg/filescanner"
	"io"
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
		_, _ = fmt.Fprintln(os.Stderr, "Failed parse config:", err)
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(os.Stderr, "Version: %s\n", version)
	_, _ = fmt.Fprintln(os.Stderr, cfg)

	files, err := filescanner.ScanDir(ctx, cfg.Path, filescanner.Filter{IncludeExts: []string{"xml"}})
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(os.Stderr, "Found files: %d\n", len(files))
	if len(files) == 0 {
		os.Exit(1)
	}

	for _, f := range files {
		if ctx.Err() != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Interrupted")
			return
		}

		fileName := filepath.Base(f.Path)

		tableName, err := sqlbuilder.ResolveTableName(fileName)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to resolve tablename: %v\n", err)
			return
		}

		_, _ = fmt.Fprintf(os.Stdout, "Started: file %q (%d bytes), table %q \n", fileName, f.Size, tableName)

		sqlBuilder := sqlbuilder.New(tableName, []string{"id"}, []string{"id"})

		func() {
			file, err := os.Open(f.Path)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
				return
			}
			defer file.Close()

			totalRows := 0
			it := itemiterator.New(file)

			for {
				items, err := it.Next(ctx, cfg.BatchSize)
				if err != nil && err != io.EOF {
					_, _ = fmt.Fprintf(os.Stderr, "Failed to read file: %v\n", err)
					return
				}

				if len(items) > 0 {
					sql, err := sqlBuilder.Build(items)
					if err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "Failed to build sql: %v\n", err)
						return
					}

					totalRows++
					_, _ = fmt.Fprintf(os.Stderr, "%s;\n", sql)
				}

				if err == io.EOF {
					_, _ = fmt.Fprintf(os.Stdout, "Ended: %q (%d rows)\n", f.Path, totalRows)
					break
				}
			}
		}()
	}
}
