package sqlbuilder

type MySQLDriver struct {
	BaseDriver
}

func (d *MySQLDriver) resolveColumnType(xsdType string) string {
	switch xsdType {
	case "xs:string":
		return "VARCHAR(255)"
	case "xs:int":
		return "INT"
	case "xs:long":
		return "BIGINT"
	case "xs:boolean":
		return "TINYINT(1)" // BOOLEAN = alias
	case "xs:date":
		return "DATE"
	case "xs:dateTime":
		return "DATETIME"
	default:
		return "VARCHAR(255)"
	}
}
