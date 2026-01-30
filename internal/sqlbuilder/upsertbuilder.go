package sqlbuilder

import (
	"fmt"
	"strings"
)

type UpsertBuilder struct {
	schema     string
	table      string
	primaryKey string
	attrs      []string
}

func NewUpsertBuilder(schema string, tablename string, primaryKey string, attrs []string) *UpsertBuilder {
	return &UpsertBuilder{
		schema:     schema,
		table:      tablename,
		primaryKey: primaryKey,
		attrs:      attrs,
	}
}

func (b *UpsertBuilder) Build(rows []map[string]string) (string, error) {
	if len(rows) == 0 {
		return "", fmt.Errorf("no rows to build")
	}

	valuesStatement, err := b.buildValues(rows)
	if err != nil {
		return "", err
	}

	res := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", b.buildTablename(), b.buildColumns(), valuesStatement)
	if b.primaryKey != "" {
		res += " " + b.buildOnConflict()
	}

	return res, nil
}

func (b *UpsertBuilder) buildValues(rows []map[string]string) (string, error) {
	var res []string

	for _, row := range rows {
		var vals []string

		// to keep order of columns
		for _, attrName := range b.attrs {
			// (to keep simple quote for all values)
			escapedValue := strings.ReplaceAll(row[attrName], "'", "''")
			vals = append(vals, fmt.Sprintf("'%s'", escapedValue))
		}

		res = append(res, fmt.Sprintf("(%s)", strings.Join(vals, ",")))
	}

	return strings.Join(res, ","), nil
}

func (b *UpsertBuilder) buildTablename() string {
	if b.schema != "" {
		return fmt.Sprintf("%s.%s", b.schema, b.table)
	}
	return b.table
}

func (b *UpsertBuilder) buildColumns() string {
	columns := make([]string, len(b.attrs))
	for i, attrName := range b.attrs {
		columns[i] = ResolveColumnName(attrName)
	}
	return strings.Join(columns, ",")
}

func (b *UpsertBuilder) buildOnConflict() string {
	var setters []string
	for _, attrName := range b.attrs {
		column := ResolveColumnName(attrName)
		if column == b.primaryKey {
			continue
		}
		setters = append(setters, fmt.Sprintf("%s=EXCLUDED.%s", column, column))
	}

	return fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET %s", b.primaryKey, strings.Join(setters, ","))
}
