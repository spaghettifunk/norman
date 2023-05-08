package model

import "encoding/json"

type OfflineTable struct {
	Table
	OfflineSegmentConfiguration *OfflineSegmentConfiguration `json:"segmentConfiguration"`
}

type OfflineSegmentConfiguration struct {
	SegmentConfiguration
	Replication int `json:"replication"` // offline only
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
