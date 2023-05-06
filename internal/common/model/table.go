package model

import "encoding/json"

type TableType string

const (
	Offline  TableType = "OFFLINE"
	Realtime TableType = "REALTIME"
)

type Table struct {
	Name                   string                  `json:"name"`
	Type                   TableType               `json:"type"`
	IndexConfiguration     *IndexConfiguration     `json:"indexConfiguration,omitempty"`
	TenantConfiguration    *TenantConfiguration    `json:"tenantConfiguration,omitempty"`
	Metadata               map[string]interface{}  `json:"metadata,omitempty"`
	IngestionConfiguration *IngestionConfiguration `json:"ingestionConfiguration,omitempty"`
}

type OfflineTable struct {
	Table
	OfflineSegmentConfiguration *OfflineSegmentConfiguration `json:"segmentConfiguration"`
}

type RealtimeTable struct {
	Table
	RealtimeSegmentConfiguration *RealtimeSegmentConfiguration `json:"segmentConfiguration"`
}

// --- Segment Configurations ---
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
}

type OfflineSegmentConfiguration struct {
	SegmentConfiguration
	Replication int `json:"replication"` // offline only
}

type RealtimeSegmentConfiguration struct {
	SegmentConfiguration
	ReplicasPerPartition    int `json:"replicas"` // realtime only
	CompletionConfiguration struct {
		Mode string `json:""`
	} `json:"completitionConfiguration"` // realtime only
}

// --- Ingestion Configurations ---
type IngestionConfiguration struct {
	FilterConfiguration struct {
		FilterFunction string `json:"filterFunction"`
	} `json:"filterConfiguration"`
	TransformConfigurations []struct {
		ColumnName string `json:"columnName"`
		Function   string `json:"function"`
	} `json:"transformConfigurations"`
}

// --- Common Configurations ---

// IndexConfiguration is used to set the way the data is indexed in the Storage
type IndexConfiguration struct {
	InvertedIndexColumns                       []string               `json:"invertedIndexColumns"`
	CreateInvertedIndexDuringSegmentGeneration bool                   `json:"createInvertedIndexDuringSegmentGeneration"`
	SortedColumn                               []string               `json:"sortedColumn"`
	BloomFilterColumns                         []string               `json:"bloomFilterColumns"`
	StarTreeIndexConfigs                       map[string]interface{} `json:"starTreeIndexConfigs"`
	NoDictionaryColumns                        []string               `json:"noDictionaryColumns"`
	RangeIndexColumns                          []string               `json:"rangeIndexColumns"`
	OnHeapDictionaryColumns                    []string               `json:"onHeapDictionaryColumns"`
	VarLengthDictionaryColumns                 []string               `json:"varLengthDictionaryColumns"`
	SegmentPartitionConfiguration              struct {
		PrimaryKey struct {
			FunctionName  string `json:"functionName"`
			NumPartitions int    `json:"numPartitions"`
		} `json:"primarykey"`
	} `json:"segmentPartitionConfiguration"`
}

type TenantConfiguration struct {
	// BrokerLabel is the name of the broker instance
	BrokerLabel string
	// StorageLabel is the name of the storage instance
	StorageLabel string
}

// NewOfflineTable creates a new Table that is used to hold the segments
// A table is created in two occasions: via an API request or from the metadata
// retrieved from Aqua
func NewOfflineTable(config []byte) (*OfflineTable, error) {
	ot := &OfflineTable{}
	if err := json.Unmarshal(config, ot); err != nil {
		return nil, err
	}
	return ot, nil
}

// NewRealtimeTable creates a new Table that is used to hold the segments
// A table is created in two occasions: via an API request or from the metadata
// retrieved from Aqua
func NewRealtimeTable(config []byte) (*RealtimeTable, error) {
	rt := &RealtimeTable{}
	if err := json.Unmarshal(config, rt); err != nil {
		return nil, err
	}
	return rt, nil
}
