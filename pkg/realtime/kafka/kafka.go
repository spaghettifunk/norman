package kafka

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/Shopify/sarama"
	"github.com/spaghettifunk/norman/internal/common/segment"
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
	wg       *sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	// segment manager
	// ....
}

func NewIngestor(kcfg *KafkaConfiguration) (*KafkaIngestor, error) {
	bs := strings.Split(kcfg.Brokers, ",")

	cfg := createConfiguration(kcfg)
	consumer, err := sarama.NewConsumerGroup(bs, kcfg.ConsumerGroup, cfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &KafkaIngestor{
		Consumer: consumer,
		Topic:    kcfg.Topic,
		Brokers:  bs,
		ready:    make(chan bool),
		wg:       &sync.WaitGroup{},
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

func createConfiguration(kcfg *KafkaConfiguration) *sarama.Config {
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
	k.wg.Add(1)
	go func() {
		defer k.wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := k.Consumer.Consume(k.ctx, []string{k.Topic}, k); err != nil {
				log.Panic().Msgf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if k.ctx.Err() != nil {
				return
			}
			k.ready = make(chan bool)
		}
	}()

	<-k.ready // Await till the consumer has been set up
	log.Info().Msg("Sarama consumer up and running!...")

	return nil
}

// GetEvents needs to be called from a goroutine
func (k *KafkaIngestor) GetEvents() error {
	keepRunning := true
	consumptionIsPaused := false

	// TODO: this should come from an API
	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	// TODO: this should come from an API
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-k.ctx.Done():
			log.Info().Msg("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Info().Msg("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			k.toggleConsumptionFlow(&consumptionIsPaused)
		}
	}

	k.cancel()
	k.wg.Wait()

	return nil
}

func (k *KafkaIngestor) toggleConsumptionFlow(isPaused *bool) {
	if *isPaused {
		k.Consumer.ResumeAll()
		log.Info().Msg("resuming consumption")
	} else {
		k.Consumer.PauseAll()
		log.Info().Msg("pausing consumption")
	}
	*isPaused = !*isPaused
}

// Shutdown closes the producer object
func (k *KafkaIngestor) Shutdown(failure bool) error {
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
			// TODO: create segment here
			buffer.InsertRow(nil)
			if len(buffer.Rows) > MaxBatchedEventsPerSegment {
				// Flush Segment
				// TODO: dir should come from Commander/Aqua
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
				// TODO: dir should come from Commander/Aqua
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
