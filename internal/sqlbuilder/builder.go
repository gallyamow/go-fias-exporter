package sqlbuilder

type Builder interface {
	Build(rows []map[string]string) (string, error)
}

type SchemaBuilderInterface interface {
	Build(data []byte) (string, error)
}
