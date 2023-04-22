package kinesis

type KinesisIngestor struct {
	Address string
	Port    int
	// configuration here...
}

func New() *KinesisIngestor {
	return &KinesisIngestor{}
}

func (k *KinesisIngestor) GetLog(log []byte) error {
	return nil
}
