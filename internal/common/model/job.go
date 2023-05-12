package model

import (
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/spaghettifunk/norman/internal/storage/ingestion"
)

type IngestionType string

const (
	Offline  IngestionType = "OFFLINE"
	Realtime IngestionType = "REALTIME"
)

type IngestionJob struct {
	ID                     uuid.UUID               `json:"-"`
	Name                   string                  `json:"name"`
	Type                   IngestionType           `json:"type"`
	IndexConfiguration     *IndexConfiguration     `json:"indexConfiguration,omitempty"`
	TenantConfiguration    *TenantConfiguration    `json:"tenantConfiguration,omitempty"`
	Metadata               map[string]interface{}  `json:"metadata,omitempty"`
	IngestionConfiguration *IngestionConfiguration `json:"ingestionConfiguration,omitempty"`
	SegmentConfiguration   *SegmentConfiguration   `json:"segmentConfiguration"`
}

type SegmentConfiguration struct {
	SchemaName          string `json:"schemaName"`
	TimeColumnName      string `json:"timeColumnName"`
	TimeType            string `json:"timeType"`
	AllowNullTimeValue  bool   `json:"allowNullTimeValue,omitempty"`
	RetentionTimeUnit   string `json:"retentionTimeUnit,omitempty"`
	RetentionTimeValue  int    `json:"retentionTimeValue,omitempty"`
	PushFrequency       string `json:"pushFrequency,omitempty"`
	PushType            string `json:"pushType,omitempty"`
	NullHandlingEnabled bool   `json:"nullHandlingEnabled,omitempty"`
	// offline only
	Replication int `json:"replication,omitempty"`
	// realtime only
	ReplicasPerPartition    int `json:"replicas,omitempty"`
	CompletionConfiguration struct {
		Mode string `json:"mode,omitempty"`
	} `json:"completitionConfiguration,omitempty"`
}

type IngestionConfiguration struct {
	FilterConfiguration struct {
		FilterFunction string `json:"filterFunction,omitempty"`
	} `json:"filterConfiguration,omitempty"`
	TransformConfigurations []struct {
		ColumnName string `json:"columnName,omitempty"`
		Function   string `json:"function,omitempty"`
	} `json:"transformConfigurations,omitempty"`
	StreamIngestionConfiguration *ingestion.StreamIngestionConfiguration `json:"streamIngestionConfigurations,omitempty"`
}

// IndexConfiguration is used to set the way the data is indexed in the Storage
type IndexConfiguration struct {
	InvertedIndexColumns                       []string               `json:"invertedIndexColumns,omitempty"`
	CreateInvertedIndexDuringSegmentGeneration bool                   `json:"createInvertedIndexDuringSegmentGeneration,omitempty"`
	SortedColumn                               []string               `json:"sortedColumn,omitempty"`
	BloomFilterColumns                         []string               `json:"bloomFilterColumns,omitempty"`
	StarTreeIndexConfigs                       map[string]interface{} `json:"starTreeIndexConfigs,omitempty"`
	NoDictionaryColumns                        []string               `json:"noDictionaryColumns,omitempty"`
	RangeIndexColumns                          []string               `json:"rangeIndexColumns,omitempty"`
	OnHeapDictionaryColumns                    []string               `json:"onHeapDictionaryColumns,omitempty"`
	VarLengthDictionaryColumns                 []string               `json:"varLengthDictionaryColumns,omitempty"`
	SegmentPartitionConfiguration              struct {
		Column struct {
			Name          string `json:"name,omitempty"`
			FunctionName  string `json:"functionName,omitempty"`
			NumPartitions int    `json:"numPartitions,omitempty"`
		} `json:"column,omitempty"`
	} `json:"segmentPartitionConfiguration,omitempty"`
}

type TenantConfiguration struct {
	// BrokerLabel is the name of the broker instance
	BrokerLabel string
	// StorageLabel is the name of the storage instance
	StorageLabel string
}

type JobStatus string

const (
	JobCreated       JobStatus = "CREATED"
	JobStarted       JobStatus = "STARTED"
	JobFailed        JobStatus = "FAILED"
	JobPartiallyDone JobStatus = "PARTIALLY_DONE"
	JobDone          JobStatus = "DONE"
)

// NewIngestionJob initiate a new job that will load events either
// from an offline method or realtime
func NewIngestionJob(config []byte) (*IngestionJob, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	ij := &IngestionJob{ID: id}
	if err := json.Unmarshal(config, ij); err != nil {
		return nil, err
	}
	return ij, nil
}
