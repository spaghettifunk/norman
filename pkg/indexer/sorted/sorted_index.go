package sortedindex

type SortedIndex struct {
	columnName string
}

func New(columnName string) *SortedIndex {
	return &SortedIndex{
		columnName: columnName,
	}
}

func (i *SortedIndex) Build(id string, value interface{}) bool {
	return true
}

func (i *SortedIndex) Search(value interface{}) []uint32 {
	return nil
}

func (i *SortedIndex) GetColumnName() string {
	return i.columnName
}
