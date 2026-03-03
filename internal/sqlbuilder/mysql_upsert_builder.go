package sqlbuilder

import (
	"fmt"
	"strings"
)

type MySQLUpsertBuilder struct {
	dbSchema   string
	table      string
	primaryKey string
	attrs      []string
}

func NewMySQLUpsertBuilder(dbSchema string, tableName string, attrs []string) *MySQLUpsertBuilder {
	return &MySQLUpsertBuilder{
		dbSchema:   dbSchema,
		table:      tableName,
		primaryKey: resolvePrimaryKey(tableName),
		attrs:      attrs,
	}
}

func (b *MySQLUpsertBuilder) Build(rows []map[string]string) (string, error) {
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
		sb.WriteString(b.buildOnConflict())
	}

	sb.WriteString(";")

	return sb.String(), nil
}

func (b *MySQLUpsertBuilder) buildValues(rows []map[string]string) (string, error) {
	var res []string

	for _, row := range rows {
		vals := make([]string, len(b.attrs))

		// to keep order of columns
		for i, attrName := range b.attrs {
			// Convert boolean string values to integers for MySQL BOOLEAN columns
			value := row[attrName]

			if resolveColumnName(attrName) == "isactive" {
				value = convertBooleanToMySQL(value)
			}
			// (to keep simple quote for all values)
			escapedValue := escapeString(value)
			vals[i] = fmt.Sprintf("'%s'", escapedValue)
		}

		res = append(res, fmt.Sprintf("(%s)", strings.Join(vals, ",")))
	}

	return strings.Join(res, ","), nil
}

func (b *MySQLUpsertBuilder) buildColumns() string {
	columns := make([]string, len(b.attrs))
	for i, attrName := range b.attrs {
		columns[i] = escapeColumnNameMySQL(resolveColumnName(attrName))
	}
	return strings.Join(columns, ",")
}

func (b *MySQLUpsertBuilder) buildOnConflict() string {
	var setters []string
	for _, attrName := range b.attrs {
		column := escapeColumnNameMySQL(resolveColumnName(attrName))
		if column == b.primaryKey {
			continue
		}
		setters = append(setters, fmt.Sprintf("%s=VALUES(%s)", column, column))
	}

	return fmt.Sprintf("ON DUPLICATE KEY UPDATE %s", strings.Join(setters, ","))
}
