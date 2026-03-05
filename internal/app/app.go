package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/gallyamow/go-fias-exporter/internal/config"
	"github.com/gallyamow/go-fias-exporter/internal/errors"
	"github.com/gallyamow/go-fias-exporter/internal/itemiterator"
	"github.com/gallyamow/go-fias-exporter/internal/model"
	"github.com/gallyamow/go-fias-exporter/internal/sqlbuilder"
	"github.com/gallyamow/go-fias-exporter/pkg/filescanner"
)

type Application struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Application {
	return &Application{cfg: cfg}
}

func (a *Application) Run(ctx context.Context) error {
	a.printHeader()

	var ext string

	switch a.cfg.Mode {
	case config.ModeBulk, config.ModeUpsert:
		ext = ".xml"
	case config.ModeSchema, config.ModeKeys:
		ext = ".xsd"
	}

	files, err := filescanner.ScanDir(ctx, a.cfg.Path, filescanner.Filter{IncludeExts: []string{ext}})
	if err != nil {
		return fmt.Errorf("failed to scan dir: %w", err)
	}

	if len(files) == 0 {
		return errors.ErrNoFilesFound
	}

	for _, fileInfo := range files {
		if ctx.Err() != nil {
			return nil
		}

		if err := a.processFile(ctx, fileInfo); err != nil {
			return fmt.Errorf("failed to process file %s: %w", fileInfo.Path, err)
		}
	}

	return nil
}

func (a *Application) printHeader() {
	fmt.Println("-- >>>")
	fmt.Printf("-- Version: %s\n", "unknown")
	fmt.Printf("-- %s\n", a.cfg)
	fmt.Printf("-- Started at: %s\n", time.Now())
	fmt.Println("-- <<<")
	fmt.Println()
}

func (a *Application) processFile(ctx context.Context, fileInfo model.FileInfo) error {
	fileName := filepath.Base(fileInfo.Path)

	tableName, err := sqlbuilder.ResolveTableName(fileName)
	if err != nil {
		return fmt.Errorf("failed to resolve table name: %w", err)
	}

	startedAt := time.Now()
	_, _ = fmt.Fprintf(os.Stderr, "Started file %q (%s) to table %q at %q.\n", fileInfo.Path, humanize.Bytes(uint64(fileInfo.Size)), tableName, startedAt.Format(time.RFC3339))

	switch a.cfg.Mode {
	case config.ModeBulk, config.ModeUpsert:
		totalRows, err := a.processDataFile(ctx, tableName, fileInfo.Path)
		if err != nil && err != io.EOF {
			if ctx.Err() != nil {
				fmt.Println("-- Canceled")
			}
			return fmt.Errorf("failed to handle data file: %w", err)
		}
		_, _ = fmt.Fprintf(os.Stderr, "Handled %d rows\n", totalRows)
	case config.ModeSchema, config.ModeKeys:
		if err := a.processSchemaFile(ctx, tableName, fileInfo.Path); err != nil {
			return fmt.Errorf("failed to handle schema file: %w", err)
		}
	}

	_, _ = fmt.Fprintf(os.Stderr, "Ended after %s.\n", time.Since(startedAt).Round(time.Second))
	return nil
}

func (a *Application) processDataFile(ctx context.Context, tableName string, filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	iterator := itemiterator.New(file)
	var sqlBuilder sqlbuilder.ImportBuilder
	totalRows := 0

	for {
		items, err := iterator.Next(ctx, a.cfg.BatchSize)
		if err != nil && err != io.EOF {
			if ctx.Err() != nil {
				return totalRows, err
			}
		}

		if len(items) == 0 {
			return totalRows, nil
		}

		if sqlBuilder == nil {
			sqlBuilder = a.resolveImportBuilder(tableName, items[0])
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

func (a *Application) processSchemaFile(ctx context.Context, tableName string, filePath string) error {
	var sqlBuilder sqlbuilder.SchemaBuilder

	switch a.cfg.DbType {
	case config.DBPostgres:
		sqlBuilder = sqlbuilder.NewPostgreSQLSchemaBuilder(a.cfg.DbSchema, tableName, a.cfg.IgnoreRequired, a.cfg.IgnorePrimaryKey)
	case config.DBMySQL:
		sqlBuilder = sqlbuilder.NewMySQLSchemaBuilder(a.cfg.DbSchema, tableName, a.cfg.IgnoreRequired, a.cfg.IgnorePrimaryKey)
	default:
		return fmt.Errorf("unsupported database type: %s", a.cfg.DbType)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var sql string

	switch a.cfg.Mode {
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

func (a *Application) resolveImportBuilder(tableName string, item map[string]string) sqlbuilder.ImportBuilder {
	attrs := sqlbuilder.ResolveAttrs(item)

	switch a.cfg.Mode {
	case config.ModeBulk:
		switch a.cfg.DbType {
		case config.DBPostgres:
			return sqlbuilder.NewPostgreSQLCopyBuilder(a.cfg.DbSchema, tableName, attrs)
		case config.DBMySQL:
			return sqlbuilder.NewMySQLLoadDataBuilder(a.cfg.DbSchema, tableName, attrs)
		}
	case config.ModeUpsert:
		switch a.cfg.DbType {
		case config.DBPostgres:
			return sqlbuilder.NewPostgreSQLUpsertBuilder(a.cfg.DbSchema, tableName, attrs)
		case config.DBMySQL:
			return sqlbuilder.NewMySQLUpsertBuilder(a.cfg.DbSchema, tableName, attrs)
		}
	}
	panic(fmt.Sprintf("failed to resolve import builder for %q", a.cfg.Mode))
}
