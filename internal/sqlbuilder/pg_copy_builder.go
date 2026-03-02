package sqlbuilder

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
)

type PostgreSQLCopyBuilder struct {
	dbSchema   string
	table      string
	primaryKey string
	attrs      []string
}

func NewPostgreSQLCopyBuilder(dbSchema string, tableName string, attrs []string) *PostgreSQLCopyBuilder {
	return &PostgreSQLCopyBuilder{
		dbSchema:   dbSchema,
		table:      tableName,
		primaryKey: resolvePrimaryKey(tableName),
		attrs:      attrs,
	}
}

func (b *PostgreSQLCopyBuilder) Build(rows []map[string]string) (string, error) {
	if len(rows) == 0 {
		return "", fmt.Errorf("no rows to build")
	}

	valuesStatement, err := b.buildValues(rows)
	if err != nil {
		return "", err
	}

	sql := fmt.Sprintf("COPY %s (%s) FROM STDIN WITH (FORMAT csv);\n%s\\.", buildFullTableName(b.dbSchema, b.table), b.buildColumns(), valuesStatement)

	return sql, nil
}

func (b *PostgreSQLCopyBuilder) buildValues(rows []map[string]string) (string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	for _, row := range rows {
		vals := make([]string, len(b.attrs))

		// to keep order of columns
		for i, attrName := range b.attrs {
			vals[i] = row[attrName]
		}

		if err := writer.Write(vals); err != nil {
			return "", err
		}
	}

	writer.Flush()

	return buf.String(), nil
}

func (b *PostgreSQLCopyBuilder) buildColumns() string {
	columns := make([]string, len(b.attrs))
	for i, attrName := range b.attrs {
		columns[i] = escapeColumnName(resolveColumnName(attrName))
	}
	return strings.Join(columns, ",")
}
