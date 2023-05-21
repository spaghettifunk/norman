package commander

import (
	"fmt"
	"os"

	"github.com/goccy/go-json"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/rs/zerolog/log"
	configuration "github.com/spaghettifunk/norman/internal/common"
	"github.com/spaghettifunk/norman/internal/common/manager"
	"github.com/spaghettifunk/norman/pkg/consul"
)

type Commander struct {
	Name                string
	ID                  uuid.UUID
	consul              *consul.Consul
	config              configuration.Configuration
	app                 *fiber.App
	schemaManager       *manager.SchemaManager
	ingestionJobManager *manager.IngestionJobManager
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

	// Create the new job manager
	ijb, err := manager.NewIngestionJobManager(cs)
	if err != nil {
		return nil, err
	}
	// initialize job manager
	ijb.Initialize()

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	c := &Commander{
		Name:                "commander",
		ID:                  id,
		consul:              cs,
		config:              config,
		app:                 app,
		schemaManager:       manager.NewSchemaManager(cs),
		ingestionJobManager: ijb,
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

	// schema endpoints
	schemaEndpoints := tenantEndpoints.Group("/:tenantId/schemas")
	schemaEndpoints.Get("/", c.GetSchemas)
	schemaEndpoints.Post("/", c.CreateSchema)
	schemaEndpoints.Get("/:schemaName", c.GetSchema)
	schemaEndpoints.Put("/:schemaName", c.UpdateSchema)
	schemaEndpoints.Patch("/:schemaName", c.PatchSchema)
	schemaEndpoints.Delete("/:schemaName", c.DeleteSchema)

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
	// register to consul
	log.Info().Msg("register and declare Commander to Consul")
	if err := c.consul.Start(c); err != nil {
		return err
	}
	if err := c.consul.Declare(c); err != nil {
		return err
	}

	// initialize api
	log.Info().Msg("Commander server is ready to handle requests")
	return c.app.Listen(address)
}

func (c *Commander) ShutdownServer() error {
	log.Info().Msg("shutting down schema manager...")
	if err := c.schemaManager.Shutdown(); err != nil {
		return err
	}
	log.Info().Msg("shutting down ingestion job manager...")
	if err := c.ingestionJobManager.Shutdown(); err != nil {
		return err
	}

	// deregister to consul
	log.Info().Msg("deregister Commander to Consul")
	if err := c.consul.Stop(c); err != nil {
		return err
	}

	log.Info().Msg("shutting down server...")
	return c.app.Shutdown()
}

func (c *Commander) GetHost() string {
	hn, err := os.Hostname()
	if err != nil {
		panic(err.Error())
	}
	return hn
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
