package model

import "encoding/json"

type RealtimeTable struct {
	Table
	RealtimeSegmentConfiguration *RealtimeSegmentConfiguration `json:"segmentConfiguration"`
}

type RealtimeSegmentConfiguration struct {
	SegmentConfiguration
	ReplicasPerPartition    int `json:"replicas"` // realtime only
	CompletionConfiguration struct {
		Mode string `json:""`
	} `json:"completitionConfiguration"` // realtime only
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
