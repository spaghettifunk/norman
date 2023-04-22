package realtime

import (
	"github.com/spaghettifunk/norman/pkg/realtime/kafka"
	"github.com/spaghettifunk/norman/pkg/realtime/kinesis"
)

type IngestionType string

const (
	// Kafka represents the Kafka type for realtime ingestion
	Kafka IngestionType = "kafka"
	// Kinesis represents the Kinesis type for realtime ingestion
	Kinesis IngestionType = "kinesis"
)

type RealtimeIngestion interface {
	GetLog(log []byte) error
}

type Realtime struct {
	Type   IngestionType
	Client RealtimeIngestion
	// configuration here...
}

func New(t IngestionType) (*Realtime, error) {
	ingestor := &Realtime{
		Type: t,
	}
	switch t {
	case Kafka:
		ingestor.Client = kafka.New()
	case Kinesis:
		ingestor.Client = kinesis.New()
	default:
		panic("invalid type")
	}
	return ingestor, nil
}

// ReadLog reads the incoming string
func (i *Realtime) ReadLog(data string) error {
	err := i.Client.GetLog([]byte(data))
	if err != nil {
		return err
	}
	return nil
}
