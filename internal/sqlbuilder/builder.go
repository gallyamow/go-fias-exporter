package sqlbuilder

type Builder interface {
	Build(rows []map[string]string) (string, error)
}
