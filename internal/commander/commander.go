package commander

import (
	"github.com/goccy/go-json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/rs/zerolog/log"
	configuration "github.com/spaghettifunk/norman/internal/common"
	"github.com/spaghettifunk/norman/internal/common/manager"
)

type Commander struct {
	Name                string
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

	// Create the new job manager
	ijb, err := manager.NewIngestionJobManager()
	if err != nil {
		return nil, err
	}

	c := &Commander{
		Name:                "commander",
		config:              config,
		app:                 app,
		schemaManager:       manager.NewSchemaManager(),
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
	log.Info().Msg("Commander Server is ready to handle requests")
	return c.app.Listen(address)
}

func (c *Commander) ShutdownServer() error {
	log.Info().Msg("Shutting down server...")
	return c.app.Shutdown()
}
