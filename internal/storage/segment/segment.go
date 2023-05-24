package segment

import (
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/spaghettifunk/norman/internal/common/schema"
	"github.com/spaghettifunk/norman/internal/common/types"
)

type Segment struct {
	ID         uuid.UUID `json:"-"`
	MapColumns map[string]*Column
	// count the inserted events
	counter uint32
}

func NewSegment(dir string, s *schema.Schema) (*Segment, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	segment := &Segment{
		ID:         id,
		MapColumns: map[string]*Column{},
	}

	for _, dimension := range s.DimensionFieldSpecs {
		col, err := NewColumn(dir, dimension.Name, types.DimensionType, types.GetDataType(dimension.Name, dimension.DataType))
		if err != nil {
			return nil, err
		}
		segment.MapColumns[dimension.Name] = col
	}

	for _, metric := range s.MetricFieldSpecs {
		col, err := NewColumn(dir, metric.Name, types.MetricType, types.GetDataType(metric.Name, metric.DataType))
		if err != nil {
			return nil, err
		}
		segment.MapColumns[metric.Name] = col
	}

	for _, dt := range s.DateTimeFieldSpecs {
		col, err := NewColumn(dir, dt.Name, types.TimeType, types.GetDataType(dt.Name, dt.DataType))
		if err != nil {
			return nil, err
		}
		segment.MapColumns[dt.Name] = col
	}

	return segment, nil
}

func (s *Segment) InsertData(column string, val interface{}) error {
	_, _, err := s.MapColumns[column].InsertData(val)
	// increase counter
	atomic.AddUint32(&s.counter, 1)
	return err
}

func (s *Segment) GetLength(colName string) int {
	return int(s.counter)
}

// Flush persist the segment on disk
func (s *Segment) Flush() error {
	for _, col := range s.MapColumns {
		if err := col.Flush(); err != nil {
			return err
		}
	}
	// reset counter now that we flushed data to disk
	atomic.SwapUint32(&s.counter, 0)
	return nil
}
