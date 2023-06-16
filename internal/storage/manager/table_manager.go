package manager

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
	"github.com/spaghettifunk/norman/internal/common/entities"
	"github.com/spaghettifunk/norman/internal/storage/segment"
)

const (
	eventIDName string = "_normanID"
)

type TableManager struct {
	Table             *entities.Table
	baseDir           string
	datetimeFieldName string
	activeSegment     *segment.Segment
	segments          []*segment.Segment
	builder           *array.RecordBuilder
	wg                sync.WaitGroup
	partition         int
	partitionStart    time.Time
	interval          time.Duration
	granularity       *entities.GranularitySpec
	eventsCounter     int
}

var (
	minTimestamp = time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
)

var buffPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

func NewTableManager(table *entities.Table) (*TableManager, error) {
	// TODO: this should depend on a folder that comes from Configuration
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// TODO: how to get the tenant name? --> default for now
	// Format: os_path + output/{tenantID}/
	baseDir := fmt.Sprintf("%s/output/default/%s", path, table.Name)
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		return nil, err
	}
	// extract the datetime field
	dtField := table.GetDatetimeField()

	granularity, err := table.Schema.GetGranularity()
	if err != nil {
		return nil, err
	}

	// create the new record builder for inserting data to arrow file
	mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
	b := array.NewRecordBuilder(mem, table.EventSchema)

	return &TableManager{
		Table:             table,
		builder:           b,
		datetimeFieldName: dtField.Name,
		baseDir:           baseDir,
		wg:                sync.WaitGroup{},
		partition:         0,
		interval:          time.Duration(granularity.Size) * granularity.UnitSpec,
		granularity:       granularity,
		eventsCounter:     0,
	}, nil
}

func (t *TableManager) CreateNewSegment() error {
	d := time.Now()

	fPath := fmt.Sprintf("%s/%s", t.baseDir, d.Format("2006-01-02T15:04:05"))
	s, err := segment.NewSegment(fPath, t.partition, t.Table.EventSchema)
	if err != nil {
		return err
	}
	t.activeSegment = s
	return nil
}

func (t *TableManager) InsertData(data []byte) error {
	event := make(map[string]interface{}, len(t.Table.EventSchema.Fields()))
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	// Get the timestamp of the consumed message
	dtVal := int64(event[t.datetimeFieldName].(float64))
	eventTimestamp := time.Unix(0, dtVal*int64(time.Millisecond))

	// initial partition time setup
	if t.partitionStart.Equal(minTimestamp) {
		t.partitionStart = eventTimestamp.Truncate(t.granularity.UnitSpec)
	}

	partitionInterval := t.partitionStart.Add(t.interval)

	// add norman ID to the event
	event[eventIDName] = t.generateEventID(partitionInterval)

	// if the interval is passed then it creates a new partition
	if !(eventTimestamp.After(t.partitionStart) && eventTimestamp.Before(partitionInterval)) {
		if err := t.FlushSegment(); err != nil {
			return err
		}
		if err := t.CreateNewSegment(); err != nil {
			return err
		}
		// Update the current partition and its start time
		t.partition++
		t.partitionStart = eventTimestamp.Truncate(t.granularity.UnitSpec)
	}

	t.processEvent(event)
	t.eventsCounter++

	return nil
}

func (t *TableManager) generateEventID(partitionInterval time.Time) string {
	bf := buffPool.Get().(*bytes.Buffer)
	defer buffPool.Put(bf)
	bf.Reset()

	bf.WriteString(t.Table.Name)
	bf.WriteString("_P")
	bf.WriteString(strconv.Itoa(t.partition))
	bf.WriteByte('_')
	bf.WriteString(t.partitionStart.Format("2006-01-02T15:04:05:00"))
	bf.WriteByte('_')
	bf.WriteString(partitionInterval.Format("2006-01-02T15:04:05:00"))
	bf.WriteByte('_')
	bf.WriteString(strconv.Itoa(t.eventsCounter))

	return bf.String()
}

// this is processed concurrently considering that there can be hundreds of columns per event
func (t *TableManager) processEvent(event map[string]interface{}) {
	for idx, field := range t.builder.Schema().Fields() {
		val, ok := event[field.Name]
		if !ok {
			log.Error().Msgf("could not find column %s in builder", field.Name)
			return
		}
		b := t.builder.Field(idx)
		t.appendValue(val, field, b)
	}
}

// FlushSegment first persist on disk the current segment
// secondly, it compresses the segment to save space and lastly
// it reset the memory object so that it can start over
func (t *TableManager) FlushSegment() error {
	record := t.builder.NewRecord()
	if err := t.activeSegment.Flush(record); err != nil {
		return err
	}
	// store the active segments in the list of segments
	t.segments = append(t.segments, t.activeSegment)
	return nil
}

// val interface{} when Unmarshalled for numbers is always float64. Initial
// assertion is necessary before the correct type casting
func (t *TableManager) appendValue(val interface{}, field arrow.Field, builder array.Builder) {
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
