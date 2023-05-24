package schema

import (
	"github.com/spaghettifunk/norman/internal/common/types"
)

type Schema struct {
	Name                string                `json:"name"`
	DimensionFieldSpecs []*DimensionFieldSpec `json:"dimensionFieldSpecs"`
	MetricFieldSpecs    []*MetricFieldSpec    `json:"metricFieldSpecs,omitempty"`
	DateTimeFieldSpecs  []*DateTimeFieldSpec  `json:"dateTimeFieldSpecs"`
}

type DimensionFieldSpec struct {
	Name             string            `json:"name"`
	DataType         types.DataTypeVal `json:"dataType"`
	SingleValueField bool              `json:"singleValueField,omitempty"`
	DefaultNullValue interface{}       `json:"defaultNullValue,omitempty"`
}

type MetricFieldSpec struct {
	Name             string            `json:"name"`
	DataType         types.DataTypeVal `json:"dataType"`
	DefaultNullValue interface{}       `json:"defaultNullValue,omitempty"`
}

type DateTimeFieldSpec struct {
	Name        string            `json:"name"`
	DataType    types.DataTypeVal `json:"dataType"`
	Format      string            `json:"format,omitempty"`
	Granularity string            `json:"granularity,omitempty"`
}

func (s *Schema) Validate(dt types.DataType) error {
	return nil
}
