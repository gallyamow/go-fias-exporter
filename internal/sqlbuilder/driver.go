package sqlbuilder

import (
	"fmt"
	"strings"
)

type Driver interface {
	CreateTable() string
	ResolveColumnType(xsdType string) string
	BuildTableComment(fullTable string, descr string) string
	BuildColumnComments(fullTable string, attrs []attribute) string
}

type BaseDriver struct {
}

func (b *SchemaBuilder) buildColumns(attrs []attribute) string {
	columns := make([]string, len(attrs))
	for i, attr := range attrs {
		columns[i] = fmt.Sprintf("\t%s", b.buildColumn(attr))
	}

	return strings.Join(columns, ",\n")
}

func (d *BaseDriver) BuildTableComment(fullTable string, descr string) string {
	return fmt.Sprintf("COMMENT ON TABLE %s IS '%s'", fullTable, descr)
}

func (d *BaseDriver) BuildColumnComments(fullTable string, attrs []attribute) string {
	columns := make([]string, len(attrs))
	for i, attr := range attrs {
		columnName := resolveColumnName(attr.Name)
		columns[i] = fmt.Sprintf("COMMENT ON COLUMN %s.%s IS '%s';", fullTable, columnName, attr.Annotation.Documentation)
	}

	return strings.Join(columns, "\n")
}
