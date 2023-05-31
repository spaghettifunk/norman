package entities

import (
	"github.com/apache/arrow/go/v12/arrow"

	"github.com/spaghettifunk/norman/internal/common/types"
)

type Table struct {
	Name        string        `json:"name"`
	Schema      *Schema       `json:"schema"`
	EventSchema *arrow.Schema `json:"-"`
}

func NewTable(name string, schema *Schema) (*Table, error) {
	return &Table{
		Name:        name,
		Schema:      schema,
		EventSchema: createEventSchema(schema),
	}, nil
}

func createEventSchema(schema *Schema) *arrow.Schema {
	fields := []arrow.Field{}

	for _, dimension := range schema.DimensionFieldSpecs {
		ty := types.GetDataType(dimension.DataType)
		fields = append(fields, arrow.Field{
			Name:     dimension.Name,
			Type:     ty.Typ,
			Nullable: dimension.Nullable,
		})
	}

	for _, metric := range schema.MetricFieldSpecs {
		ty := types.GetDataType(metric.DataType)
		fields = append(fields, arrow.Field{
			Name:     metric.Name,
			Type:     ty.Typ,
			Nullable: metric.Nullable,
		})
	}

	// datetime cannot be null
	for _, dt := range schema.DateTimeFieldSpecs {
		ty := types.GetDataType(dt.DataType)
		fields = append(fields, arrow.Field{
			Name:     dt.Name,
			Type:     ty.Typ,
			Nullable: false,
		})
	}

	return arrow.NewSchema(fields, nil)
}
