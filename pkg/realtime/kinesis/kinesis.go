package kinesis

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

func (k *KinesisIngestor) GetEvents() error {
	return nil
}

func (k *KinesisIngestor) Shutdown(failure bool) error {
	return nil
}
