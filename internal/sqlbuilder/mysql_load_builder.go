package sqlbuilder

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
)

type MySQLLoadDataBuilder struct {
	dbSchema   string
	table      string
	primaryKey string
	attrs      []string
}

func NewMySQLLoadDataBuilder(dbSchema string, tableName string, attrs []string) *MySQLLoadDataBuilder {
	return &MySQLLoadDataBuilder{
		dbSchema:   dbSchema,
		table:      tableName,
		primaryKey: resolvePrimaryKey(tableName),
		attrs:      attrs,
	}
}

func (b *MySQLLoadDataBuilder) Build(rows []map[string]string) (string, error) {
	if len(rows) == 0 {
		return "", fmt.Errorf("no rows to build")
	}

	valuesStatement, err := b.buildValues(rows)
	if err != nil {
		return "", err
	}

	sql := fmt.Sprintf(`LOAD DATA LOCAL INFILE 'stdin' INTO TABLE %s FIELDS TERMINATED BY ',' ENCLOSED BY '"' LINES TERMINATED BY '\n' (%s);`+"\n%s", buildFullTableName(b.dbSchema, b.table), b.buildColumns(), valuesStatement)

	return sql, nil
}

func (b *MySQLLoadDataBuilder) buildValues(rows []map[string]string) (string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	for _, row := range rows {
		vals := make([]string, len(b.attrs))

		// to keep order of columns
		for i, attrName := range b.attrs {
			value := row[attrName]
			// Convert boolean string values to integers for MySQL BOOLEAN columns
			if resolveColumnName(attrName) == "isactive" {
				value = convertBooleanToMySQL(value)
			}
			vals[i] = value
		}

		if err := writer.Write(vals); err != nil {
			return "", err
		}
	}

	writer.Flush()

	return buf.String(), nil
}

func (b *MySQLLoadDataBuilder) buildColumns() string {
	columns := make([]string, len(b.attrs))
	for i, attrName := range b.attrs {
		columns[i] = escapeColumnNameMySQL(resolveColumnName(attrName))
	}
	return strings.Join(columns, ",")
}
