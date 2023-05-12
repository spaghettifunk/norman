package kinesis

import (
	"context"
	"sync"
)

type KinesisIngestor struct {
	Address string
	Port    int
	// configuration here...
}

func New() *KinesisIngestor {
	return &KinesisIngestor{}
}

func (k *KinesisIngestor) Initialize() error {
	return nil
}

func (k *KinesisIngestor) GetEvents(wg *sync.WaitGroup) context.CancelFunc {
	return nil
}

func (k *KinesisIngestor) Close() error {
	return nil
}
