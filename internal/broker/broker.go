package broker

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/rs/zerolog/log"
	configuration "github.com/spaghettifunk/norman/internal/common"
	"github.com/spaghettifunk/norman/internal/common/utils"

	storageproto "github.com/spaghettifunk/norman/proto/v1/storage"
)

var (
	retriableErrors = []codes.Code{codes.Unavailable, codes.DataLoss}
	retryTimeout    = 50 * time.Millisecond
)

type Broker struct {
	config  configuration.Configuration
	server  *http.Server
	stopCtx context.Context
	stopFn  context.CancelFunc

	// grpc stuff
	storageGRPCConn   *grpc.ClientConn
	storageGRPCClient storageproto.StorageClient
}

func New(config configuration.Configuration) (*Broker, error) {
	stopCtx, stopFn := context.WithCancel(context.Background())
	return &Broker{
		config:  config,
		stopCtx: stopCtx,
		stopFn:  stopFn,
	}, nil
}

func (b *Broker) initializeGRPCClient() error {
	var err error

	unaryInterceptor := retry.UnaryClientInterceptor(
		retry.WithCodes(retriableErrors...),
		retry.WithMax(3),
		retry.WithBackoff(retry.BackoffLinear(retryTimeout)),
	)

	rpcAddr, err := utils.RPCAddr(b.config.Storage.BindAddr, b.config.Storage.RPCPort)
	if err != nil {
		return err
	}

	// initialize gRPC client of Storage Service
	log.Info().Msg("initialize gRPC client of Storage Service")
	b.storageGRPCConn, err = grpc.NewClient(rpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(unaryInterceptor),
	)
	if err != nil {
		return err
	}
	b.storageGRPCClient = storageproto.NewStorageClient(b.storageGRPCConn)

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	if _, err := b.storageGRPCClient.Ping(ctx, &storageproto.PingRequest{},
		retry.WithMax(3),
		retry.WithPerRetryTimeout(1*time.Second),
	); err != nil {
		return err
	}
	return err
}

func (b *Broker) StartServer(address string) error {
	if err := b.initializeGRPCClient(); err != nil {
		return err
	}
	log.Info().Msg("Broker Server is ready to handle gRPC requests")

	r := mux.NewRouter()
	r.HandleFunc("/broker/v1", b.apiVersion)
	r.HandleFunc("/broker/v1/query", b.runQuery)

	b.server = &http.Server{
		Addr:    address,
		Handler: r,
	}

	go func() {
		if err := b.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Info().Msgf("Server error: %v\n", err)
		}
	}()

	log.Info().Msgf("Broker server listening on %s\n", address)
	return nil

}

func (b *Broker) ShutdownServer() error {
	// closing gRPC storage service client
	log.Info().Msg("close gRCP client connection of Broker Service")
	if err := b.storageGRPCConn.Close(); err != nil {
		return fmt.Errorf("Storage gRPC client shutdown failed: %w", err)
	}

	b.stopFn()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Graceful shutdown timeout
	defer cancel()

	if err := b.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("Broker server shutdown failed: %w", err)
	}
	log.Info().Msg("Broker server stopped")
	return nil
}
