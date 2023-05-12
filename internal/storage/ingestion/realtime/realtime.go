package realtime_ingestion

import (
	"context"
	"fmt"
	"sync"

	"github.com/spaghettifunk/norman/internal/storage/ingestion"
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

// RealtimeIngestion defines the interface for the clients where the
// events are read from
type RealtimeIngestion interface {
	Initialize() error
	GetEvents(wg *sync.WaitGroup) context.CancelFunc
	Close() error
}

type Realtime struct {
	Type   IngestionType
	Client RealtimeIngestion
	// configuration here...
}

func New(t IngestionType, ic *ingestion.IngestionConfiguration) (*Realtime, error) {
	ingestor := &Realtime{
		Type: t,
	}
	var err error
	switch t {
	case Kafka:
		ingestor.Client, err = kafka.NewIngestor(ic.KafkaConfiguration)
	case Kinesis:
		ingestor.Client = kinesis.New()
	default:
		return nil, fmt.Errorf("invalid Realtime Ingestion type: %s", t)
	}
	if err != nil {
		return nil, err
	}
	return ingestor, nil
}

// ReadLog reads the incoming string
func (i *Realtime) ReadLog() error {
	wg := &sync.WaitGroup{}
	cancel := i.Client.GetEvents(wg)
	cancel()
	wg.Wait()
	return nil
}
