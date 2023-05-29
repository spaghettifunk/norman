package entities

import (
	"github.com/apache/arrow/go/arrow"
	"github.com/spaghettifunk/norman/internal/common/types"
)

type Schema struct {
	Name                string                `json:"name"`
	DimensionFieldSpecs []*DimensionFieldSpec `json:"dimensionFieldSpecs"`
	MetricFieldSpecs    []*MetricFieldSpec    `json:"metricFieldSpecs,omitempty"`
	DateTimeFieldSpecs  []*DateTimeFieldSpec  `json:"dateTimeFieldSpecs"`
}

type DimensionFieldSpec struct {
	Name             string      `json:"name"`
	DataType         string      `json:"dataType"`
	SingleValueField bool        `json:"singleValueField,omitempty"`
	Nullable         bool        `json:"nullable,omitempty"`
	DefaultNullValue interface{} `json:"defaultNullValue,omitempty"`
}

type MetricFieldSpec struct {
	Name             string      `json:"name"`
	DataType         string      `json:"dataType"`
	Nullable         bool        `json:"nullable,omitempty"`
	DefaultNullValue interface{} `json:"defaultNullValue,omitempty"`
}

type DateTimeFieldSpec struct {
	Name        string `json:"name"`
	DataType    string `json:"dataType"`
	Format      string `json:"format,omitempty"`
	Granularity string `json:"granularity,omitempty"`
}

func (s *Schema) Validate(dt types.DataType) error {
	return nil
}

// TODO: double check if it always work
func (s *Schema) GetArrowSchema() *arrow.Schema {
	fields := []arrow.Field{}
	for _, dimension := range s.DimensionFieldSpecs {
		ty := types.GetDataType(dimension.DataType)
		fields = append(fields, arrow.Field{
			Name:     dimension.Name,
			Type:     ty.Typ,
			Nullable: dimension.Nullable,
		})
	}

	for _, metric := range s.MetricFieldSpecs {
		ty := types.GetDataType(metric.DataType)
		fields = append(fields, arrow.Field{
			Name:     metric.Name,
			Type:     ty.Typ,
			Nullable: metric.Nullable,
		})
	}

	for _, dt := range s.DateTimeFieldSpecs {
		ty := types.GetDataType(dt.DataType)
		fields = append(fields, arrow.Field{
			Name:     dt.Name,
			Type:     ty.Typ,
			Nullable: false,
		})
	}

	return arrow.NewSchema(fields, nil)
}
