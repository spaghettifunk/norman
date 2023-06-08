package entities

import (
	"github.com/apache/arrow/go/v12/arrow"
)

type Table struct {
	Name        string        `json:"name"`
	Schema      *Schema       `json:"schema"`
	EventSchema *arrow.Schema `json:"-"`
}

func NewTable(name string, schema *Schema) (*Table, error) {
	s, err := schema.GetFullArrowSchema()
	if err != nil {
		return nil, err
	}
	return &Table{
		Name:        name,
		Schema:      schema,
		EventSchema: s,
	}, nil
}

func (t *Table) GetDimensionFields() []arrow.Field {
	return t.Schema.GetDimensionFields()
}

func (t *Table) GetMetricFields() []arrow.Field {
	return t.Schema.GetMetricFields()
}

func (t *Table) GetDatetimeField() arrow.Field {
	return t.Schema.GetDatetimeField()
}
