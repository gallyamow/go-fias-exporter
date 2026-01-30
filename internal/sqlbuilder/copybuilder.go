package sqlbuilder

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
)

type CopyBuilder struct {
	tableName  string
	primaryKey string
	attrs      []string
}

func NewCopyBuilder(tablename string, primaryKey string, attrs []string) *CopyBuilder {
	return &CopyBuilder{
		tableName:  tablename,
		primaryKey: primaryKey,
		attrs:      attrs,
	}
}

func (b *CopyBuilder) Build(rows []map[string]string) (string, error) {
	if len(rows) == 0 {
		return "", fmt.Errorf("no rows to build")
	}

	valuesStatement, err := b.buildValues(rows)
	if err != nil {
		return "", err
	}

	sql := fmt.Sprintf("COPY %s (%s) FROM STDIN WITH (FORMAT csv);\n\n%s\n\\.", b.tableName, b.buildColumns(), valuesStatement)

	return sql, nil
}

func (b *CopyBuilder) buildValues(rows []map[string]string) (string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	for _, row := range rows {
		vals := make([]string, len(row))
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

func (b *CopyBuilder) buildColumns() string {
	columns := make([]string, len(b.attrs))
	for i, attrName := range b.attrs {
		columns[i] = ResolveColumnName(attrName)
	}
	return strings.Join(columns, ",")
}
