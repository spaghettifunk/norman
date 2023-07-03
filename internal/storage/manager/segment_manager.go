package manager

import (
	"bytes"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
	"github.com/rs/zerolog/log"

	"github.com/spaghettifunk/norman/internal/storage/segment"
)

var (
	minTimestamp = time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
)

var buffPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

type SegmentManager struct {
	tableName      string
	baseDir        string
	activeSegment  *segment.Segment
	segments       []*segment.Segment
	builder        *array.RecordBuilder
	schema         *arrow.Schema
	partition      int
	partitionStart time.Time
	interval       time.Duration
	eventsCounter  int
}

func NewSegmentManager(tableName, baseDir string, interval time.Duration) (*SegmentManager, error) {
	return &SegmentManager{
		tableName:     tableName,
		baseDir:       baseDir,
		partition:     0,
		interval:      interval,
		eventsCounter: 0,
	}, nil
}

func (s *SegmentManager) Create(schema *arrow.Schema) error {
	d := time.Now()

	// create the new record builder for inserting data to arrow file
	mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
	b := array.NewRecordBuilder(mem, schema)

	fPath := fmt.Sprintf("%s/%s", s.baseDir, d.Format(partitionTimeFormat))
	sgm, err := segment.NewSegment(fPath, s.partition, schema)
	if err != nil {
		return err
	}
	s.activeSegment = sgm
	s.builder = b
	s.schema = schema

	return nil
}

// AppendData appends to the apache arrow record the incoming data. It returns the EventID if everything is successful
func (s *SegmentManager) AppendData(event map[string]interface{}, datetimeFieldName string, granularityUnit time.Duration) (string, error) {
	// Get the timestamp of the consumed message
	dtVal := int64(event[datetimeFieldName].(float64))
	eventTimestamp := time.Unix(0, dtVal*int64(time.Millisecond))

	// initial partition time setup
	if s.partitionStart.Equal(minTimestamp) {
		s.partitionStart = eventTimestamp.Truncate(granularityUnit)
	}

	partitionInterval := s.partitionStart.Add(s.interval)

	// add norman ID to the event
	event[eventIDName] = s.generateEventID(partitionInterval)

	// if the interval is passed then it creates a new partition
	if !(eventTimestamp.After(s.partitionStart) && eventTimestamp.Before(partitionInterval)) {
		if err := s.Flush(); err != nil {
			return "", err
		}
		if err := s.Create(s.schema); err != nil {
			return "", err
		}
		// Update the current partition and its start time
		s.partition++
		s.partitionStart = eventTimestamp.Truncate(granularityUnit)
	}

	s.processEvent(event)
	s.eventsCounter++

	return event[eventIDName].(string), nil
}

// this is processed concurrently considering that there can be hundreds of columns per event
func (s *SegmentManager) processEvent(event map[string]interface{}) {
	for idx, field := range s.builder.Schema().Fields() {
		val, ok := event[field.Name]
		if !ok {
			log.Error().Msgf("could not find column %s in builder", field.Name)
			return
		}
		b := s.builder.Field(idx)
		s.appendValue(val, field, b)
	}
}

func (s *SegmentManager) Flush() error {
	record := s.builder.NewRecord()
	if err := s.activeSegment.Flush(record); err != nil {
		return err
	}
	// store the active segments in the list of segments
	s.segments = append(s.segments, s.activeSegment)
	return nil
}

func (s *SegmentManager) generateEventID(partitionInterval time.Time) string {
	bf := buffPool.Get().(*bytes.Buffer)
	defer buffPool.Put(bf)
	bf.Reset()

	bf.WriteString(s.tableName)
	bf.WriteString("_P")
	bf.WriteString(strconv.Itoa(s.partition))
	bf.WriteByte('_')
	bf.WriteString(s.partitionStart.Format(partitionTimeFormat))
	bf.WriteByte('_')
	bf.WriteString(partitionInterval.Format(partitionTimeFormat))
	bf.WriteByte('_')
	bf.WriteString(strconv.Itoa(s.eventsCounter))

	return bf.String()
}

// val interface{} when Unmarshalled for numbers is always float64. Initial
// assertion is necessary before the correct type casting
func (s *SegmentManager) appendValue(val interface{}, field arrow.Field, builder array.Builder) {
	switch field.Type {
	case arrow.PrimitiveTypes.Int32:
		v := int32(val.(float64))
		builder.(*array.Int32Builder).Append(v)
	case arrow.PrimitiveTypes.Uint32:
		v := uint32(val.(float64))
		builder.(*array.Uint32Builder).Append(v)
	case arrow.PrimitiveTypes.Float32:
		v := float32(val.(float64))
		builder.(*array.Float32Builder).Append(v)
	case arrow.PrimitiveTypes.Float64:
		v := val.(float64)
		builder.(*array.Float64Builder).Append(v)
	case arrow.FixedWidthTypes.Boolean:
		v := val.(bool)
		builder.(*array.BooleanBuilder).Append(v)
	case arrow.PrimitiveTypes.Int64:
		v := int64(val.(float64))
		builder.(*array.Int64Builder).Append(v)
	case arrow.BinaryTypes.String:
		v := val.(string)
		builder.(*array.StringBuilder).Append(v)
	case arrow.BinaryTypes.Binary:
		v := val.(int32)
		builder.(*array.Int32Builder).Append(v)
	}
}
