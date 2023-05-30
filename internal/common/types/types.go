package types

import (
	"github.com/apache/arrow/go/v12/arrow"
)

const (
	Integer   string = "INT"
	Long      string = "LONG"
	Float     string = "FLOAT"
	Double    string = "DOUBLE"
	Boolean   string = "BOOL"
	Timestamp string = "TIMESTAMP"
	String    string = "STRING"
	JSON      string = "JSON"
	Bytes     string = "BYTES"
	Map       string = "MAP"
	List      string = "LIST"
	Unknown   string = "UNKNOWN"
)

type DataType struct {
	TypeName string
	Sortable bool
	Numberic bool
	Typ      arrow.DataType
}

func newDataType(name string, numeric, sortable bool, typ arrow.DataType) DataType {
	return DataType{
		TypeName: name,
		Sortable: sortable,
		Numberic: numeric,
		Typ:      typ,
	}
}

func GetDataType(tn string) DataType {
	switch tn {
	case Integer:
		return newDataType(tn, true, true, arrow.PrimitiveTypes.Int32)
	case Long:
		return newDataType(tn, true, true, arrow.PrimitiveTypes.Uint32)
	case Float:
		return newDataType(tn, true, true, arrow.PrimitiveTypes.Float32)
	case Double:
		return newDataType(tn, true, true, arrow.PrimitiveTypes.Float64)
	case Boolean:
		return newDataType(tn, true, true, arrow.FixedWidthTypes.Boolean)
	case Timestamp:
		return newDataType(tn, true, true, arrow.FixedWidthTypes.Date64)
	case String:
		return newDataType(tn, false, true, arrow.BinaryTypes.String)
	case JSON:
		return newDataType(tn, false, false, arrow.BinaryTypes.String)
	case Bytes:
		return newDataType(tn, true, true, arrow.BinaryTypes.Binary)
	case Map:
		return newDataType(tn, false, false, &arrow.StructType{})
	case List:
		return newDataType(tn, false, true, &arrow.ListType{})
	default:
		return newDataType(Unknown, false, false, nil)
	}
}

type FieldType string

const (
	DimensionType FieldType = "DIMENSION"
	MetricType    FieldType = "METRIC"
	TimeType      FieldType = "TIME"
	DatetimeType  FieldType = "DATE_TIME"
	ComplexType   FieldType = "COMPLEX"
)
