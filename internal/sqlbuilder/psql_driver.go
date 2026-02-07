package sqlbuilder

import "fmt"

type PSQLDriver struct {
	BaseDriver
}

func (d *PSQLDriver) CreateTable(fullTable string) string {
	res := fmt.Sprintf("CREATE TABLE %s (\n%s\n);\n%s;\n%s",
		fullTable,
		b.buildColumns(attrs),
		b.driver.BuildTableComment(b.fullTable, descr),
		b.driver.BuildColumnComments(b.fullTable, attrs),
	)
	return res
}

func (d *PSQLDriver) ResolveColumnType(xsdType string) string {
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
