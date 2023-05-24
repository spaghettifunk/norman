package segment

import (
	"fmt"
	"os"

	"encoding/json"

	"github.com/bytedance/sonic"

	"github.com/spaghettifunk/norman/internal/common/schema"
	"github.com/spaghettifunk/norman/internal/common/types"
	"github.com/spaghettifunk/norman/pkg/dynamicstruct"
)

const (
	MaxNumberOfEntriesPerSegment = 3
)

type SegmentManager struct {
	segment     *Segment
	schema      *schema.Schema
	eventStruct interface{}
	eventReader dynamicstruct.Reader
	baseDir     string
}

func NewSegmentManager(schema *schema.Schema) (*SegmentManager, error) {
	es := createEventStruct(schema)
	er := dynamicstruct.NewReader(es)

	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	baseDir := fmt.Sprintf("%s/output/default/%s/", path, schema.Name)
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		return nil, err
	}

	return &SegmentManager{
		schema: schema,
		// Format: output/{tenantID}/{schemaName}
		baseDir:     baseDir,
		eventStruct: es,
		eventReader: er,
	}, nil
}

func (sm *SegmentManager) CreateNewSegment() error {
	s, err := NewSegment(sm.baseDir, sm.schema)
	if err != nil {
		return err
	}
	sm.segment = s
	return nil
}

func (sm *SegmentManager) InsertDataInSegment(values []byte) error {
	// validate if the incoming message is a JSON -- only type supported
	if !sm.isJSON(values) {
		return fmt.Errorf("invalid JSON event")
	}
	// transform values into a map[]
	if err := sonic.Unmarshal(values, &sm.eventStruct); err != nil {
		return err
	}
	// validate if event is according to the schema
	if err := sm.validateEvent(); err != nil {
		return err
	}
	// add data to segment
	if err := sm.insertData(); err != nil {
		return err
	}
	// check if reached maximum segment size, flush otherwise
	if sm.segment.GetLength(sm.schema.DimensionFieldSpecs[0].Name) > MaxNumberOfEntriesPerSegment {
		return sm.FlushSegment()
	}
	return nil
}

func createEventStruct(schema *schema.Schema) interface{} {
	es := dynamicstruct.NewStruct()

	for _, dimension := range schema.DimensionFieldSpecs {
		ty := types.GetDataType(dimension.Name, dimension.DataType)
		es.AddField(dimension.Name, ty.Typ, ty.Tag)
	}

	for _, metric := range schema.MetricFieldSpecs {
		ty := types.GetDataType(metric.Name, metric.DataType)
		es.AddField(metric.Name, ty.Typ, ty.Tag)
	}

	for _, dt := range schema.DateTimeFieldSpecs {
		ty := types.GetDataType(dt.Name, dt.DataType)
		es.AddField(dt.Name, ty.Typ, ty.Tag)
	}

	return es.Build().New()
}

func (sm *SegmentManager) insertData() error {
	for _, f := range sm.eventReader.GetAllFields() {
		if err := sm.segment.InsertData(f.Name(), f.Interface()); err != nil {
			return err
		}
	}
	return nil
}

// TODO: super hard. Not sure how to solve it
func (sm *SegmentManager) validateEvent() error {
	return nil
}

func (sm *SegmentManager) isJSON(str []byte) bool {
	var js json.RawMessage
	return sonic.Unmarshal(str, &js) == nil
}

func (sm *SegmentManager) GetSegmentLength(column string) int {
	return sm.segment.GetLength(column)
}

// FlushSegment first persist on disk the current segment
// secondly, it compresses the segment to save space and lastly
// it reset the memory object so that it can start over
func (sm *SegmentManager) FlushSegment() error {
	if err := sm.segment.Flush(); err != nil {
		return err
	}
	return sm.compressSegment()
}

func (sm *SegmentManager) compressSegment() error {
	// zip compression
	// ...
	return nil
}
