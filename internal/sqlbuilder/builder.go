package sqlbuilder

type ImportBuilder interface {
	Build(rows []map[string]string) (string, error)
}

type SchemaBuilder interface {
	Build(data []byte) (string, error)
}
