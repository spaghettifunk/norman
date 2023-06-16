package indexer

type Indexer interface {
	GetColumnName() string
	Build(id string, value interface{}) bool
	Search(value interface{}) []uint32
}
