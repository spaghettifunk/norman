package types

type DataType struct {
	Name     string
	Sortable bool
	Numberic bool
}

var (
	IntType       DataType = NewDataType("INT", true, true)
	LongType      DataType = NewDataType("LONG", true, true)
	FloatType     DataType = NewDataType("FLOAT", true, true)
	DoubleType    DataType = NewDataType("DOUBLE", true, true)
	BooleanType   DataType = NewDataType("BOOL", true, true)
	TimestampType DataType = NewDataType("TIMESTAMP", true, true)
	StringType    DataType = NewDataType("STRING", false, true)
	JSONType      DataType = NewDataType("JSON", false, false)
	BytesType     DataType = NewDataType("BYTES", true, true)
	StructType    DataType = NewDataType("STRUCT", false, false)
	MapType       DataType = NewDataType("MAP", false, false)
	ListType      DataType = NewDataType("LIST", false, true)
	UnknownType   DataType = NewDataType("UNKNOWN", false, false)
)

func NewDataType(name string, numeric, sortable bool) DataType {
	return DataType{
		Name:     name,
		Sortable: sortable,
		Numberic: numeric,
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
