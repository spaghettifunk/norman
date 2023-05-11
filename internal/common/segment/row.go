package segment

type Row struct {
	Values map[*Column]interface{} `json:"-"`
}

func NewEmptyRow() *Row {
	return &Row{}
}

func (r *Row) AddValue(val interface{}, col *Column) error {
	r.Values[col] = val
	return nil
}
