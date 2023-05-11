package kafka

type KafkaIngestor struct {
	Address string
	Port    int
	// configuration here...
}

func New() *KafkaIngestor {
	return &KafkaIngestor{}
}

func (k *KafkaIngestor) GetEvent() error {
	return nil
}
