package broker

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-json"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
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
	Name     string
	ID       uuid.UUID
	Hostname string
	config   configuration.Configuration
	app      *fiber.App

	// grpc stuff
	storageGRPCConn   *grpc.ClientConn
	storageGRPCClient storageproto.StorageClient
}

func New(config configuration.Configuration) (*Broker, error) {
	// Create new Fiber application
	app := fiber.New(fiber.Config{
		AppName:           "broker-api-server",
		EnablePrintRoutes: true, // TODO: change this based on logger level -- DEBUG
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
	})
	// add default middleware
	app.Use(recover.New())

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	// get the hostname from the machine
	hn, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	br := &Broker{
		Name:     "broker",
		ID:       id,
		Hostname: hn,
		config:   config,
		app:      app,
	}

	br.setupRoutes()

	return br, nil
}

func (b *Broker) setupRoutes() {
	apiV1 := b.app.Group("/broker/v1")

	// query routes
	queryEndpoints := apiV1.Group("/query")
	queryEndpoints.Post("/", b.CreateQuery)
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
	b.storageGRPCConn, err = grpc.Dial(rpcAddr,
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
	var err error
	if err = b.initializeGRPCClient(); err != nil {
		return err
	}

	log.Info().Msg("Storage Server is ready to handle requests")
	return b.app.Listen(address)
}

func (b *Broker) ShutdownServer() error {
	// closing gRPC storage service client
	log.Info().Msg("close gRCP client connection of Storage Service")
	if err := b.storageGRPCConn.Close(); err != nil {
		return err
	}

	log.Info().Msg("shutting down server...")
	return b.app.Shutdown()
}

func (b *Broker) GetHost() string {
	hn, err := os.Hostname()
	if err != nil {
		panic(err.Error())
	}
	return hn
}

func (b *Broker) GetPort() string {
	return fmt.Sprint(b.config.Broker.Port)
}

func (b *Broker) GetName() string {
	return b.Name
}

func (b *Broker) GetID() string {
	return b.ID.String()
}

func (b *Broker) GetMetadata() map[string]string {
	return nil
}
