package realtime_ingestion

import (
	"context"
	"fmt"
	"sync"

	"github.com/spaghettifunk/norman/internal/common/model"
	"github.com/spaghettifunk/norman/internal/storage/ingestion"
	"github.com/spaghettifunk/norman/pkg/realtime/kafka"
	"github.com/spaghettifunk/norman/pkg/realtime/kinesis"
)

// RealtimeIngestion defines the interface for the clients where the
// events are read from
type RealtimeIngestion interface {
	Initialize() error
	GetEvents(wg *sync.WaitGroup) context.CancelFunc
	Close() error
}

type Realtime struct {
	Type         ingestion.StreamType
	Client       RealtimeIngestion
	IngestionJob *model.IngestionJob
	Schema       *model.Schema
}

func New(ic *ingestion.StreamIngestionConfiguration) (*Realtime, error) {
	ingestor := &Realtime{
		Type: *ic.Type,
	}
	var err error
	switch *ic.Type {
	case ingestion.StreamKafka:
		ingestor.Client, err = kafka.NewIngestor(ic.KafkaConfiguration)
	case ingestion.StreamKinesis:
		ingestor.Client = kinesis.New()
	default:
		return nil, fmt.Errorf("invalid Realtime Ingestion type: %s", *ic.Type)
	}
	if err != nil {
		return nil, err
	}
	return ingestor, nil
}

// ReadEvents reads the incoming events
func (i *Realtime) ReadEvents() error {
	wg := &sync.WaitGroup{}
	cancel := i.Client.GetEvents(wg)
	cancel()
	wg.Wait()
	return i.Client.Close()
}
