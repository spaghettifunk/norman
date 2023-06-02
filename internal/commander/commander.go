package commander

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/rs/zerolog/log"

	configuration "github.com/spaghettifunk/norman/internal/common"
	"github.com/spaghettifunk/norman/internal/common/utils"
	"github.com/spaghettifunk/norman/pkg/consul"
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
	Name     string
	ID       uuid.UUID
	Hostname string
	consul   *consul.Consul
	config   configuration.Configuration
	app      *fiber.App

	// grpc stuff
	storageGRPCConn   *grpc.ClientConn
	storageGRPCClient storageproto.StorageClient
}

func New(config configuration.Configuration) (*Commander, error) {
	// Create new Fiber application
	app := fiber.New(fiber.Config{
		AppName:           "commander-api-server",
		EnablePrintRoutes: true, // TODO: change this based on logger level -- DEBUG
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
	})
	// add default middleware
	app.Use(recover.New())

	// initialize consul client
	cs := consul.New()
	if err := cs.Init(); err != nil {
		return nil, err
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	// get the hostname from the machine
	hn, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	c := &Commander{
		Name:     "commander",
		ID:       id,
		Hostname: hn,
		consul:   cs,
		config:   config,
		app:      app,
	}

	c.setupRoutes()

	return c, nil
}

func (c *Commander) setupRoutes() {
	apiV1 := c.app.Group("/commander/v1")

	apiV1.Get("/", c.APIVersion)
	apiV1.Get("/swagger/*", swagger.HandlerDefault)

	// tenant routes
	tenantEndpoints := apiV1.Group("/tenants")
	tenantEndpoints.Get("/", c.GetTenants)
	tenantEndpoints.Post("/", c.CreateTenant)
	tenantEndpoints.Get("/:tenantId", c.GetTenant)
	tenantEndpoints.Put("/:tenantId", c.UpdateTenant)
	tenantEndpoints.Patch("/:tenantId", c.PatchTenant)
	tenantEndpoints.Delete("/:tenantId", c.DeleteTenant)

	// table endpoints
	tableEndpoints := tenantEndpoints.Group("/:tenantId/tables")
	tableEndpoints.Get("/", c.GetTables)
	tableEndpoints.Post("/", c.CreateTable)
	tableEndpoints.Get("/:tableName", c.GetTable)
	tableEndpoints.Put("/:tableName", c.UpdateTable)
	tableEndpoints.Patch("/:tableName", c.PatchTable)
	tableEndpoints.Delete("/:tableName", c.DeleteTable)

	// ingestion job endpoints
	jobEndpoints := tenantEndpoints.Group("/:tenantId/jobs")
	jobEndpoints.Get("/", c.GetJobs)
	jobEndpoints.Post("/", c.CreateJob)
	jobEndpoints.Get("/:jobID", c.GetJob)
	jobEndpoints.Put("/:jobID", c.UpdateJob)
	jobEndpoints.Patch("/:jobID", c.PatchJob)
	jobEndpoints.Delete("/:jobID", c.DeleteJob)
}

func (c *Commander) StartServer(address string) error {
	var err error
	// register to consul
	log.Info().Msg("register and declare Commander to Consul")
	if err = c.consul.Start(c); err != nil {
		return err
	}
	if err = c.consul.Declare(c); err != nil {
		return err
	}

	if err = c.initializeGRPCClient(); err != nil {
		return err
	}

	// initialize api
	log.Info().Msg("Commander server is ready to handle requests")
	return c.app.Listen(address)
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
	c.storageGRPCConn, err = grpc.Dial(rpcAddr,
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
	// deregister to consul
	log.Info().Msg("deregister Commander to Consul")
	if err := c.consul.Stop(c); err != nil {
		return err
	}

	// closing gRPC storage service client
	log.Info().Msg("close gRCP client connection of Storage Service")
	if err := c.storageGRPCConn.Close(); err != nil {
		return err
	}

	log.Info().Msg("shutting down server...")
	return c.app.Shutdown()
}

func (c *Commander) GetHost() string {
	return c.Hostname
}

func (c *Commander) GetPort() string {
	return fmt.Sprint(c.config.Commander.Port)
}

func (c *Commander) GetName() string {
	return c.Name
}

func (c *Commander) GetID() string {
	return c.ID.String()
}

func (c *Commander) GetMetadata() map[string]string {
	return nil
}
