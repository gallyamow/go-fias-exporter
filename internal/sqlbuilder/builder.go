package sqlbuilder

import (
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

type UpsertBuilder struct {
	tableName   string
	primaryKeys []string
	columns     []string
}

func New(tablename string, primaryKeys []string, columns []string) *UpsertBuilder {
	return &UpsertBuilder{
		tableName:   tablename,
		primaryKeys: primaryKeys,
		columns:     columns,
	}
}

func (b *UpsertBuilder) Build(rows []map[string]string) (string, error) {
	if rows == nil || len(rows) == 0 {
		return "", fmt.Errorf("no rows to build")
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES %s %s", b.tableName, b.buildColumns(), b.buildValues(rows), b.buildOnConflict()), nil
}

func (b *UpsertBuilder) buildValues(rows []map[string]string) string {
	var res []string

	for _, row := range rows {
		var vals []string

		// to keep order of columns
		for _, column := range b.columns {
			// (to keep simple quote for all values)
			vals = append(vals, fmt.Sprintf("'%s'", row[column]))
		}

		res = append(res, fmt.Sprintf("(%s)", strings.Join(vals, ",")))
	}

	return strings.Join(res, ",")
}

func (b *UpsertBuilder) buildColumns() string {
	return strings.Join(b.columns, ",")
}

func (b *UpsertBuilder) buildOnConflict() string {
	var setters []string
	for _, column := range b.columns {
		if slices.Contains(b.primaryKeys, column) {
			continue
		}
		setters = append(setters, fmt.Sprintf("%s=EXCLUDED.%s", column, column))
	}

	return fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET %s", strings.Join(b.primaryKeys, ","), strings.Join(setters, ","))
}

// ResolveTableName resolves table name from filename.
// Examples:
// - AS_ADDR_OBJ_20250626_bc6f64d9-fb28-40d6-8a99-57e44b920d07.XML => addr_obj
// - AS_CHANGE_HISTORY_20250626_d1a57485-156c-4463-8a23-2328fb0f6f9d => change_history
func ResolveTableName(filename string) (string, error) {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))

	re := regexp.MustCompile(`(?i)^AS_([A-Z_]+)_\d{8}`)
	if m := re.FindStringSubmatch(base); len(m) == 2 {
		return strings.ToLower(m[1]), nil
	}

	return "", fmt.Errorf("cannot resolve table name from filename: %s", filename)
}
