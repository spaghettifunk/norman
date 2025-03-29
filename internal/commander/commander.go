package commander

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	configuration "github.com/spaghettifunk/norman/internal/common"
	"github.com/spaghettifunk/norman/internal/common/utils"
	storageproto "github.com/spaghettifunk/norman/proto/v1/storage"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	retriableErrors = []codes.Code{codes.Unavailable, codes.DataLoss}
	retryTimeout    = 50 * time.Millisecond
)

type Commander struct {
	config  configuration.Configuration
	server  *http.Server
	stopCtx context.Context
	stopFn  context.CancelFunc

	// grpc stuff
	storageGRPCConn   *grpc.ClientConn
	storageGRPCClient storageproto.StorageClient
}

func New(config configuration.Configuration) (*Commander, error) {
	stopCtx, stopFn := context.WithCancel(context.Background())
	return &Commander{
		config:  config,
		stopCtx: stopCtx,
		stopFn:  stopFn,
	}, nil
}

func (c *Commander) StartServer(address string) error {
	if err := c.initializeGRPCClient(); err != nil {
		return err
	}
	// initialize api
	log.Info().Msg("Commander server is ready to handle grRPC requests")

	r := mux.NewRouter()
	r.HandleFunc("/commander/v1", c.apiVersion)

	// tenant routes
	r.HandleFunc("/commander/v1/tenants", c.apiVersion)

	// tenantEndpoints := apiV1.Group("/tenants")
	// tenantEndpoints.Get("/", c.GetTenants)
	// tenantEndpoints.Post("/", c.CreateTenant)
	// tenantEndpoints.Get("/:tenantId", c.GetTenant)
	// tenantEndpoints.Put("/:tenantId", c.UpdateTenant)
	// tenantEndpoints.Patch("/:tenantId", c.PatchTenant)
	// tenantEndpoints.Delete("/:tenantId", c.DeleteTenant)

	// // table endpoints
	// tableEndpoints := tenantEndpoints.Group("/:tenantId/tables")
	// tableEndpoints.Get("/", c.GetTables)
	// tableEndpoints.Post("/", c.CreateTable)
	// tableEndpoints.Get("/:tableName", c.GetTable)
	// tableEndpoints.Put("/:tableName", c.UpdateTable)
	// tableEndpoints.Patch("/:tableName", c.PatchTable)
	// tableEndpoints.Delete("/:tableName", c.DeleteTable)

	// // ingestion job endpoints
	// jobEndpoints := tenantEndpoints.Group("/:tenantId/jobs")
	// jobEndpoints.Get("/", c.GetJobs)
	// jobEndpoints.Post("/", c.CreateJob)
	// jobEndpoints.Get("/:jobID", c.GetJob)
	// jobEndpoints.Put("/:jobID", c.UpdateJob)
	// jobEndpoints.Patch("/:jobID", c.PatchJob)
	// jobEndpoints.Delete("/:jobID", c.DeleteJob)

	c.server = &http.Server{
		Addr:    address,
		Handler: r,
	}

	go func() {
		if err := c.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Info().Msgf("Server error: %v\n", err)
		}
	}()

	log.Info().Msgf("Broker server listening on %s\n", address)

	return nil
}

func (c *Commander) initializeGRPCClient() error {
	var err error
	unaryInterceptor := retry.UnaryClientInterceptor(
		retry.WithCodes(retriableErrors...),
		retry.WithMax(3),
		retry.WithBackoff(retry.BackoffLinear(retryTimeout)),
	)

	rpcAddr, err := utils.RPCAddr(c.config.Storage.BindAddr, c.config.Storage.RPCPort)
	if err != nil {
		return err
	}

	// initialize gRPC client of Storage Service
	log.Info().Msg("initialize gRPC client of Storage Service")
	c.storageGRPCConn, err = grpc.NewClient(rpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(unaryInterceptor),
	)
	if err != nil {
		return err
	}
	c.storageGRPCClient = storageproto.NewStorageClient(c.storageGRPCConn)

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	if _, err := c.storageGRPCClient.Ping(ctx, &storageproto.PingRequest{},
		retry.WithMax(3),
		retry.WithPerRetryTimeout(1*time.Second),
	); err != nil {
		return err
	}
	return err
}

func (c *Commander) ShutdownServer() error {
	// closing gRPC storage service client
	log.Info().Msg("close gRCP client connection of Storage Service")
	if err := c.storageGRPCConn.Close(); err != nil {
		return fmt.Errorf("Storage gRPC client shutdown failed: %w", err)
	}

	c.stopFn()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Graceful shutdown timeout
	defer cancel()

	if err := c.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("Broker server shutdown failed: %w", err)
	}
	log.Info().Msg("Broker server stopped")
	return nil
}
