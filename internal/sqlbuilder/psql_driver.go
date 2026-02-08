package sqlbuilder

import (
	"fmt"
	"strings"
)

type PSQLDriver struct {
	BaseDriver
}

func (d *PSQLDriver) CreateTable(fullTable string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("CREATE TABLE %s (\n%s\n);", fullTable, d.buildColumns(attrs)))
	sb.WriteRune('\n')
	sb.WriteString(fmt.Sprintf("%s;", fmt.Sprintf("COMMENT ON TABLE %s IS '%s'", fullTable, descr)))
	sb.WriteRune('\n')
	sb.WriteString(fmt.Sprintf("%s", d.BuildColumnComments(fullTable, attrs)))
	sb.WriteRune('\n')

	return sb.String()
}

func (d *PSQLDriver) resolveColumnType(xsdType string) string {
	switch xsdType {
	case "xs:string":
		return "VARCHAR"
	case "xs:int":
		return "INT"
	case "xs:long":
		return "BIGINT"
	case "xs:boolean":
		return "BOOLEAN"
	case "xs:date":
		return "DATE"
	case "xs:dateTime":
		return "TIMESTAMP"
	default:
		return "VARCHAR"
	}
}
