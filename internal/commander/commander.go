package commander

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
	configuration "github.com/spaghettifunk/norman/internal/common"
)

type Commander struct {
	Name   string
	config configuration.Configuration
	app    *fiber.App
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
		Name:   "commander",
		config: config,
		app:    app,
	}

	c.setupRoutes()

	return c
}

func (c *Commander) setupRoutes() {
	apiV1 := c.app.Group("/commander/v1")

	apiV1.Get("/", c.APIVersion)
}

func (c *Commander) StartServer(address string) error {
	log.Info().Msg("Commander Server is ready to handle requests")
	return c.app.Listen(address)
}

func (c *Commander) ShutdownServer() error {
	log.Info().Msg("Shutting down server...")
	return c.app.Shutdown()
}
