package ingestion

import (
	"github.com/google/uuid"

	"github.com/spaghettifunk/norman/internal/storage/ingestion/realtime/kafka"
	"github.com/spaghettifunk/norman/internal/storage/ingestion/realtime/kinesis"
)

type IngestionType string

const (
	OfflineLocalStorage     IngestionType = "LOCAL"
	OfflineGCPStorage       IngestionType = "GCP_CLOUD_STORAGE"
	OfflineAzureBlobStorage IngestionType = "AZURE_BLOB_STORAGE"
	OfflineS3               IngestionType = "AWS_S3"
	OfflineHDFS             IngestionType = "HDFS"
	StreamKafka             IngestionType = "KAFKA"
	StreamKinesis           IngestionType = "KINESIS"
)

type IngestionJobConfiguration struct {
	ID                     uuid.UUID               `json:"id,omitempty"`
	Name                   string                  `json:"name"`
	Type                   IngestionType           `json:"type"`
	IndexConfiguration     *IndexConfiguration     `json:"indexConfiguration,omitempty"`
	TenantConfiguration    *tenantConfiguration    `json:"tenantConfiguration,omitempty"`
	Metadata               map[string]interface{}  `json:"metadata,omitempty"`
	IngestionConfiguration *ingestionConfiguration `json:"ingestionConfiguration"`
	SegmentConfiguration   *segmentConfiguration   `json:"segmentConfiguration"`
}

type segmentConfiguration struct {
	TableName           string `json:"tableName"`
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

type ingestionConfiguration struct {
	FilterConfiguration struct {
		FilterFunction string `json:"filterFunction,omitempty"`
	} `json:"filterConfiguration,omitempty"`
	TransformConfigurations []struct {
		ColumnName string `json:"columnName,omitempty"`
		Function   string `json:"function,omitempty"`
	} `json:"transformConfigurations,omitempty"`
	Offline struct {
	} `json:"offline,omitempty"`
	Realtime struct {
		KafkaConfiguration   *kafka.KafkaConfiguration     `json:"kafka,omitempty"`
		KinesisConfiguration *kinesis.KinesisConfiguration `json:"kinesis,omitempty"`
	} `json:"realtime,omitempty"`
}

// IndexConfiguration is used to set the way the data is indexed in the Storage
type IndexConfiguration struct {
	TextInvertedIndexColumns                   []string               `json:"textInvertedIndexColumns,omitempty"`
	CreateInvertedIndexDuringSegmentGeneration bool                   `json:"createInvertedIndexDuringSegmentGeneration,omitempty"`
	SortedIndexColumn                          []string               `json:"sortedIndexColumn,omitempty"`
	RangeIndexColumns                          []string               `json:"rangeIndexColumns,omitempty"`
	BitmapIndexColumns                         []string               `json:"bitmapIndexColumns,omitempty"`
	GeospatialIndexColumns                     []string               `json:"geospatialIndexColumns,omitempty"`
	BloomFilterColumns                         []string               `json:"bloomFilterColumns,omitempty"`
	StarTreeIndexConfigs                       map[string]interface{} `json:"starTreeIndexConfigs,omitempty"`
	SegmentPartitionConfiguration              struct {
		Column struct {
			Name          string `json:"name,omitempty"`
			FunctionName  string `json:"functionName,omitempty"`
			NumPartitions int    `json:"numPartitions,omitempty"`
		} `json:"column,omitempty"`
	} `json:"segmentPartitionConfiguration,omitempty"`
}

type tenantConfiguration struct {
	// BrokerLabel is the name of the broker instance
	BrokerLabel string
	// StorageLabel is the name of the storage instance
	StorageLabel string
}

type JobStatus string

const (
	JobCreated       JobStatus = "CREATED"
	JobInitialized   JobStatus = "INITIALIZED"
	JobInProgress    JobStatus = "PROGRESS"
	JobFailed        JobStatus = "FAILED"
	JobPartiallyDone JobStatus = "PARTIALLY_DONE"
	JobDone          JobStatus = "DONE"
)

// NewIngestionJob initiate a new job that will load events either
// from an offline method or realtime
func NewIngestionJob(ij *IngestionJobConfiguration) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	ij.ID = id
	return nil
}
