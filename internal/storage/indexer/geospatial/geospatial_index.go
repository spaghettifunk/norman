package geospatial

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/spaghettifunk/norman/internal/storage/indexer"
	rtreego "github.com/spaghettifunk/norman/pkg/containers/rtree"
	"github.com/twpayne/go-geom"
)

type metadata struct {
	IndexType  indexer.IndexType `json:"type"`
	ColumnName string            `json:"column"`
}

type GeospatialIndex[G geom.T] struct {
	Metadata metadata      `json:"metadata"`
	Index    rtreego.Rtree `json:"index"`
}

type indexItem[G geom.T] struct {
	IDs      []uint32 `json:"ids"`
	Location G        `json:"location"`
}

// New creates a new Geospatial index object
func New[G geom.T](columnName string, dimensions int) *GeospatialIndex[G] {
	return &GeospatialIndex[G]{
		Metadata: metadata{
			IndexType:  indexer.GeospatialIndex,
			ColumnName: columnName,
		},
	}
}

func (i *GeospatialIndex[G]) AddValue(id string, value interface{}) bool {
	_, ok := value.(geom.T)
	if !ok {
		log.Error().Msg("value cannot be casted to ValidType")
		return false
	}

	return true
}

func (i *GeospatialIndex[G]) Search(value interface{}) []uint32 {
	_, ok := value.(float64)
	if !ok {
		log.Error().Msg("value cannot be casted to ValidType")
		return nil
	}
	return nil
}

func (i *GeospatialIndex[G]) GetColumnName() string {
	return i.Metadata.ColumnName
}

func (i *GeospatialIndex[G]) GetIndexType() indexer.IndexType {
	return i.Metadata.IndexType
}

type SortedIndexJSON[G geom.T] struct {
	Metadata metadata `json:"metadata"`
}

func (i *GeospatialIndex[G]) MarshalJSON() ([]byte, error) {
	si := &SortedIndexJSON[G]{
		Metadata: i.Metadata,
	}
	return json.Marshal(si)
}

// TODO: implement this once there is an actual index system
func (i *GeospatialIndex[G]) Deserialize(data []byte) error {
	return nil
}
