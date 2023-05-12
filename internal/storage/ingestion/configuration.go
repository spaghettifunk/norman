package ingestion

type IngestionConfiguration struct {
	KafkaConfiguration   *KafkaConfiguration   `json:"kafka"`
	KinesisConfiguration *KinesisConfiguration `json:"kinesis"`
}

type KafkaConfiguration struct {
	// Brokers are the Kafka bootstrap brokers to connect to, as a comma separated list
	Brokers string `json:"brokers"`
	// Topic is the Kafka topic to be consumed
	Topic string `json:"topic"`
	// ConsumerGroup is the Kafka consumer group definition
	ConsumerGroup string `json:"consumerGroup"`
	// Kafka cluster version
	Version string `json:"version"`
	// Assignor is the Consumer group partition assignment strategy (range, roundrobin, sticky)
	Assignor string `json:"assignor"`
	// InitialOffset is the Kafka consumer consume initial offset from oldest
	InitialOffset  string               `json:"offset"`
	Authentication *KafkaAuthentication `json:"authentication,omitempty"`
}

type KafkaAuthentication struct {
	// values accepted "OAUTHBEARER", "PLAIN", "SCRAM-SHA-256", "SCRAM-SHA-512", "GSSAPI"
	SASLMechanism string `json:"saslmechanism"`
	Username      string `json:"username"`
	Password      string `json:"password"`
}

type KinesisConfiguration struct {
	Address string
}
