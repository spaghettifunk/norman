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
	Name          string
	config        configuration.Configuration
	app           *fiber.App
	schemaManager *manager.SchemaManager
	tableManager  *manager.TableManager
}

func New(config configuration.Configuration) *Commander {
	// Create new Fiber application
	app := fiber.New(fiber.Config{
		AppName:           "commander-api-server",
		EnablePrintRoutes: true, // TODO: change this based on logger level -- DEBUG
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
	})
	// add default middleware
	app.Use(recover.New())

	c := &Commander{
		Name:          "commander",
		config:        config,
		app:           app,
		schemaManager: manager.NewSchemaManager(),
		tableManager:  manager.NewTableManager(),
	}

	c.setupRoutes()

	return c
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
	schemaEndpoints := tenantEndpoints.Group("/schemas")
	schemaEndpoints.Get("/", c.GetSchemas)
	schemaEndpoints.Post("/", c.CreateSchema)
	schemaEndpoints.Get("/:schemaName", c.GetSchema)
	schemaEndpoints.Put("/:schemaName", c.UpdateSchema)
	schemaEndpoints.Patch("/:schemaName", c.PatchSchema)
	schemaEndpoints.Delete("/:schemaName", c.DeleteSchema)
}

func (c *Commander) StartServer(address string) error {
	log.Info().Msg("Commander Server is ready to handle requests")
	return c.app.Listen(address)
}

func (c *Commander) ShutdownServer() error {
	log.Info().Msg("Shutting down server...")
	return c.app.Shutdown()
}
