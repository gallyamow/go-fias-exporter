package sqlbuilder

import (
	"fmt"
	"strings"
)

type MySQLInsertBuilder struct {
	dbSchema   string
	table      string
	primaryKey string
	attrs      []string
}

func NewMySQLInsertBuilder(dbSchema string, tableName string, attrs []string) *MySQLInsertBuilder {
	return &MySQLInsertBuilder{
		dbSchema:   dbSchema,
		table:      tableName,
		primaryKey: resolvePrimaryKey(tableName),
		attrs:      attrs,
	}
}

func (b *MySQLInsertBuilder) Build(rows []map[string]string) (string, error) {
	if len(rows) == 0 {
		return "", fmt.Errorf("no rows to build")
	}

	valuesStatement, err := b.buildValues(rows)
	if err != nil {
		return "", err
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", buildFullTableName(b.dbSchema, b.table), b.buildColumns(), valuesStatement))

	if b.primaryKey != "" {
		sb.WriteString(" ")
		sb.WriteString(b.buildOnDuplicateKeyUpdate())
	}

	sb.WriteString(";")

	return sb.String(), nil
}

func (b *MySQLInsertBuilder) buildValues(rows []map[string]string) (string, error) {
	var res []string

	for _, row := range rows {
		vals := make([]string, len(b.attrs))

		// to keep order of columns
		for i, attrName := range b.attrs {
			// (to keep simple quote for all values)
			escapedValue := escapeString(row[attrName])
			vals[i] = fmt.Sprintf("'%s'", escapedValue)
		}

		res = append(res, fmt.Sprintf("(%s)", strings.Join(vals, ",")))
	}

	return strings.Join(res, ","), nil
}

func (b *MySQLInsertBuilder) buildColumns() string {
	columns := make([]string, len(b.attrs))
	for i, attrName := range b.attrs {
		columns[i] = escapeColumnName(resolveColumnName(attrName))
	}
	return strings.Join(columns, ",")
}

func (b *MySQLInsertBuilder) buildOnDuplicateKeyUpdate() string {
	var setters []string
	for _, attrName := range b.attrs {
		column := escapeColumnName(resolveColumnName(attrName))
		if column == b.primaryKey {
			continue
		}
		setters = append(setters, fmt.Sprintf("%s=VALUES(%s)", column, column))
	}

	return fmt.Sprintf("ON DUPLICATE KEY UPDATE %s", strings.Join(setters, ","))
}
