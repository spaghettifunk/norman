package entities

import (
	"fmt"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/spaghettifunk/norman/internal/common/types"
)

type Schema struct {
	Name                string                `json:"name"`
	DimensionFieldSpecs []*DimensionFieldSpec `json:"dimensionFieldSpecs"`
	MetricFieldSpecs    []*MetricFieldSpec    `json:"metricFieldSpecs,omitempty"`
	DateTimeFieldSpecs  *DateTimeFieldSpec    `json:"dateTimeFieldSpecs"`
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

func (s *Schema) GetFullArrowSchema() (*arrow.Schema, error) {
	// datetime cannot be null
	if s.DateTimeFieldSpecs == nil {
		return nil, fmt.Errorf("datetime field in schema cannot be null")
	}

	fields := []arrow.Field{}
	fields = append(fields, s.GetDimensionFields()...)
	fields = append(fields, s.GetMetricFields()...)
	fields = append(fields, s.GetDatetimeField())

	return arrow.NewSchema(fields, nil), nil
}

func (s *Schema) GetDimensionField(name string) *DimensionFieldSpec {
	for _, dim := range s.DimensionFieldSpecs {
		if dim.Name == name {
			return dim
		}
	}
	return nil
}

func (s *Schema) GetDimensionFields() []arrow.Field {
	fields := []arrow.Field{}
	for _, dimension := range s.DimensionFieldSpecs {
		ty := types.GetDataType(dimension.DataType)
		fields = append(fields, arrow.Field{
			Name:     dimension.Name,
			Type:     ty.Typ,
			Nullable: dimension.Nullable,
		})
	}
	return fields
}

func (s *Schema) GetMetricFields() []arrow.Field {
	fields := []arrow.Field{}
	for _, metric := range s.MetricFieldSpecs {
		ty := types.GetDataType(metric.DataType)
		fields = append(fields, arrow.Field{
			Name:     metric.Name,
			Type:     ty.Typ,
			Nullable: metric.Nullable,
		})
	}
	return fields
}

func (s *Schema) GetDatetimeField() arrow.Field {
	ty := types.GetDataType(s.DateTimeFieldSpecs.DataType)
	return arrow.Field{
		Name:     s.DateTimeFieldSpecs.Name,
		Type:     ty.Typ,
		Nullable: false,
	}
}

func (s *Schema) GetDatetimeFormat() (*DateTimeFormatSpec, error) {
	return NewDateTimeFormatSpec(s.DateTimeFieldSpecs.Format)
}

func (s *Schema) GetGranularity() (*GranularitySpec, error) {
	return NewGranularitySpec(s.DateTimeFieldSpecs.Granularity)
}
