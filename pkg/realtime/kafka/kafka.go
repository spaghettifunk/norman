package kafka

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/Shopify/sarama"
	"github.com/spaghettifunk/norman/internal/common/segment"
	"github.com/spaghettifunk/norman/internal/storage/ingestion"
)

var (
	// MaxNumberOfWorkers is the max number of concurrent goroutines for uploading data
	MaxNumberOfWorkers = runtime.NumCPU()
	// FlushIntervalInSec is the amount of time before executing the Flush operation in case the buffer is not full
	FlushIntervalInSec = 10
	// MaxBatchedEvents is the maximum amount of events in a segment
	MaxBatchedEventsPerSegment = 100000
)

type KafkaIngestor struct {
	Consumer sarama.ConsumerGroup
	Topic    string
	Brokers  []string
	ready    chan bool
}

func NewIngestor(kcfg *ingestion.KafkaConfiguration) (*KafkaIngestor, error) {
	bs := strings.Split(kcfg.Brokers, ",")

	cfg := createConfiguration(kcfg)
	consumer, err := sarama.NewConsumerGroup(bs, kcfg.ConsumerGroup, cfg)
	if err != nil {
		return nil, err
	}
	return &KafkaIngestor{
		Consumer: consumer,
		Topic:    kcfg.Topic,
		Brokers:  bs,
		ready:    make(chan bool),
	}, nil
}

func createConfiguration(kcfg *ingestion.KafkaConfiguration) *sarama.Config {
	cfg := sarama.NewConfig()

	// TODO: these should be specified via the ingestion job config
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Offsets.AutoCommit.Enable = true

	switch kcfg.Assignor {
	case "sticky":
		cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "roundrobin":
		cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range":
		cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		log.Panic().Msgf("Unrecognized consumer group partition assignor: %s", kcfg.Assignor)
	}

	if kcfg.InitialOffset == "oldest" {
		cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	// TODO: refactor this in a more robust way
	if kcfg.Authentication != nil {
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.Mechanism = sarama.SASLMechanism(kcfg.Authentication.SASLMechanism)
		cfg.Net.SASL.User = kcfg.Authentication.Username
		cfg.Net.SASL.Password = kcfg.Authentication.Password
	}

	return cfg
}

func (k *KafkaIngestor) Initialize() error {
	return nil
}

func (k *KafkaIngestor) GetEvents(wg *sync.WaitGroup) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := k.Consumer.Consume(ctx, []string{k.Topic}, k); err != nil {
				log.Panic().Msgf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			k.ready = make(chan bool)
		}
	}()

	<-k.ready // Await till the consumer has been set up
	log.Info().Msg("Sarama consumer up and running!...")
	return cancel
}

// Close closes the producer object
func (k *KafkaIngestor) Close() error {
	return k.Consumer.Close()
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (k *KafkaIngestor) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(k.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (k *KafkaIngestor) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (k *KafkaIngestor) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var fflush time.Time
	buffer, err := segment.NewSegment()
	if err != nil {
		return err
	}

	// TODO: the FlushIntervalInSec should be coming from the Job Ingestion Configuration
	// Replace this with the Granularity of the segment
	flushInterval := time.Duration(FlushIntervalInSec) * time.Second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case message := <-claim.Messages():
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)

			// TODO: create segment here
			buffer.InsertRow(nil)

			if len(buffer.Rows) > MaxBatchedEventsPerSegment {
				// Flush Segment
				if err := buffer.Flush(fmt.Sprintf("/tmp/norman/%s/%s", k.Topic, time.Now().String()), true); err != nil {
					log.Panic().Err(err).Msg("error in flushing the segment")
				}
			}

			session.MarkMessage(message, "")

		case <-ticker.C:
			// Refresh pipe
			tt := time.Now()
			if tt.After(fflush) {
				log.Debug().Msg("Force flush (interval) triggered")

				// Flush Segment
				if err := buffer.Flush(fmt.Sprintf("/tmp/norman/%s/%s", k.Topic, time.Now().String()), true); err != nil {
					log.Panic().Err(err).Msg("error in flushing the segment")
				}

				fflush = tt.Add(flushInterval)
				log.Debug().Msg("Force flush (interval) finished")
			}

		case <-session.Context().Done():
			return nil
		}
	}
}
