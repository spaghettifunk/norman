package types

import (
	"fmt"
	"time"
)

type DataTypeVal string

const (
	Integer   DataTypeVal = "INT"
	Long      DataTypeVal = "LONG"
	Float     DataTypeVal = "FLOAT"
	Double    DataTypeVal = "DOUBLE"
	Boolean   DataTypeVal = "BOOL"
	Timestamp DataTypeVal = "TIMESTAMP"
	String    DataTypeVal = "STRING"
	JSON      DataTypeVal = "JSON"
	Bytes     DataTypeVal = "BYTES"
	Map       DataTypeVal = "MAP"
	List      DataTypeVal = "LIST"
	Unknown   DataTypeVal = "UNKNOWN"
)

type DataType struct {
	Name     DataTypeVal
	Sortable bool
	Numberic bool
	Typ      interface{}
	Tag      string
}

func (d DataTypeVal) ToString() string {
	return string(d)
}

func newDataType(name DataTypeVal, numeric, sortable bool, typ interface{}, tag string) DataType {
	return DataType{
		Name:     name,
		Sortable: sortable,
		Numberic: numeric,
		Typ:      typ,
		Tag:      tag,
	}
}

func GetDataType(name string, typ DataTypeVal) DataType {
	tag := fmt.Sprintf("`json:\"%s\"`", name)
	switch typ {
	case Integer:
		return newDataType(typ, true, true, 0, tag)
	case Long:
		return newDataType(typ, true, true, uint(0), tag)
	case Float:
		return newDataType(typ, true, true, float32(0), tag)
	case Double:
		return newDataType(typ, true, true, float64(0), tag)
	case Boolean:
		return newDataType(typ, true, true, false, tag)
	case Timestamp:
		return newDataType(typ, true, true, time.Time{}, tag)
	case String:
		return newDataType(typ, false, true, "", tag)
	case JSON:
		return newDataType(typ, false, false, "", tag)
	case Bytes:
		return newDataType(typ, true, true, make([]byte, 0), tag)
	case Map:
		return newDataType(typ, false, false, make(map[string]interface{}, 0), tag)
	case List:
		return newDataType(typ, false, true, make([]interface{}, 0), tag)
	default:
		return newDataType(Unknown, false, false, nil, tag)
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
